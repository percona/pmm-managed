// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

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
