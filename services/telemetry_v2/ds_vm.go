package telemetry_v2

import (
	"context"
	"time"

	pmmv1 "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type dataSourceVictoriaMetrics struct {
	l      *logrus.Entry
	config DataSourceVictoriaMetrics
	vm     v1.API
}

// check interfaces
var (
	_ TelemetryDataSource = (*dataSourceVictoriaMetrics)(nil)
)

func (d *dataSourceVictoriaMetrics) Enabled() bool {
	return d.config.Enabled
}

func NewDataSourceVictoriaMetrics(config DataSourceVictoriaMetrics, l *logrus.Entry) (TelemetryDataSource, error) {
	client, err := api.NewClient(api.Config{
		Address: config.Address,
	})

	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return &dataSourceVictoriaMetrics{
			l:      l,
			config: config,
			vm:     nil,
		}, nil
	}

	return &dataSourceVictoriaMetrics{
		l:      l,
		config: config,
		vm:     v1.NewAPI(client),
	}, nil
}

func (d *dataSourceVictoriaMetrics) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*pmmv1.ServerMetric_Metric, error) {

	localCtx, _ := context.WithTimeout(ctx, d.config.Timeout)
	result, _, err := d.vm.Query(localCtx, config.Query, time.Now())
	if err != nil {
		return nil, err
	}

	var metrics []*pmmv1.ServerMetric_Metric

	for _, v := range result.(model.Vector) {
		for _, configItem := range config.Data {
			if configItem.Label != "" {
				value, ok := v.Metric[model.LabelName(configItem.Label)]
				if ok {
					metrics = append(metrics, &pmmv1.ServerMetric_Metric{
						Key:   configItem.MetricName,
						Value: string(value),
					})
				}
			}

			if configItem.Value != "" {
				metrics = append(metrics, &pmmv1.ServerMetric_Metric{
					Key:   configItem.MetricName,
					Value: v.Value.String(),
				})
			}
		}
	}

	return metrics, nil
}
