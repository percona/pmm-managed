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

package action

import (
	"context"

	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

//nolint:unused
type gRPCServer struct {
	service service
	storage *InMemoryStorage
	db      *reform.DB
}

// NewGRPCServer creates Management Actions Server.
func NewGRPCServer(svc service, s *InMemoryStorage, db *reform.DB) managementpb.ActionsServer {
	return &gRPCServer{
		service: svc,
		storage: s,
		db:      db,
	}
}

// GetAction gets an action result.
//nolint:lll
func (s *gRPCServer) GetAction(ctx context.Context, req *managementpb.GetActionRequest) (*managementpb.GetActionResponse, error) {
	res, ok := s.storage.Load(ctx, req.ActionId)
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

// StartPTSummaryAction starts pt-summary action.
//nolint:lll
func (s *gRPCServer) StartPTSummaryAction(ctx context.Context, req *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	// TODO: Add real pt-summary arguments
	a := NewPtSummary(req.NodeId, req.PmmAgentId, []string{}, s.db.Querier)
	err := a.Prepare()
	if err != nil {
		return nil, err
	}

	err = s.service.StartPTSummaryAction(ctx, a)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartPTSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.Id,
	}, err
}

// StartPTMySQLSummaryAction starts pt-mysql-summary action.
//nolint:lll
func (s *gRPCServer) StartPTMySQLSummaryAction(ctx context.Context, req *managementpb.StartPTMySQLSummaryActionRequest) (*managementpb.StartPTMySQLSummaryActionResponse, error) {
	// TODO: Add real pt-mysql-summary arguments
	a := NewPtMySQLSummary(req.ServiceId, req.PmmAgentId, []string{}, s.db.Querier)
	err := a.Prepare()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.service.StartPTMySQLSummaryAction(ctx, a)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartPTMySQLSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.Id,
	}, err
}

// StartMySQLExplainAction starts mysql-explain action.
//nolint:lll
func (s *gRPCServer) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	a := NewMySQLExplain(req.ServiceId, req.PmmAgentId, req.Query, s.db.Querier)
	err := a.Prepare()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.service.StartMySQLExplainAction(ctx, a)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.Id,
	}, err
}

// StartMySQLExplainJSONAction starts mysql-explain json action.
//nolint:lll
func (s *gRPCServer) StartMySQLExplainJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainJSONActionRequest) (*managementpb.StartMySQLExplainJSONActionResponse, error) {
	a := NewMySQLExplainJSON(req.ServiceId, req.PmmAgentId, req.Query, s.db.Querier)
	err := a.Prepare()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = s.service.StartMySQLExplainJSONAction(ctx, a)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartMySQLExplainJSONActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.Id,
	}, err
}

// CancelAction stops an Action.
//nolint:lll
func (s *gRPCServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	err := s.service.StopAction(ctx, req.ActionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Couldn't cancel action with ID: %s, reason: %v", req.ActionId, err)
	}

	return &managementpb.CancelActionResponse{}, nil
}
