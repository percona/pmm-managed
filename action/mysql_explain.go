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

type MySQLExplain struct {
	Id         string
	ServiceID  string
	PMMAgentID string
	Dsn        string
	Query      string

	q *reform.Querier
}

// TODO: Add DSN string: findAgentForService(MYSQLD_EXPORTER, req.ServiceID)
func NewMySQLExplain(serviceId, pmmAgentID string, query string, q *reform.Querier) *MySQLExplain {
	return &MySQLExplain{
		Id:         getNewActionID(),
		ServiceID:  serviceId,
		PMMAgentID: pmmAgentID,
		Query:      query,

		q: q,
	}
}

func (exp *MySQLExplain) Prepare() error {
	var err error
	exp.PMMAgentID, err = findPmmAgentIDByServiceID(exp.q, exp.PMMAgentID, exp.ServiceID)
	if err != nil {
		return err
	}
	return nil
}
