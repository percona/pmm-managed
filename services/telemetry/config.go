package telemetry

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const (
	distributionInfoFilePath = "/srv/pmm-distribution"
	osInfoFilePath           = "/proc/version"
)

type Config struct {
	Enabled      bool            `yaml:"enabled"`
	Reporting    ReportingConfig `yaml:"reporting"`
	Endpoints    EndpointsConfig `yaml:"endpoints"`
	SaasHostname string          `yaml:"saas_hostname"`
	V1URL        string          `yaml:"v1_url"`
	V1URLEnv     string          `yaml:"v1_url_env"`
}

type EndpointsConfig struct {
	Report string `yaml:"report"`
}

type ReportingConfig struct {
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

func (c *Config) Init() error {
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

func (c *Config) ReportEndpointURL() string {
	return fmt.Sprintf(c.Endpoints.Report, c.SaasHostname)
}
