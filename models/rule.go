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
	"database/sql/driver"
	"time"

	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"gopkg.in/reform.v1"
)

//go:generate reform

// Rule represents alert rule configuration.
//reform:ia_rules
type Rule struct {
	Template     *iav1beta1.Template    `reform:"template"`
	ID           string                 `reform:"id,pk"`
	Summary      string                 `reform:"summary"`
	Disabled     bool                   `reform:"disabled"`
	Params       []*iav1beta1.RuleParam `reform:"params"`
	For          time.Duration          `reform:"for"`
	Severity     managementpb.Severity  `reform:"severity"`
	CustomLabels []byte                 `reform:"custom_labels"`
	Filters      []*Filter              `reform:"filters"`
	Channels     []*Channel             `reform:"channels"`
	CreatedAt    time.Time              `reform:"created_at"`
	UpdatedAt    time.Time              `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (r *Rule) BeforeInsert() error {
	now := Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (r *Rule) BeforeUpdate() error {
	r.UpdatedAt = Now()

	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (r *Rule) AfterFind() error {
	r.CreatedAt = r.CreatedAt.UTC()
	r.UpdatedAt = r.UpdatedAt.UTC()

	return nil
}

// TODO BeforeInsert, BeforeUpdate, AfterFind

type FilterType int32

const (
	Invalid FilterType = 0
	// =
	Equal FilterType = 1
	// !=
	NotEqual FilterType = 2
	// =~
	Regex FilterType = 3
	// !~
	NotRegex FilterType = 4
)

type Filter struct {
	Type FilterType `json:"type"`
	Key  string     `json:"key"`
	Val  string     `json:"value"`
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (f Filter) Value() (driver.Value, error) { return jsonValue(f) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (f *Filter) Scan(src interface{}) error { return jsonScan(f, src) }

// check interfaces.
var (
	_ reform.BeforeInserter = (*Rule)(nil)
	_ reform.BeforeUpdater  = (*Rule)(nil)
	_ reform.AfterFinder    = (*Rule)(nil)
)
