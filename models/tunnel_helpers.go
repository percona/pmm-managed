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

package models

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

func checkUniqueTunnelID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Tunnel ID")
	}

	tunnel := &Tunnel{TunnelID: id}
	switch err := q.Reload(tunnel); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Tunnel with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// FindTunnels returns Tunnels for given pmm-agent, or all, if pmmAgentID is empty.
func FindTunnels(q *reform.Querier, pmmAgentID string) ([]*Tunnel, error) {
	var args []interface{}
	tail := "ORDER BY tunnel_id"
	if pmmAgentID != "" {
		// TODO check that agent exist
		args = []interface{}{pmmAgentID, pmmAgentID}
		tail = "WHERE listen_agent_id = ? OR connect_agent_id = ? " + tail
	}

	structs, err := q.SelectAllFrom(TunnelTable, tail, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tunnels := make([]*Tunnel, len(structs))
	for i, s := range structs {
		tunnels[i] = s.(*Tunnel)
	}

	return tunnels, nil
}

// CreateTunnelParams TODO.
type CreateTunnelParams struct {
	ListenAgentID  string
	ListenPort     uint16
	ConnectAgentID string
	ConnectPort    uint16
}

// CreateTunnel creates Tunnel.
func CreateTunnel(q *reform.Querier, params *CreateTunnelParams) (*Tunnel, error) {
	id := "/tunnel_id/" + uuid.New().String()
	if err := checkUniqueTunnelID(q, id); err != nil {
		return nil, err
	}

	// TODO check that agents exist
	// TODO check that ports > 0

	row := &Tunnel{
		TunnelID:       id,
		ListenAgentID:  params.ListenAgentID,
		ListenPort:     params.ListenPort,
		ConnectAgentID: params.ConnectAgentID,
		ConnectPort:    params.ConnectPort,
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	return row, nil
}

// RemoveTunnel removes Tunnel by ID.
func RemoveTunnel(q *reform.Querier, id string, mode RemoveMode) (*Tunnel, error) {
	// TODO find agents
	// TODO cascade delete

	t := &Tunnel{TunnelID: id}
	if err := q.Delete(t); err != nil {
		return nil, errors.Wrap(err, "failed to delete Tunnel")
	}

	return t, nil
}
