package server

import (
	"testing"
	"time"

	"github.com/percona/pmm/api/alertmanager/amclient"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestAlertManager(t *testing.T) {
	t.Run("TestEndsAtForFailedChecksAlerts", func(t *testing.T) {
		if !pmmapitests.RunSTTTests {
			t.Skip("Skipping STT tests until we have environment: https://jira.percona.com/browse/PMM-5106")
		}

		defer restoreSettingsDefaults(t)

		// Enabling STT
		res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
			Body: server.ChangeSettingsBody{
				EnableStt: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		assert.True(t, res.Payload.Settings.SttEnabled)

		// sync with pmm-managed
		const (
			resolveTimeoutFactor  = 3
			defaultResendInterval = 2 * time.Second
		)

		// 120 sec ping for failed checks alerts to appear in alertmanager
		for i := 0; i < 120; i++ {
			res, err := amclient.Default.Alert.GetAlerts(&alert.GetAlertsParams{
				Filter:  []string{"stt_check=1"},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)
			if len(res.Payload) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			require.NotEmpty(t, res.Payload, "No alerts met")

			// TODO: Expand this test once we are silencing/removing alerts.
			alertTTL := resolveTimeoutFactor * defaultResendInterval
			for _, v := range res.Payload {
				// Since the `EndsAt` timestamp is always resolveTimeoutFactor times the
				// `resendInterval` in the future from `UpdatedAt`
				// we check whether they lie in that time alertTTL.
				assert.WithinDuration(t, time.Time(*v.EndsAt), time.Time(*v.UpdatedAt), alertTTL)
				assert.Greater(t, v.EndsAt, v.UpdatedAt)
			}
			break
		}
	})
}
