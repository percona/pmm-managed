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

// JobResult describes a job result which is storing in persistent storage.
//reform:job_results
type JobResult struct {
	ID         string    `reform:"id,pk"`
	PMMAgentID string    `reform:"pmm_agent_id"`
	Done       bool      `reform:"done"`
	Error      string    `reform:"error"`
	Result     []byte    `reform:"result"`
	CreatedAt  time.Time `reform:"created_at"`
	UpdatedAt  time.Time `reform:"updated_at"`
}

// GetEchoJobResult extracts echo job result data.
func (j *JobResult) GetEchoJobResult() (*EchoJobResult, error) {
	var result EchoJobResult

	if err := json.Unmarshal(j.Result, &result); err != nil {
		return nil, errors.Wrap(err, "failed to parse echo job result")
	}

	return &result, nil
}

// SetEchoJobResult sets echo job result data.
func (j *JobResult) SetEchoJobResult(result *EchoJobResult) error {
	var err error
	if j.Result, err = json.Marshal(result); err != nil {
		return errors.Wrap(err, "failed to marshall echo job result")
	}

	return nil
}
