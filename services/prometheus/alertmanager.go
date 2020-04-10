package prometheus

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
)

// Alertmanager is responsible for interactions with Prometheus.
type Alertmanager struct {
}

// NewAlertmanager creates new service.
func NewAlertmanager() *Alertmanager {
	return &Alertmanager{}
}

// Run runs Alertmanager configuration update loop until ctx is canceled.
func (am *Alertmanager) Run(ctx context.Context) {
	// TODO that's a temporary measure until we start generating /etc/alertmanager.yml
	// using /srv/alertmanager/alertmanager.base.yml as a base
	const path = "/srv/alertmanager/alertmanager.base.yml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultBase := strings.TrimSpace(`
---
# You can edit this file; changes will be preserved.

route:
  receiver: empty
  routes: []

receivers:
  - name: empty
`) + "\n"
		_ = ioutil.WriteFile(path, []byte(defaultBase), 0644)
	}

	<-ctx.Done()
}
