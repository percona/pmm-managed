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

	"gopkg.in/reform.v1"
)

//go:generate reform

// Rule represents alert rule configuration.
//reform:ia_rules
type Rule struct {
	Template     *Template     `reform:"template"`
	ID           string        `reform:"id,pk"`
	Summary      string        `reform:"summary"`
	Disabled     bool          `reform:"disabled"`
	Params       RuleParams    `reform:"params"`
	For          time.Duration `reform:"for"`
	Severity     string        `reform:"severity"`
	CustomLabels []byte        `reform:"custom_labels"`
	Filters      Filters       `reform:"filters"`
	Channels     Channels      `reform:"channels"`
	CreatedAt    time.Time     `reform:"created_at"`
	UpdatedAt    time.Time     `reform:"updated_at"`
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

type Filters []Filter

// Value implements database/sql/driver Valuer interface.
func (t Filters) Value() (driver.Value, error) { return jsonValue(t) }

// Scan implements database/sql Scanner interface.
func (t *Filters) Scan(src interface{}) error { return jsonScan(t, src) }

type Filter struct {
	Type FilterType `json:"type"`
	Key  string     `json:"key"`
	Val  string     `json:"value"`
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (f Filter) Value() (driver.Value, error) { return jsonValue(f) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (f *Filter) Scan(src interface{}) error { return jsonScan(f, src) }

type ParamType int32

const (
	InvalidRuleParam ParamType = 0
	BoolRuleParam    ParamType = 1
	FloatRuleParam   ParamType = 2
	StringRuleParam  ParamType = 3
)

type RuleParams []RuleParam

// Value implements database/sql/driver Valuer interface.
func (t RuleParams) Value() (driver.Value, error) { return jsonValue(t) }

// Scan implements database/sql Scanner interface.
func (t *RuleParams) Scan(src interface{}) error { return jsonScan(t, src) }

type RuleParam struct {
	Name      string    `json:"name"`
	Type      ParamType `json:"type"`
	BoolVal   bool      `json:"bval"`
	FloatVal  float32   `json:"fval"`
	StringVal string    `json:"sval"`
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (p RuleParam) Value() (driver.Value, error) { return jsonValue(p) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (p *RuleParam) Scan(src interface{}) error { return jsonScan(p, src) }

// Channels is a slice of models.Channel
type Channels []*Channel

// Value implements database/sql/driver Valuer interface.
func (t Channels) Value() (driver.Value, error) { return jsonValue(t) }

// Scan implements database/sql Scanner interface.
func (t *Channels) Scan(src interface{}) error { return jsonScan(t, src) }

// check interfaces.
var (
	_ reform.BeforeInserter = (*Rule)(nil)
	_ reform.BeforeUpdater  = (*Rule)(nil)
	_ reform.AfterFinder    = (*Rule)(nil)
)
