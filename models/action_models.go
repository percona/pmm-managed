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

	"github.com/google/uuid"
	"gopkg.in/reform.v1"
)

//go:generate reform

// ActionResult describes an action result which is storing in persistent storage.
//reform:action_results
type ActionResult struct {
	ID         string    `reform:"id,pk"`
	PmmAgentID string    `reform:"pmm_agent_id"`
	Done       bool      `reform:"done"`
	Error      string    `reform:"error"`
	Output     string    `reform:"output"`
	CreatedAt  time.Time `reform:"created_at"`
	UpdatedAt  time.Time `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (s *ActionResult) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (s *ActionResult) BeforeUpdate() error {
	s.UpdatedAt = Now()
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (s *ActionResult) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return nil
}

// PtSummaryAction represents pt-summary domain model.
type PtSummaryAction struct {
	ID         string
	PMMAgentID string
	NodeID     string

	Args []string
}

// PtMySQLSummaryAction represents pt-mysql-summary domain model.
type PtMySQLSummaryAction struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Args []string
}

// MySQLExplainAction represents mysql-explain domain model.
type MySQLExplainAction struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// MySQLExplainJSONAction represents mysql-explain-json domain model.
type MySQLExplainJSONAction struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// getActionUUID generates action uuid.
// nolint: unused
func getActionUUID() string {
	return "/action_id/" + uuid.New().String()
}

// check interfaces
var (
	_ reform.BeforeInserter = (*ActionResult)(nil)
	_ reform.BeforeUpdater  = (*ActionResult)(nil)
	_ reform.AfterFinder    = (*ActionResult)(nil)
)
