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

package agents

import (
	"sync"
)

type agentGroup string

const (
	rdsGroup agentGroup = "rds"
)

// roster groups several Agent IDs from an Inventory model to a single ID, as seen by pmm-agent.
//
// Currently, it is used only for rds_exporter.
// TODO Revisit it once we need it for something else.
type roster struct {
	rw sync.RWMutex
	m  map[string]map[agentGroup][]string
}

func newRoster() *roster {
	return &roster{
		m: make(map[string]map[agentGroup][]string),
	}
}

func (r *roster) add(pmmAgentID string, group agentGroup, agentIDs []string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.m[pmmAgentID] == nil {
		r.m[pmmAgentID] = make(map[agentGroup][]string)
	}
	r.m[pmmAgentID][group] = agentIDs
}

func (r *roster) get(pmmAgentID string, group agentGroup) []string {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.m[pmmAgentID] == nil {
		return nil
	}
	return r.m[pmmAgentID][group]
}

func (r *roster) remove(pmmAgentID string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	delete(r.m, pmmAgentID)
}
