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

	"github.com/google/uuid"
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

	Config             string
	Help               bool
	ReadSamples        string
	SaveSamples        string
	Sleep              uint32
	SummarizeMounts    bool
	SummarizeNetwork   bool
	SummarizeProcesses bool
	Version            bool
}

// Args returns arguments slice for pmm-agent actions implementation.
func (s *PtSummaryAction) Args() []string {
	var args []string
	if s.Config != "" {
		args = append(args, "--config", s.Config)
	}
	if s.Version {
		args = append(args, "--version")
	}
	if s.Help {
		args = append(args, "--help")
	}
	if s.ReadSamples != "" {
		args = append(args, "--read-samples", s.ReadSamples)
	}
	if s.SaveSamples != "" {
		args = append(args, "--save-samples", s.SaveSamples)
	}
	if s.Sleep > 0 {
		args = append(args, "--sleep", fmt.Sprintf("%d", s.Sleep))
	}
	if s.SummarizeMounts {
		args = append(args, "--summarize-mounts")
	}
	if s.SummarizeNetwork {
		args = append(args, "--summarize-network")
	}
	if s.SummarizeProcesses {
		args = append(args, "--summarize-processes")
	}
	return args
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

// GetActionUUID generates action uuid.
// nolint: unused
func GetActionUUID() string {
	return "/action_id/" + uuid.New().String()
}
