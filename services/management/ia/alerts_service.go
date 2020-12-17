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

// Package ia contains Integrated Alerting APIs implementations.
package ia

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona-platform/saas/pkg/common"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

type AlertsService struct {
	db               *reform.DB
	l                *logrus.Entry
	alertManager     alertManager
	templatesService *TemplatesService
}

func NewAlertsService(db *reform.DB, alertManager alertManager, templatesService *TemplatesService) *AlertsService {
	return &AlertsService{
		db:               db,
		alertManager:     alertManager,
		templatesService: templatesService,
		l:                logrus.WithField("component", "management/ia/alerts"),
	}
}

type Alert struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	GroupID     string            `json:"group_id"`
	Expression  string            `json:"expression"`
	ActiveAt    time.Time         `json:"activeAt"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
	State       prom.AlertState   `json:"state"`
	Value       string            `json:"value"`
}

func (s *AlertsService) ListAlerts(ctx context.Context, req *iav1beta1.ListAlertsRequest) (*iav1beta1.ListAlertsResponse, error) {
	alerts, err := s.alertManager.GetAlerts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get alerts form alertmanager")
	}

	var res []*iav1beta1.Alert
	for _, alert := range alerts {

		updatedAt, err := ptypes.TimestampProto(time.Time(*alert.UpdatedAt))
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert timestamp")
		}

		createdAt, err := ptypes.TimestampProto(time.Time(*alert.StartsAt))
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert timestamp")
		}

		st := iav1beta1.Status_STATUS_INVALID
		if *alert.Status.State == "active" {
			st = iav1beta1.Status_TRIGGERING
		}

		if len(alert.Status.SilencedBy) != 0 {
			st = iav1beta1.Status_SILENCED
		}

		ruleID, ok := alert.Labels["alertname"]
		if !ok {
			return nil, errors.New("missing 'alertname' label")
		}
		var rule *models.Rule
		var channels []*models.Channel
		e := s.db.InTransaction(func(tx *reform.TX) error {
			var err error
			rule, err = models.FindRuleByID(tx.Querier, ruleID)
			if err != nil {
				return err
			}

			channels, err = models.FindChannelsByIDs(tx.Querier, rule.ChannelIDs)
			return err
		})
		if e != nil {
			return nil, e
		}

		template, ok := s.templatesService.GetTemplates(ctx)[rule.TemplateName]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "Failed to find template with name: %s", rule.TemplateName)
		}

		r, err := convertRule(s.l, rule, template, channels)
		if err != nil {

		}

		res = append(res, &iav1beta1.Alert{
			AlertId:   getAlertID(alert),
			Summary:   alert.Annotations["summary"],
			Severity:  managementpb.Severity(common.ParseSeverity(alert.Labels["severity"])),
			Status:    st,
			Labels:    alert.Labels,
			Rule:      r,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return &iav1beta1.ListAlertsResponse{Alerts: res}, nil
}

func getAlertID(alert *ammodels.GettableAlert) string {
	return *alert.Fingerprint
}

func (s *AlertsService) ToggleAlert(ctx context.Context, req *iav1beta1.ToggleAlertRequest) (*iav1beta1.ToggleAlertResponse, error) {
	switch req.Silenced {
	case iav1beta1.BooleanFlag_TRUE:
		err := s.alertManager.Silence(ctx, req.AlertId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to silence alert with id: %s", req.AlertId)
		}
	case iav1beta1.BooleanFlag_FALSE:
		err := s.alertManager.Unsilence(ctx, req.AlertId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unsilence alert with id: %s", req.AlertId)
		}
	}

	return &iav1beta1.ToggleAlertResponse{}, nil
}

func (s *AlertsService) DeleteAlert(ctx context.Context, req *iav1beta1.DeleteAlertRequest) (*iav1beta1.DeleteAlertResponse, error) {
	panic("implement me")
}

// Check interfaces.
var (
	_ iav1beta1.AlertsServer = (*AlertsService)(nil)
)
