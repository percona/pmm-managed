package management

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

type mockRDSClient struct {
	rdsiface.RDSAPI
	count func() int
}

func (m *mockRDSClient) DescribeDBInstancesPagesWithContext(ctx context.Context, input *rds.DescribeDBInstancesInput, fn func(*rds.DescribeDBInstancesOutput, bool) bool, opts ...request.Option) error {
	if m.count() < 2 {
		fn(getMockedInstances(), true)
	}
	return nil
}

func TestRdsServiceDiscoveryUnit(t *testing.T) {
	ctx := logger.Set(context.Background(), t.Name())

	sqlDB := testdb.Open(t, models.SetupFixtures)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	r := new(mockAgentsRegistry)
	r.Test(t)

	accessKey, secretKey := tests.GetAWSKeys(t)
	rds := NewRDSService(db, r)

	l := &sync.Mutex{}
	count := 0
	mockedService := &mockRDSClient{
		count: func() int {
			l.Lock()
			defer l.Unlock()
			count++
			return count
		},
	}
	rds.serviceFunc = func(s *session.Session) rdsiface.RDSAPI {
		return mockedService
	}
	instances, err := rds.Discover(ctx, accessKey, secretKey)
	assert.NoError(t, err)
	assert.Equal(t, len(instances.RdsInstances), 1)
	require.NoError(t, sqlDB.Close())
}

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

func getMockedInstances() *rds.DescribeDBInstancesOutput {
	return &rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			&rds.DBInstance{
				DBInstanceIdentifier: pointer.ToString("database-1"),
				DbiResourceId:        pointer.ToString("db-instance-id-1"),
				Endpoint: &rds.Endpoint{
					Address:      pointer.ToString("database-1.aaaaaaaaaaaa.us-east-1.rds.amazonaws.com"),
					HostedZoneId: pointer.ToString("aaaaaaaaaaaaaa"),
					Port:         pointer.ToInt64(5432),
				},
				Engine:             pointer.ToString("postgres"),
				EngineVersion:      pointer.ToString("10.6"),
				InstanceCreateTime: pointer.ToTime(time.Now()),
			},
		},
	}
}
