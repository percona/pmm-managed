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
	"fmt"
	"sync"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/percona/pmm-managed/models"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

// Service is responsive for executing tasks and storing them to DB.
type Service struct {
	db        *reform.DB
	scheduler *gocron.Scheduler
	l         *logrus.Entry
	tasks     map[string]context.CancelFunc
	taskMx    sync.RWMutex
	jobsMx    sync.Mutex
}

// New creates new scheduler service.
func New(db *reform.DB) *Service {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()
	scheduler.WaitForScheduleAll()
	return &Service{
		db:        db,
		scheduler: scheduler,
		tasks:     make(map[string]context.CancelFunc),
		l:         logrus.WithField("component", "scheduler"),
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
	s.jobsMx.Lock()
	defer s.jobsMx.Unlock()

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

		j := s.scheduler.Cron(params.CronExpression).SingletonMode()
		if !params.StartAt.IsZero() {
			j = j.StartAt(params.StartAt)
		}
		scheduleJob, err := j.Tag(id).Do(fn)
		if err != nil {
			return err
		}

		scheduledTask, err = models.ChangeScheduledTask(tx.Querier, id, models.ChangeScheduledTaskParams{
			NextRun: scheduleJob.NextRun(),
			LastRun: scheduleJob.LastRun(),
		})
		if err != nil {
			s.l.WithField("id", id).Errorf("failed to set next run for new created task")
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
	defer s.jobsMx.Unlock()

	err := s.db.InTransaction(func(tx *reform.TX) error {
		if err := models.RemoveScheduledTask(tx.Querier, id); err != nil {
			return err
		}

		_ = s.scheduler.RemoveByTag(id)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Reload removes job from scheduler and add it again from DB.
func (s *Service) Reload(id string) error {
	dbTask, err := models.FindScheduledTaskByID(s.db.Querier, id)
	if err != nil {
		return err
	}

	if dbTask.Running {
		return fmt.Errorf("task is running")
	}

	s.jobsMx.Lock()
	defer s.jobsMx.Unlock()

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
}

func (s *Service) loadFromDB() error {
	s.jobsMx.Lock()
	defer s.jobsMx.Unlock()

	disabled := false
	dbTasks, err := models.FindScheduledTasks(s.db.Querier, models.ScheduledTasksFilter{
		Disabled: &disabled,
	})
	if err != nil {
		return err
	}

	tasks := make([]Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		task, err := s.convertDBTask(dbTask)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	s.scheduler.Clear()
	for i, task := range tasks {
		dbTask := dbTasks[i]
		fn := s.wrapTask(task, dbTask.ID)
		j := s.scheduler.Cron(dbTask.CronExpression).SingletonMode()
		if !dbTask.StartAt.IsZero() {
			j = j.StartAt(dbTask.StartAt)
		}
		if _, err := j.Tag(dbTask.ID).Do(fn); err != nil {
			return err
		}
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

		_, err = models.ChangeScheduledTask(s.db.Querier, id, models.ChangeScheduledTaskParams{
			Running: pointer.ToBool(false),
		})

		if err != nil {
			l.Errorf("failed to change running status: %v", err)
		}

		s.taskFinished(id, taskErr)
	}
}

func (s *Service) taskFinished(id string, taskErr error) {
	var job *gocron.Job
	for _, j := range s.scheduler.Jobs() {
		if len(j.Tags()) > 0 && j.Tags()[0] == id {
			job = j
			break
		}
	}
	l := s.l.WithField("id", id)

	dbTask, err := models.FindScheduledTaskByID(s.db.Querier, id)
	if err != nil {
		return
	}

	params := models.ChangeScheduledTaskParams{
		Succeeded: pointer.ToUint(dbTask.Succeeded),
		Failed:    pointer.ToUint(dbTask.Failed),
		Running:   pointer.ToBool(false),
	}

	if taskErr == nil {
		params.Succeeded = pointer.ToUint(dbTask.Succeeded + 1)
		params.Error = pointer.ToString("")
	} else {
		params.Failed = pointer.ToUint(dbTask.Failed + 1)
		params.Error = pointer.ToString(taskErr.Error())
	}

	if job != nil {
		params.NextRun = job.NextRun().UTC()
		params.LastRun = job.LastRun().UTC()
	} else {
		l.Errorf("failed to find scheduled task: %v", err)
	}

	_, err = models.ChangeScheduledTask(s.db.Querier, id, params)
	if err != nil {
		l.Errorf("failed to change scheduled task: %v", err)
	}
}

func (s *Service) convertDBTask(dbTask *models.ScheduledTask) (Task, error) {
	var task Task
	switch dbTask.Type {
	case models.ScheduledPrintTask:
		task = NewPrintTask(dbTask.Data.Print.Message)
	default:
		return task, fmt.Errorf("unknown task type: %s", dbTask.Type)
	}

	task.SetID(dbTask.ID)
	return task, nil
}
