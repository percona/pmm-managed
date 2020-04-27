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

// Package checks provides security checks functionality.
package checks

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	api "github.com/percona-platform/saas/gen/checked"
	"github.com/percona-platform/saas/pkg/check"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/utils/tlsconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

const (
	defaultHost     = "check.percona.com:443"
	defaultInterval = 24 * time.Hour

	// Environment variables that affect checks service; only for testing.
	envHost      = "PERCONA_TEST_CHECKS_HOST"
	envPublicKey = "PERCONA_TEST_CHECKS_PUBLIC_KEY"
	envInterval  = "PERCONA_TEST_CHECKS_INTERVAL"
	envCheckFile = "PERCONA_TEST_CHECKS_FILE"

	downloadTimeout   = 10 * time.Second
	actionDialTimeout = 5 * time.Second
	resultsTimeout    = time.Minute
)

var defaultPublicKeys = []string{
	"RWSKCHyoLDYxJ1k0qeayKu3/fsXVS1z8M+0deAClryiHWP99Sr4R/gPP", // PMM 2.6
}

// Service is responsible for interactions with Percona Check service.
type Service struct {
	l          *logrus.Entry
	pmmVersion string
	host       string
	publicKeys []string
	interval   time.Duration

	cm     sync.Mutex
	checks []check.Check

	registry registryService
	db       *reform.DB
}

type task struct {
	resultID   string
	pmmAgentID string
	serviceID  string
	check      *check.Check
}

// New returns Service with given PMM version.
func New(registry registryService, db *reform.DB, pmmVersion string) *Service {
	l := logrus.WithField("component", "check")
	s := &Service{
		l:          l,
		pmmVersion: pmmVersion,
		host:       defaultHost,
		publicKeys: defaultPublicKeys,
		interval:   defaultInterval,
		registry:   registry,
		db:         db,
	}

	if h := os.Getenv(envHost); h != "" {
		l.Warnf("Host changed to %s.", h)
		s.host = h
	}
	if k := os.Getenv(envPublicKey); k != "" {
		s.publicKeys = strings.Split(k, ",")
		l.Warnf("Public keys changed to %q.", k)
	}
	if d, err := time.ParseDuration(os.Getenv(envInterval)); err == nil && d > 0 {
		l.Warnf("Interval changed to %s.", d)
		s.interval = d
	}

	return s
}

// Run runs checks service that grabs checks from Percona Checks service every interval until context is canceled.
func (s *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		s.grabChecks(ctx)
		s.executeChecks()

		select {
		case <-ticker.C:
			// continue with next loop iteration
		case <-ctx.Done():
			return
		}
	}
}

// Checks returns available checks.
func (s *Service) Checks() []check.Check {
	s.cm.Lock()
	defer s.cm.Unlock()

	r := make([]check.Check, 0, len(s.checks))
	return append(r, s.checks...)
}

func (s *Service) processTasks(ctx context.Context, tasks []task) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		if len(tasks) == 0 {
			break
		}

		var retry []task
		select {
		case <-ticker.C:
			// continue with next loop iteration
		case <-ctx.Done():
			return
		}

		for _, t := range tasks {
			res, err := models.FindActionResultByID(s.db.Querier, t.resultID)
			if err != nil {
				s.l.Errorf("Can't find action result: %s.", err)
				continue
			}

			if !res.Done {
				retry = append(retry, t)
				continue
			}

			if res.Error != "" {
				s.l.Errorf("Action %s failed: %s.", t.resultID, res.Error)
				continue
			}

			_, err = agentpb.UnmarshalActionQueryResult([]byte(res.Output))
			if err != nil {
				s.l.Errorf("Failed to parse action result with id: %s, reason: %s.", t.resultID, err)
				continue
			}

			// TODO Execute script against returned data
			// fmt.Println(rr)
		}

		tasks = retry
	}
}

func (s *Service) executeChecks() {
	mySQLChecks, postgreSQLChecks, mongoDBChecks := s.groupChecksByDB(s.checks)

	var tasks []task
	mySQLTasks, err := s.executeMySQLChecks(mySQLChecks)
	if err != nil {
		s.l.Errorf("Failed to execute MySQL checks: %s.", err)
	}
	tasks = append(tasks, mySQLTasks...)

	postgreSQLTasks, err := s.executePostgreSQLChecks(postgreSQLChecks)
	if err != nil {
		s.l.Errorf("Failed to execute PostgreSQL checks: %s.", err)
	}
	tasks = append(tasks, postgreSQLTasks...)

	mongoDBTasks, err := s.executeMongoChecks(mongoDBChecks)
	if err != nil {
		s.l.Errorf("Failed to execute MongoDB checks: %s.", err)
	}
	tasks = append(tasks, mongoDBTasks...)

	ctx, cancel := context.WithTimeout(context.Background(), resultsTimeout)
	defer cancel()

	s.processTasks(ctx, tasks)
}

