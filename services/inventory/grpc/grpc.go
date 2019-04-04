// Copyright (C) 2019 Percona LLC
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

// Package inventory contains inventory business logic: Nodes, Services, Agents.
package grpc

import (
	"context"
)

//go:generate mockery -name=registry -case=snake -inpkg -testonly

// registry is a subset of methods of agents.Registry used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type registry interface {
	SendSetStateRequest(ctx context.Context, pmmAgentID string)
	IsConnected(pmmAgentID string) bool
	Kick(ctx context.Context, pmmAgentID string)
}
