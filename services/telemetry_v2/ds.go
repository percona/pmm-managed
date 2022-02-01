package telemetry_v2

import (
	"context"
	"database/sql"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
	"github.com/sirupsen/logrus"
	"time"
)

func fetchMetricsFromDB(l *logrus.Entry, timeout time.Duration, db *sql.DB, ctx context.Context, config TelemetryConfig) ([]*reporter.ServerMetric_Metric, error) {
	localCtx, _ := context.WithTimeout(ctx, timeout)
	tx, err := db.BeginTx(localCtx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	// to minimize risk of modifying DB
	defer tx.Rollback()

	rows, err := db.Query("SELECT " + config.Query)
	if err != nil {
		return nil, err
	}

	var metrics []*reporter.ServerMetric_Metric

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	strs := make([]*string, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = &strs[i]
	}
	cfgColumns := config.MapByColumn()
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			l.Error(err)
			continue
		}

		for idx, column := range columns {
			if _, ok := cfgColumns[column]; ok {
				metrics = append(metrics, &reporter.ServerMetric_Metric{
					Key:   column,
					Value: *strs[idx],
				})
			}
		}
	}

	return metrics, nil
}
