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
	"time"

	pmmv1 "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	"github.com/sirupsen/logrus"
)

func fetchMetricsFromDB(ctx context.Context, l *logrus.Entry, timeout time.Duration, db *sql.DB, config Config) ([]*pmmv1.ServerMetric_Metric, error) {
	localCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	tx, err := db.BeginTx(localCtx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	// to minimize risk of modifying DB
	defer tx.Rollback() //nolint:errcheck

	rows, err := db.Query("SELECT " + config.Query) //nolint:gosec,rowserrcheck
	if err != nil {
		return nil, err
	}

	var metrics []*pmmv1.ServerMetric_Metric

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	strs := make([]*string, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = &strs[i]
	}
	cfgColumns := config.mapByColumn()
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			l.Error(err)
			continue
		}

		for idx, column := range columns {
			if _, ok := cfgColumns[column]; ok {
				metrics = append(metrics, &pmmv1.ServerMetric_Metric{
					Key:   column,
					Value: *strs[idx],
				})
			}
		}
	}

	return metrics, nil
}
