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
)

type actionsService interface {
	RunAction(ctx context.Context, pmmAgentID string, actionName agentpb.ActionName)
	CancelAction(ctx context.Context, pmmAgentID, actionID string)
}

//nolint:unused
type actionsServer struct {
	as actionsService
}

// NewManagementActionsServer creates Management Actions Server.
func NewManagementActionsServer(s actionsService) managementpb.ActionsServer {
	return &actionsServer{as: s}
}

// RunAction runs an Action.
func (s *actionsServer) RunAction(ctx context.Context, req *managementpb.RunActionRequest) (*managementpb.RunActionResponse, error) {
	s.as.RunAction(ctx, req.PmmAgentId, agentpb.ActionName_PT_SUMMARY)
	return nil, nil
}

// CancelAction stops an Action.
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	s.as.CancelAction(ctx, req.PmmAgentId, req.ActionId)
	return nil, nil
}
