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
	"time"

	"github.com/AlekSi/pointer"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/backup"
	"github.com/percona/pmm-managed/services/scheduler"
)

// BackupsService represents backups API.
type BackupsService struct {
	db              *reform.DB
	backupService   backupService
	scheduleService scheduleService
	l               *logrus.Entry
}

// NewBackupsService creates new backups API service.
func NewBackupsService(db *reform.DB, backupService backupService, scheduleService scheduleService) *BackupsService {
	return &BackupsService{
		l:               logrus.WithField("component", "management/backup/backups"),
		db:              db,
		backupService:   backupService,
		scheduleService: scheduleService,
	}
}

// StartBackup starts on-demand backup.
func (s *BackupsService) StartBackup(ctx context.Context, req *backupv1beta1.StartBackupRequest) (*backupv1beta1.StartBackupResponse, error) {
	artifactID, err := s.backupService.PerformBackup(ctx, backup.PerformBackupParams{
		ServiceID:     req.ServiceId,
		LocationID:    req.LocationId,
		Name:          req.Name,
		Retries:       req.Retries,
		RetryInterval: req.RetryInterval.AsDuration(),
	})
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

	id, err := s.backupService.RestoreBackup(ctx, req.ServiceId, req.ArtifactId)
	if err != nil {
		return nil, err
	}

	return &backupv1beta1.RestoreBackupResponse{
		RestoreId: id,
	}, nil
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

		backupParams := scheduler.BackupTaskParams{
			ServiceID:     req.ServiceId,
			LocationID:    req.LocationId,
			Name:          req.Name,
			Description:   req.Description,
			Retention:     req.Retention,
			Retries:       req.Retries,
			RetryInterval: req.RetryInterval.AsDuration(),
		}

		var task scheduler.Task
		switch svc.ServiceType {
		case models.MySQLServiceType:
			task = scheduler.NewMySQLBackupTask(s.backupService, backupParams)
		case models.MongoDBServiceType:
			task = scheduler.NewMongoBackupTask(s.backupService, backupParams)
		case models.PostgreSQLServiceType,
			models.ProxySQLServiceType,
			models.HAProxyServiceType,
			models.ExternalServiceType:
			return status.Errorf(codes.Unimplemented, "unimplemented service: %s", svc.ServiceType)
		default:
			return status.Errorf(codes.Unknown, "unknown service: %s", svc.ServiceType)
		}

		t := req.StartTime.AsTime()
		if t.Unix() == 0 {
			t = time.Time{}
		}

		scheduledTask, err := s.scheduleService.Add(task, scheduler.AddParams{
			CronExpression: req.CronExpression,
			Disabled:       !req.Enabled,
			StartAt:        t,
		})
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Couldn't schedule backup: %v", err)
		}

		id = scheduledTask.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &backupv1beta1.ScheduleBackupResponse{ScheduledBackupId: id}, nil
}

