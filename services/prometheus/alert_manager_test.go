package prometheus

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/utils/tests"
)

func TestAlertManager(t *testing.T) {
	t.Run("ValidateRules", func(t *testing.T) {
		s := NewAlertManager("")

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
			err := s.ValidateRules(context.Background(), rules)
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
			err := s.ValidateRules(context.Background(), rules)
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "Zero Alert Manager rules found."), err)
		})

		t.Run("Invalid", func(t *testing.T) {
			rules := strings.TrimSpace(`
groups:
- name: example
  rules:
  - alert: HighRequestLatency
			`) + "\n"
			err := s.ValidateRules(context.Background(), rules)
			tests.AssertGRPCError(t, status.New(codes.InvalidArgument, "Invalid Alert Manager rules."), err)
		})
	})
}
