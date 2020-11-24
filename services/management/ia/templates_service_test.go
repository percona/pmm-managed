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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testShippedFilePath     = "../../../testdata/ia/shipped/*.yml"
	testUserDefinedFilePath = "../../../testdata/ia/userdefined/*.yml"
	testOtherRulesFilePath  = "../../../testdata/ia/others/*.yml"
	testInvalidFilePath     = "../../../testdata/ia/invalid/*.yml"
)

func TestCollect(t *testing.T) {
	t.Run("invalid template paths", func(t *testing.T) {
		svc := NewTemplatesService()
		svc.shippedRuleTemplatePath = testInvalidFilePath
		svc.userDefinedRuleTemplatePath = testInvalidFilePath
		svc.collectRuleTemplates()

		require.Empty(t, svc.rules)
	})

	t.Run("valid template paths", func(t *testing.T) {
		svc := NewTemplatesService()
		svc.shippedRuleTemplatePath = testShippedFilePath
		svc.userDefinedRuleTemplatePath = testUserDefinedFilePath
		svc.collectRuleTemplates()

		require.NotEmpty(t, svc.rules)
		require.Len(t, svc.rules, 2)
		assert.Contains(t, svc.rules, "shipped_rules")
		assert.Contains(t, svc.rules, "user_defined_rules")

		// check whether map was cleared and updated on a subsequent call
		svc.userDefinedRuleTemplatePath = testOtherRulesFilePath
		svc.collectRuleTemplates()

		require.NotEmpty(t, svc.rules)
		require.Len(t, svc.rules, 2)
		assert.NotContains(t, svc.rules, "user_defined_rules")
		assert.Contains(t, svc.rules, "shipped_rules")
		assert.Contains(t, svc.rules, "other_rules")
	})

}
