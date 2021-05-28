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

// ScheduleJobType represents schedule job type.
type ScheduleJobType string

// Supported schedule job types.
const (
	ScheduleEchoJob = ScheduleJobType("echo")
)

// ScheduleJob describes a scheduled job.
//reform:schedule_jobs
type ScheduleJob struct {
	ID             string           `reform:"id,pk"`
	CronExpression string           `reform:"cron_expression"`
	StartAt        time.Time        `reform:"start_at"`
	LastRun        time.Time        `reform:"last_run"`
	NextRun        time.Time        `reform:"next_run"`
	Type           ScheduleJobType  `reform:"type"`
	Data           *ScheduleJobData `reform:"data"`
	Retries        uint             `reform:"retries"`
	CreatedAt      time.Time        `reform:"created_at"`
	UpdatedAt      time.Time        `reform:"updated_at"`
}

// ScheduleJobData holds result data for different job types.
type ScheduleJobData struct {
	Echo *EchoJobData `json:"echo,omitempty"`
}

type EchoJobData struct {
	Value string
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (c ScheduleJobData) Value() (driver.Value, error) { return jsonValue(c) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (c *ScheduleJobData) Scan(src interface{}) error { return jsonScan(c, src) }

// BeforeInsert implements reform.BeforeInserter interface.
func (r *ScheduleJob) BeforeInsert() error {
	now := Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (r *ScheduleJob) BeforeUpdate() error {
	r.UpdatedAt = Now()

	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (r *ScheduleJob) AfterFind() error {
	r.CreatedAt = r.CreatedAt.UTC()
	r.UpdatedAt = r.UpdatedAt.UTC()

	return nil
}

// check interfaces.
var (
	_ reform.BeforeInserter = (*ScheduleJob)(nil)
	_ reform.BeforeUpdater  = (*ScheduleJob)(nil)
	_ reform.AfterFinder    = (*ScheduleJob)(nil)
)
