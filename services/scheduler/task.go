package scheduler

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/percona/pmm-managed/models"
)

// Task represents task which will be run inside scheduler.
type Task interface {
	Run(ctx context.Context) error
	Type() models.ScheduledTaskType
	Data() models.ScheduledTaskData
	ID() string
	SetID(string)
}

// common implementation for all tasks.
type common struct {
	id string
}

func (c *common) ID() string {
	return c.id
}

func (c *common) SetID(id string) {
	c.id = id
}

// PrintTask implements Task for logging mesage.
type PrintTask struct {
	*common
	Message string
}

// NewPrintTask creates new task which prints message.
func NewPrintTask(message string) *PrintTask {
	return &PrintTask{
		common:  &common{},
		Message: message,
	}
}

func (j *PrintTask) Run(ctx context.Context) error {
	logrus.Info(j.Message)
	return nil
}

func (j *PrintTask) Type() models.ScheduledTaskType {
	return models.ScheduledPrintTask
}

func (j *PrintTask) Data() models.ScheduledTaskData {
	return models.ScheduledTaskData{
		Print: &models.PrintTaskData{
			Message: j.Message,
		},
	}
}
