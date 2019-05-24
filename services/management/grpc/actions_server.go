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

	"github.com/percona/pmm-managed/services/action"
)

//nolint:unused
type actionsServer struct {
	svc action.Service
}

// NewServer creates Management Actions Server.
func NewActionsServer(svc action.Service) managementpb.ActionsServer {
	return &actionsServer{svc: svc}
}

// GetAction gets an action result.
//nolint:lll
func (s *actionsServer) GetAction(ctx context.Context, req *managementpb.GetActionRequest) (*managementpb.GetActionResponse, error) {
	res, err := s.svc.GetActionResult(ctx, req.ActionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &managementpb.GetActionResponse{
		ActionId:   res.ID,
		PmmAgentId: res.PmmAgentID,
		Done:       res.Done,
		Error:      res.Error,
		Output:     res.Output,
	}, nil
}

// StartPTSummaryAction starts pt-summary action.
//nolint:lll
func (s *actionsServer) StartPTSummaryAction(ctx context.Context, req *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	a, err := s.svc.StartPTSummaryAction(ctx, req.PmmAgentId, req.NodeId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &managementpb.StartPTSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// StartPTMySQLSummaryAction starts pt-mysql-summary action.
//nolint:lll
func (s *actionsServer) StartPTMySQLSummaryAction(ctx context.Context, req *managementpb.StartPTMySQLSummaryActionRequest) (*managementpb.StartPTMySQLSummaryActionResponse, error) {
	a, err := s.svc.StartPTMySQLSummaryAction(ctx, req.PmmAgentId, req.ServiceId, []string{})
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartPTMySQLSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// StartMySQLExplainAction starts mysql-explain action.
//nolint:lll
func (s *actionsServer) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	a, err := s.svc.StartMySQLExplainAction(ctx, req.PmmAgentId, req.ServiceId, req.Query)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// StartMySQLExplainJSONAction starts mysql-explain json action.
//nolint:lll
func (s *actionsServer) StartMySQLExplainJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainJSONActionRequest) (*managementpb.StartMySQLExplainJSONActionResponse, error) {
	a, err := s.svc.StartMySQLExplainJSONAction(ctx, req.PmmAgentId, req.ServiceId, req.Query)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartMySQLExplainJSONActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// CancelAction stops an Action.
//nolint:lll
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	err := s.svc.CancelAction(ctx, req.ActionId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &managementpb.CancelActionResponse{}, nil
}
