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

// Package agents contains business logic of working with pmm-agent.
package agents

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/percona/pmm/api/agentpb"
	"github.com/pkg/errors"

	"github.com/percona/pmm-managed/utils/logger"
)

const (
	agentRequestsCap = 32
)

// Channel encapsulates two-way communication channel between pmm-managed and pmm-agent.
//
// All exported methods are thread-safe.
type Channel struct {
	//nolint:maligned
	s       agentpb.Agent_ConnectServer
	metrics *sharedChannelMetrics

	lastSentRequestID uint32

	sendM sync.Mutex

	m         sync.Mutex
	responses map[uint32]chan agentpb.AgentMessagePayload
	requests  chan *agentpb.AgentMessage

	closeOnce sync.Once
	closeWait chan struct{}
	closeErr  error
}

// NewChannel creates new two-way communication channel with given stream.
//
// Stream should not be used by the caller after channel is created.
func NewChannel(stream agentpb.Agent_ConnectServer, m *sharedChannelMetrics) *Channel {
	s := &Channel{
		s:       stream,
		metrics: m,

		responses: make(map[uint32]chan agentpb.AgentMessagePayload),
		requests:  make(chan *agentpb.AgentMessage, agentRequestsCap),

		closeWait: make(chan struct{}),
	}

	go s.runReceiver()
	return s
}

// close marks channel as closed with given error - only once.
func (c *Channel) close(err error) {
	c.closeOnce.Do(func() {
		logger.Get(c.s.Context()).Debugf("Closing with error: %+v", err)
		c.closeErr = err

		c.m.Lock()
		for _, ch := range c.responses { // unblock all subscribers
			close(ch)
		}
		c.responses = nil // prevent future subscriptions
		c.m.Unlock()

		close(c.closeWait)
	})
}

// Wait blocks until channel is closed and returns the reason why it was closed.
//
// When Wait returns, underlying gRPC connection should be terminated to prevent goroutine leak.
func (c *Channel) Wait() error {
	<-c.closeWait
	return c.closeErr
}

// Requests returns a channel for incoming requests. It must be read. It is closed on any error (see Wait).
func (c *Channel) Requests() <-chan *agentpb.AgentMessage {
	return c.requests
}

// SendResponse sends message to pmm-managed. It is no-op once channel is closed (see Wait).
func (c *Channel) SendResponse(msg *agentpb.ServerMessage) {
	c.send(msg)
}

// SendRequest sends request to pmm-managed, blocks until response is available, and returns it.
// Response will be nil if channel is closed.
// It is no-op once channel is closed (see Wait).
func (c *Channel) SendRequest(payload agentpb.ServerMessagePayload) agentpb.AgentMessagePayload {
	id := atomic.AddUint32(&c.lastSentRequestID, 1)
	ch := c.subscribe(id)

	c.send(&agentpb.ServerMessage{
		Id:      id,
		Payload: payload,
	})

	return <-ch
}

func (c *Channel) send(msg *agentpb.ServerMessage) {
	c.sendM.Lock()
	select {
	case <-c.closeWait:
		c.sendM.Unlock()
		return
	default:
	}

	logger.Get(c.s.Context()).Debugf("Sending message: %s.", msg)
	err := c.s.Send(msg)
	c.sendM.Unlock()
	if err != nil {
		c.close(errors.Wrap(err, "failed to send message"))
		return
	}
	c.metrics.mSend.Inc()
}

// runReader receives messages from server
func (c *Channel) runReceiver() {
	defer func() {
		close(c.requests)
		logger.Get(c.s.Context()).Debug("Exiting receiver goroutine.")
	}()

	for {
		msg, err := c.s.Recv()
		if err != nil {
			c.close(errors.Wrap(err, "failed to receive message"))
			return
		}
		logger.Get(c.s.Context()).Debugf("Received message: %s.", msg)
		c.metrics.mRecv.Inc()

		switch msg.Payload.(type) {
		// requests
		case *agentpb.AgentMessage_Ping, *agentpb.AgentMessage_StateChanged, *agentpb.AgentMessage_QanCollect:
			c.requests <- msg

		// responses
		case *agentpb.AgentMessage_Pong,
			*agentpb.AgentMessage_SetState,
			*agentpb.AgentMessage_ActionRunResponse,
			*agentpb.AgentMessage_ActionCancelResponse:
			c.publish(msg.Id, msg.Payload)
		case *agentpb.AgentMessage_ActionResult:
			// TODO: PMM-3978: Doing something with ActionResult. For example push it to UI...

		default:
			c.close(errors.Errorf("failed to handle received message %s", msg))
			return
		}
	}
}

func (c *Channel) subscribe(id uint32) chan agentpb.AgentMessagePayload {
	ch := make(chan agentpb.AgentMessagePayload, 1)

	c.m.Lock()
	if c.responses == nil { // Channel is closed, no more subscriptions
		c.m.Unlock()
		close(ch)
		return ch
	}

	_, ok := c.responses[id]
	if ok {
		// it is possible only on lastSentRequestID wrap around, and we can't recover from that
		logger.Get(c.s.Context()).Panicf("Already have subscriber for ID %d.", id)
	}

	c.responses[id] = ch
	c.m.Unlock()
	return ch
}

func (c *Channel) publish(id uint32, payload agentpb.AgentMessagePayload) {
	c.m.Lock()
	if c.responses == nil { // Channel is closed, no more publishing
		c.m.Unlock()
		return
	}

	ch := c.responses[id]
	if ch == nil {
		c.m.Unlock()
		c.close(errors.WithStack(fmt.Errorf("no subscriber for ID %d", id)))
		return
	}

	delete(c.responses, id)
	c.m.Unlock()
	ch <- payload
}
