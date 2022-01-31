package telemetry_v2

import (
	"context"
	"database/sql"
	"fmt"
	reporter "github.com/percona-platform/saas/gen/telemetry/reporter"
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
	db, err := openDB(config)
	if err != nil {
		return nil, err
	}

	return &dsPmmDbSelect{
		l:      l,
		config: config,
		db:     db,
	}, nil
}

func openDB(config DSConfigPMMDB) (*sql.DB, error) {
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

func (d *dsPmmDbSelect) FetchMetrics(ctx context.Context, config TelemetryConfig) ([]*reporter.ServerMetric_Metric, error) {
	localCtx, _ := context.WithTimeout(ctx, d.config.Timeout)
	tx, err := d.db.BeginTx(localCtx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	// to minimize risk of modifying DB
	defer tx.Rollback()

	rows, err := d.db.Query("SELECT " + config.Query)
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
			d.l.Error(err)
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
