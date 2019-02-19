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
	"time"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

// NodesForAgent returns all Nodes for which Agent with given ID provides insights.
func NodesForAgent(q *reform.Querier, agentID string) ([]*NodeRow, error) {
	structs, err := q.FindAllFrom(AgentNodeView, "agent_id", agentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Node IDs")
	}

	nodeIDs := make([]interface{}, len(structs))
	for i, s := range structs {
		nodeIDs[i] = s.(*AgentNode).NodeID
	}
	if len(nodeIDs) == 0 {
		return []*NodeRow{}, nil
	}

	p := strings.Join(q.Placeholders(1, len(nodeIDs)), ", ")
	tail := fmt.Sprintf("WHERE node_id IN (%s) ORDER BY node_id", p) //nolint:gosec
	structs, err = q.SelectAllFrom(NodeRowTable, tail, nodeIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Nodes")
	}

	res := make([]*NodeRow, len(structs))
	for i, s := range structs {
		res[i] = s.(*NodeRow)
	}
	return res, nil
}

//go:generate reform

// NodeType represents Node type as stored in database.
type NodeType string

// Node types.
const (
	PMMServerNodeID string = "pmm-server" // FIXME remove

	PMMServerNodeType NodeType = "pmm-server" // FIXME remove

	GenericNodeType         NodeType = "generic"
	ContainerNodeType       NodeType = "container"
	RemoteNodeType          NodeType = "remote"
	RemoteAmazonRDSNodeType NodeType = "remote-amazon-rds"
)

// NodeRow represents Node as stored in database.
//reform:nodes
type NodeRow struct {
	NodeID    string    `reform:"node_id,pk"`
	NodeType  NodeType  `reform:"node_type"`
	NodeName  string    `reform:"node_name"`
	MachineID *string   `reform:"machine_id"`
	CreatedAt time.Time `reform:"created_at"`
	// UpdatedAt time.Time `reform:"updated_at"`

	Distro        *string `reform:"distro"`
	DistroVersion *string `reform:"distro_version"`

	DockerContainerID   *string `reform:"docker_container_id"`
	DockerContainerName *string `reform:"docker_container_name"`

	Instance *string `reform:"instance"`
	Region   *string `reform:"region"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (nr *NodeRow) BeforeInsert() error {
	now := Now()
	nr.CreatedAt = now
	// nr.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (nr *NodeRow) BeforeUpdate() error {
	// now := Now()
	// nr.UpdatedAt = now
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (nr *NodeRow) AfterFind() error {
	nr.CreatedAt = nr.CreatedAt.UTC()
	// nr.UpdatedAt = nr.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*NodeRow)(nil)
	_ reform.BeforeUpdater  = (*NodeRow)(nil)
	_ reform.AfterFinder    = (*NodeRow)(nil)
)
