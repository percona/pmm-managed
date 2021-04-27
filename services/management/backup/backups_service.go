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
	"time"

	"github.com/google/uuid"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

//go:generate mockery -name=jobsService -case=snake -inpkg -testonly

// jobsService is a subset of methods of agents.JobsService used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type jobsService interface {
	StopJob(jobID string) error
	StartMySQLBackupRestoreJob(
		jobID string,
		pmmAgentID string,
		serviceID string,
		timeout time.Duration,
		name string,
		locationConfig models.BackupLocationConfig,
	) error
}

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

// RestoreBackup starts restore backup job.
func (s *BackupsService) RestoreBackup(
	ctx context.Context,
	req *backupv1beta1.RestoreBackupRequest,
) (*backupv1beta1.RestoreBackupResponse, error) {
	var restoreID string

	var artifact *models.Artifact
	var location *models.BackupLocation
	var service *models.Service
	var job *models.JobResult

	err := s.db.InTransactionContext(ctx, nil, func(tx *reform.TX) error {
		var err error
		q := tx.Querier

		service, err = models.FindServiceByID(q, req.ServiceId)
		if err != nil {
			return err
		}

		artifact, err = models.FindArtifactByID(q, req.ArtifactId)
		if err != nil {
			return err
		}

		location, err = models.FindBackupLocationByID(q, artifact.LocationID)
		if err != nil {
			return err
		}

		agents, err := models.FindPMMAgentsForService(tx.Querier, req.ServiceId)
		if err != nil {
			return err
		}

		if len(agents) == 0 {
			return errors.Errorf("cannot find pmm agent for service %s", req.ServiceId)
		}

		// TODO: replace with restore id created in restores table
		restoreID = "/restore_id/" + uuid.New().String()

		job, err = models.CreateJobResult(tx.Querier, agents[0].AgentID, models.MySQLBackupRestoreJob, &models.JobResultData{
			MySQLBackupRestore: &models.MySQLBackupRestoreJobResult{
				RestoreID: restoreID,
			},
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	locationConfig := models.BackupLocationConfig{
		PMMServerConfig: location.PMMServerConfig,
		PMMClientConfig: location.PMMClientConfig,
		S3Config:        location.S3Config,
	}

	switch service.ServiceType {
	case models.MySQLServiceType:
		if err := s.jobsService.StartMySQLBackupRestoreJob(
			job.ID,
			job.PMMAgentID,
			service.ServiceID,
			0,
			artifact.Name,
			locationConfig,
		); err != nil {
			return nil, err
		}
	case models.PostgreSQLServiceType,
		models.MongoDBServiceType,
		models.ProxySQLServiceType,
		models.HAProxyServiceType,
		models.ExternalServiceType:
		return nil, status.Errorf(codes.Unimplemented, "unsupported service: %s", service.ServiceType)
	default:
		return nil, status.Errorf(codes.Internal, "unexpected service: %s", service.ServiceType)
	}

	return &backupv1beta1.RestoreBackupResponse{
		RestoreId: restoreID,
	}, nil
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
