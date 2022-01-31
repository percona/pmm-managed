package telemetry_v2

import (
	"context"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type dsQanDbSelect struct {
}

// check interfaces
var (
	_ TelemetryDataSource = (*dsQanDbSelect)(nil)
)

func NewDsQanDbSelect(*logrus.Entry) TelemetryDataSource {
	return &dsQanDbSelect{}
}

func (d *dsQanDbSelect) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*reporter.ServerMetric_Metric, error) {
	//TODO implement me
	return nil, errors.New("NOT IMPLEMENTED")
}
