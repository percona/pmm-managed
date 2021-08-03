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
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// FindJobByID finds Job by ID.
func FindJobByID(q *reform.Querier, id string) (*Job, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Job ID.")
	}

	res := &Job{ID: id}
	switch err := q.Reload(res); err {
	case nil:
		return res, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Job with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// CreateJob stores a job result in the storage.
func CreateJob(q *reform.Querier, pmmAgentID string, jobType JobType, data *JobData) (*Job, error) {
	result := &Job{
		ID:         "/job_id/" + uuid.New().String(),
		PMMAgentID: pmmAgentID,
		Type:       jobType,
		Data:       data,
	}
	if err := q.Insert(result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

// CleanupOldJobs deletes jobs results older than a specified date.
func CleanupOldJobs(q *reform.Querier, olderThan time.Time) error {
	_, err := q.DeleteFrom(JobTable, " WHERE updated_at <= $1", olderThan)
	return err
}
