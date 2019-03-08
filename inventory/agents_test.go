package inventory

import (
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/percona/pmm/api/inventory/json/client/agents"
	"github.com/percona/pmm/api/inventory/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Percona-Lab/pmm-api-tests"
)

func TestAgents(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, withUUID(t, "Remote node for agents list"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: withUUID(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		mySqldExporter := addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:    serviceID,
			Username:     "username",
			Password:     "password",
			RunsOnNodeID: genericNodeID,
		})
		mySqldExporterID := mySqldExporter.MysqldExporter.AgentID
		defer removeAgents(t, mySqldExporterID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{Context: pmmapitests.Context})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotZerof(t, len(res.Payload.MysqldExporter), "There should be at least one service")

		assertMySQLExporterExists(t, res, mySqldExporterID)
		assertPMMAgentExists(t, res, pmmAgentID)
	})

	t.Run("FilterList", func(t *testing.T) {
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, withUUID(t, "Remote node for agents filters"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: withUUID(t, "MySQL Service for filter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		mySqldExporter := addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:    serviceID,
			Username:     "username",
			Password:     "password",
			RunsOnNodeID: genericNodeID,
		})
		mySqldExporterID := mySqldExporter.MysqldExporter.AgentID
		defer removeAgents(t, mySqldExporterID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		// Filter by runs on node ID.
		res, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body:    agents.ListAgentsBody{RunsOnNodeID: genericNodeID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotZerof(t, len(res.Payload.MysqldExporter), "There should be at least one service")
		assertMySQLExporterExists(t, res, mySqldExporterID)
		assertPMMAgentNotExists(t, res, pmmAgentID)

		// Filter by node ID.
		res, err = client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body:    agents.ListAgentsBody{NodeID: nodeID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotZerof(t, len(res.Payload.PMMAgent), "There should be at least one service")
		assertMySQLExporterNotExists(t, res, mySqldExporterID)
		assertPMMAgentExists(t, res, pmmAgentID)

		// Filter by service ID.
		res, err = client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body:    agents.ListAgentsBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotZerof(t, len(res.Payload.MysqldExporter), "There should be at least one service")
		assertMySQLExporterExists(t, res, mySqldExporterID)
		assertPMMAgentNotExists(t, res, pmmAgentID)
	})

	t.Run("TwoOrMoreFilters", func(t *testing.T) {
		t.Skip("it doesn't return error")
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		res, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body: agents.ListAgentsBody{
				RunsOnNodeID: genericNodeID,
				NodeID:       genericNodeID,
				ServiceID:    "some-service-id",
			},
			Context: pmmapitests.Context,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestPMMAgent(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, withUUID(t, "Remote node for PMM-agent"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		res := addPMMAgent(t, nodeID)
		require.Equal(t, nodeID, res.PMMAgent.NodeID)
		agentID := res.PMMAgent.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				PMMAgent: &agents.GetAgentOKBodyPMMAgent{
					AgentID: agentID,
					NodeID:  nodeID,
				},
			},
		}, getAgentRes)
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddPMMAgent(&agents.AddPMMAgentParams{
			Body:    agents.AddPMMAgentBody{NodeID: ""},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})
}

func TestNodeExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, withUUID(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body: agents.AddNodeExporterBody{
				NodeID: nodeID,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Payload.NodeExporter)
		require.Equal(t, nodeID, res.Payload.NodeExporter.NodeID)
		agentID := res.Payload.NodeExporter.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				NodeExporter: &agents.GetAgentOKBodyNodeExporter{
					AgentID: agentID,
					NodeID:  nodeID,
				},
			},
		}, getAgentRes)
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body:    agents.AddNodeExporterBody{NodeID: ""},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})

	t.Run("NotExistNodeID", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body:    agents.AddNodeExporterBody{NodeID: "pmm-node-exporter-node"},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 404)
		assert.Nil(t, res)
	})
}

func TestMySQLdExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, withUUID(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: withUUID(t, "MySQL Service for MySQLdExporter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		mySqldExporter := addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:    serviceID,
			Username:     "username",
			Password:     "password",
			RunsOnNodeID: nodeID,
		})
		agentID := mySqldExporter.MysqldExporter.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				MysqldExporter: &agents.GetAgentOKBodyMysqldExporter{
					AgentID:      agentID,
					ServiceID:    serviceID,
					Username:     "username",
					Password:     "password",
					RunsOnNodeID: nodeID,
				},
			},
		}, getAgentRes)
	})

	t.Run("AddServiceIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:    "",
				RunsOnNodeID: genericNodeID,
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})

	t.Run("AddRunsOnNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:    "pmm-service-id",
				RunsOnNodeID: "",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})

	t.Run("NotExistServiceID", func(t *testing.T) {
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:    "pmm-service-id",
				RunsOnNodeID: genericNodeID,
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 404)
		assert.Nil(t, res)
	})

	t.Run("NotExistNodeID", func(t *testing.T) {
		t.Parallel()

		genericNode := addGenericNode(t, withUUID(t, "Test Remote Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: withUUID(t, "MySQL Service for not exists node ID"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:    serviceID,
				RunsOnNodeID: "pmm-not-exist-server",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, 404)
		assert.Nil(t, res)
	})
}

func TestRDSExporter(t *testing.T) {
	t.Skip("Not implemented yet.")
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, withUUID(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      nodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: withUUID(t, "MySQL Service for RDSExporter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddRDSExporter(&agents.AddRDSExporterParams{
			Body: agents.AddRDSExporterBody{
				RunsOnNodeID: nodeID,
				ServiceIds:   []string{serviceID},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Payload.RDSExporter)
		agentID := res.Payload.RDSExporter.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				RDSExporter: &agents.GetAgentOKBodyRDSExporter{
					AgentID:      agentID,
					RunsOnNodeID: nodeID,
					ServiceIds:   []string{serviceID},
				},
			},
		}, getAgentRes)
	})
}
