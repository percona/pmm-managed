package scheduler

import (
	"context"
	"fmt"

	"github.com/percona/pmm-managed/models"
)

type Task interface {
	Do(ctx context.Context) error
	Type() models.ScheduledTaskType
	Data() models.ScheduledTaskData
}

type EchoTask struct {
	jobsDeps
	models.EchoTaskData
}

func (j *EchoTask) Do(ctx context.Context) error {
	fmt.Println(j.Value)
	return nil
}

func (j *EchoTask) Type() models.ScheduledTaskType {
	return models.ScheduledEchoTask
}

func (j *EchoTask) Data() models.ScheduledTaskData {
	return models.ScheduledTaskData{
		Echo: &j.EchoTaskData,
	}
}
