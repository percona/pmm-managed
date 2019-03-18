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

package handlers

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	inventorypb "github.com/percona/pmm/api/inventory"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/inventory"
)

type nodesServer struct {
	s *inventory.NodesService
}

// NewNodesServer returns Inventory API handler for managing Nodes.
func NewNodesServer(s *inventory.NodesService) inventorypb.NodesServer {
	return &nodesServer{
		s: s,
	}
}

// ListNodes returns a list of all Nodes.
func (s *nodesServer) ListNodes(ctx context.Context, req *inventorypb.ListNodesRequest) (*inventorypb.ListNodesResponse, error) {
	nodes, err := s.s.List(ctx)
	if err != nil {
		return nil, err
	}

	res := new(inventorypb.ListNodesResponse)
	for _, node := range nodes {
		switch node := node.(type) {
		case *inventorypb.GenericNode:
			res.Generic = append(res.Generic, node)
		case *inventorypb.ContainerNode:
			res.Container = append(res.Container, node)
		case *inventorypb.RemoteNode:
			res.Remote = append(res.Remote, node)
		case *inventorypb.RemoteAmazonRDSNode:
			res.RemoteAmazonRds = append(res.RemoteAmazonRds, node)
		default:
			panic(fmt.Errorf("unhandled inventory Node type %T", node))
		}
	}
	return res, nil
}

// GetNode returns a single Node by ID.
func (s *nodesServer) GetNode(ctx context.Context, req *inventorypb.GetNodeRequest) (*inventorypb.GetNodeResponse, error) {
	node, err := s.s.Get(ctx, req.NodeId, nil)
	if err != nil {
		return nil, err
	}

	res := new(inventorypb.GetNodeResponse)
	switch node := node.(type) {
	case *inventorypb.GenericNode:
		res.Node = &inventorypb.GetNodeResponse_Generic{Generic: node}
	case *inventorypb.ContainerNode:
		res.Node = &inventorypb.GetNodeResponse_Container{Container: node}
	case *inventorypb.RemoteNode:
		res.Node = &inventorypb.GetNodeResponse_Remote{Remote: node}
	case *inventorypb.RemoteAmazonRDSNode:
		res.Node = &inventorypb.GetNodeResponse_RemoteAmazonRds{RemoteAmazonRds: node}
	default:
		panic(fmt.Errorf("unhandled inventory Node type %T", node))
	}
	return res, nil
}

// AddGenericNode adds Generic Node.
func (s *nodesServer) AddGenericNode(ctx context.Context, req *inventorypb.AddGenericNodeRequest) (*inventorypb.AddGenericNodeResponse, error) {
	node, err := s.s.Add(ctx, models.GenericNodeType, req.NodeName, pointer.ToStringOrNil(req.Address), nil)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.AddGenericNodeResponse{
		Generic: node.(*inventorypb.GenericNode),
	}
	return res, nil
}

// AddContainerNode adds Container Node.
func (s *nodesServer) AddContainerNode(ctx context.Context, req *inventorypb.AddContainerNodeRequest) (*inventorypb.AddContainerNodeResponse, error) {
	node, err := s.s.Add(ctx, models.ContainerNodeType, req.NodeName, nil, nil)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.AddContainerNodeResponse{
		Container: node.(*inventorypb.ContainerNode),
	}
	return res, nil
}

// AddRemoteNode adds Remote Node.
func (s *nodesServer) AddRemoteNode(ctx context.Context, req *inventorypb.AddRemoteNodeRequest) (*inventorypb.AddRemoteNodeResponse, error) {
	node, err := s.s.Add(ctx, models.RemoteNodeType, req.NodeName, nil, nil)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.AddRemoteNodeResponse{
		Remote: node.(*inventorypb.RemoteNode),
	}
	return res, nil
}

// AddRemoteAmazonRDSNode adds Amazon (AWS) RDS remote Node.
func (s *nodesServer) AddRemoteAmazonRDSNode(ctx context.Context, req *inventorypb.AddRemoteAmazonRDSNodeRequest) (*inventorypb.AddRemoteAmazonRDSNodeResponse, error) {
	node, err := s.s.Add(ctx, models.RemoteAmazonRDSNodeType, req.NodeName, &req.Instance, &req.Region)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.AddRemoteAmazonRDSNodeResponse{
		RemoteAmazonRds: node.(*inventorypb.RemoteAmazonRDSNode),
	}
	return res, nil
}

// ChangeGenericNode changes Generic Node.
func (s *nodesServer) ChangeGenericNode(ctx context.Context, req *inventorypb.ChangeGenericNodeRequest) (*inventorypb.ChangeGenericNodeResponse, error) {
	node, err := s.s.Change(ctx, req.NodeId, req.NodeName)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.ChangeGenericNodeResponse{
		Generic: node.(*inventorypb.GenericNode),
	}
	return res, nil
}

// ChangeContainerNode changes Container Node.
func (s *nodesServer) ChangeContainerNode(ctx context.Context, req *inventorypb.ChangeContainerNodeRequest) (*inventorypb.ChangeContainerNodeResponse, error) {
	node, err := s.s.Change(ctx, req.NodeId, req.NodeName)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.ChangeContainerNodeResponse{
		Container: node.(*inventorypb.ContainerNode),
	}
	return res, nil
}

// ChangeRemoteNode changes Remote Node.
func (s *nodesServer) ChangeRemoteNode(ctx context.Context, req *inventorypb.ChangeRemoteNodeRequest) (*inventorypb.ChangeRemoteNodeResponse, error) {
	node, err := s.s.Change(ctx, req.NodeId, req.NodeName)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.ChangeRemoteNodeResponse{
		Remote: node.(*inventorypb.RemoteNode),
	}
	return res, nil
}

// ChangeRemoteAmazonRDSNode changes Amazon (AWS) RDS remote Node.
func (s *nodesServer) ChangeRemoteAmazonRDSNode(ctx context.Context, req *inventorypb.ChangeRemoteAmazonRDSNodeRequest) (*inventorypb.ChangeRemoteAmazonRDSNodeResponse, error) {
	node, err := s.s.Change(ctx, req.NodeId, req.NodeName)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.ChangeRemoteAmazonRDSNodeResponse{
		RemoteAmazonRds: node.(*inventorypb.RemoteAmazonRDSNode),
	}
	return res, nil
}

// RemoveNode removes Node without any Agents and Services.
func (s *nodesServer) RemoveNode(ctx context.Context, req *inventorypb.RemoveNodeRequest) (*inventorypb.RemoveNodeResponse, error) {
	if err := s.s.Remove(ctx, req.NodeId); err != nil {
		return nil, err
	}

	return new(inventorypb.RemoveNodeResponse), nil
}
