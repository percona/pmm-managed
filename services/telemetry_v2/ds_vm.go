package telemetry_v2

import (
	"context"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type dsVm struct {
}

// check interfaces
var (
	_ TelemetryDataSource = (*dsVm)(nil)
)

func NewDsVm(*logrus.Entry) TelemetryDataSource {
	return &dsVm{}
}

func (d *dsVm) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*reporter.ServerMetric_Metric, error) {
	//TODO implement me
	return nil, errors.New("NOT IMPLEMENTED")
}
