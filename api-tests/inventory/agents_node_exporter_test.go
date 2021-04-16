package inventory

import (
	"testing"

	"github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestNodeExporter(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		node := pmmapitests.AddRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, nodeID)

		pmmAgent := pmmapitests.AddPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer pmmapitests.RemoveAgents(t, pmmAgentID)

		customLabels := map[string]string{
			"custom_label_node_exporter": "node_exporter",
		}
		res := addNodeExporter(t, pmmAgentID, customLabels)
		agentID := res.Payload.NodeExporter.AgentID
		defer pmmapitests.RemoveAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				NodeExporter: &agents.GetAgentOKBodyNodeExporter{
					AgentID:      agentID,
					PMMAgentID:   pmmAgentID,
					Disabled:     false,
					CustomLabels: customLabels,
				},
			},
		}, getAgentRes)

		// Test change API.
		changeNodeExporterOK, err := client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID: agentID,
				Common: &agents.ChangeNodeExporterParamsBodyCommon{
					Disable:            true,
					RemoveCustomLabels: true,
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
					Disabled:   true,
				},
			},
		}, changeNodeExporterOK)

		changeNodeExporterOK, err = client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID: agentID,
				Common: &agents.ChangeNodeExporterParamsBodyCommon{
					Enable: true,
					CustomLabels: map[string]string{
						"new_label": "node_exporter",
					},
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
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field PmmAgentId: value '' must not be an empty string")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveNodes(t, res.Payload.NodeExporter.AgentID)
		}
	})

	t.Run("NotExistPmmAgentID", func(t *testing.T) {
		t.Parallel()

		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body:    agents.AddNodeExporterBody{PMMAgentID: "pmm-node-exporter-node"},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Agent with ID \"pmm-node-exporter-node\" not found.")
		if !assert.Nil(t, res) {
			pmmapitests.RemoveNodes(t, res.Payload.NodeExporter.AgentID)
		}
	})

	t.Run("With PushMetrics", func(t *testing.T) {
		t.Parallel()

		node := pmmapitests.AddRemoteNode(t, pmmapitests.TestString(t, "Remote node for Node exporter"))
		nodeID := node.Remote.NodeID
		defer pmmapitests.RemoveNodes(t, nodeID)

		pmmAgent := pmmapitests.AddPMMAgent(t, nodeID)
		pmmAgentID := pmmAgent.PMMAgent.AgentID
		defer pmmapitests.RemoveAgents(t, pmmAgentID)

		customLabels := map[string]string{
			"custom_label_node_exporter": "node_exporter",
		}
		res, err := client.Default.Agents.AddNodeExporter(&agents.AddNodeExporterParams{
			Body: agents.AddNodeExporterBody{
				PMMAgentID:   pmmAgentID,
				CustomLabels: customLabels,
				PushMetrics:  true,
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Payload.NodeExporter)
		require.Equal(t, pmmAgentID, res.Payload.NodeExporter.PMMAgentID)
		agentID := res.Payload.NodeExporter.AgentID
		defer pmmapitests.RemoveAgents(t, agentID)

		getAgentRes, err := client.Default.Agents.GetAgent(&agents.GetAgentParams{
			Body:    agents.GetAgentBody{AgentID: agentID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, &agents.GetAgentOK{
			Payload: &agents.GetAgentOKBody{
				NodeExporter: &agents.GetAgentOKBodyNodeExporter{
					AgentID:            agentID,
					PMMAgentID:         pmmAgentID,
					Disabled:           false,
					CustomLabels:       customLabels,
					PushMetricsEnabled: true,
				},
			},
		}, getAgentRes)

		// Test change API.
		changeNodeExporterOK, err := client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID: agentID,
				Common: &agents.ChangeNodeExporterParamsBodyCommon{
					DisablePushMetrics: true,
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeNodeExporterOK{
			Payload: &agents.ChangeNodeExporterOKBody{
				NodeExporter: &agents.ChangeNodeExporterOKBodyNodeExporter{
					AgentID:      agentID,
					PMMAgentID:   pmmAgentID,
					Disabled:     false,
					CustomLabels: customLabels,
				},
			},
		}, changeNodeExporterOK)

		changeNodeExporterOK, err = client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID: agentID,
				Common: &agents.ChangeNodeExporterParamsBodyCommon{
					EnablePushMetrics: true,
				},
			},
			Context: pmmapitests.Context,
		})
		assert.NoError(t, err)
		assert.Equal(t, &agents.ChangeNodeExporterOK{
			Payload: &agents.ChangeNodeExporterOKBody{
				NodeExporter: &agents.ChangeNodeExporterOKBodyNodeExporter{
					AgentID:            agentID,
					PMMAgentID:         pmmAgentID,
					Disabled:           false,
					CustomLabels:       customLabels,
					PushMetricsEnabled: true,
				},
			},
		}, changeNodeExporterOK)
		_, err = client.Default.Agents.ChangeNodeExporter(&agents.ChangeNodeExporterParams{
			Body: agents.ChangeNodeExporterBody{
				AgentID: agentID,
				Common: &agents.ChangeNodeExporterParamsBodyCommon{
					EnablePushMetrics:  true,
					DisablePushMetrics: true,
				},
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "expected one of  param: enable_push_metrics or disable_push_metrics")
	})
}
