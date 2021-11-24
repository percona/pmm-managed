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
	"github.com/percona/pmm-managed/models"
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify/email"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
)

type Emailer struct {
	Settings *models.EmailAlertingSettings
}

func (e *Emailer) Send(ctx context.Context, emailTo string) error {
	host, port, err := net.SplitHostPort(e.Settings.Smarthost)
	if err != nil {
		return models.NewInvalidArgumentError("invalid smarthost: %q", err.Error())
	}

	if port == "" {
		return models.NewInvalidArgumentError("address %q: port cannot be empty", port)
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
			"Subject": `Test alert.`,
		},
		HTML:       emailTemplate,
		RequireTLS: &e.Settings.RequireTLS,
	}

	tmpl, err := template.FromGlobs()
	if err != nil {
		return err
	}
	tmpl.ExternalURL, err = url.Parse("https://example.com")
	if err != nil {
		return err
	}

	alertmanagerEmail := email.New(emailConfig, tmpl, kitlog.NewNopLogger())
	if _, err := alertmanagerEmail.Notify(ctx, &types.Alert{
		Alert: model.Alert{
			Labels: model.LabelSet{
				model.AlertNameLabel: model.LabelValue(fmt.Sprintf("Test alert %s", time.Now().String())),
				"severity":           "notice",
			},
			Annotations: model.LabelSet{
				"summary":     "This is a test alert.",
				"description": "Long description.",
				"rule":        "example-violated-rule",
			},
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(time.Minute),
		},
		Timeout: true,
	}); err != nil {
		return models.NewInvalidArgumentError("failed to send email: %s", err.Error())
	}

	return nil
}
