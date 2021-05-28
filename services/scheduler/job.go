package scheduler

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/percona/pmm-managed/models"
)

type Job interface {
	Type() models.ScheduleJobType
	Do(ctx context.Context) error
}

type EchoJob models.EchoJobData

func (j *EchoJob) Do(ctx context.Context) error {
	if rand.Intn(3) == 0 {
		return fmt.Errorf("failed")
	}
	fmt.Println(j.Value)
	return nil
}
func (j *EchoJob) Type() models.ScheduleJobType {
	return models.ScheduleEchoJob
}
