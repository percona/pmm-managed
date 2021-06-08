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

	"github.com/percona/pmm-managed/services/scheduler"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// BackupsService represents backups API.
type BackupsService struct {
	db                  *reform.DB
	backupsLogicService backupsLogicService
	jobsService         jobsService
	scheduleService     scheduleService
	l                   *logrus.Entry
}

// NewBackupsService creates new backups API service.
func NewBackupsService(db *reform.DB, jobsService jobsService, backupsLogicService backupsLogicService, scheduleService scheduleService) *BackupsService {
	return &BackupsService{
		l:                   logrus.WithField("component", "management/backup/backups"),
		db:                  db,
		jobsService:         jobsService,
		backupsLogicService: backupsLogicService,
		scheduleService:     scheduleService,
	}
}

// StartBackup starts on-demand backup.
func (s *BackupsService) StartBackup(ctx context.Context, req *backupv1beta1.StartBackupRequest) (*backupv1beta1.StartBackupResponse, error) {
	artifactID, err := s.backupsLogicService.PerformBackup(ctx, req.ServiceId, req.LocationId, req.Name)
	if err != nil {
		return nil, err
	}

	return &backupv1beta1.StartBackupResponse{
		ArtifactId: artifactID,
	}, nil
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

		job, err := models.CreateJobResult(tx.Querier, params.AgentID, models.MySQLRestoreBackupJob, &models.JobResultData{
			MySQLRestoreBackup: &models.MySQLRestoreBackupJobResult{
				RestoreID: restoreID,
			},
		})
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

type prepareRestoreJobParams struct {
	AgentID      string
	ArtifactName string
	Location     *models.BackupLocation
	ServiceType  models.ServiceType
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
	case models.PostgreSQLServiceType,
		models.MongoDBServiceType,
		models.ProxySQLServiceType,
		models.HAProxyServiceType,
		models.ExternalServiceType:
		return status.Errorf(codes.Unimplemented, "unimplemented service: %s", params.ServiceType)
	default:
		return status.Errorf(codes.Unknown, "unknown service: %s", params.ServiceType)
	}

	return nil
}

