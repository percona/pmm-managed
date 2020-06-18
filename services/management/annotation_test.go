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
	"net/http"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/grafana"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
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

	ctx, db, teardown := setup(t)
	defer teardown(t)

	r := new(mockAgentsRegistry)
	r.Test(t)
	n := NewNodeService(db, r)
	_, err := n.Register(ctx, &managementpb.RegisterNodeRequest{
		NodeType: inventorypb.NodeType_GENERIC_NODE,
		NodeName: "test-node",
		Address:  "some.address.org",
		Region:   "region",
	})
	require.NoError(t, err)

	p := new(mockPrometheusService)
	p.Test(t)
	s := NewServiceService(db, r, p)
	_, err = models.AddNewService(s.db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
		ServiceName: "test-service-mysql",
		NodeID:      models.PMMServerNodeID,
		Address:     pointer.ToString("127.0.0.1"),
		Port:        pointer.ToUint16(3306),
	})
	require.NoError(t, err)

	c := grafana.NewClient("127.0.0.1:3000")
	req, err := http.NewRequest("GET", "/dummy", nil)
	require.NoError(t, err)
	req.SetBasicAuth("admin", "admin")
	authorization := []string{req.Header.Get("Authorization")}

	t.Run("Non-existing service", func(t *testing.T) {
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"no-service"},
		})
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with name "no-service" not found.`), err)
	})

	t.Run("Non-existing node", func(t *testing.T) {
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"no-node"},
		})
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with name "no-node" not found.`), err)
	})

	t.Run("Existing service", func(t *testing.T) {
		from := time.Now()
		to := from.Add(time.Second)
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"test-service-mysql"},
		})
		assert.NoError(t, err)

		annotations, err := c.FindAnnotations(ctx, from, to, authorization[0])
		require.NoError(t, err)
		for _, a := range annotations {
			if a.Text == "Some text (Service Name: test-service-mysql)" {
				assert.Equal(t, []string{"test-service-mysql"}, a.Tags)
				assert.InDelta(t, from.Unix(), a.Time.Unix(), 1)
				return
			}
		}

		assert.Fail(t, "annotation not found", "%s", annotations)
	})

	t.Run("Existing node", func(t *testing.T) {
		from := time.Now()
		to := from.Add(time.Second)
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:     "Some text",
			NodeName: "test-node",
		})
		assert.NoError(t, err)

		annotations, err := c.FindAnnotations(ctx, from, to, authorization[0])
		require.NoError(t, err)
		for _, a := range annotations {
			if a.Text == "Some text (Node Name: test-node)" {
				assert.Equal(t, []string{"test-node"}, a.Tags)
				assert.InDelta(t, from.Unix(), a.Time.Unix(), 1)
				return
			}
		}

		assert.Fail(t, "annotation not found", "%s", annotations)
	})

	t.Run("More services, one non-existing", func(t *testing.T) {
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			ServiceNames: []string{"test-service-mysql", "no-service"},
		})
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with name "no-service" not found.`), err)
	})

	t.Run("Non-existing service, non-existing node", func(t *testing.T) {
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			NodeName:     "test-node",
			ServiceNames: []string{"no-service", "no-node"},
		})
		tests.AssertGRPCError(t, status.New(codes.NotFound, `Service with name "no-service" not found.`), err)
	})

	t.Run("Existing service, existing node", func(t *testing.T) {
		from := time.Now()
		to := from.Add(time.Second)
		a := NewAnnotationService(db, c)
		_, err := a.AddAnnotation(ctx, authorization, &managementpb.AddAnnotationRequest{
			Text:         "Some text",
			NodeName:     "test-node",
			ServiceNames: []string{"test-service-mysql"},
		})
		assert.NoError(t, err)

		annotations, err := c.FindAnnotations(ctx, from, to, authorization[0])
		require.NoError(t, err)
		for _, a := range annotations {
			if a.Text == "Some text (Service Name: test-service-mysql, Node Name: test-node)" {
				assert.Equal(t, []string{"test-service-mysql", "test-node"}, a.Tags)
				assert.InDelta(t, from.Unix(), a.Time.Unix(), 1)
				return
			}
		}

		assert.Fail(t, "annotation not found", "%s", annotations)
	})
}
