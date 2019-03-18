package management

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm-managed/services/inventory"
	inventorypb "github.com/percona/pmm/api/inventory"
	"github.com/percona/pmm/api/managementpb"
	"gopkg.in/reform.v1"
)

type MySQLService struct {
	nodesSvc    *inventory.NodesService
	servicesSvc *inventory.ServicesService
	agentsSvc   *inventory.AgentsService
}

func NewMySQLService(n *inventory.NodesService, s *inventory.ServicesService, a *inventory.AgentsService) *MySQLService {
	return &MySQLService{n, s, a}
}

func (s *MySQLService) Add(ctx context.Context, req *managementpb.AddMySQLRequest) error {
	db := &reform.DB{}

	node, err := s.nodesSvc.Get(ctx, req.NodeId)
	if err != nil {
		return err // TODO: Node not found error
	}

	nodeAgents, err := s.agentsSvc.List(ctx, db, inventory.AgentFilters{NodeID: node.ID()})
	// TODO: PMMAgent Not Found Error

	address := pointer.ToStringOrNil(req.Address)
	port := pointer.ToUint16OrNil(uint16(req.Port))
	svc, err := s.servicesSvc.AddMySQL(ctx, req.ServiceName, node.ID(), address, port)
	if err != nil {
		return err // TODO: Can't add service error
	}

	// Only if "mysqld_exporter" flag provided
	if req.MysqldExporter {
		_, err = s.agentsSvc.AddMySQLdExporter(ctx, db, &inventorypb.AddMySQLdExporterRequest{
			PmmAgentId: req.ServiceName,
			ServiceId:  svc.ID(),
			Username:   req.Username,
			Password:   req.Password,
		})

		if err != nil {
			return err // TODO: Can't add exporter error
		}
	}

	return nil
}
