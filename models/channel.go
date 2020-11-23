package models

//go:generate reform

// ChannelType represents notificaion channel type.
type ChannelType string

// Available notificaion channel types.
const (
	Email   = ChannelType("email")
	Slack   = ChannelType("slack")
	WebHook = ChannelType("webhook")
)

// Channel represents notification channel configuration.
type Channel struct {
	Id   string      `json:"name"`
	Type ChannelType `json:"type"`

	EmailConfig   *EmailConfig   `json:"email_config"`
	SlackConfig   *SlackConfig   `json:"slack_config"`
	WebHookConfig *WebHookConfig `json:"web_hook_config"`

	Disabled bool `json:"disabled"`
}

// EmailConfig is email notification channel configuration.
type EmailConfig struct {
	SendResolved bool     `json:"send_resolved"`
	To           []string `json:"to"`
}

// SlackConfig is slack notification channel configuration.
type SlackConfig struct {
	SendResolved bool   `json:"send_resolved"`
	Channel      string `json:"channel"`
}

// WebHookConfig is webhook notification channel configuration.
type WebHookConfig struct {
	SendResolved bool        ` json:"send_resolved"`
	Url          string      ` json:"url"`
	HttpConfig   *HTTPConfig ` json:"http_config"`
	MaxAlerts    int32       ` json:"max_alerts"`
}

// HTTPConfig is HTTP connection configuration.
type HTTPConfig struct {
	BasicAuth       *HTTPBasicAuth `json:"basic_auth,omitempty"`
	BearerToken     string         `json:"bearer_token,omitempty"`
	BearerTokenFile string         `json:"bearer_token_file,omitempty"`
	TlsConfig       *TLSConfig     `json:"tls_config,omitempty"`
	ProxyUrl        string         `json:"proxy_url,omitempty"`
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
