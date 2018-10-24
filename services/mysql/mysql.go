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

// Package mysql contains business logic of working with Remote MySQL instances.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services"
	"github.com/percona/pmm-managed/services/prometheus"
	"github.com/percona/pmm-managed/services/qan"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/ports"
)

const (
	defaultMySQLPort uint32 = 3306
)

type ServiceConfig struct {
	MySQLdExporterPath string

	Prometheus    *prometheus.Service
	Supervisor    services.Supervisor
	DB            *reform.DB
	PortsRegistry *ports.Registry
	QAN           *qan.Service
}

// Service is responsible for interactions with AWS RDS.
type Service struct {
	*ServiceHelper
}

// NewService creates a new service.
func NewService(config *ServiceConfig) (*Service, error) {
	var node models.Node
	err := config.DB.FindOneTo(&node, "type", models.PMMServerNodeType)
	if err != nil {
		return nil, err
	}

	serviceHelper, err := NewServiceHelper(config, &node)
	if err != nil {
		return nil, err
	}
	svc := &Service{
		ServiceHelper: serviceHelper,
	}
	return svc, nil
}

type Instance struct {
	Node    models.RemoteNode
	Service models.MySQLService
}

func (svc *Service) ApplyPrometheusConfiguration(ctx context.Context, q *reform.Querier) error {
	mySQLHR := &prometheus.ScrapeConfig{
		JobName:        "remote-mysql-hr",
		ScrapeInterval: "1s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/metrics-hr",
		HonorLabels:    true,
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mysql",
		}},
	}
	mySQLMR := &prometheus.ScrapeConfig{
		JobName:        "remote-mysql-mr",
		ScrapeInterval: "5s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/metrics-mr",
		HonorLabels:    true,
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mysql",
		}},
	}
	mySQLLR := &prometheus.ScrapeConfig{
		JobName:        "remote-mysql-lr",
		ScrapeInterval: "60s",
		ScrapeTimeout:  "5s",
		MetricsPath:    "/metrics-lr",
		HonorLabels:    true,
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "mysql",
		}},
	}

	nodes, err := q.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		var service models.MySQLService
		if e := q.SelectOneTo(&service, "WHERE node_id = ? and type = ?", node.ID, models.MySQLServiceType); e != nil {
			return errors.WithStack(e)
		}

		agents, err := models.AgentsForServiceID(q, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MySQLdExporterAgentType:
				a := models.MySQLdExporter{ID: agent.ID}
				if e := q.Reload(&a); e != nil {
					return errors.WithStack(e)
				}
				logger.Get(ctx).WithField("component", "mysql").Infof("%s %s %s %d", a.Type, node.Name, node.Region, *a.ListenPort)

				sc := prometheus.StaticConfig{
					Targets: []string{fmt.Sprintf("127.0.0.1:%d", *a.ListenPort)},
					Labels: []prometheus.LabelPair{
						{Name: "instance", Value: node.Name},
					},
				}
				mySQLHR.StaticConfigs = append(mySQLHR.StaticConfigs, sc)
				mySQLMR.StaticConfigs = append(mySQLMR.StaticConfigs, sc)
				mySQLLR.StaticConfigs = append(mySQLLR.StaticConfigs, sc)
			}
		}
	}

	// sort by instance
	sorterFor := func(sc []prometheus.StaticConfig) func(int, int) bool {
		return func(i, j int) bool {
			return sc[i].Labels[0].Value < sc[j].Labels[0].Value
		}
	}
	sort.Slice(mySQLHR.StaticConfigs, sorterFor(mySQLHR.StaticConfigs))
	sort.Slice(mySQLMR.StaticConfigs, sorterFor(mySQLMR.StaticConfigs))
	sort.Slice(mySQLLR.StaticConfigs, sorterFor(mySQLLR.StaticConfigs))

	return svc.Prometheus.SetScrapeConfigs(ctx, false, mySQLHR, mySQLMR, mySQLLR)
}

func (svc *Service) List(ctx context.Context) ([]Instance, error) {
	var res []Instance
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		structs, e := tx.SelectAllFrom(models.RemoteNodeTable, "WHERE type = ? ORDER BY id", models.RemoteNodeType)
		if e != nil {
			return e
		}
		nodes := make([]models.RemoteNode, len(structs))
		for i, str := range structs {
			nodes[i] = *str.(*models.RemoteNode)
		}

		structs, e = tx.SelectAllFrom(models.MySQLServiceTable, "WHERE type = ? ORDER BY id", models.MySQLServiceType)
		if e != nil {
			return e
		}
		services := make([]models.MySQLService, len(structs))
		for i, str := range structs {
			services[i] = *str.(*models.MySQLService)
		}

		for _, node := range nodes {
			for _, service := range services {
				if node.ID == service.NodeID {
					res = append(res, Instance{
						Node:    node,
						Service: service,
					})
				}
			}
		}
		return nil
	})
	return res, err
}

