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
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona/exporter_shared/helpers"
	"github.com/percona/pmm/api/agent"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type testServer struct {
	connect func(agent.Agent_ConnectServer) error
}

func (s *testServer) Connect(stream agent.Agent_ConnectServer) error {
	return s.connect(stream)
}

var _ agent.AgentServer = (*testServer)(nil)

func setup(t *testing.T, connect func(*Channel) error, expected ...error) (agent.Agent_ConnectClient, *grpc.ClientConn, func(*testing.T)) {
	// logrus.SetLevel(logrus.DebugLevel)

	t.Parallel()

	// start server with given connect handler
	var channel *Channel
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := grpc.NewServer()
	agent.RegisterAgentServer(server, &testServer{
		connect: func(stream agent.Agent_ConnectServer) error {
			channel = NewChannel(stream)
			return connect(channel)
		},
	})
	go func() {
		err = server.Serve(lis)
		require.NoError(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// make client and channel
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithWaitForHandshake(),
		grpc.WithInsecure(),
	}
	cc, err := grpc.DialContext(ctx, lis.Addr().String(), opts...)
	require.NoError(t, err, "failed to dial server")
	stream, err := agent.NewAgentClient(cc).Connect(ctx)
	require.NoError(t, err, "failed to create stream")

	teardown := func(t *testing.T) {
		err := channel.Wait()
		assert.Contains(t, expected, errors.Cause(err), "%+v", err)

		server.GracefulStop()
		cancel()
	}

	return stream, cc, teardown
}

func TestAgentRequest(t *testing.T) {
	const count = 50
	require.True(t, count > serverRequestsCap)

	var channel *Channel
	connect := func(ch *Channel) error { //nolint:unparam
		channel = ch // store to check metrics below

		for i := uint32(1); i <= count; i++ {
			msg := <-ch.Requests()
			require.NotNil(t, msg)
			assert.Equal(t, i, msg.Id)
			require.NotNil(t, msg.GetQanData())

			ch.SendResponse(&agent.ServerMessage{
				Id: i,
				Payload: &agent.ServerMessage_QanData{
					QanData: new(agent.QANDataResponse),
				},
			})
		}

		assert.Nil(t, <-ch.Requests())
		return nil
	}

	stream, _, teardown := setup(t, connect, io.EOF) // EOF = server exits from handler
	defer teardown(t)

	for i := uint32(1); i <= count; i++ {
		err := stream.Send(&agent.AgentMessage{
			Id: i,
			Payload: &agent.AgentMessage_QanData{
				QanData: new(agent.QANDataRequest),
			},
		})
		assert.NoError(t, err)

		msg, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, i, msg.Id)
		require.NotNil(t, msg.GetQanData())
	}

	err := stream.CloseSend()
	assert.NoError(t, err)

	// check metrics
	metrics := make([]prometheus.Metric, 0, 100)
	metricsCh := make(chan prometheus.Metric)
	go func() {
		channel.Collect(metricsCh)
		close(metricsCh)
	}()
	for m := range metricsCh {
		metrics = append(metrics, m)
	}
	expectedMetrics := strings.Split(strings.TrimSpace(`
# HELP pmm_managed_channel_messages_received_total A total number of received messages from pmm-agent.
# TYPE pmm_managed_channel_messages_received_total counter
pmm_managed_channel_messages_received_total 50
# HELP pmm_managed_channel_messages_sent_total A total number of sent messages to pmm-agent.
# TYPE pmm_managed_channel_messages_sent_total counter
pmm_managed_channel_messages_sent_total 50
`), "\n")
	assert.Equal(t, expectedMetrics, helpers.Format(metrics))

	// check that descriptions match metrics: same number, same order
	descCh := make(chan *prometheus.Desc)
	go func() {
		channel.Describe(descCh)
		close(descCh)
	}()
	var i int
	for d := range descCh {
		assert.Equal(t, metrics[i].Desc(), d)
		i++
	}
	assert.Len(t, metrics, i)
}

func TestServerRequest(t *testing.T) {
	const count = 50
	require.True(t, count > serverRequestsCap)

	connect := func(ch *Channel) error { //nolint:unparam
		for i := uint32(1); i <= count; i++ {
			res := ch.SendRequest(&agent.ServerMessage_Ping{
				Ping: new(agent.PingRequest),
			})
			ping := res.(*agent.AgentMessage_Ping)
			ts, err := ptypes.Timestamp(ping.Ping.CurrentTime)
			assert.NoError(t, err)
			assert.InDelta(t, time.Now().Unix(), ts.Unix(), 1)
		}

		assert.Nil(t, <-ch.Requests())
		return nil
	}

	stream, _, teardown := setup(t, connect, io.EOF) // EOF = server exits from handler
	defer teardown(t)

	for i := uint32(1); i <= count; i++ {
		msg, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, i, msg.Id)
		require.NotNil(t, msg.GetPing())

		err = stream.Send(&agent.AgentMessage{
			Id: i,
			Payload: &agent.AgentMessage_Ping{
				Ping: &agent.PingResponse{
					CurrentTime: ptypes.TimestampNow(),
				},
			},
		})
		assert.NoError(t, err)
	}

	err := stream.CloseSend()
	assert.NoError(t, err)
}

