// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package inventory

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	inventorypb "github.com/percona/pmm/api/inventory"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// NodesService works with inventory API Nodes.
type NodesService struct{}

// NewNodesService creates NodesService.
func NewNodesService() *NodesService {
	return &NodesService{}
}

// toInventoryNode converts database row to Inventory API Node.
func toInventoryNode(row *models.Node) (inventorypb.Node, error) {
	labels, err := row.GetCustomLabels()
	if err != nil {
		return nil, err
	}

	switch row.NodeType {
	case models.GenericNodeType:
		return &inventorypb.GenericNode{
			NodeId:        row.NodeID,
			NodeName:      row.NodeName,
			MachineId:     pointer.GetString(row.MachineID),
			Distro:        pointer.GetString(row.Distro),
			DistroVersion: pointer.GetString(row.DistroVersion),
			CustomLabels:  labels,
			Address:       pointer.GetString(row.Address),
		}, nil

	case models.ContainerNodeType:
		return &inventorypb.ContainerNode{
			NodeId:              row.NodeID,
			NodeName:            row.NodeName,
			MachineId:           pointer.GetString(row.MachineID),
			DockerContainerId:   pointer.GetString(row.DockerContainerID),
			DockerContainerName: pointer.GetString(row.DockerContainerName),
			CustomLabels:        labels,
		}, nil

	case models.RemoteNodeType:
		return &inventorypb.RemoteNode{
			NodeId:       row.NodeID,
			NodeName:     row.NodeName,
			CustomLabels: labels,
		}, nil

	case models.RemoteAmazonRDSNodeType:
		return &inventorypb.RemoteAmazonRDSNode{
			NodeId:       row.NodeID,
			NodeName:     row.NodeName,
			Instance:     pointer.GetString(row.Address),
			Region:       pointer.GetString(row.Region),
			CustomLabels: labels,
		}, nil

	default:
		panic(fmt.Errorf("unhandled Node type %s", row.NodeType))
	}
}

func toInventoryNodes(nodes []*models.Node) ([]inventorypb.Node, error) {
	var err error
	res := make([]inventorypb.Node, len(nodes))
	for i, n := range nodes {
		res[i], err = toInventoryNode(n)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// List selects all Nodes in a stable order.
func (ns *NodesService) List(ctx context.Context, q *reform.Querier) ([]inventorypb.Node, error) { //nolint:unparam
	nodes, err := models.FindAllNodes(q)
	if err != nil {
		return nil, err // TODO: Convert to gRPC errors
	}
	return toInventoryNodes(nodes)
}

// Get selects a single Node by ID.
func (ns *NodesService) Get(ctx context.Context, q *reform.Querier, id string) (inventorypb.Node, error) {
	node, err := models.FindNodeByID(q, id)
	if err != nil {
		return nil, err // TODO: Convert to gRPC errors
	}
	return toInventoryNode(node)
}

// Add inserts Node with given parameters. ID will be generated.
func (ns *NodesService) Add(ctx context.Context, q *reform.Querier, params *models.AddNodeParams) (inventorypb.Node, error) { //nolint:unparam
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// No hostname for Container, etc.
	node, err := models.AddNode(q, params)
	if err != nil {
		return nil, err // TODO: Convert to gRPC errors
	}
	return toInventoryNode(node)
}

// Remove deletes Node by ID.
//nolint:unparam
func (ns *NodesService) Remove(ctx context.Context, q *reform.Querier, id string) error { //nolint:unparam
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0.

	// TODO check absence of Services and Agents

	err := models.RemoveNode(q, id)
	if err != nil {
		return err // TODO: Convert to gRPC errors
	}

	return nil
}
