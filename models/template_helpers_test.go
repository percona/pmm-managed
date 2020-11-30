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

package models_test

import (
	"testing"

	"github.com/percona-platform/saas/pkg/alert"
	"github.com/percona-platform/saas/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestRuleTemplatesChannels(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("create", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, tx.Rollback())
		}()

		q := tx.Querier

		params := &models.CreateTemplateParams{
			Rule: &alert.Rule{
				Name:    "test",
				Version: 1,
				Summary: "test rule",
				Tiers:   []common.Tier{common.Anonymous},
				Expr:    "some expression",
				Params: []alert.Parameter{
					{
						Name:    "param",
						Summary: "test param",
						Unit:    "kg",
						Type:    alert.Float,
						Range:   []interface{}{float64(10), float64(100)},
						Value:   float64(50),
					},
				},
				For:         3,
				Severity:    common.Warning,
				Labels:      map[string]string{"foo": "bar"},
				Annotations: nil,
			},
			Source: "USER_FILE",
		}

		_, err = models.CreateTemplate(q, params)
		require.NoError(t, err)

		templates, err := models.FindTemplates(q)
		require.NoError(t, err)

		require.Len(t, templates, 1)

		actual := templates[0]
		assert.Equal(t, params.Rule.Name, actual.Name)
		assert.Equal(t, params.Rule.Version, actual.Version)
		assert.Equal(t, params.Rule.Summary, actual.Summary)
		assert.ElementsMatch(t, params.Rule.Tiers, actual.Tiers)
		assert.Equal(t, params.Rule.Expr, actual.Expr)
		assert.Equal(t,
			models.Params{
				{
					Name:    params.Rule.Params[0].Name,
					Summary: params.Rule.Params[0].Summary,
					Unit:    params.Rule.Params[0].Unit,
					Type:    string(params.Rule.Params[0].Type),
					FloatParam: &models.FloatParam{
						Default: params.Rule.Params[0].Value.(float64),
						Min:     params.Rule.Params[0].Range[0].(float64),
						Max:     params.Rule.Params[0].Range[1].(float64),
					},
				},
			},
			actual.Params)
		assert.EqualValues(t, params.Rule.For, actual.For)
		assert.Equal(t, params.Rule.Severity.String(), actual.Severity)

		labels, err := actual.GetLabels()
		require.NoError(t, err)
		assert.Equal(t, params.Rule.Labels, labels)

		annotations, err := actual.GetAnnotations()
		require.NoError(t, err)
		assert.Equal(t, params.Rule.Annotations, annotations)

		assert.Equal(t, params.Source, actual.Source)
	})
}
