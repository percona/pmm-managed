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
	"sync"
)

// InMemoryStorage in memory action results storage.
type InMemoryStorage struct {
	container map[string]Result
	mx        sync.Mutex
}

// NewInMemoryStorage created new InMemoryActionsStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

// Store stores an action result in action results storage.
func (s *InMemoryStorage) Store(ctx context.Context, result *Result) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.container[result.ID] = *result
}

// Load gets an action result from storage by action id.
func (s *InMemoryStorage) Load(ctx context.Context, id string) (*Result, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	v, ok := s.container[id]
	if !ok {
		return nil, false
	}
	return &v, true
}
