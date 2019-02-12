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

// ServiceType represents Service type as stored in database.
type ServiceType string

// Service types.
const (
	MySQLServiceType          ServiceType = "mysql"
	AmazonRDSMySQLServiceType ServiceType = "amazon-rds-mysql"
)

// ServiceRow represents Service as stored in database.
//reform:services
type ServiceRow struct {
	ServiceID   string      `reform:"service_id,pk"`
	ServiceType ServiceType `reform:"service_type"`
	ServiceName string      `reform:"service_name"`
	NodeID      string      `reform:"node_id"`
	CreatedAt   time.Time   `reform:"created_at"`
	// UpdatedAt time.Time   `reform:"updated_at"`

	Address    *string `reform:"address"`
	Port       *uint16 `reform:"port"`
	UnixSocket *string `reform:"unix_socket"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (sr *ServiceRow) BeforeInsert() error {
	now := time.Now().Truncate(time.Microsecond).UTC()
	sr.CreatedAt = now
	// sr.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (sr *ServiceRow) BeforeUpdate() error {
	// now := time.Now().Truncate(time.Microsecond).UTC()
	// sr.UpdatedAt = now
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (sr *ServiceRow) AfterFind() error {
	sr.CreatedAt = sr.CreatedAt.UTC()
	// sr.UpdatedAt = sr.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*ServiceRow)(nil)
	_ reform.BeforeUpdater  = (*ServiceRow)(nil)
	_ reform.AfterFinder    = (*ServiceRow)(nil)
)
