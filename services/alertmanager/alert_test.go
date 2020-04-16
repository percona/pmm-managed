package alertmanager

import (
	"testing"

	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/pmm-managed/models"
)

func TestMakeAlert(t *testing.T) {
	agent := &models.Agent{
		AgentID: "/agent_id/123",
	}
	node := &models.Node{
		NodeID:   "/node_id/456",
		NodeName: "nodename",
	}
	name, alert, err := makeAlertPMMAgentNotConnected(agent, node)
	require.NoError(t, err)

	assert.Equal(t, "pmm_agent_not_connected", name)

	expected := &ammodels.PostableAlert{
		Alert: ammodels.Alert{
			Labels: ammodels.LabelSet{
				"agent_id":  "/agent_id/123",
				"alertname": "pmm_agent_not_connected",
				"node_id":   "/node_id/456",
				"node_name": "nodename",
				"severity":  "warning",
			},
		},
		Annotations: ammodels.LabelSet{
			"summary":     "pmm-agent is not connected to PMM Server",
			"description": "Node name: nodename",
		},
	}
	assert.Equal(t, expected, alert)
}
