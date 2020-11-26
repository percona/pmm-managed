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

package ia

import (
	"context"

	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/ia"
)

// ChannelsService represents integrated alerting channels API.
type ChannelsService struct {
	ia alertingService
}

// NewChannelsService creates new channels API service.
func NewChannelsService(ia *ia.Service) *ChannelsService {
	return &ChannelsService{
		ia: ia,
	}
}

// ListChannels returns list of available channels.
func (s *ChannelsService) ListChannels(ctx context.Context, request *iav1beta1.ListChannelsRequest) (*iav1beta1.ListChannelsResponse, error) {
	channels, err := s.ia.ListChannels()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get notification channels")
	}

	res := make([]*iav1beta1.Channel, len(channels))
	for i, channel := range channels {
		c := &iav1beta1.Channel{
			ChannelId: channel.ID,
			Disabled:  channel.Disabled,
		}

		switch channel.Type {
		case models.Email:
			config := channel.EmailConfig
			c.Channel = &iav1beta1.Channel_EmailConfig{
				EmailConfig: &iav1beta1.EmailConfig{
					SendResolved: config.SendResolved,
					To:           config.To,
				},
			}
		case models.PagerDuty:
			config := channel.PagerDutyConfig
			c.Channel = &iav1beta1.Channel_PagerdutyConfig{
				PagerdutyConfig: &iav1beta1.PagerDutyConfig{
					SendResolved: config.SendResolved,
					RoutingKey:   config.RoutingKey,
					ServiceKey:   config.ServiceKey,
				},
			}
		case models.Slack:
			config := channel.SlackConfig
			c.Channel = &iav1beta1.Channel_SlackConfig{
				SlackConfig: &iav1beta1.SlackConfig{
					SendResolved: config.SendResolved,
					Channel:      config.Channel,
				},
			}
		case models.WebHook:
			config := channel.WebHookConfig
			c.Channel = &iav1beta1.Channel_WebhookConfig{
				WebhookConfig: &iav1beta1.WebhookConfig{
					SendResolved: config.SendResolved,
					Url:          config.URL,
					HttpConfig:   convertModelToHTTPConfig(config.HTTPConfig),
					MaxAlerts:    config.MaxAlerts,
				},
			}
		default:
			return nil, errors.Wrapf(err, "unknown notification channel type %s", channel.Type)
		}

		res[i] = c
	}

	return &iav1beta1.ListChannelsResponse{Channels: res}, nil
}

// AddChannel adds new notification channel.
func (s *ChannelsService) AddChannel(ctx context.Context, req *iav1beta1.AddChannelRequest) (*iav1beta1.AddChannelResponse, error) {
	channel := &models.Channel{
		ID:       req.GetChannelId(),
		Disabled: req.GetDisabled(),
	}

	if emailConf := req.GetEmailConfig(); emailConf != nil {
		channel.Type = models.Email
		channel.EmailConfig = &models.EmailConfig{
			SendResolved: emailConf.SendResolved,
			To:           emailConf.To,
		}
	}
	if pagerDutyConf := req.GetPagerdutyConfig(); pagerDutyConf != nil {
		if channel.Type != "" {
			return nil, status.Error(codes.InvalidArgument, "Request should contain only one type of channel configuration")
		}
		channel.Type = models.PagerDuty
		channel.PagerDutyConfig = &models.PagerDutyConfig{
			SendResolved: pagerDutyConf.SendResolved,
			RoutingKey:   pagerDutyConf.RoutingKey,
			ServiceKey:   pagerDutyConf.ServiceKey,
		}
	}
	if slackConf := req.GetSlackConfig(); slackConf != nil {
		if channel.Type != "" {
			return nil, status.Error(codes.InvalidArgument, "Request should contain only one type of channel configuration")
		}
		channel.Type = models.Slack
		channel.SlackConfig = &models.SlackConfig{
			SendResolved: slackConf.SendResolved,
			Channel:      slackConf.Channel,
		}
	}
	if webhookConf := req.GetWebhookConfig(); webhookConf != nil {
		if channel.Type != "" {
			return nil, status.Error(codes.InvalidArgument, "Request should contain only one type of channel configuration")
		}
		channel.Type = models.WebHook
		channel.WebHookConfig = &models.WebHookConfig{
			SendResolved: webhookConf.SendResolved,
			URL:          webhookConf.Url,
			MaxAlerts:    webhookConf.MaxAlerts,
			HTTPConfig:   convertHTTPConfigToModel(webhookConf.HttpConfig),
		}
	}

	err := s.ia.AddChannel(channel)
	if err != nil {
		return nil, err
	}

	return &iav1beta1.AddChannelResponse{}, nil
}

