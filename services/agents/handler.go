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
	"runtime/pprof"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents/channel"
	"github.com/percona/pmm-managed/utils/logger"
)

var (
	checkExternalExporterConnectionPMMVersion = version.MustParse("1.14.99")

	defaultActionTimeout      = durationpb.New(10 * time.Second)
	defaultQueryActionTimeout = durationpb.New(15 * time.Second) // should be less than checks.resultTimeout
	defaultPtActionTimeout    = durationpb.New(30 * time.Second) // Percona-toolkit action timeout

	mSentDesc = prom.NewDesc(
		prom.BuildFQName(prometheusNamespace, prometheusSubsystem, "messages_sent_total"),
		"A total number of messages sent to pmm-agent.",
		[]string{"agent_id"},
		nil,
	)
	mRecvDesc = prom.NewDesc(
		prom.BuildFQName(prometheusNamespace, prometheusSubsystem, "messages_received_total"),
		"A total number of messages received from pmm-agent.",
		[]string{"agent_id"},
		nil,
	)
	mResponsesDesc = prom.NewDesc(
		prom.BuildFQName(prometheusNamespace, prometheusSubsystem, "messages_response_queue_length"),
		"The current length of the response queue.",
		[]string{"agent_id"},
		nil,
	)
	mRequestsDesc = prom.NewDesc(
		prom.BuildFQName(prometheusNamespace, prometheusSubsystem, "messages_request_queue_length"),
		"The current length of the request queue.",
		[]string{"agent_id"},
		nil,
	)
)

// Handler handles agent requests.
type Handler struct {
	db          *reform.DB
	r           *Registry
	qanClient   qanClient
	jobsService *JobsService

	mConnects    prom.Counter
	mDisconnects *prom.CounterVec
	mRoundTrip   prom.Summary
	mClockDrift  prom.Summary
}

// NewHandler creates new agents handler.
func NewHandler(db *reform.DB, qanClient qanClient, registry *Registry, service *JobsService) *Handler {
	h := &Handler{
		db:          db,
		qanClient:   qanClient,
		r:           registry,
		jobsService: service,

		mConnects: prom.NewCounter(prom.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "connects_total",
			Help:      "A total number of pmm-agent connects.",
		}),
		mDisconnects: prom.NewCounterVec(prom.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "disconnects_total",
			Help:      "A total number of pmm-agent disconnects.",
		}, []string{"reason"}),
		mRoundTrip: prom.NewSummary(prom.SummaryOpts{
			Namespace:  prometheusNamespace,
			Subsystem:  prometheusSubsystem,
			Name:       "round_trip_seconds",
			Help:       "Round-trip time.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		mClockDrift: prom.NewSummary(prom.SummaryOpts{
			Namespace:  prometheusNamespace,
			Subsystem:  prometheusSubsystem,
			Name:       "clock_drift_seconds",
			Help:       "Clock drift.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}

	// initialize metrics with labels
	h.mDisconnects.WithLabelValues("unknown")
	return h

}

// Run takes over pmm-agent gRPC stream and runs it until completion.
func (h *Handler) Run(stream agentpb.Agent_ConnectServer) error {
	h.mConnects.Inc()
	disconnectReason := "unknown"
	defer func() {
		h.mDisconnects.WithLabelValues(disconnectReason).Inc()
	}()

	ctx := stream.Context()
	l := logger.Get(ctx)
	agent, err := h.r.register(stream)
	if err != nil {
		disconnectReason = "auth"
		return err
	}
	defer func() {
		l.Infof("Disconnecting client: %s.", disconnectReason)
	}()

	// run pmm-agent state update loop for the current agent.
	go h.runStateChangeHandler(ctx, agent)

	h.r.RequestStateUpdate(ctx, agent.id)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.ping(ctx, agent)

		// see unregister and Kick methods
		case <-agent.kick:
			// already unregistered, no need to call unregister method
			l.Warn("Kicked.")
			disconnectReason = "kicked"
			err = status.Errorf(codes.Aborted, "Kicked.")
			return err

		case req := <-agent.channel.Requests():
			if req == nil {
				disconnectReason = "done"
				err = agent.channel.Wait()
				h.r.unregister(agent.id)
				if err != nil {
					l.Error(errors.WithStack(err))
				}
				return h.r.updateAgentStatusForChildren(ctx, agent.id, inventorypb.AgentStatus_DONE, 0)
			}

			switch p := req.Payload.(type) {
			case *agentpb.Ping:
				agent.channel.Send(&channel.ServerResponse{
					ID: req.ID,
					Payload: &agentpb.Pong{
						CurrentTime: ptypes.TimestampNow(),
					},
				})

			case *agentpb.StateChangedRequest:
				pprof.Do(ctx, pprof.Labels("request", "StateChangedRequest"), func(ctx context.Context) {
					if err := h.r.stateChanged(ctx, p); err != nil {
						l.Errorf("%+v", err)
					}

					agent.channel.Send(&channel.ServerResponse{
						ID:      req.ID,
						Payload: new(agentpb.StateChangedResponse),
					})
				})

			case *agentpb.QANCollectRequest:
				pprof.Do(ctx, pprof.Labels("request", "QANCollectRequest"), func(ctx context.Context) {
					if err := h.qanClient.Collect(ctx, p.MetricsBucket); err != nil {
						l.Errorf("%+v", err)
					}

					agent.channel.Send(&channel.ServerResponse{
						ID:      req.ID,
						Payload: new(agentpb.QANCollectResponse),
					})
				})

			case *agentpb.ActionResultRequest:
				// TODO: PMM-3978: In the future we need to merge action parts before send it to storage.
				err := models.ChangeActionResult(h.db.Querier, p.ActionId, agent.id, p.Error, string(p.Output), p.Done)
				if err != nil {
					l.Warnf("Failed to change action: %+v", err)
				}

				if !p.Done && p.Error != "" {
					l.Warnf("Action was done with an error: %v.", p.Error)
				}

				agent.channel.Send(&channel.ServerResponse{
					ID:      req.ID,
					Payload: new(agentpb.ActionResultResponse),
				})

			case *agentpb.JobResult:

			case *agentpb.JobProgress:
				// TODO Handle job progress messages https://jira.percona.com/browse/PMM-7756

			case nil:
				l.Errorf("Unexpected request: %+v.", req)
			}
		}
	}
}

