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
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm/api/inventorypb"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services"
)

// BackupsService represents backups API.
type BackupsService struct {
	l              *logrus.Entry
	db             *reform.DB
	jobsService    jobsService
	versionService versionService
}

// NewBackupsService creates new backups API service.
func NewBackupsService(
	db *reform.DB,
	jobsService jobsService,
	versionService versionService,
) *BackupsService {
	return &BackupsService{
		l:              logrus.WithField("component", "management/backup"),
		db:             db,
		jobsService:    jobsService,
		versionService: versionService,
	}
}

type checkServiceResult struct {
	ServiceType models.ServiceType
	DBConfig    *models.DBConfig
	AgentID     string
}

func (s *BackupsService) checkService(ctx context.Context, serviceID string) (*checkServiceResult, error) {
	var result checkServiceResult
	if err := s.db.InTransactionContext(ctx, nil, func(tx *reform.TX) error {
		service, err := models.FindServiceByID(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		switch service.ServiceType {
		case models.MySQLServiceType,
			models.MongoDBServiceType:
			result.ServiceType = service.ServiceType
		case models.PostgreSQLServiceType,
			models.ProxySQLServiceType,
			models.HAProxyServiceType,
			models.ExternalServiceType:
			return status.Errorf(codes.Unimplemented, "unimplemented service: %s", service.ServiceType)
		default:
			return status.Errorf(codes.Unknown, "unknown service: %s", service.ServiceType)
		}

		dbConfig, err := models.FindDBConfigForService(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		result.DBConfig = dbConfig

		pmmAgents, err := models.FindPMMAgentsForService(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		if len(pmmAgents) == 0 {
			return errors.Errorf("pmmAgent not found for service")
		}

		result.AgentID = pmmAgents[0].AgentID

		return nil
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

// checkVersionCompatibility checks version compatibility between db and tool used for the backup.
// Returns version of the db.
func (s *BackupsService) checkVersionCompatibility(
	agentID string,
	serviceType models.ServiceType,
) (string, error) {
	// TODO: add support for mongodb, provide ticket
	if serviceType != models.MySQLServiceType {
		return "", nil
	}

	dbVersion, err := s.versionService.GetLocalMySQLVersion(agentID)
	if err != nil {
		return "", err
	}

	toolVersion, err := s.versionService.GetXtrabackupVersion(agentID)
	if err != nil {
		return "", err
	}

	_ = toolVersion
	// TODO: check compatibility between db and tool versions

	return dbVersion, nil
}

// StartBackup starts on-demand backup.
func (s *BackupsService) StartBackup(ctx context.Context, req *backupv1beta1.StartBackupRequest) (*backupv1beta1.StartBackupResponse, error) {
	csr, err := s.checkService(ctx, req.ServiceId)
	if err != nil {
		return nil, err
	}

	dbVersion, err := s.checkVersionCompatibility(csr.AgentID, csr.ServiceType)
	if err != nil {
		return nil, err
	}

	var artifact *models.Artifact
	var job *models.JobResult
	var locationConfig *models.BackupLocationConfig

	if err := s.db.InTransactionContext(ctx, nil, func(tx *reform.TX) error {
		if _, err := models.FindServiceByID(tx.Querier, req.ServiceId); err != nil {
			return err
		}

		location, err := models.FindBackupLocationByID(tx.Querier, req.LocationId)
		if err != nil {
			return err
		}

		locationConfig = &models.BackupLocationConfig{
			PMMServerConfig: location.PMMServerConfig,
			PMMClientConfig: location.PMMClientConfig,
			S3Config:        location.S3Config,
		}

		var dataModel models.DataModel
		var jobType models.JobType
		switch csr.ServiceType {
		case models.MySQLServiceType:
			dataModel = models.PhysicalDataModel
			jobType = models.MySQLBackupJob
		case models.MongoDBServiceType:
			dataModel = models.LogicalDataModel
			jobType = models.MongoDBBackupJob
		default:
			panic(fmt.Sprintf("unexpected ServiceType %v", csr.ServiceType))
		}

		artifact, err = models.CreateArtifact(tx.Querier, models.CreateArtifactParams{
			Name:       req.Name,
			Vendor:     string(csr.ServiceType),
			Version:    dbVersion,
			LocationID: location.ID,
			ServiceID:  req.ServiceId,
			DataModel:  dataModel,
			Status:     models.PendingBackupStatus,
		})
		if err != nil {
			return err
		}

		job, err = s.prepareBackupJob(tx.Querier, csr.AgentID, artifact.ID, jobType)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	switch csr.ServiceType {
	case models.MySQLServiceType:
		err = s.jobsService.StartMySQLBackupJob(job.ID, job.PMMAgentID, 0, req.Name, csr.DBConfig, locationConfig)
	case models.MongoDBServiceType:
		err = s.jobsService.StartMongoDBBackupJob(job.ID, job.PMMAgentID, 0, req.Name, csr.DBConfig, locationConfig)
	default:
		panic(fmt.Sprintf("unexpected ServiceType %v", csr.ServiceType))
	}
	if err != nil {
		return nil, err
	}

	return &backupv1beta1.StartBackupResponse{
		ArtifactId: artifact.ID,
	}, nil
}

func (s *BackupsService) findCompatibleServices(
	serviceType models.ServiceType,
	artifactDBVersion string,
) ([]*models.Service, error) {
	if artifactDBVersion == "" {
		s.l.Info("skip finding compatible services: artifact db version is unknown")
		return []*models.Service{}, nil
	}

	filters := models.ServiceFilters{
		ServiceType: &serviceType,
	}
	filteredServices, err := models.FindServices(s.db.Querier, filters)
	if err != nil {
		return nil, err
	}

	if len(filteredServices) == 0 {
		return []*models.Service{}, nil
	}

	const maxWorkers = 5
	workers := len(filteredServices)
	if workers > maxWorkers {
		workers = maxWorkers
	}

	g := &errgroup.Group{}
	// stream service for checking
	serviceCh := make(chan *models.Service, workers)
	g.Go(func() error {
		defer close(serviceCh)

		for _, s := range filteredServices {
			serviceCh <- s
		}

		return nil
	})

	// check each service for compatibility
	compatibleServicesMap := make(map[string]struct{}, len(filteredServices))
	m := &sync.Mutex{}
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for service := range serviceCh {
				pmmAgents, err := models.FindPMMAgentsForService(s.db.Querier, service.ServiceID)
				if err != nil {
					return err
				}

				if len(pmmAgents) == 0 {
					continue
				}

				agentID := pmmAgents[0].AgentID

				dbVersion, err := s.versionService.GetLocalMySQLVersion(agentID)
				if err != nil {
					s.l.WithError(err).Infof("skip incompatible service with id %q, agent id %q: "+
						"cannot get local MySQL version", service.ServiceID, agentID)
					continue
				}

				if artifactDBVersion != dbVersion {
					s.l.WithError(err).Infof("skip incompatible service with id %q, agent id %q: "+
						"version mismatch, artifact version %q != db version %q", service.ServiceID, agentID, artifactDBVersion, dbVersion)
					continue
				}

				toolVersion, err := s.versionService.GetXtrabackupVersion(agentID)
				if err != nil {
					s.l.WithError(err).Infof("skip incompatible service with id %q, agent id %q: "+
						"cannot get xtrabackup version", service.ServiceID, agentID)
					continue
				}

				_ = toolVersion
				// TODO: check compatibility between db and tool versions

				m.Lock()
				compatibleServicesMap[service.ServiceID] = struct{}{}
				m.Unlock()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	compatibleServices := make([]*models.Service, 0, len(compatibleServicesMap))
	for _, service := range filteredServices {
		if _, ok := compatibleServicesMap[service.ServiceID]; ok {
			compatibleServices = append(compatibleServices, service)
		}
	}

	return compatibleServices, nil
}

// ListServicesForRestore lists compatible service for restoring given artifact.
func (s *BackupsService) ListServicesForRestore(
	_ context.Context,
	req *backupv1beta1.ListServicesForRestoreRequest,
) (*backupv1beta1.ListServicesForRestoreResponse, error) {
	artifact, err := models.FindArtifactByID(s.db.Querier, req.ArtifactId)
	if err != nil {
		return nil, err
	}

	serviceType := models.ServiceType(artifact.Vendor)
	switch serviceType {
	case models.MySQLServiceType:
	case models.MongoDBServiceType, // TODO: add mongo support
		models.PostgreSQLServiceType,
		models.ProxySQLServiceType,
		models.HAProxyServiceType,
		models.ExternalServiceType:
		return nil, status.Errorf(codes.Unimplemented, "unimplemented service: %s", serviceType)
	default:
		return nil, status.Errorf(codes.Unknown, "unknown service: %s", serviceType)
	}

	compatibleServices, err := s.findCompatibleServices(serviceType, artifact.Version)
	if err != nil {
		return nil, err
	}

	res := &backupv1beta1.ListServicesForRestoreResponse{}
	for _, service := range compatibleServices {
		apiService, err := services.ToAPIService(service)
		if err != nil {
			return nil, err
		}

		switch s := apiService.(type) {
		case *inventorypb.MySQLService:
			res.Mysql = append(res.Mysql, s)
		default:
			panic(fmt.Errorf("unhandled inventory Service type %T", service))
		}
	}

	return res, nil
}

// RestoreBackup starts restore backup job.
func (s *BackupsService) RestoreBackup(
	ctx context.Context,
	req *backupv1beta1.RestoreBackupRequest,
) (*backupv1beta1.RestoreBackupResponse, error) {
	var params *prepareRestoreJobParams
	var jobID, restoreID string

	err := s.db.InTransactionContext(ctx, nil, func(tx *reform.TX) error {
		var err error
		params, err = s.prepareRestoreJob(tx.Querier, req.ServiceId, req.ArtifactId)
		if err != nil {
			return err
		}

		restore, err := models.CreateRestoreHistoryItem(tx.Querier, models.CreateRestoreHistoryItemParams{
			ArtifactID: req.ArtifactId,
			ServiceID:  req.ServiceId,
			Status:     models.InProgressRestoreStatus,
		})
		if err != nil {
			return err
		}

		restoreID = restore.ID

		service, err := models.FindServiceByID(tx.Querier, req.ServiceId)
		if err != nil {
			return err
		}

		var jobType models.JobType
		var jobResultData *models.JobResultData
		switch service.ServiceType {
		case models.MySQLServiceType:
			jobType = models.MySQLRestoreBackupJob
			jobResultData = &models.JobResultData{
				MySQLRestoreBackup: &models.MySQLRestoreBackupJobResult{
					RestoreID: restoreID,
				},
			}
		case models.MongoDBServiceType:
			jobType = models.MongoDBRestoreBackupJob
			jobResultData = &models.JobResultData{
				MongoDBRestoreBackup: &models.MongoDBRestoreBackupJobResult{
					RestoreID: restoreID,
				},
			}
		case models.PostgreSQLServiceType,
			models.ProxySQLServiceType,
			models.HAProxyServiceType,
			models.ExternalServiceType:
			return errors.Errorf("backup restore unimplemented for service type: %s", service.ServiceType)
		default:
			return errors.Errorf("unsupported service type: %s", service.ServiceType)
		}

		job, err := models.CreateJobResult(tx.Querier, params.AgentID, jobType, jobResultData)
		if err != nil {
			return err
		}

		jobID = job.ID

		return err
	})
	if err != nil {
		return nil, err
	}

	if err := s.startRestoreJob(jobID, req.ServiceId, params); err != nil {
		return nil, err
	}

	return &backupv1beta1.RestoreBackupResponse{
		RestoreId: restoreID,
	}, nil
}

func (s *BackupsService) prepareBackupJob(
	q *reform.Querier,
	agentID string,
	artifactID string,
	jobType models.JobType,
) (*models.JobResult, error) {
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
		models.MySQLRestoreBackupJob,
		models.MongoDBRestoreBackupJob:
		return nil, errors.Errorf("%s is not a backup job type", jobType)
	default:
		return nil, errors.Errorf("unsupported backup job type: %s", jobType)
	}

	res, err := models.CreateJobResult(q, agentID, jobType, jobResultData)

	if err != nil {
		return nil, err
	}

	return res, nil
}

type prepareRestoreJobParams struct {
	AgentID      string
	ArtifactName string
	Location     *models.BackupLocation
	ServiceType  models.ServiceType
	DBConfig     *models.DBConfig
}

func (s *BackupsService) prepareRestoreJob(
	q *reform.Querier,
	serviceID string,
	artifactID string,
) (*prepareRestoreJobParams, error) {
	service, err := models.FindServiceByID(q, serviceID)
	if err != nil {
		return nil, err
	}

	artifact, err := models.FindArtifactByID(q, artifactID)
	if err != nil {
		return nil, err
	}

	location, err := models.FindBackupLocationByID(q, artifact.LocationID)
	if err != nil {
		return nil, err
	}

	dbConfig, err := models.FindDBConfigForService(q, service.ServiceID)
	if err != nil {
		return nil, err
	}

	agents, err := models.FindPMMAgentsForService(q, serviceID)
	if err != nil {
		return nil, err
	}
	if len(agents) == 0 {
		return nil, errors.Errorf("cannot find pmm agent for service %s", serviceID)
	}

	return &prepareRestoreJobParams{
		AgentID:      agents[0].AgentID,
		ArtifactName: artifact.Name,
		Location:     location,
		ServiceType:  service.ServiceType,
		DBConfig:     dbConfig,
	}, nil
}

func (s *BackupsService) startRestoreJob(jobID, serviceID string, params *prepareRestoreJobParams) error {
	locationConfig := &models.BackupLocationConfig{
		PMMServerConfig: params.Location.PMMServerConfig,
		PMMClientConfig: params.Location.PMMClientConfig,
		S3Config:        params.Location.S3Config,
	}

	switch params.ServiceType {
	case models.MySQLServiceType:
		if err := s.jobsService.StartMySQLRestoreBackupJob(
			jobID,
			params.AgentID,
			serviceID,
			0,
			params.ArtifactName,
			locationConfig,
		); err != nil {
			return err
		}
	case models.MongoDBServiceType:
		if err := s.jobsService.StartMongoDBRestoreBackupJob(
			jobID,
			params.AgentID,
			0,
			params.ArtifactName,
			params.DBConfig,
			locationConfig,
		); err != nil {
			return err
		}
	case models.PostgreSQLServiceType,
		models.ProxySQLServiceType,
		models.HAProxyServiceType,
		models.ExternalServiceType:
		return status.Errorf(codes.Unimplemented, "unimplemented service: %s", params.ServiceType)
	default:
		return status.Errorf(codes.Unknown, "unknown service: %s", params.ServiceType)
	}

	return nil
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
