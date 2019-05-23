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
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/action"
)

//nolint:unused
type server struct {
	ss      actionsStarterStopper
	storage action.Storage
	rsr     action.PMMAgentIDResolver
	db      *reform.DB
}

// NewServer creates Management Actions Server.
func NewServer(ss actionsStarterStopper, s action.Storage, rsr action.PMMAgentIDResolver, db *reform.DB) managementpb.ActionsServer {
	return &server{
		ss:      ss,
		storage: s,
		db:      db,
		rsr:     rsr,
	}
}

// GetAction gets an action result.
//nolint:lll
func (s *server) GetAction(ctx context.Context, req *managementpb.GetActionRequest) (*managementpb.GetActionResponse, error) {
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
func (s *server) StartPTSummaryAction(ctx context.Context, req *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	a := action.NewPtSummary(s.db.Querier, s.ss, s.rsr)
	err := a.Prepare(req.NodeId, req.PmmAgentId, []string{}) // TODO: Add real pt-summary arguments
	if err != nil {
		return nil, err
	}

	err = a.Start(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartPTSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// StartPTMySQLSummaryAction starts pt-mysql-summary action.
//nolint:lll
func (s *server) StartPTMySQLSummaryAction(ctx context.Context, req *managementpb.StartPTMySQLSummaryActionRequest) (*managementpb.StartPTMySQLSummaryActionResponse, error) {
	a := action.NewPtMySQLSummary(s.db.Querier, s.ss, s.rsr)
	err := a.Prepare(req.ServiceId, req.PmmAgentId, []string{}) // TODO: Add real pt-mysql-summary arguments
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = a.Start(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartPTMySQLSummaryActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// StartMySQLExplainAction starts mysql-explain action.
//nolint:lll
func (s *server) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	a := action.NewMySQLExplain(s.db.Querier, s.ss, s.rsr)
	err := a.Prepare(req.ServiceId, req.PmmAgentId, req.Query)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = a.Start(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// StartMySQLExplainJSONAction starts mysql-explain json action.
//nolint:lll
func (s *server) StartMySQLExplainJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainJSONActionRequest) (*managementpb.StartMySQLExplainJSONActionResponse, error) {
	a := action.NewMySQLExplainJSON(s.db.Querier, s.ss, s.rsr)
	err := a.Prepare(req.ServiceId, req.PmmAgentId, req.Query)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	err = a.Start(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &managementpb.StartMySQLExplainJSONActionResponse{
		PmmAgentId: a.PMMAgentID,
		ActionId:   a.ID,
	}, err
}

// CancelAction stops an Action.
//nolint:lll
func (s *server) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	ar, ok := s.storage.Load(ctx, req.ActionId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Couldn't find action result record in storage")
	}

	err := s.ss.StopAction(ctx, ar.ID, ar.PmmAgentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Couldn't cancel action with ID: %s, reason: %v", req.ActionId, err)
	}

	return &managementpb.CancelActionResponse{}, nil
}
