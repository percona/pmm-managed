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

package agents

import (
	"context"

	"github.com/percona/pmm/api/agentpb"
)

// AgentServer provides methods for pmm-agent <-> pmm-managed interactions.
type agentGrpcServer struct {
	registry *Registry
}

// NewAgentGrpcServer creates new agents server.
func NewAgentGrpcServer(r *Registry) agentpb.AgentServer {
	return &agentGrpcServer{
		registry: r,
	}
}

// Register TODO https://jira.percona.com/browse/PMM-3453
func (s *agentGrpcServer) Register(context.Context, *agentpb.RegisterRequest) (*agentpb.RegisterResponse, error) {
	panic("not implemented yet")
}

// Connect establishes two-way communication channel between pmm-agent and pmm-managed.
func (s *agentGrpcServer) Connect(stream agentpb.Agent_ConnectServer) error {
	return s.registry.Run(stream)
}

// check interfaces
var (
	_ agentpb.AgentServer = (*agentGrpcServer)(nil)
)
