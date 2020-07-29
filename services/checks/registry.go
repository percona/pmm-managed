// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package checks

import (
	"sync"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/percona/pmm/api/alertmanager/ammodels"
)

// registry stores alerts and delay information by IDs.
type registry struct {
	rw     sync.RWMutex
	alerts ammodels.PostableAlerts
	nowF   func() time.Time // for tests
}

// newRegistry creates a new registry.
func newRegistry() *registry {
	return &registry{
		nowF: time.Now,
	}
}

// createAlert creates alert from given AlertParams.
func (r *registry) createAlert(labels, annotations map[string]string, alertTTL time.Duration) *ammodels.PostableAlert {
	return &ammodels.PostableAlert{
		Alert: ammodels.Alert{
			// GeneratorURL: "TODO",
			Labels: labels,
		},
		EndsAt:      strfmt.DateTime(r.nowF().Add(alertTTL)),
		Annotations: annotations,
	}
}

// set replaces stored alerts with a copy of given ones.
func (r *registry) set(alerts ammodels.PostableAlerts) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.alerts = make(ammodels.PostableAlerts, len(alerts))
	copy(r.alerts, alerts)
}

// collect returns a copy of stored alerts.
func (r *registry) collect() ammodels.PostableAlerts {
	r.rw.RLock()
	defer r.rw.RUnlock()

	alerts := make(ammodels.PostableAlerts, len(r.alerts))
	copy(alerts, r.alerts)

	return alerts
}
