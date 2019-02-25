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

// Package prometheus contains business logic of working with Prometheus.
package prometheus

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	"github.com/percona/pmm-managed/services/prometheus/internal/prometheus/config"
	"github.com/percona/pmm-managed/utils/logger"
)

var checkFailedRE = regexp.MustCompile(`FAILED: parsing YAML file \S+: (.+)\n`)

// Service is responsible for interactions with Prometheus.
// It assumes the following:
//   * Prometheus APIs (including lifecycle) are accessible;
//   * Prometheus configuration and rule files are accessible;
//   * promtool is available.
type Service struct {
	configPath   string
	baseURL      *url.URL
	promtoolPath string
	client       *http.Client

	m sync.Mutex // for Prometheus configuration file and, by extension, for most methods
}

// NewService creates new service.
func NewService(configPath string, baseURL string, promtool string) (*Service, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Service{
		configPath:   configPath,
		baseURL:      u,
		promtoolPath: promtool,
		client:       new(http.Client),
	}, nil
}

// loadConfig loads current Prometheus configuration from file.
func (svc *Service) loadConfig() (*config.Config, error) {
	cfg, err := config.LoadFile(svc.configPath)
	if err != nil {
		return nil, errors.Wrap(err, "can't load Prometheus configuration file")
	}
	return cfg, nil
}

// saveConfigAndReload saves given Prometheus configuration to file and reloads Prometheus.
// If configuration can't be reloaded for some reason, old file is restored, and configuration is reloaded again.
func (svc *Service) saveConfigAndReload(ctx context.Context, cfg *config.Config) error {
	l := logger.Get(ctx).WithField("component", "prometheus")

	// read existing content
	old, err := ioutil.ReadFile(svc.configPath)
	if err != nil {
		return errors.WithStack(err)
	}
	fi, err := os.Stat(svc.configPath)
	if err != nil {
		return errors.WithStack(err)
	}

	// restore old content and reload in case of error
	var restore bool
	defer func() {
		if restore {
			if err = ioutil.WriteFile(svc.configPath, old, fi.Mode()); err != nil {
				l.Error(err)
			}
			if err = svc.reload(); err != nil {
				l.Error(err)
			}
		}
	}()

	// write new content to temporary file, check it
	new, err := marshalConfig(cfg)
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile("", "pmm-managed-config-")
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err = f.Write(new); err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	args := []string{"check", "config", f.Name()}
	b, err := exec.CommandContext(ctx, svc.promtoolPath, args...).CombinedOutput() //nolint:gosec
	if err != nil {
		l.Errorf("%s", b)

		// return typed error if possible
		s := string(b)
		if m := checkFailedRE.FindStringSubmatch(s); len(m) == 2 {
			return status.Error(codes.Aborted, m[1])
		}
		return errors.Wrap(err, s)
	}
	l.Debugf("%s", b)

	// write to permanent location and reload
	restore = true
	if err = ioutil.WriteFile(svc.configPath, new, fi.Mode()); err != nil {
		return errors.WithStack(err)
	}
	if err = svc.reload(); err != nil {
		return err
	}
	restore = false
	return nil
}

// marshalConfig marshals Prometheus configuration.
func marshalConfig(cfg *config.Config) ([]byte, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal Prometheus configuration file")
	}

	// TODO add comments to each cfg.ScrapeConfigs element

	b = append([]byte("# Managed by pmm-managed. DO NOT EDIT.\n---\n"), b...)
	return b, nil
}

// reload causes Prometheus to reload configuration.
func (svc *Service) reload() error {
	u := *svc.baseURL
	u.Path = path.Join(u.Path, "-", "reload")
	resp, err := svc.client.Post(u.String(), "", nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode == 200 {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.Errorf("%d: %s", resp.StatusCode, b)
}

// Check updates Prometehus configuration using information from Consul KV.
// (During PMM update prometheus.yml is overwritten, but Consul data directory is kept.)
// It returns error if configuration is not right or Prometheus is not available.
func (svc *Service) Check(ctx context.Context) error {
	l := logger.Get(ctx)

	config, err := svc.loadConfig()
	if err != nil {
		return err
	}

	if svc.baseURL == nil {
		return errors.New("URL is not set")
	}
	u := *svc.baseURL
	u.Path = path.Join(u.Path, "version")
	resp, err := svc.client.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	l.Debugf("Prometheus: %s", b)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("expected 200, got %d", resp.StatusCode)
	}

	b, err = exec.CommandContext(ctx, svc.promtoolPath, "--version").CombinedOutput() //nolint:gosec
	if err != nil {
		return errors.Wrap(err, string(b))
	}
	l.Debugf("%s", b)

	// TODO
	changed := true

	if changed {
		l.Info("Prometheus configuration updated.")
		return svc.saveConfigAndReload(ctx, config)
	}
	l.Info("Prometheus configuration not changed.")
	return nil
}
