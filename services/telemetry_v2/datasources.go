package telemetry_v2

import (
	"context"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type TelemetryDataSourceName string

const (
	DS_VM           = TelemetryDataSourceName("VM")
	DS_PMMDB_SELECT = TelemetryDataSourceName("PMMDB_SELECT")
	DS_QANDB_SELECT = TelemetryDataSourceName("QANDB_SELECT")
)

type TelemetryDataSourceLocator interface {
	LocateTelemetryDataSource(name string) (TelemetryDataSource, error)
}

type telemetryDataSourceRegistry struct {
	l           *logrus.Entry
	dataSources map[TelemetryDataSourceName]TelemetryDataSource
}

func NewDataSourceRegistry(config ServiceConfig, l *logrus.Entry) (TelemetryDataSourceLocator, error) {
	pmmDB, err := NewDsPmmDbSelect(*config.DataSources.PMMDB_SELECT, l)
	if err != nil {
		return nil, err
	}

	qanDB, err := NewDsQanDbSelect(*config.DataSources.QANDB_SELECT, l)
	if err != nil {
		return nil, err
	}

	return &telemetryDataSourceRegistry{
		l: l,
		dataSources: map[TelemetryDataSourceName]TelemetryDataSource{
			DS_VM:           NewDsVm(l),
			DS_PMMDB_SELECT: pmmDB,
			DS_QANDB_SELECT: qanDB,
		},
	}, nil
}

func (r *telemetryDataSourceRegistry) LocateTelemetryDataSource(name string) (TelemetryDataSource, error) {
	ds, ok := r.dataSources[TelemetryDataSourceName(name)]
	if !ok {
		return nil, errors.Errorf("datasource [%s] is not supported", name)
	}
	return ds, nil
}

type TelemetryDataSource interface {
	FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*reporter.ServerMetric_Metric, error)
}
