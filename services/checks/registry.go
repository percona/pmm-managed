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
	alerts map[string]ammodels.PostableAlert
	nowF   func() time.Time // for tests
}

// newRegistry creates a new registry.
func newRegistry() *registry {
	return &registry{
		alerts: make(map[string]ammodels.PostableAlert),
		nowF:   time.Now,
	}
}

// createAlert creates alert from given AlertParams
func (r *registry) createAlert(labels, annotations map[string]string, alertTTL time.Duration) ammodels.PostableAlert {
	return ammodels.PostableAlert{
		Alert: ammodels.Alert{
			// GeneratorURL: "TODO",
			Labels: labels,
		},
		EndsAt:      strfmt.DateTime(r.nowF().Add(alertTTL)),
		Annotations: annotations,
	}
}

// set clears the previous alerts and sets a new slice of alerts in the registry
func (r *registry) set(alerts []alertWithID) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.alerts = make(map[string]ammodels.PostableAlert)

	for _, alert := range alerts {
		r.alerts[alert.id] = alert.alert
	}
}

// collect returns all firing alerts.
func (r *registry) collect() ammodels.PostableAlerts {
	r.rw.RLock()
	defer r.rw.RUnlock()

	var res ammodels.PostableAlerts
	for _, alert := range r.alerts {
		alert := alert
		res = append(res, &alert)
	}
	return res
}
