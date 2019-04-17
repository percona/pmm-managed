package management

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

var (
	errNoParamsNotFound    = status.Error(codes.InvalidArgument, "params not found")
	errOneOfParamsExpected = status.Error(codes.InvalidArgument, "service_id or service_name expected; not both")
)

type ServiceService struct {
	db   *reform.DB
	asrs agentStateRequestSender
}

// NewServiceService creates ServiceService instance.
func NewServiceService(db *reform.DB, asrs agentStateRequestSender) *ServiceService {
	return &ServiceService{
		db:   db,
		asrs: asrs,
	}
}

// RemoveService removes Service with Agents.
func (ss *ServiceService) RemoveService(ctx context.Context, req *managementpb.RemoveServiceRequest) (*managementpb.RemoveServiceResponse, error) {
	err := validateRequest(req)
	if err != nil {
		return nil, err
	}
	pmmAgentIDs := make(map[string]bool)

	if e := ss.db.InTransaction(func(tx *reform.TX) error {
		var service *models.Service
		var err error
		switch {
		case req.ServiceName != "":
			service, err = models.FindServiceByName(ss.db.Querier, req.ServiceName)
		case req.ServiceId != "":
			service, err = models.FindServiceByID(ss.db.Querier, req.ServiceId)
		}
		if err != nil {
			return err
		}

		agents, err := models.AgentsForService(ss.db.Querier, service.ServiceID)
		if err != nil {
			return err
		}
		for _, agent := range agents {
			_, err := models.AgentRemove(ss.db.Querier, agent.AgentID)
			if err != nil {
				return err
			}
			if agent.PMMAgentID != nil {
				pmmAgentIDs[pointer.GetString(agent.PMMAgentID)] = true
			}
		}
		err = models.RemoveService(ss.db.Querier, service.ServiceID)
		if err != nil {
			return err
		}
		return nil
	}); e != nil {
		return nil, e
	}
	for agentID := range pmmAgentIDs {
		ss.asrs.SendSetStateRequest(ctx, agentID)
	}
	return &managementpb.RemoveServiceResponse{}, nil
}

func validateRequest(request *managementpb.RemoveServiceRequest) error {
	if request.ServiceName == "" && request.ServiceId == "" {
		return errNoParamsNotFound
	}
	if request.ServiceName != "" && request.ServiceId != "" {
		return errOneOfParamsExpected
	}
	return nil
}
