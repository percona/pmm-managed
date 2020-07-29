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
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	t.Run("Collect Alerts", func(t *testing.T) {
		r := newRegistry()

		nowValue := time.Now().UTC().Round(0) // strip a monotonic clock reading
		r.nowF = func() time.Time { return nowValue }

		labels := map[string]string{"label": "demo"}
		annotations := map[string]string{"annotation": "test"}
		alertTTL := resolveTimeoutFactor * defaultResendInterval

		expected := &ammodels.PostableAlert{
			Annotations: annotations,
			EndsAt:      strfmt.DateTime(nowValue.Add(alertTTL)),
			Alert: ammodels.Alert{
				Labels: labels,
			},
		}

		alerts := ammodels.PostableAlerts{
			r.createAlert(labels, annotations, alertTTL),
		}
		r.set(alerts)

		collectedAlerts := r.collect()
		require.Len(t, collectedAlerts, 1)
		require.Equal(t, 1, cap(collectedAlerts))
		assert.Equal(t, expected, collectedAlerts[0])
	})
}
