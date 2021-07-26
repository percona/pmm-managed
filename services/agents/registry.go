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

// Package agents contains business logic of working with pmm-agent.
package agents

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents/channel"
	"github.com/percona/pmm-managed/utils/logger"
)

const (
	// constants for delayed batch updates
	updateBatchDelay   = time.Second
	stateChangeTimeout = 5 * time.Second

	prometheusNamespace = "pmm_managed"
	prometheusSubsystem = "agents"
)

type pmmAgentInfo struct {
	channel         *channel.Channel
	id              string
	stateChangeChan chan struct{}
	kick            chan struct{}
}

// Registry keeps track of all connected pmm-agents.
type Registry struct {
	db   *reform.DB
	vmdb prometheusService

	rw     sync.RWMutex
	agents map[string]*pmmAgentInfo // id -> info

	roster *roster

	mAgents prom.GaugeFunc
}

// NewRegistry creates a new registry with given database connection.
func NewRegistry(db *reform.DB, vmdb prometheusService) *Registry {
	agents := make(map[string]*pmmAgentInfo)
	r := &Registry{
		db:   db,
		vmdb: vmdb,

		agents: agents,

		roster: newRoster(),

		mAgents: prom.NewGaugeFunc(prom.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "connected",
			Help:      "The current number of connected pmm-agents.",
		}, func() float64 {
			return float64(len(agents))
		}),
	}

	return r
}

// IsConnected returns true if pmm-agent with given ID is currently connected, false otherwise.
func (r *Registry) IsConnected(pmmAgentID string) bool {
	_, err := r.get(pmmAgentID)
	return err == nil
}

func (r *Registry) register(stream agentpb.Agent_ConnectServer) (*pmmAgentInfo, error) {
	ctx := stream.Context()
	l := logger.Get(ctx)
	agentMD, err := agentpb.ReceiveAgentConnectMetadata(stream)
	if err != nil {
		return nil, err
	}
	var runsOnNodeID string
	err = r.db.InTransaction(func(tx *reform.TX) error {
		runsOnNodeID, err = authenticate(agentMD, tx.Querier)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Warnf("Failed to authenticate connected pmm-agent %+v.", agentMD)
		return nil, err
	}
	l.Infof("Connected pmm-agent: %+v.", agentMD)

	serverMD := agentpb.ServerConnectMetadata{
		AgentRunsOnNodeID: runsOnNodeID,
		ServerVersion:     version.Version,
	}
	l.Debugf("Sending metadata: %+v.", serverMD)
	if err = agentpb.SendServerConnectMetadata(stream, &serverMD); err != nil {
		return nil, err
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	// do not use r.get() - r.rw is already locked
	if agent := r.agents[agentMD.ID]; agent != nil {
		// pmm-agent with the same ID can still be connected in two cases:
		//   1. Someone uses the same ID by mistake, glitch, or malicious intent.
		//   2. pmm-agent detects broken connection and reconnects,
		//      but pmm-managed still thinks that the previous connection is okay.
		// In both cases, kick it.
		l.Warnf("Another pmm-agent with ID %q is already connected.", agentMD.ID)
		r.Kick(ctx, agentMD.ID)
	}

	agent := &pmmAgentInfo{
		channel:         channel.New(stream),
		id:              agentMD.ID,
		stateChangeChan: make(chan struct{}, 1),
		kick:            make(chan struct{}),
	}
	r.agents[agentMD.ID] = agent
	return agent, nil
}

func authenticate(md *agentpb.AgentConnectMetadata, q *reform.Querier) (string, error) {
	if md.ID == "" {
		return "", status.Error(codes.PermissionDenied, "Empty Agent ID.")
	}

	agent, err := models.FindAgentByID(q, md.ID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", status.Errorf(codes.PermissionDenied, "No Agent with ID %q.", md.ID)
		}
		return "", errors.Wrap(err, "failed to find agent")
	}

	if agent.AgentType != models.PMMAgentType {
		return "", status.Errorf(codes.PermissionDenied, "No pmm-agent with ID %q.", md.ID)
	}

	if pointer.GetString(agent.RunsOnNodeID) == "" {
		return "", status.Errorf(codes.PermissionDenied, "Can't get 'runs_on_node_id' for pmm-agent with ID %q.", md.ID)
	}

	agentVersion, err := version.Parse(md.Version)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Can't parse 'version' for pmm-agent with ID %q.", md.ID)
	}

	if err := addOrRemoveVMAgent(q, md.ID, pointer.GetString(agent.RunsOnNodeID), agentVersion); err != nil {
		return "", err
	}

	agent.Version = &md.Version
	if err := q.Update(agent); err != nil {
		return "", errors.Wrap(err, "failed to update agent")
	}

	return pointer.GetString(agent.RunsOnNodeID), nil
}

