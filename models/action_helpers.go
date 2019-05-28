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
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// InsertActionResult stores an action result in action results storage.
func InsertActionResult(q *reform.Querier, result *ActionResult) error {
	if err := q.Insert(result); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// UpdateActionResult updates an action result in action results storage.
func UpdateActionResult(q *reform.Querier, result *ActionResult) error {
	if err := q.Update(result); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LoadActionResult loads an action result from storage by action id.
func LoadActionResult(q *reform.Querier, id string) (*ActionResult, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty ActionResult ID.")
	}

	row := &ActionResult{ID: id}
	switch err := q.Reload(row); err {
	case nil:
		return row, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "ActionResult with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// FindPmmAgentIDToRunAction finds pmm-agent-id to run action.
func FindPmmAgentIDToRunAction(pmmAgentID string, agents []*Agent) (string, error) {
	// no explicit ID is given, and there is only one
	if pmmAgentID == "" && len(agents) == 1 {
		return agents[0].AgentID, nil
	}

	// no explicit ID is given, and there are zero or several
	if pmmAgentID == "" {
		return "", errors.New("couldn't find pmm-agent-id to run action")
	}

	// check that explicit agent id is correct
	for _, a := range agents {
		if a.AgentID == pmmAgentID {
			return a.AgentID, nil
		}
	}
	return "", errors.New("couldn't find pmm-agent-id to run action")
}
