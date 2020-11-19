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

// Package alertmanager contains business logic of working with Alertmanager.
package alertmanager

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	httptransport "github.com/go-openapi/runtime/client"
	saas "github.com/percona-platform/saas/pkg/alert"
	"github.com/percona/pmm/api/alertmanager/amclient"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"
	"github.com/percona/pmm/api/alertmanager/amclient/general"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/percona/promconfig/alertmanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
	"gopkg.in/yaml.v3"
)

const (
	alertmanagerDataDir = "/srv/alertmanager/data"
	prometheusDir       = "/srv/prometheus"
	dirPerm             = os.FileMode(0o775)

	alertmanagerConfigPath     = "/etc/alertmanager.yml"
	alertmanagerBaseConfigPath = "/srv/alertmanager/alertmanager.base.yml"

	shippedRuleTemplatePath     = "/tmp/ia1/*.yml"
	userDefinedRuleTemplatePath = "/tmp/ia2/*.yml"

	envShippedRuleTemplatePath     = "PERCONA_TEST_SHIPPED_RULE_TEMPLATE_PATH"
	envUserDefinedRuleTemplatePath = "PERCONA_TEST_USER_DEFINED_RULE_TEMPLATE_PATH"
)

// Service is responsible for interactions with Alertmanager.
type Service struct {
	db                          *reform.DB
	l                           *logrus.Entry
	shippedRuleTemplatePath     string
	userDefinedRuleTemplatePath string
	rules                       []saas.Rule
}

// New creates new service.
func New(db *reform.DB) *Service {
	s := &Service{
		db:                          db,
		l:                           logrus.WithField("component", "alertmanager"),
		shippedRuleTemplatePath:     shippedRuleTemplatePath,
		userDefinedRuleTemplatePath: userDefinedRuleTemplatePath,
	}

	if p := os.Getenv(envShippedRuleTemplatePath); p != "" {
		s.shippedRuleTemplatePath = p
	}

	if p := os.Getenv(envUserDefinedRuleTemplatePath); p != "" {
		s.userDefinedRuleTemplatePath = p
	}

	return s
}

// Run runs Alertmanager configuration update loop until ctx is canceled.
func (svc *Service) Run(ctx context.Context) {
	svc.l.Info("Starting...")
	defer svc.l.Info("Done.")

	svc.createDataDir()
	svc.generateBaseConfig()
	svc.updateConfiguration(ctx)

	svc.collectRuleTemplates()

	// we don't have "configuration update loop" yet, so do nothing
	// TODO implement loop similar to victoriametrics.Service.Run

	<-ctx.Done()
}

// createDataDir creates Alertmanager directories if not exists in the persistent volume.
func (svc *Service) createDataDir() {
	// try to create Alertmanager data directory
	if err := os.MkdirAll(alertmanagerDataDir, dirPerm); err != nil {
		svc.l.Errorf("Cannot create datadir for Alertmanager %v.", err)
		return
	}

	// check and fix directory permissions
	alertmanagerDataDirStat, err := os.Stat(alertmanagerDataDir)
	if err != nil {
		svc.l.Errorf("Cannot get stat of %q: %v.", alertmanagerDataDir, err)
		return
	}

	if alertmanagerDataDirStat.Mode()&os.ModePerm != dirPerm {
		if err := os.Chmod(alertmanagerDataDir, dirPerm); err != nil {
			svc.l.Errorf("Cannot chmod datadir for Alertmanager %v.", err)
		}
	}

	alertmanagerDataDirSysStat := alertmanagerDataDirStat.Sys().(*syscall.Stat_t)
	aUID, aGID := int(alertmanagerDataDirSysStat.Uid), int(alertmanagerDataDirSysStat.Gid)

	prometheusDirStat, err := os.Stat(prometheusDir)
	if err != nil {
		svc.l.Errorf("Cannot get stat of %q: %v.", prometheusDir, err)
		return
	}

	prometheusDirSysStat := prometheusDirStat.Sys().(*syscall.Stat_t)
	pUID, pGID := int(prometheusDirSysStat.Uid), int(prometheusDirSysStat.Gid)
	if aUID != pUID || aGID != pGID {
		if err := os.Chown(alertmanagerDataDir, pUID, pGID); err != nil {
			svc.l.Errorf("Cannot chown datadir for Alertmanager %v.", err)
		}
	}
}

// generateBaseConfig generates /srv/alertmanager/alertmanager.base.yml if it is not present.
func (svc *Service) generateBaseConfig() {
	_, err := os.Stat(alertmanagerBaseConfigPath)
	svc.l.Debugf("%s status: %v", alertmanagerBaseConfigPath, err)

	if os.IsNotExist(err) {
		defaultBase := strings.TrimSpace(`
---
# You can edit this file; changes will be preserved.

route:
  receiver: empty
  routes: []

receivers:
  - name: empty
`) + "\n"
		err = ioutil.WriteFile(alertmanagerBaseConfigPath, []byte(defaultBase), 0644) //nolint:gosec
		svc.l.Infof("%s created: %v.", alertmanagerBaseConfigPath, err)
	}
}

