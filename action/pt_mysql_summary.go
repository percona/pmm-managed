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

type PtMySQLSummary struct {
	Id         string
	ServiceID  string
	PMMAgentID string
	Args       []string

	q *reform.Querier
}

func NewPtMySQLSummary(serviceId, pmmAgentID string, args []string, q *reform.Querier) *PtMySQLSummary {
	return &PtMySQLSummary{
		Id:         getNewActionID(),
		ServiceID:  serviceId,
		PMMAgentID: pmmAgentID,
		Args:       args,

		q: q,
	}
}

func (pms *PtMySQLSummary) Prepare() error {
	var err error
	pms.PMMAgentID, err = findPmmAgentIDByServiceID(pms.q, pms.PMMAgentID, pms.ServiceID)
	if err != nil {
		return err
	}
	return nil
}
