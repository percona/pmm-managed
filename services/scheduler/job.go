package scheduler

import (
	"context"
	"fmt"

	"github.com/percona/pmm-managed/models"
)

type Job interface {
	Do(ctx context.Context) error
	Type() models.ScheduleJobType
	Data() models.ScheduleJobData
}

type EchoJob struct {
	jobsDeps
	models.EchoJobData
}

func (j *EchoJob) Do(ctx context.Context) error {
	fmt.Println(j.Value)
	return nil
}

func (j *EchoJob) Type() models.ScheduleJobType {
	return models.ScheduleEchoJob
}

func (j *EchoJob) Data() models.ScheduleJobData {
	return models.ScheduleJobData{
		Echo: &j.EchoJobData,
	}
}
