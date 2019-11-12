package management

import (
	"context"
	"testing"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
	"github.com/stretchr/testify/assert"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

func TestRdsServiceDiscovery(t *testing.T) {
	ctx := logger.Set(context.Background(), t.Name())

	sqlDB := testdb.Open(t, models.SetupFixtures)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	r := new(mockAgentsRegistry)
	r.Test(t)

	accessKey, secretKey := tests.GetAWSKeys(t)
	rds := NewRDSService(db, r)
	_, err := rds.Discover(ctx, accessKey, secretKey)
	assert.NoError(t, err)
}
