package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/percona/pmm/api/alertmanager/amclient"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestSettings(t *testing.T) {
	t.Run("GetSettings", func(t *testing.T) {
		res, err := serverClient.Default.Server.GetSettings(nil)
		require.NoError(t, err)
		assert.True(t, res.Payload.Settings.TelemetryEnabled)
		assert.False(t, res.Payload.Settings.SttEnabled)
		expected := &server.GetSettingsOKBodySettingsMetricsResolutions{
			Hr: "5s",
			Mr: "10s",
			Lr: "60s",
		}
		assert.Equal(t, expected, res.Payload.Settings.MetricsResolutions)
		assert.Equal(t, "2592000s", res.Payload.Settings.DataRetention)
		assert.Equal(t, []string{"aws"}, res.Payload.Settings.AWSPartitions)

		t.Run("ChangeSettings", func(t *testing.T) {
			restoreDefaults := func(t *testing.T) {
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
						DataRetention:           "2592000s",
						AWSPartitions:           []string{"aws"},
						RemoveAlertManagerURL:   true,
						RemoveAlertManagerRules: true,
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

			defer restoreDefaults(t)

			t.Run("InvalidBothEnableAndDisableSTT", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableStt:  true,
						DisableStt: true,
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Both enable_stt and disable_stt are present.`)
				assert.Empty(t, res)
			})

			t.Run("EnableSTTAndEnableTelemetry", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
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

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.True(t, resg.Payload.Settings.TelemetryEnabled)
				assert.True(t, resg.Payload.Settings.SttEnabled)
			})

			t.Run("EnableSTTAndDisableTelemetry", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableStt:        true,
						DisableTelemetry: true,
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Cannot enable STT while disabling telemetry.`)
				assert.Empty(t, res)
			})

			t.Run("DisableSTTAndEnableTelemetry", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
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

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.True(t, resg.Payload.Settings.TelemetryEnabled)
				assert.False(t, resg.Payload.Settings.SttEnabled)
			})

			t.Run("DisableSTTAndDisableTelemetry", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DisableStt:       true,
						DisableTelemetry: true,
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)
				assert.False(t, res.Payload.Settings.SttEnabled)
				assert.False(t, res.Payload.Settings.TelemetryEnabled)
				assert.Empty(t, err)

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.False(t, resg.Payload.Settings.TelemetryEnabled)
				assert.False(t, resg.Payload.Settings.SttEnabled)
			})

			t.Run("EnableSTTWhileTelemetryEnabled", func(t *testing.T) {
				defer restoreDefaults(t)

				// Ensure Telemetry is enabled
				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableTelemetry: true,
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)
				assert.True(t, res.Payload.Settings.TelemetryEnabled)

				res, err = serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableStt: true,
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)

				assert.True(t, res.Payload.Settings.SttEnabled)
				assert.Empty(t, err)

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.True(t, resg.Payload.Settings.TelemetryEnabled)
				assert.True(t, resg.Payload.Settings.SttEnabled)
			})

			t.Run("VerifyFailedChecksInAlertmanager", func(t *testing.T) {
				if !pmmapitests.RunSTTTests {
					t.Skip("skipping STT tests")
				}

				defer restoreDefaults(t)

				// Enabling STT
				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableStt: true,
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)
				assert.True(t, res.Payload.Settings.TelemetryEnabled)

				// 120 sec ping for failed checks alerts to appear in alertmanager
				var alertsCount int
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

					for _, v := range res.Payload {
						t.Logf("%+v", v)

						assert.Contains(t, v.Annotations, "summary")

						assert.Equal(t, "1", v.Labels["stt_check"])

						assert.Contains(t, v.Labels, "agent_id")
						assert.Contains(t, v.Labels, "agent_type")
						assert.Contains(t, v.Labels, "alert_id")
						assert.Contains(t, v.Labels, "alertname")
						assert.Contains(t, v.Labels, "node_id")
						assert.Contains(t, v.Labels, "node_name")
						assert.Contains(t, v.Labels, "node_type")
						assert.Contains(t, v.Labels, "service_id")
						assert.Contains(t, v.Labels, "service_name")
						assert.Contains(t, v.Labels, "service_type")
						assert.Contains(t, v.Labels, "severity")
					}
					alertsCount = len(res.Payload)
					break
				}
				assert.Greater(t, alertsCount, 0, "No alerts met")
			})

			t.Run("DisableSTTWhileItIsDisabled", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DisableStt: true,
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)
				assert.False(t, res.Payload.Settings.SttEnabled)
				assert.Empty(t, err)

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.True(t, resg.Payload.Settings.TelemetryEnabled)
				assert.False(t, resg.Payload.Settings.SttEnabled)
			})

			t.Run("STTEnabledState", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableStt: true,
					},
					Context: pmmapitests.Context,
				})

				require.NoError(t, err)
				assert.True(t, res.Payload.Settings.SttEnabled)
				assert.Empty(t, err)

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.True(t, resg.Payload.Settings.TelemetryEnabled)
				assert.True(t, resg.Payload.Settings.SttEnabled)

				t.Run("EnableSTTWhileItIsEnabled", func(t *testing.T) {
					res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							EnableStt: true,
						},
						Context: pmmapitests.Context,
					})
					require.NoError(t, err)
					assert.True(t, res.Payload.Settings.SttEnabled)
					assert.Empty(t, err)

					resg, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.True(t, resg.Payload.Settings.TelemetryEnabled)
					assert.True(t, resg.Payload.Settings.SttEnabled)
				})

				t.Run("DisableTelemetryWhileSTTEnabled", func(t *testing.T) {
					res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							DisableTelemetry: true,
						},
						Context: pmmapitests.Context,
					})
					pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Cannot disable telemetry while STT is enabled.`)
					assert.Empty(t, res)
				})
			})

			t.Run("TelemetryDisabledState", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DisableTelemetry: true,
					},
					Context: pmmapitests.Context,
				})

				require.NoError(t, err)
				assert.False(t, res.Payload.Settings.TelemetryEnabled)
				assert.Empty(t, err)

				resg, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.False(t, resg.Payload.Settings.TelemetryEnabled)
				assert.False(t, resg.Payload.Settings.SttEnabled)

				t.Run("EnableSTTWhileTelemetryDisabled", func(t *testing.T) {
					res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							EnableStt: true,
						},
						Context: pmmapitests.Context,
					})
					pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Cannot enable STT while telemetry is disabled.`)
					assert.Empty(t, res)

					resg, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.False(t, resg.Payload.Settings.TelemetryEnabled)
					assert.False(t, resg.Payload.Settings.SttEnabled)
				})

				t.Run("EnableTelemetryWhileItIsDisabled", func(t *testing.T) {
					defer restoreDefaults(t)

					res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							EnableTelemetry: true,
						},
						Context: pmmapitests.Context,
					})
					require.NoError(t, err)
					assert.True(t, res.Payload.Settings.TelemetryEnabled)

					resg, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.True(t, resg.Payload.Settings.TelemetryEnabled)
					assert.False(t, resg.Payload.Settings.SttEnabled)
				})
			})

			t.Run("InvalidBothEnableAndDisable", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableTelemetry:  true,
						DisableTelemetry: true,
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Both enable_telemetry and disable_telemetry are present.`)
				assert.Empty(t, res)
			})

			t.Run("InvalidPartition", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						AWSPartitions: []string{"aws-123"},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `aws_partitions: partition "aws-123" is invalid`)
				assert.Empty(t, res)
			})

			t.Run("TooManyPartitions", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						AWSPartitions: []string{"aws", "aws", "aws", "aws", "aws", "aws"},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `aws_partitions: list is too long`)
				assert.Empty(t, res)
			})

			t.Run("HRInvalid", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
							Hr: "1",
						},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `bad Duration: time: missing unit in duration 1`)
				assert.Empty(t, res)
			})

			t.Run("HRTooSmall", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
							Hr: "0.5s",
						},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `hr: minimal resolution is 1s`)
				assert.Empty(t, res)
			})

			t.Run("HRFractional", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
							Hr: "1.5s",
						},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `hr: should be a natural number of seconds`)
				assert.Empty(t, res)
			})

			t.Run("DataRetentionInvalid", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DataRetention: "1",
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `bad Duration: time: missing unit in duration 1`)
				assert.Empty(t, res)
			})

			t.Run("DataRetentionInvalidToSmall", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DataRetention: "10s",
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `data_retention: minimal resolution is 24h`)
				assert.Empty(t, res)
			})

			t.Run("DataRetentionFractional", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DataRetention: "36h",
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `data_retention: should be a natural number of days`)
				assert.Empty(t, res)
			})

			t.Run("InvalidSSHKey", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						SSHKey: "some-invalid-ssh-key",
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Invalid SSH key.`)
				assert.Empty(t, res)
			})

			t.Run("NoAdminUserForSSH", func(t *testing.T) {
				defer restoreDefaults(t)

				sshKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQClY/8sz3w03vA2bY6mBFgUzrvb2FIoHw8ZjUXGGClJzJg5HC3jW1m5df7TOIkx0bt6Da2UOhuCvS4o27IT1aiHXVFydppp6ghQRB6saiiW2TKlQ7B+mXatwVaOIkO381kEjgijAs0LJnNRGpqQW0ZEAxVMz4a8puaZmVNicYSVYs4kV3QZsHuqn7jHbxs5NGAO+uRRSjcuPXregsyd87RAUHkGmNrwNFln/XddMzdGMwqZOuZWuxIXBqSrSX927XGHAJlUaOmLz5etZXHzfAY1Zxfu39r66Sx95bpm3JBmc/Ewfr8T2WL0cqynkpH+3QQBCjweTHzBE+lpXHdR2se1 qsandbox"
				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						SSHKey: sshKey,
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 500, codes.Internal, `Internal server error.`)
				assert.Empty(t, res)
			})

			t.Run("OK", func(t *testing.T) {
				defer restoreDefaults(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DisableTelemetry: true,
						MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
							Hr: "2s",
							Mr: "15s",
							Lr: "2m",
						},
						DataRetention: "240h",
						AWSPartitions: []string{"aws-cn", "aws", "aws-cn"}, // duplicates are ok
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)
				assert.False(t, res.Payload.Settings.TelemetryEnabled)
				expected := &server.ChangeSettingsOKBodySettingsMetricsResolutions{
					Hr: "2s",
					Mr: "15s",
					Lr: "120s",
				}
				assert.Equal(t, expected, res.Payload.Settings.MetricsResolutions)
				assert.Equal(t, []string{"aws", "aws-cn"}, res.Payload.Settings.AWSPartitions)

				getRes, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.False(t, getRes.Payload.Settings.TelemetryEnabled)
				getExpected := &server.GetSettingsOKBodySettingsMetricsResolutions{
					Hr: "2s",
					Mr: "15s",
					Lr: "120s",
				}
				assert.Equal(t, getExpected, getRes.Payload.Settings.MetricsResolutions)
				assert.Equal(t, "864000s", res.Payload.Settings.DataRetention)
				assert.Equal(t, []string{"aws", "aws-cn"}, res.Payload.Settings.AWSPartitions)

				t.Run("DefaultsAreNotRestored", func(t *testing.T) {
					defer restoreDefaults(t)

					res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body:    server.ChangeSettingsBody{},
						Context: pmmapitests.Context,
					})
					require.NoError(t, err)
					assert.False(t, res.Payload.Settings.TelemetryEnabled)
					expected := &server.ChangeSettingsOKBodySettingsMetricsResolutions{
						Hr: "2s",
						Mr: "15s",
						Lr: "120s",
					}
					assert.Equal(t, expected, res.Payload.Settings.MetricsResolutions)
					assert.Equal(t, []string{"aws", "aws-cn"}, res.Payload.Settings.AWSPartitions)

					// Check if the values were persisted
					getRes, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.False(t, getRes.Payload.Settings.TelemetryEnabled)
					getExpected := &server.GetSettingsOKBodySettingsMetricsResolutions{
						Hr: "2s",
						Mr: "15s",
						Lr: "120s",
					}
					assert.Equal(t, getExpected, getRes.Payload.Settings.MetricsResolutions)
					assert.Equal(t, "864000s", res.Payload.Settings.DataRetention)
					assert.Equal(t, []string{"aws", "aws-cn"}, res.Payload.Settings.AWSPartitions)
				})
			})

			t.Run("AlertManager", func(t *testing.T) {
				t.Run("SetInvalid", func(t *testing.T) {
					defer restoreDefaults(t)

					url := "http://localhost:1234/"
					rules := `invalid rules`

					_, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							AlertManagerURL:   url,
							AlertManagerRules: rules,
						},
						Context: pmmapitests.Context,
					})
					pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Invalid alerting rules.`)

					gets, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.Empty(t, gets.Payload.Settings.AlertManagerURL)
					assert.Empty(t, gets.Payload.Settings.AlertManagerRules)
				})

				t.Run("SetAndRemoveInvalid", func(t *testing.T) {
					defer restoreDefaults(t)

					_, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							AlertManagerURL:       "invalid url",
							RemoveAlertManagerURL: true,
						},
						Context: pmmapitests.Context,
					})
					pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `Both alert_manager_url and remove_alert_manager_url are present.`)

					gets, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.Empty(t, gets.Payload.Settings.AlertManagerURL)
					assert.Empty(t, gets.Payload.Settings.AlertManagerRules)
				})

				t.Run("SetValid", func(t *testing.T) {
					defer restoreDefaults(t)

					url := "http://localhost:1234/"
					rules := strings.TrimSpace(`
groups:
- name: example
  rules:
  - alert: HighRequestLatency
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
      severity: page
    annotations:
      summary: High request latency
					`) + "\n"

					res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
						Body: server.ChangeSettingsBody{
							AlertManagerURL:   url,
							AlertManagerRules: rules,
						},
						Context: pmmapitests.Context,
					})
					require.NoError(t, err)
					assert.Equal(t, url, res.Payload.Settings.AlertManagerURL)
					assert.Equal(t, rules, res.Payload.Settings.AlertManagerRules)

					gets, err := serverClient.Default.Server.GetSettings(nil)
					require.NoError(t, err)
					assert.Equal(t, url, gets.Payload.Settings.AlertManagerURL)
					assert.Equal(t, rules, gets.Payload.Settings.AlertManagerRules)

					t.Run("EmptyShouldNotRemove", func(t *testing.T) {
						defer restoreDefaults(t)

						_, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
							Body:    server.ChangeSettingsBody{},
							Context: pmmapitests.Context,
						})
						require.NoError(t, err)

						gets, err = serverClient.Default.Server.GetSettings(nil)
						require.NoError(t, err)
						assert.Equal(t, url, gets.Payload.Settings.AlertManagerURL)
						assert.Equal(t, rules, gets.Payload.Settings.AlertManagerRules)
					})
				})
			})

			t.Run("grpc-gateway", func(t *testing.T) {
				// Test with pure JSON without swagger for tracking grpc-gateway behavior:
				// https://github.com/grpc-ecosystem/grpc-gateway/issues/400

				// do not use generated types as they can do extra work in generated methods
				type params struct {
					Settings struct {
						MetricsResolutions struct {
							LR string `json:"lr"`
						} `json:"metrics_resolutions"`
					} `json:"settings"`
				}
				changeURI := pmmapitests.BaseURL.ResolveReference(&url.URL{
					Path: "v1/Settings/Change",
				})
				getURI := pmmapitests.BaseURL.ResolveReference(&url.URL{
					Path: "v1/Settings/Get",
				})

				for change, get := range map[string]string{
					"59s": "59s",
					"60s": "60s",
					"61s": "61s",
					"61":  "", // no suffix => error
					"2m":  "120s",
					"1h":  "3600s",
					"1d":  "", // d suffix => error
					"1w":  "", // w suffix => error
				} {
					change, get := change, get
					t.Run(change, func(t *testing.T) {
						defer restoreDefaults(t)

						var p params
						p.Settings.MetricsResolutions.LR = change
						b, err := json.Marshal(p.Settings)
						require.NoError(t, err)
						req, err := http.NewRequest("POST", changeURI.String(), bytes.NewReader(b))
						require.NoError(t, err)
						if pmmapitests.Debug {
							b, err = httputil.DumpRequestOut(req, true)
							require.NoError(t, err)
							t.Logf("Request:\n%s", b)
						}

						resp, err := http.DefaultClient.Do(req)
						require.NoError(t, err)
						if pmmapitests.Debug {
							b, err = httputil.DumpResponse(resp, true)
							require.NoError(t, err)
							t.Logf("Response:\n%s", b)
						}
						b, err = ioutil.ReadAll(resp.Body)
						assert.NoError(t, err)
						resp.Body.Close() //nolint:errcheck

						if get == "" {
							assert.Equal(t, 400, resp.StatusCode, "response:\n%s", b)
							return
						}
						assert.Equal(t, 200, resp.StatusCode, "response:\n%s", b)

						p.Settings.MetricsResolutions.LR = ""
						err = json.Unmarshal(b, &p)
						require.NoError(t, err)
						assert.Equal(t, get, p.Settings.MetricsResolutions.LR, "Change")

						req, err = http.NewRequest("POST", getURI.String(), nil)
						require.NoError(t, err)
						if pmmapitests.Debug {
							b, err = httputil.DumpRequestOut(req, true)
							require.NoError(t, err)
							t.Logf("Request:\n%s", b)
						}

						resp, err = http.DefaultClient.Do(req)
						require.NoError(t, err)
						if pmmapitests.Debug {
							b, err = httputil.DumpResponse(resp, true)
							require.NoError(t, err)
							t.Logf("Response:\n%s", b)
						}
						b, err = ioutil.ReadAll(resp.Body)
						assert.NoError(t, err)
						resp.Body.Close() //nolint:errcheck
						assert.Equal(t, 200, resp.StatusCode, "response:\n%s", b)

						p.Settings.MetricsResolutions.LR = ""
						err = json.Unmarshal(b, &p)
						require.NoError(t, err)
						assert.Equal(t, get, p.Settings.MetricsResolutions.LR, "Get")
					})
				}
			})
		})
	})
}
