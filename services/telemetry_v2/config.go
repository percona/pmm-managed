package telemetry_v2

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type ServiceConfig struct {
	Enabled      bool              `yaml:"enabled"`
	telemetry    []TelemetryConfig `yaml:"-"`
	Endpoints    EndpointsConfig   `yaml:"endpoints"`
	SaasHostname string            `yaml:"saas_hostname"`
	DataSources  struct {
		VM           *DSVM          `yaml:"VM"`
		QANDB_SELECT *DSConfigQAN   `yaml:"QANDB_SELECT"`
		PMMDB_SELECT *DSConfigPMMDB `yaml:"PMMDB_SELECT"`
	} `yaml:"datasources"`
	Reporting ReportingConfig `yaml:"reporting"`
}

type EndpointsConfig struct {
	Report string `yaml:"report"`
}

func (c *ServiceConfig) ReportEndpointURL() string {
	return fmt.Sprintf(c.Endpoints.Report, c.SaasHostname)
}

type DSConfigQAN struct {
	Timeout    time.Duration `yaml:"-"`
	TimeoutStr string        `yaml:"timeout"`
	DSN        string        `yaml:"dsn"`
}

type DSVM struct {
	Timeout    time.Duration `yaml:"-"`
	TimeoutStr string        `yaml:"timeout"`
	Address    string        `yaml:"address"`
}

type DSConfigPMMDB struct {
	Timeout                time.Duration `yaml:"-"`
	TimeoutStr             string        `yaml:"timeout"`
	UseSeparateCredentials bool          `yaml:"use_separate_credentials"`
	// Credentials used by PMM
	DSN struct {
		Scheme string
		Host   string
		DB     string
		Params string
	} `yaml:"-"`
	Credentials struct {
		Username string
		Password string
	} `yaml:"-"`
	SeparateCredentials struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"separate_credentials"`
}

type TelemetryConfig struct {
	Id      string `yaml:"id"`
	Source  string `yaml:"source"`
	Query   string `yaml:"query"`
	Summary string `yaml:"summary"`
	Data    []TelemetryConfigData
}

type TelemetryConfigData struct {
	MetricName string `yaml:"metric_name"`
	Label      string `yaml:"label"`
	Value      string `yaml:"value"`
	Column     string `yaml:"column"`
}

func (c *TelemetryConfig) MapByColumn() map[string]TelemetryConfigData {
	result := make(map[string]TelemetryConfigData, len(c.Data))
	for _, each := range c.Data {
		result[each.Column] = each
	}
	return result
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
