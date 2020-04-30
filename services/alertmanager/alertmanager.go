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
	"time"

	"github.com/AlekSi/pointer"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/percona/pmm/api/alertmanager/amclient"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"
	"github.com/percona/pmm/api/alertmanager/amclient/general"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

const resendInterval = 30 * time.Second

// FIXME remove completely before release
const (
	addTestingAlerts   = true
	testingAlertsDelay = time.Second
)

// Service is responsible for interactions with Prometheus.
type Service struct {
	db             *reform.DB
	serverVersion  *version.Parsed
	agentsRegistry agentsRegistry
	r              *Registry
	l              *logrus.Entry
}

// New creates new service.
func New(db *reform.DB, v string, agentsRegistry agentsRegistry, alertsRegistry *Registry) (*Service, error) {
	serverVersion, err := version.Parse(v)
	if err != nil {
		return nil, err
	}

	return &Service{
		db:             db,
		serverVersion:  serverVersion,
		agentsRegistry: agentsRegistry,
		r:              alertsRegistry,
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
		if addTestingAlerts {
			svc.updateInventoryAlerts(ctx)
		}

		svc.sendAlerts(ctx)

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
// using /srv/alertmanager/alertmanager.base.yml as a base. See supervisord config.
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

func (svc *Service) getInventoryData(ctx context.Context) (map[string]*models.Node, map[string]*models.Service, map[string]*models.Agent, error) {
	var nodes []*models.Node
	var services []*models.Service
	var agents []*models.Agent
	err := svc.db.InTransaction(func(t *reform.TX) error {
		var e error
		nodes, e = models.FindNodes(t.Querier, models.NodeFilters{})
		if e != nil {
			return e
		}

		services, e = models.FindServices(t.Querier, models.ServiceFilters{})
		if e != nil {
			return e
		}

		agents, e = models.FindAgents(t.Querier, models.AgentFilters{})
		return e
	})
	if err != nil {
		return nil, nil, nil, err
	}

	nodesMap := make(map[string]*models.Node, len(nodes))
	for _, n := range nodes {
		nodesMap[n.NodeID] = n
	}
	servicesMap := make(map[string]*models.Service, len(services))
	for _, s := range services {
		servicesMap[s.ServiceID] = s
	}
	agentsMap := make(map[string]*models.Agent, len(agents))
	for _, a := range agents {
		agentsMap[a.AgentID] = a
	}

	return nodesMap, servicesMap, agentsMap, nil
}

// updateInventoryAlerts adds/updates alerts for inventory information in the Registry.
func (svc *Service) updateInventoryAlerts(ctx context.Context) {
	nodes, services, agents, err := svc.getInventoryData(ctx)
	if err != nil {
		svc.l.Error(err)
		return
	}

	var createdIDs []string

	for _, service := range services {
		switch service.ServiceType {
		case models.PostgreSQLServiceType:
			createdIDs = append(createdIDs, svc.updateInventoryAlertsForPostgreSQL(nodes[service.NodeID], service)...)
		}
	}

	for _, agent := range agents {
		switch agent.AgentType {
		case models.PMMAgentType:
			createdIDs = append(createdIDs, svc.updateInventoryAlertsForPMMAgent(agent, nodes[pointer.GetString(agent.RunsOnNodeID)])...)
		}
	}

	keepIDs := make(map[string]struct{})
	for _, id := range createdIDs {
		keepIDs[id] = struct{}{}
	}
	svc.r.RemovePrefix("inventory/", keepIDs)
}

func (svc *Service) updateInventoryAlertsForPostgreSQL(node *models.Node, service *models.Service) []string {
	if node == nil {
		svc.l.Error("Node not found.")
		return nil
	}

	prefix := "inventory/" + service.ServiceID + "/"
	var createdIDs []string

	name, alert, err := makeAlertPostgreSQLIsOutdated(node, service)
	if err == nil {
		id := prefix + name
		svc.r.Add(id, testingAlertsDelay, alert)
		createdIDs = append(createdIDs, id)
	} else {
		svc.l.Error(err)
	}

	return createdIDs
}

func (svc *Service) updateInventoryAlertsForPMMAgent(agent *models.Agent, node *models.Node) []string {
	if node == nil {
		svc.l.Error("Node not found.")
		return nil
	}

	prefix := "inventory/" + agent.AgentID + "/"
	var createdIDs []string

	if !svc.agentsRegistry.IsConnected(agent.AgentID) {
		name, alert, err := makeAlertPMMAgentNotConnected(agent, node)
		if err == nil {
			id := prefix + name
			svc.r.Add(id, testingAlertsDelay, alert)
			createdIDs = append(createdIDs, id)
		} else {
			svc.l.Error(err)
		}
	}

	agentVersion, err := version.Parse(pointer.GetString(agent.Version))
	if err != nil {
		svc.l.Error(err)
	}
	if agentVersion != nil && agentVersion.Less(svc.serverVersion) {
		name, alert, err := makeAlertPMMAgentIsOutdated(agent, node, svc.serverVersion.String())
		if err == nil {
			id := prefix + name
			svc.r.Add(prefix+name, testingAlertsDelay, alert)
			createdIDs = append(createdIDs, id)
		} else {
			svc.l.Error(err)
		}
	}

	return createdIDs
}

// sendAlerts sends alerts collected in the Registry.
func (svc *Service) sendAlerts(ctx context.Context) {
	alerts := svc.r.Collect()
	if len(alerts) == 0 {
		return
	}

	svc.l.Infof("Sending %d alerts...", len(alerts))
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
