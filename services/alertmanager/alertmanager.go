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
	"context"
	"io/ioutil"
	"os"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/percona/pmm/api/alertmanager/amclient"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"
	"github.com/percona/pmm/api/alertmanager/amclient/general"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/percona/promconfig"
	"github.com/percona/promconfig/alertmanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
	"gopkg.in/yaml.v3"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/dir"
)

const (
	alertmanagerDir     = "/srv/alertmanager"
	alertmanagerDataDir = "/srv/alertmanager/data"
	dirPerm             = os.FileMode(0o775)

	alertmanagerConfigPath     = "/etc/alertmanager.yml"
	alertmanagerBaseConfigPath = "/srv/alertmanager/alertmanager.base.yml"
)

// Service is responsible for interactions with Alertmanager.
type Service struct {
	db *reform.DB
	l  *logrus.Entry
}

// New creates new service.
func New(db *reform.DB) *Service {
	return &Service{
		db: db,
		l:  logrus.WithField("component", "alertmanager"),
	}
}

// Run runs Alertmanager configuration update loop until ctx is canceled.
func (svc *Service) Run(ctx context.Context) {
	svc.l.Info("Starting...")
	defer svc.l.Info("Done.")

	err := dir.CreateDataDir(alertmanagerDir, "pmm", "pmm", dirPerm)
	if err != nil {
		svc.l.Error(err)
	}
	err = dir.CreateDataDir(alertmanagerDataDir, "pmm", "pmm", dirPerm)
	if err != nil {
		svc.l.Error(err)
	}

	svc.generateBaseConfig()
	svc.updateConfiguration(ctx)

	// we don't have "configuration update loop" yet, so do nothing
	// TODO implement loop similar to victoriametrics.Service.Run

	<-ctx.Done()
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
		err = ioutil.WriteFile(alertmanagerBaseConfigPath, []byte(defaultBase), 0o644) //nolint:gosec
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

		err = svc.populateConfig(&cfg)
		if err != nil {
			svc.l.Error(err)
		}

		b, err := yaml.Marshal(cfg)
		if err != nil {
			svc.l.Errorf("Failed to marshal alertmanager config %s: %s.", alertmanagerConfigPath, err)
			return
		}

		b = append([]byte("# Managed by pmm-managed. DO NOT EDIT.\n---\n"), b...)

		err = ioutil.WriteFile(alertmanagerConfigPath, b, 0o644)
		if err != nil {
			svc.l.Errorf("Failed to write alertmanager config %s: %s.", alertmanagerConfigPath, err)
			return
		}
	}
	svc.l.Infof("%s created", alertmanagerConfigPath)
}

func (svc *Service) populateConfig(cfg *alertmanager.Config) error {
	var rules []models.Rule
	var channels []models.Channel
	e := svc.db.InTransaction(func(tx *reform.TX) error {
		var err error
		rules, err = models.FindRules(tx.Querier)
		return err
	})
	if e != nil {
		return errors.Errorf("Failed to retrieve alert rules from database: %s", e)
	}

	e = svc.db.InTransaction(func(tx *reform.TX) error {
		var err error
		channels, err = models.FindChannels(tx.Querier)
		return err
	})
	if e != nil {
		return errors.Errorf("Failed to retrieve notification channels from database: %s", e)
	}

	chanMap := make(map[string]models.Channel, len(channels))
	for _, ch := range channels {
		chanMap[ch.ID] = ch
	}

	recvSet := make(map[string]struct{}) // stores unique combinations of channel IDs
	// TODO: don't store subsets of a combination
	for _, r := range rules {
		match, _ := r.GetCustomLabels()
		match["rule_id"] = r.ID
		recv := strings.Join(r.ChannelIDs, " + ")
		recvSet[recv] = struct{}{}
		cfg.Route.Routes = append(cfg.Route.Routes, &alertmanager.Route{
			Match:          match,
			Receiver:       recv,
			RepeatInterval: promconfig.Duration(r.For),
		})
	}

	recvs, err := generateReceivers(chanMap, recvSet)
	if err != nil {
		return err
	}

	cfg.Receivers = append(cfg.Receivers, recvs...)
	return nil
}

// generateReceivers takes the channel map and a unique set of rule combinations and generates a slice of receivers
func generateReceivers(chanMap map[string]models.Channel, recvSet map[string]struct{}) ([]*alertmanager.Receiver, error) {
	var recvs []*alertmanager.Receiver
	for k := range recvSet {
		recv, err := makeReceiver(chanMap, k)
		if err != nil {
			return nil, err
		}
		recvs = append(recvs, recv)
	}
	return recvs, nil
}

// makeReceiver takes one of the unique combination of channels and turns it into a alertmanager.Receiver
func makeReceiver(chanMap map[string]models.Channel, name string) (*alertmanager.Receiver, error) {
	recv := &alertmanager.Receiver{
		Name: name,
	}

	individualChannels := strings.Split(name, " + ")

	for _, ch := range individualChannels {
		channel := chanMap[ch]
		switch channel.Type {
		case models.Email:
			recv.EmailConfigs = append(recv.EmailConfigs, &alertmanager.EmailConfig{
				// besides promconfig, To field is a slice everywhere, do we need to edit the type in promconfig?
				To: channel.EmailConfig.To[0],
			})
		case models.PagerDuty:
			pdConfig := channel.PagerDutyConfig
			if pdConfig.RoutingKey != "" {
				recv.PagerdutyConfigs = append(recv.PagerdutyConfigs, &alertmanager.PagerdutyConfig{
					RoutingKey: pdConfig.RoutingKey,
				})
				break
			}
			recv.PagerdutyConfigs = append(recv.PagerdutyConfigs, &alertmanager.PagerdutyConfig{
				ServiceKey: pdConfig.ServiceKey,
			})
		case models.Slack:
			recv.SlackConfigs = append(recv.SlackConfigs, &alertmanager.SlackConfig{
				Channel: channel.SlackConfig.Channel,
			})
		case models.WebHook:
			recv.WebhookConfigs = append(recv.WebhookConfigs, &alertmanager.WebhookConfig{
				URL:       channel.WebHookConfig.URL,
				MaxAlerts: uint64(channel.WebHookConfig.MaxAlerts),
				// TODO: add http config
			})
		default:
			return nil, errors.New("Invalid channel type")
		}

	}
	return recv, nil
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
