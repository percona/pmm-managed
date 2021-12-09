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

// Package platform provides authentication/authorization functionality.
package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm/api/platformpb"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/envvars"
)

const (
	internalServerError = "Internal server error"
)

// supervisordService is a subset of methods of supervisord.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type supervisordService interface {
	UpdateConfiguration(settings *models.Settings, ssoDetails *models.PerconaSSODetails) error
}

// Service is responsible for interactions with Percona Platform.
type Service struct {
	db          *reform.DB
	host        string
	l           *logrus.Entry
	supervisord supervisordService
}

// New returns platform Service.
func New(db *reform.DB, supervisord supervisordService) (*Service, error) {
	l := logrus.WithField("component", "auth")

	host, err := envvars.GetSAASHost()
	if err != nil {
		return nil, err
	}

	s := Service{
		host:        host,
		db:          db,
		l:           l,
		supervisord: supervisord,
	}

	return &s, nil
}

const platformAPITimeout = 10 * time.Second

// Connect connects a PMM server to the organization created on Percona Portal. That allows the user to sign in to the PMM server with their Percona Account.
func (s *Service) Connect(ctx context.Context, req *platformpb.ConnectRequest) (*platformpb.ConnectResponse, error) {
	_, err := models.GetPerconaSSODetails(ctx, s.db.Querier)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "PMM server is already connected to Portal")
	}
	settings, err := models.GetSettings(s.db)
	if err != nil {
		s.l.Errorf("Failed to fetch PMM server ID and address: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}
	if settings.PMMPublicAddress == "" {
		return nil, status.Error(codes.FailedPrecondition, "The address of PMM server is not set")
	}
	pmmServerURL := fmt.Sprintf("https://%s/graph", settings.PMMPublicAddress)

	nCtx, cancel := context.WithTimeout(ctx, platformAPITimeout)
	defer cancel()

	ssoParams, err := s.connect(nCtx, &connectPMMParams{
		serverName:                req.ServerName,
		email:                     req.Email,
		password:                  req.Password,
		pmmServerURL:              pmmServerURL,
		pmmServerOAuthCallbackURL: fmt.Sprintf("%s/login/generic_oauth", pmmServerURL),
		pmmServerID:               settings.PMMServerID,
	})
	if err != nil {
		return nil, err // this is already a status error
	}

	err = models.InsertPerconaSSODetails(s.db.Querier, &models.PerconaSSODetailsInsert{
		ClientID:     ssoParams.ClientID,
		ClientSecret: ssoParams.ClientSecret,
		IssuerURL:    ssoParams.IssuerURL,
		Scope:        ssoParams.Scope,
	})
	if err != nil {
		s.l.Errorf("Failed to insert SSO details: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}

	if err := s.UpdateSupervisordConfigurations(ctx); err != nil {
		s.l.Errorf("Failed to update configuration of grafana after connecting PMM to Portal: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}
	return &platformpb.ConnectResponse{}, nil
}

// Disconnect disconnects a PMM server from the organization created on Percona Portal.
func (s *Service) Disconnect(ctx context.Context, req *platformpb.DisconnectRequest) (*platformpb.DisconnectResponse, error) {
	ssoDetails, err := models.GetPerconaSSODetails(ctx, s.db.Querier)
	if err == nil {
		return nil, status.Error(codes.Aborted, "PMM server is not connected to Portal")
	}

	settings, err := models.GetSettings(s.db)
	if err != nil {
		s.l.Errorf("Failed to fetch PMM server ID and address: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}

	if err := s.CleanSupervisordConfigurations(ctx); err != nil {
		s.l.Errorf("Failed to clean configuration of grafana during removing PMM from Portal: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}

	err = models.DeletePerconaSSODetails(s.db.Querier)
	if err != nil {
		s.l.Errorf("Failed to delete SSO details: %s", err)
		if err := s.UpdateSupervisordConfigurations(ctx); err != nil {
			s.l.Errorf("Failed to rollback: %s", err)
		}
		return nil, status.Error(codes.Internal, internalServerError)
	}

	nCtx, cancel := context.WithTimeout(ctx, platformAPITimeout)
	defer cancel()

	err = s.disconnect(nCtx, &disconnectPMMParams{
		PMMServerID: settings.PMMServerID,
		AccessToken: ssoDetails.AccessToken.AccessToken,
	})
	if err != nil {
		if err := s.UpdateSupervisordConfigurations(ctx); err != nil {
			s.l.Errorf("Failed to rollback: %s", err)
		}
		if err := models.InsertPerconaSSODetails(s.db.Querier, &models.PerconaSSODetailsInsert{
			ClientID:     ssoDetails.ClientID,
			ClientSecret: ssoDetails.ClientSecret,
			IssuerURL:    ssoDetails.IssuerURL,
			Scope:        ssoDetails.Scope,
		}); err != nil {
			s.l.Errorf("Failed to rollback: %s", err)
		}

		return nil, err // this is already a status error
	}

	return &platformpb.DisconnectResponse{}, nil
}

func (s *Service) CleanSupervisordConfigurations(ctx context.Context) error {
	settings, err := models.GetSettings(s.db)
	if err != nil {
		return errors.Wrap(err, "failed to get settings")
	}

	if err := s.supervisord.UpdateConfiguration(settings, nil); err != nil {
		return errors.Wrap(err, "failed to update supervisord configuration")
	}
	return nil
}

func (s *Service) UpdateSupervisordConfigurations(ctx context.Context) error {
	settings, err := models.GetSettings(s.db)
	if err != nil {
		return errors.Wrap(err, "failed to get settings")
	}
	ssoDetails, err := models.GetPerconaSSODetails(ctx, s.db.Querier)
	if err != nil {
		if !errors.Is(err, reform.ErrNoRows) {
			return errors.Wrap(err, "failed to get SSO details")
		}
	}
	if err := s.supervisord.UpdateConfiguration(settings, ssoDetails); err != nil {
		return errors.Wrap(err, "failed to update supervisord configuration")
	}
	return nil
}

type connectPMMParams struct {
	pmmServerURL, pmmServerOAuthCallbackURL, pmmServerID, serverName, email, password string
}

type connectPMMRequest struct {
	PMMServerID               string `json:"pmm_server_id"`
	PMMServerName             string `json:"pmm_server_name"`
	PMMServerURL              string `json:"pmm_server_url"`
	PMMServerOAuthCallbackURL string `json:"pmm_server_oauth_callback_url"`
}

type disconnectPMMParams struct {
	PMMServerID string
	AccessToken string
}

type ssoDetails struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	IssuerURL    string `json:"issuer_url"`
}

type connectPMMResponse struct {
	SSODetails *ssoDetails `json:"sso_details"`
}

type grpcGatewayError struct {
	Message string `json:"message"`
	Code    uint32 `json:"code"`
}

func (s *Service) connect(ctx context.Context, params *connectPMMParams) (*ssoDetails, error) {
	endpoint := fmt.Sprintf("https://%s/v1/orgs/inventory", s.host)
	marshaled, err := json.Marshal(connectPMMRequest{
		PMMServerID:               params.pmmServerID,
		PMMServerName:             params.serverName,
		PMMServerURL:              params.pmmServerURL,
		PMMServerOAuthCallbackURL: params.pmmServerOAuthCallbackURL,
	})
	if err != nil {
		s.l.Errorf("Failed to marshal request data: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}

	client := http.Client{Timeout: platformAPITimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(marshaled))
	if err != nil {
		s.l.Errorf("Failed to build Connect to Platform request: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}
	req.SetBasicAuth(params.email, params.password)
	resp, err := client.Do(req)
	if err != nil {
		s.l.Errorf("Connect to Platform request failed: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}
	defer resp.Body.Close() //nolint:errcheck

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var gwErr grpcGatewayError
		if err := decoder.Decode(&gwErr); err != nil {
			s.l.Errorf("Connect to Platform request failed and we failed to decode error message: %s", err)
			return nil, status.Error(codes.Internal, internalServerError)
		}
		return nil, status.Error(codes.Code(gwErr.Code), gwErr.Message)
	}

	var response connectPMMResponse
	if err := decoder.Decode(&response); err != nil {
		s.l.Errorf("Failed to decode response into SSO details: %s", err)
		return nil, status.Error(codes.Internal, internalServerError)
	}
	return response.SSODetails, nil
}

func (s *Service) disconnect(ctx context.Context, params *disconnectPMMParams) error {
	endpoint := fmt.Sprintf("https://%s/v1/orgs/inventory/%s", s.host, params.PMMServerID)
	client := http.Client{Timeout: platformAPITimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		s.l.Errorf("Failed to build Disconnect to Platform request: %s", err)
		return status.Error(codes.Internal, internalServerError)
	}

	h := req.Header
	h.Add("Authorization", fmt.Sprintf("Bearer %s", params.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		s.l.Errorf("Disconnect to Platform request failed: %s", err)
		return status.Error(codes.Internal, internalServerError)
	}
	defer resp.Body.Close() //nolint:errcheck

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var gwErr grpcGatewayError
		if err := decoder.Decode(&gwErr); err != nil {
			s.l.Errorf("Disconnect to Platform request failed and we failed to decode error message: %s", err)
			return status.Error(codes.Internal, internalServerError)
		}
		return status.Error(codes.Code(gwErr.Code), gwErr.Message)
	}

	return nil
}
