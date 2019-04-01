package inventory

import (
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/percona/pmm/api/inventory/json/client/agents"
	"github.com/percona/pmm/api/inventory/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestAgents(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Generic node for agents list")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for agents list"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		mySqldExporter := addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,
		})
		mySqldExporterID := mySqldExporter.MysqldExporter.AgentID
		defer removeAgents(t, mySqldExporterID)

		res, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{Context: pmmapitests.Context})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotZerof(t, len(res.Payload.MysqldExporter), "There should be at least one service")

		assertMySQLExporterExists(t, res, mySqldExporterID)
		assertPMMAgentExists(t, res, pmmAgentID)
	})

	t.Run("FilterList", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Generic node for agents filters")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for agents filters"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for filter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		mySqldExporter := addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,
		})
		mySqldExporterID := mySqldExporter.MysqldExporter.AgentID
		defer removeAgents(t, mySqldExporterID)

		// Filter by pmm agent ID.
		res, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body:    agents.ListAgentsBody{PMMAgentID: pmmAgentID},
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
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body: agents.ListAgentsBody{
				PMMAgentID: pmmAgentID,
				NodeID:     genericNodeID,
				ServiceID:  "some-service-id",
			},
			Context: pmmapitests.Context,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestPMMAgent(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for PMM-agent"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		res := addPMMAgent(t, nodeID)
		require.Equal(t, nodeID, res.PMMAgent.RunsOnNodeID)
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
					AgentID:      agentID,
					RunsOnNodeID: nodeID,
				},
			},
		}, getAgentRes)
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddPMMAgent(&agents.AddPMMAgentParams{
			Body:    agents.AddPMMAgentBody{RunsOnNodeID: ""},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field RunsOnNodeId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.PMMAgent.AgentID)
		}
	})
}

func TestNodeExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body: agents.AddNodeExporterBody{
				PMMAgentID: pmmAgentID,
				CustomLabels: map[string]string{
					"custom_label_node_exporter": "node_exporter",
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Payload.NodeExporter)
		require.Equal(t, pmmAgentID, res.Payload.NodeExporter.PMMAgentID)
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
					AgentID:    agentID,
					PMMAgentID: pmmAgentID,
					Disabled:   false,
					CustomLabels: map[string]string{
						"custom_label_node_exporter": "node_exporter",
					},
				},
			},
		}, getAgentRes)

		// Test change API.
		changeNodeExporterOK, err := client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID:            agentID,
				Disabled:           true,
				RemoveCustomLabels: true,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeNodeExporterOK{
			Payload: &agents.ChangeNodeExporterOKBody{
				NodeExporter: &agents.ChangeNodeExporterOKBodyNodeExporter{
					AgentID:    agentID,
					PMMAgentID: pmmAgentID,
					Disabled:   true,
				},
			},
		}, changeNodeExporterOK)

		changeNodeExporterOK, err = client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID: agentID,
				Enabled: true,
				CustomLabels: map[string]string{
					"new_label": "node_exporter",
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeNodeExporterOK{
			Payload: &agents.ChangeNodeExporterOKBody{
				NodeExporter: &agents.ChangeNodeExporterOKBodyNodeExporter{
					AgentID:    agentID,
					PMMAgentID: pmmAgentID,
					Disabled:   false,
					CustomLabels: map[string]string{
						"new_label": "node_exporter",
					},
				},
			},
		}, changeNodeExporterOK)
	})

	t.Run("AddPMMAgentIDEmpty", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body:    agents.AddNodeExporterBody{PMMAgentID: ""},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field PmmAgentId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.NodeExporter.AgentID)
		}
	})

	t.Run("NotExistPmmAgentID", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body:    agents.AddNodeExporterBody{PMMAgentID: "pmm-node-exporter-node"},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, "Agent with ID \"pmm-node-exporter-node\" not found."})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.NodeExporter.AgentID)
		}
	})
}

func TestMySQLdExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for MySQLdExporter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		mySqldExporter := addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,
			CustomLabels: map[string]string{
				"custom_label_mysql_exporter": "mysql_exporter",
			},
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
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					CustomLabels: map[string]string{
						"custom_label_mysql_exporter": "mysql_exporter",
					},
				},
			},
		}, getAgentRes)

		// Test change API.
		changeMySQLdExporterOK, err := client.Default.Agents.ChangeMySqldExporter(&agents.ChangeMySqldExporterParams{
			Body: agents.ChangeMySqldExporterBody{
				AgentID:            agentID,
				Disabled:           true,
				RemoveCustomLabels: true,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeMySqldExporterOK{
			Payload: &agents.ChangeMySqldExporterOKBody{
				MysqldExporter: &agents.ChangeMySqldExporterOKBodyMysqldExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   true,
				},
			},
		}, changeMySQLdExporterOK)

		changeMySQLdExporterOK, err = client.Default.Agents.ChangeMySqldExporter(&agents.ChangeMySqldExporterParams{
			Body: agents.ChangeMySqldExporterBody{
				AgentID: agentID,
				Enabled: true,
				CustomLabels: map[string]string{
					"new_label": "mysql_exporter",
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeMySqldExporterOK{
			Payload: &agents.ChangeMySqldExporterOKBody{
				MysqldExporter: &agents.ChangeMySqldExporterOKBodyMysqldExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   false,
					CustomLabels: map[string]string{
						"new_label": "mysql_exporter",
					},
				},
			},
		}, changeMySQLdExporterOK)
	})

	t.Run("AddServiceIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:  "",
				PMMAgentID: pmmAgentID,
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.MysqldExporter.AgentID)
		}
	})

	t.Run("AddPMMAgentIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:  serviceID,
				PMMAgentID: "",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field PmmAgentId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MysqldExporter.AgentID)
		}
	})

	t.Run("NotExistServiceID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:  "pmm-service-id",
				PMMAgentID: pmmAgentID,
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, "Service with ID \"pmm-service-id\" not found."})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MysqldExporter.AgentID)
		}
	})

	t.Run("NotExistPMMAgentID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for not exists node ID"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
			Body: agents.AddMySqldExporterBody{
				ServiceID:  serviceID,
				PMMAgentID: "pmm-not-exist-server",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, ""})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MysqldExporter.AgentID)
		}
	})
}

func TestRDSExporter(t *testing.T) {
	t.Skip("Not implemented yet.")

	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      nodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for RDSExporter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddRDSExporter(&agents.AddRDSExporterParams{
			Body: agents.AddRDSExporterBody{
				PMMAgentID: pmmAgentID,
				ServiceIds: []string{serviceID},
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
					AgentID:    agentID,
					PMMAgentID: pmmAgentID,
					ServiceIds: []string{serviceID},
				},
			},
		}, getAgentRes)
	})
}

func TestMongoDBExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for MongoDBExporter test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		mongoDBExporter := addMongoDBExporter(t, agents.AddMongoDBExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,
			CustomLabels: map[string]string{
				"new_label": "mongodb_exporter",
			},
		})
		agentID := mongoDBExporter.MongodbExporter.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				MongodbExporter: &agents.GetAgentOKBodyMongodbExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					CustomLabels: map[string]string{
						"new_label": "mongodb_exporter",
					},
				},
			},
		}, getAgentRes)

		// Test change API.
		changeMongoDBExporterOK, err := client.Default.Agents.ChangeMongoDBExporter(&agents.ChangeMongoDBExporterParams{
			Body: agents.ChangeMongoDBExporterBody{
				AgentID:            agentID,
				Disabled:           true,
				RemoveCustomLabels: true,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeMongoDBExporterOK{
			Payload: &agents.ChangeMongoDBExporterOKBody{
				MongodbExporter: &agents.ChangeMongoDBExporterOKBodyMongodbExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   true,
				},
			},
		}, changeMongoDBExporterOK)

		changeMongoDBExporterOK, err = client.Default.Agents.ChangeMongoDBExporter(&agents.ChangeMongoDBExporterParams{
			Body: agents.ChangeMongoDBExporterBody{
				AgentID: agentID,
				Enabled: true,
				CustomLabels: map[string]string{
					"new_label": "mongodb_exporter",
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeMongoDBExporterOK{
			Payload: &agents.ChangeMongoDBExporterOKBody{
				MongodbExporter: &agents.ChangeMongoDBExporterOKBodyMongodbExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   false,
					CustomLabels: map[string]string{
						"new_label": "mongodb_exporter",
					},
				},
			},
		}, changeMongoDBExporterOK)
	})

	t.Run("AddServiceIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddMongoDBExporter(&agents.AddMongoDBExporterParams{
			Body: agents.AddMongoDBExporterBody{
				ServiceID:  "",
				PMMAgentID: pmmAgentID,
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MongodbExporter.AgentID)
		}
	})

	t.Run("AddPMMAgentIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddMongoDBExporter(&agents.AddMongoDBExporterParams{
			Body: agents.AddMongoDBExporterBody{
				ServiceID:  serviceID,
				PMMAgentID: "",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field PmmAgentId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MongodbExporter.AgentID)
		}
	})

	t.Run("NotExistServiceID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddMongoDBExporter(&agents.AddMongoDBExporterParams{
			Body: agents.AddMongoDBExporterBody{
				ServiceID:  "pmm-service-id",
				PMMAgentID: pmmAgentID,
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, "Service with ID \"pmm-service-id\" not found."})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MongodbExporter.AgentID)
		}
	})

	t.Run("NotExistPMMAgentID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for not exists node ID"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddMongoDBExporter(&agents.AddMongoDBExporterParams{
			Body: agents.AddMongoDBExporterBody{
				ServiceID:  serviceID,
				PMMAgentID: "pmm-not-exist-server",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, ""})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.MongodbExporter.AgentID)
		}
	})
}

func TestQanAgentExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Test Generic Node for Qan Agent")).NodeID
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for QanAgent test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddQANMySQLPerfSchemaAgent(
			&agents.AddQANMySQLPerfSchemaAgentParams{
				Body: agents.AddQANMySQLPerfSchemaAgentBody{
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					CustomLabels: map[string]string{
						"new_label": "QANMysqlPerfschemaAgent",
					},
				},
				Context: pmmapitests.Context,
			})
		agentID := res.Payload.QANMysqlPerfschemaAgent.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				QANMysqlPerfschemaAgent: &agents.GetAgentOKBodyQANMysqlPerfschemaAgent{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					CustomLabels: map[string]string{
						"new_label": "QANMysqlPerfschemaAgent",
					},
				},
			},
		}, getAgentRes)

		// Test change API.
		changeQANMySQLPerfSchemaAgentOK, err := client.Default.Agents.ChangeQANMySQLPerfSchemaAgent(&agents.ChangeQANMySQLPerfSchemaAgentParams{
			Body: agents.ChangeQANMySQLPerfSchemaAgentBody{
				AgentID:            agentID,
				Disabled:           true,
				RemoveCustomLabels: true,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeQANMySQLPerfSchemaAgentOK{
			Payload: &agents.ChangeQANMySQLPerfSchemaAgentOKBody{
				QANMysqlPerfschemaAgent: &agents.ChangeQANMySQLPerfSchemaAgentOKBodyQANMysqlPerfschemaAgent{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   true,
				},
			},
		}, changeQANMySQLPerfSchemaAgentOK)

		changeQANMySQLPerfSchemaAgentOK, err = client.Default.Agents.ChangeQANMySQLPerfSchemaAgent(&agents.ChangeQANMySQLPerfSchemaAgentParams{
			Body: agents.ChangeQANMySQLPerfSchemaAgentBody{
				AgentID: agentID,
				Enabled: true,
				CustomLabels: map[string]string{
					"new_label": "QANMysqlPerfschemaAgent",
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeQANMySQLPerfSchemaAgentOK{
			Payload: &agents.ChangeQANMySQLPerfSchemaAgentOKBody{
				QANMysqlPerfschemaAgent: &agents.ChangeQANMySQLPerfSchemaAgentOKBodyQANMysqlPerfschemaAgent{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   false,
					CustomLabels: map[string]string{
						"new_label": "QANMysqlPerfschemaAgent",
					},
				},
			},
		}, changeQANMySQLPerfSchemaAgentOK)
	})

	t.Run("AddServiceIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Test Generic Node for Qan Agent")).NodeID
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddQANMySQLPerfSchemaAgent(&agents.AddQANMySQLPerfSchemaAgentParams{
			Body: agents.AddQANMySQLPerfSchemaAgentBody{
				ServiceID:  "",
				PMMAgentID: pmmAgentID,
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.QANMysqlPerfschemaAgent.AgentID)
		}
	})

	t.Run("AddPMMAgentIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Test Generic Node for Qan Agent")).NodeID
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddQANMySQLPerfSchemaAgent(&agents.AddQANMySQLPerfSchemaAgentParams{
			Body: agents.AddQANMySQLPerfSchemaAgentBody{
				ServiceID:  serviceID,
				PMMAgentID: "",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field PmmAgentId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.QANMysqlPerfschemaAgent.AgentID)
		}
	})

	t.Run("NotExistServiceID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Test Generic Node for Qan Agent")).NodeID
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddQANMySQLPerfSchemaAgent(&agents.AddQANMySQLPerfSchemaAgentParams{
			Body: agents.AddQANMySQLPerfSchemaAgentBody{
				ServiceID:  "pmm-service-id",
				PMMAgentID: pmmAgentID,
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, "Service with ID \"pmm-service-id\" not found."})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.QANMysqlPerfschemaAgent.AgentID)
		}
	})

	t.Run("NotExistPMMAgentID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Test Generic Node for Qan Agent")).NodeID
		defer removeNodes(t, genericNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for not exists node ID"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddQANMySQLPerfSchemaAgent(&agents.AddQANMySQLPerfSchemaAgentParams{
			Body: agents.AddQANMySQLPerfSchemaAgentBody{
				ServiceID:  serviceID,
				PMMAgentID: "pmm-not-exist-server",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, ""})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.QANMysqlPerfschemaAgent.AgentID)
		}
	})
}

func TestPostgresExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer removeNodes(t, nodeID)

		service := addPostgreSQLService(t, services.AddPostgreSQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        5432,
			ServiceName: pmmapitests.TestString(t, "PostgreSQL Service for PostgresExporter test"),
		})
		serviceID := service.Postgresql.ServiceID
		defer removeServices(t, serviceID)

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		PostgresExporter := addPostgresExporter(t, agents.AddPostgresExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,
			CustomLabels: map[string]string{
				"custom_label_postgres_exporter": "postgres_exporter",
			},
		})
		agentID := PostgresExporter.PostgresExporter.AgentID
		defer removeAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				PostgresExporter: &agents.GetAgentOKBodyPostgresExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					CustomLabels: map[string]string{
						"custom_label_postgres_exporter": "postgres_exporter",
					},
				},
			},
		}, getAgentRes)

		// Test change API.
		changePostgresExporterOK, err := client.Default.Agents.ChangePostgresExporter(&agents.ChangePostgresExporterParams{
			Body: agents.ChangePostgresExporterBody{
				AgentID:            agentID,
				Disabled:           true,
				RemoveCustomLabels: true,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangePostgresExporterOK{
			Payload: &agents.ChangePostgresExporterOKBody{
				PostgresExporter: &agents.ChangePostgresExporterOKBodyPostgresExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   true,
				},
			},
		}, changePostgresExporterOK)

		changePostgresExporterOK, err = client.Default.Agents.ChangePostgresExporter(&agents.ChangePostgresExporterParams{
			Body: agents.ChangePostgresExporterBody{
				AgentID: agentID,
				Enabled: true,
				CustomLabels: map[string]string{
					"new_label": "postgres_exporter",
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangePostgresExporterOK{
			Payload: &agents.ChangePostgresExporterOKBody{
				PostgresExporter: &agents.ChangePostgresExporterOKBodyPostgresExporter{
					AgentID:    agentID,
					ServiceID:  serviceID,
					Username:   "username",
					Password:   "password",
					PMMAgentID: pmmAgentID,
					Disabled:   false,
					CustomLabels: map[string]string{
						"new_label": "postgres_exporter",
					},
				},
			},
		}, changePostgresExporterOK)
	})

	t.Run("AddServiceIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddPostgresExporter(&agents.AddPostgresExporterParams{
			Body: agents.AddPostgresExporterBody{
				ServiceID:  "",
				PMMAgentID: pmmAgentID,
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.PostgresExporter.AgentID)
		}
	})

	t.Run("AddPMMAgentIDEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		service := addPostgreSQLService(t, services.AddPostgreSQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        5432,
			ServiceName: pmmapitests.TestString(t, "PostgreSQL Service for agent"),
		})
		serviceID := service.Postgresql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddPostgresExporter(&agents.AddPostgresExporterParams{
			Body: agents.AddPostgresExporterBody{
				ServiceID:  serviceID,
				PMMAgentID: "",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field PmmAgentId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.PostgresExporter.AgentID)
		}
	})

	t.Run("NotExistServiceID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		pmmAgent := addPMMAgent(t, genericNodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer removeAgents(t, pmmAgentID)

		res, err := client.Default.Agents.AddPostgresExporter(&agents.AddPostgresExporterParams{
			Body: agents.AddPostgresExporterBody{
				ServiceID:  "pmm-service-id",
				PMMAgentID: pmmAgentID,
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, "Service with ID \"pmm-service-id\" not found."})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.PostgresExporter.AgentID)
		}
	})

	t.Run("NotExistPMMAgentID", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		service := addPostgreSQLService(t, services.AddPostgreSQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        5432,
			ServiceName: pmmapitests.TestString(t, "PostgreSQL Service for not exists node ID"),
		})
		serviceID := service.Postgresql.ServiceID
		defer removeServices(t, serviceID)

		res, err := client.Default.Agents.AddPostgresExporter(&agents.AddPostgresExporterParams{
			Body: agents.AddPostgresExporterBody{
				ServiceID:  serviceID,
				PMMAgentID: "pmm-not-exist-server",
				Username:   "username",
				Password:   "password",
			},
			Context: pmmapitests.Context,
		})
		assertEqualAPIError(t, err, ServerResponse{404, ""})
		if !assert.Nil(t, res) {
			removeAgents(t, res.Payload.PostgresExporter.AgentID)
		}
	})
}
