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

	t.Run("save", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		c := models.Channel{
			Id:   "some_id",
			Type: models.Email,
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		}

		err = models.SaveChannel(tx, &c)
		require.NoError(t, err)

		channels, err := models.GetChannels(tx)
		require.NoError(t, err)
		require.Len(t, channels, 1)

		assert.Equal(t, c, channels[0])
	})

	t.Run("update", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		c := models.Channel{
			Id:   "some_id",
			Type: models.Email,
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		}

		err = models.SaveChannel(tx, &c)
		require.NoError(t, err)

		c.EmailConfig.To = []string{"test2@test.test"}

		err = models.UpdateChannel(tx, &c)
		require.NoError(t, err)

		cs, err := models.GetChannels(tx)
		require.NoError(t, err)
		assert.Len(t, cs, 1)
		assert.Equal(t, c, cs[0])
	})

	t.Run("delete", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		c := models.Channel{
			Id:   "some_id",
			Type: models.Email,
			EmailConfig: &models.EmailConfig{
				To: []string{"test@test.test"},
			},
			Disabled: false,
		}

		err = models.SaveChannel(tx, &c)
		require.NoError(t, err)

		err = models.RemoveChannel(tx, c.Id)
		require.NoError(t, err)

		cs, err := models.GetChannels(tx)
		require.NoError(t, err)
		assert.Len(t, cs, 0)
	})

}

func TestChannelValidation(t *testing.T) {
	tests := []struct {
		name     string
		channel  models.Channel
		errorMsg string
	}{
		{
			name: "normal",
			channel: models.Channel{
				Id:   "some_id",
				Type: models.Email,
				EmailConfig: &models.EmailConfig{
					To: []string{"test@test.test"},
				},
				Disabled: false,
			},
			errorMsg: "",
		},
		{
			name: "missing id",
			channel: models.Channel{
				Id:   "",
				Type: models.Email,
				EmailConfig: &models.EmailConfig{
					To: []string{"test@test.test"},
				},
				Disabled: false,
			},
			errorMsg: "notification channel id is empty",
		},
		{
			name: "unknown type",
			channel: models.Channel{
				Id:   "some_id",
				Type: "qwerty",
				EmailConfig: &models.EmailConfig{
					To: []string{"test@test.test"},
				},
				Disabled: false,
			},
			errorMsg: "unknown channel type qwerty",
		},
		{
			name: "missing type",
			channel: models.Channel{
				Id:   "some_id",
				Type: "",
				EmailConfig: &models.EmailConfig{
					To: []string{"test@test.test"},
				},
				Disabled: false,
			},
			errorMsg: "notification channel type is empty",
		},
		{
			name: "missing email config",
			channel: models.Channel{
				Id:          "some_id",
				Type:        models.Email,
				EmailConfig: nil,
				Disabled:    false,
			},
			errorMsg: "email config is empty",
		},
		{
			name: "missing slack config",
			channel: models.Channel{
				Id:       "some_id",
				Type:     models.Slack,
				Disabled: false,
			},
			errorMsg: "slack config is empty",
		},
		{
			name: "missing webhook config",
			channel: models.Channel{
				Id:       "some_id",
				Type:     models.WebHook,
				Disabled: false,
			},
			errorMsg: "webhook config is empty",
		},
		{
			name: "missing to field in email configuration",
			channel: models.Channel{
				Id:   "some_id",
				Type: models.Email,
				EmailConfig: &models.EmailConfig{
					To: nil,
				},
				Disabled: false,
			},
			errorMsg: "email to field is empty",
		},
		{
			name: "missing channel in slack configuration",
			channel: models.Channel{
				Id:   "some_id",
				Type: models.Slack,
				SlackConfig: &models.SlackConfig{
					Channel: "",
				},
				Disabled: false,
			},
			errorMsg: "slack channel field is empty",
		},
		{
			name: "missing url in webhook configuration",
			channel: models.Channel{
				Id:   "some_id",
				Type: models.WebHook,
				WebHookConfig: &models.WebHookConfig{
					Url: "",
				},
				Disabled: false,
			},
			errorMsg: "webhook url field is empty",
		},
		{
			name: "type doesn't match actual configuration",
			channel: models.Channel{
				Id:   "some_id",
				Type: models.Slack,
				EmailConfig: &models.EmailConfig{
					To: []string{"test@test.test"},
				},
				Disabled: false,
			},
			errorMsg: "slack channel should has only slack configuration",
		},
		{
			name: "multiple configurations",
			channel: models.Channel{
				Id:   "some_id",
				Type: models.Email,
				EmailConfig: &models.EmailConfig{
					To: []string{"test@test.test"},
				},
				WebHookConfig: &models.WebHookConfig{
					Url: "example.com",
				},
				Disabled: false,
			},
			errorMsg: "email channel should has only email configuration",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := models.ValidateChannel(&test.channel)
			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}
