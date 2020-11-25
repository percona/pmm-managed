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

package models

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// SaveChannel persists notification channel.
func SaveChannel(q reform.DBTX, c *Channel) error {
	if err := ValidateChannel(c); err != nil {
		return err
	}

	b, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshall notification channel")
	}

	_, err = q.Exec("INSERT INTO notification_channels (id, channel) VALUES ($1, $2)", c.ID, b)
	if err != nil {
		return errors.Wrap(err, "failed to create notifications channel")
	}

	return nil
}

// UpdateChannel updates existing notifications channel.
func UpdateChannel(q reform.DBTX, c *Channel) error {
	if err := ValidateChannel(c); err != nil {
		return errors.Wrap(err, "channel validation failed")
	}
	b, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshall notification channel")
	}

	_, err = q.Exec("UPDATE notification_channels SET channel=$1 WHERE id=$2", b, c.ID)
	if err != nil {
		return errors.Wrap(err, "failed to create notifications channel")
	}

	return nil
}

// RemoveChannel removes notification channel with specified id.
func RemoveChannel(q reform.DBTX, id string) error {
	_, err := q.Exec("DELETE FROM notification_channels WHERE id=$1", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete notifications channel")
	}
	return nil
}

// GetChannels returns saved notification channels configuration.
func GetChannels(q reform.DBTX) ([]Channel, error) {
	rows, err := q.Query("SELECT channel FROM notification_channels")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select notification channels")
	}

	var channels []Channel
	for rows.Next() {
		var b []byte
		if err = rows.Scan(&b); err != nil {
			break
		}

		var channel Channel
		if err = json.Unmarshal(b, &channel); err != nil {
			break
		}
		channels = append(channels, channel)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errors.Wrap(closeErr, "failed to close rows")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to read notification channels")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to scan rows")
	}

	return channels, nil
}

// ValidateChannel validates notification channel.
func ValidateChannel(ch *Channel) error {
	if ch.ID == "" {
		return status.Error(codes.InvalidArgument, "Notification channel id is empty")
	}

	switch ch.Type {
	case Email:
		if ch.SlackConfig != nil || ch.WebHookConfig != nil || ch.PagerDutyConfig != nil {
			return status.Error(codes.InvalidArgument, "Email channel should has only email configuration")
		}

		return validateEmailConfig(ch.EmailConfig)
	case PagerDuty:
		if ch.EmailConfig != nil || ch.SlackConfig != nil || ch.WebHookConfig != nil {
			return status.Error(codes.InvalidArgument, "Pager duty channel should has only email configuration")
		}

		return validatePagerDutyConfig(ch.PagerDutyConfig)
	case Slack:
		if ch.EmailConfig != nil || ch.WebHookConfig != nil || ch.PagerDutyConfig != nil {
			return status.Error(codes.InvalidArgument, "Slack channel should has only slack configuration")
		}

		return validateSlackConfig(ch.SlackConfig)
	case WebHook:
		if ch.SlackConfig != nil || ch.EmailConfig != nil || ch.PagerDutyConfig != nil {
			return status.Error(codes.InvalidArgument, "Webhook channel should has only webhook configuration")
		}

		return validateWebHookConfig(ch.WebHookConfig)
	case "":
		return status.Error(codes.InvalidArgument, "Notification channel type is empty")
	default:
		return status.Error(codes.InvalidArgument, fmt.Sprintf("Unknown channel type %s", ch.Type))
	}
}

func validateEmailConfig(c *EmailConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Email config is empty")
	}

	if len(c.To) == 0 {
		return status.Error(codes.InvalidArgument, "Email to field is empty")
	}

	return nil
}

func validatePagerDutyConfig(c *PagerDutyConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Pager duty config is empty")
	}

	if c.RoutingKey == "" {
		return status.Error(codes.InvalidArgument, "Pager duty routing key is empty")
	}

	if c.ServiceKey == "" {
		return status.Error(codes.InvalidArgument, "Pager duty service key is empty")
	}

	return nil
}

func validateSlackConfig(c *SlackConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Slack config is empty")
	}

	if c.Channel == "" {
		return status.Error(codes.InvalidArgument, "Slack channel field is empty")
	}

	return nil
}

func validateWebHookConfig(c *WebHookConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Webhook config is empty")
	}

	if c.URL == "" {
		return status.Error(codes.InvalidArgument, "Webhook url field is empty")
	}

	return nil
}
