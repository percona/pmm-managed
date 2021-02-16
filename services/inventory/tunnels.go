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

package inventory

import (
	"context"

	"github.com/percona/pmm/api/inventorypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TunnelsService works with inventory API Tunnels.
type TunnelsService struct{}

// NewTunnelsService returns Inventory API handler for managing Tunnels.
func NewTunnelsService() inventorypb.TunnelsServer {
	return &TunnelsService{}
}

// ListTunnels returns a list of all Tunnels.
func (s *TunnelsService) ListTunnels(ctx context.Context, req *inventorypb.ListTunnelsRequest) (*inventorypb.ListTunnelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTunnels not implemented")
}

// AddTunnel adds a Tunnel.
func (s *TunnelsService) AddTunnel(ctx context.Context, req *inventorypb.AddTunnelRequest) (*inventorypb.AddTunnelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTunnel not implemented")
}

// RemoveTunnel removes a Tunnel.
func (s *TunnelsService) RemoveTunnel(ctx context.Context, req *inventorypb.RemoveTunnelRequest) (*inventorypb.RemoveTunnelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTunnel not implemented")
}
