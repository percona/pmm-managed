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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

func checkUniqueBackupLocationID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Location ID")
	}

	location := &BackupLocation{ID: id}
	switch err := q.Reload(location); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Location with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func checkUniqueBackupLocationName(q *reform.Querier, name string) error {
	if name == "" {
		panic("empty Location Name")
	}

	var location BackupLocation
	switch err := q.FindOneTo(&location, "name", name); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Location with name %q already exists.", name)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func checkPMMServerLocationConfig(c *PMMServerLocationConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "PMM server location config is empty.")
	}
	if c.Path == "" {
		return status.Error(codes.InvalidArgument, "PMM server config path field is empty.")
	}
	return nil
}

func checkPMMClientLocationConfig(c *PMMClientLocationConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "PMM client location config is empty.")
	}
	if c.Path == "" {
		return status.Error(codes.InvalidArgument, "PMM client config path field is empty.")
	}
	return nil
}

func checkS3Config(c *S3LocationConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "S3 location config is empty.")
	}

	if c.Endpoint == "" {
		return status.Error(codes.InvalidArgument, "S3 endpoint field is empty.")
	}

	if c.AccessKey == "" {
		return status.Error(codes.InvalidArgument, "S3 accessKey field is empty.")
	}

	if c.SecretKey == "" {
		return status.Error(codes.InvalidArgument, "S3 secretKey field is empty.")
	}

	return nil
}

// FindBackupLocations returns saved backup locations configuration.
func FindBackupLocations(q *reform.Querier) ([]*BackupLocation, error) {
	rows, err := q.SelectAllFrom(BackupLocationTable, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select backup locations")
	}

	locations := make([]*BackupLocation, len(rows))
	for i, s := range rows {
		locations[i] = s.(*BackupLocation)
	}

	return locations, nil
}

// FindBackupLocationByID finds a Backup Location by it's ID.
func FindBackupLocationByID(q *reform.Querier, id string) (*BackupLocation, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Location ID.")
	}

	location := &BackupLocation{ID: id}
	switch err := q.Reload(location); err {
	case nil:
		return location, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Backup location with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// BackupLocationsConfigs groups all backup locations configs.
type BackupLocationsConfigs struct {
	PMMClientConfig *PMMClientLocationConfig
	PMMServerConfig *PMMServerLocationConfig
	S3Config        *S3LocationConfig
}

// Validate checks if there is exactly one config with required fields.
func (c BackupLocationsConfigs) Validate() error {
	configCount := 0
	if c.S3Config != nil {
		configCount++
		if err := checkS3Config(c.S3Config); err != nil {
			return err
		}
	}
	if c.PMMServerConfig != nil {
		configCount++
		if err := checkPMMServerLocationConfig(c.PMMServerConfig); err != nil {
			return err
		}
	}
	if c.PMMClientConfig != nil {
		configCount++
		if err := checkPMMClientLocationConfig(c.PMMClientConfig); err != nil {
			return err
		}
	}

	if configCount > 1 {
		return status.Error(codes.InvalidArgument, "Only one config is allowed.")
	}

	return nil
}

// FillLocation fills provided location according to backup config.
func (c BackupLocationsConfigs) FillLocation(location *BackupLocation) error {
	switch {
	case c.S3Config != nil:
		location.Type = S3BackupLocationType
		location.S3Config = c.S3Config

	case c.PMMServerConfig != nil:
		location.Type = PMMServerBackupLocationType
		location.PMMServerConfig = c.PMMServerConfig

	case c.PMMClientConfig != nil:
		location.Type = PMMClientBackupLocationType
		location.PMMClientConfig = c.PMMClientConfig
	default:
		return status.Error(codes.InvalidArgument, "Missing location type.")
	}
	return nil
}

// CreateBackupLocationParams are params for creating new backup location.
type CreateBackupLocationParams struct {
	Name        string
	Description string

	BackupLocationsConfigs
}

// CreateBackupLocation creates backup location.
func CreateBackupLocation(q *reform.Querier, params CreateBackupLocationParams) (*BackupLocation, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	id := "/location_id/" + uuid.New().String()

	if err := checkUniqueBackupLocationID(q, id); err != nil {
		return nil, err
	}

	if err := checkUniqueBackupLocationName(q, params.Name); err != nil {
		return nil, err
	}

	row := &BackupLocation{
		ID:          id,
		Name:        params.Name,
		Description: params.Description,
	}

	if err := params.FillLocation(row); err != nil {
		return nil, err
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to create backup location")
	}

	return row, nil
}

// ChangeBackupLocationParams are params for updating existing backup location.
type ChangeBackupLocationParams struct {
	Name        string
	Description string

	BackupLocationsConfigs
}

// ChangeBackupLocation updates existing location by specified locationID and params.
func ChangeBackupLocation(q *reform.Querier, locationID string, params ChangeBackupLocationParams) (*BackupLocation, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	row, err := FindBackupLocationByID(q, locationID)
	if err != nil {
		return nil, err
	}

	// remove previous configuration
	row.Type = ""
	row.PMMClientConfig = nil
	row.PMMServerConfig = nil
	row.S3Config = nil

	if params.Name != "" && params.Name != row.Name {
		row.Name = params.Name
		if err := checkUniqueBackupLocationName(q, params.Name); err != nil {
			return nil, err
		}
	}

	if params.Description != "" {
		row.Description = params.Description
	}

	if err := params.FillLocation(row); err != nil {
		return nil, err
	}

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update backup location")
	}

	return row, nil
}
