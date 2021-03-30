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

	"github.com/google/uuid"
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

	id, err := s.startBackup(ctx, svc, pmmAgents[0], backupLocation)
	if err != nil {
		return nil, err
	}

	return &backupv1beta1.StartBackupResponse{
		BackupId: id,
	}, nil
}

func (s *BackupsService) startBackup(ctx context.Context, svc *models.Service, agent *models.Agent, location *models.BackupLocation) (string, error) {
	var err error
	locationConfig := models.BackupLocationConfig{
		PMMServerConfig: location.PMMServerConfig,
		PMMClientConfig: location.PMMClientConfig,
		S3Config:        location.S3Config,
	}
	id := "/job_id/" + uuid.New().String()
	switch svc.ServiceType {
	case models.MySQLServiceType:
		err = s.jobsService.StartMySQLBackupJob(id, *agent.PMMAgentID, 0, *svc.Address, locationConfig)
	default:
		return "", status.Errorf(codes.Unimplemented, "unsupported service: %s", svc.ServiceType)
	}
	if err != nil {
		return "", err
	}
	return id, err
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
