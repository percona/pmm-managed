package inventory

import (
	"fmt"
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/percona/pmm/api/inventory/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestServices(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		remoteNodeOKBody := addRemoteNode(t, pmmapitests.TestString(t, "Remote node for services test"))
		remoteNodeID := remoteNodeOKBody.Remote.NodeID
		defer removeNodes(t, remoteNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		remoteService := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      remoteNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service on remote Node"),
		})
		remoteServiceID := remoteService.Mysql.ServiceID
		defer removeServices(t, remoteServiceID)

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
		defer removeNodes(t, genericNodeID)

		remoteNodeOKBody := addRemoteNode(t, pmmapitests.TestString(t, "Remote node to check services filter"))
		remoteNodeID := remoteNodeOKBody.Remote.NodeID
		defer removeNodes(t, remoteNodeID)

		service := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      genericNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service for filters test"),
		})
		serviceID := service.Mysql.ServiceID
		defer removeServices(t, serviceID)

		remoteService := addMySQLService(t, services.AddMySQLServiceBody{
			NodeID:      remoteNodeID,
			Address:     "localhost",
			Port:        3306,
			ServiceName: pmmapitests.TestString(t, "Some MySQL Service on remote Node for filters test"),
		})
		remoteServiceID := remoteService.Mysql.ServiceID
		defer removeServices(t, remoteServiceID)

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
		assertEqualAPIError(t, err, ServerResponse{404, "Service with ID \"pmm-not-found\" not found."})
		assert.Nil(t, res)
	})

	t.Run("EmptyServiceID", func(t *testing.T) {
		t.Parallel()

		params := &services.GetServiceParams{
			Body:    services.GetServiceBody{ServiceID: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.GetService(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceId: value '' must not be an empty string"})
		assert.Nil(t, res)
	})
}

func TestMySQLService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

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
		defer removeServices(t, serviceID)

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
		assertEqualAPIError(t, err, ServerResponse{409, fmt.Sprintf("Service with name %q already exists.", serviceName)})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.Mysql.ServiceID)
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
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.Mysql.ServiceID)
		}
	})

	t.Run("AddServiceNameEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		params := &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      genericNodeID,
				ServiceName: "",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMySQLService(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.Mysql.ServiceID)
		}
	})
}

func TestAmazonRDSMySQLService(t *testing.T) {
	t.Skip("Not implemented yet.")

	t.Run("Basic", func(t *testing.T) {
		remoteNodeOKBody := addRemoteNode(t, pmmapitests.TestString(t, "Remote node to check services filter"))
		remoteNodeID := remoteNodeOKBody.Remote.NodeID
		defer removeNodes(t, remoteNodeID)

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
		defer removeServices(t, serviceID)
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
		assertEqualAPIError(t, err, ServerResponse{409, ""})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.AmazonRDSMysql.ServiceID)
		}
	})
}

func TestMongoDBService(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

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
		defer removeServices(t, serviceID)

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
		assertEqualAPIError(t, err, ServerResponse{409, fmt.Sprintf("Service with name %q already exists.", serviceName)})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.Mongodb.ServiceID)
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
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeId: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.Mongodb.ServiceID)
		}
	})

	t.Run("AddServiceNameEmpty", func(t *testing.T) {
		t.Parallel()

		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "")).NodeID
		require.NotEmpty(t, genericNodeID)
		defer removeNodes(t, genericNodeID)

		params := &services.AddMongoDBServiceParams{
			Body: services.AddMongoDBServiceBody{
				NodeID:      genericNodeID,
				ServiceName: "",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Services.AddMongoDBService(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field ServiceName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeServices(t, res.Payload.Mongodb.ServiceID)
		}
	})
}
