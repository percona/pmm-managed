package management

import (
	"testing"

	"github.com/AlekSi/pointer"
	inventoryClient "github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/percona/pmm/api/managementpb/json/client"
	mongodb "github.com/percona/pmm/api/managementpb/json/client/mongo_db"
	"github.com/percona/pmm/api/managementpb/json/client/node"
	"github.com/percona/pmm/api/managementpb/json/client/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestAddMongoDB(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-for-basic-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-name-for-basic-name")

		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        27017,

				SkipConnectionCheck: true,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		require.NoError(t, err)
		require.NotNil(t, addMongoDBOK)
		require.NotNil(t, addMongoDBOK.Payload.Service)
		serviceID := addMongoDBOK.Payload.Service.ServiceID
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
			Mongodb: &services.GetServiceOKBodyMongodb{
				ServiceID:   serviceID,
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        27017,
			},
		}, *serviceOK.Payload)

		// Check that mongodb exporter is added by default.
		listAgents, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Context: pmmapitests.Context,
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, agents.ListAgentsOKBody{
			MongodbExporter: []*agents.MongodbExporterItems0{
				{
					AgentID:    listAgents.Payload.MongodbExporter[0].AgentID,
					ServiceID:  serviceID,
					PMMAgentID: pmmAgentID,
				},
			},
		}, *listAgents.Payload)
		defer removeAllAgentsInList(t, listAgents)
	})

	t.Run("With agents", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name-for-all-fields")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		serviceName := pmmapitests.TestString(t, "service-name-for-all-fields")

		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:             nodeID,
				PMMAgentID:         pmmAgentID,
				ServiceName:        serviceName,
				Address:            "10.10.10.10",
				Port:               27017,
				Username:           "username",
				Password:           "password",
				QANMongodbProfiler: true,

				SkipConnectionCheck: true,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		require.NoError(t, err)
		require.NotNil(t, addMongoDBOK)
		require.NotNil(t, addMongoDBOK.Payload.Service)
		serviceID := addMongoDBOK.Payload.Service.ServiceID
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
			Mongodb: &services.GetServiceOKBodyMongodb{
				ServiceID:   serviceID,
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        27017,
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

		require.Len(t, listAgents.Payload.MongodbExporter, 1)
		require.Len(t, listAgents.Payload.QANMongodbProfilerAgent, 1)
		assert.Equal(t, agents.ListAgentsOKBody{
			MongodbExporter: []*agents.MongodbExporterItems0{
				{
					AgentID:    listAgents.Payload.MongodbExporter[0].AgentID,
					ServiceID:  serviceID,
					PMMAgentID: pmmAgentID,
					Username:   "username",
				},
			},
			QANMongodbProfilerAgent: []*agents.QANMongodbProfilerAgentItems0{
				{
					AgentID:    listAgents.Payload.QANMongodbProfilerAgent[0].AgentID,
					ServiceID:  serviceID,
					PMMAgentID: pmmAgentID,
					Username:   "username",
				},
			},
		}, *listAgents.Payload)
	})

	t.Run("With labels", func(realT *testing.T) {

		nodeName := pmmapitests.TestString(realT, "node-name-for-all-fields")
		nodeID, pmmAgentID := registerGenericNode(realT, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(realT, nodeID)
		defer removePMMAgentWithSubAgents(realT, pmmAgentID)

		serviceName := pmmapitests.TestString(realT, "service-name-for-all-fields")

		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:         nodeID,
				PMMAgentID:     pmmAgentID,
				ServiceName:    serviceName,
				Address:        "10.10.10.10",
				Port:           27017,
				Environment:    "some-environment",
				Cluster:        "cluster-name",
				ReplicationSet: "replication-set",
				CustomLabels:   map[string]string{"bar": "foo"},

				SkipConnectionCheck: true,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		require.NoError(realT, err)
		require.NotNil(realT, addMongoDBOK)
		require.NotNil(realT, addMongoDBOK.Payload.Service)
		serviceID := addMongoDBOK.Payload.Service.ServiceID
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
		assert.Equal(t, services.GetServiceOKBody{
			Mongodb: &services.GetServiceOKBodyMongodb{
				ServiceID:      serviceID,
				NodeID:         nodeID,
				ServiceName:    serviceName,
				Address:        "10.10.10.10",
				Port:           27017,
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

		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        27017,

				SkipConnectionCheck: true,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		require.NoError(t, err)
		require.NotNil(t, addMongoDBOK)
		require.NotNil(t, addMongoDBOK.Payload.Service)
		serviceID := addMongoDBOK.Payload.Service.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removeServiceAgents(t, serviceID)

		params = &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:      nodeID,
				PMMAgentID:  pmmAgentID,
				ServiceName: serviceName,
				Address:     "11.11.11.11",
				Port:        27017,
			},
		}
		addMongoDBOK, err = client.Default.MongoDB.AddMongoDB(params)
		require.Nil(t, addMongoDBOK)
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, `Service with name %q already exists.`, serviceName)
	})

	t.Run("Empty Service Name", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body:    mongodb.AddMongoDBBody{NodeID: nodeID},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceName: value '' must not be an empty string")
		assert.Nil(t, addMongoDBOK)
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
		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Address: value '' must not be an empty string")
		assert.Nil(t, addMongoDBOK)
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
		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Port: value '0' must be greater than '0'")
		assert.Nil(t, addMongoDBOK)
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
		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:      nodeID,
				ServiceName: serviceName,
				Address:     "10.10.10.10",
				Port:        3306,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field PmmAgentId: value '' must not be an empty string")
		assert.Nil(t, addMongoDBOK)
	})
}

