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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/percona/pmm/api/alertmanager/amclient"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"
	"github.com/percona/pmm/api/alertmanager/amclient/general"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// TODO change to smaller value
const resendInterval = 5 * time.Second

// FIXME remove completely
const test = true

// Service is responsible for interactions with Prometheus.
type Service struct {
	db             *reform.DB
	serverVersion  *version.Parsed
	agentsRegistry agentsRegistry
	l              *logrus.Entry
}

// New creates new service.
func New(db *reform.DB, v string, agentsRegistry agentsRegistry) (*Service, error) {
	serverVersion, err := version.Parse(v)
	if err != nil {
		return nil, err
	}

	return &Service{
		db:             db,
		serverVersion:  serverVersion,
		agentsRegistry: agentsRegistry,
		l:              logrus.WithField("component", "alertmanager"),
	}, nil
}

// Run runs Alertmanager configuration update loop until ctx is canceled.
func (svc *Service) Run(ctx context.Context) {
	svc.l.Info("Starting...")
	defer svc.l.Info("Done.")

	generateBaseConfig()

	t := time.NewTicker(resendInterval)
	defer t.Stop()

	for {
		if err := svc.sendAlerts(ctx); err != nil {
			svc.l.Error(err)
		}

		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// nothing, continue for loop
		}
	}
}

// generateBaseConfig generates /srv/alertmanager/alertmanager.base.yml if it is not present.
//
// TODO That's a temporary measure until we start generating /etc/alertmanager.yml
// using /srv/alertmanager/alertmanager.base.yml as a base.
func generateBaseConfig() {
	const path = "/srv/alertmanager/alertmanager.base.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultBase := strings.TrimSpace(`
---
# You can edit this file; changes will be preserved.

route:
  receiver: empty
  routes: []

receivers:
  - name: empty
`) + "\n"
		_ = ioutil.WriteFile(path, []byte(defaultBase), 0644)
	}
}

func mergeLabels(name, severity string, labels map[string]string) map[string]string {
	res := make(map[string]string, len(labels)+2)
	for k, v := range labels {
		res[k] = v
	}
	res[model.AlertNameLabel] = name
	res["severity"] = severity
	return res
}

func (svc *Service) alertsForPMMAgent(agent *models.Agent, nodes map[string]*models.Node) []*ammodels.PostableAlert {
	node := nodes[pointer.GetString(agent.RunsOnNodeID)]
	if node == nil {
		svc.l.Errorf("Node with ID %v not found.", agent.RunsOnNodeID)
		return nil
	}

	labels, err := models.MergeLabels(node, nil, agent)
	if err != nil {
		svc.l.Error(err)
		return nil
	}

	var res []*ammodels.PostableAlert

	if test {
		// to make this real, we should support a feature similar to Prometheus' `for` in alerting rules
		// to avoid alerts fluctuations

		if !svc.agentsRegistry.IsConnected(agent.AgentID) {
			res = append(res, &ammodels.PostableAlert{
				Alert: ammodels.Alert{
					// GeneratorURL: "/graph/d/pmm-inventory/pmm-inventory",
					Labels: mergeLabels("pmm_agent_not_connected", "warning", labels),
				},
				Annotations: map[string]string{
					"summary":     "pmm-agent is not connected to PMM Server",
					"description": fmt.Sprintf("Node name: %s", node.NodeName),
				},
			})
		}
	}

	agentVersion, err := version.Parse(pointer.GetString(agent.Version))
	if err != nil {
		svc.l.Error(err)
	}
	if agentVersion != nil && (test || agentVersion.Less(svc.serverVersion)) {
		res = append(res, &ammodels.PostableAlert{
			Alert: ammodels.Alert{
				// GeneratorURL: "/graph/d/pmm-inventory/pmm-inventory",
				Labels: mergeLabels("pmm_agent_outdated", "info", labels),
			},
			Annotations: map[string]string{
				"summary": "pmm-agent is outdated",
				"description": fmt.Sprintf(
					"Node name: %s\npmm-agent version: %s\nPMM Server version: %s",
					node.NodeName, agentVersion.String(), svc.serverVersion.String(),
				),
			},
		})
	}

	return res
}

func (svc *Service) sendAlerts(ctx context.Context) error {
	var nodes []*models.Node
	var agents []*models.Agent
	err := svc.db.InTransaction(func(t *reform.TX) error {
		var e error
		nodes, e = models.FindNodes(t.Querier, models.NodeFilters{})
		if e != nil {
			return e
		}

		agents, e = models.FindAgents(t.Querier, models.AgentFilters{})
		return e
	})
	if err != nil {
		return err
	}

	nodesMap := make(map[string]*models.Node, len(nodes))
	for _, n := range nodes {
		nodesMap[n.NodeID] = n
	}

	var alerts []*ammodels.PostableAlert
	for _, agent := range agents {
		switch agent.AgentType {
		case models.PMMAgentType:
			if a := svc.alertsForPMMAgent(agent, nodesMap); a != nil {
				alerts = append(alerts, a...)
			}
		}
	}

	if len(alerts) == 0 {
		return nil
	}

	svc.l.Infof("Sending %d alerts...", len(alerts))
	_, err = amclient.Default.Alert.PostAlerts(&alert.PostAlertsParams{
		Alerts:  ammodels.PostableAlerts(alerts),
		Context: ctx,
	})
	return err
}

// Check verifies that Alertmanager works.
func (svc *Service) Check(ctx context.Context) error {
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
