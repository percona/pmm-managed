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

// Package agents contains business logic of working with pmm-agents.
package agents

import (
	"sync/atomic"

	"github.com/percona/pmm/api/agent"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Conn struct {
	stream            agent.Agent_ConnectServer
	lastID            uint32
	l                 *logrus.Entry
	messageDispatcher *messageDispatcher
}

func NewConn(uuid string, stream agent.Agent_ConnectServer) *Conn {
	// Create goroutine to dispatch messages
	conn := &Conn{
		stream:            stream,
		l:                 logrus.WithField("pmm-agent", uuid),
		messageDispatcher: newMessageDispatcher(uuid),
	}
	go conn.startMessageDispatcher()
	return conn
}

func (c *Conn) SendAndRecv(toAgent agent.ServerMessagePayload) (*agent.AgentMessage, error) {
	id := atomic.AddUint32(&c.lastID, 1)
	serverMessage := &agent.ServerMessage{
		Id:      id,
		Payload: toAgent,
	}
	c.l.Debugf("Send: %s.", serverMessage)
	if err := c.stream.Send(serverMessage); err != nil {
		return nil, errors.Wrap(err, "failed to send message to agent")
	}

	agentMessage := c.messageDispatcher.WaitForMessage(id)
	c.l.Debugf("Recv: %s.", agentMessage)

	return agentMessage, nil
}
func (c *Conn) startMessageDispatcher() {
	for {
		select {
		case <-c.stream.Context().Done():
			c.l.Debugln("Close connection: ", c.lastID)
			return
		default:
			agentMessage, err := c.stream.Recv()
			if err != nil {
				errorStatus, ok := status.FromError(err)
				if ok && errorStatus.Code() == codes.Canceled {
					c.l.Debugln("Connection closed from other side")
					return
				}
				c.l.Fatal(errors.Wrap(err, "failed to receive message from agent"))
			}
			c.messageDispatcher.MessageReceived(agentMessage)
		}
	}
}
