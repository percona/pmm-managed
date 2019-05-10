package management

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	inventoryClient "github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/node"
	postgresql "github.com/percona/pmm/api/managementpb/json/client/postgre_sql"
	"github.com/percona/pmm/api/managementpb/json/client/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestAddPostgreSQL(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-for-basic-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-for-basic-name")

		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        5432,
				Username:    "username",
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		require.NoError(t, err)
		require.NotNil(t, addPostgreSQLOK)
		require.NotNil(t, addPostgreSQLOK.Payload.Service)
		serviceID := addPostgreSQLOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)

		// Check that service is created and its fields.
		serviceOK, err := inventoryClient.Default.Services.GetService(&services.GetServiceParams{
			Body: services.GetServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		require.NotNil(t, serviceOK)
		assert.Equal(t, services.GetServiceOKBody{
			Postgresql: &services.GetServiceOKBodyPostgresql{
				ServiceID:   serviceID,
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        5432,
			},
		}, *serviceOK.Payload)

		// Check that no one exporter is added.
		listAgents, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Context: pmmapitests.Context,
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, agents.ListAgentsOKBody{
			PostgresExporter: []*agents.PostgresExporterItems0{
				{
					AgentID:    listAgents.Payload.PostgresExporter[0].AgentID,
					ServiceID:  serviceID,
					PMMAgentID: pmmAgentID,
					Username:   "username",
				},
			},
		}, *listAgents.Payload)
		defer removeAllAgentsInList(t, listAgents)
	})

	t.Run("With labels", func(realT *testing.T) {
		expectedFailureTestingT := pmmapitests.ExpectFailure(realT, "https://jira.percona.com/browse/PMM-3982")
		defer expectedFailureTestingT.Check()

		nodeName := pmmapitests.TestString(realT, "node-for-all-fields-name")
		nodeID, pmmAgentID := registerGenericNode(realT, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(realT, nodeID)
		defer removePMMAgentWithSubAgents(realT, pmmAgentID)

		serviceName := pmmapitests.TestString(realT, "service-for-all-fields-name")

		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:       nodeID,
				PMMAgentID:   pmmAgentID,
				ServiceName:  serviceName,
				Address:      "10.10.10.10",
				Port:         5432,
				Username:     "username",
				Environment:  "some-environment",
				CustomLabels: map[string]string{"bar": "foo"},
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		require.NoError(realT, err)
		require.NotNil(realT, addPostgreSQLOK)
		require.NotNil(realT, addPostgreSQLOK.Payload.Service)
		serviceID := addPostgreSQLOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(realT, serviceID)
		defer removeServiceAgents(realT, serviceID)

		// Check that service is created and its fields.
		serviceOK, err := inventoryClient.Default.Services.GetService(&services.GetServiceParams{
			Body: services.GetServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(realT, err)
		assert.NotNil(realT, serviceOK)
		assert.Equal(expectedFailureTestingT, services.GetServiceOKBody{
			Postgresql: &services.GetServiceOKBodyPostgresql{
				ServiceID:    serviceID,
				NodeID:       nodeID,
				ServiceName:  serviceName,
				Address:      "10.10.10.10",
				Port:         5432,
				Environment:  "some-environment",
				CustomLabels: map[string]string{"bar": "foo"},
			},
		}, *serviceOK.Payload)
	})

	t.Run("With the same name", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-for-the-same-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-for-the-same-name")

		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Username:    "username",
				Address:     "10.10.10.10",
				Port:        5432,
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		require.NoError(t, err)
		require.NotNil(t, addPostgreSQLOK)
		require.NotNil(t, addPostgreSQLOK.Payload.Service)
		serviceID := addPostgreSQLOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removeServiceAgents(t, serviceID)

		params = &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Username:    "username",
				Address:     "11.11.11.11",
				Port:        5433,
			},
		}
		addPostgreSQLOK, err = client.Default.PostgreSQL.AddPostgreSQL(params)
		require.Nil(t, addPostgreSQLOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{409, fmt.Sprintf(`Service with name "%s" already exists.`, serviceName)})
	})

	t.Run("Empty Node ID", func(t *testing.T) {
		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body:    postgresql.AddPostgreSQLBody{},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field NodeId: value '' must not be an empty string"})
		assert.Nil(t, addPostgreSQLOK)
	})

	t.Run("Empty Service Name", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body:    postgresql.AddPostgreSQLBody{NodeID: nodeID},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field ServiceName: value '' must not be an empty string"})
		assert.Nil(t, addPostgreSQLOK)
	})

	t.Run("Empty Address", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-name")
		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field Address: value '' must not be an empty string"})
		assert.Nil(t, addPostgreSQLOK)
	})

	t.Run("Empty Port", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-name")
		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field Port: value '0' must be greater than '0'"})
		assert.Nil(t, addPostgreSQLOK)
	})

	t.Run("Empty Pmm Agent ID", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-name")
		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        5432,
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field PmmAgentId: value '' must not be an empty string"})
		assert.Nil(t, addPostgreSQLOK)
	})
}

