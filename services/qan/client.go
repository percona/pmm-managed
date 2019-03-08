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

package qan

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	qanAPI "github.com/percona/pmm/api/qan"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Client struct {
	c qanAPI.AgentClient
	l *logrus.Entry
}

func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{
		c: qanAPI.NewAgentClient(cc),
		l: logrus.WithField("component", "qan"),
	}
}

func (c *Client) TODO(ctx context.Context, m *any.Any) {
	res, err := c.c.TODO(ctx, &qanAPI.AgentMessageTODO{Data: m})
	if err != nil {
		c.l.Error(err)
	}
	c.l.Debug(res)
}
