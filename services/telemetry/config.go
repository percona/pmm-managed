package telemetry

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func LoadConfig(filename string) (*Config, error) {
	content, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func CheckConfig(cfg *Config) (updated bool, err error) {
	if cfg.URL == "" {
		cfg.URL = defaultURL
		updated = true
	}
	if cfg.Interval == 0 {
		cfg.Interval = defaultinterval
		updated = true
	}
	if cfg.UUID == "" {
		cfg.UUID, err = generateUUID()
		if err != nil {
			return false, err
		}
		updated = true
	}
	return updated, nil
}

func SaveConfig(filename string, config *Config) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "cannot marshal call home config")
	}
	err = ioutil.WriteFile(filename, content, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cannot update call home config")
	}
	return nil
}
