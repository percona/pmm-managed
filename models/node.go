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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	inventorypb "github.com/percona/pmm/api/inventory"
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
	NodeType            NodeType
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
func AddNode(q *reform.Querier, params *AddNodeParams) (*Node, error) {
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
		NodeType:            params.NodeType,
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

// RemoveNode removes single node prom persistent store.
func RemoveNode(q *reform.Querier, id string) error {
	err := q.Delete(&Node{NodeID: id})
	if err == reform.ErrNoRows {
		return status.Errorf(codes.NotFound, "Node with ID %q not found.", id)
	}
	return nil
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

// ToInventoryNode converts database row to Inventory API Node.
func ToInventoryNode(row *Node) (inventorypb.Node, error) {
	labels, err := row.GetCustomLabels()
	if err != nil {
		return nil, err
	}

	switch row.NodeType {
	case GenericNodeType:
		return &inventorypb.GenericNode{
			NodeId:        row.NodeID,
			NodeName:      row.NodeName,
			MachineId:     pointer.GetString(row.MachineID),
			Distro:        pointer.GetString(row.Distro),
			DistroVersion: pointer.GetString(row.DistroVersion),
			CustomLabels:  labels,
			Address:       pointer.GetString(row.Address),
		}, nil

	case ContainerNodeType:
		return &inventorypb.ContainerNode{
			NodeId:              row.NodeID,
			NodeName:            row.NodeName,
			MachineId:           pointer.GetString(row.MachineID),
			DockerContainerId:   pointer.GetString(row.DockerContainerID),
			DockerContainerName: pointer.GetString(row.DockerContainerName),
			CustomLabels:        labels,
		}, nil

	case RemoteNodeType:
		return &inventorypb.RemoteNode{
			NodeId:       row.NodeID,
			NodeName:     row.NodeName,
			CustomLabels: labels,
		}, nil

	case RemoteAmazonRDSNodeType:
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

// ToInventoryNode converts database rows to Inventory API Nodes.
func ToInventoryNodes(nodes []*Node) ([]inventorypb.Node, error) {
	var err error
	res := make([]inventorypb.Node, len(nodes))
	for i, n := range nodes {
		res[i], err = ToInventoryNode(n)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

//go:generate reform

// NodeType represents Node type as stored in database.
type NodeType string

// Node types.
const (
	GenericNodeType         NodeType = "generic"
	ContainerNodeType       NodeType = "container"
	RemoteNodeType          NodeType = "remote"
	RemoteAmazonRDSNodeType NodeType = "remote-amazon-rds"
)

// PMMServerNodeID is a special Node ID representing PMM Server Node.
const PMMServerNodeID string = "pmm-server"

// Node represents Node as stored in database.
//reform:nodes
type Node struct {
	NodeID       string    `reform:"node_id,pk"`
	NodeType     NodeType  `reform:"node_type"`
	NodeName     string    `reform:"node_name"`
	MachineID    *string   `reform:"machine_id"` // nil means "unknown"; non-nil value must be unique
	CustomLabels []byte    `reform:"custom_labels"`
	Address      *string   `reform:"address"` // nil means "unknown"; also stores Remote instance
	CreatedAt    time.Time `reform:"created_at"`
	UpdatedAt    time.Time `reform:"updated_at"`

	Distro        *string `reform:"distro"`
	DistroVersion *string `reform:"distro_version"`

	DockerContainerID   *string `reform:"docker_container_id"` // nil means "unknown"; non-nil value must be unique
	DockerContainerName *string `reform:"docker_container_name"`

	Region *string `reform:"region"` // nil means "not Remote"; non-nil value must be unique in combination with instance/address
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (s *Node) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (s *Node) BeforeUpdate() error {
	s.UpdatedAt = Now()
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (s *Node) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// GetCustomLabels decodes custom labels.
func (s *Node) GetCustomLabels() (map[string]string, error) {
	if len(s.CustomLabels) == 0 {
		return nil, nil
	}
	m := make(map[string]string)
	if err := json.Unmarshal(s.CustomLabels, &m); err != nil {
		return nil, errors.Wrap(err, "failed to decode custom labels")
	}
	return m, nil
}

// SetCustomLabels encodes custom labels.
func (s *Node) SetCustomLabels(m map[string]string) error {
	if len(m) == 0 {
		s.CustomLabels = nil
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "failed to encode custom labels")
	}
	s.CustomLabels = b
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Node)(nil)
	_ reform.BeforeUpdater  = (*Node)(nil)
	_ reform.AfterFinder    = (*Node)(nil)
)
