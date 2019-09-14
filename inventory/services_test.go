package inventory

import (
	"context"
	"testing"

	"github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestServices(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		remoteNodeOKBody := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for services test"))
		remoteNodeID := remoteNodeOKBody.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, remoteNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service"),
		})
		serviceID := service.Mysql.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)

		remoteService := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      remoteNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service on remote Node"),
		})
		remoteServiceID := remoteService.Mysql.ServiceID
		defer pmmapitests.RemoveServices(t, remoteServiceID)

		res, err := client.Default.Services.ListServices(&services.ListServicesParams{Context: pmmapitests.Context})
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.NotZerof(t, len(res.Payload.Mysql), "There should be at least one node")
		assertMySQLServiceExists(t, res, serviceID)
		assertMySQLServiceExists(t, res, remoteServiceID)
	})

	t.Run("FilterList", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		remoteNodeOKBody := addRemoteNode(t, pmmapitests.TestString(t, "Remote node to check services filter"))
		remoteNodeID := remoteNodeOKBody.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, remoteNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service for filters test"),
		})
		serviceID := service.Mysql.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)

		remoteService := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      remoteNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service on remote Node for filters test"),
		})
		remoteServiceID := remoteService.Mysql.ServiceID
		defer pmmapitests.RemoveServices(t, remoteServiceID)

		res, err := client.Default.Services.ListServices(&services.ListServicesParams{
			Body:    services.ListServicesBody{NodeID: remoteNodeID},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.NotZerof(t, len(res.Payload.Mysql), "There should be at least one node")
		assertMySQLServiceNotExist(t, res, serviceID)
		assertMySQLServiceExists(t, res, remoteServiceID)
	})
}

func TestGetService(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		params := &services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: "pmm-not-found"},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.GetService(params)
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Service with ID \"pmm-not-found\" not found.")
		assert.Nil(t, res)
	})

	t.Run("EmptyServiceID", func(t *testing.T) {
		t.Parallel()

		params := &services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.GetService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceId: value '' must not be an empty string")
		assert.Nil(t, res)
	})
}

func TestRemoveService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for agents list"))
		nodeID := node.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      nodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID

		params := &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.RemoveService(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Has agents", func(t *testing.T) {
		t.Parallel()

		node := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for agents list"))
		nodeID := node.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, nodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      nodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "MySQL Service for agent"),
		})
		serviceID := service.Mysql.ServiceID

		pmmAgent := addPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer pmmapitests.RemoveAgents(t, pmmAgentID)

		_ = addMySQLdExporter(t, agents.AddMySQLdExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,

			SkipConnectionCheck: true,
		})

		params := &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.RemoveService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.FailedPrecondition, `Service with ID %q has agents.`, serviceID)
		assert.Nil(t, res)

		// Remove with force flag.
		params = &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
				Force:     true,
			},
			Context: pmmapitests.Context,
		}
		res, err = client.Default.Services.RemoveService(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		// Check that the service and agents are removed.
		getServiceResp, err := client.Default.Services.GetService(&services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Service with ID %q not found.", serviceID)
		assert.Nil(t, getServiceResp)

		listAgentsOK, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Service with ID %q not found.", serviceID)
		assert.Nil(t, listAgentsOK)
	})

	t.Run("Not-exist service", func(t *testing.T) {
		t.Parallel()
		serviceID := "not-exist-service-id"

		params := &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.RemoveService(params)
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, `Service with ID %q not found.`, serviceID)
		assert.Nil(t, res)
	})

	t.Run("Empty params", func(t *testing.T) {
		t.Parallel()
		removeResp, err := client.Default.Services.RemoveService(&services.RemoveServiceParams{
			Body:    services.RemoveServiceBody{},
			Context: context.Background(),
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceId: value '' must not be an empty string")
		assert.Nil(t, removeResp)
	})
}

func TestMySQLService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		serviceName := pmmapitests.TestString(t, "Basic MySQL Service")
		params := &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				Port:        3306,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMySQLService(params)
		assert.NoError(t, err)
		require.NotNil(t, res)
		serviceID := res.Payload.Mysql.ServiceID
		assert.Equal(t, &services.AddMySQLServiceOK{
			Payload: &services.AddMySQLServiceOKBody{
				Mysql: &services.AddMySQLServiceOKBodyMysql{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					Address:     "localhost",
					Port:        3306,
					ServiceName: serviceName,
				},
			},
		}, res)
		defer pmmapitests.RemoveServices(t, serviceID)

		// Check if the service saved in pmm-managed.
		serviceRes, err := client.Default.Services.GetService(&services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.NotNil(t, serviceRes)
		assert.Equal(t, &services.GetServiceOK{
			Payload: &services.GetServiceOKBody{
				Mysql: &services.GetServiceOKBodyMysql{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					Address:     "localhost",
					Port:        3306,
					ServiceName: serviceName,
				},
			},
		}, serviceRes)

		// Check duplicates.
		params = &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "127.0.0.1",
				Port:        3336,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err = client.Default.Services.AddMySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, "Service with name %q already exists.", serviceName)
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mysql.ServiceID)
		}
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		params := &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      "",
				Address:     "localhost",
				Port:        3306,
				ServiceName: pmmapitests.TestString(t, "MySQL Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field NodeId: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mysql.ServiceID)
		}
	})

	t.Run("AddEmptyPort", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				ServiceName: pmmapitests.TestString(t, "MySQL Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Port: value '0' must be greater than '0'")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mysql.ServiceID)
		}
	})

	t.Run("AddServiceNameEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      genericNodeID,
				ServiceName: "",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceName: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mysql.ServiceID)
		}
	})
}

