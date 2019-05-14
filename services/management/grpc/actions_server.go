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

func (s *actionsServer) GetAction(ctx context.Context, req *managementpb.GetActionRequest) (*managementpb.GetActionResponse, error) {
	res, ok := s.as.GetActionResult(ctx, req.ActionId)
	if !ok {
		return nil, status.Error(codes.NotFound, "ActionResult with given ID wasn't found")
	}

	return &managementpb.GetActionResponse{
		ActionId:   res.ID,
		PmmAgentId: res.PmmAgentID,
		Done:       res.Done,
		Error:      res.Error,
		Output:     res.Output,
	}, nil
}

func (s *actionsServer) StartPTSummaryAction(ctx context.Context, req *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	p := management.RunActionParams{
		ActionName: managementpb.ActionType_PT_SUMMARY,
		PmmAgentID: req.PmmAgentId,
		NodeID:     req.NodeId,
	}

	actionID, err := s.as.RunAction(ctx, p)
	return &managementpb.StartPTSummaryActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

func (s *actionsServer) StartPTMySQLSummaryAction(ctx context.Context, req *managementpb.StartPTMySQLSummaryActionRequest) (*managementpb.StartPTMySQLSummaryActionResponse, error) {
	p := management.RunActionParams{
		ActionName: managementpb.ActionType_PT_MYSQL_SUMMARY,
		PmmAgentID: req.PmmAgentId,
		ServiceID:  req.ServiceId,
	}

	actionID, err := s.as.RunAction(ctx, p)
	return &managementpb.StartPTMySQLSummaryActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

func (s *actionsServer) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	p := management.RunActionParams{
		ActionName: managementpb.ActionType_MYSQL_EXPLAIN,
		PmmAgentID: req.PmmAgentId,
		ServiceID:  req.ServiceId,
	}

	actionID, err := s.as.RunAction(ctx, p)
	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

// NewManagementActionsServer creates Management Actions Server.
func NewManagementActionsServer(s *management.ActionsService) managementpb.ActionsServer {
	return &actionsServer{as: s}
}

// CancelAction stops an Action.
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	s.as.CancelAction(ctx, req.ActionId)
	return &managementpb.CancelActionResponse{}, nil
}
