package telemetry_v2

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"aead.dev/minisign"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	l                          *logrus.Entry
	Enabled                    bool   `yaml:"enabled"`
	ConfigLocation             string `yaml:"config_location"`
	DisableSigningVerification bool   `yaml:"disable_signing_verification"` //TODO: remove this flag after testing
	Signing                    struct {
		TrustedPublicKeys       []string             `yaml:"trusted_public_keys"`
		trustedPublicKeysParsed []minisign.PublicKey `yaml:"-"`
	} `yaml:"signing"`
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

type FileConfig struct {
	Telemetry []TelemetryConfig `yaml:"telemetry"`
}

type EndpointsConfig struct {
	Report string `yaml:"report"`
}

func (c *ServiceConfig) ReportEndpointURL() string {
	return fmt.Sprintf(c.Endpoints.Report, c.SaasHostname)
}

type DSConfigQAN struct {
	Enabled    bool          `yaml:"enabled"`
	Timeout    time.Duration `yaml:"-"`
	TimeoutStr string        `yaml:"timeout"`
	DSN        string        `yaml:"dsn"`
}

type DSVM struct {
	Enabled    bool          `yaml:"enabled"`
	Timeout    time.Duration `yaml:"-"`
	TimeoutStr string        `yaml:"timeout"`
	Address    string        `yaml:"address"`
}

type DSConfigPMMDB struct {
	Enabled                bool          `yaml:"enabled"`
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
	SkipTlsVerification bool          `yaml:"skip_tls_verification"`
	SendOnStart         bool          `yaml:"send_on_start"`
	IntervalStr         string        `yaml:"interval"`
	IntervalEnv         string        `yaml:"interval_env"`
	Interval            time.Duration `yaml:"-"`
	RetryBackoffStr     string        `yaml:"retry_backoff"`
	RetryBackoffEnv     string        `yaml:"retry_backoff_env"`
	RetryBackoff        time.Duration `yaml:"-"`
	SendTimeoutStr      string        `yaml:"send_timeout"`
	SendTimeout         time.Duration `yaml:"-"`
	RetryCount          int           `yaml:"retry_count"`
}

func (c *ServiceConfig) Init(l *logrus.Entry) error {
	c.l = l

	for _, keyText := range c.Signing.TrustedPublicKeys {
		key := minisign.PublicKey{}
		if err := key.UnmarshalText([]byte(keyText)); err != nil {
			return errors.Wrap(err, "cannot parse public key")
		}
		c.Signing.trustedPublicKeysParsed = append(c.Signing.trustedPublicKeysParsed, key)
	}

	telemetry, err := c.loadConfig(c.ConfigLocation)
	if err != nil {
		return errors.Wrapf(err, "failed to load telemetry config from [%s]", c.ConfigLocation)
	}
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

	ds := c.DataSources

	vmdb := ds.VM
	if vmdb.Enabled {
		if vmdb.TimeoutStr != "" {
			timeout, err := time.ParseDuration(vmdb.TimeoutStr)
			if err != nil {
				return errors.Wrapf(err, "failed to parse duration [%s]", ds.VM.Timeout)
			}
			vmdb.Timeout = timeout
		}
	}

	qandb := ds.QANDB_SELECT
	if qandb.Enabled {
		if qandb.TimeoutStr != "" {
			timeout, err := time.ParseDuration(qandb.TimeoutStr)
			if err != nil {
				return errors.Wrapf(err, "failed to parse duration [%s]", ds.QANDB_SELECT.Timeout)
			}
			qandb.Timeout = timeout
		}
	}

	pmmdb := ds.PMMDB_SELECT
	if pmmdb.Enabled {
		if pmmdb.TimeoutStr != "" {
			timeout, err := time.ParseDuration(pmmdb.TimeoutStr)
			if err != nil {
				return errors.Wrapf(err, "failed to parse duration [%s]", ds.PMMDB_SELECT.Timeout)
			}
			pmmdb.Timeout = timeout
		}
	}

	return nil
}

func (c *ServiceConfig) loadConfig(location string) ([]TelemetryConfig, error) {
	matches, err := filepath.Glob(location)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fileCfgs := make([]FileConfig, len(matches))
	for _, match := range matches {
		var fileCfg FileConfig
		buf, err := ioutil.ReadFile(match)
		if err != nil {
			return nil, errors.Wrapf(err, "error while reading config [%s]", match)
		}
		if !c.DisableSigningVerification {
			bufSign, err := ioutil.ReadFile(match + ".minisig")
			if err != nil {
				return nil, errors.Wrapf(err, "error while reading config [%s]", match)
			}
			valid := false
			for _, publicKey := range c.Signing.trustedPublicKeysParsed {
				if ok := minisign.Verify(publicKey, buf, bufSign); ok {
					valid = true
					break
				}
			}
			if !valid {
				return nil, errors.Errorf("signature verification failed for [%s]", match)
			}
		}
		if err := yaml.Unmarshal(buf, &fileCfg); err != nil {
			return nil, errors.Wrapf(err, "cannot unmashal config [%s]", match)
		}
		fileCfgs = append(fileCfgs, fileCfg)
	}

	if err := c.validateConfig(fileCfgs); err != nil {
		c.l.Errorf(err.Error())
	}

	return c.merge(fileCfgs), nil
}

func (c *ServiceConfig) merge(cfgs []FileConfig) []TelemetryConfig {
	var result []TelemetryConfig
	ids := make(map[string]bool)
	for _, cfg := range cfgs {
		for _, each := range cfg.Telemetry {
			_, exist := ids[each.Id]
			if !exist {
				ids[each.Id] = true
				result = append(result, each)
			}
		}
	}
	return result
}

func (c *ServiceConfig) validateConfig(cfgs []FileConfig) error {
	ids := make(map[string]bool)
	for _, cfg := range cfgs {
		for _, each := range cfg.Telemetry {
			_, exist := ids[each.Id]
			if exist {
				return errors.Errorf("telemetry config ID duplication: %s", each.Id)
			}
			ids[each.Id] = true
		}
	}
	return nil
}
