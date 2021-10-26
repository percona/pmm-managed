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

// OktaSSODetails stores everything we need to issue access_token from Okta SSO API.
// It is intended to have only one row in this table as PMM can be connected to Portal only once.
//reform:okta_sso_details
type OktaSSODetails struct {
	ClientID     string `reform:"client_id"`
	ClientSecret string `reform:"client_secret"`
	IssuerURL    string `reform:"issuer_url"`
	Scope        string `reform:"scope"`

	CreatedAt time.Time `reform:"created_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (s *OktaSSODetails) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (s *OktaSSODetails) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	return nil
}

// check interfaces.
var (
	_ reform.BeforeInserter = (*OktaSSODetails)(nil)
	_ reform.AfterFinder    = (*OktaSSODetails)(nil)
)
