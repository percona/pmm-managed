package inventory

import (
	"context"
	"fmt"
	"testing"

	"github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{404, "Service with ID \"pmm-not-found\" not found."})
		assert.Nil(t, res)
	})

	t.Run("EmptyServiceID", func(t *testing.T) {
		t.Parallel()

		params := &services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.GetService(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
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

		_ = addMySqldExporter(t, agents.AddMySqldExporterBody{
			ServiceID:  serviceID,
			Username:   "username",
			Password:   "password",
			PMMAgentID: pmmAgentID,
		})

		params := &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.RemoveService(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{412, fmt.Sprintf(`Service with ID "%s" has agents.`, serviceID)})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{404, fmt.Sprintf("Service with ID %q not found.", serviceID)})
		assert.Nil(t, getServiceResp)

		listAgentsOK, err := client.Default.Agents.ListAgents(&agents.ListAgentsParams{
			Body: agents.ListAgentsBody{
				ServiceID: serviceID,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ListAgentsOKBody{}, listAgentsOK.Payload)
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{404, fmt.Sprintf(`Service with ID "%s" not found.`, serviceID)})
		assert.Nil(t, res)
	})

	t.Run("Empty params", func(t *testing.T) {
		t.Parallel()
		removeResp, err := client.Default.Services.RemoveService(&services.RemoveServiceParams{
			Body:    services.RemoveServiceBody{},
			Context: context.Background(),
		})
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{409, fmt.Sprintf("Service with name %q already exists.", serviceName)})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field NodeId: value '' must not be an empty string"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field Port: value '0' must be greater than '0'"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field ServiceName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Mysql.ServiceID)
		}
	})
}

func TestAmazonRDSMySQLService(t *testing.T) {
	t.Skip("Not implemented yet.")

	t.Run("Basic", func(t *testing.T) {
		remoteNodeOKBody := addRemoteNode(t, pmmapitests.TestString(t, "Remote node to check services filter"))
		remoteNodeID := remoteNodeOKBody.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, remoteNodeID)

		serviceName := pmmapitests.TestString(t, "Basic AmazonRDSMySQL Service")
		params := &services.AddAmazonRDSMySQLServiceParams{
			Body: services.AddAmazonRDSMySQLServiceBody{
				NodeID:      remoteNodeID,
				Address:     "localhost",
				Port:        3306,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddAmazonRDSMySQLService(params)
		assert.NoError(t, err)
		require.NotNil(t, res)
		serviceID := res.Payload.AmazonRDSMysql.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)
		assert.Equal(t, &services.AddAmazonRDSMySQLServiceOK{
			Payload: &services.AddAmazonRDSMySQLServiceOKBody{
				AmazonRDSMysql: &services.AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql{
					ServiceID:   serviceID,
					NodeID:      remoteNodeID,
					Address:     "localhost",
					Port:        3306,
					ServiceName: serviceName,
				},
			},
		}, res)

		// Check if the service saved in pmm-managed.
		serviceRes, err := client.Default.Services.GetService(&services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: serviceID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		require.NotNil(t, serviceRes)
		assert.Equal(t, &services.GetServiceOK{
			Payload: &services.GetServiceOKBody{
				AmazonRDSMysql: &services.GetServiceOKBodyAmazonRDSMysql{
					ServiceID:   serviceID,
					NodeID:      remoteNodeID,
					Address:     "localhost",
					Port:        3306,
					ServiceName: serviceName,
				},
			},
		}, serviceRes)

		// Check duplicates.
		params = &services.AddAmazonRDSMySQLServiceParams{
			Body: services.AddAmazonRDSMySQLServiceBody{
				NodeID:      remoteNodeID,
				Address:     "127.0.0.1",
				Port:        3336,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		res, err = client.Default.Services.AddAmazonRDSMySQLService(params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{409, ""})
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.AmazonRDSMysql.ServiceID)
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{409, fmt.Sprintf("Service with name %q already exists.", serviceName)})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field NodeId: value '' must not be an empty string"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field ServiceName: value '' must not be an empty string"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{409, fmt.Sprintf("Service with name %q already exists.", serviceName)})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field NodeId: value '' must not be an empty string"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field Port: value '0' must be greater than '0'"})
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
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{400, "invalid field ServiceName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			pmmapitests.RemoveServices(t, res.Payload.Postgresql.ServiceID)
		}
	})
}
