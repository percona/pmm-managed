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

type PtSummary struct {
	ID         string
	NodeID     string
	PMMAgentID string
	Args       []string
}

type PtMySQLSummary struct {
	ID         string
	ServiceID  string
	PMMAgentID string
	Args       []string
}

type MySQLExplain struct {
	ID         string
	ServiceID  string
	PMMAgentID string
	Dsn        string
	Query      string
}

type MySQLExplainJSON struct {
	ID         string
	ServiceID  string
	PMMAgentID string
	Dsn        string
	Query      string
}

// Result describes an PMM Action result which is storing in ActionsResult storage.
type Result struct {
	ID         string
	PmmAgentID string
	Done       bool
	Error      string
	Output     string
}
