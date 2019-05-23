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

	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

type PtMySQLSummary struct {
	ID         string
	ServiceID  string
	PMMAgentID string
	Args       []string

	q *reform.Querier
	s ptMySQLSummaryStarter
}

func NewPtMySQLSummary(q *reform.Querier, s ptMySQLSummaryStarter) *PtMySQLSummary {
	return &PtMySQLSummary{
		ID: createActionID(),
		q:  q,
		s:  s,
	}
}

func (pms *PtMySQLSummary) Prepare(serviceID, pmmAgentID string, args []string) error {
	var err error
	pms.ServiceID = serviceID
	pms.PMMAgentID = pmmAgentID
	pms.Args = args

	agents, err := models.FindPMMAgentsForService(pms.q, pms.ServiceID)
	if err != nil {
		return err
	}

	pms.PMMAgentID, err = validatePMMAgentID(pms.PMMAgentID, agents)
	if err != nil {
		return err
	}
	return nil
}

func (pms *PtMySQLSummary) Start(ctx context.Context) error {
	return pms.s.StartPTMySQLSummaryAction(ctx, pms)
}
