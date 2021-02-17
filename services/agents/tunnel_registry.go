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
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm/api/agentpb"
	"gopkg.in/reform.v1"
)

// type tunnelRegistry struct {
// }

// func newTunnelRegistry() *tunnelRegistry {
// 	return &tunnelRegistry{}
// }

// getTunnels returns SetStateRequest's tunnels for given pmm-agent ID.
func getTunnels(q *reform.Querier, pmmAgentID string) (map[string]*agentpb.SetStateRequest_Tunnel, error) {
	tunnels, err := models.FindTunnels(q, pmmAgentID)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*agentpb.SetStateRequest_Tunnel, len(tunnels))
	for _, t := range tunnels {
		var req agentpb.SetStateRequest_Tunnel
		switch pmmAgentID {
		case t.ListenAgentID:
			req.ListenPort = uint32(t.ListenPort)
		case t.ConnectAgentID:
			req.ConnectPort = uint32(t.ConnectPort)
		}

		res[t.TunnelID] = &req
	}

	return res, nil
}