func (svc *Service) Add(ctx context.Context, name, address string, port uint32, username, password string) (int32, error) {
	if address == "" {
		return 0, status.Error(codes.InvalidArgument, "MySQL instance host is not given.")
	}
	if username == "" {
		return 0, status.Error(codes.InvalidArgument, "Username is not given.")
	}
	if port == 0 {
		port = defaultMySQLPort
	}
	if name == "" {
		name = address
	}

	var id int32
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		// insert node
		node := &models.RemoteNode{
			Type:   models.RemoteNodeType,
			Name:   name,
			Region: models.RemoteNodeRegion,
		}
		if err := tx.Insert(node); err != nil {
			if err, ok := err.(*mysql.MySQLError); ok && err.Number == 0x426 {
				return status.Errorf(codes.AlreadyExists, "MySQL instance %q already exists.",
					node.Name)
			}
			return errors.WithStack(err)
		}
		id = node.ID

		engine, engineVersion, err := svc.engineAndEngineVersion(ctx, address, port, username, password)
		if err != nil {
			return errors.WithStack(err)
		}

		// insert service
		service := &models.MySQLService{
			Type:   models.MySQLServiceType,
			NodeID: node.ID,

			Address:       &address,
			Port:          pointer.ToUint16(uint16(port)),
			Engine:        &engine,
			EngineVersion: &engineVersion,
		}
		if err = tx.Insert(service); err != nil {
			return errors.WithStack(err)
		}

		if err = svc.AddMySQLdExporter(ctx, tx, service, username, password); err != nil {
			return err
		}
		if err = svc.AddQanAgent(ctx, tx, service, node.Name, username, password); err != nil {
			return err
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})

	return id, err
}

func (svc *Service) Remove(ctx context.Context, id int32) error {
	var err error
	return svc.DB.InTransaction(func(tx *reform.TX) error {
		var node models.RemoteNode
		if err = tx.SelectOneTo(&node, "WHERE type = ? AND id = ?", models.RemoteNodeType, id); err != nil {
			if err == reform.ErrNoRows {
				return status.Errorf(codes.NotFound, "MySQL instance with ID %d not found.", id)
			}
			return errors.WithStack(err)
		}

		var service models.MySQLService
		if err = tx.SelectOneTo(&service, "WHERE node_id = ? and type = ?", node.ID, models.MySQLServiceType); err != nil {
			return errors.WithStack(err)
		}

		agentsForService, agentsForNode, err := svc.AgentsList(tx, service.ID, node.ID)
		if err != nil {
			return err
		}

		err = svc.RemoveAgentsAccosiations(tx, service.ID, node.ID, agentsForService, agentsForNode)
		if err != nil {
			return err
		}

		// stop agents
		agents := make(map[int32]models.Agent, len(agentsForService)+len(agentsForNode))
		for _, agent := range agentsForService {
			agents[agent.ID] = agent
		}
		for _, agent := range agentsForNode {
			agents[agent.ID] = agent
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MySQLdExporterAgentType:
				err = svc.StopMySQLExporter(ctx, tx, agent)
				if err != nil {
					return err
				}

			case models.QanAgentAgentType:
				err = svc.StopQanAgent(ctx, tx, agent)
				if err != nil {
					return err
				}
			}
		}

		// remove agents
		for _, agent := range agents {
			if err = tx.Delete(&agent); err != nil {
				return errors.WithStack(err)
			}
		}

		if err = tx.Delete(&service); err != nil {
			return errors.WithStack(err)
		}
		if err = tx.Delete(&node); err != nil {
			return errors.WithStack(err)
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})
}

// Restore configuration from database.
func (svc *Service) Restore(ctx context.Context, tx *reform.TX) error {
	nodes, err := tx.FindAllFrom(models.RemoteNodeTable, "type", models.RemoteNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.RemoteNode)

		service := &models.MySQLService{}
		if e := tx.SelectOneTo(service, "WHERE node_id = ?", node.ID); e != nil {
			return errors.WithStack(e)
		}

		agents, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.MySQLdExporterAgentType:
				err = svc.RestoreMySQLdExporter(ctx, tx, agent, service)
				if err != nil {
					return err
				}

			case models.QanAgentAgentType:
				err = svc.RestoreQanAgent(ctx, tx, agent)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (svc *Service) engineAndEngineVersion(ctx context.Context, host string, port uint32, username string, password string) (string, string, error) {
	var version string
	var versionComment string
	agent := models.MySQLdExporter{
		ServiceUsername: pointer.ToString(username),
		ServicePassword: pointer.ToString(password),
	}
	service := &models.MySQLService{
		Address: &host,
		Port:    pointer.ToUint16(uint16(port)),
	}
	dsn := agent.DSN(service)
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		err = db.QueryRowContext(ctx, "select @@version;").Scan(&version)
		err = db.QueryRowContext(ctx, "select @@version_comment;").Scan(&versionComment)
		db.Close()
	}
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	return versionComment, version, nil
}
