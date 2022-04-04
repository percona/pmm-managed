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

// DataSourceName data source name.
type DataSourceName string

// DataSourceLocator locates data source by name.
type DataSourceLocator interface {
	LocateTelemetryDataSource(name string) (DataSource, error)
}

type dataSourceRegistry struct {
	l           *logrus.Entry
	dataSources map[DataSourceName]DataSource
}

// NewDataSourceRegistry makes new data source registry
func NewDataSourceRegistry(config ServiceConfig, l *logrus.Entry) (DataSourceLocator, error) {
	pmmDB, err := NewDsPmmDbSelect(*config.DataSources.PmmDBSelect, l)
	if err != nil {
		return nil, err
	}

	qanDB, err := NewDsQanDbSelect(*config.DataSources.QanDBSelect, l)
	if err != nil {
		return nil, err
	}

	vmDB, err := NewDataSourceVictoriaMetrics(*config.DataSources.VM, l)
	if err != nil {
		return nil, err
	}

	return &dataSourceRegistry{
		l: l,
		dataSources: map[DataSourceName]DataSource{
			"VM":           vmDB,
			"PMMDB_SELECT": pmmDB,
			"QANDB_SELECT": qanDB,
		},
	}, nil
}

func (r *dataSourceRegistry) LocateTelemetryDataSource(name string) (DataSource, error) {
	ds, ok := r.dataSources[DataSourceName(name)]
	if !ok {
		return nil, errors.Errorf("data source [%s] is not supported", name)
	}
	return ds, nil
}

// DataSource telemetry data source.
type DataSource interface {
	FetchMetrics(ctx context.Context, config Config) ([]*pmmv1.ServerMetric_Metric, error)
	Enabled() bool
}
