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

// Package telemetry_v2 provides telemetry v2 functionality.
package telemetry_v2

import (
	"context"
	"time"

	//nolint:staticcheck
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

// Service reports telemetry.
type Service struct {
	db     *reform.DB
	l      *logrus.Entry
	config ServiceConfig
}

// NewService creates a new service.
func NewService(db *reform.DB, config ServiceConfig) (*Service, error) {
	l := logrus.WithField("component", "telemetry_v2")

	s := &Service{
		db:     db,
		l:      l,
		config: config,
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
		err := s.send(ctx)
		if err == nil {
			s.l.Debug("Telemetry info sent.")
		} else {
			s.l.Debugf("Telemetry info not sent, due to error: %s.", err)
		}
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

func (s *Service) send(ctx context.Context) error {
	return nil

	//TODO:
	//ctx, cancel := context.WithTimeout(ctx, s.reportInterval)
	//defer cancel()
	//
	//var settings *models.Settings
	//err := s.db.InTransaction(func(tx *reform.TX) error {
	//	var e error
	//	if settings, e = models.GetSettings(tx); e != nil {
	//		return e
	//	}
	//
	//	if settings.Telemetry.Disabled {
	//		return errors.New("disabled via settings")
	//	}
	//	if settings.Telemetry.UUID == "" {
	//		settings.Telemetry.UUID, e = generateUUID()
	//		if e != nil {
	//			return e
	//		}
	//		return models.SaveSettings(tx, settings)
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return err
	//}
	//
	//s.l.Debugf("Using %s as server UUID.", settings.Telemetry.UUID)
	//
	//var wg errgroup.Group
	//
	//wg.Go(func() error {
	//	payload := s.makeV1Payload(settings.Telemetry.UUID)
	//	return s.sendV1Request(ctx, payload)
	//})
	//
	//wg.Go(func() error {
	//	req, err := s.makeV2Payload(settings.Telemetry.UUID, settings)
	//	if err != nil {
	//		return err
	//	}
	//
	//	return s.sendV2RequestWithRetries(ctx, req)
	//})
	//
	//return wg.Wait()
}

//func (s *Service) makeV2Payload(serverUUID string, settings *models.Settings) (*reporter.ReportRequest, error) {
//	serverID, err := hex.DecodeString(serverUUID)
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to decode UUID %q", serverUUID)
//	}
//
//	event := &events.ServerUptimeEvent{
//		Id:                 serverID,
//		Version:            s.pmmVersion,
//		UpDuration:         durationpb.New(time.Since(s.initializedAt)),
//		DistributionMethod: s.tDistributionMethod,
//		SttEnabled:         wrapperspb.Bool(settings.SaaS.STTEnabled),
//		IaEnabled:          wrapperspb.Bool(settings.IntegratedAlerting.Enabled),
//	}
//
//	if err = event.Validate(); err != nil {
//		// log and ignore
//		s.l.Debugf("Failed to validate event: %s.", err)
//	}
//	eventB, err := proto.Marshal(event)
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to marshal event %+v", event)
//	}
//
//	id := uuid.New()
//	now := time.Now()
//	req := &reporter.ReportRequest{
//		Events: []*reporter.Event{{
//			Id: id[:],
//			Time: &timestamp.Timestamp{
//				Seconds: now.Unix(),
//				Nanos:   int32(now.Nanosecond()),
//			},
//			Event: &reporter.AnyEvent{
//				TypeUrl: proto.MessageName(event), //nolint:staticcheck
//				Binary:  eventB,
//			},
//		}},
//	}
//	s.l.Debugf("Request: %+v", req)
//	if err = req.Validate(); err != nil {
//		// log and ignore
//		s.l.Debugf("Failed to validate request: %s.", err)
//	}
//
//	return req, nil
//}
//
//func (s *Service) sendV2RequestWithRetries(ctx context.Context, req *reporter.ReportRequest) error {
//	if s.v2Host == "" {
//		return errors.New("v2 telemetry disabled via the empty host")
//	}
//
//	var err error
//	var attempt int
//	for {
//		err = s.sendV2Request(ctx, req)
//		attempt++
//		s.l.Debugf("sendV2Request (attempt %d/%d) result: %v", attempt, s.reportRetryCount, err)
//		if err == nil {
//			return nil
//		}
//
//		if attempt >= s.reportRetryCount {
//			s.l.Debug("Failed to send v2 event, will not retry (too much attempts).")
//			return err
//		}
//
//		retryCtx, retryCancel := context.WithTimeout(ctx, s.reportRetryBackoff)
//		<-retryCtx.Done()
//		retryCancel()
//
//		if err = ctx.Err(); err != nil {
//			s.l.Debugf("Will not retry sending v2 event: %s.", err)
//			return err
//		}
//	}
//}
//
//func (s *Service) sendV2Request(ctx context.Context, req *reporter.ReportRequest) error {
//	s.l.Debugf("Using %s as telemetry host.", s.v2Host)
//
//	var accessToken string
//	if ssoDetails, err := models.GetPerconaSSODetails(ctx, s.db.Querier); err == nil {
//		accessToken = ssoDetails.AccessToken.AccessToken
//	}
//
//	reqByte, err := protojson.Marshal(req)
//	if err != nil {
//		return err
//	}
//
//	endpoint := fmt.Sprintf("https://%s/v1/telemetry/Report", s.v2Host)
//	_, err = saasreq.MakeRequest(ctx, http.MethodPost, endpoint, accessToken, bytes.NewReader(reqByte))
//	if err != nil {
//		return errors.Wrap(err, "failed to dial")
//	}
//
//	return nil
//}

////TODO this is stub for collecting Server Metrics
//func (s *Service) makeV2ServiceMetric(serverUUID string) (*reporter.ReportRequest, error) {
//	serverID, err := hex.DecodeString(serverUUID)
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to decode UUID %q", serverUUID)
//	}
//
//	var metrics []*reporter.ServerMetric
//	var metrics2 []*reporter.ServerMetric_Metric
//	id := uuid.New()
//	metrics = append(metrics, &reporter.ServerMetric{
//		Id:                   id[:],
//		Time:                 timestamppb.Now(),
//		PmmServerTelemetryId: serverID,
//		PmmServerVersion:     "2",
//		UpDuration:           nil,
//		DistributionMethod:   0,
//		Metrics: append(metrics2, &reporter.ServerMetric_Metric{
//			Key:   "key1",
//			Value: "val1",
//		}),
//	})
//
//	req := &reporter.ReportRequest{
//		Metrics: metrics,
//	}
//	s.l.Debugf("Request: %+v", req)
//
//	return req, nil
//}
