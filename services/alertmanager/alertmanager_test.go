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
	"strings"
	"testing"
	"time"

	"github.com/percona-platform/saas/pkg/alert"
	"github.com/percona-platform/saas/pkg/common"
	"github.com/percona/promconfig"
	"github.com/percona/promconfig/alertmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/yaml.v3"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestIsReady(t *testing.T) {
	New(nil).GenerateBaseConfigs() // this method should not use database

	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	svc := New(db)

	assert.NoError(t, svc.updateConfiguration(ctx))
	assert.NoError(t, svc.IsReady(ctx))
}

// marshalAndValidate populates, marshals and validates config.
func marshalAndValidate(t *testing.T, svc *Service, base *alertmanager.Config) string {
	b, err := svc.marshalConfig(base)
	require.NoError(t, err)

	t.Logf("config:\n%s", b)

	err = svc.validateConfig(context.Background(), b)
	require.NoError(t, err)
	return string(b)
}

func TestPopulateConfig(t *testing.T) {
	New(nil).GenerateBaseConfigs() // this method should not use database

	t.Run("without receivers and routes", func(t *testing.T) {
		tests.SetTestIDReader(t)
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		svc := New(db)

		cfg := svc.loadBaseConfig()
		cfg.Global = &alertmanager.GlobalConfig{
			SlackAPIURL: "https://hooks.slack.com/services/abc/123/xyz",
		}

		actual := marshalAndValidate(t, svc, cfg)
		expected := strings.TrimSpace(`
# Managed by pmm-managed. DO NOT EDIT.
---
global:
    resolve_timeout: 0s
    smtp_require_tls: false
    slack_api_url: https://hooks.slack.com/services/abc/123/xyz
route:
    receiver: empty
    continue: false
receivers:
    - name: empty
    - name: disabled
templates: []
		`) + "\n"
		assert.Equal(t, expected, actual, "actual:\n%s", actual)
	})

	t.Run("with receivers and routes", func(t *testing.T) {
		tests.SetTestIDReader(t)
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		svc := New(db)

		channel1, err := models.CreateChannel(db.Querier, &models.CreateChannelParams{
			Summary: "channel1",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test", "test2@test.test"},
			},
			Disabled: false,
		})
		require.NoError(t, err)

		channel2, err := models.CreateChannel(db.Querier, &models.CreateChannelParams{
			Summary: "channel2",
			PagerDutyConfig: &models.PagerDutyConfig{
				RoutingKey: "ms-pagerduty-dev",
			},
			Disabled: false,
		})
		require.NoError(t, err)

		channel3, err := models.CreateChannel(db.Querier, &models.CreateChannelParams{
			Summary: "channel3",
			PagerDutyConfig: &models.PagerDutyConfig{
				RoutingKey: "ms-pagerduty-dev",
			},
			Disabled: true,
		})
		require.NoError(t, err)

		_, err = models.CreateTemplate(db.Querier, &models.CreateTemplateParams{
			Template: &alert.Template{
				Name:    "test_template",
				Version: 1,
				Summary: "summary",
				Tiers:   []common.Tier{common.Anonymous},
				Expr:    "expr",
				Params: []alert.Parameter{{
					Name:    "param",
					Summary: "param summary",
					Unit:    "%",
					Type:    alert.Float,
					Range:   []interface{}{float64(10), float64(100)},
					Value:   float64(50),
				}},
				For:         promconfig.Duration(3 * time.Second),
				Severity:    common.Warning,
				Labels:      map[string]string{"foo": "bar"},
				Annotations: nil,
			},
			Source: "USER_FILE",
		})
		require.NoError(t, err)

		rule1, err := models.CreateRule(db.Querier, &models.CreateRuleParams{
			TemplateName: "test_template",
			Disabled:     true,
			RuleParams: []models.RuleParam{{
				Name:       "test",
				Type:       models.Float,
				FloatValue: 3.14,
			}},
			For:      5 * time.Second,
			Severity: models.Severity(common.Warning),
			CustomLabels: map[string]string{
				"foo": "bar",
			},
			Filters: []models.Filter{{
				Type: models.Equal,
				Key:  "service_name",
				Val:  "mysql1",
			}},
			ChannelIDs: []string{channel1.ID, channel2.ID},
		})
		require.NoError(t, err)

		// create another rule with same channelIDs to check for redundant receivers.
		rule2, err := models.CreateRule(db.Querier, &models.CreateRuleParams{
			TemplateName: "test_template",
			Disabled:     true,
			RuleParams: []models.RuleParam{{
				Name:       "test",
				Type:       models.Float,
				FloatValue: 3.14,
			}},
			For:      5 * time.Second,
			Severity: models.Severity(common.Warning),
			CustomLabels: map[string]string{
				"foo": "baz",
			},
			Filters: []models.Filter{{
				Type: models.Equal,
				Key:  "service_name",
				Val:  "mysql2",
			}},
			ChannelIDs: []string{channel1.ID, channel2.ID, channel3.ID},
		})
		require.NoError(t, err)

		// create another rule without channelID and check if it is absent in the config.
		_, err = models.CreateRule(db.Querier, &models.CreateRuleParams{
			TemplateName: "test_template",
			Disabled:     true,
			RuleParams: []models.RuleParam{{
				Name:       "test",
				Type:       models.Float,
				FloatValue: 3.14,
			}},
			For:      5 * time.Second,
			Severity: models.Severity(common.Warning),
			CustomLabels: map[string]string{
				"foo": "baz",
			},
		})
		require.NoError(t, err)

		// CreateRule with disabled channel
		rule4, err := models.CreateRule(db.Querier, &models.CreateRuleParams{
			TemplateName: "test_template",
			Disabled:     true,
			RuleParams: []models.RuleParam{{
				Name:       "test",
				Type:       models.Float,
				FloatValue: 3.14,
			}},
			Filters: []models.Filter{{
				Type: models.Equal,
				Key:  "service_name",
				Val:  "mysql3",
			}},
			For:      5 * time.Second,
			Severity: models.Severity(common.Warning),
			CustomLabels: map[string]string{
				"foo": "baz",
			},
			ChannelIDs: []string{channel3.ID},
		})
		require.NoError(t, err)

		_, err = models.UpdateSettings(db.Querier, &models.ChangeSettingsParams{
			EmailAlertingSettings: &models.EmailAlertingSettings{
				From:      "from@test.com",
				Smarthost: "1.2.3.4:80",
				Hello:     "host",
				Username:  "user",
				Password:  "password",
				Identity:  "id",
				Secret:    "secret",
			},
			SlackAlertingSettings: &models.SlackAlertingSettings{
				URL: "https://hooks.slack.com/services/abc/456/xyz",
			},
		})
		require.NoError(t, err)

		actual := marshalAndValidate(t, svc, svc.loadBaseConfig())
		expected := strings.TrimSpace(fmt.Sprintf(`
# Managed by pmm-managed. DO NOT EDIT.
---
global:
    resolve_timeout: 0s
    smtp_from: from@test.com
    smtp_hello: host
    smtp_smarthost: 1.2.3.4:80
    smtp_auth_username: user
    smtp_auth_password: password
    smtp_auth_secret: secret
    smtp_auth_identity: id
    smtp_require_tls: false
    slack_api_url: https://hooks.slack.com/services/abc/456/xyz
route:
    receiver: empty
    continue: false
    routes:
        - receiver: %[1]s + %[2]s
          match:
            rule_id: %[3]s
            service_name: mysql1
          continue: false
        - receiver: %[1]s + %[2]s
          match:
            rule_id: %[4]s
            service_name: mysql2
          continue: false
        - receiver: disabled
          match:
            rule_id: %[5]s
            service_name: mysql3
          continue: false
receivers:
    - name: empty
    - name: disabled
    - name: %[1]s + %[2]s
      email_configs:
        - send_resolved: false
          to: test@test.test
          headers:
            Subject: '[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonLabels.alertname }}'
          html: |
            <!--
            Style and HTML derived from https://github.com/mailgun/transactional-email-templates


            The MIT License (MIT)

            Copyright (c) 2014 Mailgun

            Permission is hereby granted, free of charge, to any person obtaining a copy
            of this software and associated documentation files (the "Software"), to deal
            in the Software without restriction, including without limitation the rights
            to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
            copies of the Software, and to permit persons to whom the Software is
            furnished to do so, subject to the following conditions:

            The above copyright notice and this permission notice shall be included in all
            copies or substantial portions of the Software.

            THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
            IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
            FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
            AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
            LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
            OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
            SOFTWARE.
            -->
            <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
            <html xmlns="http://www.w3.org/1999/xhtml">
            <head>
                <meta name="viewport" content="width=device-width" />
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <title>{{ template "__subject" . }}</title>
                <style>
                    /* -------------------------------------
                        GLOBAL
                        A very basic CSS reset
                    ------------------------------------- */
                    * {
                        margin: 0;
                        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
                        box-sizing: border-box;
                        font-size: 14px;
                    }

                    img {
                        max-width: 100%%;
                    }

                    body {
                        -webkit-font-smoothing: antialiased;
                        -webkit-text-size-adjust: none;
                        width: 100%% !important;
                        height: 100%%;
                        line-height: 1.6em;
                        /* 1.6em * 14px = 22.4px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
                        /*line-height: 22px;*/
                    }

                    /* Let's make sure all tables have defaults */
                    table td {
                        vertical-align: top;
                    }

                    /* -------------------------------------
                        BODY & CONTAINER
                    ------------------------------------- */
                    body {
                        background-color: #f6f6f6;
                    }

                    .body-wrap {
                        background-color: #f6f6f6;
                        width: 100%%;
                    }

                    .container {
                        display: block !important;
                        max-width: 600px !important;
                        margin: 0 auto !important;
                        /* makes it centered */
                        clear: both !important;
                    }

                    .content {
                        max-width: 600px;
                        margin: 0 auto;
                        display: block;
                        padding: 20px;
                    }

                    /* -------------------------------------
                        HEADER, FOOTER, MAIN
                    ------------------------------------- */
                    .main {
                        background-color: #fff;
                        border: 1px solid #e9e9e9;
                        border-radius: 3px;
                    }

                    .content-wrap {
                        padding: 30px;
                    }

                    .content-block {
                        padding: 0 0 20px;
                    }

                    .footer {
                        width: 100%%;
                        clear: both;
                        color: #999;
                        padding: 20px;
                    }
                    .footer p, .footer a, .footer td {
                        color: #999;
                        font-size: 12px;
                    }

                    /* -------------------------------------
                        LINKS & BUTTONS
                    ------------------------------------- */
                    a {
                        color: #348eda;
                        text-decoration: underline;
                    }

                    .btn-primary {
                        text-decoration: none;
                        color: #FFF;
                        background-color: #3274D9;
                        border: solid #3274D9;
                        border-width: 10px 20px;
                        line-height: 2em;
                        /* 2em * 14px = 28px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
                        /*line-height: 28px;*/
                        font-weight: bold;
                        text-align: center;
                        cursor: pointer;
                        display: inline-block;
                        border-radius: 5px;
                        text-transform: capitalize;
                    }

                    /* -------------------------------------
                        OTHER STYLES THAT MIGHT BE USEFUL
                    ------------------------------------- */
                    .aligncenter {
                        text-align: center;
                    }

                    /* -------------------------------------
                        ALERTS
                        Change the class depending on warning email, good email or bad email
                    ------------------------------------- */
                    .alert {
                        font-size: 16px;
                        color: #fff;
                        font-weight: 500;
                        padding: 20px;
                        text-align: center;
                        border-radius: 3px 3px 0 0;
                    }
                    .alert a {
                        color: #fff;
                        text-decoration: none;
                        font-weight: 500;
                        font-size: 16px;
                    }
                    .alert.alert-notice {
                        background-color: #3274D9;
                    }
                    .alert.alert-warning {
                        background-color: #ECBB13;
                    }
                    .alert.alert-high {
                        background-color: #EB7B18;
                    }
                    .alert.alert-critical {
                        background-color: #E02F44;
                    }

                    /* -------------------------------------
                        RESPONSIVE AND MOBILE FRIENDLY STYLES
                    ------------------------------------- */
                    @media only screen and (max-width: 640px) {
                        body {
                            padding: 0 !important;
                        }

                        .container {
                            padding: 0 !important;
                            width: 100%% !important;
                        }

                        .content {
                            padding: 0 !important;
                        }

                        .content-wrap {
                            padding: 10px !important;
                        }
                    }
                </style>
            </head>

            <body itemscope itemtype="http://schema.org/EmailMessage">

            <table class="body-wrap">
                <tr>
                    <td></td>
                    <td class="container" width="600">
                        <div class="content">
                            <table class="main" width="100%%" cellpadding="0" cellspacing="0">
                                <tr>
                                    {{ if eq .CommonLabels.severity "notice" }}
                                    <td class="alert alert-notice">
                                        {{ else if eq .CommonLabels.severity "warning" }}
                                    <td class="alert alert-warning">
                                        {{ else if eq .CommonLabels.severity "error" }}
                                    <td class="alert alert-high">
                                        {{ else if eq .CommonLabels.severity "critical" }}
                                    <td class="alert alert-critical">
                                        {{ else }}
                                    <td class="alert alert-critical">
                                        {{ end }}
                                        You have {{ .Alerts | len }} {{ .CommonLabels.severity }} alert{{ if gt (len .Alerts) 1 }}s{{ end }} firing
                                    </td>
                                </tr>
                                <tr>
                                    <td class="content-wrap">
                                        <table width="100%%" cellpadding="0" cellspacing="0">
                                            <tr>
                                                <td align="center" class="content-block">
                                                    <a href='{{ template "__alertmanagerURL" . }}' class="btn-primary">View in {{ template "__alertmanager" . }}</a>
                                                </td>
                                            </tr>
                                            {{ range .Alerts.Firing }}
                                                <tr>
                                                    <td class="content-block">
                                                        <strong>------------------</strong><br />
                                                        <strong>Alert: </strong>
                                                        {{ if .Labels.severity }}
                                                            [{{ .Labels.severity | toUpper }}]
                                                        {{ end }}{{ .Annotations.summary }}<br />
                                                        <strong>Description: </strong>{{ .Annotations.description }}<br />
                                                        <strong>Violated rule: </strong>{{ .Annotations.rule }}<br />
                                                        <strong>Details:</strong><br />
                                                        {{ with .Labels }}
                                                        {{ with .Remove (stringSlice "alertname" "ia" "instance" "node_type" "server") }}
                                                            {{ range .SortedPairs }}
                                                                • <strong>{{ .Name }}:</strong> {{ .Value }}<br />
                                                            {{ end }}
                                                        {{ end }}
                                                        {{ end }}
                                                    </td>
                                                </tr>
                                            {{ end }}
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <div class="footer">
                                <table width="100%%">
                                    <tr>
                                        <td class="aligncenter content-block"><a href="{{ .ExternalURL }}">Sent by {{ template "__alertmanager" . }}</a></td>
                                    </tr>
                                </table>
                            </div>
                        </div>
                    </td>
                </tr>
            </table>


            </body>
            </html>
        - send_resolved: false
          to: test2@test.test
          headers:
            Subject: '[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonLabels.alertname }}'
          html: |
            <!--
            Style and HTML derived from https://github.com/mailgun/transactional-email-templates


            The MIT License (MIT)

            Copyright (c) 2014 Mailgun

            Permission is hereby granted, free of charge, to any person obtaining a copy
            of this software and associated documentation files (the "Software"), to deal
            in the Software without restriction, including without limitation the rights
            to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
            copies of the Software, and to permit persons to whom the Software is
            furnished to do so, subject to the following conditions:

            The above copyright notice and this permission notice shall be included in all
            copies or substantial portions of the Software.

            THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
            IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
            FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
            AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
            LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
            OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
            SOFTWARE.
            -->
            <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
            <html xmlns="http://www.w3.org/1999/xhtml">
            <head>
                <meta name="viewport" content="width=device-width" />
                <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
                <title>{{ template "__subject" . }}</title>
                <style>
                    /* -------------------------------------
                        GLOBAL
                        A very basic CSS reset
                    ------------------------------------- */
                    * {
                        margin: 0;
                        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
                        box-sizing: border-box;
                        font-size: 14px;
                    }

                    img {
                        max-width: 100%%;
                    }

                    body {
                        -webkit-font-smoothing: antialiased;
                        -webkit-text-size-adjust: none;
                        width: 100%% !important;
                        height: 100%%;
                        line-height: 1.6em;
                        /* 1.6em * 14px = 22.4px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
                        /*line-height: 22px;*/
                    }

                    /* Let's make sure all tables have defaults */
                    table td {
                        vertical-align: top;
                    }

                    /* -------------------------------------
                        BODY & CONTAINER
                    ------------------------------------- */
                    body {
                        background-color: #f6f6f6;
                    }

                    .body-wrap {
                        background-color: #f6f6f6;
                        width: 100%%;
                    }

                    .container {
                        display: block !important;
                        max-width: 600px !important;
                        margin: 0 auto !important;
                        /* makes it centered */
                        clear: both !important;
                    }

                    .content {
                        max-width: 600px;
                        margin: 0 auto;
                        display: block;
                        padding: 20px;
                    }

                    /* -------------------------------------
                        HEADER, FOOTER, MAIN
                    ------------------------------------- */
                    .main {
                        background-color: #fff;
                        border: 1px solid #e9e9e9;
                        border-radius: 3px;
                    }

                    .content-wrap {
                        padding: 30px;
                    }

                    .content-block {
                        padding: 0 0 20px;
                    }

                    .footer {
                        width: 100%%;
                        clear: both;
                        color: #999;
                        padding: 20px;
                    }
                    .footer p, .footer a, .footer td {
                        color: #999;
                        font-size: 12px;
                    }

                    /* -------------------------------------
                        LINKS & BUTTONS
                    ------------------------------------- */
                    a {
                        color: #348eda;
                        text-decoration: underline;
                    }

                    .btn-primary {
                        text-decoration: none;
                        color: #FFF;
                        background-color: #3274D9;
                        border: solid #3274D9;
                        border-width: 10px 20px;
                        line-height: 2em;
                        /* 2em * 14px = 28px, use px to get airier line-height also in Thunderbird, and Yahoo!, Outlook.com, AOL webmail clients */
                        /*line-height: 28px;*/
                        font-weight: bold;
                        text-align: center;
                        cursor: pointer;
                        display: inline-block;
                        border-radius: 5px;
                        text-transform: capitalize;
                    }

                    /* -------------------------------------
                        OTHER STYLES THAT MIGHT BE USEFUL
                    ------------------------------------- */
                    .aligncenter {
                        text-align: center;
                    }

                    /* -------------------------------------
                        ALERTS
                        Change the class depending on warning email, good email or bad email
                    ------------------------------------- */
                    .alert {
                        font-size: 16px;
                        color: #fff;
                        font-weight: 500;
                        padding: 20px;
                        text-align: center;
                        border-radius: 3px 3px 0 0;
                    }
                    .alert a {
                        color: #fff;
                        text-decoration: none;
                        font-weight: 500;
                        font-size: 16px;
                    }
                    .alert.alert-notice {
                        background-color: #3274D9;
                    }
                    .alert.alert-warning {
                        background-color: #ECBB13;
                    }
                    .alert.alert-high {
                        background-color: #EB7B18;
                    }
                    .alert.alert-critical {
                        background-color: #E02F44;
                    }

                    /* -------------------------------------
                        RESPONSIVE AND MOBILE FRIENDLY STYLES
                    ------------------------------------- */
                    @media only screen and (max-width: 640px) {
                        body {
                            padding: 0 !important;
                        }

                        .container {
                            padding: 0 !important;
                            width: 100%% !important;
                        }

                        .content {
                            padding: 0 !important;
                        }

                        .content-wrap {
                            padding: 10px !important;
                        }
                    }
                </style>
            </head>

            <body itemscope itemtype="http://schema.org/EmailMessage">

            <table class="body-wrap">
                <tr>
                    <td></td>
                    <td class="container" width="600">
                        <div class="content">
                            <table class="main" width="100%%" cellpadding="0" cellspacing="0">
                                <tr>
                                    {{ if eq .CommonLabels.severity "notice" }}
                                    <td class="alert alert-notice">
                                        {{ else if eq .CommonLabels.severity "warning" }}
                                    <td class="alert alert-warning">
                                        {{ else if eq .CommonLabels.severity "error" }}
                                    <td class="alert alert-high">
                                        {{ else if eq .CommonLabels.severity "critical" }}
                                    <td class="alert alert-critical">
                                        {{ else }}
                                    <td class="alert alert-critical">
                                        {{ end }}
                                        You have {{ .Alerts | len }} {{ .CommonLabels.severity }} alert{{ if gt (len .Alerts) 1 }}s{{ end }} firing
                                    </td>
                                </tr>
                                <tr>
                                    <td class="content-wrap">
                                        <table width="100%%" cellpadding="0" cellspacing="0">
                                            <tr>
                                                <td align="center" class="content-block">
                                                    <a href='{{ template "__alertmanagerURL" . }}' class="btn-primary">View in {{ template "__alertmanager" . }}</a>
                                                </td>
                                            </tr>
                                            {{ range .Alerts.Firing }}
                                                <tr>
                                                    <td class="content-block">
                                                        <strong>------------------</strong><br />
                                                        <strong>Alert: </strong>
                                                        {{ if .Labels.severity }}
                                                            [{{ .Labels.severity | toUpper }}]
                                                        {{ end }}{{ .Annotations.summary }}<br />
                                                        <strong>Description: </strong>{{ .Annotations.description }}<br />
                                                        <strong>Violated rule: </strong>{{ .Annotations.rule }}<br />
                                                        <strong>Details:</strong><br />
                                                        {{ with .Labels }}
                                                        {{ with .Remove (stringSlice "alertname" "ia" "instance" "node_type" "server") }}
                                                            {{ range .SortedPairs }}
                                                                • <strong>{{ .Name }}:</strong> {{ .Value }}<br />
                                                            {{ end }}
                                                        {{ end }}
                                                        {{ end }}
                                                    </td>
                                                </tr>
                                            {{ end }}
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <div class="footer">
                                <table width="100%%">
                                    <tr>
                                        <td class="aligncenter content-block"><a href="{{ .ExternalURL }}">Sent by {{ template "__alertmanager" . }}</a></td>
                                    </tr>
                                </table>
                            </div>
                        </div>
                    </td>
                </tr>
            </table>


            </body>
            </html>
      pagerduty_configs:
        - send_resolved: false
          routing_key: ms-pagerduty-dev
templates: []
`, channel1.ID, channel2.ID, rule1.ID, rule2.ID, rule4.ID)) + "\n"
		assert.Equal(t, expected, actual, "actual:\n%s", actual)
	})
}