func (s *Service) executeMySQLChecks(checks []check.Check) ([]task, error) {
	var res []task

	agents, services, err := s.findAgentsAndServices(models.MySQLdExporterType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find proper agents and services")
	}

	for _, agent := range agents {
		pmmAgentID := *agent.PMMAgentID
		r, err := models.CreateActionResult(s.db.Querier, pmmAgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", pmmAgentID, err)
			continue
		}
		dsn := agent.DSN(services[*agent.ServiceID], actionDialTimeout, "")

		for _, c := range checks {
			switch c.Type {
			case check.MySQLShow:
				if err := s.registry.StartMySQLQueryShowAction(context.Background(), r.ID, pmmAgentID, dsn, c.Query); err != nil {
					s.l.Errorf("Failed to start MySQL show query action for agent %s, reason: %s.", pmmAgentID, err)
					continue
				}
			case check.MySQLSelect:
				if err := s.registry.StartMySQLQuerySelectAction(context.Background(), r.ID, pmmAgentID, dsn, c.Query); err != nil {
					s.l.Errorf("Failed to start MySQL select query action for agent %s, reason: %s.", pmmAgentID, err)
					continue
				}
			default:
				s.l.Errorf("Unknown MySQL check type: %s.", c.Type)
				continue
			}

			res = append(res, task{
				resultID:   r.ID,
				pmmAgentID: pmmAgentID,
				serviceID:  *agent.ServiceID,
				check:      &c,
			})
		}
	}

	return res, nil
}

func (s *Service) executePostgreSQLChecks(checks []check.Check) ([]task, error) {
	var res []task

	agents, services, err := s.findAgentsAndServices(models.PostgresExporterType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find proper agents and services")
	}

	for _, agent := range agents {
		pmmAgentID := *agent.PMMAgentID
		r, err := models.CreateActionResult(s.db.Querier, pmmAgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", pmmAgentID, err)
			continue
		}
		dsn := agent.DSN(services[*agent.ServiceID], actionDialTimeout, "")

		for _, c := range checks {
			switch c.Type {
			case check.PostgreSQLShow:
				if err := s.registry.StartPostgreSQLQueryShowAction(context.Background(), r.ID, pmmAgentID, dsn); err != nil {
					s.l.Errorf("Failed to start PostgreSQL show query action for agent %s, reason: %s.", pmmAgentID, err)
					continue
				}
			case check.PostgreSQLSelect:
				if err := s.registry.StartPostgreSQLQuerySelectAction(context.Background(), r.ID, pmmAgentID, dsn, c.Query); err != nil {
					s.l.Errorf("Failed to start PostgreSQL select query action for agent %s, reason: %s.", pmmAgentID, err)
					continue
				}
			default:
				s.l.Errorf("Unknown PostgresSQL check type: %s.", c.Type)
				continue
			}
			res = append(res, task{
				resultID:   r.ID,
				pmmAgentID: pmmAgentID,
				serviceID:  *agent.ServiceID,
				check:      &c,
			})
		}
	}

	return res, nil
}

func (s *Service) executeMongoChecks(checks []check.Check) ([]task, error) {
	var res []task

	agents, services, err := s.findAgentsAndServices(models.MongoDBExporterType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find proper agents and services")
	}

	for _, agent := range agents {
		pmmAgentID := *agent.PMMAgentID
		r, err := models.CreateActionResult(s.db.Querier, pmmAgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", pmmAgentID, err)
			continue
		}
		dsn := agent.DSN(services[*agent.ServiceID], actionDialTimeout, "")

		for _, c := range checks {
			switch c.Type {
			case check.MongoDBGetParameter:
				if err := s.registry.StartMongoDBQueryGetParameterAction(context.Background(), r.ID, pmmAgentID, dsn); err != nil {
					s.l.Errorf("Failed to start MongoDB get parameter query action for agent %s, reason: %s.", pmmAgentID, err)
					continue
				}
			case check.MongoDBBuildInfo:
				if err := s.registry.StartMongoDBQueryBuildInfoAction(context.Background(), r.ID, pmmAgentID, dsn); err != nil {
					s.l.Errorf("Failed to start MongoDB build info query action for agent %s, reason: %s.", pmmAgentID, err)
					continue
				}

			default:
				s.l.Errorf("Unknown MongoDB check type: %s.", c.Type)
				continue
			}
			res = append(res, task{
				resultID:   r.ID,
				pmmAgentID: pmmAgentID,
				serviceID:  *agent.ServiceID,
				check:      &c,
			})
		}
	}

	return res, nil
}

