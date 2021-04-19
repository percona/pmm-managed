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

func TestRestoreHistory(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("create", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, tx.Rollback())
		})

		q := tx.Querier

		params := models.CreateRestoreHistoryItemParams{
			ArtifactID: "artifact_id",
			ServiceID:  "service_id",
			Status:     models.InProgressRestoreStatus,
		}

		i, err := models.CreateRestoreHistoryItem(q, params)
		require.NoError(t, err)
		assert.Equal(t, params.ArtifactID, i.ArtifactID)
		assert.Equal(t, params.ServiceID, i.ServiceID)
		assert.Equal(t, params.Status, i.Status)
		assert.Less(t, time.Now().UTC().Unix()-i.StartedAt.Unix(), int64(5))
	})

	t.Run("list", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, tx.Rollback())
		})

		q := tx.Querier

		params1 := models.CreateRestoreHistoryItemParams{
			ArtifactID: "artifact_id_1",
			ServiceID:  "service_id_1",
			Status:     models.InProgressRestoreStatus,
		}
		params2 := models.CreateRestoreHistoryItemParams{
			ArtifactID: "artifact_id_2",
			ServiceID:  "service_id_2",
			Status:     models.SuccessRestoreStatus,
		}

		i1, err := models.CreateRestoreHistoryItem(q, params1)
		require.NoError(t, err)
		i2, err := models.CreateRestoreHistoryItem(q, params2)
		require.NoError(t, err)

		actual, err := models.FindRestoreHistoryItems(q)
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

		assert.Condition(t, found(i1.ID), "The first restore history item not found")
		assert.Condition(t, found(i2.ID), "The second restore history item not found")
	})

	t.Run("remove", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, tx.Rollback())
		})

		q := tx.Querier

		params := models.CreateRestoreHistoryItemParams{
			ArtifactID: "artifact_id",
			ServiceID:  "service_id",
			Status:     models.SuccessRestoreStatus,
		}
		i, err := models.CreateRestoreHistoryItem(q, params)
		require.NoError(t, err)

		err = models.RemoveRestoreHistoryItem(q, i.ID)
		require.NoError(t, err)

		artifacts, err := models.FindRestoreHistoryItems(q)
		require.NoError(t, err)
		assert.Empty(t, artifacts)
	})
}

func TestRestoreHistoryValidation(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	testCases := []struct {
		name     string
		params   models.CreateRestoreHistoryItemParams
		errorMsg string
	}{
		{
			name: "normal params",
			params: models.CreateRestoreHistoryItemParams{
				ArtifactID: "artifact_id",
				ServiceID:  "service_id",
				Status:     models.SuccessRestoreStatus,
			},
			errorMsg: "",
		},
		{
			name: "artifact missing",
			params: models.CreateRestoreHistoryItemParams{
				ServiceID: "service_id",
				Status:    models.SuccessRestoreStatus,
			},
			errorMsg: "artifact_id shouldn't be empty: invalid argument",
		},
		{
			name: "service missing",
			params: models.CreateRestoreHistoryItemParams{
				ArtifactID: "artifact_id",
				Status:     models.SuccessRestoreStatus,
			},
			errorMsg: "service_id shouldn't be empty: invalid argument",
		},
		{
			name: "invalid status",
			params: models.CreateRestoreHistoryItemParams{
				ArtifactID: "artifact_id",
				ServiceID:  "service_id",
				Status:     models.RestoreStatus("invalid"),
			},
			errorMsg: "invalid status 'invalid': invalid argument",
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			tx, err := db.Begin()
			require.NoError(t, err)
			t.Cleanup(func() {
				require.NoError(t, tx.Rollback())
			})

			q := tx.Querier

			c, err := models.CreateRestoreHistoryItem(q, test.params)
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
