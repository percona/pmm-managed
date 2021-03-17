package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// FindJobResultByID finds JobResult by ID.
func FindJobResultByID(q *reform.Querier, id string) (*JobResult, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty JobResult ID.")
	}

	res := &JobResult{ID: id}
	switch err := q.Reload(res); err {
	case nil:
		return res, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "JobResult with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// CreateJobResult stores an action result in action results storage.
func CreateJobResult(q *reform.Querier, pmmAgentID string) (*JobResult, error) {
	result := &JobResult{ID: "/job_id/" + uuid.New().String(), PMMAgentID: pmmAgentID}
	if err := q.Insert(result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func ChangeJobResult(q *reform.Querier, result *JobResult)  error {
	if err := q.Update(result); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// CleanupOldJobResults deletes action results older than a specified date.
func CleanupOldJobResults(q *reform.Querier, olderThan time.Time) error {
	_, err := q.DeleteFrom(JobResultTable, " WHERE updated_at <= $1", olderThan)
	return err
}
