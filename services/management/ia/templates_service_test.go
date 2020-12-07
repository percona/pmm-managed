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

package ia

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"gopkg.in/yaml.v3"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

const (
	testBadTemplates   = "../../../testdata/ia/bad/*.yml"
	testUser2Templates = "../../../testdata/ia/user2/*.yml"
	testUserTemplates  = "../../../testdata/ia/user/*.yml"

	userRuleFilepath     = "/etc/ia/rules/user_rule.yml"
	builtinRuleFilepath1 = "/etc/ia/rules/mysql_down.yml"
	builtinRuleFilepath2 = "/etc/ia/rules/mysql_restarted.yml"
	builtinRuleFilepath3 = "/etc/ia/rules/mysql_too_many_connections.yml"
)

func TestCollect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("bad template paths", func(t *testing.T) {
		t.Parallel()

		testDir, err := ioutil.TempDir("", "")
		require.NoError(t, err)
		defer os.RemoveAll(testDir) //nolint:errcheck

		svc := NewTemplatesService(db)
		svc.userTemplatesPath = testBadTemplates
		svc.collect(ctx)

		require.Empty(t, svc.getCollected(ctx))
	})

	t.Run("valid template paths", func(t *testing.T) {
		t.Parallel()

		testDir, err := ioutil.TempDir("", "")
		require.NoError(t, err)
		defer os.RemoveAll(testDir) //nolint:errcheck

		svc := NewTemplatesService(db)
		svc.userTemplatesPath = testUserTemplates
		svc.collect(ctx)

		// todo rename rules to templates
		rules := svc.getCollected(ctx)
		require.NotEmpty(t, rules)
		require.Len(t, rules, 4)
		assert.Contains(t, rules, "user_rule")
		assert.Contains(t, rules, "mysql_down")
		assert.Contains(t, rules, "mysql_restarted")
		assert.Contains(t, rules, "mysql_too_many_connections")

		// check whether map was cleared and updated on a subsequent call
		svc.userTemplatesPath = testUser2Templates
		svc.collect(ctx)

		rules = svc.getCollected(ctx)
		require.NotEmpty(t, rules)
		require.Len(t, rules, 4)
		assert.NotContains(t, rules, "user_rule")
		assert.Contains(t, rules, "user2_rule")
	})
}

func TestConvertTemplate(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		ctx := context.Background()

		testDir, err := ioutil.TempDir("", "")
		require.NoError(t, err)
		defer os.RemoveAll(testDir) //nolint:errcheck

		svc := NewTemplatesService(db)
		svc.userTemplatesPath = testUserTemplates
		svc.collect(ctx)

		err = svc.convertTemplates(ctx)
		require.NoError(t, err)
		assert.FileExists(t, userRuleFilepath)
		assert.FileExists(t, builtinRuleFilepath1)
		assert.FileExists(t, builtinRuleFilepath2)
		assert.FileExists(t, builtinRuleFilepath3)

		testcases := []struct {
			path             string
			alert            string
			labelsCount      int
			annotationsCount int
		}{{
			path:             builtinRuleFilepath1,
			alert:            "mysql_down",
			labelsCount:      3,
			annotationsCount: 2,
		}, {
			path:             builtinRuleFilepath2,
			alert:            "mysql_restarted",
			labelsCount:      4,
			annotationsCount: 2,
		}, {
			path:             builtinRuleFilepath3,
			alert:            "mysql_too_many_connections",
			labelsCount:      4,
			annotationsCount: 2,
		}, {
			path:             userRuleFilepath,
			alert:            "user_rule",
			labelsCount:      4,
			annotationsCount: 2,
		}}

		for _, tc := range testcases {
			t.Run(tc.path, func(t *testing.T) {
				buf, err := ioutil.ReadFile(tc.path)
				require.NoError(t, err)
				var rf ruleFile
				err = yaml.Unmarshal(buf, &rf)
				require.NoError(t, err)
				rule := rf.Group[0].Rules[0]
				assert.Equal(t, tc.alert, rule.Alert)
				assert.Len(t, rule.Labels, tc.labelsCount)
				assert.Contains(t, rule.Labels, "severity")
				assert.Contains(t, rule.Labels, "ia")
				assert.NotNil(t, rule.Annotations)
				assert.Len(t, rule.Annotations, tc.annotationsCount)
			})
		}
	})
}
