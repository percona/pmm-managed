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
	t.Run("DelayFor", func(t *testing.T) {
		r := newRegistry()

		nowValue := time.Now()
		r.nowF = func() time.Time { return nowValue }

		id := "1234567890"
		labels := map[string]string{"label": "demo"}
		annotations := map[string]string{"annotation": "test"}

		r.createAlert(id, labels, annotations)
		assert.Empty(t, r.collect())

		// 1 second after
		nowValue = nowValue.Add(time.Second)

		expected := &ammodels.PostableAlert{
			Annotations: annotations,
			EndsAt:      strfmt.DateTime(nowValue.Add(resolveTimeoutFactor * r.resendInterval)),
			Alert: ammodels.Alert{
				Labels: labels,
			},
		}

		alerts := r.collect()
		require.Len(t, alerts, 1)
		assert.Equal(t, expected, alerts[0])
	})
}
