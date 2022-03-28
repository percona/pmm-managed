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

package services

import (
	"fmt"

	"github.com/percona-platform/saas/pkg/check"
	"github.com/percona/pmm/api/alertmanager/amclient/alert"

	"github.com/percona/pmm-managed/models"
)

const (
	// CheckFilter represents AlertManager filter for Checks/Advisor results.
	CheckFilter = "stt_check=1"
	// IAFilter represents AlertManager filter for Integrated Alerts.
	IAFilter = "ia=1"
)

// Target contains required info about STT check target.
type Target struct {
	AgentID       string
	ServiceID     string
	ServiceName   string
	Labels        map[string]string
	DSN           string
	Files         map[string]string
	TDP           *models.DelimiterPair
	TLSSkipVerify bool
}

// CheckResult contains the output from the check file and other information.
type CheckResult struct {
	CheckName string
	Silenced  bool
	AlertID   string
	Interval  check.Interval
	Target    Target
	Result    check.Result
}

// CheckResultSummary contains the summary of failed checks for a service.
type CheckResultSummary struct {
	ServiceName   string
	ServiceID     string
	CriticalCount uint32
	WarningCount  uint32
	NoticeCount   uint32
}

// FilterParams provides fields needed to filter alerts from AlertManager.
type FilterParams struct {
	// IsIA specifies if only Integrated Alerts should be matched.
	IsIA bool
	// IsCheck specifies if only Checks/Advisors alerts should be matched.
	IsCheck bool
	// AlertID is the ID of alert to be matched (if any).
	AlertID string
	// ServiceID is the ID of service to be matched (if any).
	ServiceID string
}

// ToAlertManagerParams returns an AlertManager-style filter for FilterParams.
func (fp FilterParams) ToAlertManagerParams() *alert.GetAlertsParams {
	alertParams := alert.NewGetAlertsParams()
	if fp.IsCheck {
		alertParams.Filter = append(alertParams.Filter, CheckFilter)
	}
	if fp.IsIA {
		alertParams.Filter = append(alertParams.Filter, IAFilter)
	}
	if fp.ServiceID != "" {
		alertParams.Filter = append(alertParams.Filter, fmt.Sprintf("service_id=\"%s\"", fp.ServiceID))
	}
	if fp.AlertID != "" {
		alertParams.Filter = append(alertParams.Filter, fmt.Sprintf("alert_id=\"%s\"", fp.AlertID))
	}
	return alertParams
}
