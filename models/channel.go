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
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

//go:generate reform

// ChannelType represents notificaion channel type.
type ChannelType string

// Available notification channel types.
const (
	Email     = ChannelType("email")
	PagerDuty = ChannelType("pagerduty")
	Slack     = ChannelType("slack")
	WebHook   = ChannelType("webhook")
)

// Channel represents notification channel configuration.
//reform:notification_channels
type Channel struct {
	ID   string      `reform:"id,pk"`
	Type ChannelType `reform:"type"`

	EmailConfig     *EmailConfig     `reform:"email_config"`
	PagerDutyConfig *PagerDutyConfig `reform:"pagerduty_config"`
	SlackConfig     *SlackConfig     `reform:"slack_config"`
	WebHookConfig   *WebHookConfig   `reform:"webhook_config"`

	Disabled bool `reform:"disabled"`
}

// EmailConfig is email notification channel configuration.
type EmailConfig struct {
	SendResolved bool     `json:"send_resolved"`
	To           []string `json:"to"`
}

// Value implements database/sql/driver Valuer interface.
func (c *EmailConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Scan implements database/sql Scanner interface.
func (c *EmailConfig) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("EmailConfi.Scan: expected []byte or string, got %T (%q)", src, src)
	}

	return json.Unmarshal(b, c)
}

// PagerDutyConfig represents PagerDuty channel configuration.
type PagerDutyConfig struct {
	SendResolved bool   `json:"send_resolved"`
	RoutingKey   string `json:"routing_key"`
	ServiceKey   string `json:"service_key"`
}

// Value implements database/sql/driver Valuer interface.
func (c *PagerDutyConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Scan implements database/sql Scanner interface.
func (c *PagerDutyConfig) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("PagerDutyConfig.Scan: expected []byte or string, got %T (%q)", src, src)
	}

	return json.Unmarshal(b, c)
}

// SlackConfig is slack notification channel configuration.
type SlackConfig struct {
	SendResolved bool   `json:"send_resolved"`
	Channel      string `json:"channel"`
}

// Value implements database/sql/driver Valuer interface.
func (c *SlackConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Scan implements database/sql Scanner interface.
func (c *SlackConfig) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("SlackConfig.Scan: expected []byte or string, got %T (%q)", src, src)
	}

	return json.Unmarshal(b, c)
}

// WebHookConfig is webhook notification channel configuration.
type WebHookConfig struct {
	SendResolved bool        ` json:"send_resolved"`
	URL          string      ` json:"url"`
	HTTPConfig   *HTTPConfig ` json:"http_config"`
	MaxAlerts    int32       ` json:"max_alerts"`
}

// Value implements database/sql/driver Valuer interface.
func (c *WebHookConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Scan implements database/sql Scanner interface.
func (c *WebHookConfig) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("WebHookConfig.Scan: expected []byte or string, got %T (%q)", src, src)
	}

	return json.Unmarshal(b, c)
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