func TestMongoDBService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		serviceName := pmmapitests.TestString(t, "Basic Mongo Service")
		params := &services.AddMongoDBServiceParams{
			Body: services.AddMongoDBServiceBody{
				NodeID:      genericNodeID,
				ServiceName: serviceName,
				Address:     "localhost",
				Port:        27017,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMongoDBService(params)
		assert.NoError(t, err)
		require.NotNil(t, res)
		serviceID := res.Payload.Mongodb.ServiceID
		assert.Equal(t, &services.AddMongoDBServiceOK{
			Payload: &services.AddMongoDBServiceOKBody{
				Mongodb: &services.AddMongoDBServiceOKBodyMongodb{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					ServiceName: serviceName,
					Address:     "localhost",
					Port:        27017,
				},
			},
		}, res)
		defer pmmapitests.RemoveServices(t, serviceID)

		// Check if the service saved in pmm-managed.
		serviceRes, err := client.Default.Services.GetService(&services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		require.NotNil(t, serviceRes)
		assert.Equal(t, &services.GetServiceOK{
			Payload: &services.GetServiceOKBody{
				Mongodb: &services.GetServiceOKBodyMongodb{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					ServiceName: serviceName,
					Address:     "localhost",
					Port:        27017,
				},
			},
		}, serviceRes)

		// Check duplicates.
		params = &services.AddMongoDBServiceParams{
			Body: services.AddMongoDBServiceBody{
				NodeID:      genericNodeID,
				ServiceName: serviceName,
				Address:     "localhost",
				Port:        27017,
			},
			Context: pmmapitests.Context,
		}
		res, err = client.Default.Services.AddMongoDBService(params)
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, "Service with name %q already exists.", serviceName)
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mongodb.ServiceID)
		}
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		params := &services.AddMongoDBServiceParams{
			Body: services.AddMongoDBServiceBody{
				NodeID:      "",
				ServiceName: pmmapitests.TestString(t, "MongoDB Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMongoDBService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field NodeId: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mongodb.ServiceID)
		}
	})

	t.Run("AddServiceNameEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddMongoDBServiceParams{
			Body: services.AddMongoDBServiceBody{
				NodeID:      genericNodeID,
				ServiceName: "",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMongoDBService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceName: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mongodb.ServiceID)
		}
	})
}

