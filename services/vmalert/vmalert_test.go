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

package vmalert

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (context.Context, *ExternalAlertingRules, *Service) {
	t.Helper()

	rules := NewExternalAlertingRules()
	err := rules.RemoveRulesFile()
	require.NoError(t, err)

	svc, err := NewVMAlert(rules, External)
	require.NoError(t, err)
	err = svc.IsReady(context.Background())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		cancel()
		err = rules.RemoveRulesFile()
		require.NoError(t, err)
	})

	return ctx, rules, svc
}

func TestVMAlert(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		ctx, _, svc := setup(t)
		err := svc.updateConfiguration(ctx)
		require.NoError(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		ctx, rules, svc := setup(t)
		err := rules.WriteRules(strings.TrimSpace(`
groups:
  - name: example
    rules:
    - alert: HighRequestLatency
      expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
      for: 10m
      labels:
          severity: page
      annotations:
          summary: High request latency
		`))
		require.NoError(t, err)
		err = svc.updateConfiguration(ctx)
		require.NoError(t, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		ctx, rules, svc := setup(t)
		err := rules.WriteRules(`foobar`)
		require.NoError(t, err)
		err = svc.updateConfiguration(ctx)
		require.NoError(t, err)
	})
}
