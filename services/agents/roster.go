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

import "sync"

const (
	rdsGroup = "rds"
)

type roster struct {
	rw sync.RWMutex
	m  map[string]map[string][]string
}

func newRoster() *roster {
	return &roster{
		m: make(map[string]map[string][]string),
	}
}

func (r *roster) add(id string, group string, ids []string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.m[id] == nil {
		r.m[id] = make(map[string][]string)
	}
	r.m[id][group] = ids
}

func (r *roster) get(id string, group string) []string {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.m[id] == nil {
		return nil
	}
	return r.m[id][group]
}

func (r *roster) remove(id string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	delete(r.m, id)
}