// ListScheduledBackups lists all tasks related to backup.
func (s *BackupsService) ListScheduledBackups(ctx context.Context, req *backupv1beta1.ListScheduledBackupsRequest) (*backupv1beta1.ListScheduledBackupsResponse, error) {
	tasks, err := models.FindScheduledTasks(s.db.Querier, models.ScheduledTasksFilter{
		Types: []models.ScheduledTaskType{
			models.ScheduledMySQLBackupTask,
			models.ScheduledMongoDBBackupTask,
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
		case models.ScheduledMongoDBBackupTask:
			serviceID = t.Data.MongoDBBackupTask.ServiceID
			locationID = t.Data.MongoDBBackupTask.LocationID
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
		scheduledBackup, err := convertTaskToScheduledBackup(task, services, locations)
		if err != nil {
			s.l.WithError(err).Warnf("convert task to scheduledBackup")
			continue
		}
		scheduledBackups = append(scheduledBackups, scheduledBackup)
	}

	return &backupv1beta1.ListScheduledBackupsResponse{
		ScheduledBackups: scheduledBackups,
	}, nil

}

// ChangeScheduledBackup changes existing scheduled backup task.
func (s *BackupsService) ChangeScheduledBackup(ctx context.Context, req *backupv1beta1.ChangeScheduledBackupRequest) (*backupv1beta1.ChangeScheduledBackupResponse, error) {
	scheduledTask, err := models.FindScheduledTaskByID(s.db.Querier, req.ScheduledBackupId)
	if err != nil {
		return nil, err
	}

	var data *models.BackupTaskData
	switch scheduledTask.Type {
	case models.ScheduledMySQLBackupTask:
		data = &scheduledTask.Data.MySQLBackupTask.BackupTaskData
	case models.ScheduledMongoDBBackupTask:
		data = &scheduledTask.Data.MongoDBBackupTask.BackupTaskData
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Unknown type: %s", scheduledTask.Type)
	}

	if req.Name != nil {
		data.Name = req.Name.Value
	}
	if req.Description != nil {
		data.Description = req.Description.Value
	}
	if req.Retention != nil {
		data.Retention = req.Retention.Value
	}
	if req.Retries != nil {
		data.Retries = req.Retries.Value
	}
	if req.RetryInterval != nil {
		data.RetryInterval = req.RetryInterval.AsDuration()
	}

	params := models.ChangeScheduledTaskParams{
		Data: scheduledTask.Data,
	}

	if req.Enabled != nil {
		params.Disable = pointer.ToBool(!req.Enabled.Value)
	}

	if req.CronExpression != nil {
		params.CronExpression = pointer.ToString(req.CronExpression.Value)
	}

	if err := s.scheduleService.Update(req.ScheduledBackupId, params); err != nil {
		return nil, err
	}

	return &backupv1beta1.ChangeScheduledBackupResponse{}, nil
}

// RemoveScheduledBackup stops and removes existing scheduled backup task.
func (s *BackupsService) RemoveScheduledBackup(ctx context.Context, req *backupv1beta1.RemoveScheduledBackupRequest) (*backupv1beta1.RemoveScheduledBackupResponse, error) {
	task, err := models.FindScheduledTaskByID(s.db.Querier, req.ScheduledBackupId)
	if err != nil {
		return nil, err
	}
	switch task.Type {
	case models.ScheduledMySQLBackupTask:
	case models.ScheduledMongoDBBackupTask:
	default:
		return nil, errors.Errorf("non-backup task: %s", task.Type)
	}

	errTx := s.db.InTransaction(func(tx *reform.TX) error {
		artifacts, err := models.FindArtifacts(tx.Querier, models.ArtifactFilters{
			ScheduleID: req.ScheduledBackupId,
		})
		if err != nil {
			return err
		}

		for _, artifact := range artifacts {
			_, err := models.UpdateArtifact(tx.Querier, artifact.ID, models.UpdateArtifactParams{
				ScheduleID: pointer.ToString(""),
			})
			if err != nil {
				return err
			}
		}

		return s.scheduleService.Remove(req.ScheduledBackupId)
	})

	if errTx != nil {
		return nil, errTx
	}

	return &backupv1beta1.RemoveScheduledBackupResponse{}, nil
}

func convertTaskToScheduledBackup(task *models.ScheduledTask,
	services map[string]*models.Service,
	locations map[string]*models.BackupLocation) (*backupv1beta1.ScheduledBackup, error) {
	scheduledBackup := &backupv1beta1.ScheduledBackup{
		ScheduledBackupId: task.ID,
		CronExpression:    task.CronExpression,
		Enabled:           !task.Disabled,
	}

	if !task.LastRun.IsZero() {
		scheduledBackup.LastRun = timestamppb.New(task.LastRun)
	}

	if !task.NextRun.IsZero() {
		scheduledBackup.NextRun = timestamppb.New(task.NextRun)
	}

	if !task.StartAt.IsZero() {
		scheduledBackup.StartTime = timestamppb.New(task.StartAt)
	}

	var commonBackupData models.BackupTaskData
	switch task.Type {
	case models.ScheduledMySQLBackupTask:
		commonBackupData = task.Data.MySQLBackupTask.BackupTaskData
		scheduledBackup.DataModel = backupv1beta1.DataModel_PHYSICAL

	case models.ScheduledMongoDBBackupTask:
		commonBackupData = task.Data.MongoDBBackupTask.BackupTaskData
		scheduledBackup.DataModel = backupv1beta1.DataModel_LOGICAL
	default:
		return nil, fmt.Errorf("unknown task type: %s", task.Type)
	}

	scheduledBackup.ServiceId = commonBackupData.ServiceID
	scheduledBackup.LocationId = commonBackupData.LocationID
	scheduledBackup.Name = commonBackupData.Name
	scheduledBackup.Description = commonBackupData.Description
	scheduledBackup.Retention = commonBackupData.Retention
	scheduledBackup.Retries = commonBackupData.Retries

	if commonBackupData.RetryInterval > 0 {
		scheduledBackup.RetryInterval = durationpb.New(commonBackupData.RetryInterval)
	}

	scheduledBackup.ServiceName = services[scheduledBackup.ServiceId].ServiceName
	scheduledBackup.Vendor = string(services[scheduledBackup.ServiceId].ServiceType)
	scheduledBackup.LocationName = locations[scheduledBackup.LocationId].Name

	return scheduledBackup, nil
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
