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

package models

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
)

// FindBackups returns performed backups list.
func FindBackups(q *reform.Querier) ([]*Backup, error) {
	rows, err := q.SelectAllFrom(BackupTable, "ORDER BY created_at DESC")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select backups")
	}

	backups := make([]*Backup, 0, len(rows))
	for _, s := range rows {
		backups = append(backups, s.(*Backup))
	}

	return backups, nil
}

func findBackupByID(q *reform.Querier, id string) (*Backup, error) {
	if id == "" {
		return nil, errors.New("provided backup id is empty")
	}

	backup := &Backup{ID: id}
	switch err := q.Reload(backup); err {
	case nil:
		return backup, nil
	case reform.ErrNoRows:
		return nil, errors.Wrapf(ErrNotFound, "backup by id '%s'", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// CreateBackupParams are params for creating a new backup.
type CreateBackupParams struct {
	Name       string
	Vendor     string
	LocationID string
	ServiceID  string
	DataModel  DataModel
	Status     BackupStatus
}

// CreateBackup creates backup entry in DB.
func CreateBackup(q *reform.Querier, params CreateBackupParams) (*Backup, error) {
	if params.Name == "" {
		return nil, errors.Wrap(ErrInvalidArgument, "name shouldn't be empty")
	}
	if params.Vendor == "" {
		return nil, errors.Wrap(ErrInvalidArgument, "vendor shouldn't be empty")
	}
	if params.LocationID == "" {
		return nil, errors.Wrap(ErrInvalidArgument, "location_id shouldn't be empty")
	}
	if params.ServiceID == "" {
		return nil, errors.Wrap(ErrInvalidArgument, "service_id shouldn't be empty")
	}
	switch params.DataModel {
	case PhysicalDataModel:
	case LogicalDataModel:
	default:
		return nil, errors.Wrapf(ErrInvalidArgument, "invalid data model '%s'", params.DataModel)
	}
	switch params.Status {
	case PendingBackupStatus:
	case InProgressBackupStatus:
	case PausedBackupStatus:
	case SuccessBackupStatus:
	case ErrorBackupStatus:
	default:
		return nil, errors.Wrapf(ErrInvalidArgument, "invalid dstatus '%s'", params.Status)
	}

	id := "/backup_id/" + uuid.New().String()
	_, err := findBackupByID(q, id)
	switch {
	case err == nil:
		return nil, errors.Errorf("backup with id '%s' already exists", id)
	case errors.Is(err, ErrNotFound):
	default:
		return nil, errors.WithStack(err)
	}

	row := &Backup{
		ID:         id,
		Name:       params.Name,
		Vendor:     params.Vendor,
		LocationID: params.LocationID,
		ServiceID:  params.ServiceID,
		DataModel:  params.DataModel,
		Status:     params.Status,
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to insert backup")
	}

	return row, nil
}

// RemoveBackup removes Backup by ID.
func RemoveBackup(q *reform.Querier, id string) error {
	if _, err := findBackupByID(q, id); err != nil {
		return err
	}

	if err := q.Delete(&Backup{ID: id}); err != nil {
		return errors.Wrap(err, "failed to delete Backup")
	}
	return nil
}
