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

package management

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"

	"github.com/percona/pmm/api/managementpb"
)

func TestAnnotations(t *testing.T) {
	setup := func(t *testing.T) (ctx context.Context, db *reform.DB, teardown func(t *testing.T)) {
		t.Helper()

		ctx = logger.Set(context.Background(), t.Name())
		uuid.SetRand(new(tests.IDReader))

		sqlDB := testdb.Open(t, models.SetupFixtures, nil)
		db = reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

		teardown = func(t *testing.T) {
			uuid.SetRand(nil)

			require.NoError(t, sqlDB.Close())
		}

		return
	}

	autorization := []string{"admin:admin"}

	t.Run("Non-existing service", func(t *testing.T) {
		ctx, db, teardown := setup(t)
		defer teardown(t)
		var grafanaClient mockGrafanaClient
		s := NewAnnotationService(db, &grafanaClient)

		_, err := s.AddAnnotation(ctx, autorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"no-service"},
		})
		require.Error(t, err)
	})

	t.Run("Existing service", func(t *testing.T) {
		ctx, db, teardown := setup(t)
		defer teardown(t)
		var grafanaClient mockGrafanaClient

		s := NewAnnotationService(db, &grafanaClient)

		_, err := models.AddNewService(s.db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
			ServiceName: "service-test",
			NodeID:      models.PMMServerNodeID,
			Address:     pointer.ToString("127.0.0.1"),
			Port:        pointer.ToUint16(3306),
		})
		require.NoError(t, err)

		expectedText := "Some text (Service Name: service-test)"
		grafanaClient.On("CreateAnnotation", ctx, mock.Anything, mock.Anything, expectedText, mock.Anything).Return("", nil)
		_, err = s.AddAnnotation(ctx, autorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"service-test"},
		})
		require.NoError(t, err)
	})

	t.Run("Two services", func(t *testing.T) {
		ctx, db, teardown := setup(t)
		defer teardown(t)
		var grafanaClient mockGrafanaClient
		s := NewAnnotationService(db, &grafanaClient)

		_, err := models.AddNewService(s.db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
			ServiceName: "service-test1",
			NodeID:      models.PMMServerNodeID,
			Address:     pointer.ToString("127.0.0.1"),
			Port:        pointer.ToUint16(3306),
		})
		require.NoError(t, err)

		_, err = models.AddNewService(s.db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
			ServiceName: "service-test2",
			NodeID:      models.PMMServerNodeID,
			Address:     pointer.ToString("127.0.0.1"),
			Port:        pointer.ToUint16(3306),
		})
		require.NoError(t, err)

		expectedText := "Some text (Service Name: service-test1,service-test2)"
		grafanaClient.On("CreateAnnotation", ctx, mock.Anything, mock.Anything, expectedText, mock.Anything).Return("", nil)
		_, err = s.AddAnnotation(ctx, autorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"service-test1", "service-test2"},
		})
		require.NoError(t, err)
	})

	t.Run("Non-existing node", func(t *testing.T) {
		ctx, db, teardown := setup(t)
		defer teardown(t)
		var grafanaClient mockGrafanaClient
		grafanaClient.On(mock.Anything).Return()
		s := NewAnnotationService(db, &grafanaClient)

		_, err := s.AddAnnotation(ctx, autorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"no-node"},
		})
		require.Error(t, err)
	})

	t.Run("Existing node", func(t *testing.T) {
		ctx, db, teardown := setup(t)
		defer teardown(t)
		var grafanaClient mockGrafanaClient
		grafanaClient.On("StartGrafana", mock.Anything).Return(errors.New("Annotation failed"))
		s := NewAnnotationService(db, &grafanaClient)

		_, err := models.CreateNode(s.db.Querier, models.GenericNodeType, &models.CreateNodeParams{
			NodeName: "node-test",
		})
		require.NoError(t, err)

		expectedText := "Some text (Node Name: node-test)"
		grafanaClient.On("CreateAnnotation", ctx, mock.Anything, mock.Anything, expectedText, mock.Anything).Return("", nil)
		_, err = s.AddAnnotation(ctx, autorization, &managementpb.AddAnnotationRequest{
			Text:     "Some text",
			NodeName: "node-test",
		})
		require.NoError(t, err)
	})

	t.Run("Existing service and non-existing node", func(t *testing.T) {
		ctx, db, teardown := setup(t)
		defer teardown(t)
		var grafanaClient mockGrafanaClient
		grafanaClient.On("StartGrafana", mock.Anything).Return(errors.New("Annotation failed"))
		s := NewAnnotationService(db, &grafanaClient)

		_, err := models.AddNewService(s.db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
			ServiceName: "service-test",
			NodeID:      models.PMMServerNodeID,
			Address:     pointer.ToString("127.0.0.1"),
			Port:        pointer.ToUint16(3306),
		})
		require.NoError(t, err)

		_, err = s.AddAnnotation(ctx, autorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"service-test"},
			NodeName:     "node-test",
		})
		require.Error(t, err)
	})
}
