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
	q := s.db.Querier

	backups, err := models.FindBackups(q)
	if err != nil {
		return nil, err
	}

	locationIDs := make([]string, 0, len(backups))
	for _, b := range backups {
		locationIDs = append(locationIDs, b.LocationID)
	}
	locations, err := models.FindBackupLocationsByIDs(q, locationIDs)
	if err != nil {
		return nil, err
	}

	serviceIDs := make([]string, 0, len(backups))
	for _, b := range backups {
		serviceIDs = append(serviceIDs, b.ServiceID)
	}

	services, err := models.FindServicesByIDs(q, serviceIDs)
	if err != nil {
		return nil, err
	}

	backupsResponse := make([]*backupv1beta1.Backup, 0, len(backups))
	for _, b := range backups {
		convertedBackup, err := convertBackup(b, services, locations)
		if err != nil {
			return nil, err
		}
		backupsResponse = append(backupsResponse, convertedBackup)
	}
	return &backupv1beta1.ListBackupsResponse{
		Backups: backupsResponse,
	}, nil
}

func convertDataModel(dataModel models.DataModel) (*backupv1beta1.DataModel, error) {
	var dm backupv1beta1.DataModel
	switch dataModel {
	case models.PhysicalDataModel:
		dm = backupv1beta1.DataModel_PHYSICAL
	case models.LogicalDataModel:
		dm = backupv1beta1.DataModel_LOGICAL
	default:
		return nil, errors.Errorf("invalid data model '%s'", dataModel)
	}

	return &dm, nil
}

func convertBackupStatus(status models.BackupStatus) (*backupv1beta1.Status, error) {
	var s backupv1beta1.Status
	switch status {
	case models.PendingBackupStatus:
		s = backupv1beta1.Status_PENDING
	case models.InProgressBackupStatus:
		s = backupv1beta1.Status_IN_PROGRESS
	case models.PausedBackupStatus:
		s = backupv1beta1.Status_PAUSED
	case models.SuccessBackupStatus:
		s = backupv1beta1.Status_SUCCESS
	case models.ErrorBackupStatus:
		s = backupv1beta1.Status_ERROR
	default:
		return nil, errors.Errorf("invalid status '%s'", status)
	}

	return &s, nil
}

func convertBackup(
	b *models.Backup,
	services map[string]*models.Service,
	locations map[string]*models.BackupLocation,
) (*backupv1beta1.Backup, error) {
	createdAt, err := ptypes.TimestampProto(b.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	l, ok := locations[b.LocationID]
	if !ok {
		return nil, errors.Errorf(
			"failed to convert backup with id '%s': no location id '%s' in the map", b.ID, b.LocationID)
	}

	s, ok := services[b.ServiceID]
	if !ok {
		return nil, errors.Errorf(
			"failed to convert backup with id '%s': no service id '%s' in the map", b.ID, b.ServiceID)
	}

	dm, err := convertDataModel(b.DataModel)
	if err != nil {
		return nil, errors.Wrapf(err, "backup id '%s'", b.ID)
	}

	status, err := convertBackupStatus(b.Status)
	if err != nil {
		return nil, errors.Wrapf(err, "backup id '%s'", b.ID)
	}

	return &backupv1beta1.Backup{
		BackupId:     b.ID,
		Name:         b.Name,
		Vendor:       b.Vendor,
		LocationId:   b.LocationID,
		LocationName: l.Name,
		ServiceId:    b.ServiceID,
		ServiceName:  s.ServiceName,
		DataModel:    *dm,
		Status:       *status,
		CreatedAt:    createdAt,
	}, nil
}

// Check interfaces.
var (
	_ backupv1beta1.BackupsServer = (*BackupsService)(nil)
)
