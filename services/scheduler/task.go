package scheduler

import (
	"context"
	"fmt"
	"time"

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

type printTask struct {
	*common
	Message string
}

// NewPrintTask creates new task which prints message.
func NewPrintTask(message string) *printTask {
	return &printTask{
		common:  &common{},
		Message: message,
	}
}

func (j *printTask) Run(ctx context.Context) error {
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
	*common
	backupsLogicService backupsLogicService
	ServiceID           string
	LocationID          string
	Name                string
	Description         string
}

func NewMySQLBackupTask(backupsLogicService backupsLogicService, serviceID, locationID, name, description string) *mySQLBackupTask {
	return &mySQLBackupTask{
		common:              &common{},
		backupsLogicService: backupsLogicService,
		ServiceID:           serviceID,
		LocationID:          locationID,
		Name:                name,
		Description:         description,
	}
}

func (t *mySQLBackupTask) Run(ctx context.Context) error {
	name := t.Name + "_" + time.Now().Format(time.RFC3339)
	_, err := t.backupsLogicService.PerformBackup(ctx, t.ServiceID, t.LocationID, name, t.ID())
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
	*common
	backupsLogicService backupsLogicService
	ServiceID           string
	LocationID          string
	Name                string
	Description         string
}

func NewMongoBackupTask(backupsLogicService backupsLogicService, serviceID, locationID, name, description string) *mongoBackupTask {
	return &mongoBackupTask{
		common:              &common{},
		backupsLogicService: backupsLogicService,
		ServiceID:           serviceID,
		LocationID:          locationID,
		Name:                name,
		Description:         description,
	}
}

func (t *mongoBackupTask) Run(ctx context.Context) error {
	name := t.Name + "_" + time.Now().Format(time.RFC3339)
	_, err := t.backupsLogicService.PerformBackup(ctx, t.ServiceID, t.LocationID, name, t.ID())
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
