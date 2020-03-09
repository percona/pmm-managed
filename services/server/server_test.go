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

package server

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/percona/pmm/api/serverpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestServer(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()

	newServer := func() *Server {
		s, err := NewServer(reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf)),
			nil, nil, nil, nil, "")
		require.NoError(t, err)
		return s
	}

	t.Run("UpdateSettingsFromEnv", func(t *testing.T) {
		t.Run("Typical", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"DISABLE_UPDATES=true",
				"DISABLE_TELEMETRY=1",
				"METRICS_RESOLUTION_HR=1s",
				"METRICS_RESOLUTION_MR=2s",
				"METRICS_RESOLUTION_LR=3s",
				"DATA_RETENTION=240h",
			})
			require.NoError(t, err)
			assert.Equal(t, true, s.envDisableUpdates)
			assert.Equal(t, true, s.envDisableTelemetry)
			assert.Equal(t, time.Second, s.envMetricsResolutionHR)
			assert.Equal(t, 2*time.Second, s.envMetricsResolutionMR)
			assert.Equal(t, 3*time.Second, s.envMetricsResolutionLR)
			assert.Equal(t, 10*24*time.Hour, s.envDataRetention)
		})

		t.Run("Untypical", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"DISABLE_TELEMETRY=TrUe",
				"METRICS_RESOLUTION=3S",
				"DATA_RETENTION=360H",
			})
			require.NoError(t, err)
			assert.Equal(t, true, s.envDisableTelemetry)
			assert.Equal(t, 3*time.Second, s.envMetricsResolutionHR)
			assert.Equal(t, 15*24*time.Hour, s.envDataRetention)
		})

		t.Run("NoValue", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"DISABLE_TELEMETRY",
			})
			require.Errorf(t, err, "validation error of environment variables")
		})

		t.Run("InvalidValue", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"DISABLE_TELEMETRY=",
			})
			require.Errorf(t, err, "validation error of environment variables")
			assert.Equal(t, false, s.envDisableTelemetry)
		})

		t.Run("MetricsLessThenMin", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"METRICS_RESOLUTION=5ns",
			})
			require.Errorf(t, err, "validation error of environment variables")
		})

		t.Run("DataRetentionLessThenMin", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"DATA_RETENTION=12h",
			})
			require.Errorf(t, err, "validation error of environment variables")
		})

		t.Run("Data retention is not a natural number of days", func(t *testing.T) {
			s := newServer()
			err := s.UpdateSettingsFromEnv([]string{
				"DATA_RETENTION=30h",
			})
			require.Errorf(t, err, "environment variable %q should be a multiple of 24h0m0s")
		})
	})

	t.Run("ValidateChangeSettingsRequest", func(t *testing.T) {
		s := newServer()

		ctx := context.TODO()

		tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "Both alert_manager_rules and remove_alert_manager_rules are present."),
			s.validateChangeSettingsRequest(ctx, &serverpb.ChangeSettingsRequest{
				AlertManagerRules:       "something",
				RemoveAlertManagerRules: true,
			}))
	})

	t.Run("ValidateAlertManagerRules", func(t *testing.T) {
		s := newServer()

		t.Run("Valid", func(t *testing.T) {
			rules := strings.TrimSpace(`
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
			`) + "\n"
			err := s.validateAlertManagerRules(context.Background(), rules)
			assert.NoError(t, err)
		})

		t.Run("Zero", func(t *testing.T) {
			rules := strings.TrimSpace(`
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
			`) + "\n"
			err := s.validateAlertManagerRules(context.Background(), rules)
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "Zero Alert Manager rules found."), err)
		})

		t.Run("Invalid", func(t *testing.T) {
			rules := strings.TrimSpace(`
groups:
- name: example
  rules:
  - alert: HighRequestLatency
			`) + "\n"
			err := s.validateAlertManagerRules(context.Background(), rules)
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "Invalid Alert Manager rules."), err)
		})
	})
}
