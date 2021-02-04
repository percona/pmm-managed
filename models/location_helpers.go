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

func checkUniqueLocationID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Channel ID")
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

func checkFSConfig(c *FSLocationConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "FS location config is empty.")
	}
	if c.Path == "" {
		return status.Error(codes.InvalidArgument, "FS path field is empty")
	}
	return nil
}

func checkS3Config(c *S3LocationConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "S3 location config is empty.")
	}

	if c.Endpoint == "" {
		return status.Error(codes.InvalidArgument, "S3 endpoint field is empty")
	}

	if c.AccessKey == "" {
		return status.Error(codes.InvalidArgument, "S3 accessKey field is empty")
	}

	if c.SecretKey == "" {
		return status.Error(codes.InvalidArgument, "S3 secretKey field is empty")
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

// CreateBackupLocationParams are params for creating new backup location.
type CreateBackupLocationParams struct {
	Name        string
	Description string

	FSConfig *FSLocationConfig
	S3Config *S3LocationConfig
}

// CreateBackupLocation persists backup location.
func CreateBackupLocation(q *reform.Querier, params CreateBackupLocationParams) (*BackupLocation, error) {
	id := "/location_id/" + uuid.New().String()

	if err := checkUniqueLocationID(q, id); err != nil {
		return nil, err
	}

	if params.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Location name can't be empty")
	}

	row := &BackupLocation{
		ID:          id,
		Name:        params.Name,
		Description: params.Description,
	}

	if params.FSConfig != nil && params.S3Config != nil {
		return nil, status.Error(codes.InvalidArgument, "Only one config is allowed")

	}
	switch {
	case params.FSConfig != nil:
		if err := checkFSConfig(params.FSConfig); err != nil {
			return nil, err
		}
		row.Type = FSBackupLocationType
		row.FSConfig = params.FSConfig
	case params.S3Config != nil:
		if err := checkS3Config(params.S3Config); err != nil {
			return nil, err
		}
		row.Type = S3BackupLocationType
		row.S3Config = params.S3Config
	default:
		return nil, status.Error(codes.InvalidArgument, "Missing location type")
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to create backup location")
	}

	return row, nil

}
