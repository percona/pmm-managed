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
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// TunnelsService works with inventory API Tunnels.
type TunnelsService struct {
	db *reform.DB
	r  agentsRegistry
}

// NewTunnelsService returns Inventory API handler for managing Tunnels.
func NewTunnelsService(db *reform.DB, r agentsRegistry) inventorypb.TunnelsServer {
	return &TunnelsService{
		db: db,
		r:  r,
	}
}

// ListTunnels returns a list of all Tunnels.
func (s *TunnelsService) ListTunnels(ctx context.Context, req *inventorypb.ListTunnelsRequest) (*inventorypb.ListTunnelsResponse, error) {
	tunnels, err := models.FindTunnels(s.db.Querier, req.AgentId)
	if err != nil {
		return nil, err
	}

	res := &inventorypb.ListTunnelsResponse{
		Tunnel: make([]*inventorypb.Tunnel, len(tunnels)),
	}
	for i, t := range tunnels {
		res.Tunnel[i] = &inventorypb.Tunnel{
			TunnelId:       t.TunnelID,
			ListenAgentId:  t.ListenAgentID,
			ListenPort:     uint32(t.ListenPort),
			ConnectAgentId: t.ConnectAgentID,
			ConnectPort:    uint32(t.ConnectPort),
		}
	}
	return res, nil
}

// AddTunnel adds a Tunnel.
func (s *TunnelsService) AddTunnel(ctx context.Context, req *inventorypb.AddTunnelRequest) (*inventorypb.AddTunnelResponse, error) {
	params := &models.CreateTunnelParams{
		ListenAgentID:  req.ListenAgentId,
		ListenPort:     uint16(req.ListenPort),
		ConnectAgentID: req.ConnectAgentId,
		ConnectPort:    uint16(req.ConnectPort),
	}

	t, err := models.CreateTunnel(s.db.Querier, params)
	if err != nil {
		return nil, err
	}

	s.r.RequestStateUpdate(ctx, req.ListenAgentId)
	s.r.RequestStateUpdate(ctx, req.ConnectAgentId)

	return &inventorypb.AddTunnelResponse{
		TunnelId: t.TunnelID,
	}, nil
}

// RemoveTunnel removes a Tunnel.
func (s *TunnelsService) RemoveTunnel(ctx context.Context, req *inventorypb.RemoveTunnelRequest) (*inventorypb.RemoveTunnelResponse, error) {
	t, err := models.RemoveTunnel(s.db.Querier, req.TunnelId, models.RemoveCascade)
	if err != nil {
		return nil, err
	}

	s.r.RequestStateUpdate(ctx, t.ListenAgentID)
	s.r.RequestStateUpdate(ctx, t.ConnectAgentID)

	return new(inventorypb.RemoveTunnelResponse), nil
}
