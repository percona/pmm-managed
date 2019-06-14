package management

import (
	"context"
	"testing"

	inventoryClient "github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/nodes"
	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func registerGenericNode(t *testing.T, body node.RegisterBody) (string, string) {
	t.Helper()
	params := node.RegisterParams{
		Context: pmmapitests.Context,
		Body:    body,
	}
	registerOK, err := client.Default.Node.Register(&params)
	require.NoError(t, err)
	require.NotNil(t, registerOK)
	require.NotNil(t, registerOK.Payload.PMMAgent)
	require.NotNil(t, registerOK.Payload.PMMAgent.AgentID)
	require.NotNil(t, registerOK.Payload.GenericNode)
	require.NotNil(t, registerOK.Payload.GenericNode.NodeID)
	return registerOK.Payload.GenericNode.NodeID, registerOK.Payload.PMMAgent.AgentID
}

func registerContainerNode(t *testing.T, body node.RegisterBody) (string, string) {
	t.Helper()
	params := node.RegisterParams{
		Context: pmmapitests.Context,
		Body:    body,
	}
	registerOK, err := client.Default.Node.Register(&params)
	require.NoError(t, err)
	require.NotNil(t, registerOK)
	require.NotNil(t, registerOK.Payload.PMMAgent)
	require.NotNil(t, registerOK.Payload.PMMAgent.AgentID)
	require.NotNil(t, registerOK.Payload.ContainerNode)
	require.NotNil(t, registerOK.Payload.ContainerNode.NodeID)
	return registerOK.Payload.ContainerNode.NodeID, registerOK.Payload.PMMAgent.AgentID
}

func assertNodeExporterCreated(t *testing.T, pmmAgentID string) (string, bool) {
	t.Helper()
	listAgentsOK, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
		Body: agents.ListAgentsBody{
			PMMAgentID: pmmAgentID,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	require.Len(t, listAgentsOK.Payload.NodeExporter, 1)
	nodeExporterAgentID := listAgentsOK.Payload.NodeExporter[0].AgentID
	asserted := assert.Equal(t, agents.NodeExporterItems0{
		PMMAgentID: pmmAgentID,
		AgentID:    nodeExporterAgentID,
	}, *listAgentsOK.Payload.NodeExporter[0])
	return nodeExporterAgentID, asserted
}

func assertPMMAgentCreated(t *testing.T, nodeID string, pmmAgentID string) {
	t.Helper()
	agentOK, err := inventoryClient.Default.Agents.GetAgent(&agents.GetAgentParams{
		Body: agents.GetAgentBody{
			AgentID: pmmAgentID,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	assert.Equal(t, agents.GetAgentOKBody{
		PMMAgent: &agents.GetAgentOKBodyPMMAgent{
			AgentID:      pmmAgentID,
			RunsOnNodeID: nodeID,
		},
	}, *agentOK.Payload)
}

func assertNodeCreated(t assert.TestingT, nodeID string, expectedResult nodes.GetNodeOKBody) {
	if n, ok := t.(interface {
		Helper()
	}); ok {
		n.Helper()
	}
	nodeOK, err := inventoryClient.Default.Nodes.GetNode(&nodes.GetNodeParams{
		Body: nodes.GetNodeBody{
			NodeID: nodeID,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, *nodeOK.Payload)
}

func removePMMAgentWithSubAgents(t *testing.T, pmmAgentID string) {
	t.Helper()
	listAgentsOK, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
		Body: agents.ListAgentsBody{
			PMMAgentID: pmmAgentID,
		},
		Context: context.Background(),
	})
	assert.NoError(t, err)
	removeAllAgentsInList(t, listAgentsOK)
	pmmapitests.RemoveAgents(t, pmmAgentID)
}

func removeServiceAgents(t *testing.T, serviceID string) {
	t.Helper()
	listAgentsOK, err := inventoryClient.Default.Agents.ListAgents(&agents.ListAgentsParams{
		Body: agents.ListAgentsBody{
			ServiceID: serviceID,
		},
		Context: context.Background(),
	})
	assert.NoError(t, err)
	removeAllAgentsInList(t, listAgentsOK)
}

func removeAllAgentsInList(t *testing.T, listAgentsOK *agents.ListAgentsOK) {
	t.Helper()
	var agentIDs []string
	for _, agent := range listAgentsOK.Payload.NodeExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.PMMAgent {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.PostgresExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.MysqldExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.ProxysqlExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.QANMysqlPerfschemaAgent {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.MongodbExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.RDSExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.ExternalExporter {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.QANMongodbProfilerAgent {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	for _, agent := range listAgentsOK.Payload.QANMysqlSlowlogAgent {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	pmmapitests.RemoveAgents(t, agentIDs...)
}
