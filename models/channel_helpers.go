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
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

var (
	invalidConfigurationError = status.Error(codes.InvalidArgument, "Channel should contain only one type of channel configuration.")
)

func checkUniqueChannelID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Channel ID")
	}

	agent := &Channel{ID: id}
	switch err := q.Reload(agent); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Channel with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func checkEmailConfig(c *EmailConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Email config is empty.")
	}

	if len(c.To) == 0 {
		return status.Error(codes.InvalidArgument, "Email to field is empty.")
	}

	return nil
}

func checkPagerDutyConfig(c *PagerDutyConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Pager duty config is empty.")
	}

	if (c.RoutingKey == "" && c.ServiceKey == "") || (c.RoutingKey != "" && c.ServiceKey != "") {
		return status.Error(codes.InvalidArgument, "Exactly one key should be present in pager duty configuration.")
	}

	return nil
}

func checkSlackConfig(c *SlackConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Slack config is empty.")
	}

	if c.Channel == "" {
		return status.Error(codes.InvalidArgument, "Slack channel field is empty.")
	}

	return nil
}

func checkWebHookConfig(c *WebHookConfig) error {
	if c == nil {
		return status.Error(codes.InvalidArgument, "Webhook config is empty.")
	}

	if c.URL == "" {
		return status.Error(codes.InvalidArgument, "Webhook url field is empty.")
	}

	return nil
}

// FindChannels returns saved notification channels configuration.
func FindChannels(q *reform.Querier) ([]Channel, error) {
	rows, err := q.SelectAllFrom(ChannelTable, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select notification channels")
	}

	channels := make([]Channel, len(rows))
	for i, s := range rows {
		c := s.(*Channel)

		channels[i] = *c
	}

	return channels, nil
}

// FindChannelByID finds Channel by ID.
func FindChannelByID(q *reform.Querier, id string) (*Channel, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Channel ID.")
	}

	channel := &Channel{ID: id}
	switch err := q.Reload(channel); err {
	case nil:
		return channel, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Channel with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// FindChannelsByIDs finds channels by IDs.
func FindChannelsByIDs(q *reform.Querier, ids []string) ([]*Channel, error) {
	p := strings.Join(q.Placeholders(1, len(ids)), ", ")
	tail := fmt.Sprintf("WHERE id IN (%s)", p) //nolint:gosec
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	structs, err := q.SelectAllFrom(ChannelTable, tail, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]*Channel, len(structs))
	for i, s := range structs {
		res[i] = s.(*Channel)
	}
	return res, nil
}

// CreateChannelParams are params for creating new channel.
type CreateChannelParams struct {
	Summary string

	EmailConfig     *EmailConfig
	PagerDutyConfig *PagerDutyConfig
	SlackConfig     *SlackConfig
	WebHookConfig   *WebHookConfig

	Disabled bool
}

// CreateChannel persists notification channel.
func CreateChannel(q *reform.Querier, params *CreateChannelParams) (*Channel, error) {
	id := "/channel_id/" + uuid.New().String()

	if err := checkUniqueChannelID(q, id); err != nil {
		return nil, err
	}

	if params.Summary == "" {
		return nil, status.Error(codes.InvalidArgument, "Channel summary can't be empty.")
	}

	row := &Channel{
		ID:       id,
		Summary:  params.Summary,
		Disabled: params.Disabled,
	}

	if params.EmailConfig != nil {
		if err := checkEmailConfig(params.EmailConfig); err != nil {
			return nil, err
		}
		row.Type = Email
		row.EmailConfig = params.EmailConfig
	}

	if params.PagerDutyConfig != nil {
		if row.Type != "" {
			return nil, invalidConfigurationError
		}

		if err := checkPagerDutyConfig(params.PagerDutyConfig); err != nil {
			return nil, err
		}
		row.Type = PagerDuty
		row.PagerDutyConfig = params.PagerDutyConfig
	}

	if params.SlackConfig != nil {
		if row.Type != "" {
			return nil, invalidConfigurationError
		}
		if err := checkSlackConfig(params.SlackConfig); err != nil {
			return nil, err
		}
		row.Type = Slack
		row.SlackConfig = params.SlackConfig
	}

	if params.WebHookConfig != nil {
		if row.Type != "" {
			return nil, invalidConfigurationError
		}
		if err := checkWebHookConfig(params.WebHookConfig); err != nil {
			return nil, err
		}
		row.Type = WebHook
		row.WebHookConfig = params.WebHookConfig
	}

	if row.Type == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing channel configuration.")
	}

	err := q.Insert(row)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create notifications channel")
	}

	return row, nil
}

// ChangeChannelParams is params for changing existing channel.
type ChangeChannelParams struct {
	EmailConfig     *EmailConfig
	PagerDutyConfig *PagerDutyConfig
	SlackConfig     *SlackConfig
	WebHookConfig   *WebHookConfig

	Disabled bool
}

// ChangeChannel updates existing notifications channel.
func ChangeChannel(q *reform.Querier, channelID string, params *ChangeChannelParams) (*Channel, error) {
	row, err := FindChannelByID(q, channelID)
	if err != nil {
		return nil, err
	}

	// remove previous configuration
	row.EmailConfig = nil
	row.PagerDutyConfig = nil
	row.SlackConfig = nil
	row.WebHookConfig = nil

	if params.EmailConfig != nil {
		if err := checkEmailConfig(params.EmailConfig); err != nil {
			return nil, err
		}
		row.Type = Email
		row.EmailConfig = params.EmailConfig
	}

	if params.PagerDutyConfig != nil {
		if row.Type != "" {
			return nil, invalidConfigurationError
		}

		if err := checkPagerDutyConfig(params.PagerDutyConfig); err != nil {
			return nil, err
		}
		row.Type = PagerDuty
		row.PagerDutyConfig = params.PagerDutyConfig
	}

	if params.SlackConfig != nil {
		if row.Type != "" {
			return nil, invalidConfigurationError
		}
		if err := checkSlackConfig(params.SlackConfig); err != nil {
			return nil, err
		}
		row.Type = Slack
		row.SlackConfig = params.SlackConfig
	}

	if params.WebHookConfig != nil {
		if row.Type != "" {
			return nil, invalidConfigurationError
		}
		if err := checkWebHookConfig(params.WebHookConfig); err != nil {
			return nil, err
		}
		row.Type = WebHook
		row.WebHookConfig = params.WebHookConfig
	}

	row.Disabled = params.Disabled

	err = q.Update(row)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update notifications channel")
	}

	return row, nil
}

// RemoveChannel removes notification channel with specified id.
func RemoveChannel(q *reform.Querier, id string) error {
	if _, err := FindChannelByID(q, id); err != nil {
		return err
	}

	err := q.Delete(&Channel{ID: id})
	if err != nil {
		return errors.Wrap(err, "failed to delete notification channel")
	}
	return nil
}
