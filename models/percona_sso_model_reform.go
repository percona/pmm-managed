// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package models

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type perconaSSODetailsTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *perconaSSODetailsTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("percona_sso_details").
func (v *perconaSSODetailsTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *perconaSSODetailsTableType) Columns() []string {
	return []string{
		"client_id",
		"client_secret",
		"issuer_url",
		"scope",
		"access_token",
		"organization_id",
		"created_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *perconaSSODetailsTableType) NewStruct() reform.Struct {
	return new(PerconaSSODetails)
}

// NewRecord makes a new record for that table.
func (v *perconaSSODetailsTableType) NewRecord() reform.Record {
	return new(PerconaSSODetails)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *perconaSSODetailsTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// PerconaSSODetailsTable represents percona_sso_details view or table in SQL database.
var PerconaSSODetailsTable = &perconaSSODetailsTableType{
	s: parse.StructInfo{
		Type:    "PerconaSSODetails",
		SQLName: "percona_sso_details",
		Fields: []parse.FieldInfo{
			{Name: "ClientID", Type: "string", Column: "client_id"},
			{Name: "ClientSecret", Type: "string", Column: "client_secret"},
			{Name: "IssuerURL", Type: "string", Column: "issuer_url"},
			{Name: "Scope", Type: "string", Column: "scope"},
			{Name: "AccessToken", Type: "*PerconaSSOAccessToken", Column: "access_token"},
			{Name: "OrganizationID", Type: "string", Column: "organization_id"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(PerconaSSODetails).Values(),
}

// String returns a string representation of this struct or record.
func (s PerconaSSODetails) String() string {
	res := make([]string, 7)
	res[0] = "ClientID: " + reform.Inspect(s.ClientID, true)
	res[1] = "ClientSecret: " + reform.Inspect(s.ClientSecret, true)
	res[2] = "IssuerURL: " + reform.Inspect(s.IssuerURL, true)
	res[3] = "Scope: " + reform.Inspect(s.Scope, true)
	res[4] = "AccessToken: " + reform.Inspect(s.AccessToken, true)
	res[5] = "OrganizationID: " + reform.Inspect(s.OrganizationID, true)
	res[6] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *PerconaSSODetails) Values() []interface{} {
	return []interface{}{
		s.ClientID,
		s.ClientSecret,
		s.IssuerURL,
		s.Scope,
		s.AccessToken,
		s.OrganizationID,
		s.CreatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *PerconaSSODetails) Pointers() []interface{} {
	return []interface{}{
		&s.ClientID,
		&s.ClientSecret,
		&s.IssuerURL,
		&s.Scope,
		&s.AccessToken,
		&s.OrganizationID,
		&s.CreatedAt,
	}
}

// View returns View object for that struct.
func (s *PerconaSSODetails) View() reform.View {
	return PerconaSSODetailsTable
}

// Table returns Table object for that record.
func (s *PerconaSSODetails) Table() reform.Table {
	return PerconaSSODetailsTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *PerconaSSODetails) PKValue() interface{} {
	return s.ClientID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *PerconaSSODetails) PKPointer() interface{} {
	return &s.ClientID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *PerconaSSODetails) HasPK() bool {
	return s.ClientID != PerconaSSODetailsTable.z[PerconaSSODetailsTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ClientID = pk.
func (s *PerconaSSODetails) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = PerconaSSODetailsTable
	_ reform.Struct = (*PerconaSSODetails)(nil)
	_ reform.Table  = PerconaSSODetailsTable
	_ reform.Record = (*PerconaSSODetails)(nil)
	_ fmt.Stringer  = (*PerconaSSODetails)(nil)
)

func init() {
	parse.AssertUpToDate(&PerconaSSODetailsTable.s, new(PerconaSSODetails))
}
