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
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

type DSNResolver struct {
	db *reform.DB
}

func NewDSNResolver(db *reform.DB) *DSNResolver {
	return &DSNResolver{
		db: db,
	}
}

func (r *DSNResolver) ResolveDSNByServiceID(serviceID string) (string, error) {
	var result string

	err := r.db.InTransaction(func(t *reform.TX) error {
		svc, err := FindServiceByID(r.db.Querier, serviceID)
		if err != nil {
			return err
		}

		pmmAgents, err := FindPMMAgentsForService(r.db.Querier, serviceID)
		if err != nil {
			return err
		}

		if len(pmmAgents) != 1 {
			return errors.New("couldn't resolve pmm-agent")
		}

		pmmAgentID := pointer.GetString(pmmAgents[0].PMMAgentID)
		var agentType AgentType
		switch svc.ServiceType {
		case MySQLServiceType:
			agentType = MySQLdExporterType
		case MongoDBServiceType:
			agentType = MongoDBExporterType
		case PostgreSQLServiceType:
			agentType = PostgresExporterType
		default:
			return errors.New("couldn't resolve service type")
		}

		agents, err := FindAgentsByPmmAgentIDAndAgentType(r.db.Querier, pmmAgentID, agentType)
		if err != nil {
			return err
		}

		if len(agents) != 1 {
			return errors.New("couldn't resolve agent")
		}

		resolvedAgent := agents[0]

		switch svc.ServiceType {
		case MySQLServiceType:
			cfg := mysql.NewConfig()
			cfg.Addr = fmt.Sprintf("%s:%d", pointer.GetString(svc.Address), pointer.GetUint16(svc.Port))
			cfg.User = pointer.GetString(resolvedAgent.Username)
			cfg.Passwd = pointer.GetString(resolvedAgent.Password)
			result = cfg.FormatDSN()
			return nil

		case MongoDBServiceType:
			// TODO: Implement DSN resolver for MongoDB
		case PostgreSQLServiceType:
			// TODO: Implement DSN resolver for MongoDB
		}

		return errors.New("couldn't resolve service type")
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
