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

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// BackupsLogicService represents core logic for backuping.
type BackupsLogicService struct {
	db          *reform.DB
	jobsService jobsService
	l           *logrus.Entry
}

// NewBackupsLogicService creates new backups logic service.
func NewBackupsLogicService(db *reform.DB, jobsService jobsService) *BackupsLogicService {
	return &BackupsLogicService{
		l:           logrus.WithField("component", "management/backup/backups-logic"),
		db:          db,
		jobsService: jobsService,
	}
}

// PerformBackup starts on-demand backup.
func (s *BackupsLogicService) PerformBackup(ctx context.Context, serviceID, locationID, name string) (string, error) {
	var err error
	var artifact *models.Artifact
	var location *models.BackupLocation
	var svc *models.Service
	var job *models.JobResult
	var config *models.DBConfig

	errTX := s.db.InTransactionContext(ctx, nil, func(tx *reform.TX) error {
		svc, err = models.FindServiceByID(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		location, err = models.FindBackupLocationByID(tx.Querier, locationID)
		if err != nil {
			return err
		}

		var dataModel models.DataModel
		var jobType models.JobType
		switch svc.ServiceType {
		case models.MySQLServiceType:
			dataModel = models.PhysicalDataModel
			jobType = models.MySQLBackupJob
		case models.MongoDBServiceType:
			dataModel = models.LogicalDataModel
			jobType = models.MongoDBBackupJob
		case models.PostgreSQLServiceType:
			fallthrough
		case models.ProxySQLServiceType:
			fallthrough
		case models.HAProxyServiceType:
			fallthrough
		case models.ExternalServiceType:
			return status.Errorf(codes.Unimplemented, "unimplemented service: %s", svc.ServiceType)
		default:
			return status.Errorf(codes.Unknown, "unknown service: %s", svc.ServiceType)
		}

		artifact, err = models.CreateArtifact(tx.Querier, models.CreateArtifactParams{
			Name:       name,
			Vendor:     string(svc.ServiceType),
			LocationID: location.ID,
			ServiceID:  svc.ServiceID,
			DataModel:  dataModel,
			Status:     models.PendingBackupStatus,
		})
		if err != nil {
			return err
		}

		job, config, err = s.prepareBackupJob(tx.Querier, svc, artifact.ID, jobType)
		if err != nil {
			return err
		}
		return nil
	})
	if errTX != nil {
		return "", errTX
	}

	locationConfig := &models.BackupLocationConfig{
		PMMServerConfig: location.PMMServerConfig,
		PMMClientConfig: location.PMMClientConfig,
		S3Config:        location.S3Config,
	}

	switch svc.ServiceType {
	case models.MySQLServiceType:
		err = s.jobsService.StartMySQLBackupJob(job.ID, job.PMMAgentID, 0, name, config, locationConfig)
	case models.MongoDBServiceType:
		err = s.jobsService.StartMongoDBBackupJob(job.ID, job.PMMAgentID, 0, name, config, locationConfig)
	case models.PostgreSQLServiceType,
		models.ProxySQLServiceType,
		models.HAProxyServiceType,
		models.ExternalServiceType:
		return "", status.Errorf(codes.Unimplemented, "unimplemented service: %s", svc.ServiceType)
	default:
		return "", status.Errorf(codes.Unknown, "unknown service: %s", svc.ServiceType)
	}
	if err != nil {
		return "", err
	}

	return artifact.ID, nil
}

func (s *BackupsLogicService) prepareBackupJob(
	q *reform.Querier,
	service *models.Service,
	artifactID string,
	jobType models.JobType,
) (*models.JobResult, *models.DBConfig, error) {
	dbConfig, err := models.FindDBConfigForService(q, service.ServiceID)
	if err != nil {
		return nil, nil, err
	}

	pmmAgents, err := models.FindPMMAgentsForService(q, service.ServiceID)
	if err != nil {
		return nil, nil, err
	}

	if len(pmmAgents) == 0 {
		return nil, nil, errors.Errorf("pmmAgent not found for service")
	}

	var jobResultData *models.JobResultData
	switch jobType {
	case models.MySQLBackupJob:
		jobResultData = &models.JobResultData{
			MySQLBackup: &models.MySQLBackupJobResult{
				ArtifactID: artifactID,
			},
		}
	case models.MongoDBBackupJob:
		jobResultData = &models.JobResultData{
			MongoDBBackup: &models.MongoDBBackupJobResult{
				ArtifactID: artifactID,
			},
		}
	case models.Echo,
		models.MySQLRestoreBackupJob:
		return nil, nil, errors.Errorf("%s is not a backup job type", jobType)
	default:
		return nil, nil, errors.Errorf("unsupported backup job type: %s", jobType)
	}

	res, err := models.CreateJobResult(q, pmmAgents[0].AgentID, jobType, jobResultData)

	if err != nil {
		return nil, nil, err
	}

	return res, dbConfig, nil
}
