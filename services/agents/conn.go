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

// Package agents contains business logic of working with pmm-agents.
package agents

import (
	"fmt"
	"sync/atomic"

	"github.com/Percona-Lab/pmm-api/agent"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Conn struct {
	stream agent.Agent_ConnectServer
	lastID uint32
}

func NewConn(stream agent.Agent_ConnectServer) *Conn {
	return &Conn{
		stream: stream,
	}
}

func (c *Conn) SendAndRecv(toAgent agent.ServerMessagePayload) (*agent.AgentMessage, error) {
	id := atomic.AddUint32(&c.lastID, 1)
	err := c.stream.Send(&agent.ServerMessage{
		Id:      id,
		Payload: toAgent,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send message to agent")
	}

	// FIXME Instead of reading the next message and dropping it if it is not what we expect,
	//       we should wait for the right message.
	//       We should have a single stream receiver in a separate goroutine,
	//       and internal subscriptions for IDs.

	fromAgent, err := c.stream.Recv()
	if err != nil {
		return nil, errors.Wrap(err, "failed to receive message from agent")
	}
	if fromAgent.Id != id {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("expected message from agent with ID %d, got ID %d", id, fromAgent.Id))
	}
	return fromAgent, nil
}
