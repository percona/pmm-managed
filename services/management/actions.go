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

package management

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/percona/pmm/api/agentpb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

type agentsRegistry interface {
	// SendRequest sends request to pmm-agent with given id.
	SendRequest(ctx context.Context, pmmAgentID string, payload agentpb.ServerRequestPayload) agentpb.AgentResponsePayload
}

// ActionsService describes an Actions Application Service.
// Provides functions for PMM Actions manipulation.
type ActionsService struct {
	agentsRegistry agentsRegistry
	actionsStorage *InMemoryActionsStorage
	logger         logrus.FieldLogger
	db             *reform.DB
}

// NewActionsService creates new actions service.
func NewActionsService(r agentsRegistry, s *InMemoryActionsStorage, db *reform.DB) *ActionsService {
	return &ActionsService{
		agentsRegistry: r,
		actionsStorage: s,
		logger:         logrus.WithField("component", "actions-service"),
		db:             db,
	}
}

// RunActionParams parameters for run actions.
type RunActionParams struct {
	ActionName   agentpb.ActionName
	ActionParams []string
	PmmAgentID   string
	NodeID       string
	ServiceID    string
}

// RunAction runs PMM Action on the given client.
func (a *ActionsService) RunAction(ctx context.Context, rp RunActionParams) (string, error) {
	actionID := "/action_id/" + uuid.New().String()
	pmmAgentId := ""
	params := rp.ActionParams
	var err error

	switch rp.ActionName {
	case agentpb.ActionName_PT_SUMMARY:
		pmmAgentId, err = findPmmAgentIdByNodeId(a.db.Querier, rp.PmmAgentID, rp.NodeID)
		if err != nil {
			return "", err
		}

	case agentpb.ActionName_PT_MYSQL_SUMMARY:
		pmmAgentId, err = findPmmAgentIdByServiceId(a.db.Querier, rp.PmmAgentID, rp.ServiceID)
		if err != nil {
			return "", err
		}

	case agentpb.ActionName_MYSQL_EXPLAIN:
		pmmAgentId, err = findPmmAgentIdByServiceId(a.db.Querier, rp.PmmAgentID, rp.ServiceID)
		if err != nil {
			return "", err
		}

	case agentpb.ActionName_ACTION_NAME_INVALID:
		err = errors.New("unknown action name")
	}

	if pmmAgentId == "" {
		return "", errors.New("can't find proper pmm_agent_id")
	}

	res := a.agentsRegistry.SendRequest(ctx, pmmAgentId, &agentpb.StartActionRequest{
		Id:         actionID,
		Name:       rp.ActionName,
		Parameters: params,
	})
	a.logger.Infof("RunAction response: %+v.", res)
	return actionID, nil
}

// CancelAction stops PMM Action with the given ID on the given client.
func (a *ActionsService) CancelAction(ctx context.Context, pmmAgentID, actionID string) {
	res := a.agentsRegistry.SendRequest(ctx, pmmAgentID, &agentpb.StopActionRequest{
		Id: actionID,
	})
	a.logger.Infof("CancelAction response: %+v.", res)
}

// GetActionResult gets PMM Action with the given ID from action results storage.
func (a *ActionsService) GetActionResult(ctx context.Context, actionID string) (ActionResult, bool) {
	return a.actionsStorage.Load(actionID)
}

func findPmmAgentIdByNodeId(q *reform.Querier, pmmAgentId, nodeID string) (string, error) {
	agents, err := models.PMMAgentsForNode(q, nodeID)
	if err != nil {
		return "", err
	}
	return findPmmAgentIdInAgents(pmmAgentId, agents)
}

func findPmmAgentIdByServiceId(q *reform.Querier, pmmAgentId, serviceID string) (string, error) {
	agents, err := models.PMMAgentsForService(q, serviceID)
	if err != nil {
		return "", err
	}
	return findPmmAgentIdInAgents(pmmAgentId, agents)
}

func findPmmAgentIdInAgents(pmmAgentId string, agents []*models.Agent) (string, error) {
	// no explicit ID is given, and there is only one
	if pmmAgentId == "" && len(agents) == 1 {
		return agents[0].AgentID, nil
	}

	// no explicit ID is given, and there are zero or several
	if pmmAgentId == "" {
		return "", errors.New("can't detect pmm_agent_id")
	}

	// check that explicit agent id is correct
	for _, a := range agents {
		if a.AgentID == pmmAgentId {
			return a.AgentID, nil
		}
	}
	return "", errors.New("contradiction")
}

// ActionResult describes an PMM Action result which is storing in ActionsResult storage.
type ActionResult struct {
	ID         string
	PmmAgentID string
	Output     string
}

// InMemoryActionsStorage in memory action results storage.
type InMemoryActionsStorage struct {
	container sync.Map
}

// NewInMemoryActionsStorage created new InMemoryActionsStorage.
func NewInMemoryActionsStorage() *InMemoryActionsStorage {
	return &InMemoryActionsStorage{}
}

// TODO: Store action result first, then store action output.

// Store stores an action result in action results storage.
func (s *InMemoryActionsStorage) Store(result ActionResult) {
	s.container.Store(result.ID, result)
}

// Load gets an action result from storage by action id.
func (s *InMemoryActionsStorage) Load(id string) (ActionResult, bool) {
	v, ok := s.container.Load(id)
	if !ok {
		return ActionResult{}, false
	}
	return v.(ActionResult), true
}
