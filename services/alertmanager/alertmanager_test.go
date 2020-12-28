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
	err = svc.validateConfig(context.Background(), b)
	require.NoError(t, err)
	return string(b)
}

func TestPopulateConfig(t *testing.T) {
	New(nil).GenerateBaseConfigs() // this method should not use database

	tests.SetTestIDReader(t)

	t.Run("without receivers and routes", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		svc := New(db)

		// add fake setting to check for overwrite
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
templates: []
		`) + "\n"
		assert.Equal(t, expected, actual, "actual:\n%s", actual)
	})

	t.Run("with receivers and routes", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		svc := New(db)

		channel1, err := models.CreateChannel(db.Querier, &models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test", "test2@test.test"},
			},
			Disabled: false,
		})
		require.NoError(t, err)

		channel2, err := models.CreateChannel(db.Querier, &models.CreateChannelParams{
			Summary: "some summary",
			PagerDutyConfig: &models.PagerDutyConfig{
				RoutingKey: "ms-pagerduty-dev",
			},
			Disabled: false,
		})
		require.NoError(t, err)

		templateName := gofakeit.UUID()
		_, err = models.CreateTemplate(db.Querier, &models.CreateTemplateParams{
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

		_, err = models.CreateRule(db.Querier, &models.CreateRuleParams{
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

		// create another rule with same channelIDs to check for redundant receivers.
		_, err = models.CreateRule(db.Querier, &models.CreateRuleParams{
			TemplateName: templateName,
			Disabled:     true,
			RuleParams: []models.RuleParam{{
				Name:       "test",
				Type:       models.Float,
				FloatValue: 3.14,
			}},
			For:          5 * time.Second,
			Severity:     common.Warning,
			CustomLabels: map[string]string{"foo": "bar"},
			Filters:      []models.Filter{{Type: models.Equal, Key: "value", Val: "10"}},
			ChannelIDs:   []string{channel1.ID, channel2.ID},
		})
		require.NoError(t, err)

		_, err = models.UpdateSettings(db.Querier, &models.ChangeSettingsParams{
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

		cfg := svc.loadBaseConfig()
		err = svc.populateConfig(cfg)
		require.NoError(t, err)

		actual := marshalAndValidate(t, svc, svc.loadBaseConfig())
		expected := strings.TrimSpace(`
		`) + "\n"
		assert.Equal(t, expected, actual, "actual:\n%s", actual)
	})
}
