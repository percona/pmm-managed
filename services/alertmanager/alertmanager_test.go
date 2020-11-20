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
	testInvalidFilePath     = "../../testdata/ia/invalid/*.yml"
)

func TestAlertmanager(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	svc := New(db)

	require.NoError(t, svc.IsReady(context.Background()))
}

func TestCollect(t *testing.T) {
	t.Run("invalid template paths", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

		svc := New(db)
		svc.shippedRuleTemplatePath = testInvalidFilePath
		svc.userDefinedRuleTemplatePath = testInvalidFilePath
		svc.collectRuleTemplates()

		require.Empty(t, svc.rules)
	})

	t.Run("valid template paths", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

		svc := New(db)
		svc.shippedRuleTemplatePath = testShippedFilePath
		svc.userDefinedRuleTemplatePath = testUserDefinedFilePath
		svc.collectRuleTemplates()

		require.NotEmpty(t, svc.rules)
		require.Len(t, svc.rules, 2)
		assert.Equal(t, svc.rules[0].Name, "shipped_rules")
		assert.Equal(t, svc.rules[1].Name, "user_defined_rules")
	})
}
