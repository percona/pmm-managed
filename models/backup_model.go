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

// BackupLocationType represents BackupLocation type as stored in database.
type BackupLocationType string

// BackupLocation types
const (
	S3BackupLocationType NodeType = "s3"
	FSBackupLocationType NodeType = "fs"
)

// BackupLocation represents destination for backup.
//reform:backup_locations
type BackupLocation struct {
	ID          string             `reform:"id,pk"`
	Name        string             `reform:"name"`
	Description string             `reform:"description"`
	Type        BackupLocationType `reform:"type"`
	Endpoint    *string            `reform:"endpoint"`
	AccessKey   *string            `reform:"access_key"`
	SecretKey   *string            `reform:"secret_key"`
	Path        *string            `reform:"path"`

	CreatedAt time.Time `reform:"created_at"`
	UpdatedAt time.Time `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (s *BackupLocation) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (s *BackupLocation) BeforeUpdate() error {
	s.UpdatedAt = Now()
	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (s *BackupLocation) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*BackupLocation)(nil)
	_ reform.BeforeUpdater  = (*BackupLocation)(nil)
	_ reform.AfterFinder    = (*BackupLocation)(nil)
)
