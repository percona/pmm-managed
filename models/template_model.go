package models

import (
	"database/sql/driver"
	"time"

	"github.com/percona-platform/saas/pkg/common"
	"gopkg.in/reform.v1"
)

//go:generate reform

//reform:notification_rule_templates
type Template struct {
	Name        string   `reform:"name,pk"`
	Version     uint32   `reform:"version"`
	Summary     string   `reform:"summary"`
	Tiers       Tiers    `reform:"tiers"`
	Expr        string   `reform:"expr"`
	Params      Params   `reform:"params"`
	For         Duration `reform:"for"`
	Severity    string   `reform:"severity"`
	Labels      Map      `reform:"labels"`
	Annotations Map      `reform:"annotations"`
	Source      string   `reform:"source"`

	CreatedAt time.Time `reform:"created_at"`
	UpdatedAt time.Time `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (t *Template) BeforeInsert() error {
	now := Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (t *Template) BeforeUpdate() error {
	t.UpdatedAt = Now()

	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (t *Template) AfterFind() error {
	t.CreatedAt = t.CreatedAt.UTC()
	t.UpdatedAt = t.UpdatedAt.UTC()

	return nil
}

type Tiers []common.Tier

// Value implements database/sql/driver Valuer interface.
func (t Tiers) Value() (driver.Value, error) { return jsonValue(t) }

// Scan implements database/sql Scanner interface.
func (t *Tiers) Scan(src interface{}) error { return jsonScan(t, src) }

type Map map[string]string

// Value implements database/sql/driver Valuer interface.
func (m Map) Value() (driver.Value, error) { return jsonValue(m) }

// Scan implements database/sql Scanner interface.
func (m *Map) Scan(src interface{}) error { return jsonScan(m, src) }

type Duration time.Duration

type Params []Param

// Value implements database/sql/driver Valuer interface.
func (p Params) Value() (driver.Value, error) { return jsonValue(p) }

// Scan implements database/sql Scanner interface.
func (p *Params) Scan(src interface{}) error { return jsonScan(p, src) }

type Param struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Unit    string `json:"unit"`
	Type    string `json:"type"`

	FloatParam *FloatParam `json:"float_param"`
	// BoolParam   *BoolParam   `json:"bool_param"`
	// StringParam *StringParam `json:"string_param"`
}

type BoolParam struct {
	Default bool `json:"default"`
}

type FloatParam struct {
	HasDefault bool    `json:"has_default"`
	Default    float64 `json:"default"`

	HasMin bool    `json:"has_min"`
	Min    float64 `json:"min"`

	HaxMax bool    `json:"hax_max"`
	Max    float64 `json:"max"`
}

type StringParam struct {
	HasDefault bool   `json:"has_default"`
	Default    string `json:"default"`
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Template)(nil)
	_ reform.BeforeUpdater  = (*Template)(nil)
	_ reform.AfterFinder    = (*Template)(nil)
)