// updateConfiguration updates Alertmanager configuration.
func (svc *Service) updateConfiguration(ctx context.Context) {
	// TODO split into marshalConfig and configAndReload like in victoriametrics.Service

	// if /etc/alertmanager.yml already exists, read its contents.
	var content []byte
	_, err := os.Stat(alertmanagerConfigPath)
	if err == nil {
		svc.l.Infof("%s exists, checking content", alertmanagerConfigPath)
		content, err = ioutil.ReadFile(alertmanagerConfigPath)
		if err != nil {
			svc.l.Errorf("Failed to load alertmanager config %s: %s", alertmanagerConfigPath, err)
		}
	}

	// copy the base config if `/etc/alertmanager.yml` is not present or
	// is already present but does not have any config.
	if os.IsNotExist(err) || string(content) == "---\n" {
		var cfg alertmanager.Config
		buf, err := ioutil.ReadFile(alertmanagerBaseConfigPath)
		if err != nil {
			svc.l.Errorf("Failed to load alertmanager base config %s: %s", alertmanagerBaseConfigPath, err)
			return
		}
		if err := yaml.Unmarshal(buf, &cfg); err != nil {
			svc.l.Errorf("Failed to parse alertmanager base config %s: %s.", alertmanagerBaseConfigPath, err)
			return
		}

		// TODO add custom information to this config.
		b, err := yaml.Marshal(cfg)
		if err != nil {
			svc.l.Errorf("Failed to marshal alertmanager config %s: %s.", alertmanagerConfigPath, err)
			return
		}
		b = append([]byte("# Managed by pmm-managed. DO NOT EDIT.\n---\n"), b...)

		err = ioutil.WriteFile(alertmanagerConfigPath, b, 0644)
		if err != nil {
			svc.l.Errorf("Failed to write alertmanager config %s: %s.", alertmanagerConfigPath, err)
			return
		}
	}
	svc.l.Infof("%s created", alertmanagerConfigPath)
}

func (svc *Service) collectRuleTemplates() {
	rules := make([]saas.Rule, 0)

	shippedFilePaths, err := getRuleTemplateFilePaths(svc.shippedRuleTemplatePath)
	if err != nil {
		svc.l.Errorf("Failed to get paths of template files shipped with PMM: %s.", err)
		return // keep previously loaded rules
	}

	for _, path := range shippedFilePaths {
		r, err := svc.loadRuleTemplates(path)
		if err != nil {
			svc.l.Errorf("Failed to load shipped rule template file: %s, reason: %s.", path, err)
			return // keep previously loaded rules
		}
		rules = append(rules, r...)
	}

	userDefinedFilePaths, err := getRuleTemplateFilePaths(svc.userDefinedRuleTemplatePath)
	if err != nil {
		svc.l.Errorf("Failed to get paths of user-defined template files: %s.", err)
		return // keep previously loaded rules
	}

	for _, path := range userDefinedFilePaths {
		r, err := svc.loadRuleTemplates(path)
		if err != nil {
			svc.l.Errorf("Failed to load user-defined rule template file: %s, reason: %s.", path, err)
			return // keep previously loaded rules
		}
		rules = append(rules, r...)
	}

	// TODO add method to download templates from SAAS.

	svc.rules = rules
}

func getRuleTemplateFilePaths(pattern string) ([]string, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func (svc *Service) loadRuleTemplates(file string) ([]saas.Rule, error) {
	data, err := ioutil.ReadFile(file) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "failed to read test rule template file")
	}

	// be strict about local files
	params := &saas.ParseParams{
		DisallowUnknownFields: true,
		DisallowInvalidRules:  true,
	}
	rules, err := saas.Parse(bytes.NewReader(data), params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse test rule template file")
	}

	return rules, nil
}

// SendAlerts sends given alerts. It is the caller's responsibility
// to call this method every now and then.
func (svc *Service) SendAlerts(ctx context.Context, alerts ammodels.PostableAlerts) {
	if len(alerts) == 0 {
		svc.l.Debug("0 alerts to send, exiting.")
		return
	}

	svc.l.Debugf("Sending %d alerts...", len(alerts))
	_, err := amclient.Default.Alert.PostAlerts(&alert.PostAlertsParams{
		Alerts:  alerts,
		Context: ctx,
	})
	if err != nil {
		svc.l.Error(err)
	}
}

// IsReady verifies that Alertmanager works.
func (svc *Service) IsReady(ctx context.Context) error {
	_, err := amclient.Default.General.GetStatus(&general.GetStatusParams{
		Context: ctx,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// configure default client; we use it mainly because we can't remove it from generated code
//nolint:gochecknoinits
func init() {
	amclient.Default.SetTransport(httptransport.New("127.0.0.1:9093", "/alertmanager/api/v2", []string{"http"}))
}
