package management

import (
	"context"

	"github.com/AlekSi/pointer"
	inventorypb "github.com/percona/pmm/api/inventory"
	"github.com/percona/pmm/api/managementpb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/services/inventory"
)

// MySQLService MySQL Management Service
type MySQLService struct {
	db          *reform.DB
	servicesSvc *inventory.ServicesService
	agentsSvc   *inventory.AgentsService
}

// NewMySQLService creates new MySQL Management Service
func NewMySQLService(db *reform.DB, s *inventory.ServicesService, a *inventory.AgentsService) *MySQLService {
	return &MySQLService{db, s, a}
}

// Add adds "MySQL Service", "MySQL Exporter Agent" and "QAN MySQL PerfSchema Agent"
func (s *MySQLService) Add(ctx context.Context, req *managementpb.AddMySQLRequest) (res *managementpb.AddMySQLResponse, err error) {
	res = &managementpb.AddMySQLResponse{}

	if e := s.db.InTransaction(func(tx *reform.TX) error {
		address := pointer.ToStringOrNil(req.Address)
		port := pointer.ToUint16OrNil(uint16(req.Port))
		service, err := s.servicesSvc.AddMySQL(ctx, req.ServiceName, req.NodeId, address, port, tx.Querier)
		if err != nil {
			return err
		}
		res.Service = service

		if req.MysqldExporter {
			request := &inventorypb.AddMySQLdExporterRequest{
				PmmAgentId: req.PmmAgentId,
				ServiceId:  service.ID(),
				Username:   req.Username,
				Password:   req.Password,
			}

			agent, err := s.agentsSvc.AddMySQLdExporter(ctx, request, tx.Querier)
			if err != nil {
				return err
			}

			res.MysqldExporter = agent
		}

		if req.QanMysqlPerfschema {
			request := &inventorypb.AddQANMySQLPerfSchemaAgentRequest{
				PmmAgentId: req.PmmAgentId,
				ServiceId:  service.ID(),
				Username:   req.QanUsername,
				Password:   req.QanPassword,
			}

			qAgent, err := s.agentsSvc.AddQANMySQLPerfSchemaAgent(ctx, request, tx.Querier)
			if err != nil {
				return err
			}

			res.QanMysqlPerfschema = qAgent
		}

		return nil
	}); e != nil {
		return nil, e
	}

	return res, nil
}
