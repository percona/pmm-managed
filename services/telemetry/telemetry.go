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

// Package telemetry provides Call Home functionality.
package telemetry

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	pmmv1beta1 "github.com/Percona-Platform/saas/gen/telemetry/events/pmm"
	reporterv1beta1 "github.com/Percona-Platform/saas/gen/telemetry/reporter"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/percona/pmm/api/serverpb"
	"github.com/percona/pmm/utils/tlsconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

const (
	// FIXME
	interval     = 1 * time.Minute
	defaultV1URL = ""

	// interval      = 24 * time.Hour
	timeout = 5 * time.Second
	// defaultV1URL  = "https://v.percona.com/"
	defaultV2Host = "callhome-staging.percona.com:443" // protocol is always https

	// environment variables that affect telemetry service
	envV1URL  = "PERCONA_VERSION_CHECK_URL" // the same name as for the Toolkit
	envV2Host = "PERCONA_TELEMETRY_HOST"
)

// Service is responsible for interactions with Percona Call Home service.
type Service struct {
	db         *reform.DB
	pmmVersion string
	l          *logrus.Entry
	start      time.Time

	initOnce            sync.Once
	v1OS                string
	sDistributionMethod serverpb.DistributionMethod
	tDistributionMethod pmmv1beta1.DistributionMethod
	v1URL               string
	v2Host              string
}

// NewService creates a new service with given UUID and PMM version.
func NewService(db *reform.DB, pmmVersion string) *Service {
	return &Service{
		db:         db,
		pmmVersion: pmmVersion,
		l:          logrus.WithField("component", "telemetry"),
		start:      time.Now(),
	}
}

func (s *Service) init() {
	b, err := ioutil.ReadFile("/srv/pmm-distribution")
	if err != nil {
		s.l.Debugf("Failed to read /srv/pmm-distribution: %s", err)
	}

	b = bytes.ToLower(bytes.TrimSpace(b))
	switch string(b) {
	case "ovf":
		s.v1OS = "ovf"
		s.sDistributionMethod = serverpb.DistributionMethod_OVF
		s.tDistributionMethod = pmmv1beta1.DistributionMethod_OVF
	case "ami":
		s.v1OS = "ami"
		s.sDistributionMethod = serverpb.DistributionMethod_AMI
		s.tDistributionMethod = pmmv1beta1.DistributionMethod_AMI
	case "docker", "": // /srv/pmm-distribution does not exist in PMM 2.0.
		b, err = ioutil.ReadFile("/proc/version")
		if err != nil {
			s.l.Debugf("Failed to read /proc/version: %s", err)
		}
		s.v1OS = getLinuxDistribution(string(b))

		s.sDistributionMethod = serverpb.DistributionMethod_DOCKER
		s.tDistributionMethod = pmmv1beta1.DistributionMethod_DOCKER
	}

	s.v1URL = defaultV1URL
	if u := os.Getenv(envV1URL); u != "" {
		s.v1URL = u
	}

	s.v2Host = defaultV2Host
	if u := os.Getenv(envV2Host); u != "" {
		s.v2Host = u
	}

	s.l.Debugf("Telemetry settings: v1OS=%q, sDistributionMethod=%q, v1URL=%q, v2Host=%q.",
		s.v1OS, s.sDistributionMethod, s.v1URL, s.v2Host)
}

// DistributionMethod returns PMM Server distribution method where this pmm-managed runs.
func (s *Service) DistributionMethod() serverpb.DistributionMethod {
	s.initOnce.Do(s.init)
	return s.sDistributionMethod
}

// Run runs telemetry service, sending data every interval until context is canceled.
func (s *Service) Run(ctx context.Context) {
	s.initOnce.Do(s.init)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		sendOnceCtx, sendOnceCancel := context.WithTimeout(ctx, 5*time.Second)
		if err := s.sendOnce(sendOnceCtx); err != nil {
			s.l.Debugf("Telemetry info not send: %s.", err)
		}
		sendOnceCancel()

		select {
		case <-ticker.C:
			// continue with next loop iteration
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) sendOnce(ctx context.Context) error {
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

	var wg errgroup.Group

	wg.Go(func() error {
		payload := s.makeV1Payload(settings.Telemetry.UUID)
		return s.sendV1Request(ctx, payload)
	})

	wg.Go(func() error {
		req, err := s.makeV2Payload(settings.Telemetry.UUID)
		if err != nil {
			return err
		}
		err = s.sendV2Request(ctx, req)
		s.l.Debugf("sendV2Request: %+v", err)
		return err
	})

	return wg.Wait()
}

func (s *Service) makeV1Payload(uuid string) []byte {
	var w bytes.Buffer
	fmt.Fprintf(&w, "%s;%s;%s\n", uuid, "OS", s.v1OS)
	fmt.Fprintf(&w, "%s;%s;%s\n", uuid, "PMM", s.pmmVersion)
	return w.Bytes()
}

func (s *Service) sendV1Request(ctx context.Context, data []byte) error {
	if s.v1URL == "" {
		return errors.New("v1 telemetry disabled via the empty URL")
	}

	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", s.v1URL, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("X-Percona-Toolkit-Tool", "pmm")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) makeV2Payload(serverUUID string) (*reporterv1beta1.ReportRequest, error) {
	serverID, err := hex.DecodeString(serverUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode UUID %q", serverUUID)
	}

	event := &pmmv1beta1.ServerUptimeEvent{
		Id:                 serverID,
		Version:            s.pmmVersion,
		UpDuration:         ptypes.DurationProto(time.Since(s.start)),
		DistributionMethod: s.tDistributionMethod,
	}
	if err = event.Validate(); err != nil {
		// log and ignore
		s.l.Debugf("Failed to validate event: %s.", err)
	}
	eventB, err := proto.Marshal(event)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal event %+v", event)
	}

	id := uuid.New()
	req := &reporterv1beta1.ReportRequest{
		Events: []*reporterv1beta1.Event{{
			Id:   id[:],
			Time: ptypes.TimestampNow(),
			Event: &reporterv1beta1.AnyEvent{
				TypeUrl: proto.MessageName(event),
				Binary:  eventB,
			},
		}},
	}
	s.l.Debugf("Request: %+v", req)
	if err = req.Validate(); err != nil {
		// log and ignore
		s.l.Debugf("Failed to validate request: %s.", err)
	}

	return req, nil
}

func (s *Service) sendV2Request(ctx context.Context, req *reporterv1beta1.ReportRequest) error {
	if s.v2Host == "" {
		return errors.New("v2 telemetry disabled via the empty host")
	}

	host, _, err := net.SplitHostPort(s.v2Host)
	if err != nil {
		host = s.v2Host
	}
	tlsConfig := tlsconfig.Get()
	tlsConfig.ServerName = host

	opts := []grpc.DialOption{
		// replacement is marked as experimental
		grpc.WithBackoffMaxDelay(time.Second), //nolint:staticcheck

		grpc.WithBlock(),
		grpc.WithUserAgent("pmm-managed/" + s.pmmVersion),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}
	cc, err := grpc.DialContext(ctx, s.v2Host, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}
	defer cc.Close() //nolint:errcheck

	if _, err = reporterv1beta1.NewReporterAPIClient(cc).Report(ctx, req); err != nil {
		return errors.Wrap(err, "failed to report")
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
