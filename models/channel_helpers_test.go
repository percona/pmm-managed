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

package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestNotificationChannels(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("create", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params := models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		}

		expected, err := models.CreateChannel(q, &params)
		require.NoError(t, err)
		assert.Equal(t, models.Email, expected.Type)
		assert.Equal(t, params.Summary, expected.Summary)
		assert.Equal(t, params.Disabled, expected.Disabled)
		assert.Equal(t, params.EmailConfig.SendResolved, expected.EmailConfig.SendResolved)
		assert.EqualValues(t, params.EmailConfig.SendResolved, expected.EmailConfig.SendResolved)

		actual, err := models.FindChannelByID(q, expected.ID)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("change", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		createParams := &models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To:           []string{"test@test.test"},
				SendResolved: false,
			},
			Disabled: false,
		}

		created, err := models.CreateChannel(q, createParams)
		require.NoError(t, err)

		updateParams := &models.ChangeChannelParams{
			Summary: "completely new summary",
			SlackConfig: &models.SlackConfig{
				SendResolved: true,
				Channel:      "general",
			},
			Disabled: true,
		}

		updated, err := models.ChangeChannel(q, created.ID, updateParams)
		require.NoError(t, err)
		assert.Equal(t, models.Slack, updated.Type)
		assert.Equal(t, updateParams.Summary, updated.Summary)
		assert.Equal(t, updateParams.Disabled, updated.Disabled)
		assert.Nil(t, updated.EmailConfig)
		assert.Equal(t, updateParams.SlackConfig.Channel, updated.SlackConfig.Channel)
		assert.EqualValues(t, updateParams.SlackConfig.SendResolved, updated.SlackConfig.SendResolved)

		actual, err := models.FindChannelByID(q, created.ID)
		require.NoError(t, err)
		assert.Equal(t, updated, actual)
	})

	t.Run("remove", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params := &models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		}

		c, err := models.CreateChannel(q, params)
		require.NoError(t, err)

		err = models.RemoveChannel(q, c.ID)
		require.NoError(t, err)

		cs, err := models.FindChannels(q)
		require.NoError(t, err)
		assert.Len(t, cs, 0)
	})

	t.Run("find", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params1 := models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		}
		expected1, err := models.CreateChannel(q, &params1)
		require.NoError(t, err)

		params2 := models.CreateChannelParams{
			Summary: "another summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test2@test.test"},
			},
			Disabled: true,
		}
		expected2, err := models.CreateChannel(q, &params2)
		require.NoError(t, err)

		actual, err := models.FindChannels(q)
		require.NoError(t, err)
		var found1, found2 bool
		for _, channel := range actual {
			if channel.ID == expected1.ID {
				found1 = true
			}
			if channel.ID == expected2.ID {
				found2 = true
			}
		}

		assert.True(t, found1, "Fist channel not found")
		assert.True(t, found2, "Second channel not found")
	})
}

func TestChannelValidation(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	tests := []struct {
		name     string
		channel  models.CreateChannelParams
		errorMsg string
	}{{
		name: "normal email config",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		},
		errorMsg: "",
	}, {
		name: "normal pager duty config",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			PagerDutyConfig: &models.PagerDutyConfig{
				SendResolved: false,
				RoutingKey:   "some key",
			},
			Disabled: false,
		},
		errorMsg: "",
	}, {
		name: "normal slack config",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			SlackConfig: &models.SlackConfig{
				Channel: "channel",
			},
			Disabled: false,
		},
		errorMsg: "",
	}, {
		name: "normal webhook config",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			WebHookConfig: &models.WebHookConfig{
				URL: "test.test",
			},
			Disabled: false,
		},
		errorMsg: "",
	}, {
		name: "missing summary",
		channel: models.CreateChannelParams{
			Summary: "",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Channel summary can't be empty.",
	}, {
		name: "missing email config",
		channel: models.CreateChannelParams{
			Summary:     "some summary",
			EmailConfig: nil,
			Disabled:    false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Missing channel configuration.",
	}, {
		name: "missing pager duty config",
		channel: models.CreateChannelParams{
			Summary:         "some summary",
			PagerDutyConfig: nil,
			Disabled:        false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Missing channel configuration.",
	}, {
		name: "missing slack config",
		channel: models.CreateChannelParams{
			Summary:  "some summary",
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Missing channel configuration.",
	}, {
		name: "missing webhook config",
		channel: models.CreateChannelParams{
			Summary:  "some summary",
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Missing channel configuration.",
	}, {
		name: "missing to field in email configuration",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: nil,
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Email to field is empty.",
	}, {
		name: "no keys set in pager duty config",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			PagerDutyConfig: &models.PagerDutyConfig{
				SendResolved: false,
				RoutingKey:   "",
				ServiceKey:   "",
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Exactly one key should be present in pager duty configuration.",
	}, {
		name: "both keys set in pager duty config",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			PagerDutyConfig: &models.PagerDutyConfig{
				SendResolved: false,
				RoutingKey:   "some key",
				ServiceKey:   "some key",
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Exactly one key should be present in pager duty configuration.",
	}, {
		name: "missing channel in slack configuration",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			SlackConfig: &models.SlackConfig{
				Channel: "",
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Slack channel field is empty.",
	}, {
		name: "missing url in webhook configuration",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			WebHookConfig: &models.WebHookConfig{
				URL: "",
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Webhook url field is empty.",
	}, {
		name: "multiple configurations",
		channel: models.CreateChannelParams{
			Summary: "some summary",
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			WebHookConfig: &models.WebHookConfig{
				URL: "example.com",
			},
			Disabled: false,
		},
		errorMsg: "rpc error: code = InvalidArgument desc = Channel should contain only one type of channel configuration.",
	}}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			tx, err := db.Begin()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, tx.Rollback())
			}()

			q := tx.Querier

			c, err := models.CreateChannel(q, &test.channel)
			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, c)
		})
	}
}
