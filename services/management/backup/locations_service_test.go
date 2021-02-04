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
	"testing"

	"github.com/brianvoe/gofakeit"
	backupv1beta1 "github.com/percona/pmm/api/managementpb/backup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/utils/tests"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestCreateBackupLocation(t *testing.T) {
	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	svc := NewLocationsService(db)
	t.Run("add fs", func(t *testing.T) {
		loc, err := svc.AddLocation(ctx, &backupv1beta1.AddLocationRequest{
			Name: gofakeit.Name(),
			FsConfig: &backupv1beta1.FSConfig{
				Path: "/tmp",
			},
		})
		assert.Nil(t, err)

		assert.NotEmpty(t, loc.LocationId)
	})

	t.Run("add s3", func(t *testing.T) {
		loc, err := svc.AddLocation(ctx, &backupv1beta1.AddLocationRequest{
			Name: gofakeit.Name(),
			S3Config: &backupv1beta1.S3Config{
				Endpoint:  gofakeit.URL(),
				AccessKey: "access_key",
				SecretKey: "secret_key",
			},
		})
		assert.Nil(t, err)

		assert.NotEmpty(t, loc.LocationId)
	})

	t.Run("multiple configs", func(t *testing.T) {
		_, err := svc.AddLocation(ctx, &backupv1beta1.AddLocationRequest{
			Name: gofakeit.Name(),
			FsConfig: &backupv1beta1.FSConfig{
				Path: "/tmp",
			},
			S3Config: &backupv1beta1.S3Config{
				Endpoint:  gofakeit.URL(),
				AccessKey: "access_key",
				SecretKey: "secret_key",
			},
		})
		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "Only one config is allowed"), err)

	})
}

func TestListBackupLocations(t *testing.T) {
	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	svc := NewLocationsService(db)

	_, err := svc.AddLocation(ctx, &backupv1beta1.AddLocationRequest{
		Name: gofakeit.Name(),
		FsConfig: &backupv1beta1.FSConfig{
			Path: "/tmp",
		},
	})
	require.Nil(t, err)
	_, err = svc.AddLocation(ctx, &backupv1beta1.AddLocationRequest{
		Name: gofakeit.Name(),
		S3Config: &backupv1beta1.S3Config{
			Endpoint:  gofakeit.URL(),
			AccessKey: "access_key",
			SecretKey: "secret_key",
		},
	})
	require.Nil(t, err)

	t.Run("list", func(t *testing.T) {
		res, err := svc.ListLocations(ctx, &backupv1beta1.ListLocationsRequest{})
		assert.Nil(t, err)

		assert.Len(t, res.Locations, 2)
	})

}
