package inventory

import (
	"testing"

	"github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/nodes"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func addGenericNode(t *testing.T, nodeName string) *nodes.AddGenericNodeOKBodyGeneric {
	t.Helper()

	params := &nodes.AddGenericNodeParams{
		Body: nodes.AddGenericNodeBody{
			NodeName: nodeName,
			Address:  "10.10.10.10",
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

func addProxySQLService(t *testing.T, body services.AddProxySQLServiceBody) *services.AddProxySQLServiceOKBody {
	t.Helper()

	params := &services.AddProxySQLServiceParams{
		Body:    body,
		Context: pmmapitests.Context,
	}
	res, err := client.Default.Services.AddProxySQLService(params)
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
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

func addNodeExporter(t *testing.T, pmmAgentID string, customLabels map[string]string) *agents.AddNodeExporterOK {
	res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
		Body: agents.AddNodeExporterBody{
			PMMAgentID:   pmmAgentID,
			CustomLabels: customLabels,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Payload.NodeExporter)
	require.Equal(t, pmmAgentID, res.Payload.NodeExporter.PMMAgentID)
	return res
}

func addMySQLdExporter(t *testing.T, body agents.AddMySQLdExporterBody) *agents.AddMySQLdExporterOKBody {
	t.Helper()

	res, err := client.Default.Agents.AddMySQLdExporter(&agents.AddMySQLdExporterParams{
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

func addProxySQLExporter(t *testing.T, body agents.AddProxySQLExporterBody) *agents.AddProxySQLExporterOKBody {
	t.Helper()

	res, err := client.Default.Agents.AddProxySQLExporter(&agents.AddProxySQLExporterParams{
		Body:    body,
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.NotNil(t, res)
	return res.Payload
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

func assertNodeExporterExists(t *testing.T, res *agents.ListAgentsOK, nodeExporterID string) bool {
	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.NodeExporter {
			if v.AgentID == nodeExporterID {
				return true
			}
		}
		return false
	}, "There should be Node exporter with id `%s`", nodeExporterID)
}

func assertNodeExporterNotExists(t *testing.T, res *agents.ListAgentsOK, nodeExporterID string) bool {
	return assert.Conditionf(t, func() bool {
		for _, v := range res.Payload.NodeExporter {
			if v.AgentID == nodeExporterID {
				return false
			}
		}
		return true
	}, "There should be Node exporter with id `%s`", nodeExporterID)
}
