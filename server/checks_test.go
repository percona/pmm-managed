package server

import (
	"testing"

	managementClient "github.com/percona/pmm/api/managementpb/json/client"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestStartChecks(t *testing.T) {
	client := serverClient.Default.Server

	t.Run("StartSecurityChecksWithEnabledSTT", func(t *testing.T) {
		defer restoreSettingsDefaults(t)
		// Enabled STT
		res, err := client.ChangeSettings(&server.ChangeSettingsParams{
			Body: server.ChangeSettingsBody{
				EnableStt:       true,
				EnableTelemetry: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.True(t, res.Payload.Settings.SttEnabled)
		assert.True(t, res.Payload.Settings.TelemetryEnabled)
		assert.Empty(t, err)

		resp, err := managementClient.Default.SecurityChecks.StartSecurityChecks(nil)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("StartSecurityChecksWithDisabledSTT", func(t *testing.T) {
		defer restoreSettingsDefaults(t)
		// Disabled STT
		res, err := client.ChangeSettings(&server.ChangeSettingsParams{
			Body: server.ChangeSettingsBody{
				DisableStt:      true,
				EnableTelemetry: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.False(t, res.Payload.Settings.SttEnabled)
		assert.True(t, res.Payload.Settings.TelemetryEnabled)
		assert.Empty(t, err)

		resp, err := managementClient.Default.SecurityChecks.StartSecurityChecks(nil)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.FailedPrecondition, `STT is disabled.`)
		assert.Nil(t, resp)
	})
}
