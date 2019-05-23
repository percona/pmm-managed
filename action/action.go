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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/models"
)

type StarterStopper interface {
	ptSummaryStarter
	ptMySQLSummaryStarter
	mysqlExplainStarter
	mysqlExplainJSONStarter
	Stopper
}

type Stopper interface {
	StopAction(ctx context.Context, actionID string) error
}

type ptSummaryStarter interface {
	StartPTSummaryAction(context.Context, *PtSummary) error
}

type ptMySQLSummaryStarter interface {
	StartPTMySQLSummaryAction(context.Context, *PtMySQLSummary) error
}

type mysqlExplainStarter interface {
	StartMySQLExplainAction(context.Context, *MySQLExplain) error
}

type mysqlExplainJSONStarter interface {
	StartMySQLExplainJSONAction(context.Context, *MySQLExplainJSON) error
}

var (
	errPmmAgentIDNotFound = status.Error(codes.Internal, "can't detect pmm_agent_id")
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
