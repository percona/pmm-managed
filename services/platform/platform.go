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

// Package auth provides authentication/authorization functionality.
package platform

import (
	"context"
	"net"
	"os"
	"time"

	api "github.com/percona-platform/saas/gen/auth/external"
	"github.com/percona/pmm/utils/tlsconfig"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

const (
	defaultHost = "check.percona.com:443"
	envHost     = "PERCONA_TEST_AUTH_HOST"
	timeout     = 5 * time.Second
	authType    = "PP-v1beta1" // TODO Change to PP-1 after auth API release
)

// Service is responsible for interactions with Percona platform.
type Service struct {
	db   *reform.DB
	host string
	l    *logrus.Entry
}

// New returns platform Service.
func New(db *reform.DB) *Service {
	l := logrus.WithField("component", "auth")

	s := Service{
		host: defaultHost,
		db:   db,
		l:    l,
	}

	if h := os.Getenv(envHost); h != "" {
		l.Warnf("Host changed to %s.", h)
		s.host = h
	}

	return &s
}

// SignUp creates new Percona Platform user with given email and password.
func (s *Service) SignUp(ctx context.Context, email, password string) error {
	cc, err := s.getConnection(ctx)
	if err != nil {
		return errors.Wrap(err, "failed establish connection with Percona")
	}
	defer cc.Close() //nolint:errcheck

	_, err = api.NewAuthAPIClient(cc).SignUp(ctx, &api.SignUpRequest{Email: email, Password: password})
	if err != nil {
		return errors.Wrap(err, "failed to sign up user")
	}

	return nil
}

// SignIn checks Percona Platform user authentication and creates session.
func (s *Service) SignIn(ctx context.Context, email, password string) error {
	cc, err := s.getConnection(ctx)
	if err != nil {
		return errors.Wrap(err, "failed establish connection with Percona")
	}
	defer cc.Close() //nolint:errcheck

	resp, err := api.NewAuthAPIClient(cc).SignIn(ctx, &api.SignInRequest{Email: email, Password: password})
	if err != nil {
		return errors.Wrap(err, "failed to sign in user")
	}

	err = s.db.InTransaction(func(tx *reform.TX) error {
		_, err := models.UpdateSettings(tx.Querier, &models.ChangeSettingsParams{SessionID: resp.SessionId})
		return err
	})
	if err != nil {
		return errors.Wrap(err, "failed to save session id")
	}

	return nil
}

// RefreshSession resets session timeout.
func (s *Service) RefreshSession(ctx context.Context) error {
	settings, err := models.GetSettings(s.db)
	if err != nil {
		return err
	}

	if settings.SaaS.SessionID == "" {
		return errors.New("no active sessions")
	}

	cc, err := s.getConnection(ctx)
	if err != nil {
		return errors.Wrap(err, "failed establish connection with Percona")
	}
	defer cc.Close() //nolint:errcheck

	md := metadata.New(map[string]string{"authorization": authType + " " + settings.SaaS.SessionID})
	_, err = api.NewAuthAPIClient(cc).RefreshSession(metadata.NewOutgoingContext(ctx, md), &api.RefreshSessionRequest{})
	if err != nil {
		return errors.Wrap(err, "failed to refresh session")
	}

	return nil
}

func (s *Service) getConnection(ctx context.Context) (*grpc.ClientConn, error) {
	host, _, err := net.SplitHostPort(s.host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set checks host")
	}
	tlsConfig := tlsconfig.Get()
	tlsConfig.ServerName = host

	opts := []grpc.DialOption{
		// replacement is marked as experimental
		grpc.WithBackoffMaxDelay(timeout), //nolint:staticcheck

		grpc.WithBlock(),
		grpc.WithUserAgent("pmm-managed/" + version.Version),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cc, err := grpc.DialContext(ctx, s.host, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}

	return cc, nil
}