func TestRemoveMongoDB(t *testing.T) {
	addMongoDB := func(t *testing.T, serviceName, nodeName string, withAgents bool) (nodeID string, pmmAgentID string, serviceID string) {
		t.Helper()
		nodeID, pmmAgentID = registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})

		params := &mongodb.AddMongoDBParams{
			Context: pmmapitests.Context,
			Body: mongodb.AddMongoDBBody{
				NodeID:             nodeID,
				PMMAgentID:         pmmAgentID,
				ServiceName:        serviceName,
				Address:            "10.10.10.10",
				Port:               27017,
				Username:           "username",
				Password:           "password",
				QANMongodbProfiler: withAgents,

				SkipConnectionCheck: true,
			},
		}
		addMongoDBOK, err := client.Default.MongoDB.AddMongoDB(params)
		require.NoError(t, err)
		require.NotNil(t, addMongoDBOK)
		require.NotNil(t, addMongoDBOK.Payload.Service)
		serviceID = addMongoDBOK.Payload.Service.ServiceID
		return
	}

	t.Run("By name", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-by-name")
		nodeName := pmmapitests.TestString(t, "node-remove-by-name")
		nodeID, pmmAgentID, serviceID := addMongoDB(t, serviceName, nodeName, true)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceName: serviceName,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypeMONGODBSERVICE),
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
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Service with ID %q not found.", serviceID)
		assert.Nil(t, listAgents)
	})

	t.Run("By ID", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-by-id")
		nodeName := pmmapitests.TestString(t, "node-remove-by-id")
		nodeID, pmmAgentID, serviceID := addMongoDB(t, serviceName, nodeName, true)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypeMONGODBSERVICE),
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
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Service with ID %q not found.", serviceID)
		assert.Nil(t, listAgents)
	})

	t.Run("Both params", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-both-params")
		nodeName := pmmapitests.TestString(t, "node-remove-both-params")
		nodeID, pmmAgentID, serviceID := addMongoDB(t, serviceName, nodeName, false)
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer pmmapitests.RemoveServices(t, serviceID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		removeServiceOK, err := client.Default.Service.RemoveService(&service.RemoveServiceParams{
			Body: service.RemoveServiceBody{
				ServiceID:   serviceID,
				ServiceName: serviceName,
				ServiceType: pointer.ToString(service.RemoveServiceBodyServiceTypeMYSQLSERVICE),
			},
			Context: pmmapitests.Context,
		})
		assert.Nil(t, removeServiceOK)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "service_id or service_name expected; not both")
	})

	t.Run("Wrong type", func(t *testing.T) {
		serviceName := pmmapitests.TestString(t, "service-remove-wrong-type")
		nodeName := pmmapitests.TestString(t, "node-remove-wrong-type")
		nodeID, pmmAgentID, serviceID := addMongoDB(t, serviceName, nodeName, false)
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
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "wrong service type")
	})
}
