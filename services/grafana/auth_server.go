// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"path"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

// rules maps original URL prefix to minimal required role.
var rules = map[string]role{
	"/agent.Agent/Connect": none,

	"/inventory.Agents/Get":    admin,
	"/inventory.Agents/List":   admin,
	"/inventory.Nodes/Get":     admin,
	"/inventory.Nodes/List":    admin,
	"/inventory.Services/Get":  admin,
	"/inventory.Services/List": admin,
	"/inventory.":              admin,

	"/management.": admin,

	"/server.": admin,

	"/v0/inventory/Agents/Get":    admin,
	"/v0/inventory/Agents/List":   admin,
	"/v0/inventory/Nodes/Get":     admin,
	"/v0/inventory/Nodes/List":    admin,
	"/v0/inventory/Services/Get":  admin,
	"/v0/inventory/Services/List": admin,
	"/v0/inventory/":              admin,

	"/v0/management/": admin,

	"/v1/Updates/Check":   admin,
	"/v1/Updates/Perform": admin,

	"/v1/Settings/Change": admin,
	"/v1/Settings/Get":    admin,

	"/v0/qan/": viewer,

	"/qan/":        viewer,
	"/prometheus/": viewer,

	// TODO cleanup
	"/v1/readyz": none,
	"/ping":      none, // PMM 1.x variant

	"/v1/version":         viewer,
	"/managed/v1/version": viewer, // PMM 1.x variant

	// "/" is a special case
}

// clientError contains authentication error response details.
type authError struct {
	code    codes.Code
	message string
}

// AuthServer authenticates incoming requests via Grafana API.
type AuthServer struct {
	c *Client
	l *logrus.Entry

	// TODO server metrics should be provided by middleware https://jira.percona.com/browse/PMM-4326
}

// NewAuthServer creates new AuthServer.
func NewAuthServer(c *Client) *AuthServer {
	return &AuthServer{
		c: c,
		l: logrus.WithField("component", "grafana/auth"),
	}
}

// ServeHTTP serves internal location /auth_request.
func (s *AuthServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// fail-safe
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	if err := s.authenticate(ctx, req); err != nil {
		// nginx completely ignores auth_request subrequest response body;
		// out nginx configuration then sends the same request as a normal request
		// and returns response body to the client

		// copy grpc-gateway behavior: set correct codes, set both "error" and "message"
		m := map[string]interface{}{
			"code":    int(err.code),
			"error":   err.message,
			"message": err.message,
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(runtime.HTTPStatusFromCode(err.code))
		if err := json.NewEncoder(rw).Encode(m); err != nil {
			s.l.Warnf("%s", err)
		}
	}
}

func (s *AuthServer) authenticate(ctx context.Context, req *http.Request) *authError {
	// TODO l := logger.Get(ctx) once we have it after https://jira.percona.com/browse/PMM-4326
	l := s.l

	if l.Logger.GetLevel() >= logrus.DebugLevel {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			l.Errorf("Failed to dump request: %v.", err)
		}
		l.Debugf("Request:\n%s", b)
	}

	if req.URL.Path != "/auth_request" {
		l.Errorf("Unexpected path %s.", req.URL.Path)
		return &authError{code: codes.Internal, message: "Internal server error."}
	}

	origURI := req.Header.Get("X-Original-Uri")
	if origURI == "" {
		l.Errorf("Empty X-Original-Uri.")
		return &authError{code: codes.Internal, message: "Internal server error."}
	}
	l = l.WithField("req", fmt.Sprintf("%s %s", req.Header.Get("X-Original-Method"), origURI))

	// find the longest prefix present in rules:
	// /foo/bar -> /foo/ -> /foo -> /
	prefix := origURI
	for prefix != "/" {
		if _, ok := rules[prefix]; ok {
			break
		}

		if strings.HasSuffix(prefix, "/") {
			prefix = strings.TrimSuffix(prefix, "/")
		} else {
			prefix = path.Dir(prefix) + "/"
		}
	}

	// fallback to Grafana admin if there is no explicit rule
	// TODO https://jira.percona.com/browse/PMM-4338
	minRole, ok := rules[prefix]
	if ok {
		l = l.WithField("prefix", prefix)
	} else {
		l.Warnf("No explicit rule for %q, falling back to Grafana admin.", origURI)
		minRole = grafanaAdmin
	}

	if minRole == none {
		l.Debugf("Minimal required role is %q, granting access without checking Grafana.", minRole)
		return nil
	}

	// check Grafana with some headers from request
	authHeaders := make(http.Header)
	for _, k := range []string{
		"Authorization",
		"Cookie",
	} {
		if v := req.Header.Get(k); v != "" {
			authHeaders.Set(k, v)
		}
	}
	role, err := s.c.getRole(ctx, authHeaders)
	if err != nil {
		l.Warnf("%s", err)
		if cErr, ok := errors.Cause(err).(*clientError); ok {
			code := codes.Internal
			if cErr.Code == 401 || cErr.Code == 403 {
				code = codes.Unauthenticated
			}
			return &authError{code: code, message: cErr.ErrorMessage}
		}
		return &authError{code: codes.Internal, message: "Internal server error."}
	}
	l = l.WithField("role", role.String())

	if role == grafanaAdmin {
		l.Debugf("Grafana admin, allowing access.")
		return nil
	}

	if minRole <= role {
		l.Debugf("Minimal required role is %q, granting access.", minRole)
		return nil
	}

	l.Warnf("Minimal required role is %q.", minRole)
	return &authError{code: codes.PermissionDenied, message: "Access denied."}
}
