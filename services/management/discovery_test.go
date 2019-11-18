package management

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm/api/managementpb"

	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestDiscoveryService(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	s := NewDiscoveryService()

	t.Run("RDS", func(t *testing.T) {
		t.Run("InvalidClientTokenId", func(t *testing.T) {
			ctx := logger.Set(context.Background(), t.Name())
			accessKey, secretKey := "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" //nolint:gosec

			instances, err := s.DiscoverRDS(ctx, &managementpb.DiscoverRDSRequest{
				AwsAccessKey: accessKey,
				AwsSecretKey: secretKey,
			})

			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "The security token included in the request is invalid."), err)
			assert.Empty(t, instances)
		})

		t.Run("DeadlineExceeded", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
			defer cancel()
			ctx = logger.Set(ctx, t.Name())
			accessKey, secretKey := "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" //nolint:gosec

			instances, err := s.DiscoverRDS(ctx, &managementpb.DiscoverRDSRequest{
				AwsAccessKey: accessKey,
				AwsSecretKey: secretKey,
			})

			tests.AssertGRPCError(t, status.New(codes.DeadlineExceeded, "Request timeout."), err)
			assert.Empty(t, instances)
		})

		t.Run("Normal", func(t *testing.T) {
			ctx := logger.Set(context.Background(), t.Name())
			accessKey, secretKey := tests.GetAWSKeys(t)

			instances, err := s.DiscoverRDS(ctx, &managementpb.DiscoverRDSRequest{
				AwsAccessKey: accessKey,
				AwsSecretKey: secretKey,
			})

			// TODO: Improve this test. https://jira.percona.com/browse/PMM-4896
			// In our current testing env with current AWS keys, 2 regions are returning errors but we don't know why for sure
			// Also, probably we can have more than 1 instance or none. PLEASE UPDATE THIS TESTS !
			assert.NoError(t, err)
			t.Logf("%+v", instances)
			assert.GreaterOrEqualf(t, len(instances.RdsInstances), 1, "Should have at least one instance")
		})
	})
}
