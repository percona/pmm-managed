package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"

	"github.com/percona/pmm-managed/models"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

type Service struct {
	db        *reform.DB
	scheduler *gocron.Scheduler
	l         *logrus.Entry
	jobs      map[string]context.CancelFunc
	jobsMx    sync.RWMutex
}

func New(db *reform.DB) *Service {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.TagsUnique()
	return &Service{
		db:        db,
		scheduler: scheduler,
		jobs:      make(map[string]context.CancelFunc),
		l:         logrus.WithField("component", "scheduler"),
	}
}

func (s *Service) Run() {
	if err := s.loadFromDB(); err != nil {
		panic(err)
	}
	s.scheduler.StartBlocking()
}

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
	_, err = s.db.DeleteFrom(models.ScheduleJobTable, " WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Add(job Job, cronExpr string, startAt time.Time, retry uint) error {
	j := s.scheduler.Cron(cronExpr).SingletonMode()
	if !startAt.IsZero() {
		j = j.StartAt(startAt)
	}
	var data models.ScheduleJobData
	switch job.Type() {
	case models.ScheduleEchoJob:
		val := models.EchoJobData(*job.(*EchoJob))
		data.Echo = &val
	}
	err := s.db.InTransaction(func(tx *reform.TX) error {
		id := uuid.Must(uuid.NewV4()).String()
		fn := s.wrapJob(job, id, int(retry))
		scheduledJob, err := j.Tag(id).Do(fn)
		if err != nil {
			return err
		}

		dbJob := &models.ScheduleJob{
			ID:             id,
			CronExpression: cronExpr,
			StartAt:        startAt,
			LastRun:        scheduledJob.LastRun(),
			NextRun:        scheduledJob.NextRun(),
			Type:           job.Type(),
			Data:           &data,
			Retries:        retry,
		}

		if err := s.db.Insert(dbJob); err != nil {
			s.scheduler.RemoveByReference(scheduledJob)
			return err
		}
		return nil
	})
	return err
}

func (s *Service) loadFromDB() error {
	s.scheduler.Clear()
	structs, err := s.db.SelectAllFrom(models.ScheduleJobTable, "")
	if err != nil {
		return err
	}
	for _, entry := range structs {
		scheduleJob := entry.(*models.ScheduleJob)
		var job Job
		switch scheduleJob.Type {
		case models.ScheduleEchoJob:
			val := EchoJob(*scheduleJob.Data.Echo)
			job = &val
		}
		fn := s.wrapJob(job, scheduleJob.ID, int(scheduleJob.Retries))
		j := s.scheduler.Cron(scheduleJob.CronExpression).SingletonMode()
		if !scheduleJob.StartAt.IsZero() {
			j = j.StartAt(scheduleJob.StartAt)
		}
		if _, err := j.Tag(scheduleJob.ID).Do(fn); err != nil {
			return err
		}

	}
	return nil
}
func (s *Service) wrapJob(job Job, id string, retry int) func() {
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

		for once := true; once || retry > 0; once = false {
			t := time.Now()
			l.Debug("Starting job")
			err = job.Do(ctx)
			l.WithField("duration", time.Since(t)).Debug("Ended  job")
			if err == nil {
				break
			} else {
				l.Error(err)
				if err == context.Canceled {
					break
				}
			}
			retry--
		}
		s.jobFinished(id, err)
	}
}

func (s *Service) jobFinished(id string, err error) {
	var job *gocron.Job
	for _, j := range s.scheduler.Jobs() {
		if len(j.Tags()) > 0 && j.Tags()[0] == id {
			job = j
			break
		}
	}
	if job == nil {
		return
	}
	fmt.Println(job.LastRun(), job.NextRun())
}
