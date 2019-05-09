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
	"github.com/sirupsen/logrus"
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
}

// NewActionsService creates new actions service.
func NewActionsService(r agentsRegistry, s *InMemoryActionsStorage) *ActionsService {
	return &ActionsService{
		agentsRegistry: r,
		actionsStorage: s,
		logger:         logrus.WithField("component", "actions-service"),
	}
}

// RunAction runs PMM Action on the given client.
func (a *ActionsService) RunAction(ctx context.Context, pmmAgentID string, actionName agentpb.ActionName, params []string) string {
	actionID := "/action_id/" + uuid.New().String()
	res := a.agentsRegistry.SendRequest(ctx, pmmAgentID, &agentpb.StartActionRequest{
		Id:         actionID,
		Name:       actionName,
		Parameters: params,
	})
	a.logger.Infof("RunAction response: %+v.", res)
	return actionID
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

// Store stores an action result in action results storage.
func (s *InMemoryActionsStorage) Store(result ActionResult) {
	s.container.Store(result.ID, result)
}

// Get gets an action result from storage by action id.
func (s *InMemoryActionsStorage) Load(id string) (ActionResult, bool) {
	v, ok := s.container.Load(id)
	if !ok {
		return ActionResult{}, false
	}
	return v.(ActionResult), true
}
