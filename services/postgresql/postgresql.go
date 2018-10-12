package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	servicelib "github.com/percona/kardianos-service"
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services"
	"github.com/percona/pmm-managed/services/prometheus"
	"github.com/percona/pmm-managed/services/qan"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/ports"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
)

type ServiceConfig struct {
	PostgreSQLExporterPath string

	DB            *reform.DB
	Prometheus    *prometheus.Service
	QAN           *qan.Service
	Supervisor    services.Supervisor
	PortsRegistry *ports.Registry
}

// Service is responsible for interactions with PostgreSQL.
type Service struct {
	*ServiceConfig
	httpClient    *http.Client
	pmmServerNode *models.Node
}

// NewService creates a new service.
func NewService(config *ServiceConfig) (*Service, error) {
	var node models.Node
	err := config.DB.FindOneTo(&node, "type", models.PMMServerNodeType)
	if err != nil {
		return nil, err
	}

	for _, path := range []*string{
		&config.PostgreSQLExporterPath,
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

	svc := &Service{
		ServiceConfig: config,
		httpClient:    new(http.Client),
		pmmServerNode: &node,
	}
	return svc, nil
}

func (svc *Service) ApplyPrometheusConfiguration(ctx context.Context, q *reform.Querier) error {
	postgreSQLConfig := &prometheus.ScrapeConfig{
		JobName:        "postgres",
		ScrapeInterval: "1s",
		ScrapeTimeout:  "1s",
		MetricsPath:    "/metrics",
		HonorLabels:    true,
		RelabelConfigs: []prometheus.RelabelConfig{{
			TargetLabel: "job",
			Replacement: "postgres",
		}},
	}

	nodes, err := q.FindAllFrom(models.PostgreSQLNodeTable, "type", models.PostgreSQLNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.PostgreSQLNode)

		var service models.PostgreSQLService
		if e := q.SelectOneTo(&service, "WHERE node_id = ?", node.ID); e != nil {
			return errors.WithStack(e)
		}

		agents, err := models.AgentsForServiceID(q, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.PostgreSQLExporterAgentType:
				a := models.PostgreSQLExporter{ID: agent.ID}
				if e := q.Reload(&a); e != nil {
					return errors.WithStack(e)
				}
				logger.Get(ctx).WithField("component", "postgresql").Infof("%s %s %d", a.Type, node.Name, *a.ListenPort)

				sc := prometheus.StaticConfig{
					Targets: []string{fmt.Sprintf("127.0.0.1:%d", *a.ListenPort)},
					Labels: []prometheus.LabelPair{
						{Name: "instance", Value: node.Name},
					},
				}
				postgreSQLConfig.StaticConfigs = append(postgreSQLConfig.StaticConfigs, sc)
			}
		}
	}

	// sort by region and name
	sorterFor := func(sc []prometheus.StaticConfig) func(int, int) bool {
		return func(i, j int) bool {
			if sc[i].Labels[0].Value != sc[j].Labels[0].Value {
				return sc[i].Labels[0].Value < sc[j].Labels[0].Value
			}
			return sc[i].Labels[1].Value < sc[j].Labels[1].Value
		}
	}
	sort.Slice(postgreSQLConfig.StaticConfigs, sorterFor(postgreSQLConfig.StaticConfigs))

	return svc.Prometheus.SetScrapeConfigs(ctx, false, postgreSQLConfig)
}

type Instance struct {
	Node    models.PostgreSQLNode
	Service models.PostgreSQLService
}