// ScheduleBackup add new backup task to scheduler.
func (s *BackupsService) ScheduleBackup(ctx context.Context, req *backupv1beta1.ScheduleBackupRequest) (*backupv1beta1.ScheduleBackupResponse, error) {
	var id string
	err := s.db.InTransaction(func(tx *reform.TX) error {
		svc, err := models.FindServiceByID(tx.Querier, req.ServiceId)
		if err != nil {
			return err
		}

		_, err = models.FindBackupLocationByID(tx.Querier, req.LocationId)
		if err != nil {
			return err
		}
		var task scheduler.Task
		switch svc.ServiceType {
		case models.MySQLServiceType:
			task = scheduler.NewMySQLBackupTask(s.backupsLogicService, req.ServiceId, req.LocationId, req.Name, req.Description)
		case models.MongoDBServiceType:
			task = scheduler.NewMongoBackupTask(s.backupsLogicService, req.ServiceId, req.LocationId, req.Name, req.Description)
		case models.PostgreSQLServiceType,
			models.ProxySQLServiceType,
			models.HAProxyServiceType,
			models.ExternalServiceType:
			return status.Errorf(codes.Unimplemented, "unimplemented service: %s", svc.ServiceType)
		default:
			return status.Errorf(codes.Unknown, "unknown service: %s", svc.ServiceType)

		}
		scheduledTask, err := s.scheduleService.Add(task, req.CronExpression, req.StartTime.AsTime(), uint(req.RetryTimes), req.RetryInterval.AsDuration())
		if err != nil {
			return err
		}

		id = scheduledTask.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &backupv1beta1.ScheduleBackupResponse{ScheduleBackupId: id}, nil
}

// ListScheduledBackups lists all tasks related to backup.
func (s *BackupsService) ListScheduledBackups(ctx context.Context, req *backupv1beta1.ListScheduledBackupsRequest) (*backupv1beta1.ListScheduledBackupsResponse, error) {
	tasks, err := models.FindScheduledTasks(s.db.Querier, models.ScheduledTasksFilter{
		Types: []models.ScheduledTaskType{
			models.ScheduledMySQLBackupTask,
			models.ScheduledMongoBackupTask,
		},
	})
	if err != nil {
		return nil, err
	}

	locationIDs := make([]string, 0, len(tasks))
	serviceIDs := make([]string, 0, len(tasks))
	for _, t := range tasks {
		var serviceID string
		var locationID string
		switch t.Type {
		case models.ScheduledMySQLBackupTask:
			serviceID = t.Data.MySQLBackupTask.ServiceID
			locationID = t.Data.MySQLBackupTask.LocationID
		case models.ScheduledMongoBackupTask:
			serviceID = t.Data.MongoBackupTask.ServiceID
			locationID = t.Data.MongoBackupTask.LocationID
		default:
			continue
		}
		serviceIDs = append(serviceIDs, serviceID)
		locationIDs = append(locationIDs, locationID)
	}
	locations, err := models.FindBackupLocationsByIDs(s.db.Querier, locationIDs)
	if err != nil {
		return nil, err
	}

	services, err := models.FindServicesByIDs(s.db.Querier, serviceIDs)
	if err != nil {
		return nil, err
	}

	scheduledBackups := make([]*backupv1beta1.ScheduledBackup, 0, len(tasks))
	for _, task := range tasks {
		backup, err := convertTaskToScheduledBackup(task, services, locations)
		if err != nil {
			s.l.WithError(err).Warnf("convert task to scheduled backup")
			continue
		}
		scheduledBackups = append(scheduledBackups, backup)
	}

	return &backupv1beta1.ListScheduledBackupsResponse{
		ScheduledBackups: scheduledBackups,
	}, nil

}

// ChangeScheduledBackup changes existing scheduled backup task.
func (s *BackupsService) ChangeScheduledBackup(ctx context.Context, req *backupv1beta1.ChangeScheduledBackupRequest) (*backupv1beta1.ChangeScheduledBackupResponse, error) {
	scheduledTask, err := models.FindScheduledTaskByID(s.db.Querier, req.ScheduleBackupId)
	if err != nil {
		return nil, err
	}
	switch scheduledTask.Type {
	case models.ScheduledMySQLBackupTask:
		data := scheduledTask.Data.MySQLBackupTask
		if req.Name != nil {
			data.Name = req.Name.Value
		}
		if req.Description != nil {
			data.Description = req.Description.Value
		}
	case models.ScheduledMongoBackupTask:
		data := scheduledTask.Data.MongoBackupTask
		if req.Name != nil {
			data.Name = req.Name.Value
		}
		if req.Description != nil {
			data.Description = req.Description.Value
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Unsupported type: %s", scheduledTask.Type)
	}

	params := models.ChangeScheduledTaskParams{
		Data: scheduledTask.Data,
	}

	if req.Enabled != nil {
		disabled := !req.Enabled.Value
		params.Disable = &disabled
	}

	if req.CronExpression != nil {
		val := req.CronExpression.Value
		params.CronExpression = &val
	}

	if req.RetryTimes != nil {
		val := uint(req.RetryTimes.Value)
		params.Retries = &val
	}

	if req.RetryInterval != nil {
		val := req.RetryInterval.AsDuration()
		params.RetryInterval = &val
	}

	err = s.db.InTransaction(func(tx *reform.TX) error {
		_, err := models.ChangeScheduledTask(s.db.Querier, req.ScheduleBackupId, params)
		if err != nil {
			return err
		}

		return s.scheduleService.Reload(req.ScheduleBackupId)
	})

	if err != nil {
		return nil, err
	}

	return &backupv1beta1.ChangeScheduledBackupResponse{}, nil
}

// RemoveScheduledBackup stops and removes existing scheduled backup task.
func (s *BackupsService) RemoveScheduledBackup(ctx context.Context, req *backupv1beta1.RemoveScheduledBackupRequest) (*backupv1beta1.RemoveScheduledBackupResponse, error) {
	err := s.scheduleService.Remove(req.ScheduleBackupId)
	if err != nil {
		return nil, err
	}
	return &backupv1beta1.RemoveScheduledBackupResponse{}, nil
}

func convertTaskToScheduledBackup(task *models.ScheduledTask,
	services map[string]*models.Service,
	locations map[string]*models.BackupLocation) (*backupv1beta1.ScheduledBackup, error) {
	backup := &backupv1beta1.ScheduledBackup{
		ScheduledBackupId: task.ID,
		CronExpression:    task.CronExpression,
		StartTime:         timestamppb.New(task.StartAt),
		RetryMode:         backupv1beta1.RetryMode_MANUAL,
		RetryInterval:     durationpb.New(task.RetryInterval),
		RetryTimes:        uint32(task.Retries),
		Enabled:           !task.Disabled,
		LastRun:           timestamppb.New(task.LastRun),
		NextRun:           timestamppb.New(task.NextRun)}
	if task.Retries > 0 {
		backup.RetryMode = backupv1beta1.RetryMode_AUTO
	}
	switch task.Type {
	case models.ScheduledMySQLBackupTask:
		data := task.Data.MySQLBackupTask
		backup.ServiceId = data.ServiceID
		backup.LocationId = data.LocationID
		backup.Name = data.Name
		backup.Description = data.Description
		backup.DataModel = backupv1beta1.DataModel_PHYSICAL
	case models.ScheduledMongoBackupTask:
		data := task.Data.MongoBackupTask
		backup.ServiceId = data.ServiceID
		backup.LocationId = data.LocationID
		backup.Name = data.Name
		backup.Description = data.Description
		backup.DataModel = backupv1beta1.DataModel_LOGICAL
	default:
		return nil, fmt.Errorf("unsupported task type: %s", task.Type)
	}

	backup.ServiceName = services[backup.ServiceId].ServiceName
	backup.Vendor = string(services[backup.ServiceId].ServiceType)
	backup.LocationName = locations[backup.LocationId].Name

	return backup, nil
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
