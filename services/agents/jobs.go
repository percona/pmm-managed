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

// Package jobs provides jobs functionality.
package agents

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona/pmm/api/agentpb"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// JobsService provides methods for managing jobs.
type JobsService struct {
	r  *Registry
	db *reform.DB
}

// NewJobsService returns new jobs service.
func NewJobsService(db *reform.DB, registry *Registry) *JobsService {
	return &JobsService{
		r:  registry,
		db: db,
	}
}

// StartEchoJob starts echo job on the pmm-agent.
func (s *JobsService) StartEchoJob(id, pmmAgentID string, timeout time.Duration, message string, delay time.Duration) error {
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

	agent, err := s.r.get(pmmAgentID)
	if err != nil {
		return err
	}

	resp := agent.channel.SendAndWaitResponse(req)

	if e := resp.(*agentpb.StartJobResponse).Error; e != "" {
		return errors.Errorf("failed to start echo job: %s", e)
	}

	return nil
}

// StartMySQLBackupJob starts mysql backup job on the pmm-agent.
func (s *JobsService) StartMySQLBackupJob(id, pmmAgentID string, timeout time.Duration, name, dsn string, locationConfig models.BackupLocationConfig) error {
	mysqlConfig, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}
	mySQLReq := &agentpb.StartJobRequest_MySQLBackup{
		Name:     name,
		User:     mysqlConfig.User,
		Password: mysqlConfig.Passwd,
	}

	switch mysqlConfig.Net {
	case "tcp":
		splitAddr := strings.Split(mysqlConfig.Addr, ":")
		mySQLReq.Address = splitAddr[0]
		if len(splitAddr) > 1 {
			port, err := strconv.ParseInt(splitAddr[1], 10, strconv.IntSize)
			if err != nil {
				return err
			}
			mySQLReq.Port = int32(port)
		}
	case "unix":
		mySQLReq.Socket = mysqlConfig.Addr
	}

	switch {
	case locationConfig.S3Config != nil:
		mySQLReq.LocationConfig = &agentpb.StartJobRequest_MySQLBackup_S3Config{
			S3Config: &agentpb.S3LocationConfig{
				Endpoint:   locationConfig.S3Config.Endpoint,
				AccessKey:  locationConfig.S3Config.AccessKey,
				SecretKey:  locationConfig.S3Config.SecretKey,
				BucketName: locationConfig.S3Config.BucketName,
				// TODO: replace hardcoded region with real value https://jira.percona.com/browse/PMM-7832
				BucketRegion: "us-east-2",
			},
		}
	default:
		return errors.Errorf("unsupported location config")
	}
	req := &agentpb.StartJobRequest{
		JobId:   id,
		Timeout: ptypes.DurationProto(timeout),
		Job: &agentpb.StartJobRequest_MysqlBackup{
			MysqlBackup: mySQLReq,
		},
	}

	agent, err := s.r.get(pmmAgentID)
	if err != nil {
		return err
	}

	resp := agent.channel.SendAndWaitResponse(req)
	if e := resp.(*agentpb.StartJobResponse).Error; e != "" {
		return errors.Errorf("failed to start MySQL job: %s", e)
	}

	return nil
}

// StopJob stops job with given given id.
func (s *JobsService) StopJob(jobID string) error {
	jobResult, err := models.FindJobResultByID(s.db.Querier, jobID)
	if err != nil {
		return errors.WithStack(err)
	}

	if jobResult.Done {
		// Job already finished
		return nil
	}

	agent, err := s.r.get(jobResult.PMMAgentID)
	if err != nil {
		return errors.WithStack(err)
	}

	agent.channel.SendAndWaitResponse(&agentpb.StopJobRequest{JobId: jobID})

	return nil
}
