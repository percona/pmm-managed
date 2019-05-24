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

// Package inmemory is an infrastructure level package that provides in memory storage implementations.
package inmemory

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/percona/pmm-managed/services/action"
)

// ActionStorage in memory action results storage.
type ActionStorage struct {
	container map[string]*action.Result
	mx        sync.Mutex
}

// NewActionStorage created new InMemoryActionsStorage.
func NewActionStorage() *ActionStorage {
	return &ActionStorage{}
}

// Store stores an action result in action results storage.
//nolint:unparam
func (s *ActionStorage) Store(ctx context.Context, result *action.Result) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.container[result.ID]
	if ok {
		return errors.New("ActionResult already exists")
	}
	s.container[result.ID] = result
	return nil
}

// Update updates an action result in action results storage.
//nolint:unparam
func (s *ActionStorage) Update(ctx context.Context, result *action.Result) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, ok := s.container[result.ID]
	if !ok {
		return errors.New("ActionResult doesn't exists")
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
func (s *ActionStorage) Load(ctx context.Context, id string) (*action.Result, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	v, ok := s.container[id]
	if !ok {
		return nil, errors.New("ActionResult not found")
	}
	return v, nil
}
