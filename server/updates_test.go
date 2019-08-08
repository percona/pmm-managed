package server

import (
	"strings"
	"testing"
	"time"

	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestUpdates(t *testing.T) {
	t.Run("CheckUpdates", func(t *testing.T) {
		t.Parallel()

		version, err := serverClient.Default.Server.Version(nil)
		require.NoError(t, err)
		if version.Payload.Server == nil || version.Payload.Server.Version == "" {
			t.Skip("skipping test in developer's environment")
		}

		res, err := serverClient.Default.Server.CheckUpdates(nil)
		require.NoError(t, err)

		require.NotEmpty(t, res.Payload.Installed)
		assert.True(t, strings.HasPrefix(res.Payload.Installed.Version, "2.0.0-"),
			"installed.version = %q should have '2.0.0-' prefix", res.Payload.Installed.Version)
		assert.NotEmpty(t, res.Payload.Installed.FullVersion)
		require.NotEmpty(t, res.Payload.Installed.Timestamp)
		ts := time.Time(res.Payload.Installed.Timestamp)
		hour, min, _ := ts.Clock()
		assert.Zero(t, hour, "installed.timestamp should contain only date")
		assert.Zero(t, min, "installed.timestamp should contain only date")

		require.NotEmpty(t, res.Payload.Latest)
		assert.True(t, strings.HasPrefix(res.Payload.Latest.Version, "2.0.0-"),
			"latest.version = %q should have '2.0.0-' prefix", res.Payload.Latest.Version)
		assert.NotEmpty(t, res.Payload.Latest.FullVersion)
		require.NotEmpty(t, res.Payload.Latest.Timestamp)
		ts = time.Time(res.Payload.Latest.Timestamp)
		hour, min, _ = ts.Clock()
		assert.Zero(t, hour, "latest.timestamp should contain only date")
		assert.Zero(t, min, "latest.timestamp should contain only date")

		assert.Equal(t, res.Payload.Installed.FullVersion != res.Payload.Latest.FullVersion, res.Payload.UpdateAvailable)

		assert.NotEmpty(t, res.Payload.LastCheck)

		resForce, err := serverClient.Default.Server.CheckUpdates(&server.CheckUpdatesParams{
			Body: server.CheckUpdatesBody{
				Force: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.Equal(t, res.Payload.Installed, resForce.Payload.Installed)
		assert.Equal(t, resForce.Payload.Installed.FullVersion != resForce.Payload.Latest.FullVersion, resForce.Payload.UpdateAvailable)
		assert.NotEmpty(t, resForce.Payload.LastCheck)
		assert.NotEqual(t, res.Payload.LastCheck, resForce.Payload.LastCheck)
	})
}
