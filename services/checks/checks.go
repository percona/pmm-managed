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
	"fmt"
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
	defaultInterval = 10 * time.Second // TODO Debug value

	// Environment variables that affect checks service; only for testing.
	envHost      = "PERCONA_TEST_CHECKS_HOST"
	envPublicKey = "PERCONA_TEST_CHECKS_PUBLIC_KEY"
	envInterval  = "PERCONA_TEST_CHECKS_INTERVAL"
	envCheckFile = "PERCONA_TEST_CHECKS_FILE"

	timeout = 5 * time.Second
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

	tm    sync.Mutex
	tasks map[string]checkTask

	registry registryService
	db       *reform.DB
}

type checkTask struct {
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
		tasks:      make(map[string]checkTask),
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
	time.Sleep(5 * time.Second) // TODO to let agents connect/reconnect
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	go s.checkResults(ctx)

	for {
		if f := os.Getenv(envCheckFile); f != "" {
			s.l.Warnf("Use local test checks file: %s", f)
			if err := s.loadLocalChecks(f); err != nil {
				s.l.Error("Failed to load local checks file: %s.", err)
			}
		} else {
			if err := s.downloadChecks(ctx); err != nil {
				s.l.Error("Failed to download checks, %s.", err)
			}
		}

		if err := s.executeChecks(); err != nil {
			s.l.Error("Failed to execute security checks: %s.", err)
		}

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

func (s *Service) checkResults(ctx context.Context) {
	tic := time.NewTicker(30 * time.Second) // TODO Scan interval

	for {
		select {
		case <-tic.C:
			// continue with next loop iteration
		case <-ctx.Done():
			return
		}

		for id := range s.getResults() {
			res, err := models.FindActionResultByID(s.db.Querier, id)
			if err != nil {
				s.l.Error(err)
			}

			if !res.Done {
				continue
			}

			if res.Error != "" {
				s.l.Warn("Action %s failed: %s", id, res.Error) // TODO better log message
				s.removeResult(id)
				continue
			}

			rr, err := agentpb.UnmarshalActionQueryResult([]byte(res.Output))
			if err != nil {
				s.l.Error(err) // TODO log message
			}

			// TODO Execute script against returned data
			// TODO Throw away results with expired TTL (they probably should have TTL)
			fmt.Println(rr)
			s.removeResult(id)
		}

	}
}

func (s *Service) addResults(res []checkTask) {
	s.tm.Lock()
	defer s.tm.Unlock()

	for _, r := range res {
		s.tasks[r.resultID] = r
	}
}

func (s *Service) removeResult(id string) {
	s.tm.Lock()
	defer s.tm.Unlock()

	delete(s.tasks, id)
}

func (s *Service) getResults() map[string]checkTask {
	s.tm.Lock()
	defer s.tm.Unlock()

	res := make(map[string]checkTask, len(s.tasks))
	for k, v := range s.tasks {
		res[k] = v
	}

	return res
}

func (s *Service) executeChecks() error {
	mySQLChecks, _, _ := groupChecksByDB(s.checks)

	mySQLRes, err := s.executeMySQLChecks(mySQLChecks)
	if err != nil {
		s.l.Error("Failed to execute mySQL checks: %s.", err)
	}
	s.addResults(mySQLRes)
	return nil
}

func (s *Service) executeMySQLChecks(checks []check.Check) ([]checkTask, error) {
	var res []checkTask
	var agents []*models.Agent
	var services []*models.Service

	e := s.db.InTransaction(func(t *reform.TX) error {
		var err error
		typ := models.MySQLdExporterType
		if agents, err = models.FindAgents(s.db.Querier, models.AgentFilters{AgentType: &typ}); err != nil {
			return err
		}
		if services, err = models.FindServicesByIDs(s.db.Querier, getServicesIDs(agents)); err != nil {
			return err
		}
		return nil
	})
	if e != nil {
		return nil, e
	}

	sMap := servicesToMap(services)
	for _, agent := range agents {
		r, err := models.CreateActionResult(s.db.Querier, *agent.PMMAgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", *agent.PMMAgentID, err)
			continue
		}
		dsn := agent.DSN(sMap[*agent.ServiceID], 2*time.Second, "") // TODO Do we need DB name for some checks?

		for _, c := range checks {
			switch c.Type {
			case check.MySQLShow:
				if err := s.registry.StartMySQLQueryShowAction(context.Background(), r.ID, *agent.PMMAgentID, dsn, c.Query); err != nil {
					s.l.Error("Failed to start mySQL action: %s.", err)
					continue
				}
				res = append(res, checkTask{
					resultID:   r.ID,
					pmmAgentID: *agent.PMMAgentID,
					serviceID:  *agent.ServiceID,
					check:      &c})
			}
		}
	}

	return res, nil
}

func groupChecksByDB(checks []check.Check) ([]check.Check, []check.Check, []check.Check) {
	var mySQLChecks []check.Check
	var postgreSQLChecks []check.Check
	var mongoChecks []check.Check

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
			mongoChecks = append(mongoChecks, c)
		}
	}

	return mySQLChecks, postgreSQLChecks, mongoChecks
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
		grpc.WithBackoffMaxDelay(timeout), //nolint:staticcheck
		grpc.WithBlock(),
		grpc.WithUserAgent("pmm-managed/" + s.pmmVersion),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
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
