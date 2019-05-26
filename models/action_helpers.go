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
	"context"
	"sync"

	"github.com/pkg/errors"
)

// InMemoryActionsStorage in memory action results storage implementation.
//nolint:unused
type InMemoryActionsStorage struct {
	container map[string]*ActionResult
	mx        sync.Mutex
}

// NewInMemoryActionsStorage created new InMemoryActionsStorage.
func NewInMemoryActionsStorage() *InMemoryActionsStorage {
	return &InMemoryActionsStorage{
		container: make(map[string]*ActionResult),
	}
}

// Store stores an action result in action results storage.
//nolint:unparam
func (s *InMemoryActionsStorage) Store(ctx context.Context, result *ActionResult) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.container[result.ID]
	if ok {
		return errors.New("already exists")
	}
	s.container[result.ID] = result
	return nil
}

// Update updates an action result in action results storage.
//nolint:unparam
func (s *InMemoryActionsStorage) Update(ctx context.Context, result *ActionResult) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.container[result.ID]
	if !ok {
		return errors.New("not found")
	}

	a := s.container[result.ID]

	a.PmmAgentID = result.PmmAgentID
	a.Error = result.Error
	a.Done = result.Done
	a.Output = result.Output
	return nil
}

// Load loads an action result from storage by action id.
//nolint:unparam
func (s *InMemoryActionsStorage) Load(ctx context.Context, id string) (*ActionResult, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	v, ok := s.container[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
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
