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
	"time"

	"gopkg.in/reform.v1"
)

//go:generate reform

type TunnelType string

const (
	ConnectTunnelType TunnelType = "connect"
	ListenTunnelType  TunnelType = "listen"
)

//reform:tunnels
type Tunnel struct {
	TunnelID    string     `reform:"tunnel_id,pk"`
	TunnelType  TunnelType `reform:"tunnel_type"`
	PMMAgentID  string     `reform:"pmm_agent_id"`
	ConnectPort uint16     `reform:"connect_port"`
	ListenPort  uint16     `reform:"listen_port"`
	CreatedAt   time.Time  `reform:"created_at"`
	UpdatedAt   time.Time  `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (s *Tunnel) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (s *Tunnel) BeforeUpdate() error {
	s.UpdatedAt = Now()
	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (s *Tunnel) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Tunnel)(nil)
	_ reform.BeforeUpdater  = (*Tunnel)(nil)
	_ reform.AfterFinder    = (*Tunnel)(nil)
)
