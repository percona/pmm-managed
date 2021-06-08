package scheduler

import (
	"context"
	"fmt"

	"github.com/percona/pmm-managed/models"
)

// Task represents task which will be run inside scheduler.
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
	backupsLogicService backupsLogicService
	ServiceID           string
	LocationID          string
	Name                string
	Description         string
}

func NewMySQLBackupTask(backupsLogicService backupsLogicService, serviceID, locationID, name, description string) *mySQLBackupTask {
	return &mySQLBackupTask{
		backupsLogicService: backupsLogicService,
		ServiceID:           serviceID,
		LocationID:          locationID,
		Name:                name,
		Description:         description,
	}
}

func (t *mySQLBackupTask) Do(ctx context.Context) error {
	_, err := t.backupsLogicService.PerformBackup(ctx, t.ServiceID, t.LocationID, t.Name)
	return err
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

type mongoBackupTask struct {
	backupsLogicService backupsLogicService
	ServiceID           string
	LocationID          string
	Name                string
	Description         string
}

func NewMongoBackupTask(backupsLogicService backupsLogicService, serviceID, locationID, name, description string) *mongoBackupTask {
	return &mongoBackupTask{
		backupsLogicService: backupsLogicService,
		ServiceID:           serviceID,
		LocationID:          locationID,
		Name:                name,
		Description:         description,
	}
}

func (t *mongoBackupTask) Do(ctx context.Context) error {
	_, err := t.backupsLogicService.PerformBackup(ctx, t.ServiceID, t.LocationID, t.Name)
	return err
}

func (t *mongoBackupTask) Type() models.ScheduledTaskType {
	return models.ScheduledMongoBackupTask
}

func (t *mongoBackupTask) Data() models.ScheduledTaskData {
	return models.ScheduledTaskData{
		MongoBackupTask: &models.MongoBackupTaskData{
			ServiceID:   t.ServiceID,
			LocationID:  t.LocationID,
			Name:        t.Name,
			Description: t.Description,
		},
	}
}
