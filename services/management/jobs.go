package management

import (
	"context"

	jobsAPI "github.com/percona/pmm/api/managementpb/jobs"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/jobs"
)

// JobsAPIService provides methods for Jobs starting and management.
type JobsAPIService struct {
	db          *reform.DB
	jobsService *jobs.Service
}

// NewJobsAPIServer creates new jobs service.
func NewJobsAPIServer(db *reform.DB, service *jobs.Service) *JobsAPIService {
	return &JobsAPIService{
		db:          db,
		jobsService: service,
	}
}

// GetJob returns job result.
func (s *JobsAPIService) GetJob(ctx context.Context, req *jobsAPI.GetJobRequest) (*jobsAPI.GetJobResponse, error) {
	result, err := models.FindJobResultByID(s.db.Querier, req.JobId)
	if err != nil {
		return nil, err
	}

	resp := &jobsAPI.GetJobResponse{
		JobId:      result.ID,
		PmmAgentId: result.PMMAgentID,
		Done:       result.Done,
	}

	if !result.Done {
		return resp, nil
	}

	if result.Error != "" {
		resp.Result = &jobsAPI.GetJobResponse_Error_{
			Error: &jobsAPI.GetJobResponse_Error{
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

		resp.Result = &jobsAPI.GetJobResponse_Echo_{
			Echo: &jobsAPI.GetJobResponse_Echo{
				Message: echoResult.Message,
			}}
	default:
		return nil, errors.Errorf("Unexpected job type: %s", result.Type)
	}

	return resp, nil
}

// StartEchoJob starts echo job. Its purpose is testing.
func (s *JobsAPIService) StartEchoJob(ctx context.Context, req *jobsAPI.StartEchoJobRequest) (*jobsAPI.StartEchoJobResponse, error) {
	res, err := s.prepareAgentJob(req.PmmAgentId, models.Echo)
	if err != nil {
		return nil, err
	}

	if err := s.jobsService.StartEchoJob(res.ID, res.PMMAgentID, req.Timeout.AsDuration(), req.Message, req.Delay.AsDuration()); err != nil {
		return nil, err
	}

	return &jobsAPI.StartEchoJobResponse{
		PmmAgentId: req.PmmAgentId,
		JobId:      res.ID,
	}, nil
}

// CancelJob terminates job.
func (s *JobsAPIService) CancelJob(ctx context.Context, req *jobsAPI.CancelJobRequest) (*jobsAPI.CancelJobResponse, error) {
	if err := s.jobsService.StopJob(req.JobId); err != nil {
		return nil, err
	}

	return &jobsAPI.CancelJobResponse{}, nil
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