func TestGenerateReceivers(t *testing.T) {
	t.Parallel()

	chanMap := map[string]*models.Channel{
		"1": {
			ID:   "1",
			Type: models.Slack,
			SlackConfig: &models.SlackConfig{
				Channel: "channel1",
			},
		},
		"2": {
			ID:   "2",
			Type: models.Slack,
			SlackConfig: &models.SlackConfig{
				Channel: "channel2",
			},
		},
		"3": {
			ID:   "3",
			Type: models.Slack,
			SlackConfig: &models.SlackConfig{
				Channel: "channel3",
			},
			Disabled: true,
		},
	}
	recvSet := map[string]models.ChannelIDs{
		"1":   {"1"},
		"2":   {"2"},
		"1+2": {"1", "2"},
	}
	s := New(nil)
	actualR, err := s.generateReceivers(chanMap, recvSet)
	require.NoError(t, err)
	actual, err := yaml.Marshal(actualR)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
- name: "1"
  slack_configs:
    - send_resolved: false
      channel: channel1
      short_fields: false
      link_names: false
- name: 1+2
  slack_configs:
    - send_resolved: false
      channel: channel1
      short_fields: false
      link_names: false
    - send_resolved: false
      channel: channel2
      short_fields: false
      link_names: false
- name: "2"
  slack_configs:
    - send_resolved: false
      channel: channel2
      short_fields: false
      link_names: false
`) + "\n"
	assert.Equal(t, expected, string(actual), "actual:\n%s", actual)
}
