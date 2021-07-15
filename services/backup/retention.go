package backup

import (
	"context"
	"github.com/percona/pmm-managed/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

// RetentionService represents core logic for db backup.
type RetentionService struct {
	db *reform.DB
	l  *logrus.Entry
}

// NewRetentionService creates new backups logic service.
func NewRetentionService(db *reform.DB) *RetentionService {
	return &RetentionService{
		l:  logrus.WithField("component", "management/backup/retention"),
		db: db,
	}
}

func (s *RetentionService) EnforceRetention(ctx context.Context, scheduleID string) error {
	if scheduleID == "" {
		return nil
	}

	var artifacts []*models.Artifact
	var retention uint32

	txErr := s.db.InTransaction(func(tx *reform.TX) error {
		task, err := models.FindScheduledTaskByID(tx.Querier, scheduleID)
		if err != nil {
			return err
		}

		switch task.Type {
		case models.ScheduledMySQLBackupTask:
			retention = task.Data.MySQLBackupTask.Retention
		case models.ScheduledMongoDBBackupTask:
			retention = task.Data.MongoDBBackupTask.Retention
		default:
			return errors.Errorf("invalid backup type %s", task.Type)
		}

		if retention == 0 {
			return nil
		}

		artifacts, err = models.FindArtifacts(tx.Querier, &models.ArtifactFilters{
			ScheduleID: scheduleID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	if retention == 0 || int(retention) > len(artifacts) {
		return nil
	}

	for _, artifact := range artifacts[retention:] {
		_ = artifact
	}

	return nil

}
