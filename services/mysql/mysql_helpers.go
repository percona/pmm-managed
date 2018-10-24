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

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"sort"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	servicelib "github.com/percona/kardianos-service"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

const (
	qanAgentPort uint16 = 9000
)

// ServiceHelper adds, restores and stops MySQLdExporter and QANAgent.
type ServiceHelper struct {
	*ServiceConfig
	pmmServerNode *models.Node
}

// NewServiceHelper creates a new ServiceHelper and checks path to MySQLdExporterPath.
func NewServiceHelper(config *ServiceConfig, pmmServerNode *models.Node) (*ServiceHelper, error) {
	for _, path := range []*string{
		&config.MySQLdExporterPath,
	} {
		if *path == "" {
			continue
		}
		p, err := exec.LookPath(*path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		*path = p
	}

	sh := &ServiceHelper{
		ServiceConfig: config,
		pmmServerNode: pmmServerNode,
	}
	return sh, nil
}

// AgentsList returns a list of agents for service and for node.
func (svc *ServiceHelper) AgentsList(tx *reform.TX, serviceId int32, nodeId int32) ([]models.Agent, []models.Agent, error) {
	agentsForService, err := models.AgentsForServiceID(tx.Querier, serviceId)
	if err != nil {
		return nil, nil, err
	}
	agentsForNode, err := models.AgentsForNodeID(tx.Querier, nodeId)
	if err != nil {
		return nil, nil, err
	}
	return agentsForService, agentsForNode, nil
}

// RemoveAgentsAccosiations removes associations of the service and agents and the node and agents.
func (svc *ServiceHelper) RemoveAgentsAccosiations(tx *reform.TX, serviceId, nodeId int32, agentsForService, agentsForNode []models.Agent) error {
	var err error
	// remove associations of the service and agents
	for _, agent := range agentsForService {
		var deleted uint
		deleted, err = tx.DeleteFrom(models.AgentServiceView, "WHERE service_id = ? AND agent_id = ?", serviceId, agent.ID)
		if err != nil {
			return errors.WithStack(err)
		}
		if deleted != 1 {
			return errors.Errorf("expected to delete 1 record, deleted %d", deleted)
		}
	}
	// remove associations of the node and agents
	for _, agent := range agentsForNode {
		var deleted uint
		deleted, err = tx.DeleteFrom(models.AgentNodeView, "WHERE node_id = ? AND agent_id = ?", nodeId, agent.ID)
		if err != nil {
			return errors.WithStack(err)
		}
		if deleted != 1 {
			return errors.Errorf("expected to delete 1 record, deleted %d", deleted)
		}
	}
	return nil
}

// AddMySQLdExporter starts MySQLdExporter and inserts agent to DB.
func (svc *ServiceHelper) AddMySQLdExporter(ctx context.Context, tx *reform.TX, service *models.MySQLService, username, password string) error {
	// insert mysqld_exporter agent and association
	port, err := svc.PortsRegistry.Reserve()
	if err != nil {
		return err
	}
	agent := &models.MySQLdExporter{
		Type:         models.MySQLdExporterAgentType,
		RunsOnNodeID: svc.pmmServerNode.ID,

		ServiceUsername: &username,
		ServicePassword: &password,
		ListenPort:      &port,
	}
	if err = tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// check connection and a number of tables
	var tableCount int
	dsn := agent.DSN(service)
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables").Scan(&tableCount)
		db.Close()
		agent.MySQLDisableTablestats = pointer.ToBool(tableCount > 1000)
	}
	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok {
			switch err.Number {
			case 0x414: // 1044
				return status.Error(codes.PermissionDenied, err.Message)
			case 0x415: // 1045
				return status.Error(codes.Unauthenticated, err.Message)
			}
		}
		return errors.WithStack(err)
	}

	// start mysqld_exporter agent
	if svc.MySQLdExporterPath != "" {
		cfg := svc.MysqlExporterCfg(agent, port, dsn)
		if err = svc.Supervisor.Start(ctx, cfg); err != nil {
			return err
		}
	}

	return nil
}

