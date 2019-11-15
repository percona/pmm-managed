package management

import (
	"context"
	"testing"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

func TestRdsServiceDiscoveryIntegration(t *testing.T) {
	ctx := logger.Set(context.Background(), t.Name())

	sqlDB := testdb.Open(t, models.SetupFixtures)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	r := new(mockAgentsRegistry)
	r.Test(t)

	accessKey, secretKey := tests.GetAWSKeys(t)
	rds := NewRDSService(db, r)

	instances, err := rds.Discover(ctx, accessKey, secretKey)

	//TODO: Improve this test.
	// In our current testing env with current AWS keys, 2 regions are returning errors but we don't know why for sure
	// Also, probably we can have more than 1 instance or none. PLEASE UPDATE THIS TESTS !
	assert.NotNil(t, err)
	assert.GreaterOrEqualf(t, len(instances.RdsInstances), 1, "Should have at least one instance")
	require.NoError(t, sqlDB.Close())
}