func (svc *Service) List(ctx context.Context) ([]Instance, error) {
	res := []Instance{}
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		structs, e := tx.SelectAllFrom(models.PostgreSQLNodeTable, "WHERE type = ? ORDER BY id", models.PostgreSQLNodeType)
		if e != nil {
			return e
		}
		nodes := make([]models.PostgreSQLNode, len(structs))
		for i, str := range structs {
			nodes[i] = *str.(*models.PostgreSQLNode)
		}

		structs, e = tx.SelectAllFrom(models.PostgreSQLServiceTable, "WHERE type = ? ORDER BY id", models.PostgreSQLServiceType)
		if e != nil {
			return e
		}
		services := make([]models.PostgreSQLService, len(structs))
		for i, str := range structs {
			services[i] = *str.(*models.PostgreSQLService)
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
		return 0, status.Error(codes.InvalidArgument, "PostgreSQL instance host is not given.")
	}
	if port == 0 {
		return 0, status.Error(codes.InvalidArgument, "PostgreSQL instance port is not given.")
	}
	if username == "" {
		return 0, status.Error(codes.InvalidArgument, "Username is not given.")
	}
	if name == "" {
		name = fmt.Sprintf("%s:%d", address, port)
	}

	var id int32
	err := svc.DB.InTransaction(func(tx *reform.TX) error {
		// insert node
		node := &models.PostgreSQLNode{
			Type: models.PostgreSQLNodeType,
			Name: name,
		}
		if err := tx.Insert(node); err != nil {
			if err, ok := err.(*mysql.MySQLError); ok && err.Number == 0x426 {
				return status.Errorf(codes.AlreadyExists, "PostgreSQL instance %q already exists",
					node.Name)
			}
			return errors.WithStack(err)
		}
		//		cocroachDB := `CockroachDB CCL v2.0.6 (x86_64-unknown-linux-gnu, built 2018/10/01       |
		//| 13:59:40, go1.10)`
		//		postgresql := `PostgreSQL 10.5 (Ubuntu 10.5-0ubuntu0.18.04) on x86_64-pc-linux-gnu, compiled by gcc (Ubuntu 7.3.0-16ubuntu3) 7.3.0, 64-bit`

		engine, engineVersion := svc.engineAndEngineVersion(address, port, username, password, ctx)

		// insert service
		service := &models.PostgreSQLService{
			Type:   models.PostgreSQLServiceType,
			NodeID: node.ID,

			Address:       &address,
			Port:          pointer.ToUint16(uint16(port)),
			Engine:        &engine,
			EngineVersion: &engineVersion,
		}
		if err := tx.Insert(service); err != nil {
			return errors.WithStack(err)
		}
		id = service.PKValue().(int32)

		if err := svc.addPostgreSQLExporter(ctx, tx, service, username, password); err != nil {
			return err
		}

		return svc.ApplyPrometheusConfiguration(ctx, tx.Querier)
	})

	return id, err
}

func (svc *Service) engineAndEngineVersion(host string, port uint32, username string, password string, ctx context.Context) (string, string) {
	var databaseVersion string
	address := net.JoinHostPort(host, strconv.Itoa(int(port)))
	dsn := fmt.Sprintf(`postgres://%s:%s@%s`, username, password, address)
	db, err := sql.Open("postgres", dsn)
	if err == nil {
		err = db.QueryRowContext(ctx, "SELECT Version();").Scan(&databaseVersion)
		db.Close()
	}
	var engine string
	var engineVersion string
	cocroachDBRegexp := regexp.MustCompile(`CockroachDB CCL (v[\d\.]]+)`)
	postgresDBRegexp := regexp.MustCompile(`PostgreSQL ([\d\.]]+)`)
	if cocroachDBRegexp.MatchString(databaseVersion) {
		engine = "CockroachDB"
		submatch := cocroachDBRegexp.FindStringSubmatch(databaseVersion)
		engineVersion = submatch[1]
	} else if postgresDBRegexp.MatchString(databaseVersion) {
		engine = "PostgreSQL"
		submatch := postgresDBRegexp.FindStringSubmatch(databaseVersion)
		engineVersion = submatch[1]
	}
	return engine, engineVersion
}

