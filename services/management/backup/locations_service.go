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

	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// LocationsService represents backup locations API.
type LocationsService struct {
	db *reform.DB
	s3 s3
	l  *logrus.Entry
}

// NewLocationsService creates new backup locations API service.
func NewLocationsService(db *reform.DB, s3 s3) *LocationsService {
	return &LocationsService{
		l:  logrus.WithField("component", "management/backup/locations"),
		db: db,
		s3: s3,
	}
}

// Enabled returns if service is enabled and can be used.
func (s *LocationsService) Enabled() bool {
	settings, err := models.GetSettings(s.db)
	if err != nil {
		s.l.WithError(err).Error("can't get settings")
		return false
	}
	return settings.BackupManagement.Enabled
}

// ListLocations returns list of all available backup locations.
func (s *LocationsService) ListLocations(ctx context.Context, req *backupv1beta1.ListLocationsRequest) (*backupv1beta1.ListLocationsResponse, error) {
	locations, err := models.FindBackupLocations(s.db.Querier)
	if err != nil {
		return nil, err
	}
	res := make([]*backupv1beta1.Location, len(locations))
	for i, location := range locations {
		loc, err := convertLocation(location)
		if err != nil {
			return nil, err
		}
		res[i] = loc
	}
	return &backupv1beta1.ListLocationsResponse{
		Locations: res,
	}, nil
}

// AddLocation adds new backup location.
func (s *LocationsService) AddLocation(ctx context.Context, req *backupv1beta1.AddLocationRequest) (*backupv1beta1.AddLocationResponse, error) {
	params := models.CreateBackupLocationParams{
		Name:        req.Name,
		Description: req.Description,
	}

	if req.S3Config != nil {
		params.S3Config = &models.S3LocationConfig{
			Endpoint:   req.S3Config.Endpoint,
			AccessKey:  req.S3Config.AccessKey,
			SecretKey:  req.S3Config.SecretKey,
			BucketName: req.S3Config.BucketName,
		}
	}
	if req.PmmServerConfig != nil {
		params.PMMServerConfig = &models.PMMServerLocationConfig{
			Path: req.PmmServerConfig.Path,
		}
	}

	if req.PmmClientConfig != nil {
		params.PMMClientConfig = &models.PMMClientLocationConfig{
			Path: req.PmmClientConfig.Path,
		}
	}

	if err := params.Validate(true, false); err != nil {
		return nil, err
	}

	if params.S3Config != nil {
		bucketLocation, err := s.getBucketLocation(params.S3Config)
		if err != nil {
			return nil, err
		}

		params.S3Config.BucketLocation = bucketLocation
	}

	loc, err := models.CreateBackupLocation(s.db.Querier, params)
	if err != nil {
		return nil, err
	}

	return &backupv1beta1.AddLocationResponse{
		LocationId: loc.ID,
	}, nil
}

// ChangeLocation changes existing backup location.
func (s *LocationsService) ChangeLocation(ctx context.Context, req *backupv1beta1.ChangeLocationRequest) (*backupv1beta1.ChangeLocationResponse, error) {
	params := models.ChangeBackupLocationParams{
		Name:        req.Name,
		Description: req.Description,
	}

	if req.S3Config != nil {
		params.S3Config = &models.S3LocationConfig{
			Endpoint:   req.S3Config.Endpoint,
			AccessKey:  req.S3Config.AccessKey,
			SecretKey:  req.S3Config.SecretKey,
			BucketName: req.S3Config.BucketName,
		}
	}

	if req.PmmServerConfig != nil {
		params.PMMServerConfig = &models.PMMServerLocationConfig{
			Path: req.PmmServerConfig.Path,
		}
	}

	if req.PmmClientConfig != nil {
		params.PMMClientConfig = &models.PMMClientLocationConfig{
			Path: req.PmmClientConfig.Path,
		}
	}
	if err := params.Validate(false, false); err != nil {
		return nil, err
	}

	if params.S3Config != nil {
		bucketLocation, err := s.getBucketLocation(params.S3Config)
		if err != nil {
			return nil, err
		}

		params.S3Config.BucketLocation = bucketLocation
	}

	_, err := models.ChangeBackupLocation(s.db.Querier, req.LocationId, params)
	if err != nil {
		return nil, err
	}

	return &backupv1beta1.ChangeLocationResponse{}, nil
}

