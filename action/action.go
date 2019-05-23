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

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/percona/pmm-managed/models"
)

type Stopper interface {
	StopAction(ctx context.Context, actionID, pmmAgentID string) error
}

type PtSummaryStarter interface {
	StartPTSummaryAction(context.Context, *PtSummary) error
}

type PtMySQLSummaryStarter interface {
	StartPTMySQLSummaryAction(context.Context, *PtMySQLSummary) error
}

type MySQLExplainStarter interface {
	StartMySQLExplainAction(context.Context, *MySQLExplain) error
}

type MySQLExplainJSONStarter interface {
	StartMySQLExplainJSONAction(context.Context, *MySQLExplainJSON) error
}

// Storage is an interface represents abstract storage for action results.
type Storage interface {
	// Store an action result to persistent storage.
	Store(context.Context, *Result)
	Load(context.Context, string) (*Result, bool)
}

var (
	errPmmAgentIDNotFound = errors.New("can't detect pmm_agent_id")
)

// Result describes an PMM Action result which is storing in ActionsResult storage.
//nolint:unused
type Result struct {
	ID         string
	PmmAgentID string
	Done       bool
	Error      string
	Output     string
}

func createActionID() string {
	return "/action_id/" + uuid.New().String()
}

func validatePMMAgentID(pmmAgentID string, agents []*models.Agent) (string, error) {
	// no explicit ID is given, and there is only one
	if pmmAgentID == "" && len(agents) == 1 {
		return agents[0].AgentID, nil
	}

	// no explicit ID is given, and there are zero or several
	if pmmAgentID == "" {
		return "", errPmmAgentIDNotFound
	}

	// check that explicit agent id is correct
	for _, a := range agents {
		if a.AgentID == pmmAgentID {
			return a.AgentID, nil
		}
	}
	return "", errPmmAgentIDNotFound
}
