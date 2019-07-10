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
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// AuthServer authenticates incoming requests via Grafana API.
type AuthServer struct {
	c *Client
	l *logrus.Entry
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
		s.l.Errorf("%+v", err)
		rw.WriteHeader(500)
		return
	}
}

func (s *AuthServer) authenticate(ctx context.Context, req *http.Request) error {
	// TODO l := logger.Get(ctx) once we have it after https://jira.percona.com/browse/PMM-4326

	if s.l.Logger.GetLevel() >= logrus.DebugLevel {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			s.l.Errorf("Failed to dump request: %v.", err)
		}
		s.l.Debugf("Request:\n%s", b)
	}

	if req.URL.Path != "/auth_request" {
		return errors.Errorf("Unexpected path %s.", req.URL.Path)
	}

	h := make(http.Header)
	for _, k := range []string{
		"Authorization",
		"Cookie",
	} {
		if v := req.Header.Get(k); v != "" {
			h.Set(k, v)
		}
	}

	if isGrafanaAdmin, _ := s.c.isGrafanaAdmin(h); isGrafanaAdmin {
		return nil
	}

	role, err := s.c.getRole(h)
	if err != nil {
		return err
	}

	s.l.Warnf("Unhandled request(%q, %s), authenticating anyway.", role, req.URL.Path)
	return nil
}
