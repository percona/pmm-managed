package management

import (
	"context"

	"github.com/percona/pmm/api/managementpb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents"
)

type JobsAPIService struct {
	db *reform.DB
	r  *agents.Registry
}

func NewJobsAPIServer(db *reform.DB, registry *agents.Registry) *JobsAPIService {
	return &JobsAPIService{
		db: db,
		r:  registry,
	}
}

func (s *JobsAPIService) GetAction(ctx context.Context, req *managementpb.GetJobRequest) (*managementpb.GetJobResponse, error) {
	panic("implement me")
}

func (s *JobsAPIService) StartEchoJob(ctx context.Context, req *managementpb.StartEchoJobRequest) (*managementpb.StartEchoJobResponse, error) {
	res, err := s.prepareAgentJob(req.PmmAgentId)
	if err != nil {
		return nil, err
	}

	if err := s.r.StartEchoJob(ctx, res.ID, res.PMMAgentID, req.Message, req.Delay.AsDuration()); err != nil {
		return nil, err
	}

	return &managementpb.StartEchoJobResponse{
		PmmAgentId: req.PmmAgentId,
		JobId:      res.ID,
	}, nil
}

func (s *JobsAPIService) CancelAction(ctx context.Context, req *managementpb.CancelJobRequest) (*managementpb.CancelJobResponse, error) {
	panic("implement me")
}

func (s *JobsAPIService) prepareAgentJob(pmmAgentID string) (*models.JobResult, error) {
	var res *models.JobResult
	e := s.db.InTransaction(func(tx *reform.TX) error {
		_, err := models.FindAgentByID(tx.Querier, pmmAgentID)
		if err != nil {
			return err
		}

		res, err = models.CreateJobResult(tx.Querier, pmmAgentID)
		return err
	})
	if e != nil {
		return nil, e
	}
	return res, nil
}

func (s *JobsAPIService) prepareServiceJob(serviceID, pmmAgentID, database string) (*models.JobResult, string, error) {
	var res *models.JobResult
	var dsn string
	e := s.db.InTransaction(func(tx *reform.TX) error {
		agents, err := models.FindPMMAgentsForService(tx.Querier, serviceID)
		if err != nil {
			return err
		}

		if pmmAgentID, err = models.FindPmmAgentIDToRunActionOrJob(pmmAgentID, agents); err != nil {
			return err
		}

		if dsn, _, err = models.FindDSNByServiceIDandPMMAgentID(tx.Querier, serviceID, pmmAgentID, database); err != nil {
			return err
		}

		res, err = models.CreateJobResult(tx.Querier, pmmAgentID)
		return err
	})
	if e != nil {
		return nil, "", e
	}
	return res, dsn, nil
}
