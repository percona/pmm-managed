package inventory

import (
	"fmt"
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/percona/pmm/api/inventory/json/client/nodes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestNodes(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Parallel()

		remoteNode := addRemoteNode(t, pmmapitests.TestString(t, "Test Remote Node for List"))
		remoteNodeID := remoteNode.Remote.NodeID
		defer removeNodes(t, remoteNodeID)
		genericNodeID := addGenericNode(t, pmmapitests.TestString(t, "Test Generic Node for List")).NodeID
		require.NotEmpty(t, genericNodeID)
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
		t.Parallel()

		nodeName := pmmapitests.TestString(t, "TestGenericNode")
		nodeID := addGenericNode(t, nodeName).NodeID
		require.NotEmpty(t, nodeID)
		defer removeNodes(t, nodeID)

		expectedResponse := &nodes.GetNodeOK{
			Payload: &nodes.GetNodeOKBody{
				Generic: &nodes.GetNodeOKBodyGeneric{
					NodeID:   nodeID,
					NodeName: nodeName,
				},
			},
		}

		params := &nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: nodeID},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.GetNode(params)
		require.NoError(t, err)
		require.Equal(t, res, expectedResponse)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		params := &nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{NodeID: "pmm-not-found"},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.GetNode(params)
		assertEqualAPIError(t, err, ServerResponse{404, "Node with ID \"pmm-not-found\" not found."})
		assert.Nil(t, res)
	})

	t.Run("EmptyNodeID", func(t *testing.T) {
		t.Parallel()

		params := &nodes.GetNodeParams{
			Body:    nodes.GetNodeBody{},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.GetNode(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeId: value '' must not be an empty string"})
		assert.Nil(t, res)
	})
}

func TestGenericNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		nodeName := pmmapitests.TestString(t, "Test Generic Node")
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
		assertEqualAPIError(t, err, ServerResponse{409, fmt.Sprintf("Node with name %q already exists.", nodeName)})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.Generic.NodeID)
		}
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		t.Parallel()

		params := &nodes.AddGenericNodeParams{
			Body:    nodes.AddGenericNodeBody{NodeName: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddGenericNode(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.Generic.NodeID)
		}
	})
}

func TestContainerNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		nodeName := pmmapitests.TestString(t, "Test Container Node")
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
					NodeID:              res.Payload.Container.NodeID,
					NodeName:            nodeName,
					DockerContainerID:   "docker-id",
					DockerContainerName: "docker-name",
					MachineID:           "machine-id",
				},
			},
		}
		require.Equal(t, expectedResponse, getNodeRes)

		// Check duplicates.
		res, err = client.Default.Nodes.AddContainerNode(params)
		assertEqualAPIError(t, err, ServerResponse{409, fmt.Sprintf("Node with name %q already exists.", nodeName)})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.Container.NodeID)
		}
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		t.Parallel()

		params := &nodes.AddContainerNodeParams{
			Body:    nodes.AddContainerNodeBody{NodeName: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddContainerNode(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.Container.NodeID)
		}
	})
}

func TestRemoteNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		nodeName := pmmapitests.TestString(t, "Test Remote Node")
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
		assertEqualAPIError(t, err, ServerResponse{409, fmt.Sprintf("Node with name %q already exists.", nodeName)})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.Remote.NodeID)
		}
	})

	t.Run("AddNameEmpty", func(t *testing.T) {
		t.Parallel()

		params := &nodes.AddRemoteNodeParams{
			Body:    nodes.AddRemoteNodeBody{NodeName: ""},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteNode(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.Remote.NodeID)
		}
	})
}

func TestRemoteAmazonRDSNode(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Parallel()

		nodeName := pmmapitests.TestString(t, "Test RemoteAmazonRDS Node")
		instanceName := pmmapitests.TestString(t, "some-instance")
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
		assertEqualAPIError(t, err, ServerResponse{409, fmt.Sprintf("Node with name %q already exists.", nodeName)})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.RemoteAmazonRDS.NodeID)
		}
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
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field NodeName: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.RemoteAmazonRDS.NodeID)
		}
	})

	t.Run("AddInstanceEmpty", func(t *testing.T) {
		t.Parallel()

		params := &nodes.AddRemoteAmazonRDSNodeParams{
			Body: nodes.AddRemoteAmazonRDSNodeBody{
				NodeName: pmmapitests.TestString(t, "Remote AmazonRDSNode without instance"),
				Region:   "us-west-1",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field Instance: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.RemoteAmazonRDS.NodeID)
		}
	})

	t.Run("AddRegionEmpty", func(t *testing.T) {
		t.Parallel()

		params := &nodes.AddRemoteAmazonRDSNodeParams{
			Body: nodes.AddRemoteAmazonRDSNodeBody{
				NodeName: pmmapitests.TestString(t, "Remote AmazonRDSNode without instance"),
				Instance: "instance-without-region",
			},
			Context: pmmapitests.Context,
		}
		res, err := client.Default.Nodes.AddRemoteAmazonRDSNode(params)
		assertEqualAPIError(t, err, ServerResponse{400, "invalid field Region: value '' must not be an empty string"})
		if !assert.Nil(t, res) {
			removeNodes(t, res.Payload.RemoteAmazonRDS.NodeID)
		}
	})
}
