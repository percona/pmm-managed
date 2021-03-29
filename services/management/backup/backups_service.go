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

package backup

import (
	"context"

	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// BackupsService represents backups API.
type BackupsService struct {
	db          *reform.DB
	jobsService jobsService
}

// NewBackupsService creates new backups API service.
func NewBackupsService(db *reform.DB, jobsService jobsService) *BackupsService {
	return &BackupsService{
		db:          db,
		jobsService: jobsService,
	}
}

// PerformBackups starts on-demand backup.
func (s *BackupsService) StartBackup(ctx context.Context, req *backupv1beta1.StartBackupRequest) (*backupv1beta1.StartBackupResponse, error) {
	svc, err := models.FindServiceByID(s.db.Querier, req.ServiceId)
	if err != nil {
		return nil, err
	}
	pmmAgents, err := models.FindPMMAgentsForService(s.db.Querier, req.ServiceId)
	if err != nil {
		return nil, err
	}
	if len(pmmAgents) == 0 {
		return nil, status.Errorf(codes.NotFound, "No pmm-agent running on service %s", req.ServiceId)
	}
	backupLocation, err := models.FindBackupLocationByID(s.db.Querier, req.LocationId)
	if err != nil {
		return nil, err
	}
	id, err := s.startBackup(ctx, Options{
		Service:        svc,
		Agent:          pmmAgents[0],
		BackupLocation: backupLocation,
	})
	return &backupv1beta1.StartBackupResponse{
		BackupId: id,
	}, nil
}

type Options struct {
	Service        *models.Service
	Agent          *models.Agent
	BackupLocation *models.BackupLocation
}

func (s *BackupsService) startBackup(ctx context.Context, options Options) (string, error) {
	id := "/jobs/backup/mysql/dummy"
	err := s.jobsService.StartMySQLBackupJob(id, *options.Agent.PMMAgentID, 0, *options.Service.Address, models.BackupLocationConfig{
		PMMServerConfig: options.BackupLocation.PMMServerConfig,
		PMMClientConfig: options.BackupLocation.PMMClientConfig,
		S3Config:        options.BackupLocation.S3Config,
	})
	if err != nil {
		return "", err
	}
	return id, err
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
