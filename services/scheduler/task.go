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
	ID() string
	SetID(string)
}

type common struct {
	id string
}

func (c *common) ID() string {
	return c.id
}

func (c *common) SetID(id string) {
	c.id = id
}

type printTask struct {
	*common
	Message string
}

func NewPrintTask(message string) *printTask {
	return &printTask{
		common:  &common{},
		Message: message,
	}
}

func (j *printTask) Do(ctx context.Context) error {
	fmt.Println(j.Message)
	return nil
}

func (j *printTask) Type() models.ScheduledTaskType {
	return models.ScheduledPrintTask
}

func (j *printTask) Data() models.ScheduledTaskData {
	return models.ScheduledTaskData{
		Print: &models.PrintTaskData{
			Message: j.Message,
		},
	}
}
