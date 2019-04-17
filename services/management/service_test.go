package management

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm-managed/models"

	"github.com/google/uuid"
	"github.com/percona/pmm/api/managementpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestRemoveService(t *testing.T) {

	setup := func(t *testing.T) (ss *ServiceService, teardown func(t *testing.T)) {
		uuid.SetRand(new(tests.IDReader))
		sqlDB := tests.OpenTestDB(t)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		r := new(mockRegistry)
		r.Test(t)
		teardown = func(t *testing.T) {
			require.NoError(t, sqlDB.Close())
			r.AssertExpectations(t)
		}
		ss = NewServiceService(db, r)
		return ss, teardown
	}

	t.Run("No params", func(t *testing.T) {
		ctx := logger.Set(context.Background(), t.Name())
		ss, teardown := setup(t)
		defer teardown(t)

		response, err := ss.RemoveService(ctx, &managementpb.RemoveServiceRequest{})
		assert.EqualError(t, err, errNoParamsNotFound.Error())
		assert.Nil(t, response)
	})

	t.Run("Both params", func(t *testing.T) {
		ctx := logger.Set(context.Background(), t.Name())
		ss, teardown := setup(t)
		defer teardown(t)

		response, err := ss.RemoveService(ctx, &managementpb.RemoveServiceRequest{ServiceId: "some-id", ServiceName: "some-service-name"})
		assert.EqualError(t, err, errOneOfParamsExpected.Error())
		assert.Nil(t, response)
	})

	t.Run("Not found", func(t *testing.T) {
		ctx := logger.Set(context.Background(), t.Name())
		ss, teardown := setup(t)
		defer teardown(t)

		response, err := ss.RemoveService(ctx, &managementpb.RemoveServiceRequest{ServiceName: "some-service-name"})
		assert.EqualError(t, err, "rpc error: code = NotFound desc = Service with Name \"some-service-name\" not found.")
		assert.Nil(t, response)
	})

	t.Run("Basic", func(t *testing.T) {
		ctx := logger.Set(context.Background(), t.Name())
		ss, teardown := setup(t)
		defer teardown(t)

		service, err := models.AddNewService(ss.db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
			ServiceName: "test-mysql",
			NodeID:      models.PMMServerNodeID,
			Address:     pointer.ToString("127.0.0.1"),
			Port:        pointer.ToUint16(3306),
		})
		require.NoError(t, err)

		pmmAgent, err := models.AgentAddPmmAgent(ss.db.Querier, models.PMMServerNodeID, nil)
		require.NoError(t, err)
		defer models.AgentRemove(ss.db.Querier, pmmAgent.AgentID)

		mysqldExporter, err := models.AgentAddExporter(ss.db.Querier, models.MySQLdExporterType, &models.AddExporterAgentParams{
			PMMAgentID: pmmAgent.AgentID,
			ServiceID:  service.ServiceID,
			Password:   "password",
			Username:   "username",
		})
		require.NoError(t, err)
		ss.asrs.(*mockRegistry).On("SendSetStateRequest", ctx, pmmAgent.AgentID)

		response, err := ss.RemoveService(ctx, &managementpb.RemoveServiceRequest{ServiceName: service.ServiceName})
		assert.NoError(t, err)
		assert.NotNil(t, response)

		ss.asrs.(*mockRegistry).AssertCalled(t, "SendSetStateRequest", ctx, pmmAgent.AgentID)

		agent, err := models.AgentFindByID(ss.db.Querier, mysqldExporter.AgentID)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = Agent with ID \"/agent_id/00000000-0000-4000-8000-000000000003\" not found.")
		assert.Nil(t, agent)

		service, err = models.FindServiceByID(ss.db.Querier, service.ServiceID)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = Service with ID \"/service_id/00000000-0000-4000-8000-000000000001\" not found.")
		assert.Nil(t, service)
	})
}
