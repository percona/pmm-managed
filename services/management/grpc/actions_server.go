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

	"github.com/percona/pmm-managed/action"
)

//nolint:unused
type actionsServer struct {
	r action.Runner
	s action.Storage
	f action.Factory
}

// NewActionsServer creates Management Actions Server.
func NewActionsServer(r action.Runner, s action.Storage, f action.Factory) managementpb.ActionsServer {
	return &actionsServer{r, s, f}
}

// GetAction gets an action result.
//nolint:lll
func (s *actionsServer) GetAction(ctx context.Context, req *managementpb.GetActionRequest) (*managementpb.GetActionResponse, error) {
	res, err := s.s.Load(ctx, req.ActionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
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
//nolint:lll,dupl
func (s *actionsServer) StartPTSummaryAction(ctx context.Context, req *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	a, err := s.f.NewPTSummary(ctx, req.NodeId, req.PmmAgentId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.r.StartPTSummaryAction(ctx, a)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartPTSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, nil
}

// StartPTMySQLSummaryAction starts pt-mysql-summary action.
//nolint:lll,dupl
func (s *actionsServer) StartPTMySQLSummaryAction(ctx context.Context, req *managementpb.StartPTMySQLSummaryActionRequest) (*managementpb.StartPTMySQLSummaryActionResponse, error) {
	a, err := s.f.NewPTMySQLSummary(ctx, req.ServiceId, req.PmmAgentId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.r.StartPTMySQLSummaryAction(ctx, a)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartPTMySQLSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, nil
}

// StartMySQLExplainAction starts mysql-explain action.
//nolint:lll,dupl
func (s *actionsServer) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	a, err := s.f.NewMySQLExplain(ctx, req.ServiceId, req.PmmAgentId, req.Query)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.r.StartMySQLExplainAction(ctx, a)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, nil
}

// StartMySQLExplainJSONAction starts mysql-explain json action.
//nolint:lll,dupl
func (s *actionsServer) StartMySQLExplainJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainJSONActionRequest) (*managementpb.StartMySQLExplainJSONActionResponse, error) {
	a, err := s.f.NewMySQLExplainJSON(ctx, req.ServiceId, req.PmmAgentId, req.Query)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.r.StartMySQLExplainJSONAction(ctx, a)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.StartMySQLExplainJSONActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, nil
}

// CancelAction stops an Action.
//nolint:lll
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	ar, err := s.s.Load(ctx, req.ActionId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.r.StopAction(ctx, ar.ID, ar.PmmAgentID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &managementpb.CancelActionResponse{}, nil
}