func TestRemovePostgreSQL(t *testing.T) {
	addPostgreSQL := func(serviceName, nodeName string) (nodeID string, pmmAgentID string, serviceID string) {
		nodeID, pmmAgentID = registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})

		params := &postgresql.AddPostgreSQLParams{
			Context: pmmapitests.Context,
			Body: postgresql.AddPostgreSQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        5432,
				Username:    "username",
				Password:    "password",
			},
		}
		addPostgreSQLOK, err := client.Default.PostgreSQL.AddPostgreSQL(params)
		require.NoError(t, err)
		require.NotNil(t, addPostgreSQLOK)
		require.NotNil(t, addPostgreSQLOK.Payload.Service)
		serviceID = addPostgreSQLOK.Payload.Service.ServiceID
		return
	}

	t.Run("By name", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-by-name")
		nodeName := pmmapitests.TestString(t, "node-remove-by-name")
		nodeID, pmmAgentID, serviceID := addPostgreSQL(serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceName: serviceName,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePOSTGRESQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		noError := assert.NoError(t, err)
		notNil := assert.NotNil(t, removeServiceOK)
		if !noError || !notNil {
			defer pmmapitests.RemoveServices(t, serviceID)
		}

		// Check that the service removed with agents.
		listAgents, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Context: pmmapitests.Context,
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, agents.ListAgentsOKBody{}, *listAgents.Payload)
		defer removeAllAgentsInList(t, listAgents)
	})

	t.Run("By ID", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-by-id")
		nodeName := pmmapitests.TestString(t, "node-remove-by-id")
		nodeID, pmmAgentID, serviceID := addPostgreSQL(serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePOSTGRESQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		noError := assert.NoError(t, err)
		notNil := assert.NotNil(t, removeServiceOK)
		if !noError || !notNil {
			defer pmmapitests.RemoveServices(t, serviceID)
		}
		// Check that the service removed with agents.
		listAgents, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Context: pmmapitests.Context,
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, agents.ListAgentsOKBody{}, *listAgents.Payload)
		defer removeAllAgentsInList(t, listAgents)
	})

	t.Run("Both params", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-both-params")
		nodeName := pmmapitests.TestString(t, "node-remove-both-params")
		nodeID, pmmAgentID, serviceID := addPostgreSQL(serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceName: serviceName,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePOSTGRESQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		assert.Nil(t, removeServiceOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "service_id or service_name expected; not both"})
	})

	t.Run("Wrong type", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-wrong-type")
		nodeName := pmmapitests.TestString(t, "node-remove-wrong-type")
		nodeID, pmmAgentID, serviceID := addPostgreSQL(serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypeMYSQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		assert.Nil(t, removeServiceOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "wrong service type"})
	})
}
