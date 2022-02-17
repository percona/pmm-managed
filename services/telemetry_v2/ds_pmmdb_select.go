package telemetry_v2

import (
	"context"
	"database/sql"
	"fmt"
	pmmv1 "github.com/percona-platform/saas/gen/telemetry/events/pmm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/url"
)

type dsPmmDbSelect struct {
	l      *logrus.Entry
	config DSConfigPMMDB
	db     *sql.DB
}

// check interfaces
var (
	_ TelemetryDataSource = (*dsPmmDbSelect)(nil)
)

func NewDsPmmDbSelect(config DSConfigPMMDB, l *logrus.Entry) (TelemetryDataSource, error) {
	db, err := openPMMDBConnection(config)
	if err != nil {
		return nil, err
	}

	return &dsPmmDbSelect{
		l:      l,
		config: config,
		db:     db,
	}, nil
}

func openPMMDBConnection(config DSConfigPMMDB) (*sql.DB, error) {
	var user *url.Userinfo
	if config.UseSeparateCredentials {
		user = url.UserPassword(config.SeparateCredentials.Username, config.SeparateCredentials.Password)
	} else {
		user = url.UserPassword(config.Credentials.Username, config.Credentials.Password)
	}
	uri := url.URL{
		Scheme:   config.DSN.Scheme,
		User:     user,
		Host:     config.DSN.Host,
		Path:     config.DSN.DB,
		RawQuery: config.DSN.Params,
	}
	dsn := uri.String()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a connection pool to PostgreSQL")
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection to PMM DB failed: %s", err)
	}

	return db, nil
}

func (d *dsPmmDbSelect) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*pmmv1.ServerMetric_Metric, error) {
	return fetchMetricsFromDB(d.l, d.config.Timeout, d.db, ctx, config)
}
