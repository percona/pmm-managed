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

func TestDeleteArtifact(t *testing.T) {
	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

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

	artifact, err := models.CreateArtifact(db.Querier, models.CreateArtifactParams{
		Name:       "artifact_name",
		Vendor:     "MySQL",
		LocationID: locationRes.ID,
		ServiceID:  *agent.ServiceID,
		DataModel:  "physical",
		Status:     models.SuccessBackupStatus,
	})
	require.NoError(t, err)

	findArtifact := func(artifactID string) *backupv1beta1.Artifact {
		artifacts, err := artifactsSvc.ListArtifacts(ctx, &backupv1beta1.ListArtifactsRequest{})
		require.NoError(t, err)
		require.NotNil(t, artifacts)

		var artifact *backupv1beta1.Artifact
		for _, a := range artifacts.Artifacts {
			if a.ArtifactId == artifactID {
				artifact = a
				break
			}
		}

		return artifact
	}

	mockedS3.On("RemoveRecursive", mock.Anything, endpoint, accessKey, secretKey, bucketName,
		artifact.Name+"/",
	).Return(errors.Errorf("failed to remove")).Run(func(args mock.Arguments) {
		artifact := findArtifact(artifact.ID)
		require.NotNil(t, artifact)
		assert.Equal(t, artifact.Status, backupv1beta1.BackupStatus_BACKUP_STATUS_DELETING)
	}).Once()

	_, err = artifactsSvc.DeleteArtifact(ctx, &backupv1beta1.DeleteArtifactRequest{
		ArtifactId:  artifact.ID,
		RemoveFiles: true,
	})
	require.EqualError(t, err, "failed to remove")

	foundArtifact := findArtifact(artifact.ID)
	require.NotNil(t, foundArtifact)
	assert.Equal(t, foundArtifact.Status, backupv1beta1.BackupStatus_BACKUP_STATUS_FAILED_TO_DELETE)

	mockedS3.On("RemoveRecursive", mock.Anything, endpoint, accessKey, secretKey, bucketName,
		artifact.Name+"/",
	).Return(nil).Once()

	_, err = artifactsSvc.DeleteArtifact(ctx, &backupv1beta1.DeleteArtifactRequest{
		ArtifactId:  artifact.ID,
		RemoveFiles: true,
	})
	assert.NoError(t, err)

	foundArtifact = findArtifact(artifact.ID)
	require.Nil(t, foundArtifact)

	mock.AssertExpectationsForObjects(t, mockedS3)
}
