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
	"time"

	"gopkg.in/reform.v1"
)

//go:generate reform

// Backup represents destination for backup.
//reform:backups
type Backup struct {
	ID           string `reform:"id,pk"`
	Name         string `reform:"name"`
	LocationName string `reform:"location_name"`

	CreatedAt time.Time `reform:"created_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (s *Backup) BeforeInsert() error {
	s.CreatedAt = Now()
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (s *Backup) BeforeUpdate() error {
	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (s *Backup) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Backup)(nil)
	_ reform.BeforeUpdater  = (*Backup)(nil)
	_ reform.AfterFinder    = (*Backup)(nil)
)