// unregister removes pmm-agent with given ID from the registry.
func (r *Registry) unregister(pmmAgentID string) *pmmAgentInfo {
	r.rw.Lock()
	defer r.rw.Unlock()

	// We do not check that pmmAgentID is in fact ID of existing pmm-agent because
	// it may be already deleted from the database, that's why we unregister it.

	agent := r.agents[pmmAgentID]
	if agent == nil {
		return nil
	}

	delete(r.agents, pmmAgentID)
	r.roster.clear(pmmAgentID)
	return agent
}

// addOrRemoveVMAgent - creates vmAgent agentType if pmm-agent's version supports it and agent not exists yet,
// otherwise ensures that vmAgent not exist for pmm-agent and pmm-agent's agents don't have push_metrics mode,
// removes it if needed.
func addOrRemoveVMAgent(q *reform.Querier, pmmAgentID, runsOnNodeID string, pmmAgentVersion *version.Parsed) error {
	if pmmAgentVersion.Less(models.PMMAgentWithPushMetricsSupport) {
		// ensure that vmagent not exists and agents dont have push_metrics.
		return removeVMAgentFromPMMAgent(q, pmmAgentID)
	}
	return addVMAgentToPMMAgent(q, pmmAgentID, runsOnNodeID)
}

func addVMAgentToPMMAgent(q *reform.Querier, pmmAgentID, runsOnNodeID string) error {
	// TODO remove it after fix
	// https://jira.percona.com/browse/PMM-4420
	if runsOnNodeID == "pmm-server" {
		return nil
	}
	vmAgentType := models.VMAgentType
	vmAgent, err := models.FindAgents(q, models.AgentFilters{PMMAgentID: pmmAgentID, AgentType: &vmAgentType})
	if err != nil {
		return status.Errorf(codes.Internal, "Can't get 'vmAgent' for pmm-agent with ID %q", pmmAgentID)
	}
	if len(vmAgent) == 0 {
		if _, err := models.CreateAgent(q, models.VMAgentType, &models.CreateAgentParams{
			PMMAgentID:  pmmAgentID,
			PushMetrics: true,
			NodeID:      runsOnNodeID,
		}); err != nil {
			return errors.Wrapf(err, "Can't create 'vmAgent' for pmm-agent with ID %q", pmmAgentID)
		}
	}
	return nil
}

func removeVMAgentFromPMMAgent(q *reform.Querier, pmmAgentID string) error {
	vmAgentType := models.VMAgentType
	vmAgent, err := models.FindAgents(q, models.AgentFilters{PMMAgentID: pmmAgentID, AgentType: &vmAgentType})
	if err != nil {
		return status.Errorf(codes.Internal, "Can't get 'vmAgent' for pmm-agent with ID %q", pmmAgentID)
	}
	if len(vmAgent) != 0 {
		for _, agent := range vmAgent {
			if _, err := models.RemoveAgent(q, agent.AgentID, models.RemoveRestrict); err != nil {
				return errors.Wrapf(err, "Can't remove 'vmAgent' for pmm-agent with ID %q", pmmAgentID)
			}
		}
	}
	agents, err := models.FindAgents(q, models.AgentFilters{PMMAgentID: pmmAgentID})
	if err != nil {
		return errors.Wrapf(err, "Can't find agents for pmm-agent with ID %q", pmmAgentID)
	}
	for _, agent := range agents {
		if agent.PushMetrics {
			logrus.Warnf("disabling push_metrics for agent with unsupported version ID %q with pmm-agent ID %q", agent.AgentID, pmmAgentID)
			agent.PushMetrics = false
			if err := q.Update(agent); err != nil {
				return errors.Wrapf(err, "Can't set push_metrics=false for agent %q at pmm-agent with ID %q", agent.AgentID, pmmAgentID)
			}
		}
	}
	return nil
}

// Kick unregisters and forcefully disconnects pmm-agent with given ID.
func (r *Registry) Kick(ctx context.Context, pmmAgentID string) {
	agent := r.unregister(pmmAgentID)
	if agent == nil {
		return
	}

	l := logger.Get(ctx)
	l.Debugf("pmm-agent with ID %q will be kicked in a moment.", pmmAgentID)

	// see Run method
	close(agent.kick)

	// Do not close agent.stateChangeChan to avoid breaking RequestStateUpdate;
	// closing agent.kick is enough to exit runStateChangeHandler goroutine.
}

