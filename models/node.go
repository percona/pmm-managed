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
	"time"

	"gopkg.in/reform.v1"
)

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
	now := time.Now().Truncate(time.Microsecond).UTC()
	nr.CreatedAt = now
	// nr.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (nr *NodeRow) BeforeUpdate() error {
	// now := time.Now().Truncate(time.Microsecond).UTC()
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
