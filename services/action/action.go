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

import (
	"context"

	"github.com/google/uuid"
)

type Runner interface {
	StartMySQLExplainJSONAction(context.Context, *MySQLExplainJSON) error
	StartMySQLExplainAction(context.Context, *MySQLExplain) error
	StartPTMySQLSummaryAction(context.Context, *PtMySQLSummary) error
	StartPTSummaryAction(context.Context, *PtSummary) error
	StopAction(ctx context.Context, actionID, pmmAgentID string) error
}

type Storage interface {
	Store(context.Context, *Result) error
	Update(context.Context, *Result) error
	Load(context.Context, string) (*Result, error)
}

type PMMAgentIDResolver interface {
	ResolvePMMAgentIDByServiceID(serviceID, pmmAgentID string) (string, error)
	ResolvePMMAgentIDByNodeID(nodeID, pmmAgentID string) (string, error)
}

type DSNResolver interface {
	ResolveDSNByServiceID(serviceID string) (string, error)
}

type Service interface {
	GetActionResult(ctx context.Context, actionID string) (*Result, error)
	StartPTSummaryAction(ctx context.Context, pmmAgentID, nodeID string) (*PtSummary, error)
	StartPTMySQLSummaryAction(ctx context.Context, pmmAgentID, serviceID string) (*PtMySQLSummary, error)
	StartMySQLExplainAction(ctx context.Context, pmmAgentID, serviceID, query string) (*MySQLExplain, error)
	StartMySQLExplainJSONAction(ctx context.Context, pmmAgentID, serviceID, query string) (*MySQLExplainJSON, error)
	CancelAction(ctx context.Context, actionID string) error
}

func getUUID() string {
	return "/action_id/" + uuid.New().String()
}
