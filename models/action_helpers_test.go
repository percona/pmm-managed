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

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestActionHelpers(t *testing.T) {
	now, origNowF := models.Now(), models.Now
	models.Now = func() time.Time {
		return now
	}
	sqlDB := testdb.Open(t)
	defer func() {
		models.Now = origNowF
		require.NoError(t, sqlDB.Close())
	}()

	setup := func(t *testing.T) (q *reform.Querier, teardown func(t *testing.T)) {
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		tx, err := db.Begin()
		require.NoError(t, err)
		q = tx.Querier

		for _, str := range []reform.Struct{
			&models.ActionResult{
				ID:         "A1",
				PmmAgentID: "A2",
			},
		} {
			require.NoError(t, q.Insert(str))
		}

		teardown = func(t *testing.T) {
			require.NoError(t, tx.Rollback())
		}
		return
	}

	t.Run("FindActionResultByID", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		a, err := models.FindActionResultByID(q, "A1")
		require.NoError(t, err)
		assert.NotEmpty(t, a.ID)

		_, err = models.FindActionResultByID(q, "A2")
		require.Error(t, err)
	})

	t.Run("FindPmmAgentIDToRunAction", func(t *testing.T) {
		_, teardown := setup(t)
		defer teardown(t)

		a := []*models.Agent{
			{AgentID: "A1", AgentType: models.PMMAgentType},
			{AgentID: "A2", AgentType: models.MySQLdExporterType, PMMAgentID: pointer.ToString("A1")},
		}

		id, err := models.FindPmmAgentIDToRunAction("A1", a)
		require.NoError(t, err)
		assert.Equal(t, "A1", id)

		a2 := []*models.Agent{
			{AgentID: "A1", AgentType: models.PMMAgentType},
			{AgentID: "A2", AgentType: models.MySQLdExporterType, PMMAgentID: pointer.ToString("A1")},
			{AgentID: "A3", AgentType: models.MySQLdExporterType, PMMAgentID: pointer.ToString("A1")},
		}

		id, err = models.FindPmmAgentIDToRunAction("A3", a2)
		require.NoError(t, err)
		assert.Equal(t, "A3", id)

		_, err = models.FindPmmAgentIDToRunAction("A4", a2)
		require.Error(t, err)
		tests.AssertGRPCError(t, status.New(codes.FailedPrecondition, "couldn't find pmm-agent-id to run action"), err)

		_, err = models.FindPmmAgentIDToRunAction("", a2)
		require.Error(t, err)
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "couldn't find pmm-agent-id to run action"), err)
	})

}
