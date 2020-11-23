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
