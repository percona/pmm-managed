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

type PMMAgentIDResolver interface {
	ResolveByServiceID(serviceID, pmmAgentID string) (string, error)
	ResolveByNodeID(nodeID, pmmAgentID string) (string, error)
}

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
