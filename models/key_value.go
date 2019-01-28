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
	"gopkg.in/reform.v1"
)

//go:generate reform

// KeyValue represents Service as stored in database.
//reform:key_values
type KeyValue struct {
	Key   string `reform:"key,pk"`
	Value string `reform:"value"`
	// CreatedAt time.Time   `reform:"created_at"`
	// UpdatedAt time.Time   `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (sr *KeyValue) BeforeInsert() error {
	// now := time.Now().Truncate(time.Microsecond).UTC()
	// sr.CreatedAt = now
	// sr.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (sr *KeyValue) BeforeUpdate() error {
	// now := time.Now().Truncate(time.Microsecond).UTC()
	// sr.UpdatedAt = now
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (sr *KeyValue) AfterFind() error {
	// sr.CreatedAt = sr.CreatedAt.UTC()
	// sr.UpdatedAt = sr.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*KeyValue)(nil)
	_ reform.BeforeUpdater  = (*KeyValue)(nil)
	_ reform.AfterFinder    = (*KeyValue)(nil)
)
