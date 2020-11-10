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

package grpc

import (
	"context"
	"time"

	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/managementpb"
	"github.com/percona/pmm/version"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents"
)

type actionsServer struct {
	r  *agents.Registry
	db *reform.DB
	l  *logrus.Entry
}

// structure to keep id with timestamp
type idWithTimeStamp struct {
	id        string
	timestamp int64
}

var pmmAgent2100 = version.MustParse("2.10.0-HEAD") // TODO: Remove HEAD later once 2.11.0 is released.

// ****** TBD: Define where the period will be stored or configured
// Period that pt_summary action in seconds is cached for. After this time the pt-summary will be refreshed again.
const ptSummaryRefreshPeriod = 60

// Dictionary that of agent_id againts a structure of action_id and timestamp
var dicPtSummaryLastAction = make(map[string]idWithTimeStamp)

// NewActionsServer creates Management Actions Server.
func NewActionsServer(r *agents.Registry, db *reform.DB) managementpb.ActionsServer {
	l := logrus.WithField("component", "actions")
	return &actionsServer{r, db, l}
}

// GetAction gets an action result.
func (s *actionsServer) GetAction(ctx context.Context, req *managementpb.GetActionRequest) (*managementpb.GetActionResponse, error) {
	res, err := models.FindActionResultByID(s.db.Querier, req.ActionId)
	if err != nil {
		return nil, err
	}

	return &managementpb.GetActionResponse{
		ActionId:   res.ID,
		PmmAgentId: res.PMMAgentID,
		Done:       res.Done,
		Error:      res.Error,
		Output:     res.Output,
	}, nil
}

