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

// Package telemetry provides telemetry functionality.
package telemetry

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	events "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
	"github.com/percona/pmm/api/serverpb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto" //nolint:staticcheck
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/envvars"
	"github.com/percona/pmm-managed/utils/saasreq"
)

// Service is responsible for interactions with Percona Check / Telemetry service.
type Service struct {
	config     Config
	db         *reform.DB
	pmmVersion string
	start      time.Time
	l          *logrus.Entry

	os                  string
	sDistributionMethod serverpb.DistributionMethod
	tDistributionMethod events.DistributionMethod
}

// NewService creates a new service with given UUID and PMM version.
func NewService(db *reform.DB, pmmVersion string, config Config) (*Service, error) {
	l := logrus.WithField("component", "telemetry")

	if config.SaasHostname == "" {
		host, err := envvars.GetSAASHost()
		if err != nil {
			return nil, err
		}
		config.SaasHostname = host
	}

	s := &Service{
		config:     config,
		db:         db,
		pmmVersion: pmmVersion,
		start:      time.Now(),
		l:          l,
	}

	s.sDistributionMethod, s.tDistributionMethod, s.os = getDistributionMethodAndOS(l)

	if u := os.Getenv(s.config.V1URLEnv); u != "" {
		l.Warnf("v1URL changed to %q.", u)
		s.config.V1URL = u
	}

	if d, err := time.ParseDuration(os.Getenv(s.config.Reporting.IntervalEnv)); err == nil && d > 0 {
		l.Warnf("Interval changed to %s.", d)
		s.config.Reporting.Interval = d
	}
	if d, err := time.ParseDuration(os.Getenv(s.config.Reporting.RetryBackoffEnv)); err == nil && d > 0 {
		l.Warnf("Retry backoff changed to %s.", d)
		s.config.Reporting.RetryBackoff = d
	}

	s.l.Debugf("Telemetry settings: os=%q, sDistributionMethod=%q, tDistributionMethod=%q.",
		s.os, s.sDistributionMethod, s.tDistributionMethod)

	return s, nil
}

func getDistributionMethodAndOS(l *logrus.Entry) (serverpb.DistributionMethod, events.DistributionMethod, string) {
	b, err := ioutil.ReadFile(distributionInfoFilePath)
	if err != nil {
		l.Debugf("Failed to read %s: %s", distributionInfoFilePath, err)
	}

	b = bytes.ToLower(bytes.TrimSpace(b))
	switch string(b) {
	case "ovf":
		return serverpb.DistributionMethod_OVF, events.DistributionMethod_OVF, "ovf"
	case "ami":
		return serverpb.DistributionMethod_AMI, events.DistributionMethod_AMI, "ami"
	case "azure":
		return serverpb.DistributionMethod_AZURE, events.DistributionMethod_AZURE, "azure"
	case "digitalocean":
		return serverpb.DistributionMethod_DO, events.DistributionMethod_DO, "digitalocean"
	case "docker", "": // /srv/pmm-distribution does not exist in PMM 2.0.
		if b, err = ioutil.ReadFile(osInfoFilePath); err != nil {
			l.Debugf("Failed to read %s: %s", osInfoFilePath, err)
		}
		return serverpb.DistributionMethod_DOCKER, events.DistributionMethod_DOCKER, getLinuxDistribution(string(b))
	default:
		return serverpb.DistributionMethod_DISTRIBUTION_METHOD_INVALID, events.DistributionMethod_DISTRIBUTION_METHOD_INVALID, ""
	}
}

// DistributionMethod returns PMM Server distribution method where this pmm-managed runs.
func (s *Service) DistributionMethod() serverpb.DistributionMethod {
	return s.sDistributionMethod
}

// Run runs telemetry service after delay, sending data every interval until context is canceled.
func (s *Service) Run(ctx context.Context) {
	if !s.config.Enabled {
		s.l.Warn("service is disabled, skip Run")
		return
	}

	ticker := time.NewTicker(s.config.Reporting.Interval)
	defer ticker.Stop()

	// delay the very first report too to let users opt-out
	for {
		select {
		case <-ticker.C:
			// continue with next loop iteration
		case <-ctx.Done():
			return
		}

		err := s.sendOneEvent(ctx)
		if err == nil {
			s.l.Debug("Telemetry info sent.")
		} else {
			s.l.Debugf("Telemetry info not sent: %s.", err)
		}
	}
}

