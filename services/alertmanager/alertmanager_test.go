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
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/percona-platform/saas/pkg/alert"
	"github.com/percona-platform/saas/pkg/common"
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

func TestAlertmanager(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	svc := New(db)

	require.NoError(t, svc.IsReady(context.Background()))
}

func TestPopulateConfig(t *testing.T) {
	var cfg alertmanager.Config
	buf, err := ioutil.ReadFile(alertmanagerBaseConfigPath)
	require.NoError(t, err)
	err = yaml.Unmarshal(buf, &cfg)
	require.NoError(t, err)

	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	tx, err := db.Begin()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback())
	}()

	q := tx.Querier

	channel1, err := models.CreateChannel(q, &models.CreateChannelParams{
		Summary: "some summary",
		EmailConfig: &models.EmailConfig{
			To: []string{"test@test.test"},
		},
		Disabled: false,
	})
	require.NoError(t, err)

	channel2, err := models.CreateChannel(q, &models.CreateChannelParams{
		Summary: "some summary",
		PagerDutyConfig: &models.PagerDutyConfig{
			RoutingKey: "ms-pagerduty-dev",
		},
		Disabled: false,
	})
	require.NoError(t, err)

	templateName := gofakeit.UUID()
	_, err = models.CreateTemplate(q, &models.CreateTemplateParams{
		Template: &alert.Template{
			Name:    templateName,
			Version: 1,
			Summary: gofakeit.Quote(),
			Tiers:   []common.Tier{common.Anonymous},
			Expr:    gofakeit.Quote(),
			Params: []alert.Parameter{{
				Name:    gofakeit.UUID(),
				Summary: gofakeit.Quote(),
				Unit:    gofakeit.Letter(),
				Type:    alert.Float,
				Range:   []interface{}{float64(10), float64(100)},
				Value:   float64(50),
			}},
			For:         3,
			Severity:    common.Warning,
			Labels:      map[string]string{"foo": "bar"},
			Annotations: nil,
		},
		Source: "USER_FILE",
	})
	require.NoError(t, err)

	rule, err := models.CreateRule(q, &models.CreateRuleParams{
		TemplateName: templateName,
		Disabled:     true,
		RuleParams: []models.RuleParam{
			{
				Name:       "test",
				Type:       models.Float,
				FloatValue: 3.14,
			},
		},
		For:          5 * time.Second,
		Severity:     common.Warning,
		CustomLabels: map[string]string{"foo": "bar"},
		Filters:      []models.Filter{{Type: models.Equal, Key: "value", Val: "10"}},
		ChannelIDs:   []string{channel1.ID, channel2.ID},
	})
	require.NoError(t, err)

	settings, err := models.UpdateSettings(q, &models.ChangeSettingsParams{
		EmailAlertingSettings: &models.EmailAlertingSettings{
			From:      tests.GenEmail(t),
			Smarthost: "0.0.0.0:80",
			Hello:     "host",
			Username:  "user",
			Password:  "password",
			Identity:  "id",
			Secret:    "secret",
		},
		SlackAlertingSettings: &models.SlackAlertingSettings{
			URL: gofakeit.URL(),
		},
	})
	require.NoError(t, err)

	svc := New(db)
	err = svc.populateConfig(&cfg)
	require.NoError(t, err)

	assert.Len(t, cfg.Receivers, 2)
	assert.Equal(t, "empty", cfg.Receivers[0].Name)
	assert.Equal(t, channel1.ID+" + "+channel2.ID, cfg.Receivers[1].Name)
	assert.NotNil(t, cfg.Receivers[1].EmailConfigs)
	assert.NotNil(t, cfg.Receivers[1].PagerdutyConfigs)
	assert.Equal(t, "empty", cfg.Route.Receiver)
	assert.Len(t, cfg.Route.Routes, 1)
	assert.Equal(t, rule.ID, cfg.Route.Routes[0].Match["rule_id"])
	// check global config
	assert.Equal(t, cfg.Global.SMTPFrom, settings.IntegratedAlerting.EmailAlertingSettings.From)
	assert.Equal(t, cfg.Global.SMTPHello, settings.IntegratedAlerting.EmailAlertingSettings.Hello)
	assert.Equal(t, cfg.Global.SMTPAuthUsername, settings.IntegratedAlerting.EmailAlertingSettings.Username)
	assert.Equal(t, cfg.Global.SMTPAuthPassword, settings.IntegratedAlerting.EmailAlertingSettings.Password)
	assert.Equal(t, cfg.Global.SMTPAuthIdentity, settings.IntegratedAlerting.EmailAlertingSettings.Identity)
	assert.Equal(t, cfg.Global.SMTPAuthSecret, settings.IntegratedAlerting.EmailAlertingSettings.Secret)

	smarthost := strings.Split(settings.IntegratedAlerting.EmailAlertingSettings.Smarthost, ":")
	assert.Equal(t, cfg.Global.SMTPSmarthost.Host, smarthost[0])
	assert.Equal(t, cfg.Global.SMTPSmarthost.Port, smarthost[1])
	assert.Equal(t, cfg.Global.SlackAPIURL, settings.IntegratedAlerting.SlackAlertingSettings.URL)
}
