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

// JobType represents job type.
type JobType string

// Supported job types.
const (
	Echo                    = JobType("echo")
	MySQLBackupJob          = JobType("mysql_backup")
	MySQLRestoreBackupJob   = JobType("mysql_restore_backup")
	MongoDBBackupJob        = JobType("mongodb_backup")
	MongoDBRestoreBackupJob = JobType("mongodb_restore_backup")
)

// EchoJobResult stores echo job specific result data.
type EchoJobResult struct {
	Message string `json:"message"`
}

// MySQLBackupJobResult stores MySQL job specific result data.
type MySQLBackupJobResult struct {
}

// MySQLRestoreBackupJobResult stores MySQL restore backup job specific result data.
type MySQLRestoreBackupJobResult struct {
}

// MongoDBBackupJobResult stores MongoDB job specific result data.
type MongoDBBackupJobResult struct {
}

// MongoDBRestoreBackupJobResult stores MongoDB restore backup job specific result data.
type MongoDBRestoreBackupJobResult struct {
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (r JobResult) Value() (driver.Value, error) { return jsonValue(r) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (r *JobResult) Scan(src interface{}) error { return jsonScan(r, src) }

// JobResult holds result data for different job types.
type JobResult struct {
	Echo                 *EchoJobResult                 `json:"echo,omitempty"`
	MySQLBackup          *MySQLBackupJobResult          `json:"mysql_backup,omitempty"`
	MySQLRestoreBackup   *MySQLRestoreBackupJobResult   `json:"mysql_restore_backup,omitempty"`
	MongoDBBackup        *MongoDBBackupJobResult        `json:"mongo_db_backup,omitempty"`
	MongoDBRestoreBackup *MongoDBRestoreBackupJobResult `json:"mongo_db_restore_backup,omitempty"`
}

// EchoJobData stores echo job specific result data.
type EchoJobData struct {
	Message string        `json:"message"`
	Delay   time.Duration `json:"delay"`
}

// MySQLBackupJobData stores MySQL job specific result data.
type MySQLBackupJobData struct {
	ServiceID  string `json:"service_id"`
	ArtifactID string `json:"artifact_id"`
}

// MySQLRestoreBackupJobData stores MySQL restore backup job specific result data.
type MySQLRestoreBackupJobData struct {
	ServiceID string `json:"service_id"`
	RestoreID string `json:"restore_id,omitempty"`
}

// MongoDBBackupJobData stores MongoDB job specific result data.
type MongoDBBackupJobData struct {
	ServiceID  string `json:"service_id"`
	ArtifactID string `json:"artifact_id"`
}

// MongoDBRestoreBackupJobData stores MongoDB restore backup job specific result data.
type MongoDBRestoreBackupJobData struct {
	ServiceID string `json:"service_id"`
	RestoreID string `json:"restore_id,omitempty"`
}

type JobData struct {
	Echo                 *EchoJobData                 `json:"echo,omitempty"`
	MySQLBackup          *MySQLBackupJobData          `json:"mySQLBackup,omitempty"`
	MySQLRestoreBackup   *MySQLRestoreBackupJobData   `json:"mysql_restore_backup,omitempty"`
	MongoDBBackup        *MongoDBBackupJobData        `json:"mongoDBBackup,omitempty"`
	MongoDBRestoreBackup *MongoDBRestoreBackupJobData `json:"mongoDBRestoreBackup,omitempty"`
}

// Value implements database/sql/driver.Valuer interface. Should be defined on the value.
func (c JobData) Value() (driver.Value, error) { return jsonValue(c) }

// Scan implements database/sql.Scanner interface. Should be defined on the pointer.
func (c *JobData) Scan(src interface{}) error { return jsonScan(c, src) }

// Job describes a job result which is storing in persistent storage.
//reform:jobs
type Job struct {
	ID         string        `reform:"id,pk"`
	PMMAgentID string        `reform:"pmm_agent_id"`
	Type       JobType       `reform:"type"`
	Data       *JobData      `reform:"data"`
	Timeout    time.Duration `reform:"timeout"`
	Retries    uint32        `reform:"retries"`
	Interval   time.Duration `reform:"interval"`
	Done       bool          `reform:"done"`
	Error      string        `reform:"error"`
	Result     *JobResult    `reform:"result"`
	CreatedAt  time.Time     `reform:"created_at"`
	UpdatedAt  time.Time     `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (r *Job) BeforeInsert() error {
	now := Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (r *Job) BeforeUpdate() error {
	r.UpdatedAt = Now()

	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (r *Job) AfterFind() error {
	r.CreatedAt = r.CreatedAt.UTC()
	r.UpdatedAt = r.UpdatedAt.UTC()

	return nil
}

// check interfaces.
var (
	_ reform.BeforeInserter = (*Job)(nil)
	_ reform.BeforeUpdater  = (*Job)(nil)
	_ reform.AfterFinder    = (*Job)(nil)
)
