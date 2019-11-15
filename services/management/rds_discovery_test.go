package management

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestRdsServiceDiscoveryIntegration(t *testing.T) {
	ctx := logger.Set(context.Background(), t.Name())

	accessKey, secretKey := tests.GetAWSKeys(t)
	rds := NewRDSService()

	instances, err := rds.Discover(ctx, accessKey, secretKey)

	//TODO: Improve this test.
	// In our current testing env with current AWS keys, 2 regions are returning errors but we don't know why for sure
	// Also, probably we can have more than 1 instance or none. PLEASE UPDATE THIS TESTS !
	assert.NotNil(t, err)
	assert.GreaterOrEqualf(t, len(instances.RdsInstances), 1, "Should have at least one instance")
}
