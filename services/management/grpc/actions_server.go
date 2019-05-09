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

package grpc

import (
	"context"

	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/services/management"
)

//nolint:unused
type actionsServer struct {
	as *management.ActionsService
}

// NewManagementActionsServer creates Management Actions Server.
func NewManagementActionsServer(s *management.ActionsService) managementpb.ActionsServer {
	return &actionsServer{as: s}
}

// RunAction runs an Action.
func (s *actionsServer) RunAction(ctx context.Context, req *managementpb.RunActionRequest) (*managementpb.RunActionResponse, error) {
	p := management.RunActionParams{
		ActionName:   req.ActionName,
		ActionParams: req.ActionParams,
		PmmAgentID:   req.PmmAgentId,
		NodeID:       req.NodeId,
		ServiceID:    req.ServiceId,
	}
	// TODO: Handle errors.
	actionID, _ := s.as.RunAction(ctx, p)
	return &managementpb.RunActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, nil
}

// CancelAction stops an Action.
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	s.as.CancelAction(ctx, req.PmmAgentId, req.ActionId)
	return &managementpb.CancelActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   req.ActionId,
	}, nil
}

// GetActionResult gets an Action result.
func (s *actionsServer) GetActionResult(ctx context.Context, req *managementpb.GetActionResultRequest) (*managementpb.GetActionResultResponse, error) {
	res, ok := s.as.GetActionResult(ctx, req.ActionId)
	if !ok {
		return nil, status.Error(codes.NotFound, "ActionResult with given ID wasn't found")
	}

	return &managementpb.GetActionResultResponse{
		Id:         res.ID,
		PmmAgentId: res.PmmAgentID,
		Output:     res.Output,
	}, nil
}
