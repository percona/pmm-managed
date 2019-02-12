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
	"time"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

//go:generate reform

// AgentsRunningOnNode returns all Agents running on Node.
// TODO Remove after https://jira.percona.com/browse/PMM-3478.
func AgentsRunningOnNode(q *reform.Querier, nodeID string) ([]*AgentRow, error) {
	tail := fmt.Sprintf("WHERE runs_on_node_id = %s ORDER BY agent_id", q.Placeholder(1)) //nolint:gosec
	structs, err := q.SelectAllFrom(AgentRowTable, tail, nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agents")
	}

	res := make([]*AgentRow, len(structs))
	for i, s := range structs {
		res[i] = s.(*AgentRow)
	}
	return res, nil
}

// AgentType represents Agent type as stored in database.
type AgentType string

// Agent types.
const (
	PMMAgentType         AgentType = "pmm-agent"
	NodeExporterType     AgentType = "node_exporter"
	MySQLdExporterType   AgentType = "mysqld_exporter"
	RDSExporterType      AgentType = "rds_exporter"
	ExternalExporterType AgentType = "external"
)

// AgentRow represents Agent as stored in database.
//reform:agents
type AgentRow struct {
	AgentID      string    `reform:"agent_id,pk"`
	AgentType    AgentType `reform:"agent_type"`
	RunsOnNodeID string    `reform:"runs_on_node_id"`
	CreatedAt    time.Time `reform:"created_at"`
	// UpdatedAt    time.Time `reform:"updated_at"`

	Version    *string `reform:"version"`
	Status     *string `reform:"status"`
	ListenPort *uint16 `reform:"listen_port"`

	Username *string `reform:"username"`
	Password *string `reform:"password"`

	MetricsURL *string `reform:"metrics_url"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (ar *AgentRow) BeforeInsert() error {
	now := time.Now().Truncate(time.Microsecond).UTC()
	ar.CreatedAt = now
	// ar.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (ar *AgentRow) BeforeUpdate() error {
	// now := time.Now().Truncate(time.Microsecond).UTC()
	// ar.UpdatedAt = now
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (ar *AgentRow) AfterFind() error {
	ar.CreatedAt = ar.CreatedAt.UTC()
	// ar.UpdatedAt = ar.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*AgentRow)(nil)
	_ reform.BeforeUpdater  = (*AgentRow)(nil)
	_ reform.AfterFinder    = (*AgentRow)(nil)
)
