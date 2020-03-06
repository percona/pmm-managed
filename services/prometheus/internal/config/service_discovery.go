package config

// ServiceDiscoveryConfig configures lists of different service discovery mechanisms.
type ServiceDiscoveryConfig struct {
	// List of labeled target groups for this job.
	StaticConfigs []*Group `yaml:"static_configs,omitempty"`

	//TODO
	// code removed
}

// Group is a set of targets with a common label set(production , test, staging etc.).
type Group struct {
	// Targets is a list of targets identified by a label set. Each target is
	// uniquely identifiable in the group by its address label.
	Targets []string `yaml:"targets,omitempty"`
	// Labels is a set of labels that is common across all targets in the group.
	Labels map[string]string  `yaml:"labels,omitempty"`
}
