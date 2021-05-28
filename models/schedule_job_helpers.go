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
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// FindScheduleJobByID finds ScheduleJob by ID.
func FindScheduleJobByID(q *reform.Querier, id string) (*ScheduleJob, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty JobResult ID.")
	}

	res := &ScheduleJob{ID: id}
	switch err := q.Reload(res); err {
	case nil:
		return res, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "ScheduleJob with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

type ScheduleJobsFilter struct {
	Disabled *bool
}

// FindScheduleJobs returns all scheduled satisfying filter.
func FindScheduleJobs(q *reform.Querier, filters ScheduleJobsFilter) ([]*ScheduleJob, error) {
	tail := ""
	if filters.Disabled != nil {
		tail = "WHERE disabled IS"
		if *filters.Disabled {
			tail += "TRUE"
		} else {
			tail += "FALSE"
		}
	}
	structs, err := q.SelectAllFrom(ScheduleJobTable, "")
	if err != nil {
		return nil, err
	}
	jobs := make([]*ScheduleJob, len(structs))
	for i, s := range structs {
		jobs[i] = s.(*ScheduleJob)
	}
	return jobs, nil
}

// CreateScheduleJobParams are params for creating new schedule job.
type CreateScheduleJobParams struct {
	CronExpression string
	StartAt        time.Time
	NextRun        time.Time
	Type           ScheduleJobType
	Data           ScheduleJobData
	Retries        uint
	RetryInterval  time.Duration
}

// Validate checks if required params are set and valid.
func (p CreateScheduleJobParams) Validate() error {
	switch p.Type {
	case ScheduleEchoJob:
	default:
		return status.Errorf(codes.InvalidArgument, "Unknown type: %s", p.Type)
	}
	_, err := cron.ParseStandard(p.CronExpression)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid cron expression: %v", err)
	}

	return nil
}

// CreateScheduleJob creates schedule job.
func CreateScheduleJob(q *reform.Querier, params CreateScheduleJobParams) (*ScheduleJob, error) {
	id := "/schedule_job_id/" + uuid.New().String()
	if err := checkUniqueScheduleJobID(q, id); err != nil {
		return nil, err
	}

	job := &ScheduleJob{
		ID:             id,
		Disabled:       false,
		CronExpression: params.CronExpression,
		StartAt:        params.StartAt,
		NextRun:        params.NextRun,
		Type:           params.Type,
		Data:           &params.Data,
		Retries:        params.Retries,
		RetryInterval:  params.RetryInterval,
	}
	if err := q.Insert(job); err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

// ChangeScheduleJobParams are params for updating existing schedule job.
type ChangeScheduleJobParams struct {
	NextRun time.Time
	LastRun time.Time
	Disable *bool
}

// ChangeScheduleJob updates existing schedule job.
func ChangeScheduleJob(q *reform.Querier, scheduleJobID string, params ChangeScheduleJobParams) (*ScheduleJob, error) {
	row, err := FindScheduleJobByID(q, scheduleJobID)
	if err != nil {
		return nil, err
	}
	row.NextRun = params.NextRun
	row.LastRun = params.LastRun

	if params.Disable != nil {
		row.Disabled = *params.Disable
	}

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update schedule job")
	}

	return row, nil
}

func RemoveScheduleJob(q *reform.Querier, id string) error {
	if _, err := FindScheduleJobByID(q, id); err != nil {
		return err
	}
	if err := q.Delete(&ScheduleJob{ID: id}); err != nil {
		return errors.Wrap(err, "failed to delete ScheduleJob")
	}

	return nil
}

func checkUniqueScheduleJobID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty schedule job ID")
	}

	location := &BackupLocation{ID: id}
	switch err := q.Reload(location); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Location with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}
