package management

import (
	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

func nodeID(tx *reform.TX, nodeID, nodeName string, addNodeParams *managementpb.AddNodeParams, address string) (string, error) {
	if err := validateNodeParamsOneOf(nodeID, nodeName, addNodeParams); err != nil {
		return "", err
	}
	switch {
	case nodeID != "":
		return nodeID, nil
	case nodeName != "":
		node, err := models.FindNodeByName(tx.Querier, nodeName)
		if err != nil {
			return "", err
		}
		return node.NodeID, nil
	case addNodeParams != nil:
		node, err := addNode(tx, addNodeParams, address)
		if err != nil {
			return "", err
		}
		nodeID = node.NodeID
	}
	return nodeID, nil
}

func addNode(tx *reform.TX, addNodeParams *managementpb.AddNodeParams, address string) (*models.Node, error) {
	nodeType, err := nodeType(addNodeParams.NodeType)
	if err != nil {
		return nil, err
	}
	node, err := models.CreateNode(tx.Querier, nodeType, &models.CreateNodeParams{
		NodeName:      addNodeParams.NodeName,
		MachineID:     pointer.ToStringOrNil(addNodeParams.MachineId),
		Distro:        addNodeParams.Distro,
		NodeModel:     addNodeParams.NodeModel,
		AZ:            addNodeParams.Az,
		ContainerID:   pointer.ToStringOrNil(addNodeParams.ContainerId),
		ContainerName: pointer.ToStringOrNil(addNodeParams.ContainerName),
		CustomLabels:  addNodeParams.CustomLabels,
		Address:       address,
		Region:        pointer.ToStringOrNil(addNodeParams.Region),
	})
	if err != nil {
		return nil, err
	}
	return node, nil
}

func nodeType(inputNodeType inventorypb.NodeType) (models.NodeType, error) {
	var nodeType models.NodeType
	switch inputNodeType {
	case inventorypb.NodeType_GENERIC_NODE:
		nodeType = models.GenericNodeType
	case inventorypb.NodeType_CONTAINER_NODE:
		nodeType = models.ContainerNodeType
	case inventorypb.NodeType_REMOTE_NODE:
		nodeType = models.RemoteNodeType
	default:
		return "", status.Errorf(codes.InvalidArgument, "Unsupported Node type %q.", inputNodeType)
	}
	return nodeType, nil
}

func validateNodeParamsOneOf(nodeID, nodeName string, addNodeParams *managementpb.AddNodeParams) error {
	got := 0
	if nodeID != "" {
		got++
	}
	if nodeName != "" {
		got++
	}
	if addNodeParams != nil {
		got++
	}
	if got != 1 {
		return status.Errorf(codes.InvalidArgument, "expected only one param; node id, node name or register node params")
	}
	return nil
}
