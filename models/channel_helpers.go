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

// validateChannel validates notification channel.
func validateChannel(ch *Channel) error {
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

// FindChannels returns saved notification channels configuration.
func FindChannels(q *reform.Querier) ([]Channel, error) {
	structs, err := q.SelectAllFrom(notificationChannelTable, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select notification channels")

	}

	channels := make([]Channel, len(structs))
	for i, s := range structs {
		c, err := notificationChannelToChannel(s.(*notificationChannel))
		if err != nil {
			return nil, err
		}
		channels[i] = *c
	}

	return channels, nil
}

// CreateChannel persists notification channel.
func CreateChannel(q *reform.Querier, c *Channel) error {
	if err := validateChannel(c); err != nil {
		return err
	}

	nc, err := channelToNotificationChannel(c)
	if err != nil {
		return err
	}

	err = q.Insert(nc)
	if err != nil {
		return errors.Wrap(err, "failed to create notifications channel")
	}

	return nil
}

// ChangeChannel updates existing notifications channel.
func ChangeChannel(q *reform.Querier, c *Channel) error {
	if err := validateChannel(c); err != nil {
		return errors.Wrap(err, "channel validation failed")
	}

	nc, err := channelToNotificationChannel(c)
	if err != nil {
		return err
	}

	err = q.Update(nc)
	if err != nil {
		return errors.Wrap(err, "failed to create notifications channel")
	}

	return nil
}

// RemoveChannel removes notification channel with specified id.
func RemoveChannel(q *reform.Querier, id string) error {
	err := q.Delete(&notificationChannel{ID: id})
	if err != nil {
		return errors.Wrap(err, "failed to delete notifications channel")
	}
	return nil
}

func channelToNotificationChannel(c *Channel) (*notificationChannel, error) {
	nc := &notificationChannel{
		ID:       c.ID,
		Type:     c.Type,
		Disabled: false,
	}

	switch c.Type {
	case Email:
		b, err := json.Marshal(c.EmailConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall email configuration")
		}
		nc.EmailConfig = &b
	case PagerDuty:
		b, err := json.Marshal(c.PagerDutyConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall pager duty configuration")
		}
		nc.PagerDutyConfig = &b
	case Slack:
		b, err := json.Marshal(c.SlackConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall slack configuration")
		}
		nc.SlackConfig = &b
	case WebHook:
		b, err := json.Marshal(c.WebHookConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall webhook configuration")
		}
		nc.WebHookConfig = &b
	}

	return nc, nil
}

func notificationChannelToChannel(nc *notificationChannel) (*Channel, error) {
	c := &Channel{
		ID:       nc.ID,
		Type:     nc.Type,
		Disabled: nc.Disabled,
	}

	switch nc.Type {
	case Email:
		c.EmailConfig = &EmailConfig{}
		err := json.Unmarshal(*nc.EmailConfig, c.EmailConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall email configuration")
		}
	case PagerDuty:
		c.PagerDutyConfig = &PagerDutyConfig{}
		err := json.Unmarshal(*nc.PagerDutyConfig, c.PagerDutyConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall pager duty configuration")
		}
	case Slack:
		c.SlackConfig = &SlackConfig{}
		err := json.Unmarshal(*nc.SlackConfig, c.SlackConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall slack configuration")
		}
	case WebHook:
		c.WebHookConfig = &WebHookConfig{}
		err := json.Unmarshal(*nc.WebHookConfig, c.WebHookConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshall webhook configuration")
		}
	}

	return c, nil
}
