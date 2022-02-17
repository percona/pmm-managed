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

package config

import (
	_ "embed"
	"github.com/percona/pmm-managed/services/alertmanager"
	"github.com/percona/pmm-managed/services/management/ia"
	"github.com/percona/pmm-managed/services/supervisord"
	"github.com/percona/pmm-managed/services/telemetry"
	"github.com/percona/pmm-managed/services/telemetry_v2"
	"github.com/percona/pmm-managed/services/victoriametrics"
	"github.com/percona/pmm-managed/services/vmalert"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	ENV_CONFIG_PATH   = "PERCONA_PMM_CONFIG_PATH"
	defaultConfigPath = "/etc/percona/pmm/pmm-managed.yml"
)

//go:embed pmm-managed.yaml
var defaultConfig string

type Service struct {
	l      *logrus.Entry
	Config Config
}

type Config struct {
	Services struct {
		AlertManager    alertmanager.Config    `yaml:"alert_manager"`
		VMAlert         vmalert.Config         `yaml:"vmalert"`
		VictoriaMetrics victoriametrics.Config `yaml:"victoria_metrics"`
		Supervisord     supervisord.Config     `yaml:"supervisord"`
		Management      struct {
			IntegratedAlerting ia.Config `yaml:"ia"`
		} `yaml:"management"`
		Telemetry   telemetry.Config           `yaml:"telemetry"`
		TelemetryV2 telemetry_v2.ServiceConfig `yaml:"telemetry_v2"`
	} `yaml:"services"`
	Telemetry    []telemetry_v2.TelemetryConfig `yaml:"telemetry"`
	ExtraHeaders struct {
		Enabled   bool `yaml:"enabled"`
		Endpoints []struct {
			Method   string            `yaml:"method"`
			Endpoint string            `yaml:"endpoint"`
			Headers  map[string]string `yaml:"headers"`
		} `yaml:"endpoints"`
	} `yaml:"extra_headers"`
}

func NewService() *Service {
	l := logrus.WithField("component", "config")

	return &Service{
		l: l,
	}
}

func (s *Service) Load() error {
	configPath, present := os.LookupEnv(ENV_CONFIG_PATH)
	if present {
		if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "config file [%s] doen't not exit", configPath)
		}
	} else {
		s.l.Debugf("[%s] is not set, using default location [%s]", ENV_CONFIG_PATH, defaultConfigPath)
	}

	var cfg Config

	if _, err := os.Stat(configPath); err == nil {
		s.l.Trace("config exist, reading file")
		buf, err := ioutil.ReadFile(configPath)
		if err != nil {
			return errors.Wrapf(err, "error while reading config [%s]", configPath)
		}
		if err := yaml.Unmarshal(buf, &cfg); err != nil {
			return errors.Wrapf(err, "cannot unmashal config [%s]", configPath)
		}
	} else {
		s.l.Trace("config does not exist, fallback to embedded config")
		if err := yaml.Unmarshal([]byte(defaultConfig), &cfg); err != nil {
			return errors.Wrapf(err, "cannot unmashal config [%s]", configPath)
		}
	}

	cfg.Services.AlertManager.Init()
	cfg.Services.VMAlert.Init()
	cfg.Services.VictoriaMetrics.Init()
	cfg.Services.Supervisord.Init()
	cfg.Services.Management.IntegratedAlerting.Init()
	if err := cfg.Services.Telemetry.Init(); err != nil {
		return err
	}
	if err := cfg.Services.TelemetryV2.Init(cfg.Telemetry); err != nil {
		return err
	}

	cfg.configureSaasReqEnrichment()

	err := validate(cfg, configPath)
	if err != nil {
		return err
	}

	s.Config = cfg

	return nil
}

func (s *Service) Update(updater func(s *Service) error) error {
	return updater(s)
}

func validate(cfg Config, configPath string) error {
	allIds := make(map[string]bool)
	for _, config := range cfg.Telemetry {
		if _, value := allIds[config.Id]; !value {
			allIds[config.Id] = true
		} else {
			return errors.Errorf("Duplicated id [%s] found in the config file [%s]", config.Id, configPath)
		}
	}
	return nil
}
