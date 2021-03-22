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
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

//go:generate reform

// EchoJobResult stores echo job specific result data.
type EchoJobResult struct {
	Message string `json:"message"`
}

// JobType represents job type.
type JobType string

// Supported job types.
const (
	Echo = JobType("echo")
)

// JobResult describes a job result which is storing in persistent storage.
//reform:job_results
type JobResult struct {
	ID         string    `reform:"id,pk"`
	PMMAgentID string    `reform:"pmm_agent_id"`
	Type       JobType   `reform:"type"`
	Done       bool      `reform:"done"`
	Error      string    `reform:"error"`
	Result     []byte    `reform:"result"`
	CreatedAt  time.Time `reform:"created_at"`
	UpdatedAt  time.Time `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (r *JobResult) BeforeInsert() error {
	now := Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (r *JobResult) BeforeUpdate() error {
	r.UpdatedAt = Now()

	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (r *JobResult) AfterFind() error {
	r.CreatedAt = r.CreatedAt.UTC()
	r.UpdatedAt = r.UpdatedAt.UTC()

	return nil
}

// GetEchoJobResult extracts echo job result data.
func (r *JobResult) GetEchoJobResult() (*EchoJobResult, error) {
	var result EchoJobResult

	if err := json.Unmarshal(r.Result, &result); err != nil {
		return nil, errors.Wrap(err, "failed to parse echo job result")
	}

	return &result, nil
}

// SetEchoJobResult sets echo job result data.
func (r *JobResult) SetEchoJobResult(result *EchoJobResult) error {
	var err error
	if r.Result, err = json.Marshal(result); err != nil {
		return errors.Wrap(err, "failed to marshall echo job result")
	}

	return nil
}
