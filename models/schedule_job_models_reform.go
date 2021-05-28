// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package models

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type scheduleJobTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *scheduleJobTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("schedule_jobs").
func (v *scheduleJobTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *scheduleJobTableType) Columns() []string {
	return []string{
		"id",
		"cron_expression",
		"start_at",
		"last_run",
		"next_run",
		"type",
		"data",
		"retries",
		"created_at",
		"updated_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *scheduleJobTableType) NewStruct() reform.Struct {
	return new(ScheduleJob)
}

// NewRecord makes a new record for that table.
func (v *scheduleJobTableType) NewRecord() reform.Record {
	return new(ScheduleJob)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *scheduleJobTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ScheduleJobTable represents schedule_jobs view or table in SQL database.
var ScheduleJobTable = &scheduleJobTableType{
	s: parse.StructInfo{
		Type:    "ScheduleJob",
		SQLName: "schedule_jobs",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "CronExpression", Type: "string", Column: "cron_expression"},
			{Name: "StartAt", Type: "time.Time", Column: "start_at"},
			{Name: "LastRun", Type: "time.Time", Column: "last_run"},
			{Name: "NextRun", Type: "time.Time", Column: "next_run"},
			{Name: "Type", Type: "ScheduleJobType", Column: "type"},
			{Name: "Data", Type: "*ScheduleJobData", Column: "data"},
			{Name: "Retries", Type: "uint", Column: "retries"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(ScheduleJob).Values(),
}

// String returns a string representation of this struct or record.
func (s ScheduleJob) String() string {
	res := make([]string, 10)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "CronExpression: " + reform.Inspect(s.CronExpression, true)
	res[2] = "StartAt: " + reform.Inspect(s.StartAt, true)
	res[3] = "LastRun: " + reform.Inspect(s.LastRun, true)
	res[4] = "NextRun: " + reform.Inspect(s.NextRun, true)
	res[5] = "Type: " + reform.Inspect(s.Type, true)
	res[6] = "Data: " + reform.Inspect(s.Data, true)
	res[7] = "Retries: " + reform.Inspect(s.Retries, true)
	res[8] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[9] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *ScheduleJob) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.CronExpression,
		s.StartAt,
		s.LastRun,
		s.NextRun,
		s.Type,
		s.Data,
		s.Retries,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *ScheduleJob) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.CronExpression,
		&s.StartAt,
		&s.LastRun,
		&s.NextRun,
		&s.Type,
		&s.Data,
		&s.Retries,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *ScheduleJob) View() reform.View {
	return ScheduleJobTable
}

// Table returns Table object for that record.
func (s *ScheduleJob) Table() reform.Table {
	return ScheduleJobTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *ScheduleJob) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *ScheduleJob) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *ScheduleJob) HasPK() bool {
	return s.ID != ScheduleJobTable.z[ScheduleJobTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (s *ScheduleJob) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = ScheduleJobTable
	_ reform.Struct = (*ScheduleJob)(nil)
	_ reform.Table  = ScheduleJobTable
	_ reform.Record = (*ScheduleJob)(nil)
	_ fmt.Stringer  = (*ScheduleJob)(nil)
)

func init() {
	parse.AssertUpToDate(&ScheduleJobTable.s, new(ScheduleJob))
}
