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
	"github.com/percona/pmm/api/managementpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

var (
	errUnsupportedAction  = status.Error(codes.InvalidArgument, "unsupported action")
	errPmmAgentIDNotFound = status.Error(codes.Internal, "can't detect pmm_agent_id")
)

// ActionsService describes an Actions Application Service.
// Provides functions for PMM Actions manipulation.
type ActionsService struct {
	registry registry
	storage  *InMemoryActionsStorage
	logger   *logrus.Entry
	db       *reform.DB
}

// NewActionsService creates new actions service.
func NewActionsService(r registry, s *InMemoryActionsStorage, db *reform.DB) *ActionsService {
	return &ActionsService{
		registry: r,
		storage:  s,
		logger:   logrus.WithField("component", "actions-service"),
		db:       db,
	}
}

// RunActionParams parameters for run actions.
type RunActionParams struct {
	ActionName   managementpb.ActionType
	ActionParams []string
	PmmAgentID   string
	NodeID       string
	ServiceID    string
}

// RunAction runs PMM Action on the given client.
// First parameter returned by this method is "ActionID".
func (a *ActionsService) RunAction(ctx context.Context, rp *RunActionParams) (string, error) {
	action, err := a.prepareAction(rp)
	if err != nil {
		return "", err
	}

	req := &agentpb.StartActionRequest{
		ActionId: action.ID,
		Type:     action.Name,
	}
	switch rp.ActionName {
	case managementpb.ActionType_PT_SUMMARY:
		req.Params = &agentpb.StartActionRequest_ProcessParams_{
			ProcessParams: &agentpb.StartActionRequest_ProcessParams{
				Args: action.Params,
			},
		}
	case managementpb.ActionType_PT_MYSQL_SUMMARY:
		req.Params = &agentpb.StartActionRequest_ProcessParams_{
			ProcessParams: &agentpb.StartActionRequest_ProcessParams{
				Args: action.Params,
			},
		}
	case managementpb.ActionType_MYSQL_EXPLAIN:
		return "", errUnsupportedAction
	case managementpb.ActionType_ACTION_TYPE_INVALID:
		return "", errUnsupportedAction
	}

	res := a.registry.SendRequest(ctx, action.PmmAgentID, req)
	a.logger.Infof("RunAction response: %+v.", res)
	return action.ID, nil
}

// CancelAction stops PMM Action with the given ID on the given client.
func (a *ActionsService) CancelAction(ctx context.Context, actionID string) {
	action, ok := a.storage.Load(ctx, actionID)
	if !ok {
		a.logger.Errorf("Unknown action with ID: %s.", actionID)
		return
	}

	res := a.registry.SendRequest(ctx, action.PmmAgentID, &agentpb.StopActionRequest{
		ActionId: actionID,
	})
	a.logger.Infof("CancelAction response: %+v.", res)
}

// GetActionResult gets PMM Action with the given ID from action results storage.
//nolint:unparam
func (a *ActionsService) GetActionResult(ctx context.Context, actionID string) (*models.ActionResult, bool) {
	return a.storage.Load(ctx, actionID)
}

type preparedAction struct {
	ID         string
	Name       managementpb.ActionType
	Params     []string
	PmmAgentID string
}

func (a *ActionsService) prepareAction(rp *RunActionParams) (preparedAction, error) {
	action := preparedAction{
		ID:         "/action_id/" + uuid.New().String(),
		PmmAgentID: rp.PmmAgentID,
		Name:       rp.ActionName,
		Params:     rp.ActionParams,
	}
	var err error

	switch action.Name {
	case managementpb.ActionType_PT_SUMMARY:
		action.PmmAgentID, err = findPmmAgentIDByNodeID(a.db.Querier, rp.PmmAgentID, rp.NodeID)

	case managementpb.ActionType_PT_MYSQL_SUMMARY:
		action.PmmAgentID, err = findPmmAgentIDByServiceID(a.db.Querier, rp.PmmAgentID, rp.ServiceID)

	case managementpb.ActionType_MYSQL_EXPLAIN:
		action.PmmAgentID, err = findPmmAgentIDByServiceID(a.db.Querier, rp.PmmAgentID, rp.ServiceID)

	case managementpb.ActionType_ACTION_TYPE_INVALID:
		return action, errUnsupportedAction
	}

	return action, err
}

func findPmmAgentIDByNodeID(q *reform.Querier, pmmAgentID, nodeID string) (string, error) {
	agents, err := models.FindPMMAgentsForNode(q, nodeID)
	if err != nil {
		return "", err
	}
	return validatePmmAgentID(pmmAgentID, agents)
}

func findPmmAgentIDByServiceID(q *reform.Querier, pmmAgentID, serviceID string) (string, error) {
	agents, err := models.FindPMMAgentsForService(q, serviceID)
	if err != nil {
		return "", err
	}
	return validatePmmAgentID(pmmAgentID, agents)
}

func validatePmmAgentID(pmmAgentID string, agents []*models.Agent) (string, error) {
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

// InMemoryActionsStorage in memory action results storage.
type InMemoryActionsStorage struct {
	container map[string]models.ActionResult
	mx        sync.Mutex
}

// NewInMemoryActionsStorage created new InMemoryActionsStorage.
func NewInMemoryActionsStorage() *InMemoryActionsStorage {
	return &InMemoryActionsStorage{}
}

// Store stores an action result in action results storage.
func (s *InMemoryActionsStorage) Store(ctx context.Context, result *models.ActionResult) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.container[result.ID] = *result
}

// Load gets an action result from storage by action id.
func (s *InMemoryActionsStorage) Load(ctx context.Context, id string) (*models.ActionResult, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	v, ok := s.container[id]
	if !ok {
		return nil, false
	}
	return &v, true
}
