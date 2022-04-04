// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package telemetry provides telemetry functionality.
package telemetry

import (
	_ "embed" //nolint:golint
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"aead.dev/minisign"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ServiceConfig telemetry config.
type ServiceConfig struct {
	l                          *logrus.Entry
	Enabled                    bool   `yaml:"enabled"`
	LoadDefaults               bool   `yaml:"load_defaults"`                //nolint:tagliatelle
	ConfigLocation             string `yaml:"config_location"`              //nolint:tagliatelle
	DisableSigningVerification bool   `yaml:"disable_signing_verification"` //nolint:tagliatelle
	Signing                    struct {
		TrustedPublicKeys       []string             `yaml:"trusted_public_keys"` //nolint:tagliatelle
		trustedPublicKeysParsed []minisign.PublicKey `yaml:"-"`
	} `yaml:"signing"`
	telemetry    []Config        `yaml:"-"`
	Endpoints    EndpointsConfig `yaml:"endpoints"`
	SaasHostname string          `yaml:"saas_hostname"` //nolint:tagliatelle
	DataSources  struct {
		VM          *DataSourceVictoriaMetrics `yaml:"VM"`
		QanDBSelect *DSConfigQAN               `yaml:"QANDB_SELECT"`
		PmmDBSelect *DSConfigPMMDB             `yaml:"PMMDB_SELECT"`
	} `yaml:"datasources"`
	Reporting ReportingConfig `yaml:"reporting"`
}

// FileConfig telemetry config.
type FileConfig struct {
	Telemetry []Config `yaml:"telemetry"`
}

// EndpointsConfig telemetry config.
type EndpointsConfig struct {
	Report string `yaml:"report"`
}

// ReportEndpointURL reporting endpoint URL.
func (c *ServiceConfig) ReportEndpointURL() string {
	return fmt.Sprintf(c.Endpoints.Report, c.SaasHostname)
}

// DSConfigQAN telemetry config.
type DSConfigQAN struct {
	Enabled    bool          `yaml:"enabled"`
	Timeout    time.Duration `yaml:"-"`
	TimeoutStr string        `yaml:"timeout"`
	DSN        string        `yaml:"dsn"`
}

// DataSourceVictoriaMetrics telemetry config.
type DataSourceVictoriaMetrics struct {
	Enabled    bool          `yaml:"enabled"`
	Timeout    time.Duration `yaml:"-"`
	TimeoutStr string        `yaml:"timeout"`
	Address    string        `yaml:"address"`
}

// DSConfigPMMDB telemetry config.
type DSConfigPMMDB struct {
	Enabled                bool          `yaml:"enabled"`
	Timeout                time.Duration `yaml:"-"`
	TimeoutStr             string        `yaml:"timeout"`
	UseSeparateCredentials bool          `yaml:"use_separate_credentials"` //nolint:tagliatelle
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
	} `yaml:"separate_credentials"` //nolint:tagliatelle
}

// Config telemetry config.
type Config struct {
	ID      string `yaml:"id"`
	Source  string `yaml:"source"`
	Query   string `yaml:"query"`
	Summary string `yaml:"summary"`
	Data    []ConfigData
}

// ConfigData telemetry config.
type ConfigData struct {
	MetricName string `yaml:"metric_name"` //nolint:tagliatelle
	Label      string `yaml:"label"`
	Value      string `yaml:"value"`
	Column     string `yaml:"column"`
}

func (c *Config) mapByColumn() map[string]ConfigData {
	result := make(map[string]ConfigData, len(c.Data))
	for _, each := range c.Data {
		result[each.Column] = each
	}
	return result
}

// ReportingConfig reporting config.
type ReportingConfig struct {
	SkipTlsVerification bool          `yaml:"skip_tls_verification"` //nolint:tagliatelle
	SendOnStart         bool          `yaml:"send_on_start"`         //nolint:tagliatelle
	IntervalStr         string        `yaml:"interval"`
	IntervalEnv         string        `yaml:"interval_env"` //nolint:tagliatelle
	Interval            time.Duration `yaml:"-"`
	RetryBackoffStr     string        `yaml:"retry_backoff"`     //nolint:tagliatelle
	RetryBackoffEnv     string        `yaml:"retry_backoff_env"` //nolint:tagliatelle
	RetryBackoff        time.Duration `yaml:"-"`
	SendTimeoutStr      string        `yaml:"send_timeout"` //nolint:tagliatelle
	SendTimeout         time.Duration `yaml:"-"`
	RetryCount          int           `yaml:"retry_count"` //nolint:tagliatelle
}

//go:embed config.default.yml
var defaultConfig string

func (c *ServiceConfig) Init(l *logrus.Entry) error { //nolint:gocognit
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

	if d, err := time.ParseDuration(os.Getenv(c.Reporting.IntervalEnv)); err == nil && d > 0 {
		l.Warnf("Interval changed to %s.", d)
		c.Reporting.Interval = d
	}
	if d, err := time.ParseDuration(os.Getenv(c.Reporting.RetryBackoffEnv)); err == nil && d > 0 {
		l.Warnf("Retry backoff changed to %s.", d)
		c.Reporting.RetryBackoff = d
	}

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

	qandb := ds.QanDBSelect
	if qandb.Enabled {
		if qandb.TimeoutStr != "" {
			timeout, err := time.ParseDuration(qandb.TimeoutStr)
			if err != nil {
				return errors.Wrapf(err, "failed to parse duration [%s]", ds.QanDBSelect.Timeout)
			}
			qandb.Timeout = timeout
		}
	}

	pmmdb := ds.PmmDBSelect
	if pmmdb.Enabled {
		if pmmdb.TimeoutStr != "" {
			timeout, err := time.ParseDuration(pmmdb.TimeoutStr)
			if err != nil {
				return errors.Wrapf(err, "failed to parse duration [%s]", ds.PmmDBSelect.Timeout)
			}
			pmmdb.Timeout = timeout
		}
	}

	return nil
}

func (c *ServiceConfig) loadConfig(location string) ([]Config, error) { //nolint:cyclop
	matches, err := filepath.Glob(location)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var fileConfigs []FileConfig
	var fileCfg FileConfig
	for _, match := range matches {
		buf, err := ioutil.ReadFile(match) //nolint:gosec
		if err != nil {
			return nil, errors.Wrapf(err, "error while reading config [%s]", match)
		}
		if !c.DisableSigningVerification {
			bufSign, err := ioutil.ReadFile(match + ".minisig") //nolint:gosec
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
		fileConfigs = append(fileConfigs, fileCfg)
	}

	if c.LoadDefaults {
		defaultConfigBytes := []byte(defaultConfig)
		if err := yaml.Unmarshal(defaultConfigBytes, &fileCfg); err != nil {
			return nil, errors.Wrap(err, "cannot unmashal default config")
		}
		fileConfigs = append(fileConfigs, fileCfg)
	}

	if err := c.validateConfig(fileConfigs); err != nil {
		c.l.Errorf(err.Error())
	}

	return c.merge(fileConfigs), nil
}

func (c *ServiceConfig) merge(cfgs []FileConfig) []Config {
	var result []Config
	ids := make(map[string]bool)
	for _, cfg := range cfgs {
		for _, each := range cfg.Telemetry {
			_, exist := ids[each.ID]
			if !exist {
				ids[each.ID] = true
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
			_, exist := ids[each.ID]
			if exist {
				return errors.Errorf("telemetry config ID duplication: %s", each.ID)
			}
			ids[each.ID] = true
		}
	}
	return nil
}
