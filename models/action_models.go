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

	"github.com/google/uuid"
)

// PtSummary represents pt-summary domain model.
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

// PtMySQLSummary represents pt-mysql-summary domain model.
type PtMySQLSummaryAction struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Args []string
}

// MySQLExplain represents mysql-explain domain model.
type MySQLExplainAction struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// MySQLExplainJSON represents mysql-explain-json domain model.
type MySQLExplainJSONAction struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// Result describes an action result which is storing in persistent storage.
type ActionResult struct {
	ID         string
	PmmAgentID string

	Done   bool
	Error  string
	Output string
}

// nolint: unused
func GetActionUUID() string {
	return "/action_id/" + uuid.New().String()
}
