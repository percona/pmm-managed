package inventory

import (
	"context"
	"reflect"
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/percona/pmm/api/inventory/json/client/agents"
	"github.com/percona/pmm/api/inventory/json/client/nodes"
	"github.com/percona/pmm/api/inventory/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

type ErrorResponse interface {
	Code() int
}

type ServerResponse struct {
	Code  int
	Error string
}

func removeNodes(t *testing.T, nodeIDs ...string) {
	t.Helper()

	for _, nodeID := range nodeIDs {
		params := &nodes.RemoveNodeParams{
			Body: nodes.RemoveNodeBody{
				NodeID: nodeID,
			},
			Context: context.Background(),
		}
		res, err := client.Default.Nodes.RemoveNode(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

func addGenericNode(t *testing.T, nodeName string) *nodes.AddGenericNodeOKBodyGeneric {
	t.Helper()

	params := &nodes.AddGenericNodeParams{
		Body: nodes.AddGenericNodeBody{
			NodeName: nodeName,
		},
		Context: pmmapitests.Context,
	}
	res, err := client.Default.Nodes.AddGenericNode(params)
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Payload)
	require.NotNil(t, res.Payload.Generic)
	return res.Payload.Generic
}

func addRemoteNode(t *testing.T, nodeName string) *nodes.AddRemoteNodeOKBody {
	t.Helper()

	params := &nodes.AddRemoteNodeParams{
		Body: nodes.AddRemoteNodeBody{
			NodeName: nodeName,
		},
		Context: pmmapitests.Context,
	}
	res, err := client.Default.Nodes.AddRemoteNode(params)
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func removeServices(t *testing.T, serviceIDs ...string) {
	t.Helper()

	for _, serviceID := range serviceIDs {
		params := &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
			},
			Context: context.Background(),
		}
		res, err := client.Default.Services.RemoveService(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

func addMySQLService(t *testing.T, body services.AddMySQLServiceBody) *services.AddMySQLServiceOKBody {
	t.Helper()

	params := &services.AddMySQLServiceParams{
		Body:    body,
		Context: pmmapitests.Context,
	}
	res, err := client.Default.Services.AddMySQLService(params)
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func addPostgreSQLService(t *testing.T, body services.AddPostgreSQLServiceBody) *services.AddPostgreSQLServiceOKBody {
	t.Helper()

	params := &services.AddPostgreSQLServiceParams{
		Body:    body,
		Context: pmmapitests.Context,
	}
	res, err := client.Default.Services.AddPostgreSQLService(params)
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func removeAgents(t *testing.T, agentIDs ...string) {
	t.Helper()

	for _, agentID := range agentIDs {
		params := &agents.RemoveAgentParams{
			Body: agents.RemoveAgentBody{
				AgentID: agentID,
			},
			Context: context.Background(),
		}
		res, err := client.Default.Agents.RemoveAgent(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

func addPMMAgent(t *testing.T, nodeID string) *agents.AddPMMAgentOKBody {
	t.Helper()

	res, err := client.Default.Agents.AddPMMAgent(&agents.AddPMMAgentParams{
		Body: agents.AddPMMAgentBody{
			RunsOnNodeID: nodeID,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func addMySqldExporter(t *testing.T, body agents.AddMySqldExporterBody) *agents.AddMySqldExporterOKBody {
	t.Helper()

	res, err := client.Default.Agents.AddMySqldExporter(&agents.AddMySqldExporterParams{
		Body:    body,
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func addMongoDBExporter(t *testing.T, body agents.AddMongoDBExporterBody) *agents.AddMongoDBExporterOKBody {
	t.Helper()

	res, err := client.Default.Agents.AddMongoDBExporter(&agents.AddMongoDBExporterParams{
		Body:    body,
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func addPostgresExporter(t *testing.T, body agents.AddPostgresExporterBody) *agents.AddPostgresExporterOKBody {
	t.Helper()

	res, err := client.Default.Agents.AddPostgresExporter(&agents.AddPostgresExporterParams{
		Body:    body,
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
}

func assertEqualAPIError(t *testing.T, err error, expected ServerResponse) bool {
	t.Helper()

	if !assert.Error(t, err) {
		return false
	}

	assert.Equal(t, expected.Code, err.(ErrorResponse).Code())

	// Have to use reflect because there are a lot of types with the same structure and different names.
	val := reflect.ValueOf(err)

	payload := val.Elem().FieldByName("Payload")
	if !assert.True(t, payload.IsValid(), "Wrong response structure. There is no field Payload.") {
		return false
	}

	errorField := payload.Elem().FieldByName("Error")
	if !assert.True(t, errorField.IsValid(), "Wrong response structure. There is no field Error in Payload.") {
		return false
	}

	return assert.Equal(t, expected.Error, errorField.String())
}

func assertMySQLServiceExists(t *testing.T, res *services.ListServicesOK, serviceID string) bool {
	t.Helper()

	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.Mysql {
			if v.ServiceID == serviceID {
				return true
			}
		}
		return false
	}, "There should be MySQL service with id `%s`", serviceID)
}

func assertMySQLServiceNotExist(t *testing.T, res *services.ListServicesOK, serviceID string) bool {
	t.Helper()

	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.Mysql {
			if v.ServiceID == serviceID {
				return false
			}
		}
		return true
	}, "There should not be MySQL service with id `%s`", serviceID)
}

func assertMySQLExporterExists(t *testing.T, res *agents.ListAgentsOK, mySqldExporterID string) bool {
	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.MysqldExporter {
			if v.AgentID == mySqldExporterID {
				return true
			}
		}
		return false
	}, "There should be MySQL agent with id `%s`", mySqldExporterID)
}

func assertMySQLExporterNotExists(t *testing.T, res *agents.ListAgentsOK, mySqldExporterID string) bool {
	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.MysqldExporter {
			if v.AgentID == mySqldExporterID {
				return false
			}
		}
		return true
	}, "There should be MySQL agent with id `%s`", mySqldExporterID)
}

func assertPMMAgentExists(t *testing.T, res *agents.ListAgentsOK, pmmAgentID string) bool {
	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.PMMAgent {
			if v.AgentID == pmmAgentID {
				return true
			}
		}
		return false
	}, "There should be PMM-agent with id `%s`", pmmAgentID)
}

func assertPMMAgentNotExists(t *testing.T, res *agents.ListAgentsOK, pmmAgentID string) bool {
	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.PMMAgent {
			if v.AgentID == pmmAgentID {
				return false
			}
		}
		return true
	}, "There should be PMM-agent with id `%s`", pmmAgentID)
}
