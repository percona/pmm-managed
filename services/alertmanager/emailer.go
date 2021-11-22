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

package alertmanager

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/pkg/errors"
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"

	"github.com/percona/pmm-managed/models"
	"github.com/prometheus/alertmanager/notify/email"
)

type Emailer struct {
	Settings *models.EmailAlertingSettings
	Logger   *logrus.Entry
}

func (e *Emailer) Send(ctx context.Context, emailTo string) error {
	host, port, err := net.SplitHostPort(e.Settings.Smarthost)
	if err != nil {
		return err
	}

	if port == "" {
		return errors.Errorf("address %q: port cannot be empty", port)
	}

	emailConfig := &config.EmailConfig{
		NotifierConfig: config.NotifierConfig{},
		To:             emailTo,
		From:           e.Settings.From,
		Hello:          e.Settings.Hello,
		Smarthost: config.HostPort{
			Host: host,
			Port: port,
		},
		AuthUsername: e.Settings.Username,
		AuthPassword: config.Secret(e.Settings.Password),
		AuthSecret:   config.Secret(e.Settings.Secret),
		AuthIdentity: e.Settings.Identity,
		Headers: map[string]string{
			"Subject": `[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}]`,
		},
		HTML: emailTemplate,
		// RequireTLS: e.settings.RequireTLS TODO: implement once https://jira.percona.com/browse/PMM-9068 get merged
	}

	tmpl, err := template.FromGlobs()
	if err != nil {
		return err
	}
	tmpl.ExternalURL, _ = url.Parse("https://example.com")

	alertmanagerEmail := email.New(emailConfig, tmpl, kitlog.NewNopLogger())
	if _, err := alertmanagerEmail.Notify(ctx, &types.Alert{
		Alert: model.Alert{
			Labels: model.LabelSet{
				model.AlertNameLabel: model.LabelValue(fmt.Sprintf("This is a test alert %s", time.Now().String())),
			},
			Annotations: model.LabelSet{},
			StartsAt:    time.Now(),
			EndsAt:      time.Now().Add(time.Minute),
		},
		UpdatedAt: time.Time{},
		Timeout:   false,
	}); err != nil {
		return err
	}

	return nil
}
