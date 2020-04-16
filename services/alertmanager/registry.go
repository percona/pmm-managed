package alertmanager

import (
	"strings"
	"sync"
	"time"

	"github.com/percona/pmm/api/alertmanager/ammodels"
)

// registry stores alerts and delay information by IDs.
type registry struct {
	rw     sync.RWMutex
	alerts map[string]*ammodels.PostableAlert
	times  map[string]time.Time
}

// newRegistry creates a new registry.
func newRegistry() *registry {
	return &registry{
		alerts: make(map[string]*ammodels.PostableAlert),
		times:  make(map[string]time.Time),
	}
}

// Add adds or replaces alert with given ID. If that ID wasn't present before,
// alert is added in the pending state. It we be transitioned to the firing state after delayFor interval.
// This is similar to `for` field of Prometheus alerting rule:
// https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/
func (r *registry) Add(id string, delayFor time.Duration, alert *ammodels.PostableAlert) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.alerts[id] = alert
	if r.times[id].IsZero() {
		r.times[id] = time.Now().Add(delayFor)
	}
}

// Remove removes alert by ID.
func (r *registry) Remove(id string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	delete(r.alerts, id)
	delete(r.times, id)
}

// RemovePrefix removes all alerts with given ID prefix.
func (r *registry) RemovePrefix(prefix string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id := range r.alerts {
		if strings.HasPrefix(id, prefix) {
			delete(r.alerts, id)
			delete(r.times, id)
		}
	}
}

// Collect returns all firing alerts.
func (r *registry) Collect() ammodels.PostableAlerts {
	r.rw.RLock()
	defer r.rw.RUnlock()

	var res ammodels.PostableAlerts
	now := time.Now()
	for id, t := range r.times {
		if t.After(now) {
			res = append(res, r.alerts[id])
		}
	}
	return res
}
