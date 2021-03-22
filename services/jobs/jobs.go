package jobs

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona/pmm/api/agentpb"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents"
)

// Service provides methods for managing jobs.
type Service struct {
	r  *agents.Registry
	db *reform.DB
}

// New returns new jobs service.
func New(db *reform.DB, registry *agents.Registry) *Service {
	return &Service{
		r:  registry,
		db: db,
	}
}

// StartEchoJob starts echo job on the pmm-agent.
func (s *Service) StartEchoJob(id, pmmAgentID string, timeout time.Duration, message string, delay time.Duration) error {
	req := &agentpb.StartJobRequest{
		JobId:   id,
		Timeout: ptypes.DurationProto(timeout),
		Job: &agentpb.StartJobRequest_Echo_{
			Echo: &agentpb.StartJobRequest_Echo{
				Message: message,
				Delay:   ptypes.DurationProto(delay),
			},
		},
	}

	channel, err := s.r.GetAgentChannel(pmmAgentID)
	if err != nil {
		return err
	}

	channel.SendAndWaitResponse(req)

	return nil
}

// StopJob stops job with given given id.
func (s *Service) StopJob(jobID string) error {
	jobResult, err := models.FindJobResultByID(s.db.Querier, jobID)
	if err != nil {
		return errors.WithStack(err)
	}

	if jobResult.Done {
		// Job already finished
		return nil
	}

	channel, err := s.r.GetAgentChannel(jobResult.PMMAgentID)
	if err != nil {
		return errors.WithStack(err)
	}

	channel.SendAndWaitResponse(&agentpb.StopJobRequest{JobId: jobID})

	return nil
}
