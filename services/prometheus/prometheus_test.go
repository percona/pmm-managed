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

package prometheus

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/tests"
)

func TestDefaultConfig(t *testing.T) {
	sqlDB := tests.OpenTestDB(t)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	// always restore original file after test
	configPath := filepath.Join("..", "..", "testdata", "prometheus", "prometheus.yml")
	original, err := ioutil.ReadFile(configPath) //nolint:gosec
	require.NoError(t, err)
	defer func() {
		require.NoError(t, ioutil.WriteFile(configPath, original, 0644))
	}()

	ctx := logger.Set(context.Background(), t.Name())
	svc, err := NewService(configPath, "promtool", db, "http://127.0.0.1:9090/prometheus/")
	require.NoError(t, err)
	require.NoError(t, svc.Check(ctx))

	assert.NoError(t, svc.UpdateConfiguration(ctx))

	b, err := ioutil.ReadFile(configPath) //nolint:gosec
	require.NoError(t, err)
	assert.Equal(t, string(original), string(b))
}