func (s *actionsServer) prepareServiceAction(serviceID, pmmAgentID, database string) (*models.ActionResult, string, error) {
	var res *models.ActionResult
	var dsn string
	e := s.db.InTransaction(func(tx *reform.TX) error {
		agents, err := models.FindPMMAgentsForService(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		if pmmAgentID, err = models.FindPmmAgentIDToRunAction(pmmAgentID, agents); err != nil {
			return err
		}

		if dsn, err = models.FindDSNByServiceIDandPMMAgentID(tx.Querier, serviceID, pmmAgentID, database); err != nil {
			return err
		}

		res, err = models.CreateActionResult(tx.Querier, pmmAgentID)
		return err
	})
	if e != nil {
		return nil, "", e
	}
	return res, dsn, nil
}

// StartMySQLExplainAction starts MySQL EXPLAIN Action with traditional output.
//nolint:lll
func (s *actionsServer) StartMySQLExplainAction(ctx context.Context, req *managementpb.StartMySQLExplainActionRequest) (*managementpb.StartMySQLExplainActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLExplainAction(ctx, res.ID, res.PMMAgentID, dsn, req.Query, agentpb.MysqlExplainOutputFormat_MYSQL_EXPLAIN_OUTPUT_FORMAT_DEFAULT)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMySQLExplainActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartMySQLExplainJSONAction starts MySQL EXPLAIN Action with JSON output.
//nolint:lll
func (s *actionsServer) StartMySQLExplainJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainJSONActionRequest) (*managementpb.StartMySQLExplainJSONActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLExplainAction(ctx, res.ID, res.PMMAgentID, dsn, req.Query, agentpb.MysqlExplainOutputFormat_MYSQL_EXPLAIN_OUTPUT_FORMAT_JSON)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMySQLExplainJSONActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartMySQLExplainTraditionalJSONAction starts MySQL EXPLAIN Action with traditional JSON output.
//nolint:lll
func (s *actionsServer) StartMySQLExplainTraditionalJSONAction(ctx context.Context, req *managementpb.StartMySQLExplainTraditionalJSONActionRequest) (*managementpb.StartMySQLExplainTraditionalJSONActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLExplainAction(ctx, res.ID, res.PMMAgentID, dsn, req.Query, agentpb.MysqlExplainOutputFormat_MYSQL_EXPLAIN_OUTPUT_FORMAT_TRADITIONAL_JSON)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMySQLExplainTraditionalJSONActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartMySQLShowCreateTableAction starts MySQL SHOW CREATE TABLE Action.
//nolint:lll
func (s *actionsServer) StartMySQLShowCreateTableAction(ctx context.Context, req *managementpb.StartMySQLShowCreateTableActionRequest) (*managementpb.StartMySQLShowCreateTableActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLShowCreateTableAction(ctx, res.ID, res.PMMAgentID, dsn, req.TableName)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMySQLShowCreateTableActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartMySQLShowTableStatusAction starts MySQL SHOW TABLE STATUS Action.
//nolint:lll
func (s *actionsServer) StartMySQLShowTableStatusAction(ctx context.Context, req *managementpb.StartMySQLShowTableStatusActionRequest) (*managementpb.StartMySQLShowTableStatusActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLShowTableStatusAction(ctx, res.ID, res.PMMAgentID, dsn, req.TableName)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMySQLShowTableStatusActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartMySQLShowIndexAction starts MySQL SHOW INDEX Action.
//nolint:lll
func (s *actionsServer) StartMySQLShowIndexAction(ctx context.Context, req *managementpb.StartMySQLShowIndexActionRequest) (*managementpb.StartMySQLShowIndexActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLShowIndexAction(ctx, res.ID, res.PMMAgentID, dsn, req.TableName)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMySQLShowIndexActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartPostgreSQLShowCreateTableAction starts PostgreSQL SHOW CREATE TABLE Action.
//nolint:lll
func (s *actionsServer) StartPostgreSQLShowCreateTableAction(ctx context.Context, req *managementpb.StartPostgreSQLShowCreateTableActionRequest) (*managementpb.StartPostgreSQLShowCreateTableActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartPostgreSQLShowCreateTableAction(ctx, res.ID, res.PMMAgentID, dsn, req.TableName)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartPostgreSQLShowCreateTableActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartPostgreSQLShowIndexAction starts PostgreSQL SHOW INDEX Action.
//nolint:lll
func (s *actionsServer) StartPostgreSQLShowIndexAction(ctx context.Context, req *managementpb.StartPostgreSQLShowIndexActionRequest) (*managementpb.StartPostgreSQLShowIndexActionResponse, error) {
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, req.Database)
	if err != nil {
		return nil, err
	}

	err = s.r.StartPostgreSQLShowIndexAction(ctx, res.ID, res.PMMAgentID, dsn, req.TableName)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartPostgreSQLShowIndexActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartMongoDBExplainAction starts MongoDB Explain action
//nolint:lll
func (s *actionsServer) StartMongoDBExplainAction(ctx context.Context, req *managementpb.StartMongoDBExplainActionRequest) (
	*managementpb.StartMongoDBExplainActionResponse, error) {
	// Explain action must be executed against the admin database
	res, dsn, err := s.prepareServiceAction(req.ServiceId, req.PmmAgentId, "admin")
	if err != nil {
		return nil, err
	}

	err = s.r.StartMongoDBExplainAction(ctx, res.ID, res.PMMAgentID, dsn, req.Query)
	if err != nil {
		return nil, err
	}

	return &managementpb.StartMongoDBExplainActionResponse{
		PmmAgentId: req.PmmAgentId,
		ActionId:   res.ID,
	}, nil
}

// StartPTSummaryAction starts pt-summary action. If the time since the last successfull start of pt-summary action is lower than ptSummaryRefreshPeriod,
// the response from the last action will be used. If the time is longer, the new pt-summary action will be called.
//nolint:lll
func (s *actionsServer) StartPTSummaryAction(ctx context.Context, psReq *managementpb.StartPTSummaryActionRequest) (*managementpb.StartPTSummaryActionResponse, error) {
	// Gets current timestamp
	timeNow := time.Now().Unix()

	// Gets pointers to agents running on the node
	psAgents, err := models.FindPMMAgentsRunningOnNode(s.db.Querier, psReq.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "No pmm-agent running on this node")
	}

	// Filter by the version
	psAgents = models.FindPMMAgentsForVersion(s.l, psAgents, pmmAgent2100)

	// No agent found
	if len(psAgents) == 0 {
		return nil, status.Error(codes.NotFound, "all available agents are outdated")
	}

	// Gets the agent by ID
	agentID, err := models.FindPmmAgentIDToRunAction(psReq.PmmAgentId, psAgents)
	if err != nil {
		return nil, err
	}

	// If founds an old action record for the agent
	if sAction, bFound := dicPtSummaryLastAction[agentID]; bFound {
		// If the time since the last call is less than 30 s
		if timeNow-sAction.timestamp < ptSummaryRefreshPeriod {
			// If found the last pt-summary response
			if _, err := models.FindActionResultByID(s.db.Querier, sAction.id); err == nil {
				// Returns the pointer to the found action response
				return &managementpb.StartPTSummaryActionResponse{PmmAgentId: agentID, ActionId: sAction.id}, nil
			}
		}
	}

	// Gets a pointer to the created action result structure
	pActRes, err := models.CreateActionResult(s.db.Querier, agentID)
	if err != nil {
		return nil, err
	}

	// PT summary action
	err = s.r.StartPTSummaryAction(ctx, pActRes.ID, agentID)
	if err != nil {
		return nil, err
	}

	// Saves the created action_id and timestand to the particular agentID
	dicPtSummaryLastAction[agentID] = idWithTimeStamp{pActRes.ID, timeNow}

	// Returns the pointer to the action response
	return &managementpb.StartPTSummaryActionResponse{PmmAgentId: agentID, ActionId: pActRes.ID}, nil
}

// CancelAction stops an Action.
func (s *actionsServer) CancelAction(ctx context.Context, req *managementpb.CancelActionRequest) (*managementpb.CancelActionResponse, error) {
	ar, err := models.FindActionResultByID(s.db.Querier, req.ActionId)
	if err != nil {
		return nil, err
	}

	err = s.r.StopAction(ctx, ar.ID)
	if err != nil {
		return nil, err
	}

	return &managementpb.CancelActionResponse{}, nil
}
