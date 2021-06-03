package scheduler

import (
	"context"
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"testing"
	"time"
)

func setup(t *testing.T) *Service {
	t.Helper()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	return New(db)

}
func TestService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc := setup(t)
	go func() {
		svc.Run(ctx)
	}()
	for !svc.scheduler.IsRunning() {
		// Wait a while, so scheduler is running
		time.Sleep(time.Millisecond * 10)
	}

	job := &EchoJob{EchoJobData: models.EchoJobData{Value: "test"}}
	cronExpr := "* * * * *"
	startAt := time.Now().Truncate(time.Second).UTC()
	retries := uint(3)
	retryInterval := time.Millisecond
	dbJob, err := svc.Add(job, "* * * * *", startAt, retries, retryInterval)
	assert.NoError(t, err)

	assert.Len(t, svc.scheduler.Jobs(), 1)
	findJob, err := models.FindScheduleJobByID(svc.db.Querier, dbJob.ID)
	assert.NoError(t, err)

	assert.Equal(t, startAt, dbJob.StartAt)
	assert.Equal(t, retries, dbJob.Retries)
	assert.Equal(t, retryInterval, dbJob.RetryInterval)
	assert.Equal(t, cronExpr, findJob.CronExpression)
	assert.Truef(t, dbJob.NextRun.After(startAt), "next run %s is before startAt %s", dbJob.NextRun, startAt)

}