func (s *Service) sendOneEvent(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Reporting.SendTimeout)
	defer cancel()

	var settings *models.Settings
	err := s.db.InTransaction(func(tx *reform.TX) error {
		var e error
		if settings, e = models.GetSettings(tx); e != nil {
			return e
		}

		if settings.Telemetry.Disabled {
			return errors.New("disabled via settings")
		}
		if settings.Telemetry.UUID == "" {
			settings.Telemetry.UUID, e = generateUUID()
			if e != nil {
				return e
			}
			return models.SaveSettings(tx, settings)
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.l.Debugf("Using %s as server UUID.", settings.Telemetry.UUID)

	var wg errgroup.Group

	wg.Go(func() error {
		payload := s.makeV1Payload(settings.Telemetry.UUID)
		return s.sendV1Request(ctx, payload)
	})

	wg.Go(func() error {
		req, err := s.makeV2Payload(settings.Telemetry.UUID, settings)
		if err != nil {
			return err
		}

		return s.sendV2RequestWithRetries(ctx, req)
	})

	return wg.Wait()
}

func (s *Service) makeV1Payload(uuid string) []byte {
	var w bytes.Buffer
	fmt.Fprintf(&w, "%s;%s;%s\n", uuid, "OS", s.os)
	fmt.Fprintf(&w, "%s;%s;%s\n", uuid, "PMM", s.pmmVersion)
	return w.Bytes()
}

func (s *Service) sendV1Request(ctx context.Context, data []byte) error {
	if s.config.V1URL == "" {
		return errors.New("v1 telemetry disabled via the empty URL")
	}

	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", s.config.V1URL, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("X-Percona-Toolkit-Tool", "pmm")

	ctx, cancel := context.WithTimeout(ctx, s.config.Reporting.SendTimeout)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) makeV2Payload(serverUUID string, settings *models.Settings) (*reporter.ReportRequest, error) {
	serverID, err := hex.DecodeString(serverUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode UUID %q", serverUUID)
	}

	event := &events.ServerUptimeEvent{
		Id:                 serverID,
		Version:            s.pmmVersion,
		UpDuration:         durationpb.New(time.Since(s.start)),
		DistributionMethod: s.tDistributionMethod,
		SttEnabled:         wrapperspb.Bool(settings.SaaS.STTEnabled),
		IaEnabled:          wrapperspb.Bool(settings.IntegratedAlerting.Enabled),
	}

	eventB, err := proto.Marshal(event)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal event %+v", event)
	}

	id := uuid.New()
	now := time.Now()
	req := &reporter.ReportRequest{
		Events: []*reporter.Event{{
			Id:   id[:],
			Time: timestamppb.New(now),
			Event: &reporter.AnyEvent{
				TypeUrl: string(event.ProtoReflect().Descriptor().FullName()),
				Binary:  eventB,
			},
		}},
	}
	s.l.Debugf("Request: %+v", req)

	return req, nil
}

func (s *Service) sendV2RequestWithRetries(ctx context.Context, req *reporter.ReportRequest) error {
	if s.config.SaasHostname == "" {
		return errors.New("v2 telemetry disabled via the empty host")
	}

	var err error
	var attempt int
	for {
		err = s.sendV2Request(ctx, req)
		attempt++
		s.l.Debugf("sendV2Request (attempt %d/%d) result: %v", attempt, s.config.Reporting.RetryCount, err)
		if err == nil {
			return nil
		}

		if attempt >= s.config.Reporting.RetryCount {
			s.l.Debug("Failed to send v2 event, will not retry (too much attempts).")
			return err
		}

		retryCtx, retryCancel := context.WithTimeout(ctx, s.config.Reporting.RetryBackoff)
		<-retryCtx.Done()
		retryCancel()

		if err = ctx.Err(); err != nil {
			s.l.Debugf("Will not retry sending v2 event: %s.", err)
			return err
		}
	}
}

func (s *Service) sendV2Request(ctx context.Context, req *reporter.ReportRequest) error {
	s.l.Debugf("Using %s as telemetry host.", s.config.SaasHostname)

	var accessToken string
	if ssoDetails, err, _ := models.GetPerconaSSODetails(ctx, s.db.Querier); err == nil {
		accessToken = ssoDetails.AccessToken.AccessToken
	}

	reqByte, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	_, err = saasreq.MakeRequest(ctx, http.MethodPost, s.config.ReportEndpointURL(), accessToken, bytes.NewReader(reqByte))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	return nil
}

func generateUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "can't generate UUID")
	}

	// Old telemetry IDs have only 32 chars in the table but UUIDs + "-" = 36
	cleanUUID := strings.Replace(uuid.String(), "-", "", -1)
	return cleanUUID, nil
}

// Currently, we only detect OS (Linux distribution) version from the kernel version (/proc/version).
// For both AMI and OVF images this value is fixed by the environment variable and not autodetected –
// we know OS for them because we make those images ourselves.
// If/when we decide to support installation with "normal" Linux package managers (apt, yum, etc.),
// we could use the code that was there. See PMM-1333 and PMM-1507 in both git logs and Jira for details.

type pair struct {
	re *regexp.Regexp
	t  string
}

var procVersionRegexps = []pair{
	{regexp.MustCompile(`ubuntu\d+~(?P<version>\d+\.\d+)`), "Ubuntu ${version}"},
	{regexp.MustCompile(`ubuntu`), "Ubuntu"},
	{regexp.MustCompile(`Debian`), "Debian"},
	{regexp.MustCompile(`\.fc(?P<version>\d+)\.`), "Fedora ${version}"},
	{regexp.MustCompile(`\.centos\.`), "CentOS"},
	{regexp.MustCompile(`\-ARCH`), "Arch"},
	{regexp.MustCompile(`\-moby`), "Moby"},
	{regexp.MustCompile(`\.amzn\d+\.`), "Amazon"},
	{regexp.MustCompile(`Microsoft`), "Microsoft"},
}

// getLinuxDistribution detects Linux distribution and version from /proc/version information.
func getLinuxDistribution(procVersion string) string {
	for _, p := range procVersionRegexps {
		match := p.re.FindStringSubmatchIndex(procVersion)
		if match != nil {
			return string(p.re.ExpandString(nil, p.t, procVersion, match))
		}
	}
	return "unknown"
}
