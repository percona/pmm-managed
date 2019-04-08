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
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

func newAgentID(q *reform.Querier) (string, error) {
	agentID := "/agent_id/" + uuid.New().String()
	agent := &Agent{AgentID: agentID}
	switch err := q.Reload(agent); err {
	case nil:
		return "", status.Errorf(codes.AlreadyExists, "Agent with ID %q already exists.", agentID)
	case reform.ErrNoRows:
		return agentID, nil
	default:
		return "", errors.WithStack(err)
	}
}

// FindAllAgents returns all Agents.
func FindAllAgents(q *reform.Querier) ([]*Agent, error) {
	structs, err := q.SelectAllFrom(AgentTable, "ORDER BY agent_id")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	agents := make([]*Agent, len(structs))
	for i, s := range structs {
		agents[i] = s.(*Agent)
	}

	return agents, nil
}

// FindAgentsForNodeID returns all Agents providing insights for given Node.
func FindAgentsForNodeID(q *reform.Querier, nodeID string) ([]*Agent, error) {
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

// FindAgentsForPMMAgentID returns all Agents running by PMMAgent.
func FindAgentsForPMMAgentID(q *reform.Querier, pmmAgentID string) ([]*Agent, error) {
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

// FindAgentsForServiceID returns all Agents providing insights for given Service.
func FindAgentsForServiceID(q *reform.Querier, serviceID string) ([]*Agent, error) {
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

// PMMAgentsForChangedNode returns pmm-agents IDs that are affected
// by the change of the Node with given ID.
// It may return (nil, nil) if no such pmm-agents are found.
// It returns wrapped reform.ErrNoRows if Service with given ID is not found.
func PMMAgentsForChangedNode(q *reform.Querier, nodeID string) ([]string, error) {
	// TODO Real code.
	// Returning all pmm-agents is currently safe, but not optimal for large number of Agents.
	_ = nodeID

	structs, err := q.SelectAllFrom(AgentTable, "ORDER BY agent_id")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agents")
	}

	var res []string
	for _, str := range structs {
		agent := str.(*Agent)
		if agent.AgentType == PMMAgentType {
			res = append(res, agent.AgentID)
		}
	}
	return res, nil
}

// PMMAgentsForChangedService returns pmm-agents IDs that are affected
// by the change of the Service with given ID.
// It may return (nil, nil) if no such pmm-agents are found.
// It returns wrapped reform.ErrNoRows if Service with given ID is not found.
func PMMAgentsForChangedService(q *reform.Querier, serviceID string) ([]string, error) {
	// TODO Real code. We need to returns IDs of pmm-agents that:
	// * run Agents providing insights for this Service;
	// * run Agents providing insights for Node that hosts this Service.
	// Returning all pmm-agents is currently safe, but not optimal for large number of Agents.
	_ = serviceID

	structs, err := q.SelectAllFrom(AgentTable, "ORDER BY agent_id")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Agents")
	}

	var res []string
	for _, str := range structs {
		agent := str.(*Agent)
		if agent.AgentType == PMMAgentType {
			res = append(res, agent.AgentID)
		}
	}
	return res, nil
}

// FindAgentByID finds agent by ID.
func FindAgentByID(q *reform.Querier, id string) (*Agent, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Agent ID.")
	}

	agent := &Agent{AgentID: id}
	switch err := q.Reload(agent); err {
	case nil:
		return agent, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Agent with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// AgentAddPmmAgent creates PMMAgent.
func AgentAddPmmAgent(q *reform.Querier, runsOnNodeID string, customLabels map[string]string) (*Agent, error) {
	id, err := newAgentID(q)
	if err != nil {
		return nil, err
	}

	if _, err := FindNodeByID(q, runsOnNodeID); err != nil {
		return nil, err
	}

	agent := &Agent{
		AgentID:      id,
		AgentType:    PMMAgentType,
		RunsOnNodeID: &runsOnNodeID,
	}
	if err := agent.SetCustomLabels(customLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(agent); err != nil {
		return nil, errors.WithStack(err)
	}

	return agent, nil
}

// AgentAddNodeExporter creates NodeExporter agent.
func AgentAddNodeExporter(q *reform.Querier, pmmAgentID string, customLabels map[string]string) (*Agent, error) {
	id, err := newAgentID(q)
	if err != nil {
		return nil, err
	}

	pmmAgent, err := FindAgentByID(q, pmmAgentID)
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		AgentID:    id,
		AgentType:  NodeExporterType,
		PMMAgentID: &pmmAgentID,
	}
	if err := agent.SetCustomLabels(customLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(agent); err != nil {
		return nil, errors.WithStack(err)
	}

	err = q.Insert(&AgentNode{
		AgentID: agent.AgentID,
		NodeID:  pointer.GetString(pmmAgent.RunsOnNodeID),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return agent, nil
}

// CreateAgentParams params for add common exporter.
type CreateAgentParams struct {
	PMMAgentID   string
	ServiceID    string
	Username     string
	Password     string
	CustomLabels map[string]string
}

// CreateAgent adds exporter with given type.
func CreateAgent(q *reform.Querier, agentType AgentType, params *CreateAgentParams) (*Agent, error) {
	id, err := newAgentID(q)
	if err != nil {
		return nil, err
	}

	if _, err := FindServiceByID(q, params.ServiceID); err != nil {
		return nil, err
	}

	agent := &Agent{
		AgentID:    id,
		AgentType:  agentType,
		PMMAgentID: &params.PMMAgentID,
		Username:   pointer.ToStringOrNil(params.Username),
		Password:   pointer.ToStringOrNil(params.Password),
	}
	if err := agent.SetCustomLabels(params.CustomLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(agent); err != nil {
		return nil, errors.WithStack(err)
	}

	err = q.Insert(&AgentService{
		AgentID:   agent.AgentID,
		ServiceID: params.ServiceID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return agent, nil
}

// UpdateAgentParams describe common change params for exporters.
type UpdateAgentParams struct {
	AgentID            string
	CustomLabels       map[string]string
	Disabled           bool
	RemoveCustomLabels bool
}

// UpdateAgent changes common params for given agent.
func UpdateAgent(q *reform.Querier, params *UpdateAgentParams) (*Agent, error) {
	agent, err := FindAgentByID(q, params.AgentID)
	if err != nil {
		return nil, err
	}

	agent.Disabled = params.Disabled

	if params.RemoveCustomLabels {
		if err = agent.SetCustomLabels(nil); err != nil {
			return nil, err
		}
	}
	if len(params.CustomLabels) != 0 {
		if err = agent.SetCustomLabels(params.CustomLabels); err != nil {
			return nil, err
		}
	}

	if err = q.Update(agent); err != nil {
		return nil, errors.WithStack(err)
	}

	return agent, nil
}

// RemoveAgent removes Agent by ID and returns it.
func RemoveAgent(q *reform.Querier, agentID string) (*Agent, error) {
	agent, err := FindAgentByID(q, agentID)
	if err != nil {
		return nil, err
	}

	if _, err = q.DeleteFrom(AgentServiceView, "WHERE agent_id = "+q.Placeholder(1), agentID); err != nil { //nolint:gosec
		return nil, errors.WithStack(err)
	}
	if _, err = q.DeleteFrom(AgentNodeView, "WHERE agent_id = "+q.Placeholder(1), agentID); err != nil { //nolint:gosec
		return nil, errors.WithStack(err)
	}

	if err = q.Delete(agent); err != nil {
		return nil, errors.WithStack(err)
	}

	return agent, nil
}
