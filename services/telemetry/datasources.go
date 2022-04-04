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

// Package telemetry provides telemetry functionality.
package telemetry

import (
	"context"

	pmmv1 "github.com/percona-platform/saas/gen/telemetry/events/pmm"
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

	vmDB, err := NewDataSourceVictoriaMetrics(*config.DataSources.VM, l)
	if err != nil {
		return nil, err
	}

	return &telemetryDataSourceRegistry{
		l: l,
		dataSources: map[TelemetryDataSourceName]TelemetryDataSource{
			DS_VM:           vmDB,
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
	FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*pmmv1.ServerMetric_Metric, error)
	Enabled() bool
}
