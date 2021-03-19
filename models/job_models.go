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
