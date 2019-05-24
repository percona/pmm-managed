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
	"gopkg.in/reform.v1"
)

var (
	errPmmAgentIDNotFound = errors.New("can't detect pmm_agent_id")
)

type PMMAgentIDSQLResolver struct {
	db *reform.DB
}

func NewPMMAgentIDSQLResolver(db *reform.DB) *PMMAgentIDSQLResolver {
	return &PMMAgentIDSQLResolver{
		db: db,
	}
}

func (r *PMMAgentIDSQLResolver) ResolvePMMAgentIDByServiceID(serviceID, pmmAgentID string) (string, error) {
	agents, err := FindPMMAgentsForService(r.db.Querier, serviceID)
	if err != nil {
		return "", err
	}

	return validatePMMAgentID(pmmAgentID, agents)
}

func (r *PMMAgentIDSQLResolver) ResolvePMMAgentIDByNodeID(nodeID, pmmAgentID string) (string, error) {
	agents, err := FindPMMAgentsForNode(r.db.Querier, nodeID)
	if err != nil {
		return "", err
	}

	return validatePMMAgentID(pmmAgentID, agents)
}

func validatePMMAgentID(pmmAgentID string, agents []*Agent) (string, error) {
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
