package telemetry_v2

import (
	"context"
	pmmv1 "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"time"
)

type dsVm struct {
	l      *logrus.Entry
	config DSVM
	vm     v1.API
}

// check interfaces
var (
	_ TelemetryDataSource = (*dsVm)(nil)
)

func NewDsVm(config DSVM, l *logrus.Entry) (TelemetryDataSource, error) {
	client, err := api.NewClient(api.Config{
		Address: config.Address,
	})

	if err != nil {
		return nil, err
	}

	return &dsVm{
		l:      l,
		config: config,
		vm:     v1.NewAPI(client),
	}, nil
}

func (d *dsVm) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*pmmv1.ServerMetric_Metric, error) {

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
						Key:   configItem.Label,
						Value: string(value),
					})
				}
			}
			//TODO: verify if impl is correct
			if configItem.Value != "" {
				metrics = append(metrics, &pmmv1.ServerMetric_Metric{
					Key:   configItem.MetricName,
					Value: configItem.Value,
				})
			}
		}
	}

	return metrics, nil
}
