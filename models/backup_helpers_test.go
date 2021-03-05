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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestBackup(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("create backup", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params := models.CreateBackupParams{
			Name:         "name",
			LocationName: "location name",
		}

		backup, err := models.CreateBackup(q, params)
		require.NoError(t, err)
		assert.Equal(t, params.Name, backup.Name)
		assert.Equal(t, params.LocationName, backup.LocationName)
		assert.Less(t, time.Now().UTC().Unix(), backup.CreatedAt.Unix())
	})

	t.Run("list", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params1 := models.CreateBackupParams{
			Name:         "name 1",
			LocationName: "location name 1",
		}
		params2 := models.CreateBackupParams{
			Name:         "name 2",
			LocationName: "location name 2",
		}

		b1, err := models.CreateBackup(q, params1)
		require.NoError(t, err)
		b2, err := models.CreateBackup(q, params2)
		require.NoError(t, err)

		actual, err := models.FindBackups(q)
		require.NoError(t, err)

		found := func(id string) func() bool {
			return func() bool {
				for _, b := range actual {
					if b.ID == id {
						return true
					}
				}
				return false
			}
		}

		assert.Condition(t, found(b1.ID), "The first backup not found")
		assert.Condition(t, found(b2.ID), "The second backup not found")
	})

	t.Run("remove", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params := models.CreateBackupParams{
			Name:         "some name",
			LocationName: "some location name",
		}

		b, err := models.CreateBackup(q, params)
		require.NoError(t, err)

		err = models.RemoveBackup(q, b.ID)
		require.NoError(t, err)

		backups, err := models.FindBackups(q)
		require.NoError(t, err)
		assert.Empty(t, backups)
	})
}

func TestBackupValidation(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	testCases := []struct {
		name     string
		params   models.CreateBackupParams
		errorMsg string
	}{
		{
			name: "normal params",
			params: models.CreateBackupParams{
				Name:         "name",
				LocationName: "location name",
			},
			errorMsg: "",
		},
		{
			name: "name missing",
			params: models.CreateBackupParams{
				LocationName: "location name",
			},
			errorMsg: "backup name shouldn't be empty: invalid argument",
		},
		{
			name: "location name missing",
			params: models.CreateBackupParams{
				Name: "name",
			},
			errorMsg: "backup location name shouldn't be empty: invalid argument",
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			tx, err := db.Begin()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, tx.Rollback())
			}()

			q := tx.Querier

			c, err := models.CreateBackup(q, test.params)
			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
				assert.Nil(t, c)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, c)
		})
	}
}
