package models_test

import (
	"testing"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

func TestChecksState(t *testing.T) {
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

		actual, err := models.CreateCheckState(q, "check-name", models.Standard)
		require.NoError(t, err)
		assert.Equal(t, "check-name", actual.Name)
		assert.Equal(t, models.Standard, actual.Interval)
	})

	t.Run("change", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		oldState, err := models.CreateCheckState(q, "check-name", models.Standard)
		require.NoError(t, err)
		assert.Equal(t, "check-name", oldState.Name)
		assert.Equal(t, models.Standard, oldState.Interval)

		newState, err := models.ChangeCheckState(q, "check-name", models.Rare)
		require.NoError(t, err)
		assert.Equal(t, oldState.Name, newState.Name)
		assert.NotEqual(t, oldState.Interval, newState.Interval)
		assert.Equal(t, models.Rare, newState.Interval)
	})

	t.Run("find", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		expected, err := models.CreateCheckState(q, "check-name", models.Standard)
		require.NoError(t, err)

		actual, err := models.FindCheckStateByName(q, "check-name")
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
