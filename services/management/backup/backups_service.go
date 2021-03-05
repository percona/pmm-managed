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

	"github.com/golang/protobuf/ptypes"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// BackupsService represents backups API.
type BackupsService struct {
	db *reform.DB
}

// NewBackupsService creates new backup API service.
func NewBackupsService(db *reform.DB) *BackupsService {
	return &BackupsService{
		db: db,
	}
}

// ListBackups returns a list of all backups.
func (s *BackupsService) ListBackups(context.Context, *backupv1beta1.ListBackupsRequest) (*backupv1beta1.ListBackupsResponse, error) {
	backups, err := models.FindBackups(s.db.Querier)
	if err != nil {
		return nil, err
	}

	backupsResponse := make([]*backupv1beta1.Backup, 0, len(backups))
	for _, b := range backups {
		convertedBackup, err := convertBackup(b)
		if err != nil {
			return nil, err
		}
		backupsResponse = append(backupsResponse, convertedBackup)
	}
	return &backupv1beta1.ListBackupsResponse{
		Backups: backupsResponse,
	}, nil
}

func convertBackup(b *models.Backup) (*backupv1beta1.Backup, error) {
	createdAt, err := ptypes.TimestampProto(b.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	return &backupv1beta1.Backup{
		BackupId:     b.ID,
		Name:         b.Name,
		LocationName: b.LocationName,
		CreatedAt:    createdAt,
	}, nil
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
