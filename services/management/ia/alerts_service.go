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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona-platform/saas/pkg/common"
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
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

var generatorURLRegexp = regexp.MustCompile("http://localhost:9090/prometheus/api/v1/\\d*/(\\d*)/status")

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
	vmAlerts, err := getAlertsFromVM(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get alerts from victoriametrics")
	}

	alerts, err := s.alertManager.GetAlerts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get alerts form alertmanager")
	}

	var res []*iav1beta1.Alert
	for _, alert := range alerts {
		match := generatorURLRegexp.FindStringSubmatch(alert.GeneratorURL.String())
		alertID := match[1]

		vmAlert, ok := vmAlerts[alertID]
		if !ok {
			return nil, errors.Errorf("failed to find alert with id: %s", alertID)
		}

		updatedAt, err := ptypes.TimestampProto(time.Time(*alert.UpdatedAt))
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert timestamp")
		}

		createdAt, err := ptypes.TimestampProto(vmAlert.ActiveAt) // TODO not sure about it, alternative is alert.StartsAt
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert timestamp")
		}

		var status iav1beta1.Status
		switch vmAlert.State {
		case prom.AlertStateFiring:
			status = iav1beta1.Status_TRIGGERING
		case prom.AlertStatePending:
			status = iav1beta1.Status_PENDING
		}

		if len(alert.Status.SilencedBy) != 0 || len(alert.Status.InhibitedBy) != 0 { // TODO do we interpretate inhibition as silencing?
			status = iav1beta1.Status_SILENCED
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
			// TODO How to handle that case?
			return nil, errors.Errorf("failed to find template with name: %s", rule.TemplateName)
		}

		r, err := convertRule(s.l, rule, template, channels)
		if err != nil {

		}

		res = append(res, &iav1beta1.Alert{
			AlertId:   alertID,
			Summary:   alert.Annotations["summary"],
			Severity:  managementpb.Severity(common.ParseSeverity(alert.Labels["severity"])),
			Status:    status,
			Labels:    alert.Labels,
			Rule:      r,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return &iav1beta1.ListAlertsResponse{Alerts: res}, nil
}

func getAlertsFromVM(ctx context.Context) (map[string]*Alert, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8880/api/v1/alerts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	alerts := struct {
		Data struct {
			Alerts []*Alert `json:"alerts"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &alerts)
	if err != nil {
		return nil, err
	}

	res := make(map[string]*Alert, len(alerts.Data.Alerts))
	for _, alert := range alerts.Data.Alerts {
		res[alert.ID] = alert
	}

	return res, nil
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
