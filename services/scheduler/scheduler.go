package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/percona/pmm-managed/models"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

type jobsDeps struct {
}

// Service is responsive for executing jobs and storing them to DB.
type Service struct {
	db        *reform.DB
	scheduler *gocron.Scheduler
	l         *logrus.Entry
	jobs      map[string]context.CancelFunc
	jobsMx    sync.RWMutex
	jobsDeps  jobsDeps
}

// New creates new scheduler service.
func New(db *reform.DB) *Service {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()
	// @TODO accept deps and fill jobsDeps
	return &Service{
		db:        db,
		scheduler: scheduler,
		jobs:      make(map[string]context.CancelFunc),
		l:         logrus.WithField("component", "scheduler"),
	}
}

// Run loads jobs from DB and starts scheduler.
func (s *Service) Run(ctx context.Context) {
	if err := s.loadFromDB(); err != nil {
		s.l.Warn(err)
	}
	s.scheduler.StartAsync()
	<-ctx.Done()
	s.scheduler.Stop()
}

// Add adds job to scheduler and save it to DB.
func (s *Service) Add(job Job, cronExpr string, startAt time.Time, retry uint, retryInterval time.Duration) error {
	j := s.scheduler.Cron(cronExpr).SingletonMode()
	if !startAt.IsZero() {
		j = j.StartAt(startAt)
	}

	err := s.db.InTransaction(func(tx *reform.TX) error {
		dbJob, err := models.CreateScheduleJob(tx.Querier, models.CreateScheduleJobParams{
			CronExpression: cronExpr,
			StartAt:        startAt,
			Type:           job.Type(),
			Data:           job.Data(),
			Retries:        retry,
			RetryInterval:  retryInterval,
		})
		if err != nil {
			return err
		}
		fn := s.wrapJob(job, dbJob.ID, int(retry), retryInterval)
		scheduleJob, err := j.Tag(dbJob.ID).Do(fn)
		if err != nil {
			return err
		}

		_, err = models.ChangeScheduleJob(tx.Querier, dbJob.ID, models.ChangeScheduleJobParams{
			NextRun: scheduleJob.NextRun(),
			LastRun: scheduleJob.LastRun(),
		})
		if err != nil {
			s.l.WithField("id", dbJob.ID).Errorf("failed to set next run for new created job")
		}

		return nil
	})
	return err
}

// Remove stops job specified by id and removes it from DB and scheduler.
func (s *Service) Remove(id string) error {
	s.jobsMx.RLock()
	if cancel, ok := s.jobs[id]; ok {
		cancel()
	}
	s.jobsMx.RUnlock()
	err := s.scheduler.RemoveByTag(id)
	if err != nil {
		return err
	}

	if err := models.RemoveScheduleJob(s.db.Querier, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) loadFromDB() error {
	s.scheduler.Clear()
	// Reset tags
	s.scheduler.TagsUnique()

	disabled := false
	jobs, err := models.FindScheduleJobs(s.db.Querier, models.ScheduleJobsFilter{
		Disabled: &disabled,
	})
	if err != nil {
		return err
	}

	for _, dbJob := range jobs {
		job, err := s.convertDBJob(dbJob)
		if err != nil {
			return err
		}
		fn := s.wrapJob(job, dbJob.ID, int(dbJob.Retries), dbJob.RetryInterval)
		j := s.scheduler.Cron(dbJob.CronExpression).SingletonMode()
		if !dbJob.StartAt.IsZero() {
			j = j.StartAt(dbJob.StartAt)
		}
		if _, err := j.Tag(dbJob.ID).Do(fn); err != nil {
			return err
		}

	}
	return nil
}
func (s *Service) wrapJob(job Job, id string, retry int, retryInterval time.Duration) func() {
	return func() {
		var err error
		l := s.l.WithField("jobType", job.Type())
		ctx, cancel := context.WithCancel(context.Background())

		s.jobsMx.Lock()
		s.jobs[id] = cancel
		s.jobsMx.Unlock()

		defer func() {
			cancel()
			s.jobsMx.Lock()
			delete(s.jobs, id)
			s.jobsMx.Unlock()
		}()

		for {
			t := time.Now()
			l.Debug("Starting job")
			err = job.Do(ctx)
			l.WithField("duration", time.Since(t)).Debug("Ended  job")
			if err == nil || err == context.Canceled {
				break
			} else {
				l.Error(err)
			}
			if retry <= 0 {
				break
			}
			retry--
			time.Sleep(retryInterval)
		}
		s.jobFinished(id)
	}
}

func (s *Service) jobFinished(id string) {
	var job *gocron.Job
	for _, j := range s.scheduler.Jobs() {
		if len(j.Tags()) > 0 && j.Tags()[0] == id {
			job = j
			break
		}
	}
	l := s.l.WithField("schedule_job_id", id)
	if job == nil {
		l.Warn("couldn't find finished job in scheduler")
		return
	}

	_, err := models.ChangeScheduleJob(s.db.Querier, id, models.ChangeScheduleJobParams{
		NextRun: job.NextRun(),
		LastRun: job.LastRun(),
	})
	if err != nil {
		l.Error("failed to change schedule job")
	}
}

func (s *Service) convertDBJob(dbJob *models.ScheduleJob) (Job, error) {
	var job Job
	switch dbJob.Type {
	case models.ScheduleEchoJob:
		val := EchoJob{
			jobsDeps:    s.jobsDeps,
			EchoJobData: *dbJob.Data.Echo,
		}
		job = &val
	default:
		return job, fmt.Errorf("unknown job type: %s", dbJob.Type)
	}
	return job, nil
}
