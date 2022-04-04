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
	"database/sql"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	pmmv1 "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type dsQanDbSelect struct {
	l      *logrus.Entry
	config DSConfigQAN
	db     *sql.DB
}

// check interfaces
var (
	_ TelemetryDataSource = (*dsQanDbSelect)(nil)
)

func (d *dsQanDbSelect) Enabled() bool {
	return d.config.Enabled
}

func NewDsQanDbSelect(config DSConfigQAN, l *logrus.Entry) (TelemetryDataSource, error) {
	db, err := openQANDBConnection(config.DSN, config.Enabled)
	if err != nil {
		return nil, err
	}
	return &dsQanDbSelect{
		l:      l,
		config: config,
		db:     db,
	}, nil
}

func openQANDBConnection(dsn string, enabled bool) (*sql.DB, error) {
	if !enabled {
		return nil, nil
	}

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open connection to QAN DB")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Failed to ping QAN DB")
	}
	return db, nil
}

func (d *dsQanDbSelect) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*pmmv1.ServerMetric_Metric, error) {
	return fetchMetricsFromDB(d.l, d.config.Timeout, d.db, ctx, config)
}
