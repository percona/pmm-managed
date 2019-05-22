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
type RunProcessActionParams struct {
	ActionName   managementpb.ActionType
	ActionParams []string
	PmmAgentID   string
	NodeID       string
	ServiceID    string
}

// RunProcessAction runs Process PMM Action on the given client.
// First parameter returned by this method is "ActionID".
func (a *ActionsService) RunProcessAction(ctx context.Context, rp *RunProcessActionParams) (string, error) {
	actionID := getNewActionID()
	pmmAgentID, err := a.resolvePmmAgentID(rp.ActionName, rp.NodeID, rp.ServiceID, rp.PmmAgentID)
	if err != nil {
		return "", err
	}

	req := &agentpb.StartActionRequest{
		ActionId: actionID,
		Type:     rp.ActionName,
	}
	switch rp.ActionName {
	case managementpb.ActionType_PT_SUMMARY:
		req.Params = &agentpb.StartActionRequest_ProcessParams_{
			ProcessParams: &agentpb.StartActionRequest_ProcessParams{
				Args: rp.ActionParams,
			},
		}
	case managementpb.ActionType_PT_MYSQL_SUMMARY:
		req.Params = &agentpb.StartActionRequest_ProcessParams_{
			ProcessParams: &agentpb.StartActionRequest_ProcessParams{
				Args: rp.ActionParams,
			},
		}
	default:
		return "", errUnsupportedAction
	}

	res, err := a.registry.SendRequest(ctx, pmmAgentID, req)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	a.logger.Infof("RunProcessAction response: %+v.", res)
	return actionID, nil
}

// RunActionParams parameters for run actions.
type StartMySQLExplainActionParams struct {
	OutputFormat agentpb.MysqlExplainOutputFormat
	Query        string
	PmmAgentID   string
	ServiceID    string
}

// StartMySQLExplainAction runs MySQL Explain PMM Action on the given client.
// First parameter returned by this method is "ActionID".
func (a *ActionsService) StartMySQLExplainAction(ctx context.Context, rp *StartMySQLExplainActionParams) (string, error) {
	actionID := getNewActionID()
	pmmAgentID, err := a.resolvePmmAgentID(managementpb.ActionType_MYSQL_EXPLAIN, "", rp.ServiceID, rp.PmmAgentID)
	if err != nil {
		return "", err
	}

	req := &agentpb.StartActionRequest{
		ActionId: actionID,
		Type:     managementpb.ActionType_MYSQL_EXPLAIN,
	}

	req.Params = &agentpb.StartActionRequest_MysqlExplainParams{
		MysqlExplainParams: &agentpb.StartActionRequest_MySQLExplainParams{
			Dsn:          "", // TODO: Add DSN string: findAgentForService(MYSQLD_EXPORTER, req.ServiceID)
			Query:        rp.Query,
			OutputFormat: rp.OutputFormat,
		},
	}

	res, err := a.registry.SendRequest(ctx, pmmAgentID, req)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	a.logger.Infof("StartMySQLExplainAction response: %+v.", res)
	return actionID, nil
}

// CancelAction stops PMM Action with the given ID on the given client.
func (a *ActionsService) CancelAction(ctx context.Context, actionID string) {
	action, ok := a.storage.Load(ctx, actionID)
	if !ok {
		a.logger.Errorf("Unknown action with ID: %s.", actionID)
		return
	}

	res, err := a.registry.SendRequest(ctx, action.PmmAgentID, &agentpb.StopActionRequest{ActionId: actionID})
	if err != nil {
		a.logger.Error(err)
		return
	}

	a.logger.Infof("CancelAction response: %+v.", res)
}

// GetActionResult gets PMM Action with the given ID from action results storage.
//nolint:unparam
func (a *ActionsService) GetActionResult(ctx context.Context, actionID string) (*models.ActionResult, bool) {
	return a.storage.Load(ctx, actionID)
}

func (a *ActionsService) resolvePmmAgentID(actionType managementpb.ActionType, nodeID, serviceID, pmmAgentID string) (string, error) {
	switch actionType {
	case managementpb.ActionType_PT_SUMMARY:
		return findPmmAgentIDByNodeID(a.db.Querier, pmmAgentID, nodeID)

	case managementpb.ActionType_PT_MYSQL_SUMMARY:
		return findPmmAgentIDByServiceID(a.db.Querier, pmmAgentID, serviceID)

	case managementpb.ActionType_MYSQL_EXPLAIN:
		return findPmmAgentIDByServiceID(a.db.Querier, pmmAgentID, serviceID)

	case managementpb.ActionType_ACTION_TYPE_INVALID:
		return "", errUnsupportedAction
	}

	return "", errUnsupportedAction
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

func getNewActionID() string {
	return "/action_id/" + uuid.New().String()
}
