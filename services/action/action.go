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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/models"
)

// PtSummary represents pt-summary domain model.
type PtSummary struct {
	ID         string
	PMMAgentID string
	NodeID     string

	Args []string
}

// PtMySQLSummary represents pt-mysql-summary domain model.
type PtMySQLSummary struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Args []string
}

// MySQLExplain represents mysql-explain domain model.
type MySQLExplain struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// MySQLExplainJSON represents mysql-explain-json domain model.
type MySQLExplainJSON struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// FindPmmAgentIDToRunAction finds pmm-agent-id to run action.
func FindPmmAgentIDToRunAction(pmmAgentID string, agents []*models.Agent) (string, error) {
	// no explicit ID is given, and there is only one
	if pmmAgentID == "" && len(agents) == 1 {
		return agents[0].AgentID, nil
	}

	// no explicit ID is given, and there are zero or several
	if pmmAgentID == "" {
		return "", status.Errorf(codes.InvalidArgument, "Couldn't find pmm-agent-id to run action")
	}

	// check that explicit agent id is correct
	for _, a := range agents {
		if a.AgentID == pmmAgentID {
			return a.AgentID, nil
		}
	}
	return "", status.Errorf(codes.FailedPrecondition, "Couldn't find pmm-agent-id to run action")
}
