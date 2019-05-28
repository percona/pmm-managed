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

package action

// PtSummary represents pt-summary domain model.
type PtSummary struct {
	ID         string
	PMMAgentID string
	NodeID     string

	Args []string
}

// PtMySQLSummary represents pt-mysql-summary domain model.
type PtMySQLSummary struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Args []string
}

// MySQLExplain represents mysql-explain domain model.
type MySQLExplain struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}

// MySQLExplainJSON represents mysql-explain-json domain model.
type MySQLExplainJSON struct {
	ID         string
	PMMAgentID string
	ServiceID  string

	Dsn   string
	Query string
}
