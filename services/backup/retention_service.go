package backup

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// RetentionService represents core logic for db backup.
type RetentionService struct {
	db         *reform.DB
	l          *logrus.Entry
	removalSVC removalService
}

// NewRetentionService creates new backups logic service.
func NewRetentionService(db *reform.DB, removalSVC removalService) *RetentionService {
	return &RetentionService{
		l:          logrus.WithField("component", "management/backup/retention"),
		db:         db,
		removalSVC: removalSVC,
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

		artifacts, err = models.FindArtifacts(tx.Querier, models.ArtifactFilters{
			ScheduleID: scheduleID,
			Status:     models.SuccessBackupStatus,
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
		if err := s.removalSVC.DeleteArtifact(ctx, artifact.ID); err != nil {
			return err
		}
	}

	return nil

}