func updateAgentStatus(ctx context.Context, q *reform.Querier, agentID string, status inventorypb.AgentStatus, listenPort uint32) error {
	l := logger.Get(ctx)
	l.Debugf("updateAgentStatus: %s %s %d", agentID, status, listenPort)

	agent := &models.Agent{AgentID: agentID}
	err := q.Reload(agent)

	// FIXME that requires more investigation: https://jira.percona.com/browse/PMM-4932
	if err == reform.ErrNoRows {
		l.Warnf("Failed to select Agent by ID for (%s, %s).", agentID, status)

		switch status {
		case inventorypb.AgentStatus_STOPPING, inventorypb.AgentStatus_DONE:
			return nil
		}
	}
	if err != nil {
		return errors.Wrap(err, "failed to select Agent by ID")
	}

	agent.Status = status.String()
	agent.ListenPort = pointer.ToUint16(uint16(listenPort))
	if err = q.Update(agent); err != nil {
		return errors.Wrap(err, "failed to update Agent")
	}
	return nil
}

func (r *Registry) stateChanged(ctx context.Context, req *agentpb.StateChangedRequest) error {
	e := r.db.InTransaction(func(tx *reform.TX) error {
		agentIDs := r.roster.get(req.AgentId)
		if agentIDs == nil {
			agentIDs = []string{req.AgentId}
		}

		for _, agentID := range agentIDs {
			if err := updateAgentStatus(ctx, tx.Querier, agentID, req.Status, req.ListenPort); err != nil {
				return err
			}
		}
		return nil
	})
	if e != nil {
		return e
	}
	r.vmdb.RequestConfigurationUpdate()
	agent, err := models.FindAgentByID(r.db.Querier, req.AgentId)
	if err != nil {
		return err
	}
	if agent.PMMAgentID == nil {
		return nil
	}
	r.RequestStateUpdate(ctx, *agent.PMMAgentID)
	return nil
}

// UpdateAgentsState sends SetStateRequest to all pmm-agents with push metrics agents.
func (r *Registry) UpdateAgentsState(ctx context.Context) error {
	pmmAgents, err := models.FindPMMAgentsIDsWithPushMetrics(r.db.Querier)
	if err != nil {
		return errors.Wrap(err, "cannot find pmmAgentsIDs for AgentsState update")
	}
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 10)
	for _, pmmAgentID := range pmmAgents {
		wg.Add(1)
		limiter <- struct{}{}
		go func(pmmAgentID string) {
			defer wg.Done()
			r.RequestStateUpdate(ctx, pmmAgentID)
			<-limiter
		}(pmmAgentID)
	}
	wg.Wait()
	return nil
}

