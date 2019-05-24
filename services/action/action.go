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

// Package action provides models and service for interacting with Actions.
// It's application layer package, so it contains only abstract application logic
// and separated with infrastructure and data layers through interfaces.
// Infrastructure and data layers packages should implement those interfaces.
package action

import (
	"context"
)

// Runner provides service that can Run actions of different types.
// Action parameters should be prepared before action will be started.
type Runner interface {
	StartMySQLExplainJSONAction(context.Context, *MySQLExplainJSON) error
	StartMySQLExplainAction(context.Context, *MySQLExplain) error
	StartPTMySQLSummaryAction(context.Context, *PtMySQLSummary) error
	StartPTSummaryAction(context.Context, *PtSummary) error
	StopAction(ctx context.Context, actionID, pmmAgentID string) error
}

// Storage provides persistent storage service for saving and loading action results.
type Storage interface {
	Store(context.Context, *Result) error
	Update(context.Context, *Result) error
	Load(context.Context, string) (*Result, error)
}

// Factory provides factory service to create different action types.
// After creation by this factory actions is ready to use and to run by Runner service.
type Factory interface {
	NewMySQLExplain(ctx context.Context, serviceID, pmmAgentID, query string) (*MySQLExplain, error)
	NewMySQLExplainJSON(ctx context.Context, serviceID, pmmAgentID, query string) (*MySQLExplainJSON, error)
	NewPTMySQLSummary(ctx context.Context, serviceID, pmmAgentID string) (*PtMySQLSummary, error)
	NewPTSummary(ctx context.Context, nodeID, pmmAgentID string) (*PtSummary, error)
}
