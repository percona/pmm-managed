package server

import (
	"testing"

	"google.golang.org/grpc/codes"

	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestPanics(t *testing.T) {
	for _, mode := range []string{"panic-error", "panic-fmterror", "panic-string"} {
		mode := mode
		t.Run(mode, func(t *testing.T) {
			t.Parallel()

			res, err := serverClient.Default.Server.Version(&server.VersionParams{
				Dummy:   &mode,
				Context: pmmapitests.Context,
			})
			assert.Empty(t, res)
			pmmapitests.AssertAPIErrorf(t, err, 500, codes.Internal, "Internal server error.")
		})
	}
}
