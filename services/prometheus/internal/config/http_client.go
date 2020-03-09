package config

// HTTPClientConfig configures an HTTP client.
type HTTPClientConfig struct {
	// The HTTP basic authentication credentials for the targets.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty"`
	// The bearer token for the targets.
	BearerToken string `yaml:"bearer_token,omitempty"`
	// The bearer token file for the targets.
	BearerTokenFile string `yaml:"bearer_token_file,omitempty"`
	// HTTP proxy server to use to connect to the targets.
	ProxyURL string `yaml:"proxy_url,omitempty"`
	// TLSConfig to use to connect to the targets.
	TLSConfig TLSConfig `yaml:"tls_config,omitempty"`
}

// BasicAuth contains basic HTTP authentication credentials.
type BasicAuth struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password,omitempty"`
	PasswordFile string `yaml:"password_file,omitempty"`
}

// TLSConfig configures the options for TLS connections.
type TLSConfig struct {
	// The CA cert to use for the targets.
	CAFile string `yaml:"ca_file,omitempty"`
	// The client cert file for the targets.
	CertFile string `yaml:"cert_file,omitempty"`
	// The client key file for the targets.
	KeyFile string `yaml:"key_file,omitempty"`
	// Used to verify the hostname for the targets.
	ServerName string `yaml:"server_name,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
}
