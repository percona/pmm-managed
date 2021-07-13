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
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/services/backup"

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

func TestBackup(t *testing.T) {
	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	mockedJobsService := &mockJobsService{}
	mockedJobsService.On("StartMySQLBackupJob", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	s := backup.NewService(db, mockedJobsService)
	backupSvc := NewBackupsService(db, s, nil)

	mockedS3 := &mockAwsS3{}
	artifactsSvc := NewArtifactsService(db, mockedS3)

	agent := setup(t, db.Querier, "test-service")
	endpoint := "https://s3.us-west-2.amazonaws.com/"
	accessKey, secretKey, bucketName, bucketRegion := "access_key", "secret_key", "example_bucket", "us-east-2"

	locationRes, err := models.CreateBackupLocation(db.Querier, models.CreateBackupLocationParams{
		Name:        "Test location",
		Description: "Test description",
		BackupLocationConfig: models.BackupLocationConfig{
			S3Config: &models.S3LocationConfig{
				Endpoint:     endpoint,
				AccessKey:    accessKey,
				SecretKey:    secretKey,
				BucketName:   bucketName,
				BucketRegion: bucketRegion,
			},
		},
	})
	require.NoError(t, err)

	backupRes, err := backupSvc.StartBackup(ctx, &backupv1beta1.StartBackupRequest{
		ServiceId:  pointer.GetString(agent.ServiceID),
		LocationId: locationRes.ID,
		Name:       "Test backup",
	})
	require.NoError(t, err)

	artifact, err := models.FindArtifactByID(db.Querier, backupRes.ArtifactId)
	require.NoError(t, err)

	require.NotNil(t, artifact)
	assert.Equal(t, locationRes.ID, artifact.LocationID)
	assert.Equal(t, *agent.ServiceID, artifact.ServiceID)
	assert.EqualValues(t, models.MySQLServiceType, artifact.Vendor)

	_, err = models.UpdateArtifact(db.Querier, backupRes.ArtifactId, models.UpdateArtifactParams{
		Status: models.BackupStatusPointer(models.SuccessBackupStatus),
	})
	require.NoError(t, err)

	mockedS3.On(
		"RemoveRecursive",
		mock.Anything,
		endpoint,
		accessKey,
		secretKey,
		bucketName,
		artifact.Name+"/",
	).Return(nil).Once()

	_, err = artifactsSvc.DeleteArtifact(ctx, &backupv1beta1.DeleteArtifactRequest{
		ArtifactId:  backupRes.ArtifactId,
		RemoveFiles: true,
	})
	require.NoError(t, err)

	artifact, err = models.FindArtifactByID(db.Querier, backupRes.ArtifactId)
	assert.Nil(t, artifact)
	assert.True(t, assert.True(t, errors.Is(err, models.ErrNotFound)))

	mock.AssertExpectationsForObjects(t, mockedS3, mockedJobsService)
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

		_, err = models.CreateArtifact(db.Querier, models.CreateArtifactParams{
			Name:       "artifact",
			Vendor:     "mysql",
			LocationID: locationRes.ID,
			ServiceID:  *agent.ServiceID,
			DataModel:  "physical",
			Status:     "pending",
			ScheduleID: id,
		})
		require.NoError(t, err)

		_, err = backupSvc.RemoveScheduledBackup(ctx, &backupv1beta1.RemoveScheduledBackupRequest{
			ScheduledBackupId: task.ID,
		})
		assert.NoError(t, err)

		task, err = models.FindScheduledTaskByID(db.Querier, task.ID)
		assert.Nil(t, task)
		tests.AssertGRPCError(t, status.Newf(codes.NotFound, `ScheduledTask with ID "%s" not found.`, id), err)

		artifacts, err := models.FindArtifacts(db.Querier, &models.ArtifactFilters{
			ScheduleID: id,
		})

		assert.NoError(t, err)
		assert.Len(t, artifacts, 0)
	})

}
