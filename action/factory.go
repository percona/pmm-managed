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

type Factory struct {
	db *reform.DB
	s  *InMemoryStorage
}

// NewFactory creates new actions Factory.
func NewFactory(db *reform.DB, s *InMemoryStorage) *Factory {
	return &Factory{
		db: db,
		s:  s,
	}
}

func (f *Factory) NewPTSummary(ctx context.Context, nodeID, pmmAgentID string) (*PtSummary, error) {
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
		return nil, errors.Wrap(err, "couldn't create action")
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	err = f.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	return a, nil
}

func (f *Factory) NewPTMySQLSummary(ctx context.Context, serviceID, pmmAgentID string) (*PtMySQLSummary, error) {
	a := &PtMySQLSummary{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Args:       []string{},
	}

	agents, err := models.FindPMMAgentsForService(f.db.Querier, a.ServiceID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	err = f.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	return a, nil
}

//nolint:dupl
func (f *Factory) NewMySQLExplain(ctx context.Context, serviceID, pmmAgentID, query string) (*MySQLExplain, error) {
	a := &MySQLExplain{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Query:      query,
	}

	agents, err := models.FindPMMAgentsForService(f.db.Querier, a.ServiceID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	a.Dsn, err = f.resolveDSNByServiceID(a.ServiceID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	err = f.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	return a, nil
}

//nolint:dupl
func (f *Factory) NewMySQLExplainJSON(ctx context.Context, serviceID, pmmAgentID, query string) (*MySQLExplainJSON, error) {
	a := &MySQLExplainJSON{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Query:      query,
	}

	agents, err := models.FindPMMAgentsForService(f.db.Querier, a.ServiceID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	a.PMMAgentID, err = findPmmAgentIDToRunAction(a.PMMAgentID, agents)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	a.Dsn, err = f.resolveDSNByServiceID(a.ServiceID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	err = f.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create action")
	}

	return a, nil
}

func (f *Factory) resolveDSNByServiceID(serviceID string) (string, error) {
	var result string

	err := f.db.InTransaction(func(tx *reform.TX) error {
		svc, err := models.FindServiceByID(tx.Querier, serviceID)
		if err != nil {
			return errors.Wrap(err, "couldn't resolve dsn")
		}

		pmmAgents, err := models.FindPMMAgentsForService(tx.Querier, serviceID)
		if err != nil {
			return errors.Wrap(err, "couldn't resolve dsn")
		}

		if len(pmmAgents) != 1 {
			return errors.New("couldn't resolve dsn, as there should be only one pmm-agent")
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
			return errors.New("couldn't resolve dsn, as service is unsupported")
		}

		exporters, err := models.FindAgentsByPmmAgentIDAndAgentType(tx.Querier, pmmAgentID, agentType)
		if err != nil {
			return errors.Wrap(err, "couldn't resolve dsn")
		}

		if len(exporters) != 1 {
			return errors.New("couldn't resolve dsn, as there should be only one exporter")
		}

		switch svc.ServiceType {
		case models.MySQLServiceType:
			result = models.DSNforMySQL(svc, exporters[0])

		case models.MongoDBServiceType:
			result = models.DSNforMongoDB(svc, exporters[0])

		case models.PostgreSQLServiceType:
			result = models.DSNforPostgreSQL(svc, exporters[0])
		}

		return nil
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
		return "", errors.New("couldn't find pmm-agent-id to run action")
	}

	// check that explicit agent id is correct
	for _, a := range agents {
		if a.AgentID == pmmAgentID {
			return a.AgentID, nil
		}
	}
	return "", errors.New("couldn't find pmm-agent-id to run action")
}

// nolint: unused
func getUUID() string {
	return "/action_id/" + uuid.New().String()
}
