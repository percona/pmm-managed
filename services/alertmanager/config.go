package alertmanager

import (
	"os"
	"time"
)

const (
	updateBatchDelay           = time.Second
	configurationUpdateTimeout = 3 * time.Second

	alertmanagerDir     = "/srv/alertmanager"
	alertmanagerCertDir = "/srv/alertmanager/cert"
	alertmanagerDataDir = "/srv/alertmanager/data"
	dirPerm             = os.FileMode(0o775)

	alertmanagerConfigPath     = "/etc/alertmanager.yml"
	alertmanagerBaseConfigPath = "/srv/alertmanager/alertmanager.base.yml"

	receiverNameSeparator = " + "
)

var notificationLabels = []string{"node_name", "node_id", "service_name", "service_id", "service_type", "rule_id",
	"alertgroup", "template_name", "severity", "agent_id", "agent_type", "job"}

type Config struct {
	Enabled bool `yaml:"enabled"`
}

func (c *Config) Init() {
}
