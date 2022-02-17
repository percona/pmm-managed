package vmalert

type Config struct {
	Enabled bool `yaml:"enabled"`
}

func (c *Config) Init() {
}
