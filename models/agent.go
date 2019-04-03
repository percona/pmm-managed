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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// AgentFindByID finds agent by ID.
func AgentFindByID(q *reform.Querier, id string) (*Agent, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Agent ID.")
	}

	row := &Agent{AgentID: id}
	switch err := q.Reload(row); err {
	case nil:
		return row, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Agent with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

func agentCheckUniqueID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Agent ID")
	}

	row := &Agent{AgentID: id}
	switch err := q.Reload(row); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Agent with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func agentNewID(q *reform.Querier) (string, error) {
	id := "/agent_id/" + uuid.New().String()
	if err := agentCheckUniqueID(q, id); err != nil {
		return id, err
	}
	return id, nil
}

// AgentRemove removes agent by ID.
func AgentRemove(q *reform.Querier, id string) (*Agent, error) {
	row, err := AgentFindByID(q, id)
	if err != nil {
		return row, err
	}

	if _, err = q.DeleteFrom(AgentServiceView, "WHERE agent_id = "+q.Placeholder(1), id); err != nil { //nolint:gosec
		return row, errors.WithStack(err)
	}
	if _, err = q.DeleteFrom(AgentNodeView, "WHERE agent_id = "+q.Placeholder(1), id); err != nil { //nolint:gosec
		return row, errors.WithStack(err)
	}

	if err = q.Delete(row); err != nil {
		return row, errors.WithStack(err)
	}

	return row, nil
}

// AgentFindAll finds all agents.
func AgentFindAll(q *reform.Querier) ([]*Agent, error) {
	var structs []reform.Struct
	structs, err := q.SelectAllFrom(AgentTable, "ORDER BY agent_id")
	err = errors.Wrap(err, "failed to select Agents")
	agents := make([]*Agent, len(structs))
	for i, s := range structs {
		agents[i] = s.(*Agent)
	}
	return agents, err
}

