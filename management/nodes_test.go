package management

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/inventorypb/json/client/nodes"
	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestNodeRegister(t *testing.T) {
	t.Run("Generic Node", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			nodeName := pmmapitests.TestString(t, "node-name")
			nodeID, pmmAgentID := registerGenericNode(t, node.RegisterNodeBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterNodeBodyNodeTypeGENERICNODE),
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
			nodeName := pmmapitests.TestString(t, "node-name")
			machineID := pmmapitests.TestString(t, "machine-id")
			nodeModel := pmmapitests.TestString(t, "node-model")
			body := node.RegisterNodeBody{
				NodeName:     nodeName,
				NodeType:     pointer.ToString(node.RegisterNodeBodyNodeTypeGENERICNODE),
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
			nodeID, pmmAgentID := registerGenericNode(t, node.RegisterNodeBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterNodeBodyNodeTypeGENERICNODE),
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
			newNodeID, newPMMAgentID := registerGenericNode(t, node.RegisterNodeBody{
				NodeName:     nodeName,
				NodeType:     pointer.ToString(node.RegisterNodeBodyNodeTypeGENERICNODE),
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
			nodeID, pmmAgentID := registerContainerNode(t, node.RegisterNodeBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterNodeBodyNodeTypeCONTAINERNODE),
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
			nodeName := pmmapitests.TestString(t, "node-name")
			nodeModel := pmmapitests.TestString(t, "node-model")
			containerID := pmmapitests.TestString(t, "container-id")
			containerName := pmmapitests.TestString(t, "container-name")
			body := node.RegisterNodeBody{
				NodeName:      nodeName,
				NodeType:      pointer.ToString(node.RegisterNodeBodyNodeTypeCONTAINERNODE),
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
			assertNodeCreated(t, nodeID, nodes.GetNodeOKBody{
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
			nodeID, pmmAgentID := registerContainerNode(t, node.RegisterNodeBody{
				NodeName: nodeName,
				NodeType: pointer.ToString(node.RegisterNodeBodyNodeTypeCONTAINERNODE),
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
			newNodeID, newPMMAgentID := registerContainerNode(t, node.RegisterNodeBody{
				NodeName:      nodeName,
				NodeType:      pointer.ToString(node.RegisterNodeBodyNodeTypeCONTAINERNODE),
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

	t.Run("Empty node name", func(t *testing.T) {
		params := node.RegisterNodeParams{
			Context: pmmapitests.Context,
			Body:    node.RegisterNodeBody{},
		}
		registerOK, err := client.Default.Node.RegisterNode(&params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field NodeName: value '' must not be an empty string")
		require.Nil(t, registerOK)
	})

	t.Run("Unsupported node type", func(t *testing.T) {
		params := node.RegisterNodeParams{
			Context: pmmapitests.Context,
			Body: node.RegisterNodeBody{
				NodeName: pmmapitests.TestString(t, "node-name"),
			},
		}
		registerOK, err := client.Default.Node.RegisterNode(&params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Unsupported Node type "NODE_TYPE_INVALID".`)
		require.Nil(t, registerOK)
	})
}
