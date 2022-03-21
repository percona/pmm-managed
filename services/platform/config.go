package platform

type Config struct {
	SkipTlsVerification bool `yaml:"skip_tls_verification"`
}

func (c *Config) Init() {
}
