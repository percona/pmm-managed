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

package alertmanager

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

const (
	testShippedFilePath     = "../../testdata/ia/shipped/*.yml"
	testUserDefinedFilePath = "../../testdata/ia/userdefined/*.yml"
)

func TestAlertmanager(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	svc := New(db)

	require.NoError(t, svc.IsReady(context.Background()))
}

func TestCollect(t *testing.T) {
	err := os.Setenv("PERCONA_TEST_SHIPPED_RULE_TEMPLATE_PATH", testShippedFilePath)
	require.NoError(t, err)
	err = os.Setenv("PERCONA_TEST_USER_DEFINED_RULE_TEMPLATE_PATH", testUserDefinedFilePath)
	require.NoError(t, err)

	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	svc := New(db)
	svc.collectRuleTemplates()

	require.NotNil(t, svc.rules)
	require.Len(t, svc.rules, 2)
	assert.Equal(t, svc.rules[0].Name, "shipped_rules")
	assert.Equal(t, svc.rules[1].Name, "user_defined_rules")
}