// ChangeChannel changes existing notification channel.
func (s *ChannelsService) ChangeChannel(ctx context.Context, req *iav1beta1.ChangeChannelRequest) (*iav1beta1.ChangeChannelResponse, error) {
	channel := &models.Channel{
		ID:       req.GetChannelId(),
		Disabled: req.GetDisabled(),
	}

	if emailConf := req.GetEmailConfig(); emailConf != nil {
		channel.Type = models.Email
		channel.EmailConfig = &models.EmailConfig{
			SendResolved: emailConf.SendResolved,
			To:           emailConf.To,
		}
	}
	if pagerDutyConf := req.GetPagerdutyConfig(); pagerDutyConf != nil {
		if channel.Type != "" {
			return nil, status.Error(codes.InvalidArgument, "Request should contain only one type of channel configuration")
		}
		channel.Type = models.PagerDuty
		channel.PagerDutyConfig = &models.PagerDutyConfig{
			SendResolved: pagerDutyConf.SendResolved,
			RoutingKey:   pagerDutyConf.RoutingKey,
			ServiceKey:   pagerDutyConf.ServiceKey,
		}
	}
	if slackConf := req.GetSlackConfig(); slackConf != nil {
		if channel.Type != "" {
			return nil, status.Error(codes.InvalidArgument, "Request should contain only one type of channel configuration")
		}
		channel.Type = models.Slack
		channel.SlackConfig = &models.SlackConfig{
			SendResolved: slackConf.SendResolved,
			Channel:      slackConf.Channel,
		}
	}
	if webhookConf := req.GetWebhookConfig(); webhookConf != nil {
		if channel.Type != "" {
			return nil, status.Error(codes.InvalidArgument, "Request should contain only one type of channel configuration")
		}
		channel.Type = models.WebHook
		channel.WebHookConfig = &models.WebHookConfig{
			SendResolved: webhookConf.SendResolved,
			URL:          webhookConf.Url,
			MaxAlerts:    webhookConf.MaxAlerts,
			HTTPConfig:   convertHTTPConfigToModel(webhookConf.HttpConfig),
		}
	}

	err := s.ia.ChangeChannel(channel)
	if err != nil {
		return nil, err
	}

	return &iav1beta1.ChangeChannelResponse{}, nil
}

// RemoveChannel removes notification channel.
func (s *ChannelsService) RemoveChannel(ctx context.Context, req *iav1beta1.RemoveChannelRequest) (*iav1beta1.RemoveChannelResponse, error) {
	if err := s.ia.RemoveChannel(req.ChannelId); err != nil {
		return nil, errors.Wrap(err, "failed to remove notification channel")
	}

	return &iav1beta1.RemoveChannelResponse{}, nil
}

func convertHTTPConfigToModel(config *iav1beta1.HTTPConfig) *models.HTTPConfig {
	var res *models.HTTPConfig
	if config != nil {
		res = &models.HTTPConfig{
			BearerToken:     config.BearerToken,
			BearerTokenFile: config.BearerTokenFile,
			ProxyURL:        config.ProxyUrl,
		}

		if basicAuthConf := config.BasicAuth; basicAuthConf != nil {
			res.BasicAuth = &models.HTTPBasicAuth{
				Username:     basicAuthConf.Username,
				Password:     basicAuthConf.Password,
				PasswordFile: basicAuthConf.PasswordFile,
			}
		}

		if tlsConfig := config.TlsConfig; tlsConfig != nil {
			res.TLSConfig = &models.TLSConfig{
				CaFile:             tlsConfig.CaFile,
				CertFile:           tlsConfig.CertFile,
				KeyFile:            tlsConfig.KeyFile,
				ServerName:         tlsConfig.ServerName,
				InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
			}
		}
		return res
	}
	return nil
}

func convertModelToHTTPConfig(config *models.HTTPConfig) *iav1beta1.HTTPConfig {
	var res *iav1beta1.HTTPConfig
	if config != nil {
		res = &iav1beta1.HTTPConfig{
			BearerToken:     config.BearerToken,
			BearerTokenFile: config.BearerTokenFile,
			ProxyUrl:        config.ProxyURL,
		}

		if basicAuthConf := config.BasicAuth; basicAuthConf != nil {
			res.BasicAuth = &iav1beta1.BasicAuth{
				Username:     basicAuthConf.Username,
				Password:     basicAuthConf.Password,
				PasswordFile: basicAuthConf.PasswordFile,
			}
		}

		if tlsConfig := config.TLSConfig; tlsConfig != nil {
			res.TlsConfig = &iav1beta1.TLSConfig{
				CaFile:             tlsConfig.CaFile,
				CertFile:           tlsConfig.CertFile,
				KeyFile:            tlsConfig.KeyFile,
				ServerName:         tlsConfig.ServerName,
				InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
			}
		}
		return res
	}
	return nil
}

// Check interfaces.
var (
	_ iav1beta1.ChannelsServer = (*ChannelsService)(nil)
)
