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

	"github.com/percona/pmm/api/agentpb"
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

// GetAction gets an action result.
//nolint:lll
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

// StartPTSummaryAction starts pt-summary action.
//nolint:lll
func (s *actionsServer) StartPTSummaryAction(ctx context.Context, req *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	p := management.RunProcessActionParams{
		ActionName: managementpb.ActionType_PT_SUMMARY,
		PmmAgentID: req.PmmAgentId,
		NodeID:     req.NodeId,
	}

	actionID, err := s.as.RunProcessAction(ctx, &p)
	return &managementpb.StartPTSummaryActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

// StartPTMySQLSummaryAction starts pt-mysql-summary action.
//nolint:lll
func (s *actionsServer) StartPTMySQLSummaryAction(ctx context.Context, req *managementpb.StartPTMySQLSummaryActionRequest) (*managementpb.StartPTMySQLSummaryActionResponse, error) {
	p := management.RunProcessActionParams{
		ActionName: managementpb.ActionType_PT_MYSQL_SUMMARY,
		PmmAgentID: req.PmmAgentId,
		ServiceID:  req.ServiceId,
	}

	actionID, err := s.as.RunProcessAction(ctx, &p)
	return &managementpb.StartPTMySQLSummaryActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

// StartMySQLExplainAction starts mysql-explain action.
//nolint:lll
func (s *actionsServer) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	p := management.StartMySQLExplainActionParams{
		PmmAgentID:   req.PmmAgentId,
		ServiceID:    req.ServiceId,
		Query:        req.Query,
		OutputFormat: agentpb.MysqlExplainOutputFormat_MYSQL_EXPLAIN_OUTPUT_FORMAT_DEFAULT,
	}

	actionID, err := s.as.StartMySQLExplainAction(ctx, &p)
	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

// StartMySQLExplainJSONAction starts mysql-explain json action.
//nolint:lll
func (s *actionsServer) StartMySQLExplainJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainJSONActionRequest) (*managementpb.StartMySQLExplainJSONActionResponse, error) {
	p := management.StartMySQLExplainActionParams{
		PmmAgentID:   req.PmmAgentId,
		ServiceID:    req.ServiceId,
		Query:        req.Query,
		OutputFormat: agentpb.MysqlExplainOutputFormat_MYSQL_EXPLAIN_OUTPUT_FORMAT_JSON,
	}

	actionID, err := s.as.StartMySQLExplainAction(ctx, &p)
	return &managementpb.StartMySQLExplainJSONActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   actionID,
	}, err
}

// CancelAction stops an Action.
//nolint:lll
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	s.as.CancelAction(ctx, req.ActionId)
	return &managementpb.CancelActionResponse{}, nil
}
