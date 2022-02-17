// Package telemetry_v2 provides telemetry v2 functionality.
package telemetry_v2

import (
	//nolint:staticcheck
	"github.com/percona/pmm-managed/utils/saasreq"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"

	"bytes"
	"context"
	"encoding/hex"
	"github.com/google/uuid"
	events "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm/api/serverpb"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	//nolint:staticcheck
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

const (
	distributionInfoFilePath = "/srv/pmm-distribution"
	osInfoFilePath           = "/proc/version"
)

// Service reports telemetry.
type Service struct {
	db         *reform.DB
	l          *logrus.Entry
	start      time.Time
	config     ServiceConfig
	dsRegistry TelemetryDataSourceLocator
	pmmVersion string
}

func (s *Service) LocateTelemetryDataSource(name string) (TelemetryDataSource, error) {
	return s.dsRegistry.LocateTelemetryDataSource(name)
}

// check interfaces
var (
	_ TelemetryDataSourceLocator = (*Service)(nil)
)

// NewService creates a new service.
func NewService(db *reform.DB, pmmVersion string, config ServiceConfig) (*Service, error) {
	if config.SaasHostname == "" {
		return nil, errors.New("empty host")
	}

	l := logrus.WithField("component", "telemetry_v2")

	registry, err := NewDataSourceRegistry(config, l)
	if err != nil {
		return nil, err
	}
	s := &Service{
		db:         db,
		l:          l,
		pmmVersion: pmmVersion,
		start:      time.Now(),
		config:     config,
		dsRegistry: registry,
	}

	return s, nil
}

// Run start sending telemetry to SaaS.
func (s *Service) Run(ctx context.Context) {
	if !s.config.Enabled {
		s.l.Warn("service is disabled, skip Run")
		return
	}

	ticker := time.NewTicker(s.config.Reporting.Interval)
	defer ticker.Stop()

	doSend := func() {
		report, err := s.prepareReport(ctx)
		if err != nil {
			s.l.Debugf("Failed to prepare report: %s.", err)
			return
		}

		err = s.send(ctx, report)
		if err != nil {
			s.l.Debugf("Telemetry info not sent, due to error: %s.", err)
			return
		}
		s.l.Debug("Telemetry info sent.")
	}

	if s.config.Reporting.SendOnStart {
		s.l.Debug("Telemetry on start is enabled, sending...")
		doSend()
	}

	for {
		select {
		case <-ticker.C:
			doSend()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) prepareReport(ctx context.Context) (*reporter.ReportRequest, error) {
	current, err := s.makeMetric()
	if err != nil {
		return nil, err
	}

	for _, telemetry := range s.config.telemetry {
		ds, err := s.LocateTelemetryDataSource(telemetry.Source)
		if err != nil {
			s.l.Debugf("failed to lookup telemetry datasource for [%s]:[%s]", telemetry.Source, telemetry.Id)
			continue
		}

		metrics, err := ds.FetchMetrics(ctx, telemetry)
		if err != nil {
			s.l.Debugf("failed to extract metric from datasource for [%s]:[%s]: %v", telemetry.Source, telemetry.Id, err)
			continue
		}

		current.Metrics = append(current.Metrics, metrics...)
	}

	return &reporter.ReportRequest{
		Metrics: []*reporter.ServerMetric{current},
	}, nil
}

func (s *Service) makeMetric() (*reporter.ServerMetric, error) {
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
		return nil, err
	}

	serverID, err := hex.DecodeString(settings.Telemetry.UUID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode UUID %q", settings.Telemetry.UUID)
	}

	_, distMethod, _ := getDistributionMethodAndOS(s.l)

	eventId := uuid.New()
	return &reporter.ServerMetric{
		Id:                   eventId[:],
		Time:                 timestamppb.New(time.Now()),
		PmmServerTelemetryId: serverID,
		PmmServerVersion:     s.pmmVersion,
		UpDuration:           durationpb.New(time.Since(s.start)),
		DistributionMethod:   distMethod,
	}, nil
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

func (s *Service) send(ctx context.Context, report *reporter.ReportRequest) error {
	var err error
	var attempt int
	for {
		err = s.sendRequest(ctx, report)
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

func (s *Service) sendRequest(ctx context.Context, req *reporter.ReportRequest) error {
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