func TestPostgreSQLService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		serviceName := pmmapitests.TestString(t, "Basic PostgreSQL Service")
		params := &services.AddPostgreSQLServiceParams{
			Body: services.AddPostgreSQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				Port:        5432,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddPostgreSQLService(params)
		assert.NoError(t, err)
		require.NotNil(t, res)
		serviceID := res.Payload.Postgresql.ServiceID
		assert.Equal(t, &services.AddPostgreSQLServiceOK{
			Payload: &services.AddPostgreSQLServiceOKBody{
				Postgresql: &services.AddPostgreSQLServiceOKBodyPostgresql{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					Address:     "localhost",
					Port:        5432,
					ServiceName: serviceName,
				},
			},
		}, res)
		defer pmmapitests.RemoveServices(t, serviceID)

		// Check if the service saved in pmm-managed.
		serviceRes, err := client.Default.Services.GetService(&services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.NotNil(t, serviceRes)
		assert.Equal(t, &services.GetServiceOK{
			Payload: &services.GetServiceOKBody{
				Postgresql: &services.GetServiceOKBodyPostgresql{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					Address:     "localhost",
					Port:        5432,
					ServiceName: serviceName,
				},
			},
		}, serviceRes)

		// Check duplicates.
		params = &services.AddPostgreSQLServiceParams{
			Body: services.AddPostgreSQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "127.0.0.1",
				Port:        3336,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err = client.Default.Services.AddPostgreSQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, "Service with name %q already exists.", serviceName)
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Postgresql.ServiceID)
		}
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		params := &services.AddPostgreSQLServiceParams{
			Body: services.AddPostgreSQLServiceBody{
				NodeID:      "",
				Address:     "localhost",
				Port:        5432,
				ServiceName: pmmapitests.TestString(t, "PostgreSQL Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddPostgreSQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field NodeId: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Postgresql.ServiceID)
		}
	})

	t.Run("AddEmptyPort", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddPostgreSQLServiceParams{
			Body: services.AddPostgreSQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				ServiceName: pmmapitests.TestString(t, "PostgreSQL Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddPostgreSQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Port: value '0' must be greater than '0'")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Postgresql.ServiceID)
		}
	})

	t.Run("AddServiceNameEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddPostgreSQLServiceParams{
			Body: services.AddPostgreSQLServiceBody{
				NodeID:      genericNodeID,
				ServiceName: "",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddPostgreSQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceName: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Postgresql.ServiceID)
		}
	})
}

func TestProxySQLService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		serviceName := pmmapitests.TestString(t, "Basic ProxySQL Service")
		params := &services.AddProxySQLServiceParams{
			Body: services.AddProxySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				Port:        5432,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddProxySQLService(params)
		assert.NoError(t, err)
		require.NotNil(t, res)
		serviceID := res.Payload.Proxysql.ServiceID
		assert.Equal(t, &services.AddProxySQLServiceOK{
			Payload: &services.AddProxySQLServiceOKBody{
				Proxysql: &services.AddProxySQLServiceOKBodyProxysql{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					Address:     "localhost",
					Port:        5432,
					ServiceName: serviceName,
				},
			},
		}, res)
		defer pmmapitests.RemoveServices(t, serviceID)

		// Check if the service saved in pmm-managed.
		serviceRes, err := client.Default.Services.GetService(&services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.NotNil(t, serviceRes)
		assert.Equal(t, &services.GetServiceOK{
			Payload: &services.GetServiceOKBody{
				Proxysql: &services.GetServiceOKBodyProxysql{
					ServiceID:   serviceID,
					NodeID:      genericNodeID,
					Address:     "localhost",
					Port:        5432,
					ServiceName: serviceName,
				},
			},
		}, serviceRes)

		// Check duplicates.
		params = &services.AddProxySQLServiceParams{
			Body: services.AddProxySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "127.0.0.1",
				Port:        3336,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err = client.Default.Services.AddProxySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, "Service with name %q already exists.", serviceName)
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Proxysql.ServiceID)
		}
	})

	t.Run("AddNodeIDEmpty", func(t *testing.T) {
		t.Parallel()

		params := &services.AddProxySQLServiceParams{
			Body: services.AddProxySQLServiceBody{
				NodeID:      "",
				Address:     "localhost",
				Port:        5432,
				ServiceName: pmmapitests.TestString(t, "ProxySQL Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddProxySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field NodeId: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Proxysql.ServiceID)
		}
	})

	t.Run("AddEmptyPort", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddProxySQLServiceParams{
			Body: services.AddProxySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				ServiceName: pmmapitests.TestString(t, "ProxySQL Service with empty node id"),
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddProxySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Port: value '0' must be greater than '0'")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Proxysql.ServiceID)
		}
	})

	t.Run("AddServiceNameEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		params := &services.AddProxySQLServiceParams{
			Body: services.AddProxySQLServiceBody{
				NodeID:      genericNodeID,
				ServiceName: "",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddProxySQLService(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServiceName: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Proxysql.ServiceID)
		}
	})
}