// MysqlExporterCfg returns configs for MySQLdExporter.
func (svc *ServiceHelper) MysqlExporterCfg(agent *models.MySQLdExporter, port uint16, dsn string) *servicelib.Config {
	name := agent.NameForSupervisor()

	arguments := []string{
		"-collect.binlog_size",
		"-collect.global_status",
		"-collect.global_variables",
		"-collect.info_schema.innodb_metrics",
		"-collect.info_schema.processlist",
		"-collect.info_schema.query_response_time",
		"-collect.info_schema.userstats",
		"-collect.perf_schema.eventswaits",
		"-collect.perf_schema.file_events",
		"-collect.slave_status",
		fmt.Sprintf("-web.listen-address=127.0.0.1:%d", port),
	}
	if agent.MySQLDisableTablestats == nil || !*agent.MySQLDisableTablestats {
		// enable tablestats and a few related collectors just like pmm-admin
		// https://github.com/percona/pmm-client/blob/e94b61ed0e5482a27039f0d1b0b36076731f0c29/pmm/plugin/mysql/metrics/metrics.go#L98-L105
		arguments = append(arguments, "-collect.auto_increment.columns")
		arguments = append(arguments, "-collect.info_schema.tables")
		arguments = append(arguments, "-collect.info_schema.tablestats")
		arguments = append(arguments, "-collect.perf_schema.indexiowaits")
		arguments = append(arguments, "-collect.perf_schema.tableiowaits")
		arguments = append(arguments, "-collect.perf_schema.tablelocks")
	}
	sort.Strings(arguments)

	return &servicelib.Config{
		Name:        name,
		DisplayName: name,
		Description: name,
		Executable:  svc.MySQLdExporterPath,
		Arguments:   arguments,
		Environment: []string{fmt.Sprintf("DATA_SOURCE_NAME=%s", dsn)},
	}
}

// RestoreMySQLdExporter restores MySQLdExporter.
func (svc *ServiceHelper) RestoreMySQLdExporter(ctx context.Context, tx *reform.TX, agent models.Agent, service *models.MySQLService) error {
	a := &models.MySQLdExporter{ID: agent.ID}
	if err := tx.Reload(a); err != nil {
		return errors.WithStack(err)
	}
	dsn := a.DSN(service)
	port := *a.ListenPort
	if svc.MySQLdExporterPath != "" {
		name := a.NameForSupervisor()

		// Checks if init script already running.
		err := svc.Supervisor.Status(ctx, name)
		if err == nil {
			// Stops init script.
			if err = svc.Supervisor.Stop(ctx, name); err != nil {
				return err
			}
		}

		// Installs new version of the script.
		cfg := svc.MysqlExporterCfg(a, port, dsn)
		if err = svc.Supervisor.Start(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

// StopMySQLdExporter stops MySQLdExporter process.
func (svc *ServiceHelper) StopMySQLdExporter(ctx context.Context, tx *reform.TX, agent models.Agent) error {
	a := models.MySQLdExporter{ID: agent.ID}
	if err := tx.Reload(&a); err != nil {
		return errors.WithStack(err)
	}
	if svc.MySQLdExporterPath != "" {
		if err := svc.Supervisor.Stop(ctx, a.NameForSupervisor()); err != nil {
			return err
		}
	}
	return nil
}

// AddQanAgent adds new MySQL service to QAN and inserts agent to DB.
func (svc *ServiceHelper) AddQanAgent(ctx context.Context, tx *reform.TX, service *models.MySQLService, nodeName string, username, password string) error {
	// Despite running a single qan-agent process on PMM Server, we use one database record per MySQL instance
	// to store username/password and UUID.

	// insert qan-agent agent and association
	agent := &models.QanAgent{
		Type:         models.QanAgentAgentType,
		RunsOnNodeID: svc.pmmServerNode.ID,

		ServiceUsername: &username,
		ServicePassword: &password,
		ListenPort:      pointer.ToUint16(qanAgentPort),
	}
	var err error
	if err = tx.Insert(agent); err != nil {
		return errors.WithStack(err)
	}
	if err = tx.Insert(&models.AgentService{AgentID: agent.ID, ServiceID: service.ID}); err != nil {
		return errors.WithStack(err)
	}

	// DSNs for mysqld_exporter and qan-agent are currently identical,
	// so we do not check connection again

	// start or reconfigure qan-agent
	if svc.QAN != nil {
		if err = svc.QAN.AddMySQL(ctx, nodeName, service, agent); err != nil {
			return err
		}

		// re-save agent with set QANDBInstanceUUID
		if err = tx.Save(agent); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// StopQanAgent removes a MySQL service from QAN.
func (svc *ServiceHelper) StopQanAgent(ctx context.Context, tx *reform.TX, agent models.Agent) error {
	a := models.QanAgent{ID: agent.ID}
	if err := tx.Reload(&a); err != nil {
		return errors.WithStack(err)
	}
	if svc.QAN != nil {
		if err := svc.QAN.RemoveMySQL(ctx, &a); err != nil {
			return err
		}
	}
	return nil
}

// RestoreQanAgent ensures QAN agent runs.
func (svc *ServiceHelper) RestoreQanAgent(ctx context.Context, tx *reform.TX, agent models.Agent) error {
	a := models.QanAgent{ID: agent.ID}
	if err := tx.Reload(&a); err != nil {
		return errors.WithStack(err)
	}
	if svc.QAN != nil {
		name := a.NameForSupervisor()

		// Checks if init script already running.
		err := svc.Supervisor.Status(ctx, name)
		if err == nil {
			// Stops init script.
			if err = svc.Supervisor.Stop(ctx, name); err != nil {
				return err
			}
		}

		// Installs new version of the script.
		if err = svc.QAN.EnsureAgentRuns(ctx, name, *a.ListenPort); err != nil {
			return err
		}
	}
	return nil
}
