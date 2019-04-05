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

package models

import (
	"fmt"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

func checkIsUniqueNodeID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Node ID")
	}

	row := &Node{NodeID: id}
	switch err := q.Reload(row); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Node with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func checkUniqueName(q *reform.Querier, name string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "Empty Node name.")
	}

	_, err := q.FindOneFrom(NodeTable, "node_name", name)
	switch err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Node with name %q already exists.", name)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func checkUniqueNodeInstanceRegion(q *reform.Querier, instance, region string) error {
	if instance == "" {
		return status.Error(codes.InvalidArgument, "Empty Node instance.")
	}
	if region == "" {
		return status.Error(codes.InvalidArgument, "Empty Node region.")
	}

	tail := fmt.Sprintf("WHERE address = %s AND region = %s LIMIT 1", q.Placeholder(1), q.Placeholder(2)) //nolint:gosec
	_, err := q.SelectOneFrom(NodeTable, tail, instance, region)
	switch err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Node with instance %q and region %q already exists.", instance, region)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// AddNodeParams contains parameters for adding Nodes.
type AddNodeParams struct {
	NodeName            string
	MachineID           *string
	Distro              *string
	DistroVersion       *string
	DockerContainerID   *string
	DockerContainerName *string
	CustomLabels        map[string]string
	Address             *string
	Region              *string
}

// AddNode adds new node to persistent store.
func AddNode(q *reform.Querier, nodeType NodeType, params *AddNodeParams) (*Node, error) {
	id := "/node_id/" + uuid.New().String()
	if err := checkIsUniqueNodeID(q, id); err != nil {
		return nil, err
	}

	if err := checkUniqueName(q, params.NodeName); err != nil {
		return nil, err
	}

	if params.Address != nil && params.Region != nil {
		if err := checkUniqueNodeInstanceRegion(q, *params.Address, *params.Region); err != nil {
			return nil, err
		}
	}

	row := &Node{
		NodeID:              id,
		NodeType:            nodeType,
		NodeName:            params.NodeName,
		MachineID:           params.MachineID,
		Distro:              params.Distro,
		DistroVersion:       params.DistroVersion,
		DockerContainerID:   params.DockerContainerID,
		DockerContainerName: params.DockerContainerName,
		Address:             params.Address,
		Region:              params.Region,
	}
	if err := row.SetCustomLabels(params.CustomLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(row); err != nil {
		return nil, err
	}

	return row, nil
}

// FindAllNodes finds all nodes and loads it from persistent store.
func FindAllNodes(q *reform.Querier) ([]*Node, error) {
	structs, err := q.SelectAllFrom(NodeTable, "ORDER BY node_id")
	if err != nil {
		return nil, err
	}

	nodes := make([]*Node, len(structs))
	for i, s := range structs {
		nodes[i] = s.(*Node)
	}

	return nodes, nil
}

// FindNodeByID finds a node by ID and loads it from persistent store.
func FindNodeByID(q *reform.Querier, id string) (*Node, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Node ID.")
	}

	row := &Node{NodeID: id}
	switch err := q.Reload(row); err {
	case nil:
		return row, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Node with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// UpdateNodeParams describe editable node parameters.
type UpdateNodeParams struct {
	Address         string
	MachineID       string
	CustomLabels    map[string]string
	RemoveLabels    bool
	RemoveMachineID bool
}

// UpdateNode updates Node.
func UpdateNode(q *reform.Querier, nodeID string, params *UpdateNodeParams) (*Node, error) {
	row, err := FindNodeByID(q, nodeID)
	if err != nil {
		return nil, err
	}

	if params.Address != "" {
		row.Address = &params.Address
	}

	if params.RemoveLabels {
		row.CustomLabels = nil
	} else {
		err := row.SetCustomLabels(params.CustomLabels)
		if err != nil {
			return nil, err
		}
	}

	if params.RemoveMachineID {
		row.MachineID = nil
	} else {
		row.MachineID = pointer.ToStringOrNil(params.MachineID)
	}

	if err := q.Update(row); err != nil {
		return nil, err
	}

	return row, nil
}

// NodesForAgent returns all Nodes for which Agent with given ID provides insights.
func NodesForAgent(q *reform.Querier, agentID string) ([]*Node, error) {
	structs, err := q.FindAllFrom(AgentNodeView, "agent_id", agentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Node IDs")
	}

	nodeIDs := make([]interface{}, len(structs))
	for i, s := range structs {
		nodeIDs[i] = s.(*AgentNode).NodeID
	}
	if len(nodeIDs) == 0 {
		return []*Node{}, nil
	}

	p := strings.Join(q.Placeholders(1, len(nodeIDs)), ", ")
	tail := fmt.Sprintf("WHERE node_id IN (%s) ORDER BY node_id", p) //nolint:gosec
	structs, err = q.SelectAllFrom(NodeTable, tail, nodeIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Nodes")
	}

	res := make([]*Node, len(structs))
	for i, s := range structs {
		res[i] = s.(*Node)
	}
	return res, nil
}

// RemoveNode removes single node prom persistent store.
func RemoveNode(q *reform.Querier, id string) error {
	err := q.Delete(&Node{NodeID: id})
	if err == reform.ErrNoRows {
		return status.Errorf(codes.NotFound, "Node with ID %q not found.", id)
	}
	return nil
}
