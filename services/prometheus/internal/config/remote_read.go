package config

import "github.com/prometheus/common/model"

// RemoteReadConfig is the configuration for reading from remote storage.
type RemoteReadConfig struct {
	URL           string         `yaml:"url"`
	RemoteTimeout model.Duration `yaml:"remote_timeout,omitempty"`
	ReadRecent    bool           `yaml:"read_recent,omitempty"`
	Name          string         `yaml:"name,omitempty"`

	// We cannot do proper Go type embedding below as the parser will then parse
	// values arbitrarily into the overflow maps of further-down types.
	HTTPClientConfig HTTPClientConfig `yaml:",inline"`

	// RequiredMatchers is an optional list of equality matchers which have to
	// be present in a selector to query the remote read endpoint.
	RequiredMatchers map[string]string `yaml:"required_matchers,omitempty"`
}
