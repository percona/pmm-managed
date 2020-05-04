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

	api "github.com/percona-platform/saas/gen/check/retrieval"
	"github.com/percona-platform/saas/pkg/check"
	"github.com/percona-platform/saas/pkg/starlark"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/percona/pmm/utils/tlsconfig"
	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
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

	checksTimeout       = time.Hour
	downloadTimeout     = 10 * time.Second
	resultTimeout       = 15 * time.Second
	resultCheckInterval = time.Second

	prometheusNamespace = "pmm_managed"
	prometheusSubsystem = "checks"
)

var defaultPublicKeys = []string{
	"RWTfyQTP3R7VzZggYY7dzuCbuCQWqTiGCqOvWRRAMVEiw0eSxHMVBBE5", // PMM 2.6
}

// Service is responsible for interactions with Percona Check service.
type Service struct {
	agentsRegistry agentsRegistry
	alertsRegistry alertRegistry
	db             *reform.DB
	pmmVersion     string

	l          *logrus.Entry
	host       string
	publicKeys []string
	interval   time.Duration

	cm     sync.Mutex
	checks []check.Check
}

// New returns Service with given PMM version.
func New(agentsRegistry agentsRegistry, alertsRegistry alertRegistry, db *reform.DB, pmmVersion string) *Service {
	l := logrus.WithField("component", "checks")
	s := &Service{
		agentsRegistry: agentsRegistry,
		alertsRegistry: alertsRegistry,
		db:             db,
		pmmVersion:     pmmVersion,

		l:          l,
		host:       defaultHost,
		publicKeys: defaultPublicKeys,
		interval:   defaultInterval,
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
		var sttEnabled bool
		settings, err := models.GetSettings(s.db)
		if err != nil {
			s.l.Error(err)
		}
		if settings != nil && settings.SaaS.STTEnabled {
			sttEnabled = true
		}

		if sttEnabled {
			nCtx, cancel := context.WithTimeout(ctx, checksTimeout)
			s.grabChecks(nCtx)
			s.executeChecks(nCtx)
			cancel()
		} else {
			s.l.Info("STT is not enabled, doing nothing.")
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

// waitForResult periodically checks result state and returns it when complete.
func (s *Service) waitForResult(ctx context.Context, resultID string) ([]map[string]interface{}, error) {
	ticker := time.NewTicker(resultCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil, errors.WithStack(ctx.Err())
		}

		res, err := models.FindActionResultByID(s.db.Querier, resultID)
		if err != nil {
			return nil, err
		}

		if !res.Done {
			continue
		}

		// FIXME we still need to delete old result - they may never be done: https://jira.percona.com/browse/PMM-5840
		if err = s.db.Delete(res); err != nil {
			s.l.Warnf("Failed to delete action result %s: %s.", resultID, err)
		}

		if res.Error != "" {
			return nil, errors.Errorf("Action %s failed: %s.", resultID, res.Error)
		}

		out, err := agentpb.UnmarshalActionQueryResult([]byte(res.Output))
		if err != nil {
			return nil, errors.Errorf("Failed to parse action result : %s.", err)
		}

		return out, nil
	}
}

// executeChecks runs all available checks for all reachable services.
func (s *Service) executeChecks(ctx context.Context) {
	mySQLChecks, postgreSQLChecks, mongoDBChecks := s.groupChecksByDB(s.checks)

	if err := s.executeMySQLChecks(ctx, mySQLChecks); err != nil {
		s.l.Errorf("Failed to execute MySQL checks: %s.", err)
	}

	if err := s.executePostgreSQLChecks(ctx, postgreSQLChecks); err != nil {
		s.l.Errorf("Failed to execute PostgreSQL checks: %s.", err)
	}

	if err := s.executeMongoDBChecks(ctx, mongoDBChecks); err != nil {
		s.l.Errorf("Failed to execute MongoDB checks: %s.", err)
	}
}

// executeMySQLChecks runs MySQL checks for available MySQL services.
func (s *Service) executeMySQLChecks(ctx context.Context, checks []check.Check) error {
	targets, err := s.findTargets(models.MySQLServiceType)
	if err != nil {
		return errors.Wrap(err, "failed to find proper agents and services")
	}

	for _, target := range targets {
		r, err := models.CreateActionResult(s.db.Querier, target.agent.AgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", target.agent.AgentID, err)
			continue
		}

		for _, c := range checks {
			switch c.Type {
			case check.MySQLShow:
				if err := s.agentsRegistry.StartMySQLQueryShowAction(ctx, r.ID, target.agent.AgentID, target.dsn, c.Query); err != nil {
					s.l.Errorf("Failed to start MySQL show query action for agent %s, reason: %s.", target.agent.AgentID, err)
					continue
				}
			case check.MySQLSelect:
				if err := s.agentsRegistry.StartMySQLQuerySelectAction(ctx, r.ID, target.agent.AgentID, target.dsn, c.Query); err != nil {
					s.l.Errorf("Failed to start MySQL select query action for agent %s, reason: %s.", target.agent.AgentID, err)
					continue
				}
			default:
				s.l.Errorf("Unknown MySQL check type: %s.", c.Type)
				continue
			}

			if err = s.processResults(ctx, c, target, r.ID); err != nil {
				s.l.Errorf("failed to process action result: %s", err)
			}
		}
	}

	return nil
}

// executePostgreSQLChecks runs PostgreSQL checks for available PostgreSQL services.
func (s *Service) executePostgreSQLChecks(ctx context.Context, checks []check.Check) error {
	targets, err := s.findTargets(models.PostgreSQLServiceType)
	if err != nil {
		return errors.Wrap(err, "failed to find proper agents and services")
	}

	for _, target := range targets {
		r, err := models.CreateActionResult(s.db.Querier, target.agent.AgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", target.agent.AgentID, err)
			continue
		}

		for _, c := range checks {
			switch c.Type {
			case check.PostgreSQLShow:
				if err := s.agentsRegistry.StartPostgreSQLQueryShowAction(ctx, r.ID, target.agent.AgentID, target.dsn); err != nil {
					s.l.Errorf("Failed to start PostgreSQL show query action for agent %s, reason: %s.", target.agent.AgentID, err)
					continue
				}
			case check.PostgreSQLSelect:
				if err := s.agentsRegistry.StartPostgreSQLQuerySelectAction(ctx, r.ID, target.agent.AgentID, target.dsn, c.Query); err != nil {
					s.l.Errorf("Failed to start PostgreSQL select query action for agent %s, reason: %s.", target.agent.AgentID, err)
					continue
				}
			default:
				s.l.Errorf("Unknown PostgresSQL check type: %s.", c.Type)
				continue
			}

			if err = s.processResults(ctx, c, target, r.ID); err != nil {
				s.l.Errorf("failed to process action result: %s", err)
			}
		}
	}

	return nil
}

// executeMongoDBChecks runs MongoDB checks for available MongoDB services.
func (s *Service) executeMongoDBChecks(ctx context.Context, checks []check.Check) error {
	targets, err := s.findTargets(models.MongoDBServiceType)
	if err != nil {
		return errors.Wrap(err, "failed to find proper agents and services")
	}

	for _, target := range targets {
		r, err := models.CreateActionResult(s.db.Querier, target.agent.AgentID)
		if err != nil {
			s.l.Errorf("Failed to prepare action result for agent %s: %s.", target.agent.AgentID, err)
			continue
		}

		for _, c := range checks {
			switch c.Type {
			case check.MongoDBGetParameter:
				if err := s.agentsRegistry.StartMongoDBQueryGetParameterAction(context.Background(), r.ID, target.agent.AgentID, target.dsn); err != nil {
					s.l.Errorf("Failed to start MongoDB get parameter query action for agent %s, reason: %s.", target.agent.AgentID, err)
					continue
				}
			case check.MongoDBBuildInfo:
				if err := s.agentsRegistry.StartMongoDBQueryBuildInfoAction(context.Background(), r.ID, target.agent.AgentID, target.dsn); err != nil {
					s.l.Errorf("Failed to start MongoDB build info query action for agent %s, reason: %s.", target.agent.AgentID, err)
					continue
				}

			default:
				s.l.Errorf("Unknown MongoDB check type: %s.", c.Type)
				continue
			}

			if err = s.processResults(ctx, c, target, r.ID); err != nil {
				s.l.Errorf("failed to process action result: %s", err)
			}
		}
	}

	return nil
}

// TODO find better name
func (s *Service) processResults(ctx context.Context, check check.Check, target target, resID string) error {
	nCtx, cancel := context.WithTimeout(ctx, resultTimeout)
	r, err := s.waitForResult(nCtx, resID)
	cancel()
	if err != nil {
		return errors.Wrap(err, "failed to get check result")
	}

	funcs, err := getFuncsForVersion(check.Version)
	if err != nil {
		return err
	}

	env, err := starlark.NewEnv(check.Name, check.Script, funcs)
	if err != nil {
		return errors.Wrap(err, "failed to prepare starlark environment")
	}

	results, err := env.Run(check.Name, r, s.l.Debug)
	if err != nil {
		return errors.Wrap(err, "failed to execute script")
	}

	for _, result := range results {
		alert, err := makeAlert(check.Name, target, &result)
		if err != nil {
			return errors.Wrap(err, "failed to create alert")
		}

		s.alertsRegistry.Add(check.Name, time.Second, alert)
	}

	return nil
}

func makeAlert(name string, target target, result *check.Result) (*ammodels.PostableAlert, error) {
	labels, err := models.MergeLabels(target.node, target.service, target.agent)
	if err != nil {
		return nil, err
	}

	labels[model.AlertNameLabel] = name
	labels["severity"] = result.Severity.String()
	labels["stt_check"] = "1"

	return &ammodels.PostableAlert{
		Alert: ammodels.Alert{
			// GeneratorURL: "TODO",
			Labels: labels,
		},

		// StartsAt and EndAt can't be added there without changes in Registry

		Annotations: map[string]string{
			"summary":     result.Summary,
			"description": result.Description,
		},
	}, nil
}

// target contains required info about check target
type target struct {
	agent   *models.Agent
	service *models.Service
	node    *models.Node
	dsn     string
}

// findTargets returns slice of available targets for specified service type.
func (s *Service) findTargets(serviceType models.ServiceType) ([]target, error) {
	var targets []target
	services, err := models.FindServices(s.db.Querier, models.ServiceFilters{ServiceType: &serviceType})
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		e := s.db.InTransaction(func(tx *reform.TX) error {
			a, err := models.FindPMMAgentsForService(s.db.Querier, service.ServiceID)
			if err != nil {
				return err
			}
			if len(a) == 0 {
				return errors.New("no available pmm agents")
			}

			dsn, err := models.FindDSNByServiceIDandPMMAgentID(s.db.Querier, service.ServiceID, a[0].AgentID, "")
			if err != nil {
				return err
			}

			n, err := models.FindNodeByID(s.db.Querier, service.NodeID)
			if err != nil {
				return err
			}

			targets = append(targets,
				target{
					agent:   a[0],
					service: service,
					node:    n,
					dsn:     dsn,
				},
			)
			return nil
		})
		if e != nil {
			s.l.Errorf("Failed to find agents for service %s, reason: %s", service.ServiceID, err)
		}
	}

	return targets, nil
}

// groupChecksByDB splits provided checks by database and returns three slices: for MySQL, for PostgreSQL and for MongoDB.
func (s *Service) groupChecksByDB(checks []check.Check) ([]check.Check, []check.Check, []check.Check) {
	var mySQLChecks, postgreSQLChecks, mongoDBChecks []check.Check

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
			mongoDBChecks = append(mongoDBChecks, c)

		default:
			s.l.Warnf("Unknown check type %s, skip it.", c.Type)
		}
	}

	return mySQLChecks, postgreSQLChecks, mongoDBChecks
}

