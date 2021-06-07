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
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// FindScheduledTaskByID finds ScheduledTask by ID.
func FindScheduledTaskByID(q *reform.Querier, id string) (*ScheduledTask, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty ScheduledTask ID.")
	}

	res := &ScheduledTask{ID: id}
	switch err := q.Reload(res); err {
	case nil:
		return res, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "ScheduledTask with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

type ScheduledTasksFilter struct {
	Disabled *bool
	Types    []ScheduledTaskType
}

// FindScheduledTasks returns all scheduled tasks satisfying filter.
func FindScheduledTasks(q *reform.Querier, filters ScheduledTasksFilter) ([]*ScheduledTask, error) {
	var args []interface{}
	var andConds []string
	if len(filters.Types) > 0 {
		p := strings.Join(q.Placeholders(1, len(filters.Types)), ", ")
		for _, fType := range filters.Types {
			args = append(args, fType)
		}
		andConds = append(andConds, fmt.Sprintf("type IN (%s)", p))
	}

	if filters.Disabled != nil {
		cond := "disabled IS "
		if *filters.Disabled {
			cond += "TRUE"
		} else {
			cond += "FALSE"
		}
		andConds = append(andConds, cond)
	}

	var tail strings.Builder
	if len(andConds) > 0 {
		tail.WriteString("WHERE ")
		tail.WriteString(strings.Join(andConds, " AND "))
		tail.WriteRune(' ')
	}
	tail.WriteString("ORDER BY created_at")

	structs, err := q.SelectAllFrom(ScheduledTaskTable, tail.String(), args...)
	if err != nil {
		return nil, err
	}
	tasks := make([]*ScheduledTask, len(structs))
	for i, s := range structs {
		tasks[i] = s.(*ScheduledTask)
	}
	return tasks, nil
}

// CreateScheduledTaskParams are params for creating new scheduled task.
type CreateScheduledTaskParams struct {
	CronExpression string
	StartAt        time.Time
	NextRun        time.Time
	Type           ScheduledTaskType
	Data           ScheduledTaskData
	Retries        uint
	RetryInterval  time.Duration
	Disabled       bool
}

// Validate checks if required params are set and valid.
func (p CreateScheduledTaskParams) Validate() error {
	switch p.Type {
	case ScheduledPrintTask:
	default:
		return status.Errorf(codes.InvalidArgument, "Unknown type: %s", p.Type)
	}
	_, err := cron.ParseStandard(p.CronExpression)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid cron expression: %v", err)
	}

	return nil
}

// CreateScheduledTask creates scheduled task.
func CreateScheduledTask(q *reform.Querier, params CreateScheduledTaskParams) (*ScheduledTask, error) {
	id := "/scheduled_task_id/" + uuid.New().String()
	if err := checkUniqueScheduledTaskID(q, id); err != nil {
		return nil, err
	}

	task := &ScheduledTask{
		ID:               id,
		CronExpression:   params.CronExpression,
		Disabled:         params.Disabled,
		StartAt:          params.StartAt,
		NextRun:          params.NextRun,
		Type:             params.Type,
		Data:             &params.Data,
		Retries:          params.Retries,
		RetryInterval:    params.RetryInterval,
		RetriesRemaining: params.Retries,
		Succeeded:        0,
		Failed:           0,
	}
	if err := q.Insert(task); err != nil {
		return nil, errors.WithStack(err)
	}
	return task, nil
}

// ChangeScheduledTaskParams are params for updating existing schedule task.
type ChangeScheduledTaskParams struct {
	NextRun          time.Time
	LastRun          time.Time
	Disable          *bool
	Retries          *uint
	RetriesRemaining *uint
	Succeeded        *uint
	Failed           *uint
}

// ChangeScheduledTask updates existing scheduled task.
func ChangeScheduledTask(q *reform.Querier, id string, params ChangeScheduledTaskParams) (*ScheduledTask, error) {
	row, err := FindScheduledTaskByID(q, id)
	if err != nil {
		return nil, err
	}
	row.NextRun = params.NextRun
	row.LastRun = params.LastRun

	if params.Disable != nil {
		row.Disabled = *params.Disable
	}

	if params.Retries != nil {
		row.Retries = *params.Retries
	}

	if params.RetriesRemaining != nil {
		row.RetriesRemaining = *params.RetriesRemaining
	}

	if params.Succeeded != nil {
		row.Succeeded = *params.Succeeded
	}

	if params.Failed != nil {
		row.Failed = *params.Failed
	}

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update scheduled task")
	}

	return row, nil
}

func RemoveScheduledTask(q *reform.Querier, id string) error {
	if _, err := FindScheduledTaskByID(q, id); err != nil {
		return err
	}
	if err := q.Delete(&ScheduledTask{ID: id}); err != nil {
		return errors.Wrap(err, "failed to delete scheduled task")
	}

	return nil
}

func checkUniqueScheduledTaskID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty schedule task ID")
	}

	location := &ScheduledTask{ID: id}
	switch err := q.Reload(location); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Scheduled task with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}
