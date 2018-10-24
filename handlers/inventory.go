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

package handlers

import (
	"context"

	"github.com/Percona-Lab/pmm-api/inventory"

	"github.com/percona/pmm-managed/services/agents"
)

type InventoryServer struct {
	Store  *agents.Store
	Agents map[uint32]*agents.Conn
}

func (s *InventoryServer) AddBareMetal(ctx context.Context, req *inventory.AddBareMetalNodeRequest) (*inventory.AddBareMetalNodeResponse, error) {
	node := s.Store.AddBareMetalNode(req)
	return &inventory.AddBareMetalNodeResponse{
		Node: node,
	}, nil
}

func (s *InventoryServer) AddMySQLdExporter(ctx context.Context, req *inventory.AddMySQLdExporterRequest) (*inventory.AddMySQLdExporterResponse, error) {
	exporter := s.Store.AddMySQLdExporter(req)
	return &inventory.AddMySQLdExporterResponse{
		Agent: exporter,
	}, nil
}

// check interfaces
var (
	_ inventory.NodesServer    = (*InventoryServer)(nil)
	_ inventory.ServicesServer = (*InventoryServer)(nil)
	_ inventory.AgentsServer   = (*InventoryServer)(nil)
)
