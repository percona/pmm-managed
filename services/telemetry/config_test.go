package telemetry

import (
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
	"time"
)

func TestServiceConfigUnmarshal(t *testing.T) {
	input := `
enabled: true
load_defaults: true
# priority is as follows, from the highest priority
#   1. this config value
#   2. PERCONA_TEST_SAAS_HOST env variable
#   3. check.percona.com
saas_hostname: "check.localhost"
endpoints:
  # %s is substituted with 'saas_hostname'
  report: https://%s/v1/telemetry/Report
datasources:
  VM:
    enabled: true
    timeout: 2s
    address: http://localhost:80/victoriametrics/
  QANDB_SELECT:
    enabled: true
    timeout: 2s
    dsn: tcp://127.0.0.1:9000?database=pmm&block_size=10000&pool_size=2
  PMMDB_SELECT:
    enabled: true
    timeout: 2s
    credentials:
      username: pmm-managed
      password: pmm-managed

reporting:
  skip_tls_verification: true
  send_on_start: true
  interval: 10s
  interval_env: "PERCONA_TEST_TELEMETRY_INTERVAL"
  retry_backoff: 1s
  retry_backoff_env: "PERCONA_TEST_TELEMETRY_RETRY_BACKOFF"
  retry_count: 2
  send_timeout: 10s`
	var actual ServiceConfig
	err := yaml.Unmarshal([]byte(input), &actual)
	require.Nil(t, err)
	assert.Equal(t, actual, &ServiceConfig{
		Enabled:      true,
		LoadDefaults: true,
		SaasHostname: "check.localhost",
		Endpoints: EndpointsConfig{
			Report: "https://%s/v1/telemetry/Report",
		},
		DataSources: struct {
			VM          *DataSourceVictoriaMetrics `yaml:"VM"`
			QanDBSelect *DSConfigQAN               `yaml:"QANDB_SELECT"` //nolint:tagliatelle
			PmmDBSelect *DSConfigPMMDB             `yaml:"PMMDB_SELECT"` //nolint:tagliatelle
		}{
			VM: &DataSourceVictoriaMetrics{
				Enabled: true,
				Timeout: time.Second * 2,
				Address: "http://localhost:80/victoriametrics/",
			},
			QanDBSelect: &DSConfigQAN{
				Enabled: true,
				Timeout: time.Second * 2,
				DSN:     "tcp://127.0.0.1:9000?database=pmm&block_size=10000&pool_size=2",
			},
			PmmDBSelect: &DSConfigPMMDB{
				Enabled: true,
				Timeout: time.Second * 2,
				Credentials: struct {
					Username string
					Password string
				}{
					Username: "pmm-managed",
					Password: "pmm-managed",
				},
			},
		},
	})
	logger, _ := test.NewNullLogger()
	err = actual.Init(logger.WithField("test", t.Name()))
	require.Nil(t, err)
}
