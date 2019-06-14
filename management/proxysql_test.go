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
	proxysql "github.com/percona/pmm/api/managementpb/json/client/proxy_sql"
	"github.com/percona/pmm/api/managementpb/json/client/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestAddProxySQL(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-for-basic-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-for-basic-name")

		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
				Username:    "username",

				SkipConnectionCheck: true,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		require.NoError(t, err)
		require.NotNil(t, addProxySQLOK)
		require.NotNil(t, addProxySQLOK.Payload.Service)
		serviceID := addProxySQLOK.Payload.Service.ServiceID
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
			Proxysql: &services.GetServiceOKBodyProxysql{
				ServiceID:   serviceID,
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
			},
		}, *serviceOK.Payload)

		// Check that proxysqld exporter is added by default.
		listAgents, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Context: pmmapitests.Context,
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, agents.ListAgentsOKBody{
			ProxysqlExporter: []*agents.ProxysqlExporterItems0{
				{
					AgentID:    listAgents.Payload.ProxysqlExporter[0].AgentID,
					ServiceID:  serviceID,
					PMMAgentID: pmmAgentID,
					Username:   "username",
				},
			},
		}, *listAgents.Payload)
		defer removeAllAgentsInList(t, listAgents)
	})

	t.Run("With agents", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-for-all-fields-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-for-all-fields-name")

		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
				Username:    "username",
				Password:    "password",

				SkipConnectionCheck: true,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		require.NoError(t, err)
		require.NotNil(t, addProxySQLOK)
		require.NotNil(t, addProxySQLOK.Payload.Service)
		serviceID := addProxySQLOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)

		// Check that service is created and its fields.
		serviceOK, err := inventoryClient.Default.Services.GetService(&services.GetServiceParams{
			Body: services.GetServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.NotNil(t, serviceOK)
		assert.Equal(t, services.GetServiceOKBody{
			Proxysql: &services.GetServiceOKBodyProxysql{
				ServiceID:   serviceID,
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
			},
		}, *serviceOK.Payload)

		// Check that exporters are added.
		listAgents, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Context: pmmapitests.Context,
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
		})
		assert.NoError(t, err)
		require.NotNil(t, listAgents)
		defer removeAllAgentsInList(t, listAgents)
		require.Len(t, listAgents.Payload.ProxysqlExporter, 1)
		assert.Equal(t, agents.ListAgentsOKBody{
			ProxysqlExporter: []*agents.ProxysqlExporterItems0{
				{
					AgentID:    listAgents.Payload.ProxysqlExporter[0].AgentID,
					ServiceID:  serviceID,
					PMMAgentID: pmmAgentID,
					Username:   "username",
					Password:   "password",
				},
			},
		}, *listAgents.Payload)
	})

	t.Run("With labels", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-for-all-fields-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-for-all-fields-name")

		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:         nodeID,
				PMMAgentID:     pmmAgentID,
				ServiceName:    serviceName,
				Address:        "10.10.10.10",
				Port:           3306,
				Username:       "username",
				Password:       "password",
				Environment:    "some-environment",
				Cluster:        "cluster-name",
				ReplicationSet: "replication-set",
				CustomLabels:   map[string]string{"bar": "foo"},

				SkipConnectionCheck: true,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		require.NoError(t, err)
		require.NotNil(t, addProxySQLOK)
		require.NotNil(t, addProxySQLOK.Payload.Service)
		serviceID := addProxySQLOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removeServiceAgents(t, serviceID)

		// Check that service is created and its fields.
		serviceOK, err := inventoryClient.Default.Services.GetService(&services.GetServiceParams{
			Body: services.GetServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.NotNil(t, serviceOK)
		assert.Equal(t, services.GetServiceOKBody{
			Proxysql: &services.GetServiceOKBodyProxysql{
				ServiceID:      serviceID,
				NodeID:         nodeID,
				ServiceName:    serviceName,
				Address:        "10.10.10.10",
				Port:           3306,
				Environment:    "some-environment",
				Cluster:        "cluster-name",
				ReplicationSet: "replication-set",
				CustomLabels:   map[string]string{"bar": "foo"},
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

		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
				Username:    "username",

				SkipConnectionCheck: true,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		require.NoError(t, err)
		require.NotNil(t, addProxySQLOK)
		require.NotNil(t, addProxySQLOK.Payload.Service)
		serviceID := addProxySQLOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removeServiceAgents(t, serviceID)

		params = &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "11.11.11.11",
				Port:        3307,
				Username:    "username",
			},
		}
		addProxySQLOK, err = client.Default.ProxySQL.AddProxySQL(params)
		require.Nil(t, addProxySQLOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{409, fmt.Sprintf(`Service with name "%s" already exists.`, serviceName)})
	})

	t.Run("Empty Node ID", func(t *testing.T) {
		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body:    proxysql.AddProxySQLBody{},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field NodeId: value '' must not be an empty string"})
		assert.Nil(t, addProxySQLOK)
	})

	t.Run("Empty Service Name", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body:    proxysql.AddProxySQLBody{NodeID: nodeID},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field ServiceName: value '' must not be an empty string"})
		assert.Nil(t, addProxySQLOK)
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
		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field Address: value '' must not be an empty string"})
		assert.Nil(t, addProxySQLOK)
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
		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field Port: value '0' must be greater than '0'"})
		assert.Nil(t, addProxySQLOK)
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
		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field PmmAgentId: value '' must not be an empty string"})
		assert.Nil(t, addProxySQLOK)
	})

	t.Run("Empty username", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-name")
		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
				PMMAgentID:  pmmAgentID,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field Username: value '' must not be an empty string"})
		assert.Nil(t, addProxySQLOK)
	})
}

func TestRemoveProxySQL(t *testing.T) {
	addProxySQL := func(t *testing.T, serviceName, nodeName string) (nodeID string, pmmAgentID string, serviceID string) {
		t.Helper()
		nodeID, pmmAgentID = registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})

		params := &proxysql.AddProxySQLParams{
			Context: pmmapitests.Context,
			Body: proxysql.AddProxySQLBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
				Username:    "username",
				Password:    "password",

				SkipConnectionCheck: true,
			},
		}
		addProxySQLOK, err := client.Default.ProxySQL.AddProxySQL(params)
		require.NoError(t, err)
		require.NotNil(t, addProxySQLOK)
		require.NotNil(t, addProxySQLOK.Payload.Service)
		serviceID = addProxySQLOK.Payload.Service.ServiceID
		return
	}

	t.Run("By name", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-by-name")
		nodeName := pmmapitests.TestString(t, "node-remove-by-name")
		nodeID, pmmAgentID, serviceID := addProxySQL(t, serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceName: serviceName,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePROXYSQLSERVICE),
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
		nodeID, pmmAgentID, serviceID := addProxySQL(t, serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePROXYSQLSERVICE),
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
		nodeID, pmmAgentID, serviceID := addProxySQL(t, serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceName: serviceName,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePROXYSQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		assert.Nil(t, removeServiceOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "service_id or service_name expected; not both"})
	})

	t.Run("Wrong type", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-wrong-type")
		nodeName := pmmapitests.TestString(t, "node-remove-wrong-type")
		nodeID, pmmAgentID, serviceID := addProxySQL(t, serviceName, nodeName)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypePOSTGRESQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		assert.Nil(t, removeServiceOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "wrong service type"})
	})

	t.Run("No params", func(t *testing.T) {
		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body:    service.RemoveServiceBody{},
			Context: pmmapitests.Context,
		})
		assert.Nil(t, removeServiceOK)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "params not found"})
	})
}
