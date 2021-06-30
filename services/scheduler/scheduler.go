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

package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/percona/pmm-managed/models"

	"github.com/AlekSi/pointer"
	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

// Service is responsible for executing tasks and storing them to DB.
type Service struct {
	db                  *reform.DB
	l                   *logrus.Entry
	backupsLogicService backupsLogicService

	mx        sync.Mutex
	scheduler *gocron.Scheduler

	taskMx sync.RWMutex
	tasks  map[string]context.CancelFunc

	jobsMx sync.RWMutex
	jobs   map[string]*gocron.Job
}

// New creates new scheduler service.
func New(db *reform.DB, backupsLogicService backupsLogicService) *Service {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()
	scheduler.WaitForScheduleAll()
	return &Service{
		db:                  db,
		scheduler:           scheduler,
		l:                   logrus.WithField("component", "scheduler"),
		backupsLogicService: backupsLogicService,
		tasks:               make(map[string]context.CancelFunc),
		jobs:                make(map[string]*gocron.Job),
	}
}

// Run loads tasks from DB and starts scheduler.
func (s *Service) Run(ctx context.Context) {
	if err := s.loadFromDB(); err != nil {
		s.l.Warn(err)
	}
	s.scheduler.StartAsync()
	<-ctx.Done()
	s.scheduler.Stop()
}

// AddParams contains parameters for adding new add to service.
type AddParams struct {
	CronExpression string
	Disabled       bool
	StartAt        time.Time
}

// Add adds task to scheduler and save it to DB.
func (s *Service) Add(task Task, params AddParams) (*models.ScheduledTask, error) {
	var scheduledTask *models.ScheduledTask
	var err error

	err = s.db.InTransaction(func(tx *reform.TX) error {
		scheduledTask, err = models.CreateScheduledTask(tx.Querier, models.CreateScheduledTaskParams{
			CronExpression: params.CronExpression,
			StartAt:        params.StartAt,
			Type:           task.Type(),
			Data:           task.Data(),
			Disabled:       params.Disabled,
		})
		if err != nil {
			return err
		}

		id := scheduledTask.ID
		task.SetID(id)

		// Don't add job to scheduler if task is disabled.
		if scheduledTask.Disabled {
			return nil
		}

		fn := s.wrapTask(task, id)

		s.mx.Lock()
		j := s.scheduler.Cron(params.CronExpression).SingletonMode()
		if !params.StartAt.IsZero() {
			j = j.StartAt(params.StartAt)
		}
		scheduleJob, err := j.Tag(id).Do(fn)
		s.mx.Unlock()
		if err != nil {
			return err
		}

		s.jobsMx.Lock()
		s.jobs[id] = scheduleJob
		s.jobsMx.Unlock()

		scheduledTask, err = models.ChangeScheduledTask(tx.Querier, id, models.ChangeScheduledTaskParams{
			NextRun: pointer.ToTime(scheduleJob.NextRun().UTC()),
			LastRun: pointer.ToTime(scheduleJob.LastRun().UTC()),
		})
		if err != nil {
			s.l.WithField("id", id).Errorf("failed to set next run for new created task")
			s.mx.Lock()
			s.scheduler.RemoveByReference(scheduleJob)
			s.mx.Unlock()
			return err
		}

		return nil
	})
	return scheduledTask, err
}

