package inventory

import (
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/percona/pmm/api/inventory/json/client/nodes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Percona-Lab/pmm-api-tests"
)

func TestNodes(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		remoteNode := addRemoteNode(t, withUUID(t, "Test Remote Node for List"))
		remoteNodeID := remoteNode.Remote.NodeID
		defer removeNodes(t, remoteNodeID)
		genericNode := addGenericNode(t, withUUID(t, "Test Generic Node for List"))
		genericNodeID := genericNode.Generic.NodeID
		defer removeNodes(t, genericNodeID)

		res, err := client.Default.Nodes.ListNodes(nil)
		require.NoError(t, err)
		require.NotZerof(t, len(res.Payload.Generic), "There should be at least one node")
		require.Conditionf(t, func() (success bool) {
			for _, v := range res.Payload.Generic {
				if v.NodeID == genericNodeID {
					return true
				}
			}
			return false
		}, "There should be generic node with id `%s`", genericNodeID)
		require.NotZerof(t, len(res.Payload.Remote), "There should be at least one node")
		require.Conditionf(t, func() (success bool) {
			for _, v := range res.Payload.Remote {
				if v.NodeID == remoteNodeID {
					return true
				}
			}
			return false
		}, "There should be remote node with id `%s`", remoteNodeID)
	})
}

func TestGetNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := withUUID(t, "TestGenericNode")
		node := addGenericNode(t, nodeName)
		defer removeNodes(t, node.Generic.NodeID)

		expectedResponse := &nodes.GetNodeOK{
			Payload: &nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:   node.Generic.NodeID,
					NodeName: nodeName,
				},
			},
		}

		params := &nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: node.Generic.NodeID},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.GetNode(params)
		require.NoError(t, err)
		require.Equal(t, res, expectedResponse)
	})

	t.Run("NotFound", func(t *testing.T) {
		params := &nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: "pmm-not-found"},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.GetNode(params)
		assertEqualAPIError(t, err, 404)
		assert.Nil(t, res)
	})

	t.Run("EmptyNodeID", func(t *testing.T) {
		params := &nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.GetNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})
}

func TestGenericNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := withUUID(t, "Test Generic Node")
		params := &nodes.AddGenericNodeParams{
			Body:    nodes.AddGenericNodeBody{NodeName: nodeName},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddGenericNode(params)
		assert.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Payload.Generic)
		nodeID := res.Payload.Generic.NodeID
		defer removeNodes(t, nodeID)

		// Check node exists in DB.
		getNodeRes, err := client.Default.Nodes.GetNode(&nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: nodeID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedResponse := &nodes.GetNodeOK{
			Payload: &nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:   res.Payload.Generic.NodeID,
					NodeName: nodeName,
				},
			},
		}
		require.Equal(t, expectedResponse, getNodeRes)

		// Check duplicates.
		res, err = client.Default.Nodes.AddGenericNode(params)
		assertEqualAPIError(t, err, 409)
		assert.Nil(t, res)

		// Change node.
		changedNodeName := withUUID(t, "Changed Generic Node")
		changeRes, err := client.Default.Nodes.ChangeGenericNode(&nodes.ChangeGenericNodeParams{
			Body:    nodes.ChangeGenericNodeBody{NodeID: nodeID, NodeName: changedNodeName},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedChangeResponse := &nodes.ChangeGenericNodeOK{
			Payload: &nodes.ChangeGenericNodeOKBody{
				Generic: &nodes.ChangeGenericNodeOKBodyGeneric{
					NodeID:   nodeID,
					NodeName: changedNodeName,
				},
			},
		}
		require.Equal(t, expectedChangeResponse, changeRes)
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		params := &nodes.AddGenericNodeParams{
			Body:    nodes.AddGenericNodeBody{NodeName: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddGenericNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})
}

func TestContainerNode(t *testing.T) {
	t.Skip("Haven't implemented yet.")
	t.Run("Basic", func(t *testing.T) {
		nodeName := withUUID(t, "Test Container Node")
		params := &nodes.AddContainerNodeParams{
			Body: nodes.AddContainerNodeBody{
				NodeName:            nodeName,
				DockerContainerID:   "docker-id",
				DockerContainerName: "docker-name",
				MachineID:           "machine-id",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddContainerNode(params)
		require.NoError(t, err)
		require.NotNil(t, res.Payload.Container)
		nodeID := res.Payload.Container.NodeID
		defer removeNodes(t, nodeID)

		// Check node exists in DB.
		getNodeRes, err := client.Default.Nodes.GetNode(&nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: nodeID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedResponse := &nodes.GetNodeOK{
			Payload: &nodes.GetNodeOKBody{
				Container: &nodes.GetNodeOKBodyContainer{
					NodeID:   res.Payload.Container.NodeID,
					NodeName: nodeName,
				},
			},
		}
		require.Equal(t, expectedResponse, getNodeRes)

		// Check duplicates.
		res, err = client.Default.Nodes.AddContainerNode(params)
		assertEqualAPIError(t, err, 409)
		assert.Nil(t, res)

		// Change node.
		changedNodeName := withUUID(t, "Changed Container Node")
		changeRes, err := client.Default.Nodes.ChangeContainerNode(&nodes.ChangeContainerNodeParams{
			Body:    nodes.ChangeContainerNodeBody{NodeID: nodeID, NodeName: changedNodeName},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedChangeResponse := &nodes.ChangeContainerNodeOK{
			Payload: &nodes.ChangeContainerNodeOKBody{
				Container: &nodes.ChangeContainerNodeOKBodyContainer{
					NodeID:   nodeID,
					NodeName: changedNodeName,
				},
			},
		}
		require.Equal(t, expectedChangeResponse, changeRes)
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		params := &nodes.AddContainerNodeParams{
			Body:    nodes.AddContainerNodeBody{NodeName: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddContainerNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})
}

func TestRemoteNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := withUUID(t, "Test Remote Node")
		params := &nodes.AddRemoteNodeParams{
			Body: nodes.AddRemoteNodeBody{
				NodeName: nodeName,
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteNode(params)
		require.NoError(t, err)
		require.NotNil(t, res.Payload.Remote)
		nodeID := res.Payload.Remote.NodeID
		defer removeNodes(t, nodeID)

		// Check node exists in DB.
		getNodeRes, err := client.Default.Nodes.GetNode(&nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: nodeID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedResponse := &nodes.GetNodeOK{
			Payload: &nodes.GetNodeOKBody{
				Remote: &nodes.GetNodeOKBodyRemote{
					NodeID:   res.Payload.Remote.NodeID,
					NodeName: nodeName,
				},
			},
		}
		require.Equal(t, expectedResponse, getNodeRes)

		// Check duplicates.
		res, err = client.Default.Nodes.AddRemoteNode(params)
		assertEqualAPIError(t, err, 409)
		assert.Nil(t, res)

		// Change node.
		changedNodeName := withUUID(t, "Changed Remote Node")
		changeRes, err := client.Default.Nodes.ChangeRemoteNode(&nodes.ChangeRemoteNodeParams{
			Body:    nodes.ChangeRemoteNodeBody{NodeID: nodeID, NodeName: changedNodeName},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedChangeResponse := &nodes.ChangeRemoteNodeOK{
			Payload: &nodes.ChangeRemoteNodeOKBody{
				Remote: &nodes.ChangeRemoteNodeOKBodyRemote{
					NodeID:   nodeID,
					NodeName: changedNodeName,
				},
			},
		}
		require.Equal(t, expectedChangeResponse, changeRes)
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		params := &nodes.AddRemoteNodeParams{
			Body:    nodes.AddRemoteNodeBody{NodeName: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})
}

func TestRemoteAmazonRDSNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		nodeName := withUUID(t, "Test RemoteAmazonRDS Node")
		instanceName := withUUID(t, "some-instance")
		params := &nodes.AddRemoteAmazonRDSNodeParams{
			Body: nodes.AddRemoteAmazonRDSNodeBody{
				NodeName: nodeName,
				Instance: instanceName,
				Region:   "us-east-1",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		require.NoError(t, err)
		require.NotNil(t, res.Payload.RemoteAmazonRDS)
		nodeID := res.Payload.RemoteAmazonRDS.NodeID
		defer removeNodes(t, nodeID)

		// Check if the node saved in PMM-Managed.
		getNodeRes, err := client.Default.Nodes.GetNode(&nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: nodeID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedResponse := &nodes.GetNodeOK{
			Payload: &nodes.GetNodeOKBody{
				RemoteAmazonRDS: &nodes.GetNodeOKBodyRemoteAmazonRDS{
					NodeID:   nodeID,
					NodeName: nodeName,
					Region:   "us-east-1",
					Instance: instanceName,
				},
			},
		}
		assert.Equal(t, expectedResponse, getNodeRes)

		// Check duplicates.
		res, err = client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		assertEqualAPIError(t, err, 409)
		assert.Nil(t, res)

		// Change node.
		changedNodeName := withUUID(t, "Changed RemoteAmazonRDS Node")
		changeRes, err := client.Default.Nodes.ChangeRemoteAmazonRDSNode(&nodes.ChangeRemoteAmazonRDSNodeParams{
			Body:    nodes.ChangeRemoteAmazonRDSNodeBody{NodeID: nodeID, NodeName: changedNodeName},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		expectedChangeResponse := &nodes.ChangeRemoteAmazonRDSNodeOK{
			Payload: &nodes.ChangeRemoteAmazonRDSNodeOKBody{
				RemoteAmazonRDS: &nodes.ChangeRemoteAmazonRDSNodeOKBodyRemoteAmazonRDS{
					NodeID:   nodeID,
					NodeName: changedNodeName,
					Region:   "us-east-1",
					Instance: instanceName,
				},
			},
		}
		require.Equal(t, expectedChangeResponse, changeRes)
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		params := &nodes.AddRemoteAmazonRDSNodeParams{
			Body: nodes.AddRemoteAmazonRDSNodeBody{
				NodeName: "",
				Instance: "some-instance-without-name",
				Region:   "us-east-1",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})

	t.Run("AddInstanceEmpty", func(t *testing.T) {
		params := &nodes.AddRemoteAmazonRDSNodeParams{
			Body: nodes.AddRemoteAmazonRDSNodeBody{
				NodeName: withUUID(t, "Remote AmazonRDSNode without instance"),
				Region:   "us-west-1",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})

	t.Run("AddRegionEmpty", func(t *testing.T) {
		params := &nodes.AddRemoteAmazonRDSNodeParams{
			Body: nodes.AddRemoteAmazonRDSNodeBody{
				NodeName: withUUID(t, "Remote AmazonRDSNode without instance"),
				Instance: "instance-without-region",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		assertEqualAPIError(t, err, 400)
		assert.Nil(t, res)
	})
}