func (svc *Service) Remove(ctx context.Context, id int32) error {
	var err error
	return svc.DB.InTransaction(func(tx *reform.TX) error {
		var node models.PostgreSQLNode
		if err = tx.SelectOneTo(&node, "WHERE type = ? AND id = ?", models.PostgreSQLNodeType, id); err != nil {
			if err == reform.ErrNoRows {
				return status.Errorf(codes.NotFound, "PostgreSQL instance %q not found.", id)
			}
			return errors.WithStack(err)
		}

		var service models.PostgreSQLService
		if err = tx.SelectOneTo(&service, "WHERE node_id = ?", node.ID); err != nil {
			return errors.WithStack(err)
		}

		// remove associations of the service and agents
		agentsForService, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agentsForService {
			var deleted uint
			deleted, err = tx.DeleteFrom(models.AgentServiceView, "WHERE service_id = ? AND agent_id = ?", service.ID, agent.ID)
			if err != nil {
				return errors.WithStack(err)
			}
			if deleted != 1 {
				return errors.Errorf("expected to delete 1 record, deleted %d", deleted)
			}
		}

		// remove associations of the node and agents
		agentsForNode, err := models.AgentsForNodeID(tx.Querier, node.ID)
		if err != nil {
			return err
		}
		for _, agent := range agentsForNode {
			var deleted uint
			deleted, err = tx.DeleteFrom(models.AgentNodeView, "WHERE node_id = ? AND agent_id = ?", node.ID, agent.ID)
			if err != nil {
				return errors.WithStack(err)
			}
			if deleted != 1 {
				return errors.Errorf("expected to delete 1 record, deleted %d", deleted)
			}
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
			case models.PostgreSQLExporterAgentType:
				a := models.MySQLdExporter{ID: agent.ID}
				if err = tx.Reload(&a); err != nil {
					return errors.WithStack(err)
				}
				if svc.PostgreSQLExporterPath != "" {
					if err = svc.Supervisor.Stop(ctx, a.NameForSupervisor()); err != nil {
						return err
					}
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

func (svc *Service) addPostgreSQLExporter(ctx context.Context, tx *reform.TX, service *models.PostgreSQLService, username, password string) error {
	// insert postgres_exporter agent and association
	port, err := svc.PortsRegistry.Reserve()
	if err != nil {
		return err
	}
	agent := &models.PostgreSQLExporter{
		Type:         models.PostgreSQLExporterAgentType,
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
	db, err := sql.Open("postgres", dsn)
	if err == nil {
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables").Scan(&tableCount)
		db.Close()
	}
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {
			case "42501":
				return status.Error(codes.PermissionDenied, err.Message)
			case "28P01":
				return status.Error(codes.Unauthenticated, err.Message)
			}
		}
		return errors.WithStack(err)
	}

	// start mysqld_exporter agent
	if svc.PostgreSQLExporterPath != "" {
		cfg := svc.postgresExporterCfg(agent, port, dsn)
		if err = svc.Supervisor.Start(ctx, cfg); err != nil {
			return err
		}
	}

	return nil
}

// Restore configuration from database.
func (svc *Service) Restore(ctx context.Context, tx *reform.TX) error {
	nodes, err := tx.FindAllFrom(models.PostgreSQLNodeTable, "type", models.PostgreSQLNodeType)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range nodes {
		node := n.(*models.PostgreSQLNode)

		service := &models.PostgreSQLService{}
		if e := tx.SelectOneTo(service, "WHERE node_id = ?", node.ID); e != nil {
			return errors.WithStack(e)
		}

		agents, err := models.AgentsForServiceID(tx.Querier, service.ID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			switch agent.Type {
			case models.PostgreSQLExporterAgentType:
				a := &models.PostgreSQLExporter{ID: agent.ID}
				if err = tx.Reload(a); err != nil {
					return errors.WithStack(err)
				}
				dsn := a.DSN(service)
				port := *a.ListenPort
				if svc.PostgreSQLExporterPath != "" {
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
					cfg := svc.postgresExporterCfg(a, port, dsn)
					if err = svc.Supervisor.Start(ctx, cfg); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (svc *Service) postgresExporterCfg(agent *models.PostgreSQLExporter, port uint16, dsn string) *servicelib.Config {
	name := agent.NameForSupervisor()

	arguments := []string{
		fmt.Sprintf("-web.listen-address=127.0.0.1:%d", port),
		//TODO: set arguments
	}
	sort.Strings(arguments)

	return &servicelib.Config{
		Name:        name,
		DisplayName: name,
		Description: name,
		Executable:  svc.PostgreSQLExporterPath,
		Arguments:   arguments,
		Environment: []string{fmt.Sprintf("DATA_SOURCE_NAME=%s", dsn)},
	}
}
