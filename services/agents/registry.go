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
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	api "github.com/percona/pmm/api/agent"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
)

const (
	prometheusNamespace = "pmm_managed"
	prometheusSubsystem = "agents"
)

type agentInfo struct {
	connectedAt time.Time
	channel     *Channel
	l           *logrus.Entry
	kick        chan struct{}
}

type sharedChannelMetrics struct {
	mRecv prometheus.Counter
	mSend prometheus.Counter
}

type Registry struct {
	db *reform.DB

	rw     sync.RWMutex
	agents map[string]*agentInfo

	sharedMetrics *sharedChannelMetrics
	mConnects     prometheus.Counter
	mDisconnects  *prometheus.CounterVec
	mLatency      prometheus.Summary
	mTimeDrift    prometheus.Summary
}

func NewRegistry(db *reform.DB) *Registry {
	return &Registry{
		db:     db,
		agents: make(map[string]*agentInfo),
		sharedMetrics: &sharedChannelMetrics{
			mRecv: prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: prometheusNamespace,
				Subsystem: prometheusSubsystem,
				Name:      "messages_received_total",
				Help:      "A total number of messages received from pmm-agents.",
			}),
			mSend: prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: prometheusNamespace,
				Subsystem: prometheusSubsystem,
				Name:      "messages_sent_total",
				Help:      "A total number of messages sent to pmm-agents.",
			}),
		},
		mConnects: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "connects_total",
			Help:      "A total number of pmm-agent connects.",
		}),
		mDisconnects: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "disconnects_total",
			Help:      "A total number of pmm-agent disconnects.",
		}, []string{"reason"}),
		mLatency: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  prometheusNamespace,
			Subsystem:  prometheusSubsystem,
			Name:       "latency_seconds",
			Help:       "Ping latency.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		mTimeDrift: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  prometheusNamespace,
			Subsystem:  prometheusSubsystem,
			Name:       "time_drift_seconds",
			Help:       "Time drift.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
}

func (r *Registry) Run(stream api.Agent_ConnectServer) error {
	r.mConnects.Inc()
	disconnectReason := "unknown"
	defer func() {
		r.mDisconnects.WithLabelValues(disconnectReason).Inc()
	}()

	l := logger.Get(stream.Context())
	md := api.GetAgentConnectMetadata(stream.Context())
	if err := r.authenticate(&md); err != nil {
		l.Warnf("Failed to authenticate connected pmm-agent %+v.", md)
		disconnectReason = "auth"
		return err
	}
	l.Infof("Connected pmm-agent: %+v.", md)

	r.rw.Lock()
	defer r.rw.Unlock()

	if agent := r.agents[md.ID]; agent != nil {
		close(agent.kick)
	}

	agent := &agentInfo{
		connectedAt: time.Now(),
		channel:     NewChannel(stream, l.WithField("component", "channel"), r.sharedMetrics),
		l:           l,
		kick:        make(chan struct{}),
	}
	r.agents[md.ID] = agent

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.ping(agent.channel, l)

		case <-agent.kick:
			l.Warnf("Kicked.")
			disconnectReason = "kicked"
			return nil

		case msg := <-agent.channel.Requests():
			if msg == nil {
				disconnectReason = "done"
				return agent.channel.Wait()
			}

			switch req := msg.Payload.(type) {
			case *api.AgentMessage_QanData:
				// TODO
				agent.channel.SendResponse(&api.ServerMessage{
					Id: msg.Id,
					Payload: &api.ServerMessage_QanData{
						QanData: new(api.QANDataResponse),
					},
				})
			default:
				l.Warnf("Unexpected request: %s.", req)
				disconnectReason = "bad_request"
				return nil
			}
		}
	}
}

func (r *Registry) authenticate(md *api.AgentConnectMetadata) error {
	if md.ID == "" {
		return status.Error(codes.Unauthenticated, "Empty Agent ID.")
	}

	row := &models.AgentRow{ID: md.ID}
	if err := r.db.Reload(row); err != nil {
		if err == reform.ErrNoRows {
			return status.Errorf(codes.Unauthenticated, "No Agent with ID %q.", md.ID)
		}
		return errors.Wrap(err, "failed to find agent")
	}

	if row.Type != models.PMMAgentType {
		return status.Errorf(codes.Unauthenticated, "No pmm-agent with ID %q.", md.ID)
	}

	row.Version = &md.Version
	if err := r.db.Update(row); err != nil {
		return errors.Wrap(err, "failed to update agent")
	}
	return nil
}

func (r *Registry) ping(channel *Channel, l *logrus.Entry) {
	start := time.Now()
	res := channel.SendRequest(&api.ServerMessage_Ping{
		Ping: new(api.PingRequest),
	})
	if res == nil {
		return
	}
	latency := time.Since(start) / 2
	t, err := ptypes.Timestamp(res.(*api.AgentMessage_Ping).Ping.CurrentTime)
	if err != nil {
		l.Errorf("Failed to decode PingResponse.current_time: %s.", err)
		return
	}
	timeDrift := t.Sub(start.Add(latency))
	l.Infof("Latency: %s. Time drift: %s.", latency, timeDrift)
	r.mLatency.Observe(latency.Seconds())
	r.mTimeDrift.Observe(timeDrift.Seconds())
}

// Describe implements prometheus.Collector.
func (r *Registry) Describe(ch chan<- *prometheus.Desc) {
	r.sharedMetrics.mRecv.Describe(ch)
	r.sharedMetrics.mSend.Describe(ch)
	r.mConnects.Describe(ch)
	r.mDisconnects.Describe(ch)
}

// Collect implement prometheus.Collector.
func (r *Registry) Collect(ch chan<- prometheus.Metric) {
	r.sharedMetrics.mRecv.Collect(ch)
	r.sharedMetrics.mSend.Collect(ch)
	r.mConnects.Collect(ch)
	r.mDisconnects.Collect(ch)
}

// check interfaces
var (
	_ prometheus.Collector = (*Registry)(nil)
)