/*
func TestServerExitsWithGRPCError(t *testing.T) {
	errUnimplemented := status.Error(codes.Unimplemented, "Test error")
	connect := func(stream agent.Agent_ConnectServer) error { //nolint:unparam
		msg, err := stream.Recv()
		require.NoError(t, err)
		assert.EqualValues(t, 1, msg.Id)
		require.NotNil(t, msg.GetQanData())

		return errUnimplemented
	}

	channel, _, teardown := setup(t, connect, errUnimplemented)
	defer teardown(t)

	resp := channel.SendRequest(&agent.AgentMessage_QanData{
		QanData: new(agent.QANDataRequest),
	})
	assert.Nil(t, resp)
}

func TestServerExitsWithUnknownError(t *testing.T) {
	connect := func(stream agent.Agent_ConnectServer) error { //nolint:unparam
		msg, err := stream.Recv()
		require.NoError(t, err)
		assert.EqualValues(t, 1, msg.Id)
		require.NotNil(t, msg.GetQanData())

		return io.EOF // any error without GRPCStatus() method
	}

	channel, _, teardown := setup(t, connect, status.Error(codes.Unknown, "EOF"))
	defer teardown(t)

	resp := channel.SendRequest(&agent.AgentMessage_QanData{
		QanData: new(agent.QANDataRequest),
	})
	assert.Nil(t, resp)
}

func TestAgentClosesStream(t *testing.T) {
	connect := func(stream agent.Agent_ConnectServer) error { //nolint:unparam
		err := stream.Send(&agent.ServerMessage{
			Id: 1,
			Payload: &agent.ServerMessage_Ping{
				Ping: new(agent.PingRequest),
			},
		})
		assert.NoError(t, err)

		msg, err := stream.Recv()
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, msg)

		return nil
	}

	channel, _, teardown := setup(t, connect, io.EOF)
	defer teardown(t)

	req := <-channel.Requests()
	ping := req.GetPing()
	assert.NotNil(t, ping)

	err := channel.s.CloseSend()
	assert.NoError(t, err)
}

func TestAgentClosesConnection(t *testing.T) {
	connect := func(stream agent.Agent_ConnectServer) error { //nolint:unparam
		err := stream.Send(&agent.ServerMessage{
			Id: 1,
			Payload: &agent.ServerMessage_Ping{
				Ping: new(agent.PingRequest),
			},
		})
		assert.NoError(t, err)

		msg, err := stream.Recv()
		assert.Equal(t, status.Error(codes.Canceled, "context canceled"), err)
		assert.Nil(t, msg)

		return nil
	}

	// gRPC library has a race in that case, so we can get two errors
	errClientConnClosing := status.Error(codes.Canceled, "grpc: the client connection is closing") // == grpc.ErrClientConnClosing
	errConnClosing := status.Error(codes.Unavailable, "transport is closing")
	channel, cc, teardown := setup(t, connect, errClientConnClosing, errConnClosing)
	defer teardown(t)

	req := <-channel.Requests()
	ping := req.GetPing()
	assert.NotNil(t, ping)

	err := cc.Close()
	assert.NoError(t, err)
}

func TestUnexpectedMessageFromServer(t *testing.T) {
	connect := func(stream agent.Agent_ConnectServer) error { //nolint:unparam
		// this message triggers "no subscriber for ID" error
		err := stream.Send(&agent.ServerMessage{
			Id: 111,
			Payload: &agent.ServerMessage_QanData{
				QanData: new(agent.QANDataResponse),
			},
		})
		assert.NoError(t, err)

		// this message should not trigger new error
		err = stream.Send(&agent.ServerMessage{
			Id: 222,
			Payload: &agent.ServerMessage_QanData{
				QanData: new(agent.QANDataResponse),
			},
		})
		assert.NoError(t, err)

		return nil
	}

	channel, _, teardown := setup(t, connect, fmt.Errorf("no subscriber for ID 111"))
	defer teardown(t)

	// after receiving unexpected response, channel is closed
	resp := channel.SendRequest(&agent.AgentMessage_QanData{
		QanData: new(agent.QANDataRequest),
	})
	assert.Nil(t, resp)
	msg := <-channel.Requests()
	assert.Nil(t, msg)

	// future requests are ignored
	resp = channel.SendRequest(&agent.AgentMessage_QanData{
		QanData: new(agent.QANDataRequest),
	})
	assert.Nil(t, resp)
	msg = <-channel.Requests()
	assert.Nil(t, msg)
}
*/
