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

func TestAgentHelpers(t *testing.T) {
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
			&models.Node{
				NodeID:   "N1",
				NodeType: models.GenericNodeType,
				NodeName: "Node with Service",
			},

			&models.Service{
				ServiceID:   "S1",
				ServiceType: models.MySQLServiceType,
				ServiceName: "Service on N1",
				NodeID:      "N1",
			},

			&models.Agent{
				AgentID:      "A1",
				AgentType:    models.PMMAgentType,
				RunsOnNodeID: pointer.ToString("N1"),
			},
			&models.Agent{
				AgentID:      "A2",
				AgentType:    models.MySQLdExporterType,
				PMMAgentID:   pointer.ToString("A1"),
				RunsOnNodeID: nil,
			},
			&models.Agent{
				AgentID:      "A3",
				AgentType:    models.NodeExporterType,
				PMMAgentID:   pointer.ToString("A1"),
				RunsOnNodeID: nil,
			},

			&models.AgentNode{
				AgentID: "A3",
				NodeID:  "N1",
			},

			&models.AgentService{
				AgentID:   "A2",
				ServiceID: "S1",
			},
		} {
			require.NoError(t, q.Insert(str))
		}

		teardown = func(t *testing.T) {
			require.NoError(t, tx.Rollback())
		}
		return
	}

	t.Run("AgentsForNode", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agents, err := models.AgentsForNode(q, "N1")
		require.NoError(t, err)
		expected := []*models.Agent{{
			AgentID:      "A3",
			AgentType:    models.NodeExporterType,
			PMMAgentID:   pointer.ToStringOrNil("A1"),
			RunsOnNodeID: nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}}
		assert.Equal(t, expected, agents)
	})

	t.Run("AgentsRunningByPMMAgent", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agents, err := models.AgentsRunningByPMMAgent(q, "A1")
		require.NoError(t, err)
		expected := []*models.Agent{{
			AgentID:      "A2",
			AgentType:    models.MySQLdExporterType,
			PMMAgentID:   pointer.ToStringOrNil("A1"),
			RunsOnNodeID: nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}, {
			AgentID:      "A3",
			AgentType:    models.NodeExporterType,
			PMMAgentID:   pointer.ToStringOrNil("A1"),
			RunsOnNodeID: nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}}
		assert.Equal(t, expected, agents)
	})

	t.Run("AgentsForService", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agents, err := models.AgentsForService(q, "S1")
		require.NoError(t, err)
		expected := []*models.Agent{{
			AgentID:      "A2",
			AgentType:    models.MySQLdExporterType,
			PMMAgentID:   pointer.ToStringOrNil("A1"),
			RunsOnNodeID: nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}}
		assert.Equal(t, expected, agents)
	})

	t.Run("PMMAgentsForChangedNode", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		ids, err := models.PMMAgentsForChangedNode(q, "N1")
		require.NoError(t, err)
		assert.Equal(t, []string{"A1", models.PMMServerAgentID}, ids)
	})

	t.Run("PMMAgentsForChangedService", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		ids, err := models.PMMAgentsForChangedService(q, "S1")
		require.NoError(t, err)
		assert.Equal(t, []string{"A1", models.PMMServerAgentID}, ids)
	})

	t.Run("RemoveAgent", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agent, err := models.RemoveAgent(q, "", models.RemoveRestrict)
		assert.Nil(t, agent)
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, `Empty Agent ID.`), err)

		agent, err = models.RemoveAgent(q, "A0", models.RemoveRestrict)
		assert.Nil(t, agent)
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Agent with ID "A0" not found.`), err)

		agent, err = models.RemoveAgent(q, "A1", models.RemoveRestrict)
		assert.Nil(t, agent)
		tests.AssertGRPCError(t, status.New(codes.FailedPrecondition, `pmm-agent with ID "A1" has agents.`), err)

		expected := &models.Agent{
			AgentID:      "A1",
			AgentType:    models.PMMAgentType,
			RunsOnNodeID: pointer.ToString("N1"),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		agent, err = models.RemoveAgent(q, "A1", models.RemoveCascade)
		assert.Equal(t, expected, agent)
		assert.NoError(t, err)
		_, err = models.AgentFindByID(q, "A1")
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Agent with ID "A1" not found.`), err)
	})

	t.Run("FindPMMAgentsForNode", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agents, err := models.FindPMMAgentsForNode(q, "N1")
		require.NoError(t, err)
		assert.Equal(t, "A1", agents[0].AgentID)

		// find with non existing node.
		_, err = models.FindPMMAgentsForNode(q, "X1")
		require.Error(t, err)
	})

	t.Run("FindPMMAgentsForService", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agents, err := models.FindPMMAgentsForService(q, "S1")
		require.NoError(t, err)
		t.Log(agents, err)
		assert.Equal(t, "A1", agents[0].AgentID)

		// find with non existing service.
		_, err = models.FindPMMAgentsForService(q, "X1")
		require.Error(t, err)
	})

	t.Run("FindAgentsByServiceIDAndAgentType", func(t *testing.T) {
		q, teardown := setup(t)
		defer teardown(t)

		agents, err := models.FindAgentsByServiceIDAndAgentType(q, "S1", models.MySQLdExporterType)
		require.NoError(t, err)
		t.Log(agents, err)
		assert.Equal(t, "A2", agents[0].AgentID)

		// find with no existing pmm-agent-id
		_, err = models.FindAgentsByServiceIDAndAgentType(q, "X1", models.MySQLdExporterType)
		require.Error(t, err)
	})
}
