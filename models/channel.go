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

// ChannelType represents notificaion channel type.
type ChannelType string

// Available notificaion channel types.
const (
	Email     = ChannelType("email")
	PagerDuty = ChannelType("pagerduty")
	Slack     = ChannelType("slack")
	WebHook   = ChannelType("webhook")
)

// Channel represents notification channel configuration.
type Channel struct {
	ID   string      `json:"name"`
	Type ChannelType `json:"type"`

	EmailConfig     *EmailConfig     `json:"email_config"`
	PagerDutyConfig *PagerDutyConfig `json:"pager_duty_config"`
	SlackConfig     *SlackConfig     `json:"slack_config"`
	WebHookConfig   *WebHookConfig   `json:"web_hook_config"`

	Disabled bool `json:"disabled"`
}

// EmailConfig is email notification channel configuration.
type EmailConfig struct {
	SendResolved bool     `json:"send_resolved"`
	To           []string `json:"to"`
}

// PagerDutyConfig represents PagerDuty channel configuration.
type PagerDutyConfig struct {
	SendResolved bool   `json:"send_resolved"`
	RoutingKey   string `json:"routing_key"`
	ServiceKey   string `json:"service_key"`
}

// SlackConfig is slack notification channel configuration.
type SlackConfig struct {
	SendResolved bool   `json:"send_resolved"`
	Channel      string `json:"channel"`
}

// WebHookConfig is webhook notification channel configuration.
type WebHookConfig struct {
	SendResolved bool        ` json:"send_resolved"`
	URL          string      ` json:"url"`
	HTTPConfig   *HTTPConfig ` json:"http_config"`
	MaxAlerts    int32       ` json:"max_alerts"`
}

// HTTPConfig is HTTP connection configuration.
type HTTPConfig struct {
	BasicAuth       *HTTPBasicAuth `json:"basic_auth,omitempty"`
	BearerToken     string         `json:"bearer_token,omitempty"`
	BearerTokenFile string         `json:"bearer_token_file,omitempty"`
	TLSConfig       *TLSConfig     `json:"tls_config,omitempty"`
	ProxyURL        string         `json:"proxy_url,omitempty"`
}

// HTTPBasicAuth is HTTP basic authentication configuration.
type HTTPBasicAuth struct {
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"password_file,omitempty"`
}

// TLSConfig is TLS configuration.
type TLSConfig struct {
	CaFile             string `json:"ca_file,omitempty"`
	CertFile           string `json:"cert_file,omitempty"`
	KeyFile            string `json:"key_file,omitempty"`
	ServerName         string `json:"server_name"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}