// TestLocationConfig tests backup location and credentials.
func (s *LocationsService) TestLocationConfig(
	_ context.Context,
	req *backupv1beta1.TestLocationConfigRequest,
) (*backupv1beta1.TestLocationConfigResponse, error) {
	var params models.BackupLocationConfig

	if req.S3Config != nil {
		params.S3Config = &models.S3LocationConfig{
			Endpoint:   req.S3Config.Endpoint,
			AccessKey:  req.S3Config.AccessKey,
			SecretKey:  req.S3Config.SecretKey,
			BucketName: req.S3Config.BucketName,
		}

		if err := s.checkBucket(params.S3Config); err != nil {
			return nil, err
		}
	}

	if req.PmmServerConfig != nil {
		params.PMMServerConfig = &models.PMMServerLocationConfig{
			Path: req.PmmServerConfig.Path,
		}
	}

	if req.PmmClientConfig != nil {
		params.PMMClientConfig = &models.PMMClientLocationConfig{
			Path: req.PmmClientConfig.Path,
		}
	}

	if err := params.Validate(true, false); err != nil {
		return nil, err
	}

	return &backupv1beta1.TestLocationConfigResponse{}, nil
}

// RemoveLocation removes backup location.
func (s *LocationsService) RemoveLocation(ctx context.Context, req *backupv1beta1.RemoveLocationRequest) (*backupv1beta1.RemoveLocationResponse, error) {
	mode := models.RemoveRestrict
	if req.Force {
		mode = models.RemoveCascade
	}
	if err := models.RemoveBackupLocation(s.db.Querier, req.LocationId, mode); err != nil {
		return nil, err
	}
	return &backupv1beta1.RemoveLocationResponse{}, nil
}

func convertLocation(location *models.BackupLocation) (*backupv1beta1.Location, error) {
	loc := &backupv1beta1.Location{
		LocationId:  location.ID,
		Name:        location.Name,
		Description: location.Description,
	}
	switch location.Type {
	case models.PMMClientBackupLocationType:
		config := location.PMMClientConfig
		loc.Config = &backupv1beta1.Location_PmmClientConfig{
			PmmClientConfig: &backupv1beta1.PMMClientLocationConfig{
				Path: config.Path,
			},
		}
	case models.PMMServerBackupLocationType:
		config := location.PMMServerConfig
		loc.Config = &backupv1beta1.Location_PmmServerConfig{
			PmmServerConfig: &backupv1beta1.PMMServerLocationConfig{
				Path: config.Path,
			},
		}
	case models.S3BackupLocationType:
		config := location.S3Config
		loc.Config = &backupv1beta1.Location_S3Config{
			S3Config: &backupv1beta1.S3LocationConfig{
				Endpoint:   config.Endpoint,
				AccessKey:  config.AccessKey,
				SecretKey:  config.SecretKey,
				BucketName: config.BucketName,
			},
		}
	default:
		return nil, errors.Errorf("unknown backup location type %s", location.Type)
	}
	return loc, nil
}

func (s *LocationsService) getBucketLocation(c *models.S3LocationConfig) (string, error) {
	url, err := models.ParseEndpoint(c.Endpoint)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "%s", err)
	}

	secure := true
	if url.Scheme == "https" {
		secure = false
	}

	bucketLocation, err := s.s3.GetBucketLocation(url.Host, secure, c.AccessKey, c.SecretKey, c.BucketName)
	if err != nil {
		return "", err
	}

	return bucketLocation, nil
}

func (s *LocationsService) checkBucket(c *models.S3LocationConfig) error {
	url, err := models.ParseEndpoint(c.Endpoint)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "%s", err)
	}

	secure := true
	if url.Scheme == "https" {
		secure = false
	}

	exists, err := s.s3.BucketExists(url.Host, secure, c.AccessKey, c.SecretKey, c.BucketName)
	if err != nil {
		return err
	}

	if !exists {
		return status.Errorf(codes.InvalidArgument, "Bucket doesn't exist")
	}

	return nil
}

// Check interfaces.
var (
	_ backupv1beta1.LocationsServer = (*LocationsService)(nil)
)
