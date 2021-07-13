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

	"github.com/AlekSi/pointer"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
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
	backupSvc := NewBackupsService(db, mockedJobsService)

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
