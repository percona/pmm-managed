package config

import "github.com/prometheus/common/model"

// RemoteWriteConfig is the configuration for writing to remote storage.
type RemoteWriteConfig struct {
	URL                 string           `yaml:"url"`
	RemoteTimeout       model.Duration   `yaml:"remote_timeout,omitempty"`
	WriteRelabelConfigs []*RelabelConfig `yaml:"write_relabel_configs,omitempty"`
	Name                string           `yaml:"name,omitempty"`

	// We cannot do proper Go type embedding below as the parser will then parse
	// values arbitrarily into the overflow maps of further-down types.
	HTTPClientConfig HTTPClientConfig `yaml:",inline"`
	QueueConfig      QueueConfig      `yaml:"queue_config,omitempty"`
}

// QueueConfig is the configuration for the queue used to write to remote
// storage.
type QueueConfig struct {
	// Number of samples to buffer per shard before we block. Defaults to
	// MaxSamplesPerSend.
	Capacity int `yaml:"capacity,omitempty"`

	// Max number of shards, i.e. amount of concurrency.
	MaxShards int `yaml:"max_shards,omitempty"`

	// Min number of shards, i.e. amount of concurrency.
	MinShards int `yaml:"min_shards,omitempty"`

	// Maximum number of samples per send.
	MaxSamplesPerSend int `yaml:"max_samples_per_send,omitempty"`

	// Maximum time sample will wait in buffer.
	BatchSendDeadline model.Duration `yaml:"batch_send_deadline,omitempty"`

	// On recoverable errors, backoff exponentially.
	MinBackoff model.Duration `yaml:"min_backoff,omitempty"`
	MaxBackoff model.Duration `yaml:"max_backoff,omitempty"`
}
