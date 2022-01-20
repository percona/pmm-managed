package telemetry_v2

import (
	"github.com/pkg/errors"
	"time"
)

const (
	ENV_DISABLED        = "PERCONA_TEST_TELEMETRY_DISABLE_SEND"
	ENV_CONFIG          = "PERCONA_TEST_TELEMETRY_FILE"
	ENV_REPORT_INTERVAL = "PERCONA_TEST_TELEMETRY_INTERVAL"
	ENV_REPORT_NO_DELAY = "PERCONA_TEST_TELEMETRY_DISABLE_START_DELAY"
)

type ServiceConfig struct {
	Enabled   bool `yaml:"enabled"`
	telemetry []TelemetryConfig
	Reporting ReportingConfig `yaml:"reporting"`
}

type TelemetryConfig struct {
	Id      string `yaml:"id"`
	Source  string `yaml:"source"`
	Query   string `yaml:"query"`
	Summary string `yaml:"summary"`
	Data    []struct {
		Name   string `yaml:"metric_name"`
		Label  string `yaml:"label"`
		Value  string `yaml:"value"`
		Column string `yaml:"column"`
	}
}

type ReportingConfig struct {
	SendOnStart     bool          `yaml:"send_on_start"`
	IntervalStr     string        `yaml:"interval"`
	IntervalEnv     string        `yaml:"interval_env"`
	Interval        time.Duration `yaml:"-"`
	RetryBackoffStr string        `yaml:"retry_backoff"`
	RetryBackoffEnv string        `yaml:"retry_backoff_env"`
	RetryBackoff    time.Duration `yaml:"-"`
	SendTimeoutStr  string        `yaml:"send_timeout"`
	SendTimeout     time.Duration `yaml:"-"`
	RetryCount      int           `yaml:"retry_count"`
}

func (c *ServiceConfig) Init(telemetry []TelemetryConfig) error {
	c.telemetry = telemetry

	reportingInterval, err := time.ParseDuration(c.Reporting.IntervalStr)
	if err != nil {
		return errors.Wrapf(err, "failed to parse duration [%s]", c.Reporting.IntervalStr)
	}
	c.Reporting.Interval = reportingInterval

	retryBackoff, err := time.ParseDuration(c.Reporting.RetryBackoffStr)
	if err != nil {
		return errors.Wrapf(err, "failed to parse duration [%s]", c.Reporting.RetryBackoffStr)
	}
	c.Reporting.RetryBackoff = retryBackoff

	sendTimeout, err := time.ParseDuration(c.Reporting.SendTimeoutStr)
	if err != nil {
		return errors.Wrapf(err, "failed to parse duration [%s]", c.Reporting.SendTimeoutStr)
	}
	c.Reporting.SendTimeout = sendTimeout

	return nil
}
