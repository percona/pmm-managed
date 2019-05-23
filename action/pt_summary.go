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

import "gopkg.in/reform.v1"

type PtSummary struct {
	ID         string
	NodeID     string
	PMMAgentID string
	Args       []string

	q *reform.Querier
}

func NewPtSummary(q *reform.Querier) *PtSummary {
	return &PtSummary{
		ID: getNewActionID(),
		q:  q,
	}
}

func (ps *PtSummary) Prepare(nodeID, pmmAgentID string, args []string) error {
	var err error
	ps.NodeID = nodeID
	ps.PMMAgentID = pmmAgentID
	ps.Args = args

	ps.PMMAgentID, err = findPmmAgentIDByNodeID(ps.q, ps.PMMAgentID, ps.NodeID)
	if err != nil {
		return err
	}
	return nil
}
