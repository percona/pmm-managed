package config

import "github.com/prometheus/common/model"

// ServiceDiscoveryConfig configures lists of different service discovery mechanisms.
type ServiceDiscoveryConfig struct {
	// List of labeled target groups for this job.
	StaticConfigs []*Group `yaml:"static_configs,omitempty"`
	// List of file service discovery configurations.
	FileSDConfigs []*FilesSDConfig `yaml:"file_sd_configs,omitempty"`
	// List of Kubernetes service discovery configurations.
	KubernetesSDConfigs []*KubernetesSDConfig `yaml:"kubernetes_sd_configs,omitempty"`
}

// Group is a set of targets with a common label set(production , test, staging etc.).
type Group struct {
	// Targets is a list of targets identified by a label set. Each target is
	// uniquely identifiable in the group by its address label.
	Targets []string `yaml:"targets,omitempty"`
	// Labels is a set of labels that is common across all targets in the group.
	Labels map[string]string `yaml:"labels,omitempty"`
}

// FilesSDConfig is the configuration for file based discovery.
type FilesSDConfig struct {
	Files           []string       `yaml:"files"`
	RefreshInterval model.Duration `yaml:"refresh_interval,omitempty"`
}

// KubernetesSDConfig is the configuration for Kubernetes service discovery.
type KubernetesSDConfig struct {
	APIServer          string           `yaml:"api_server,omitempty"`
	Role               string           `yaml:"role"`
	HTTPClientConfig   HTTPClientConfig `yaml:",inline"`
	NamespaceDiscovery []string         `yaml:"namespaces,omitempty"`
}
