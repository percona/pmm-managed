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
	ID         string
	ServiceID  string
	PMMAgentID string
	Dsn        string
	Query      string

	q *reform.Querier
}

func NewMySQLExplain(q *reform.Querier) *MySQLExplain {
	return &MySQLExplain{
		ID: getNewActionID(),
		q:  q,
	}
}

func (exp *MySQLExplain) Prepare(serviceID, pmmAgentID, query string) error {
	var err error
	exp.Query = query
	exp.ServiceID = serviceID
	exp.PMMAgentID = pmmAgentID

	exp.PMMAgentID, err = findPmmAgentIDByServiceID(exp.q, exp.PMMAgentID, exp.ServiceID)
	if err != nil {
		return err
	}
	// TODO: Add DSN string: findAgentForService(MYSQLD_EXPORTER, req.ServiceID)

	return nil
}
