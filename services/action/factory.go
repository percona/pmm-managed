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
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

var (
	errPmmAgentIDNotFound = errors.New("can't detect pmm_agent_id")
)

type factory struct {
	db *reform.DB
}

// NewFactory creates new actions factory.
func NewFactory(db *reform.DB) Factory {
	return &factory{
		db: db,
	}
}

func (f *factory) NewPTSummary(ctx context.Context, nodeID, pmmAgentID string) (*PtSummary, error) {
	a := &PtSummary{
		ID:                 getUUID(),
		NodeID:             nodeID,
		PMMAgentID:         pmmAgentID,
		SummarizeMounts:    true,
		SummarizeNetwork:   true,
		SummarizeProcesses: true,
		Sleep:              5,
		Help:               false,
	}

	agents, err := models.FindPMMAgentsForNode(f.db.Querier, a.NodeID)
	if err != nil {
		return nil, err
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (f *factory) NewPTMySQLSummary(ctx context.Context, serviceID, pmmAgentID string) (*PtMySQLSummary, error) {
	a := &PtMySQLSummary{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Args:       []string{},
	}

	agents, err := models.FindPMMAgentsForService(f.db.Querier, a.ServiceID)
	if err != nil {
		return nil, err
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (f *factory) NewMySQLExplain(ctx context.Context, serviceID, pmmAgentID, query string) (*MySQLExplain, error) {
	a := &MySQLExplain{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Query:      query,
	}

	agents, err := models.FindPMMAgentsForService(f.db.Querier, a.ServiceID)
	if err != nil {
		return nil, err
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, err
	}

	a.Dsn, err = f.resolveDSNByServiceID(a.ServiceID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (f *factory) NewMySQLExplainJSON(ctx context.Context, serviceID, pmmAgentID, query string) (*MySQLExplainJSON, error) {
	a := &MySQLExplainJSON{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Query:      query,
	}

	agents, err := models.FindPMMAgentsForService(f.db.Querier, a.ServiceID)
	if err != nil {
		return nil, err
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, err
	}

	a.Dsn, err = f.resolveDSNByServiceID(a.ServiceID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (f *factory) resolveDSNByServiceID(serviceID string) (string, error) {
	var result string

	err := f.db.InTransaction(func(tx *reform.TX) error {
		svc, err := models.FindServiceByID(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		pmmAgents, err := models.FindPMMAgentsForService(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		if len(pmmAgents) != 1 {
			return errors.New("couldn't resolve pmm-agent")
		}

		pmmAgentID := pmmAgents[0].AgentID
		var agentType models.AgentType
		switch svc.ServiceType {
		case models.MySQLServiceType:
			agentType = models.MySQLdExporterType
		case models.MongoDBServiceType:
			agentType = models.MongoDBExporterType
		case models.PostgreSQLServiceType:
			agentType = models.PostgresExporterType
		default:
			return errors.New("couldn't resolve service type")
		}

		agents, err := models.FindAgentsByPmmAgentIDAndAgentType(tx.Querier, pmmAgentID, agentType)
		if err != nil {
			return err
		}

		if len(agents) != 1 {
			return errors.New("couldn't resolve agent")
		}

		resolvedAgent := agents[0]

		switch svc.ServiceType {
		case models.MySQLServiceType:
			result = models.DSNforMySQL(svc, resolvedAgent)

		case models.MongoDBServiceType:
			result = models.DSNforMongoDB(svc, resolvedAgent)

		case models.PostgreSQLServiceType:
			result = models.DSNforPostgreSQL(svc, resolvedAgent)
		}

		return errors.New("couldn't resolve service type")
	})
	if err != nil {
		return result, err
	}

	return result, nil
}

func findPmmAgentIDToRunAction(pmmAgentID string, agents []*models.Agent) (string, error) {
	// no explicit ID is given, and there is only one
	if pmmAgentID == "" && len(agents) == 1 {
		return agents[0].AgentID, nil
	}

	// no explicit ID is given, and there are zero or several
	if pmmAgentID == "" {
		return "", errPmmAgentIDNotFound
	}

	// check that explicit agent id is correct
	for _, a := range agents {
		if a.AgentID == pmmAgentID {
			return a.AgentID, nil
		}
	}
	return "", errPmmAgentIDNotFound
}

// nolint: unused
func getUUID() string {
	return "/action_id/" + uuid.New().String()
}
