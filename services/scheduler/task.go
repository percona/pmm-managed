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

type printTask struct {
	Message string
}

func NewPrintTask(message string) *printTask {
	return &printTask{
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

type mySQLBackupTask struct {
	ServiceID   string
	LocationID  string
	Name        string
	Description string
}

func NewMySQLBackupTask(serviceID, locationID, name, description string) *mySQLBackupTask {
	return &mySQLBackupTask{
		ServiceID:   serviceID,
		LocationID:  locationID,
		Name:        name,
		Description: description,
	}
}

func (t *mySQLBackupTask) Do(ctx context.Context) error {
	panic("implement me")
}

func (t *mySQLBackupTask) Type() models.ScheduledTaskType {
	return models.ScheduledMySQLBackupTask
}

func (t *mySQLBackupTask) Data() models.ScheduledTaskData {
	return models.ScheduledTaskData{
		MySQLBackupTask: &models.MySQLBackupTaskData{
			ServiceID:   t.ServiceID,
			LocationID:  t.LocationID,
			Name:        t.Name,
			Description: t.Description,
		},
	}
}
