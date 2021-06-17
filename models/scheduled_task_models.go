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

// ScheduledTaskType represents scheduled task type.
type ScheduledTaskType string

// Supported scheduled task types.
const (
	ScheduledPrintTask       = ScheduledTaskType("print")
	ScheduledMySQLBackupTask = ScheduledTaskType("mysql_backup")
	ScheduledMongoBackupTask = ScheduledTaskType("mongo_backup")
)

// ScheduledTask describes a scheduled task.
//reform:scheduled_tasks
type ScheduledTask struct {
	ID               string             `reform:"id,pk"`
	CronExpression   string             `reform:"cron_expression"`
	Disabled         bool               `reform:"disabled"`
	StartAt          time.Time          `reform:"start_at"`
	LastRun          time.Time          `reform:"last_run"`
	NextRun          time.Time          `reform:"next_run"`
	Type             ScheduledTaskType  `reform:"type"`
	Data             *ScheduledTaskData `reform:"data"`
	Retries          uint               `reform:"retries"`
	RetryInterval    time.Duration      `reform:"retry_interval"`
	RetriesRemaining uint               `reform:"retries_remaining"`
	Succeeded        uint               `reform:"succeeded"`
	Failed           uint               `reform:"failed"`
	Running          bool               `reform:"running"`
	CreatedAt        time.Time          `reform:"created_at"`
	UpdatedAt        time.Time          `reform:"updated_at"`
}

// ScheduledTaskData holds result data for different task types.
type ScheduledTaskData struct {
	Print           *PrintTaskData       `json:"print,omitempty"`
	MySQLBackupTask *MySQLBackupTaskData `json:"mySQLBackup,omitempty"`
	MongoBackupTask *MongoBackupTaskData `json:"mongoBackup,omitempty"`
}

// PrintTaskData holds data for print task.
type PrintTaskData struct {
	Message string
}

// BackupRetryData holds common data for backup retrying.
type BackupRetryData struct {
	Retries  uint          `json:"retries"`
	Interval time.Duration `json:"interval"`
}

// MySQLBackupTaskData holds data for mysql backup task.
type MySQLBackupTaskData struct {
	Retry       BackupRetryData `json:"retry"`
	ServiceID   string          `json:"serviceID"`
	LocationID  string          `json:"locationID"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
}

// MongoBackupTaskData holds data for mysql backup task.
type MongoBackupTaskData struct {
	Retry       BackupRetryData `json:"retry"`
	ServiceID   string          `json:"serviceID"`
	LocationID  string          `json:"locationID"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (c ScheduledTaskData) Value() (driver.Value, error) { return jsonValue(c) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (c *ScheduledTaskData) Scan(src interface{}) error { return jsonScan(c, src) }

// BeforeInsert implements reform.BeforeInserter interface.
func (r *ScheduledTask) BeforeInsert() error {
	now := Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (r *ScheduledTask) BeforeUpdate() error {
	r.UpdatedAt = Now()

	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (r *ScheduledTask) AfterFind() error {
	r.CreatedAt = r.CreatedAt.UTC()
	r.UpdatedAt = r.UpdatedAt.UTC()
	r.StartAt = r.StartAt.UTC()
	r.NextRun = r.NextRun.UTC()
	r.LastRun = r.LastRun.UTC()

	return nil
}

// check interfaces.
var (
	_ reform.BeforeInserter = (*ScheduledTask)(nil)
	_ reform.BeforeUpdater  = (*ScheduledTask)(nil)
	_ reform.AfterFinder    = (*ScheduledTask)(nil)
)
