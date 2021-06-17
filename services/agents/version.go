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

package agents

import (
	"github.com/pkg/errors"

	"github.com/percona/pmm/api/agentpb"

	"github.com/percona/pmm-managed/models"
)

// VersionService provides methods for retrieving versions of different software.
type VersionService struct {
	r *Registry
}

// NewVersionService returns new version service.
func NewVersionService(registry *Registry) *VersionService {
	return &VersionService{
		r: registry,
	}
}

// GetRemoteMySQLVersion retrieves remote MySQL version using provided credentials.
func (s *VersionService) GetRemoteMySQLVersion(agentID string, dbConfig *models.DBConfig) (string, error) {
	agent, err := s.r.get(agentID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	req := &agentpb.GetVersionRequest{
		Software: &agentpb.GetVersionRequest_RemoteMysql{
			RemoteMysql: &agentpb.GetVersionRequest_RemoteMySQL{
				User:     dbConfig.User,
				Password: dbConfig.Password,
				Address:  dbConfig.Address,
				Port:     int32(dbConfig.Port),
				Socket:   dbConfig.Socket,
			},
		},
	}
	resp, err := agent.channel.SendAndWaitResponse(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.(*agentpb.GetVersionResponse).Version, nil
}

// GetLocalMySQLVersion retrieves local MySQL version.
func (s *VersionService) GetLocalMySQLVersion(agentID string) (string, error) {
	agent, err := s.r.get(agentID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	req := &agentpb.GetVersionRequest{
		Software: &agentpb.GetVersionRequest_LocalMysql{
			LocalMysql: &agentpb.GetVersionRequest_LocalMySQL{},
		},
	}
	resp, err := agent.channel.SendAndWaitResponse(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.(*agentpb.GetVersionResponse).Version, nil
}

// GetXtrabackupVersion retrieves xtrabackup version.
func (s *VersionService) GetXtrabackupVersion(agentID string) (string, error) {
	agent, err := s.r.get(agentID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	req := &agentpb.GetVersionRequest{
		Software: &agentpb.GetVersionRequest_Xtrabackup{
			Xtrabackup: &agentpb.GetVersionRequest_XtraBackup{},
		},
	}
	resp, err := agent.channel.SendAndWaitResponse(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.(*agentpb.GetVersionResponse).Version, nil
}
