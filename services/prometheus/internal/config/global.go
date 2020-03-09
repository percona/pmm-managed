package config

import "github.com/prometheus/common/model"

// GlobalConfig configures values that are used across other configuration
// objects.
type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval model.Duration `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout model.Duration `yaml:"scrape_timeout,omitempty"`
	// How frequently to evaluate rules by default.
	EvaluationInterval model.Duration `yaml:"evaluation_interval,omitempty"`
	// File to which PromQL queries are logged.
	QueryLogFile string `yaml:"query_log_file,omitempty"`
	// The labels to add to any timeseries that this Prometheus instance scrapes.
	ExternalLabels map[string]string `yaml:"external_labels,omitempty"`
}
