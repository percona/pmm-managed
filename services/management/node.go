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

package management

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"

	// FIXME Refactor, as service shouldn't depend on other service in one abstraction level.
	// https://jira.percona.com/browse/PMM-3541
	// See also main_test.go
	"github.com/percona/pmm-managed/services/inventory"
)

// NodeService represents service for working with nodes.
type NodeService struct {
	db       *reform.DB
	registry registry
}

// NewNodeService creates NodeService instance.
func NewNodeService(db *reform.DB, registry registry) *NodeService {
	return &NodeService{
		db:       db,
		registry: registry,
	}
}

func (s *NodeService) Register(ctx context.Context, req *managementpb.RegisterNodeRequest) (*managementpb.RegisterNodeResponse, error) {
	var res managementpb.RegisterNodeResponse
	if e := s.db.InTransaction(func(tx *reform.TX) error {
		node, err := createOrUpdateNode(req, tx.Querier)
		if err != nil {
			return err
		}

		pmmAgent, err := s.findPmmAgentByNodeID(tx.Querier, node.NodeID)
		switch err {
		case errAgentNotFound:
			pmmAgent, err = models.AgentAddPmmAgent(tx.Querier, node.NodeID, nil)
			if err != nil {
				return err
			}
		case nil:
			// noop
		default:
			return err
		}

		if err := s.addPmmAgentToResponse(tx.Querier, pmmAgent, &res); err != nil {
			return err
		}

		_, err = s.findNodeExporterByPmmAgentID(tx.Querier, pmmAgent.AgentID)
		switch err {
		case errAgentNotFound:
			_, err := models.AgentAddNodeExporter(tx.Querier, pmmAgent.AgentID, nil)
			if err != nil {
				return err
			}
		case nil:
			// noop
		default:
			return err
		}

		n, err := inventory.ToInventoryNode(node)
		if err != nil {
			return err
		}
		switch n := n.(type) {
		case *inventorypb.GenericNode:
			res.GenericNode = n
		case *inventorypb.ContainerNode:
			res.ContainerNode = n
		}

		return nil
	}); e != nil {
		return nil, e
	}

	s.registry.SendSetStateRequest(ctx, res.PmmAgent.AgentId)

	return &res, nil
}

func createOrUpdateNode(req *managementpb.RegisterNodeRequest, q *reform.Querier) (*models.Node, error) {
	node, err := models.FindNodeByName(q, req.NodeName)
	switch status.Code(err) {
	case codes.OK:
		var nodeType inventorypb.NodeType
		switch node.NodeType {
		case models.GenericNodeType:
			nodeType = inventorypb.NodeType_GENERIC_NODE
		case models.ContainerNodeType:
			nodeType = inventorypb.NodeType_CONTAINER_NODE
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unsupported Node type %q", req.NodeType)
		}

		if nodeType != req.NodeType {
			return nil, status.Errorf(codes.InvalidArgument, "unexpected Node type %q", req.NodeType)
		}

		params := &models.UpdateNodeParams{
			Address:      req.Address,
			MachineID:    req.MachineId,
			CustomLabels: req.CustomLabels,
			// TODO distro, node_model, region, az, container_id, container_name
		}
		return models.UpdateNode(q, node.NodeID, params)

	case codes.NotFound:
		var nodeType models.NodeType
		switch req.NodeType {
		case inventorypb.NodeType_GENERIC_NODE:
			nodeType = models.GenericNodeType
		case inventorypb.NodeType_CONTAINER_NODE:
			nodeType = models.ContainerNodeType
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unsupported Node type %q", req.NodeType)
		}

		params := &models.CreateNodeParams{
			NodeName:      req.NodeName,
			MachineID:     pointer.ToStringOrNil(req.MachineId),
			Distro:        req.Distro,
			NodeModel:     req.NodeModel, // TODO
			AZ:            req.Az,        // TODO
			ContainerID:   pointer.ToStringOrNil(req.ContainerId),
			ContainerName: pointer.ToStringOrNil(req.ContainerName),
			CustomLabels:  req.CustomLabels,
			Address:       req.Address,
			Region:        pointer.ToStringOrNil(req.Region),
		}
		return models.CreateNode(q, nodeType, params)

	default:
		return nil, err
	}
}

func findOrCreatePMMAgent(q *reform.Querier, node *models.Node, customLabels map[string]string) (*models.Agent, error) {
	agents, err := models.AgentFindAll(q)
	if err != nil {
		return nil, err
	}

	var res []*models.Agent
	for _, a := range agents {
		if a.AgentType == models.PMMAgentType && pointer.GetString(a.RunsOnNodeID) == node.NodeID {
			res = append(res, a)
		}
	}
	switch len(res) {
	case 0:
		return models.AgentAddPmmAgent(q, node.NodeID, nil)
	case 1:
		return res[0], nil
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "Found %d pmm-agents for Node %q.", len(res), node.NodeID)
	}
}

func (s *NodeService) findNodeExporterByPmmAgentID(q *reform.Querier, pmmAgentID string) (nodeExporter *inventorypb.NodeExporter, err error) {
	agents, err := models.AgentsRunningByPMMAgent(q, pmmAgentID)
	if err != nil {
		return nil, err
	}

	for _, a := range agents {
		if pointer.GetString(a.PMMAgentID) == pmmAgentID {
			invAgent, err := inventory.ToInventoryAgent(q, a, s.registry)
			if err != nil {
				return nodeExporter, err
			}
			nodeExporter = invAgent.(*inventorypb.NodeExporter)
			return nodeExporter, nil
		}
	}

	return nodeExporter, errAgentNotFound
}

func (s *NodeService) addPmmAgentToResponse(q *reform.Querier, model *models.Agent, res *managementpb.RegisterNodeResponse) error {
	invAgent, err := inventory.ToInventoryAgent(q, model, s.registry)
	if err != nil {
		return err
	}
	res.PmmAgent = invAgent.(*inventorypb.PMMAgent)
	return nil
}
