package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

func setup(t *testing.T) *Service {
	t.Helper()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	backupsLogicService := &mockBackupsLogicService{}
	return New(db, backupsLogicService)

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

	task := NewPrintTask("test")
	cronExpr := "* * * * *"
	startAt := time.Now().Truncate(time.Second).UTC()
	retries := uint(3)
	retryInterval := time.Millisecond
	dbTask, err := svc.Add(task, AddParams{
		CronExpression: cronExpr,
		Retry:          retries,
		RetryInterval:  retryInterval,
	})
	assert.NoError(t, err)

	assert.Len(t, svc.scheduler.Jobs(), 1)
	findJob, err := models.FindScheduledTaskByID(svc.db.Querier, dbTask.ID)
	assert.NoError(t, err)

	assert.Equal(t, startAt, dbTask.StartAt)
	assert.Equal(t, retries, dbTask.Retries)
	assert.Equal(t, retryInterval, dbTask.RetryInterval)
	assert.Equal(t, cronExpr, findJob.CronExpression)
	assert.Truef(t, dbTask.NextRun.After(startAt), "next run %s is before startAt %s", dbTask.NextRun, startAt)

	err = svc.Remove(dbTask.ID)
	assert.NoError(t, err)
	assert.Len(t, svc.scheduler.Jobs(), 0)
	_, err = models.FindScheduledTaskByID(svc.db.Querier, dbTask.ID)
	tests.AssertGRPCError(t, status.Newf(codes.NotFound, `ScheduledTask with ID "%s" not found.`, dbTask.ID), err)

}
