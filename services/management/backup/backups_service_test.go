package backup

import (
	"context"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/scheduler"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

func setup(t *testing.T, q *reform.Querier, serviceName string) *models.Agent {
	t.Helper()
	node, err := models.CreateNode(q, models.GenericNodeType, &models.CreateNodeParams{
		NodeName: "test-node",
	})
	require.NoError(t, err)

	pmmAgent, err := models.CreatePMMAgent(q, node.NodeID, nil)
	require.NoError(t, err)
	require.NoError(t, q.Update(pmmAgent))

	mysql, err := models.AddNewService(q, models.MySQLServiceType, &models.AddDBMSServiceParams{
		ServiceName: serviceName,
		NodeID:      node.NodeID,
		Address:     pointer.ToString("127.0.0.1"),
		Port:        pointer.ToUint16(3306),
	})
	require.NoError(t, err)

	agent, err := models.CreateAgent(q, models.MySQLdExporterType, &models.CreateAgentParams{
		PMMAgentID: pmmAgent.AgentID,
		ServiceID:  mysql.ServiceID,
		Username:   "user",
		Password:   "password",
	})
	require.NoError(t, err)
	return agent
}

func TestScheduledBackups(t *testing.T) {
	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	backupService := &mockBackupService{}
	schedulerService := scheduler.New(db, backupService)
	backupSvc := NewBackupsService(db, backupService, schedulerService)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	agent := setup(t, db.Querier, t.Name())
	locationRes, err := models.CreateBackupLocation(db.Querier, models.CreateBackupLocationParams{
		Name:        "Test location",
		Description: "Test description",
		BackupLocationConfig: models.BackupLocationConfig{
			S3Config: &models.S3LocationConfig{
				Endpoint:     "https://s3.us-west-2.amazonaws.com/",
				AccessKey:    "access_key",
				SecretKey:    "secret_key",
				BucketName:   "example_bucket",
				BucketRegion: "us-east-2",
			},
		},
	})
	require.NoError(t, err)

	t.Run("schedule/change", func(t *testing.T) {
		req := &backupv1beta1.ScheduleBackupRequest{
			ServiceId:      pointer.GetString(agent.ServiceID),
			LocationId:     locationRes.ID,
			CronExpression: "1 * * * *",
			StartTime:      timestamppb.New(time.Now()),
			Name:           t.Name(),
			Description:    t.Name(),
			Enabled:        true,
		}
		res, err := backupSvc.ScheduleBackup(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.ScheduledBackupId)

		task, err := models.FindScheduledTaskByID(db.Querier, res.ScheduledBackupId)
		require.NoError(t, err)
		assert.Equal(t, models.ScheduledMySQLBackupTask, task.Type)
		assert.Equal(t, req.CronExpression, task.CronExpression)
		data := task.Data.MySQLBackupTask
		assert.Equal(t, req.Name, data.Name)
		assert.Equal(t, req.Description, data.Description)
		assert.Equal(t, req.ServiceId, data.ServiceID)
		assert.Equal(t, req.LocationId, data.LocationID)

		changeReq := &backupv1beta1.ChangeScheduledBackupRequest{
			ScheduledBackupId: task.ID,
			Enabled:           wrapperspb.Bool(false),
			CronExpression:    wrapperspb.String("2 * * * *"),
			StartTime:         timestamppb.New(time.Now()),
			Name:              wrapperspb.String("test"),
			Description:       wrapperspb.String("test"),
		}
		_, err = backupSvc.ChangeScheduledBackup(ctx, changeReq)

		assert.NoError(t, err)
		task, err = models.FindScheduledTaskByID(db.Querier, res.ScheduledBackupId)
		require.NoError(t, err)
		data = task.Data.MySQLBackupTask
		assert.Equal(t, changeReq.CronExpression.GetValue(), task.CronExpression)
		assert.Equal(t, changeReq.Enabled.GetValue(), !task.Disabled)
		assert.Equal(t, changeReq.Name.GetValue(), data.Name)
		assert.Equal(t, changeReq.Description.GetValue(), data.Description)
	})

	t.Run("list", func(t *testing.T) {
		// Add non-backup task
		_, err := models.CreateScheduledTask(db.Querier, models.CreateScheduledTaskParams{
			CronExpression: "* * * * *",
			Type:           models.ScheduledPrintTask,
			Data: models.ScheduledTaskData{
				Print: &models.PrintTaskData{
					Message: "42",
				},
			},
		})
		require.NoError(t, err)

		res, err := backupSvc.ListScheduledBackups(ctx, &backupv1beta1.ListScheduledBackupsRequest{})

		assert.NoError(t, err)
		assert.Len(t, res.ScheduledBackups, 1)
	})

	t.Run("remove", func(t *testing.T) {
		task, err := models.CreateScheduledTask(db.Querier, models.CreateScheduledTaskParams{
			CronExpression: "* * * * *",
			Type:           models.ScheduledMySQLBackupTask,
			Data:           models.ScheduledTaskData{},
		})
		require.NoError(t, err)

		id := task.ID
		_, err = backupSvc.RemoveScheduledBackup(ctx, &backupv1beta1.RemoveScheduledBackupRequest{
			ScheduledBackupId: task.ID,
		})
		assert.NoError(t, err)

		task, err = models.FindScheduledTaskByID(db.Querier, task.ID)
		assert.Nil(t, task)
		tests.AssertGRPCError(t, status.Newf(codes.NotFound, `ScheduledTask with ID "%s" not found.`, id), err)
	})

}
