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

	jobsAPI "github.com/percona/pmm/api/managementpb/jobs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// JobsAPIService provides methods for Jobs starting and management.
type JobsAPIService struct {
	l *logrus.Entry

	db          *reform.DB
	jobsService jobsService
}

// NewJobsAPIServer creates new jobs service.
func NewJobsAPIServer(db *reform.DB, service jobsService) *JobsAPIService {
	return &JobsAPIService{
		l: logrus.WithField("component", "management/jobs"),

		db:          db,
		jobsService: service,
	}
}

// GetJob returns job result.
func (s *JobsAPIService) GetJob(_ context.Context, req *jobsAPI.GetJobRequest) (*jobsAPI.GetJobResponse, error) {
	result, err := models.FindJobByID(s.db.Querier, req.JobId)
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
	case models.MySQLBackupJob, models.MySQLRestoreBackupJob, models.MongoDBBackupJob, models.MongoDBRestoreBackupJob:
	default:
		return nil, errors.Errorf("Unexpected job type: %s", result.Type)
	}

	return resp, nil
}

// CancelJob terminates job.
func (s *JobsAPIService) CancelJob(_ context.Context, req *jobsAPI.CancelJobRequest) (*jobsAPI.CancelJobResponse, error) {
	if err := s.jobsService.StopJob(req.JobId); err != nil {
		return nil, err
	}

	return &jobsAPI.CancelJobResponse{}, nil
}