// AgentAddPmmAgent creates PMMAgent.
func AgentAddPmmAgent(q *reform.Querier, runsOnNodeID string, customLabels map[string]string) (*Agent, error) {
	id, err := agentNewID(q)
	if err != nil {
		return nil, err
	}

	if _, err := FindNodeByID(q, runsOnNodeID); err != nil {
		return nil, err
	}

	row := &Agent{
		AgentID:      id,
		AgentType:    PMMAgentType,
		RunsOnNodeID: &runsOnNodeID,
	}
	if err := row.SetCustomLabels(customLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	return row, nil
}

// AgentAddNodeExporter creates NodeExporter agent.
func AgentAddNodeExporter(q *reform.Querier, pmmAgentID string, customLabels map[string]string) (*Agent, error) {
	id, err := agentNewID(q)
	if err != nil {
		return nil, err
	}

	pmmAgent, err := AgentFindByID(q, pmmAgentID)
	if err != nil {
		return nil, err
	}

	row := &Agent{
		AgentID:    id,
		AgentType:  NodeExporterType,
		PMMAgentID: &pmmAgentID,
	}
	if err := row.SetCustomLabels(customLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	err = q.Insert(&AgentNode{
		AgentID: row.AgentID,
		NodeID:  pointer.GetString(pmmAgent.RunsOnNodeID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return row, nil
}

// ChangeCommonExporterParams describe common change params for exporters.
type ChangeCommonExporterParams struct {
	AgentID            string
	CustomLabels       map[string]string
	Disabled           bool
	RemoveCustomLabels bool
}

// AgentChangeExporter changes common params for given agent.
func AgentChangeExporter(q *reform.Querier, params *ChangeCommonExporterParams) (*Agent, error) {
	row, err := AgentFindByID(q, params.AgentID)
	if err != nil {
		return nil, err
	}

	row.Disabled = params.Disabled

	if params.RemoveCustomLabels {
		if err = row.SetCustomLabels(nil); err != nil {
			return nil, err
		}
	}
	if len(params.CustomLabels) != 0 {
		if err = row.SetCustomLabels(params.CustomLabels); err != nil {
			return nil, err
		}
	}

	if err = q.Update(row); err != nil {
		return nil, errors.WithStack(err)
	}

	return row, nil
}

// AddExporterAgentParams params for add common exporter.
type AddExporterAgentParams struct {
	PMMAgentID   string
	ServiceID    string
	Username     string
	Password     string
	CustomLabels map[string]string
}

// AgentAddExporter adds exporter with given type.
func AgentAddExporter(q *reform.Querier, agentType AgentType, params *AddExporterAgentParams) (*Agent, error) {
	id, err := agentNewID(q)
	if err != nil {
		return nil, err
	}

	if _, err := FindServiceByID(q, params.ServiceID); err != nil {
		return nil, err
	}

	row := &Agent{
		AgentID:    id,
		AgentType:  agentType,
		PMMAgentID: &params.PMMAgentID,
		Username:   pointer.ToStringOrNil(params.Username),
		Password:   pointer.ToStringOrNil(params.Password),
	}
	if err := row.SetCustomLabels(params.CustomLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	err = q.Insert(&AgentService{
		AgentID:   row.AgentID,
		ServiceID: params.ServiceID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return row, nil
}

// AgentsForNode returns all Agents providing insights for given Node.
func AgentsForNode(q *reform.Querier, nodeID string) ([]*Agent, error) {
	structs, err := q.FindAllFrom(AgentNodeView, "node_id", nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agent IDs")
	}

	agentIDs := make([]interface{}, len(structs))
	for i, s := range structs {
		agentIDs[i] = s.(*AgentNode).AgentID
	}
	if len(agentIDs) == 0 {
		return []*Agent{}, nil
	}

	p := strings.Join(q.Placeholders(1, len(agentIDs)), ", ")
	tail := fmt.Sprintf("WHERE agent_id IN (%s) ORDER BY agent_id", p) //nolint:gosec
	structs, err = q.SelectAllFrom(AgentTable, tail, agentIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agents")
	}

	res := make([]*Agent, len(structs))
	for i, s := range structs {
		res[i] = s.(*Agent)
	}
	return res, nil
}

// AgentsRunningByPMMAgent returns all Agents running by PMMAgent.
func AgentsRunningByPMMAgent(q *reform.Querier, pmmAgentID string) ([]*Agent, error) {
	tail := fmt.Sprintf("WHERE pmm_agent_id = %s ORDER BY agent_id", q.Placeholder(1)) //nolint:gosec
	structs, err := q.SelectAllFrom(AgentTable, tail, pmmAgentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agents")
	}

	res := make([]*Agent, len(structs))
	for i, s := range structs {
		res[i] = s.(*Agent)
	}
	return res, nil
}

// AgentsForService returns all Agents providing insights for given Service.
func AgentsForService(q *reform.Querier, serviceID string) ([]*Agent, error) {
	structs, err := q.FindAllFrom(AgentServiceView, "service_id", serviceID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agent IDs")
	}

	agentIDs := make([]interface{}, len(structs))
	for i, s := range structs {
		agentIDs[i] = s.(*AgentService).AgentID
	}
	if len(agentIDs) == 0 {
		return []*Agent{}, nil
	}

	p := strings.Join(q.Placeholders(1, len(agentIDs)), ", ")
	tail := fmt.Sprintf("WHERE agent_id IN (%s) ORDER BY agent_id", p) //nolint:gosec
	structs, err = q.SelectAllFrom(AgentTable, tail, agentIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agents")
	}

	res := make([]*Agent, len(structs))
	for i, s := range structs {
		res[i] = s.(*Agent)
	}
	return res, nil
}

//go:generate reform

// AgentType represents Agent type as stored in database.
type AgentType string

// Agent types.
const (
	PMMAgentType                AgentType = "pmm-agent"
	NodeExporterType            AgentType = "node_exporter"
	MySQLdExporterType          AgentType = "mysqld_exporter"
	MongoDBExporterType         AgentType = "mongodb_exporter"
	QANMySQLPerfSchemaAgentType AgentType = "qan-mysql-perfschema-agent"
	PostgresExporterType        AgentType = "postgres_exporter"
)

// Agent represents Agent as stored in database.
//reform:agents
type Agent struct {
	AgentID      string    `reform:"agent_id,pk"`
	AgentType    AgentType `reform:"agent_type"`
	RunsOnNodeID *string   `reform:"runs_on_node_id"`
	PMMAgentID   *string   `reform:"pmm_agent_id"`
	CustomLabels []byte    `reform:"custom_labels"`
	CreatedAt    time.Time `reform:"created_at"`
	UpdatedAt    time.Time `reform:"updated_at"`

	Disabled   bool    `reform:"disabled"`
	Status     string  `reform:"status"`
	ListenPort *uint16 `reform:"listen_port"`
	Version    *string `reform:"version"`

	Username   *string `reform:"username"`
	Password   *string `reform:"password"`
	MetricsURL *string `reform:"metrics_url"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (s *Agent) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (s *Agent) BeforeUpdate() error {
	s.UpdatedAt = Now()
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (s *Agent) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// GetCustomLabels decodes custom labels.
func (s *Agent) GetCustomLabels() (map[string]string, error) {
	if len(s.CustomLabels) == 0 {
		return nil, nil
	}
	m := make(map[string]string)
	if err := json.Unmarshal(s.CustomLabels, &m); err != nil {
		return nil, errors.Wrap(err, "failed to decode custom labels")
	}
	return m, nil
}

// SetCustomLabels encodes custom labels.
func (s *Agent) SetCustomLabels(m map[string]string) error {
	if len(m) == 0 {
		s.CustomLabels = nil
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "failed to encode custom labels")
	}
	s.CustomLabels = b
	return nil
}

func (s *Agent) IsChild() bool {
	return pointer.GetString(s.PMMAgentID) != ""
}

func (s *Agent) IsPMMAgent() bool {
	return s.AgentType == PMMAgentType
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Agent)(nil)
	_ reform.BeforeUpdater  = (*Agent)(nil)
	_ reform.AfterFinder    = (*Agent)(nil)
)