// SetAllAgentsStatusUnknown goes through all pmm-agents and sets status to UNKNOWN.
func (r *Registry) SetAllAgentsStatusUnknown(ctx context.Context) error {
	agentType := models.PMMAgentType
	agents, err := models.FindAgents(r.db.Querier, models.AgentFilters{AgentType: &agentType})
	if err != nil {
		return errors.Wrap(err, "failed to get pmm-agents")

	}
	for _, agent := range agents {
		if !r.IsConnected(agent.AgentID) {
			err = r.updateAgentStatusForChildren(ctx, agent.AgentID, inventorypb.AgentStatus_UNKNOWN, 0)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Registry) updateAgentStatusForChildren(ctx context.Context, agentID string, status inventorypb.AgentStatus, listenPort uint32) error {
	return r.db.InTransaction(func(t *reform.TX) error {
		agents, err := models.FindAgents(t.Querier, models.AgentFilters{
			PMMAgentID: agentID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get pmm-agent's child agents")
		}
		for _, agent := range agents {
			if err := updateAgentStatus(ctx, t.Querier, agent.AgentID, status, listenPort); err != nil {
				return errors.Wrap(err, "failed to update agent's status")
			}
		}
		return nil
	})
}

// RequestStateUpdate requests state update on pmm-agent with given ID. It sets
// the status to done if the agent is not connected.
func (r *Registry) RequestStateUpdate(ctx context.Context, pmmAgentID string) {
	l := logger.Get(ctx)

	agent, err := r.get(pmmAgentID)
	if err != nil {
		l.Infof("RequestStateUpdate: %s.", err)
		return
	}

	select {
	case agent.stateChangeChan <- struct{}{}:
	default:
	}
}

// sendSetStateRequest sends SetStateRequest to given pmm-agent.
func (r *Registry) sendSetStateRequest(ctx context.Context, agent *pmmAgentInfo) error {
	l := logger.Get(ctx)
	start := time.Now()
	defer func() {
		if dur := time.Since(start); dur > time.Second {
			l.Warnf("sendSetStateRequest took %s.", dur)
		}
	}()
	pmmAgent, err := models.FindAgentByID(r.db.Querier, agent.id)
	if err != nil {
		return errors.Wrap(err, "failed to get PMM Agent")
	}
	pmmAgentVersion, err := version.Parse(*pmmAgent.Version)
	if err != nil {
		return errors.Wrapf(err, "failed to parse PMM agent version %q", *pmmAgent.Version)
	}

	agents, err := models.FindAgents(r.db.Querier, models.AgentFilters{PMMAgentID: agent.id})
	if err != nil {
		return errors.Wrap(err, "failed to collect agents")
	}

	redactMode := redactSecrets
	if l.Logger.GetLevel() >= logrus.DebugLevel {
		redactMode = exposeSecrets
	}

	rdsExporters := make(map[*models.Node]*models.Agent)
	agentProcesses := make(map[string]*agentpb.SetStateRequest_AgentProcess)
	builtinAgents := make(map[string]*agentpb.SetStateRequest_BuiltinAgent)
	for _, row := range agents {
		if row.Disabled {
			continue
		}

		// in order of AgentType consts
		switch row.AgentType {
		case models.PMMAgentType:
			continue
		case models.VMAgentType:
			scrapeCfg, err := r.vmdb.BuildScrapeConfigForVMAgent(agent.id)
			if err != nil {
				return errors.Wrapf(err, "cannot get agent scrape config for agent: %s", agent.id)
			}
			agentProcesses[row.AgentID] = vmAgentConfig(string(scrapeCfg))

		case models.NodeExporterType:
			node, err := models.FindNodeByID(r.db.Querier, pointer.GetString(row.NodeID))
			if err != nil {
				return err
			}
			agentProcesses[row.AgentID] = nodeExporterConfig(node, row)

		case models.RDSExporterType:
			node, err := models.FindNodeByID(r.db.Querier, pointer.GetString(row.NodeID))
			if err != nil {
				return err
			}
			rdsExporters[node] = row
		case models.ExternalExporterType:
			// ignore

		case models.AzureDatabaseExporterType:
			service, err := models.FindServiceByID(r.db.Querier, pointer.GetString(row.ServiceID))
			if err != nil {
				return err
			}
			config, err := azureDatabaseExporterConfig(row, service, redactMode)
			if err != nil {
				return err
			}
			agentProcesses[row.AgentID] = config

		// Agents with exactly one Service
		case models.MySQLdExporterType, models.MongoDBExporterType, models.PostgresExporterType, models.ProxySQLExporterType,
			models.QANMySQLPerfSchemaAgentType, models.QANMySQLSlowlogAgentType, models.QANMongoDBProfilerAgentType, models.QANPostgreSQLPgStatementsAgentType,
			models.QANPostgreSQLPgStatMonitorAgentType:

			service, err := models.FindServiceByID(r.db.Querier, pointer.GetString(row.ServiceID))
			if err != nil {
				return err
			}

			switch row.AgentType {
			case models.MySQLdExporterType:
				agentProcesses[row.AgentID] = mysqldExporterConfig(service, row, redactMode)
			case models.MongoDBExporterType:
				agentProcesses[row.AgentID] = mongodbExporterConfig(service, row, redactMode, pmmAgentVersion)
			case models.PostgresExporterType:
				agentProcesses[row.AgentID] = postgresExporterConfig(service, row, redactMode, pmmAgentVersion)
			case models.ProxySQLExporterType:
				agentProcesses[row.AgentID] = proxysqlExporterConfig(service, row, redactMode)
			case models.QANMySQLPerfSchemaAgentType:
				builtinAgents[row.AgentID] = qanMySQLPerfSchemaAgentConfig(service, row)
			case models.QANMySQLSlowlogAgentType:
				builtinAgents[row.AgentID] = qanMySQLSlowlogAgentConfig(service, row)
			case models.QANMongoDBProfilerAgentType:
				builtinAgents[row.AgentID] = qanMongoDBProfilerAgentConfig(service, row)
			case models.QANPostgreSQLPgStatementsAgentType:
				builtinAgents[row.AgentID] = qanPostgreSQLPgStatementsAgentConfig(service, row)
			case models.QANPostgreSQLPgStatMonitorAgentType:
				builtinAgents[row.AgentID] = qanPostgreSQLPgStatMonitorAgentConfig(service, row)
			}

		default:
			return errors.Errorf("unhandled Agent type %s", row.AgentType)
		}
	}

	if len(rdsExporters) > 0 {
		rdsExporterIDs := make([]string, 0, len(rdsExporters))
		for _, rdsExporter := range rdsExporters {
			rdsExporterIDs = append(rdsExporterIDs, rdsExporter.AgentID)
		}
		sort.Strings(rdsExporterIDs)

		groupID := r.roster.add(agent.id, rdsGroup, rdsExporterIDs)
		c, err := rdsExporterConfig(rdsExporters, redactMode)
		if err != nil {
			return err
		}
		agentProcesses[groupID] = c
	}
	state := &agentpb.SetStateRequest{
		AgentProcesses: agentProcesses,
		BuiltinAgents:  builtinAgents,
	}
	l.Debugf("sendSetStateRequest:\n%s", proto.MarshalTextString(state))
	resp, err := agent.channel.SendAndWaitResponse(state)
	if err != nil {
		return err
	}
	l.Infof("SetState response: %+v.", resp)
	return nil
}

func (r *Registry) isExternalExporterConnectionCheckSupported(q *reform.Querier, pmmAgentID string) (bool, error) {
	pmmAgent, err := models.FindAgentByID(r.db.Querier, pmmAgentID)
	if err != nil {
		return false, fmt.Errorf("failed to get PMM Agent: %s.", err)
	}
	pmmAgentVersion, err := version.Parse(*pmmAgent.Version)
	if err != nil {
		return false, fmt.Errorf("failed to parse PMM agent version %q: %s", *pmmAgent.Version, err)
	}

	if pmmAgentVersion.Less(checkExternalExporterConnectionPMMVersion) {
		return false, nil
	}
	return true, nil
}

// CheckConnectionToService sends request to pmm-agent to check connection to service.
func (r *Registry) CheckConnectionToService(ctx context.Context, q *reform.Querier, service *models.Service, agent *models.Agent) error {
	// TODO: extract to a separate struct to keep Single Responsibility principles: https://jira.percona.com/browse/PMM-4932
	l := logger.Get(ctx)
	start := time.Now()
	defer func() {
		if dur := time.Since(start); dur > 4*time.Second {
			l.Warnf("CheckConnectionToService took %s.", dur)
		}
	}()

	pmmAgentID := pointer.GetString(agent.PMMAgentID)
	if !agent.PushMetrics && (service.ServiceType == models.ExternalServiceType || service.ServiceType == models.HAProxyServiceType) {
		pmmAgentID = models.PMMServerAgentID
	}

	// Skip check connection to external exporter with old pmm-agent.
	if service.ServiceType == models.ExternalServiceType || service.ServiceType == models.HAProxyServiceType {
		isCheckConnSupported, err := r.isExternalExporterConnectionCheckSupported(q, pmmAgentID)
		if err != nil {
			return err
		}

		if !isCheckConnSupported {
			return nil
		}
	}

	pmmAgent, err := r.get(pmmAgentID)
	if err != nil {
		return err
	}

	var request *agentpb.CheckConnectionRequest
	switch service.ServiceType {
	case models.MySQLServiceType:
		tdp := agent.TemplateDelimiters(service)
		request = &agentpb.CheckConnectionRequest{
			Type:    inventorypb.ServiceType_MYSQL_SERVICE,
			Dsn:     agent.DSN(service, 2*time.Second, "", nil),
			Timeout: durationpb.New(3 * time.Second),
			TextFiles: &agentpb.TextFiles{
				Files:              agent.Files(),
				TemplateLeftDelim:  tdp.Left,
				TemplateRightDelim: tdp.Right,
			},
			TlsSkipVerify: agent.TLSSkipVerify,
		}
	case models.PostgreSQLServiceType:
		request = &agentpb.CheckConnectionRequest{
			Type:    inventorypb.ServiceType_POSTGRESQL_SERVICE,
			Dsn:     agent.DSN(service, 2*time.Second, "postgres", nil),
			Timeout: durationpb.New(3 * time.Second),
		}
	case models.MongoDBServiceType:
		tdp := agent.TemplateDelimiters(service)
		request = &agentpb.CheckConnectionRequest{
			Type:    inventorypb.ServiceType_MONGODB_SERVICE,
			Dsn:     agent.DSN(service, 2*time.Second, "", nil),
			Timeout: durationpb.New(3 * time.Second),
			TextFiles: &agentpb.TextFiles{
				Files:              agent.Files(),
				TemplateLeftDelim:  tdp.Left,
				TemplateRightDelim: tdp.Right,
			},
		}
	case models.ProxySQLServiceType:
		request = &agentpb.CheckConnectionRequest{
			Type:    inventorypb.ServiceType_PROXYSQL_SERVICE,
			Dsn:     agent.DSN(service, 2*time.Second, "", nil),
			Timeout: durationpb.New(3 * time.Second),
		}
	case models.ExternalServiceType:
		exporterURL, err := agent.ExporterURL(q)
		if err != nil {
			return err
		}

		request = &agentpb.CheckConnectionRequest{
			Type:    inventorypb.ServiceType_EXTERNAL_SERVICE,
			Dsn:     exporterURL,
			Timeout: durationpb.New(3 * time.Second),
		}
	case models.HAProxyServiceType:
		exporterURL, err := agent.ExporterURL(q)
		if err != nil {
			return err
		}

		request = &agentpb.CheckConnectionRequest{
			Type:    inventorypb.ServiceType_HAPROXY_SERVICE,
			Dsn:     exporterURL,
			Timeout: durationpb.New(3 * time.Second),
		}
	default:
		return errors.Errorf("unhandled Service type %s", service.ServiceType)
	}

	var sanitizedDSN string
	for _, word := range redactWords(agent) {
		sanitizedDSN = strings.ReplaceAll(request.Dsn, word, "****")
	}
	l.Infof("CheckConnectionRequest: type: %s, DSN: %s timeout: %s.", request.Type, sanitizedDSN, request.Timeout)
	resp, err := pmmAgent.channel.SendAndWaitResponse(request)
	if err != nil {
		return err
	}
	l.Infof("CheckConnection response: %+v.", resp)

	switch service.ServiceType {
	case models.MySQLServiceType:
		tableCount := resp.(*agentpb.CheckConnectionResponse).GetStats().GetTableCount()
		agent.TableCount = &tableCount
		l.Debugf("Updating table count: %d.", tableCount)
		if err = q.Update(agent); err != nil {
			return errors.Wrap(err, "failed to update table count")
		}
	case models.ExternalServiceType, models.HAProxyServiceType:
	case models.PostgreSQLServiceType:
	case models.MongoDBServiceType:
	case models.ProxySQLServiceType:
		// nothing yet

	default:
		return errors.Errorf("unhandled Service type %s", service.ServiceType)
	}

	msg := resp.(*agentpb.CheckConnectionResponse).Error
	switch msg {
	case "":
		return nil
	case context.Canceled.Error(), context.DeadlineExceeded.Error():
		msg = fmt.Sprintf("timeout (%s)", msg)
	}
	return status.Error(codes.FailedPrecondition, fmt.Sprintf("Connection check failed: %s.", msg))
}

func (r *Registry) get(pmmAgentID string) (*pmmAgentInfo, error) {
	r.rw.RLock()
	pmmAgent := r.agents[pmmAgentID]
	r.rw.RUnlock()
	if pmmAgent == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "pmm-agent with ID %q is not currently connected", pmmAgentID)
	}
	return pmmAgent, nil
}

// Describe implements prometheus.Collector.
func (r *Registry) Describe(ch chan<- *prom.Desc) {
	r.mAgents.Describe(ch)
}

// Collect implement prometheus.Collector.
func (r *Registry) Collect(ch chan<- prom.Metric) {
	r.rw.RLock()

	for _, agent := range r.agents {
		m := agent.channel.Metrics()

		ch <- prom.MustNewConstMetric(mSentDesc, prom.CounterValue, m.Sent, agent.id)
		ch <- prom.MustNewConstMetric(mRecvDesc, prom.CounterValue, m.Recv, agent.id)
		ch <- prom.MustNewConstMetric(mResponsesDesc, prom.GaugeValue, m.Responses, agent.id)
		ch <- prom.MustNewConstMetric(mRequestsDesc, prom.GaugeValue, m.Requests, agent.id)
	}

	r.rw.RUnlock()

	r.mAgents.Collect(ch)
}

// check interfaces
var (
	_ prom.Collector = (*Registry)(nil)
)
