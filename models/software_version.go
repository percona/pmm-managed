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

	"database/sql/driver"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

//go:generate reform

// SoftwareName represents software name.
type SoftwareName string

// SoftwareName types of different software.
const (
	MysqldSoftwareName     SoftwareName = "mysqld"
	XtrabackupSoftwareName SoftwareName = "xtrabackup"
	XbcloudSoftwareName    SoftwareName = "xbcloud"
	QpressSoftwareName     SoftwareName = "qpress"
)

// SoftwareVersion represents version of the given software.
type SoftwareVersion struct {
	Name    SoftwareName `reform:"name"`
	Version string       `reform:"version"`
}

// SoftwareVersions represents slice of SoftwareVersion.
type SoftwareVersions []SoftwareVersion

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (sv SoftwareVersions) Value() (driver.Value, error) {
	return sv, nil
}

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (sv *SoftwareVersions) Scan(src interface{}) error {
	switch s := src.(type) {
	case SoftwareVersions:
		*sv = s
	default:
		return errors.Errorf("unexpected type for scanning software versions: %T", s)
	}

	return nil
}

// ServiceSoftwareVersions represents service software versions.
//reform:service_software_versions
type ServiceSoftwareVersions struct {
	ServiceID        string           `reform:"service_id,pk"`
	SoftwareVersions SoftwareVersions `reform:"software_versions"`
	NextCheckAt      time.Time        `reform:"next_check_at"`
	CreatedAt        time.Time        `reform:"created_at"`
	UpdatedAt        time.Time        `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (s *ServiceSoftwareVersions) BeforeInsert() error {
	s.CreatedAt = Now()
	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (s *ServiceSoftwareVersions) AfterFind() error {
	s.NextCheckAt = s.NextCheckAt.UTC()
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (s *ServiceSoftwareVersions) BeforeUpdate() error {
	s.UpdatedAt = Now()
	return nil
}

// check interfaces.
var (
	_ reform.BeforeInserter = (*ServiceSoftwareVersions)(nil)
	_ reform.AfterFinder    = (*ServiceSoftwareVersions)(nil)
	_ reform.BeforeUpdater  = (*ServiceSoftwareVersions)(nil)
)
