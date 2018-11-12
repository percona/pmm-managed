package agents

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/percona/pmm/api/agent"
)

type messageDispatcher struct {
	sync.RWMutex
	l           *logrus.Entry
	subscribers map[uint32]chan *agent.AgentMessage
}

func newMessageDispatcher(uuid string) *messageDispatcher {
	return &messageDispatcher{
		subscribers: make(map[uint32]chan *agent.AgentMessage),
		l:           logrus.WithField("message dispatcher", uuid),
	}
}

func (m *messageDispatcher) MessageReceived(message *agent.AgentMessage) {
	m.Lock()
	defer m.Unlock()
	m.subscribers[message.Id] <- message
	delete(m.subscribers, message.Id)
}

func (m *messageDispatcher) WaitForMessage(id uint32) *agent.AgentMessage {
	m.Lock()
	agentChan := make(chan *agent.AgentMessage)
	defer close(agentChan)
	m.subscribers[id] = agentChan
	m.Unlock()
	return <-agentChan
}
