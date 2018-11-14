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
	"sync"
	"sync/atomic"

	"github.com/percona/pmm/api/agent"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Conn struct {
	stream      agent.Agent_ConnectServer
	lastID      uint32
	l           *logrus.Entry
	rw          sync.RWMutex
	subscribers map[uint32][]chan *agent.AgentMessage
	requestChan chan *agent.AgentMessage
}

func NewConn(uuid string, stream agent.Agent_ConnectServer) *Conn {
	conn := &Conn{
		stream:      stream,
		l:           logrus.WithField("pmm-agent", uuid),
		subscribers: make(map[uint32][]chan *agent.AgentMessage),
		requestChan: make(chan *agent.AgentMessage),
	}
	// create goroutine to dispatch messages
	go conn.startResponseDispatcher()
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

	agentChan := make(chan *agent.AgentMessage)
	defer close(agentChan)

	c.addSubscriber(id, agentChan)
	defer c.removeSubscriber(id, agentChan)

	agentMessage := <-agentChan
	c.l.Debugf("Recv: %s.", agentMessage)

	return agentMessage, nil
}

func (c *Conn) RecvRequestMessage() *agent.AgentMessage {
	agentMessage := <-c.requestChan
	c.l.Debugf("Recv: %s.", agentMessage)
	return agentMessage
}

func (c *Conn) startResponseDispatcher() {
	for c.stream.Context().Err() != nil {
		agentMessage, err := c.stream.Recv()
		if err != nil {
			c.l.Warnln("Connection closed", err)
			return
		}

		switch agentMessage.GetPayload().(type) {
		case *agent.AgentMessage_Ping, *agent.AgentMessage_State:
			go func(agentMessage *agent.AgentMessage) {
				c.requestChan <- agentMessage
			}(agentMessage)
		case *agent.AgentMessage_Auth, *agent.AgentMessage_QanData:
			c.emit(agentMessage)
		}
	}
}

func (c *Conn) emit(message *agent.AgentMessage) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	if _, ok := c.subscribers[message.Id]; ok {
		for i := range c.subscribers[message.Id] {
			go func(subscriber chan *agent.AgentMessage) {
				subscriber <- message
			}(c.subscribers[message.Id][i])
		}
	} else {
		c.l.Warnf("Unexpected message: %T %s", message, message)
	}
}

func (c *Conn) removeSubscriber(id uint32, agentChan chan *agent.AgentMessage) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, ok := c.subscribers[id]; ok {
		for i := range c.subscribers[id] {
			if c.subscribers[id][i] == agentChan {
				c.subscribers[id] = append(c.subscribers[id][:i], c.subscribers[id][i+1:]...)
				break
			}
		}
	}
}

func (c *Conn) addSubscriber(id uint32, agentChan chan *agent.AgentMessage) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, ok := c.subscribers[id]; !ok {
		c.subscribers[id] = []chan *agent.AgentMessage{}
	}
	c.subscribers[id] = append(c.subscribers[id], agentChan)
}
