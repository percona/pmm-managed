package telemetry_v2

import (
	"context"
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
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

func NewDsQanDbSelect(config DSConfigQAN, l *logrus.Entry) (TelemetryDataSource, error) {
	db, err := openQANDBConnection(config.DSN)
	if err != nil {
		return nil, err
	}
	return &dsQanDbSelect{
		l:      l,
		config: config,
		db:     db,
	}, nil
}

func openQANDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open connection to QAN DB")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Failed to ping QAN DB")
	}
	return db, nil
}

func (d *dsQanDbSelect) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*reporter.ServerMetric_Metric, error) {
	return fetchMetricsFromDB(d.l, d.config.Timeout, d.db, ctx, config)
}
