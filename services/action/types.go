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

package action

import "fmt"

type PtSummary struct {
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

func NewPtSummary(pmmAgentID, nodeID string) *PtSummary {
	return &PtSummary{
		ID:                 getUUID(),
		NodeID:             nodeID,
		PMMAgentID:         pmmAgentID,
		SummarizeMounts:    true,
		SummarizeNetwork:   true,
		SummarizeProcesses: true,
		Sleep:              5,
		Help:               false,
	}
}

func (s *PtSummary) Args() []string {
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

type PtMySQLSummary struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Args []string
}

type MySQLExplain struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

type MySQLExplainJSON struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// Result describes an PMM Action result which is storing in ActionsResult storage.
type Result struct {
	ID         string
	PmmAgentID string

	Done   bool
	Error  string
	Output string
}
