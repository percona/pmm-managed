package server

import (
	"testing"

	"github.com/AlekSi/pointer"
	managementClient "github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/security_checks"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func restoreSettingsDefaults(t *testing.T) {
	t.Helper()

	res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
		Body: server.ChangeSettingsBody{
			DisableStt:      true,
			EnableTelemetry: true,
			MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
				Hr: "5s",
				Mr: "10s",
				Lr: "60s",
			},
			DataRetention:               "2592000s",
			AWSPartitions:               []string{"aws"},
			RemoveAlertManagerURL:       true,
			RemoveAlertManagerRules:     true,
			RemoveEmailAlertingSettings: true,
			RemoveSlackAlertingSettings: true,
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	assert.Equal(t, true, res.Payload.Settings.TelemetryEnabled)
	assert.Equal(t, false, res.Payload.Settings.SttEnabled)
	expected := &server.ChangeSettingsOKBodySettingsMetricsResolutions{
		Hr: "5s",
		Mr: "10s",
		Lr: "60s",
	}
	assert.Equal(t, expected, res.Payload.Settings.MetricsResolutions)
	assert.Equal(t, "2592000s", res.Payload.Settings.DataRetention)
	assert.Equal(t, []string{"aws"}, res.Payload.Settings.AWSPartitions)
	assert.Equal(t, "", res.Payload.Settings.AlertManagerURL)
	assert.Equal(t, "", res.Payload.Settings.AlertManagerRules)
}

func restoreCheckIntervalDefaults(t *testing.T) {
	t.Helper()

	resp, err := managementClient.Default.SecurityChecks.ListSecurityChecks(nil)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Payload.Checks)

	var params *security_checks.ChangeSecurityChecksParams

	for _, check := range resp.Payload.Checks {
		params = &security_checks.ChangeSecurityChecksParams{
			Body: security_checks.ChangeSecurityChecksBody{
				Params: []*security_checks.ParamsItems0{
					{
						Name:     check.Name,
						Interval: pointer.ToString(security_checks.ParamsItems0IntervalSTANDARD),
					},
				},
			},
			Context: pmmapitests.Context,
		}

		_, err = managementClient.Default.SecurityChecks.ChangeSecurityChecks(params)
		require.NoError(t, err)
	}
}
