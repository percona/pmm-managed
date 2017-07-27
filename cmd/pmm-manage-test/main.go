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

package main

import (
	"context"
	"flag"
	"io"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/Percona-Lab/pmm-managed/api"
)

var (
	gRPCAddrF = flag.String("grpc-addr", "127.0.0.1:7771", "gRPC server address")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*gRPCAddrF, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	c := api.NewBaseClient(conn)
	stream, err := c.Ping(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	// start pinger
	go func() {
		for {
			req := &api.BasePingRequest{
				Type:   api.PingType_PING,
				Cookie: uint64(time.Now().UnixNano()),
			}
			if pingErr := stream.Send(req); pingErr != nil {
				logrus.Error(pingErr)
				return
			}
			time.Sleep(time.Duration(rand.Intn(int(time.Second))))
		}
	}()

	// ponger
	for {
		var resp *api.BasePingResponse
		resp, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		switch resp.Type {
		case api.PingType_PING:
			logrus.Printf("Received ping: %d", resp.Cookie)
			req := &api.BasePingRequest{
				Type:   api.PingType_PONG,
				Cookie: resp.Cookie,
			}
			if err = stream.Send(req); err != nil {
				return
			}
		case api.PingType_PONG:
			d := time.Since(time.Unix(0, int64(resp.Cookie)))
			logrus.Printf("Received pong: %d (latency %v)", resp.Cookie, d)
		}
	}
}
