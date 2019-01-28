package models

import "gopkg.in/reform.v1"

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