// runStateChangeHandler runs pmm-agent state update loop for given pmm-agent until ctx is canceled or agent is kicked.
func (h *Handler) runStateChangeHandler(ctx context.Context, agent *pmmAgentInfo) {
	l := logger.Get(ctx).WithField("agent_id", agent.id)

	l.Info("Starting runStateChangeHandler ...")
	defer l.Info("Done runStateChangeHandler.")

	// stateChangeChan, state update loop, and RequestStateUpdate method ensure that state
	// is reloaded when requested, but several requests are batched together to avoid too often reloads.
	// That allows the caller to just call RequestStateUpdate when it seems fit.
	if cap(agent.stateChangeChan) != 1 {
		panic("stateChangeChan should have capacity 1")
	}

	for {
		select {
		case <-ctx.Done():
			return

		case <-agent.kick:
			return

		case <-agent.stateChangeChan:
			// batch several update requests together by delaying the first one
			sleepCtx, sleepCancel := context.WithTimeout(ctx, updateBatchDelay)
			<-sleepCtx.Done()
			sleepCancel()

			if ctx.Err() != nil {
				return
			}

			nCtx, cancel := context.WithTimeout(ctx, stateChangeTimeout)
			err := h.r.sendSetStateRequest(nCtx, agent)
			if err != nil {
				l.Error(err)
				h.r.RequestStateUpdate(ctx, agent.id)
			}
			cancel()
		}
	}
}

// ping sends Ping message to given Agent, waits for Pong and observes round-trip time and clock drift.
func (h *Handler) ping(ctx context.Context, agent *pmmAgentInfo) error {
	l := logger.Get(ctx)
	start := time.Now()
	resp, err := agent.channel.SendAndWaitResponse(new(agentpb.Ping))
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	roundtrip := time.Since(start)
	agentTime, err := ptypes.Timestamp(resp.(*agentpb.Pong).CurrentTime)
	if err != nil {
		return errors.Wrap(err, "failed to decode Pong.current_time")
	}
	clockDrift := agentTime.Sub(start) - roundtrip/2
	if clockDrift < 0 {
		clockDrift = -clockDrift
	}
	l.Infof("Round-trip time: %s. Estimated clock drift: %s.", roundtrip, clockDrift)
	h.mRoundTrip.Observe(roundtrip.Seconds())
	h.mClockDrift.Observe(clockDrift.Seconds())
	return nil
}

// Describe implements prometheus.Collector.
func (h *Handler) Describe(ch chan<- *prom.Desc) {
	ch <- mSentDesc
	ch <- mRecvDesc
	ch <- mResponsesDesc
	ch <- mRequestsDesc

	h.mConnects.Describe(ch)
	h.mDisconnects.Describe(ch)
	h.mRoundTrip.Describe(ch)
	h.mClockDrift.Describe(ch)
}

// Collect implement prometheus.Collector.
func (h *Handler) Collect(ch chan<- prom.Metric) {
	h.mConnects.Collect(ch)
	h.mDisconnects.Collect(ch)
	h.mRoundTrip.Collect(ch)
	h.mClockDrift.Collect(ch)
}

// check interfaces
var (
	_ prom.Collector = (*Handler)(nil)
)
