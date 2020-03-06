package config

import "github.com/prometheus/common/model"

const (
	// AlertmanagerAPIVersionV1 represents
	// github.com/prometheus/alertmanager/api/v1.
	AlertmanagerAPIVersionV1 = "v1"
	// AlertmanagerAPIVersionV2 represents
	// github.com/prometheus/alertmanager/api/v2.
	AlertmanagerAPIVersionV2 = "v2"
)

type AlertingConfig struct {
	AlertRelabelConfigs []*RelabelConfig      `yaml:"alert_relabel_configs,omitempty"`
	AlertmanagerConfigs []*AlertmanagerConfig `yaml:"alertmanagers,omitempty"`
}

// AlertmanagerConfig configures how Alertmanagers can be discovered and communicated with.
type AlertmanagerConfig struct {
	// We cannot do proper Go type embedding below as the parser will then parse
	// values arbitrarily into the overflow maps of further-down types.

	ServiceDiscoveryConfig ServiceDiscoveryConfig `yaml:",inline"`
	HTTPClientConfig       HTTPClientConfig     `yaml:",inline"`

	// The URL scheme to use when talking to Alertmanagers.
	Scheme string `yaml:"scheme,omitempty"`
	// Path prefix to add in front of the push endpoint path.
	PathPrefix string `yaml:"path_prefix,omitempty"`
	// The timeout used when sending alerts.
	Timeout model.Duration `yaml:"timeout,omitempty"`

	// The api version of Alertmanager.
	APIVersion string `yaml:"api_version"`

	// List of Alertmanager relabel configurations.
	RelabelConfigs []*RelabelConfig `yaml:"relabel_configs,omitempty"`
}