// Remove stops task specified by id and removes it from DB and scheduler.
func (s *Service) Remove(id string) error {
	s.taskMx.RLock()
	if cancel, ok := s.tasks[id]; ok {
		cancel()
	}
	s.taskMx.RUnlock()

	s.jobsMx.Lock()
	delete(s.jobs, id)
	s.jobsMx.Unlock()

	err := s.db.InTransaction(func(tx *reform.TX) error {
		if err := models.RemoveScheduledTask(tx.Querier, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.mx.Lock()
	_ = s.scheduler.RemoveByTag(id)
	s.mx.Unlock()

	return nil
}

// Update changes scheduled task in DB and re-add it to scheduler.
func (s *Service) Update(id string, params models.ChangeScheduledTaskParams) error {
	txErr := s.db.InTransaction(func(tx *reform.TX) error {
		dbTask, err := models.ChangeScheduledTask(tx.Querier, id, params)
		if err != nil {
			return err
		}
		s.mx.Lock()
		defer s.mx.Unlock()

		task, err := s.convertDBTask(dbTask)
		if err != nil {
			return err
		}

		_ = s.scheduler.RemoveByTag(id)

		// Don't add it to scheduler, if it's disabled.
		if dbTask.Disabled {
			return nil
		}

		j := s.scheduler.Cron(dbTask.CronExpression).SingletonMode()
		if !dbTask.StartAt.IsZero() {
			j = j.StartAt(dbTask.StartAt)
		}

		fn := s.wrapTask(task, dbTask.ID)
		if _, err := j.Tag(dbTask.ID).Do(fn); err != nil {
			return err
		}

		return nil
	})

	return txErr
}

func (s *Service) loadFromDB() error {
	s.mx.Lock()
	defer s.mx.Unlock()

	dbTasks, err := models.FindScheduledTasks(s.db.Querier, models.ScheduledTasksFilter{
		Disabled: pointer.ToBool(false),
	})
	if err != nil {
		return err
	}

	s.scheduler.Clear()

	s.jobsMx.Lock()
	defer s.jobsMx.Unlock()

	for _, dbTask := range dbTasks {
		task, err := s.convertDBTask(dbTask)
		if err != nil {
			return err
		}

		fn := s.wrapTask(task, dbTask.ID)
		j := s.scheduler.Cron(dbTask.CronExpression).SingletonMode()
		if !dbTask.StartAt.IsZero() {
			j = j.StartAt(dbTask.StartAt)
		}
		scheduleJob, err := j.Tag(dbTask.ID).Do(fn)
		if err != nil {
			return err
		}

		s.jobs[dbTask.ID] = scheduleJob
	}

	return nil
}
func (s *Service) wrapTask(task Task, id string) func() {
	return func() {
		var err error
		l := s.l.WithFields(logrus.Fields{
			"id":       id,
			"taskType": task.Type(),
		})
		ctx, cancel := context.WithCancel(context.Background())

		s.taskMx.Lock()
		s.tasks[id] = cancel
		s.taskMx.Unlock()

		defer func() {
			cancel()
			s.taskMx.Lock()
			delete(s.tasks, id)
			s.taskMx.Unlock()
		}()

		t := time.Now()
		l.Debug("Starting task")
		_, err = models.ChangeScheduledTask(s.db.Querier, id, models.ChangeScheduledTaskParams{
			Running: pointer.ToBool(true),
		})

		if err != nil {
			l.Errorf("failed to change running state: %v", err)
		}

		taskErr := task.Run(ctx)
		if taskErr != nil {
			l.Error(taskErr)
		}
		l.WithField("duration", time.Since(t)).Debug("Ended task")

		s.taskFinished(id, taskErr)
	}
}

func (s *Service) taskFinished(id string, taskErr error) {
	s.jobsMx.RLock()
	job := s.jobs[id]
	s.jobsMx.RUnlock()

	l := s.l.WithField("id", id)

	txErr := s.db.InTransaction(func(tx *reform.TX) error {
		params := models.ChangeScheduledTaskParams{
			Running: pointer.ToBool(false),
		}

		if taskErr != nil {
			params.Error = pointer.ToString(taskErr.Error())
		} else {
			params.Error = pointer.ToString("")
		}

		if job != nil {
			params.NextRun = pointer.ToTime(job.NextRun().UTC())
			params.LastRun = pointer.ToTime(job.LastRun().UTC())
		} else {
			l.Errorf("failed to find scheduled task")
		}

		_, err := models.ChangeScheduledTask(tx.Querier, id, params)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		l.Errorf("failed to commit finished task: %v", txErr)
	}
}

func (s *Service) convertDBTask(dbTask *models.ScheduledTask) (Task, error) {
	var task Task
	switch dbTask.Type {
	case models.ScheduledPrintTask:
		task = NewPrintTask(dbTask.Data.Print.Message)
	case models.ScheduledMySQLBackupTask:
		data := dbTask.Data.MySQLBackupTask
		task = NewMySQLBackupTask(s.backupsLogicService, data.ServiceID, data.LocationID, data.Name, data.Description)
	case models.ScheduledMongoDBBackupTask:
		data := dbTask.Data.MongoDBBackupTask
		task = NewMongoBackupTask(s.backupsLogicService, data.ServiceID, data.LocationID, data.Name, data.Description)
	default:
		return task, errors.Errorf("unknown task type: %s", dbTask.Type)
	}

	task.SetID(dbTask.ID)
	return task, nil
}
