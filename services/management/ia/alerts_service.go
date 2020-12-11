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
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

type AlertsService struct {
	db               *reform.DB
	l                *logrus.Entry
	alertManager     alertManager
	templatesService *TemplatesService
}

func NewAlertsService(db *reform.DB) *AlertsService {
	return &AlertsService{
		db: db,
		l:  logrus.WithField("component", "management/ia/alerts"),
	}
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

		status := iav1beta1.Status_TRIGGERING
		if len(alert.Status.SilencedBy) != 0 {
			status = iav1beta1.Status_SILENCED
		}
		if len(alert.Status.InhibitedBy) != 0 {
			status = iav1beta1.Status_PENDING
		}

		ruleID := alert.Labels["alertname"] // TODO is that field confirmed?
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
			// TODO How to handle that case?
		}

		r, err := convertRule(s.l, rule, template, channels)
		if err != nil {

		}

		res = append(res, &iav1beta1.Alert{
			AlertId:  alert.Labels["alert_id"], // TODO missing in vmalert generated alerts.
			Summary:  alert.Annotations["summary"],
			Severity: managementpb.Severity(common.ParseSeverity(alert.Labels["severity"])),
			Status:   status,
			Labels:   alert.Labels,
			Rule:     r,
			// CreatedAt: nil, // TODO ???
			UpdatedAt: updatedAt,
		})
	}

	return &iav1beta1.ListAlertsResponse{Alerts: res}, nil
}

func (s *AlertsService) ToggleAlert(ctx context.Context, req *iav1beta1.ToggleAlertRequest) (*iav1beta1.ToggleAlertResponse, error) {
	panic("implement me")
}

func (s *AlertsService) DeleteAlert(ctx context.Context, req *iav1beta1.DeleteAlertRequest) (*iav1beta1.DeleteAlertResponse, error) {
	panic("implement me")
}

// Check interfaces.
var (
	_ iav1beta1.AlertsServer = (*AlertsService)(nil)
)
