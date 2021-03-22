package management

import (
	"context"

	jobs "github.com/percona/pmm/api/managementpb/jobs"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents"
)

// JobsAPIService provides methods for Jobs starting and management.
type JobsAPIService struct {
	db *reform.DB
	r  *agents.Registry
}

// NewJobsAPIServer creates new jobs service.
func NewJobsAPIServer(db *reform.DB, registry *agents.Registry) *JobsAPIService {
	return &JobsAPIService{
		db: db,
		r:  registry,
	}
}

// GetJob returns job result.
func (s *JobsAPIService) GetJob(ctx context.Context, req *jobs.GetJobRequest) (*jobs.GetJobResponse, error) {
	result, err := models.FindJobResultByID(s.db.Querier, req.JobId)
	if err != nil {
		return nil, err
	}

	resp := &jobs.GetJobResponse{
		JobId:      result.ID,
		PmmAgentId: result.PMMAgentID,
		Done:       result.Done,
	}

	if !result.Done {
		return resp, nil
	}

	if result.Error != "" {
		resp.Result = &jobs.GetJobResponse_Error_{
			Error: &jobs.GetJobResponse_Error{
				Message: result.Error,
			},
		}

		return resp, nil
	}

	switch result.Type {
	case models.Echo:
		echoResult, err := result.GetEchoJobResult()
		if err != nil {
			return nil, err
		}

		resp.Result = &jobs.GetJobResponse_Echo_{
			Echo: &jobs.GetJobResponse_Echo{
				Message: echoResult.Message,
			}}
	default:
		return nil, errors.Errorf("Unexpected job type: %s", result.Type)
	}

	return resp, nil
}

// StartEchoJob starts echo job. Its purpose is testing.
func (s *JobsAPIService) StartEchoJob(ctx context.Context, req *jobs.StartEchoJobRequest) (*jobs.StartEchoJobResponse, error) {
	res, err := s.prepareAgentJob(req.PmmAgentId, models.Echo)
	if err != nil {
		return nil, err
	}

	if err := s.r.StartEchoJob(res.ID, res.PMMAgentID, req.Timeout.AsDuration(), req.Message, req.Delay.AsDuration()); err != nil {
		return nil, err
	}

	return &jobs.StartEchoJobResponse{
		PmmAgentId: req.PmmAgentId,
		JobId:      res.ID,
	}, nil
}

// CancelJob terminates job.
func (s *JobsAPIService) CancelJob(ctx context.Context, req *jobs.CancelJobRequest) (*jobs.CancelJobResponse, error) {
	if err := s.r.StopJob(req.JobId); err != nil {
		return nil, err
	}

	return &jobs.CancelJobResponse{}, nil
}

func (s *JobsAPIService) prepareAgentJob(pmmAgentID string, jobType models.JobType) (*models.JobResult, error) {
	var res *models.JobResult
	e := s.db.InTransaction(func(tx *reform.TX) error {
		_, err := models.FindAgentByID(tx.Querier, pmmAgentID)
		if err != nil {
			return err
		}

		res, err = models.CreateJobResult(tx.Querier, pmmAgentID, jobType)
		return err
	})
	if e != nil {
		return nil, e
	}
	return res, nil
}

func (s *JobsAPIService) prepareServiceJob(serviceID, pmmAgentID, database string, jobType models.JobType) (*models.JobResult, string, error) {
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

		res, err = models.CreateJobResult(tx.Querier, pmmAgentID, jobType)
		return err
	})
	if e != nil {
		return nil, "", e
	}
	return res, dsn, nil
}