func (s *Service) findAgentsAndServices(agentType models.AgentType) ([]*models.Agent, map[string]*models.Service, error) {
	var agents []*models.Agent
	var services []*models.Service

	e := s.db.InTransaction(func(t *reform.TX) error {
		var err error
		if agents, err = models.FindAgents(s.db.Querier, models.AgentFilters{AgentType: &agentType}); err != nil {
			return err
		}
		if services, err = models.FindServicesByIDs(s.db.Querier, getServicesIDs(agents)); err != nil {
			return err
		}
		return nil
	})
	if e != nil {
		return nil, nil, e
	}

	return agents, servicesToMap(services), nil
}

func (s *Service) groupChecksByDB(checks []check.Check) ([]check.Check, []check.Check, []check.Check) {
	var mySQLChecks, postgreSQLChecks, mongoChecks []check.Check

	for _, c := range checks {
		switch c.Type {
		case check.MySQLSelect:
			fallthrough
		case check.MySQLShow:
			mySQLChecks = append(mySQLChecks, c)

		case check.PostgreSQLSelect:
			fallthrough
		case check.PostgreSQLShow:
			postgreSQLChecks = append(postgreSQLChecks, c)

		case check.MongoDBGetParameter:
			fallthrough
		case check.MongoDBBuildInfo:
			mongoChecks = append(mongoChecks, c)

		default:
			s.l.Warnf("Unknown check type %s, skip it.", c.Type)
		}
	}

	return mySQLChecks, postgreSQLChecks, mongoChecks
}

func (s *Service) grabChecks(ctx context.Context) {
	if f := os.Getenv(envCheckFile); f != "" {
		s.l.Warnf("Use local test checks file: %s.", f)
		if err := s.loadLocalChecks(f); err != nil {
			s.l.Errorf("Failed to load local checks file: %s.", err)
		}
	} else {
		if err := s.downloadChecks(ctx); err != nil {
			s.l.Errorf("Failed to download checks, %s.", err)
		}
	}
}

func (s *Service) loadLocalChecks(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read test checks file")
	}
	checks, err := check.Parse(bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "failed to parse test checks file")
	}

	s.updateChecks(checks)
	return nil
}

func (s *Service) downloadChecks(ctx context.Context) error {
	s.l.Infof("Downloading checks from %s ...", s.host)

	host, _, err := net.SplitHostPort(s.host)
	if err != nil {
		return errors.Wrap(err, "failed to set checks host")
	}
	tlsConfig := tlsconfig.Get()
	tlsConfig.ServerName = host

	opts := []grpc.DialOption{
		// replacement is marked as experimental
		grpc.WithBackoffMaxDelay(downloadTimeout), //nolint:staticcheck

		grpc.WithBlock(),
		grpc.WithUserAgent("pmm-managed/" + s.pmmVersion),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	defer cancel()
	cc, err := grpc.DialContext(ctx, s.host, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}
	defer cc.Close() //nolint:errcheck

	resp, err := api.NewCheckedAPIClient(cc).GetAllChecks(ctx, &api.GetAllChecksRequest{})
	if err != nil {
		return errors.Wrap(err, "failed to request checks service")
	}

	if err = s.verifySignatures(resp); err != nil {
		return err
	}

	checks, err := check.Parse(strings.NewReader(resp.File))
	if err != nil {
		return err
	}

	s.updateChecks(checks)
	return nil
}

func (s *Service) updateChecks(checks []check.Check) {
	s.cm.Lock()
	defer s.cm.Unlock()

	s.checks = checks
}

func (s *Service) verifySignatures(resp *api.GetAllChecksResponse) error {
	if len(resp.Signatures) == 0 {
		return errors.New("zero signatures received")
	}

	var err error
	for _, sign := range resp.Signatures {
		for _, key := range s.publicKeys {
			if err = check.Verify([]byte(resp.File), key, sign); err == nil {
				return nil
			}
			s.l.Debugf("Key %q doesn't match signature %q: %s.", key, sign, err)
		}
	}

	return errors.New("no verified signatures")
}

func servicesToMap(services []*models.Service) map[string]*models.Service {
	res := make(map[string]*models.Service, len(services))
	for _, service := range services {
		res[service.ServiceID] = service
	}

	return res
}

func getServicesIDs(agents []*models.Agent) []string {
	res := make([]string, len(agents))
	for i, agent := range agents {
		res[i] = *agent.ServiceID
	}

	return res
}
