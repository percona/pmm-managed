package management

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/inventorypb/json/client/nodes"
	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/node"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestNodeRegister(t *testing.T) {
	t.Run("Generic Node", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			nodeName := pmmapitests.TestString(t, "node-name")
			nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
			})
			defer pmmapitests.RemoveNodes(t, nodeID)
			defer pmmapitests.RemoveAgents(t, pmmAgentID)

			// Check Node is created
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:   nodeID,
					NodeName: nodeName,
				},
			})

			// Check PMM Agent is created
			assertPMMAgentCreated(t, nodeID, pmmAgentID)

			// Check Node Exporter is created
			nodeExporterAgentID, ok := assertNodeExporterCreated(t, pmmAgentID)
			if ok {
				defer pmmapitests.RemoveAgents(t, nodeExporterAgentID)
			}
		})

		t.Run("With all fields", func(t *testing.T) {
			tt := pmmapitests.ExpectFailure(t, "https://jira.percona.com/browse/PMM-3982")
			defer tt.Check()

			nodeName := pmmapitests.TestString(t, "node-name")
			machineID := pmmapitests.TestString(t, "machine-id")
			nodeModel := pmmapitests.TestString(t, "node-model")
			body := node.RegisterBody{
				NodeName:     nodeName,
				NodeType:     pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
				MachineID:    machineID,
				NodeModel:    nodeModel,
				Az:           "eu",
				Region:       "us-west",
				Address:      "10.10.10.10",
				Distro:       "Linux",
				CustomLabels: map[string]string{"foo": "bar"},
			}
			nodeID, pmmAgentID := registerGenericNode(t, body)
			defer pmmapitests.RemoveNodes(t, nodeID)
			defer pmmapitests.RemoveAgents(t, pmmAgentID)

			// Check Node is created
			assertNodeCreated(tt, nodeID, nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:       nodeID,
					NodeName:     nodeName,
					MachineID:    machineID,
					NodeModel:    nodeModel,
					Az:           "eu",
					Region:       "us-west",
					Address:      "10.10.10.10",
					Distro:       "Linux",
					CustomLabels: map[string]string{"foo": "bar"},
				},
			})

			// Check PMM Agent is created
			assertPMMAgentCreated(t, nodeID, pmmAgentID)

			// Check Node Exporter is created
			nodeExporterAgentID, ok := assertNodeExporterCreated(t, pmmAgentID)
			if ok {
				defer pmmapitests.RemoveAgents(t, nodeExporterAgentID)
			}
		})

		t.Run("Re-register", func(t *testing.T) {
			t.Skip("Re-register logic is not defined yet. https://jira.percona.com/browse/PMM-3717")

			nodeName := pmmapitests.TestString(t, "node-name")
			nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
			})
			defer pmmapitests.RemoveNodes(t, nodeID)
			defer removePMMAgentWithSubAgents(t, pmmAgentID)

			// Check Node is created
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:   nodeID,
					NodeName: nodeName,
				},
			})

			// Re-register node
			machineID := pmmapitests.TestString(t, "machine-id")
			nodeModel := pmmapitests.TestString(t, "node-model")
			newNodeID, newPMMAgentID := registerGenericNode(t, node.RegisterBody{
				NodeName:     nodeName,
				NodeType:     pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
				MachineID:    machineID,
				NodeModel:    nodeModel,
				Az:           "eu",
				Region:       "us-west",
				Address:      "10.10.10.10",
				Distro:       "Linux",
				CustomLabels: map[string]string{"foo": "bar"},
			})
			if !assert.Equal(t, nodeID, newNodeID) {
				defer pmmapitests.RemoveNodes(t, newNodeID)
			}
			if !assert.Equal(t, pmmAgentID, newPMMAgentID) {
				defer pmmapitests.RemoveAgents(t, newPMMAgentID)
			}

			// Check Node fields is updated
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:       nodeID,
					NodeName:     nodeName,
					MachineID:    machineID,
					NodeModel:    nodeModel,
					Az:           "eu",
					Region:       "us-west",
					Address:      "10.10.10.10",
					Distro:       "Linux",
					CustomLabels: map[string]string{"foo": "bar"},
				},
			})
		})
	})

	t.Run("Container Node", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			nodeName := pmmapitests.TestString(t, "node-name")
			nodeID, pmmAgentID := registerContainerNode(t, node.RegisterBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterBodyNodeTypeCONTAINERNODE),
			})
			defer pmmapitests.RemoveNodes(t, nodeID)
			defer pmmapitests.RemoveAgents(t, pmmAgentID)

			// Check Node is created
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
				Container: &nodes.GetNodeOKBodyContainer{
					NodeID:   nodeID,
					NodeName: nodeName,
				},
			})

			// Check PMM Agent is created
			assertPMMAgentCreated(t, nodeID, pmmAgentID)

			// Check Node Exporter is created
			nodeExporterAgentID, ok := assertNodeExporterCreated(t, pmmAgentID)
			if ok {
				defer pmmapitests.RemoveAgents(t, nodeExporterAgentID)
			}
		})

		t.Run("With all fields", func(t *testing.T) {
			tt := pmmapitests.ExpectFailure(t, "https://jira.percona.com/browse/PMM-3982")
			defer tt.Check()

			nodeName := pmmapitests.TestString(t, "node-name")
			nodeModel := pmmapitests.TestString(t, "node-model")
			containerID := pmmapitests.TestString(t, "container-id")
			containerName := pmmapitests.TestString(t, "container-name")
			body := node.RegisterBody{
				NodeName:      nodeName,
				NodeType:      pointer.ToString(node.RegisterBodyNodeTypeCONTAINERNODE),
				NodeModel:     nodeModel,
				ContainerID:   containerID,
				ContainerName: containerName,
				Az:            "eu",
				Region:        "us-west",
				Address:       "10.10.10.10",
				CustomLabels:  map[string]string{"foo": "bar"},
			}
			nodeID, pmmAgentID := registerContainerNode(t, body)
			defer pmmapitests.RemoveNodes(t, nodeID)
			defer pmmapitests.RemoveAgents(t, pmmAgentID)

			// Check Node is created
			assertNodeCreated(tt, nodeID, nodes.GetNodeOKBody{
				Container: &nodes.GetNodeOKBodyContainer{
					NodeID:        nodeID,
					NodeName:      nodeName,
					NodeModel:     nodeModel,
					ContainerID:   containerID,
					ContainerName: containerName,
					Az:            "eu",
					Region:        "us-west",
					Address:       "10.10.10.10",
					CustomLabels:  map[string]string{"foo": "bar"},
				},
			})

			// Check PMM Agent is created
			assertPMMAgentCreated(t, nodeID, pmmAgentID)

			// Check Node Exporter is created
			nodeExporterAgentID, ok := assertNodeExporterCreated(t, pmmAgentID)
			if ok {
				defer pmmapitests.RemoveAgents(t, nodeExporterAgentID)
			}
		})

		t.Run("Re-register", func(t *testing.T) {
			t.Skip("Re-register logic is not defined yet. https://jira.percona.com/browse/PMM-3717")

			nodeName := pmmapitests.TestString(t, "node-name")
			nodeID, pmmAgentID := registerContainerNode(t, node.RegisterBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterBodyNodeTypeCONTAINERNODE),
			})
			defer pmmapitests.RemoveNodes(t, nodeID)
			defer removePMMAgentWithSubAgents(t, pmmAgentID)

			// Check Node is created
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:   nodeID,
					NodeName: nodeName,
				},
			})

			// Re-register node
			nodeModel := pmmapitests.TestString(t, "node-model")
			containerID := pmmapitests.TestString(t, "container-id")
			containerName := pmmapitests.TestString(t, "container-name")
			newNodeID, newPMMAgentID := registerContainerNode(t, node.RegisterBody{
				NodeName:      nodeName,
				NodeType:      pointer.ToString(node.RegisterBodyNodeTypeCONTAINERNODE),
				ContainerID:   containerID,
				ContainerName: containerName,
				NodeModel:     nodeModel,
				Az:            "eu",
				Region:        "us-west",
				Address:       "10.10.10.10",
				CustomLabels:  map[string]string{"foo": "bar"},
			})
			if !assert.Equal(t, nodeID, newNodeID) {
				defer pmmapitests.RemoveNodes(t, newNodeID)
			}
			if !assert.Equal(t, pmmAgentID, newPMMAgentID) {
				defer pmmapitests.RemoveAgents(t, newPMMAgentID)
			}

			// Check Node fields is updated
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
				Container: &nodes.GetNodeOKBodyContainer{
					NodeID:        nodeID,
					NodeName:      nodeName,
					ContainerID:   containerID,
					ContainerName: containerName,
					NodeModel:     nodeModel,
					Az:            "eu",
					Region:        "us-west",
					Address:       "10.10.10.10",
					CustomLabels:  map[string]string{"foo": "bar"},
				},
			})
		})
	})

	t.Run("Re-register node with different type", func(t *testing.T) {
		t.Skip("Re-register logic is not defined yet. https://jira.percona.com/browse/PMM-3717")

		nodeName := pmmapitests.TestString(t, "node-name")
		nodeID, pmmAgentID := registerGenericNode(t, node.RegisterBody{
			NodeName: nodeName,
			NodeType: pointer.ToString(node.RegisterBodyNodeTypeGENERICNODE),
		})
		defer pmmapitests.RemoveNodes(t, nodeID)
		defer removePMMAgentWithSubAgents(t, pmmAgentID)

		params := node.RegisterParams{
			Context: pmmapitests.Context,
			Body: node.RegisterBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterBodyNodeTypeCONTAINERNODE),
			},
		}
		registerOK, err := client.Default.Node.Register(&params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: ""})
		require.Nil(t, registerOK)
	})

	t.Run("Empty node name", func(t *testing.T) {
		params := node.RegisterParams{
			Context: pmmapitests.Context,
			Body:    node.RegisterBody{},
		}
		registerOK, err := client.Default.Node.Register(&params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "invalid field NodeName: value '' must not be an empty string"})
		require.Nil(t, registerOK)
	})

	t.Run("Unsupported node type", func(t *testing.T) {
		params := node.RegisterParams{
			Context: pmmapitests.Context,
			Body: node.RegisterBody{
				NodeName: pmmapitests.TestString(t, "node-name"),
			},
		}
		registerOK, err := client.Default.Node.Register(&params)
		pmmapitests.AssertEqualAPIError(t, err, pmmapitests.ServerResponse{Code: 400, Error: "unsupported node type"})
		require.Nil(t, registerOK)
	})
}
