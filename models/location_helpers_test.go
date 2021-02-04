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

package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestBackupLocations(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("create", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params := models.CreateBackupLocationParams{
			Name:        "some name",
			Description: "some desc",
			FSConfig: &models.FSLocationConfig{
				Path: "/tmp",
			},
		}

		location, err := models.CreateBackupLocation(q, params)
		require.NoError(t, err)
		assert.Equal(t, models.FSBackupLocationType, location.Type)
		assert.Equal(t, params.Name, location.Name)
		assert.Equal(t, params.Description, location.Description)
		assert.Equal(t, params.FSConfig.Path, location.FSConfig.Path)
		assert.NotEmpty(t, location.ID)
	})

	t.Run("list", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params1 := models.CreateBackupLocationParams{
			Name:        "some name",
			Description: "some desc",
			FSConfig: &models.FSLocationConfig{
				Path: "/tmp",
			},
		}
		params2 := models.CreateBackupLocationParams{
			Name:        "some name2",
			Description: "some desc",
			FSConfig: &models.FSLocationConfig{
				Path: "/tmp",
			},
		}

		loc1, err := models.CreateBackupLocation(q, params1)
		require.NoError(t, err)
		loc2, err := models.CreateBackupLocation(q, params2)
		require.NoError(t, err)

		actual, err := models.FindBackupLocations(q)
		require.NoError(t, err)
		var found1, found2 bool
		for _, channel := range actual {
			if channel.ID == loc1.ID {
				found1 = true
			}
			if channel.ID == loc2.ID {
				found2 = true
			}
		}

		assert.True(t, found1, "Fist location not found")
		assert.True(t, found2, "Second location not found")
	})
}

func TestBackupLocationValidation(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	tests := []struct {
		name     string
		location models.CreateBackupLocationParams
		errorMsg string
	}{
		{
			name: "normal fs config",
			location: models.CreateBackupLocationParams{
				Name: "fs-1",
				FSConfig: &models.FSLocationConfig{
					Path: "/tmp",
				},
			},
			errorMsg: "",
		},
		{
			name: "fs config - missing path",
			location: models.CreateBackupLocationParams{
				Name: "fs-1",
				FSConfig: &models.FSLocationConfig{
					Path: "",
				},
			},
			errorMsg: "rpc error: code = InvalidArgument desc = FS path field is empty",
		},
		{
			name: "normal s3 config",
			location: models.CreateBackupLocationParams{
				Name: "s3-1",
				S3Config: &models.S3LocationConfig{
					Endpoint:  "https://s3.us-west-2.amazonaws.com/mybucket",
					AccessKey: "access_key",
					SecretKey: "secret_key",
				},
			},
			errorMsg: "",
		},
		{
			name: "s3 config - missing endpoint",
			location: models.CreateBackupLocationParams{
				Name: "s3-2",
				S3Config: &models.S3LocationConfig{
					Endpoint:  "",
					AccessKey: "access_key",
					SecretKey: "secret_key",
				},
			},
			errorMsg: "rpc error: code = InvalidArgument desc = S3 endpoint field is empty",
		},
		{
			name: "s3 config - missing access key",
			location: models.CreateBackupLocationParams{
				Name: "s3-3",
				S3Config: &models.S3LocationConfig{
					Endpoint:  "https://s3.us-west-2.amazonaws.com/mybucket",
					AccessKey: "",
					SecretKey: "secret_key",
				},
			},
			errorMsg: "rpc error: code = InvalidArgument desc = S3 accessKey field is empty",
		},
		{
			name: "s3 config - missing secret key",
			location: models.CreateBackupLocationParams{
				Name: "s3-4",
				S3Config: &models.S3LocationConfig{
					Endpoint:  "https://s3.us-west-2.amazonaws.com/mybucket",
					AccessKey: "secret_key",
					SecretKey: "",
				},
			},
			errorMsg: "rpc error: code = InvalidArgument desc = S3 secretKey field is empty",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			tx, err := db.Begin()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, tx.Rollback())
			}()

			q := tx.Querier

			c, err := models.CreateBackupLocation(q, test.location)
			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, c)
		})
	}
}
