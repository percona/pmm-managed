// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package handlers

import (
	"context"
	"time"

	"github.com/Percona-Lab/pmm-api/agent"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/percona/pmm-managed/services/agents"
	"github.com/percona/pmm-managed/utils/logger"
)

type AgentServer struct {
}

func (s *AgentServer) Register(ctx context.Context, req *agent.RegisterRequest) (*agent.RegisterResponse, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &agent.RegisterResponse{
		Uuid: uuid.String(),
	}, nil
}

func (s *AgentServer) Connect(stream agent.Agent_ConnectServer) error {
	l := logger.Get(stream.Context())

	// connect request/response
	agentMessage, err := stream.Recv()
	if err != nil {
		return err
	}
	connect := agentMessage.GetConnect()
	if connect == nil {
		return errors.Errorf("Expected ConnectRequest, got %T.", agentMessage.Payload)
	}
	l.Infof("Got %T %s.", connect, connect)
	serverMessage := &agent.ServerMessage{
		Payload: &agent.ServerMessage_Connect{
			Connect: &agent.ConnectResponse{},
		},
	}
	if err = stream.Send(serverMessage); err != nil {
		return err
	}

	conn := agents.NewConn(stream)
	for stream.Context().Err() == nil {
		time.Sleep(time.Second)

		start := time.Now()
		fromAgent, err := conn.SendAndRecv(&agent.ServerMessage_Ping{
			Ping: &agent.PingRequest{},
		})
		if err != nil {
			return err
		}
		latency := time.Since(start) / 2
		ping := fromAgent.GetPing()
		if ping == nil {
			return errors.Errorf("Expected PingResponse, got %T.", agentMessage.Payload)
		}
		t, err := ptypes.Timestamp(ping.CurrentTime)
		if err != nil {
			return err
		}
		l.Infof("Latency: %s. Time drift: %s.", latency, t.Sub(start.Add(latency)))
	}

	return nil
}

// check interfaces
var (
	_ agent.AgentServer = (*AgentServer)(nil)
)
