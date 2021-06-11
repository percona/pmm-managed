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

// Add adds task to scheduler and save it to DB.
func (s *Service) Add(task Task, enabled bool, cronExpr string, startAt time.Time, retry uint, retryInterval time.Duration) (*models.ScheduledTask, error) {
	var scheduledTask *models.ScheduledTask
	var err error
	s.jobsMx.Lock()
	defer s.jobsMx.Unlock()

	err = s.db.InTransaction(func(tx *reform.TX) error {
		scheduledTask, err = models.CreateScheduledTask(tx.Querier, models.CreateScheduledTaskParams{
			CronExpression: cronExpr,
			StartAt:        startAt,
			Type:           task.Type(),
			Data:           task.Data(),
			Retries:        retry,
			RetryInterval:  retryInterval,
			Disabled:       !enabled,
		})
		if err != nil {
			return err
		}

		id := scheduledTask.ID
		task.SetID(id)
		fn := s.wrapTask(task, id, int(retry), retryInterval)

		j := s.scheduler.Cron(cronExpr).SingletonMode()
		if !startAt.IsZero() {
			j = j.StartAt(startAt)
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

	return models.RemoveScheduledTask(s.db.Querier, id)
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

	fn := s.wrapTask(task, dbTask.ID, int(dbTask.RetriesRemaining), dbTask.RetryInterval)
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
	// Reset tags
	s.scheduler.TagsUnique()
	for i, task := range tasks {
		dbTask := dbTasks[i]
		fn := s.wrapTask(task, dbTask.ID, int(dbTask.RetriesRemaining), dbTask.RetryInterval)
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
func (s *Service) wrapTask(task Task, id string, retry int, retryInterval time.Duration) func() {
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
		retriesRemaining := retry
		succeeded := false
		for {
			t := time.Now()
			l.Debug("Starting task")
			_, err = models.ChangeScheduledTask(s.db.Querier, id, models.ChangeScheduledTaskParams{
				Running: pointer.ToBool(true),
			})

			if err != nil {
				l.Errorf("failed to change running state: %v", err)
			}

			err = task.Run(ctx)
			l.WithField("duration", time.Since(t)).Debug("Ended task")
			if err == nil {
				succeeded = true
				break
			}

			if err == context.Canceled {
				break
			}
			l.Error(err)

			if retriesRemaining <= 0 {
				break
			}
			retriesRemaining--
			_, err = models.ChangeScheduledTask(s.db.Querier, id, models.ChangeScheduledTaskParams{
				RetriesRemaining: pointer.ToUint(uint(retriesRemaining)),
				Running:          pointer.ToBool(false),
			})

			if err != nil {
				l.Errorf("failed to change retries remaining: %v", err)
			}

			timer := time.NewTimer(retryInterval)
			select {
			case <-ctx.Done():
			case <-timer.C:
			}
			timer.Stop()
		}
		s.taskFinished(id, succeeded)
	}
}

func (s *Service) taskFinished(id string, succeeded bool) {
	var job *gocron.Job
	for _, j := range s.scheduler.Jobs() {
		if len(j.Tags()) > 0 && j.Tags()[0] == id {
			job = j
			break
		}
	}
	l := s.l.WithField("id", id)
	if job == nil {
		l.Warn("couldn't find finished job in scheduler")
		return
	}

	dbTask, err := models.FindScheduledTaskByID(s.db.Querier, id)
	if err != nil {
		l.Errorf("failed to find scheduled task: %v", err)
		return
	}

	if succeeded {
		dbTask.Succeeded++
	} else {
		dbTask.Failed++
	}

	params := models.ChangeScheduledTaskParams{
		RetriesRemaining: pointer.ToUint(dbTask.Retries),
		NextRun:          job.NextRun(),
		LastRun:          job.LastRun(),
		Succeeded:        pointer.ToUint(dbTask.Succeeded),
		Failed:           pointer.ToUint(dbTask.Failed),
		Running:          pointer.ToBool(false),
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
