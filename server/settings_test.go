package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"

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
		expected := &server.GetSettingsOKBodySettingsMetricsResolutions{
			Hr: "1s",
			Mr: "5s",
			Lr: "60s",
		}
		assert.Equal(t, expected, res.Payload.Settings.MetricsResolutions)
		expectedDataRetention := &server.GetSettingsOKBodySettingsQAN{
			DataRetention: "2592000s",
		}
		assert.Equal(t, expectedDataRetention, res.Payload.Settings.QAN)

		t.Run("ChangeSettings", func(t *testing.T) {
			teardown := func(t *testing.T) {
				t.Helper()

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						EnableTelemetry: true,
						MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
							Hr: "1s",
							Mr: "5s",
							Lr: "60s",
						},
						QAN: &server.ChangeSettingsParamsBodyQAN{
							DataRetention: "720h",
						},
					},
					Context: pmmapitests.Context,
				})
				require.NoError(t, err)
				assert.True(t, res.Payload.Settings.TelemetryEnabled)
				expected := &server.ChangeSettingsOKBodySettingsMetricsResolutions{
					Hr: "1s",
					Mr: "5s",
					Lr: "60s",
				}
				assert.Equal(t, expected, res.Payload.Settings.MetricsResolutions)
				expectedDataRetention := &server.ChangeSettingsOKBodySettingsQAN{
					DataRetention: "2592000s",
				}
				assert.Equal(t, expectedDataRetention, res.Payload.Settings.QAN)
			}

			defer teardown(t)

			t.Run("Invalid", func(t *testing.T) {
				t.Run("BothEnableAndDisable", func(t *testing.T) {
					defer teardown(t)

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
			})

			t.Run("HRInvalid", func(t *testing.T) {
				defer teardown(t)

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
				defer teardown(t)

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
				defer teardown(t)

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
				defer teardown(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						QAN: &server.ChangeSettingsParamsBodyQAN{
							DataRetention: "1",
						},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `bad Duration: time: missing unit in duration 1`)
				assert.Empty(t, res)
			})

			t.Run("DataRetentionInvalidToSmall", func(t *testing.T) {
				defer teardown(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						QAN: &server.ChangeSettingsParamsBodyQAN{
							DataRetention: "10s",
						},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `data_retention: minimal resolution is 24h`)
				assert.Empty(t, res)
			})

			t.Run("DataRetentionFractional", func(t *testing.T) {
				defer teardown(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						QAN: &server.ChangeSettingsParamsBodyQAN{
							DataRetention: "36h",
						},
					},
					Context: pmmapitests.Context,
				})
				pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `data_retention: should be a natural number of days`)
				assert.Empty(t, res)
			})

			t.Run("OK", func(t *testing.T) {
				defer teardown(t)

				res, err := serverClient.Default.Server.ChangeSettings(&server.ChangeSettingsParams{
					Body: server.ChangeSettingsBody{
						DisableTelemetry: true,
						MetricsResolutions: &server.ChangeSettingsParamsBodyMetricsResolutions{
							Hr: "2s",
							Mr: "15s",
							Lr: "2m",
						},
						QAN: &server.ChangeSettingsParamsBodyQAN{
							DataRetention: "240h",
						},
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

				getRes, err := serverClient.Default.Server.GetSettings(nil)
				require.NoError(t, err)
				assert.False(t, getRes.Payload.Settings.TelemetryEnabled)
				getExpected := &server.GetSettingsOKBodySettingsMetricsResolutions{
					Hr: "2s",
					Mr: "15s",
					Lr: "120s",
				}
				assert.Equal(t, getExpected, getRes.Payload.Settings.MetricsResolutions)
				expectedDataRetention := &server.GetSettingsOKBodySettingsQAN{
					DataRetention: "864000s",
				}
				assert.Equal(t, expectedDataRetention, getRes.Payload.Settings.QAN)
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
						defer teardown(t)

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