// grabChecks loads checks list.
func (s *Service) grabChecks(ctx context.Context) {
	// TODO once checks are collected, filter out unknown versions
	// collectCheck(ctx context.Context) ([]check.Check, error) might be a good function signature

	if f := os.Getenv(envCheckFile); f != "" {
		s.l.Warnf("Using local test checks file: %s.", f)
		if err := s.loadLocalChecks(f); err != nil {
			s.l.Errorf("Failed to load local checks file: %s.", err)
		}
		return
	}

	if err := s.downloadChecks(ctx); err != nil {
		s.l.Errorf("Failed to download checks: %s.", err)
	}
}

// loadLocalCheck loads checks form local file.
func (s *Service) loadLocalChecks(file string) error {
	data, err := ioutil.ReadFile(file) //nolint:gosec
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

// downloadChecks downloads checks form percona service endpoint.
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

	resp, err := api.NewRetrievalAPIClient(cc).GetAllChecks(ctx, &api.GetAllChecksRequest{})
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

// updateChecks update service checks filed value under mutex.
func (s *Service) updateChecks(checks []check.Check) {
	s.cm.Lock()
	defer s.cm.Unlock()

	s.checks = checks
}

// verifySignatures verifies checks signatures and returns error in case of verification problem.
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

// Describe implements prom.Collector.
func (s *Service) Describe(ch chan<- *prom.Desc) {
	// TODO
}

// Collect implements prom.Collector.
func (s *Service) Collect(ch chan<- prom.Metric) {
	// TODO
}

// check interfaces
var (
	_ prom.Collector = (*Service)(nil)
)
