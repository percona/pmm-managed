package server

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/api"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestSTTMetrics(t *testing.T) {
	if !pmmapitests.RunSTTTests {
		t.Skip("Skipping STT tests until we have environment: https://jira.percona.com/browse/PMM-5106")
	}

	t.Run("StartSTTChecksAndRecordMetrics", func(t *testing.T) {
		client, err := api.NewClient(api.Config{
			Address: pmmapitests.BaseURL.ResolveReference(&url.URL{
				Path: "/prometheus",
			}).String(),
		})
		require.NoError(t, err)
		promClient := promapi.NewAPI(client)

		testCases := []struct {
			query          string
			metricType     string
			expectedValues []string
		}{
			{
				query:      "pmm_managed_checks_alerts_generated_total",
				metricType: "vector",
				expectedValues: []string{
					`pmm_managed_checks_alerts_generated_total{check_type="MONGODB_BUILDINFO", instance="pmm-server", job="pmm-managed", service_type="mongodb"} => 0`,
					`pmm_managed_checks_alerts_generated_total{check_type="MONGODB_GETCMDLINEOPTS", instance="pmm-server", job="pmm-managed", service_type="mongodb"} => 0`,
					`pmm_managed_checks_alerts_generated_total{check_type="MONGODB_GETPARAMETER", instance="pmm-server", job="pmm-managed", service_type="mongodb"} => 0`,
					`pmm_managed_checks_alerts_generated_total{check_type="MYSQL_SELECT", instance="pmm-server", job="pmm-managed", service_type="mysql"} => 0`,
					`pmm_managed_checks_alerts_generated_total{check_type="MYSQL_SHOW", instance="pmm-server", job="pmm-managed", service_type="mysql"} => 0`,
					`pmm_managed_checks_alerts_generated_total{check_type="POSTGRESQL_SELECT", instance="pmm-server", job="pmm-managed", service_type="postgresql"} => 0`,
					`pmm_managed_checks_alerts_generated_total{check_type="POSTGRESQL_SHOW", instance="pmm-server", job="pmm-managed", service_type="postgresql"} => 0`},
			},
			{
				query:      "pmm_managed_checks_scripts_executed_total",
				metricType: "vector",
				expectedValues: []string{
					`pmm_managed_checks_scripts_executed_total{instance="pmm-server", job="pmm-managed", service_type="mongodb"} => 0`,
					`pmm_managed_checks_scripts_executed_total{instance="pmm-server", job="pmm-managed", service_type="mysql"} => 0`,
					`pmm_managed_checks_scripts_executed_total{instance="pmm-server", job="pmm-managed", service_type="postgresql"} => 0`},
			},
		}

		for _, tc := range testCases {
			result, _, err := promClient.Query(context.Background(),
				tc.query, time.Now())

			var actualValues []string
			for _, s := range strings.Split(result.String(), "\n") {
				// remove the timestamp from the values
				metric := strings.Split(s, " @")
				actualValues = append(actualValues, metric[0])
			}

			require.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Len(t, result, len(tc.expectedValues))
			assert.Equal(t, tc.metricType, result.Type().String())
			assert.Equal(t, tc.expectedValues, actualValues)
		}
	})
}
