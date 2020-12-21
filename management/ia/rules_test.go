package ia

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/brianvoe/gofakeit"
	"github.com/percona/pmm/api/managementpb/ia/json/client"
	"github.com/percona/pmm/api/managementpb/ia/json/client/channels"
	"github.com/percona/pmm/api/managementpb/ia/json/client/rules"
	"github.com/percona/pmm/api/managementpb/ia/json/client/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

// Note: Even though the IA services check for alerting enabled or disabled before returning results
// we don't enable or disable IA explicit in our tests since it is enabled by default through
// ENABLE_ALERTING env var.
func TestRulesAPI(t *testing.T) {
	templateName := createTemplate(t)
	defer deleteTemplate(t, client.Default.Templates, templateName)

	channelID := createChannel(t)
	defer deleteChannel(t, client.Default.Channels, channelID)

	client := client.Default.Rules

	t.Run("add", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			params := createAlertRuleParams(templateName, channelID)
			rule, err := client.CreateAlertRule(params)
			require.NoError(t, err)
			defer deleteRule(t, client, rule.Payload.RuleID)

			assert.NotEmpty(t, rule.Payload.RuleID)
		})

		t.Run("unknown template", func(t *testing.T) {
			templateName := gofakeit.UUID()
			params := createAlertRuleParams(templateName, channelID)
			_, err := client.CreateAlertRule(params)
			pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Unknown template %s.", templateName)
		})

		t.Run("unknown channel", func(t *testing.T) {
			channelID := gofakeit.UUID()
			params := createAlertRuleParams(templateName, channelID)
			_, err := client.CreateAlertRule(params)
			pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Failed to find all required channels: [%s].", channelID)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			cParams := createAlertRuleParams(templateName, channelID)
			rule, err := client.CreateAlertRule(cParams)
			require.NoError(t, err)
			defer deleteRule(t, client, rule.Payload.RuleID)

			newChannelID := createChannel(t)

			params := &rules.UpdateAlertRuleParams{
				Body: rules.UpdateAlertRuleBody{
					RuleID:   rule.Payload.RuleID,
					Disabled: false,
					Params: []*rules.ParamsItems0{
						{
							Name:  "threshold",
							Type:  pointer.ToString("FLOAT"),
							Float: 21,
						},
					},
					For:          "10s",
					Severity:     pointer.ToString("SEVERITY_ERROR"),
					CustomLabels: map[string]string{"foo": "bar", "baz": "faz"},
					Filters: []*rules.FiltersItems0{
						{
							Type:  pointer.ToString("NOT_EQUAL"),
							Key:   "threshold",
							Value: "21",
						},
					},
					ChannelIds: []string{channelID, newChannelID},
				},
				Context: pmmapitests.Context,
			}
			_, err = client.UpdateAlertRule(params)
			require.NoError(t, err)

			list, err := client.ListAlertRules(&rules.ListAlertRulesParams{Context: pmmapitests.Context})
			require.NoError(t, err)

			var found bool
			for _, r := range list.Payload.Rules {
				if r.RuleID == rule.Payload.RuleID {
					assert.False(t, r.Disabled)
					assert.Equal(t, "10s", r.For)
					assert.Len(t, r.Params, 1)
					assert.Equal(t, params.Body.Params[0].Type, r.Params[0].Type)
					assert.Equal(t, params.Body.Params[0].Name, r.Params[0].Name)
					assert.Equal(t, params.Body.Params[0].Float, r.Params[0].Float)
					assert.Equal(t, params.Body.Params[0].Bool, r.Params[0].Bool)
					assert.Equal(t, params.Body.Params[0].String, r.Params[0].String)
					found = true
				}
			}
			assert.Truef(t, found, "Rule with id %s not found", rule.Payload.RuleID)
		})

		t.Run("unknown channel", func(t *testing.T) {
			cParams := createAlertRuleParams(templateName, channelID)
			rule, err := client.CreateAlertRule(cParams)
			require.NoError(t, err)
			defer deleteRule(t, client, rule.Payload.RuleID)

			unknownChannelID := gofakeit.UUID()
			params := &rules.UpdateAlertRuleParams{
				Body: rules.UpdateAlertRuleBody{
					RuleID:   rule.Payload.RuleID,
					Disabled: false,
					Params: []*rules.ParamsItems0{
						{
							Name:  "threshold",
							Type:  pointer.ToString("FLOAT"),
							Float: 21,
						},
					},
					For:          "10s",
					Severity:     pointer.ToString("SEVERITY_ERROR"),
					CustomLabels: map[string]string{"foo": "bar", "baz": "faz"},
					Filters: []*rules.FiltersItems0{
						{
							Type:  pointer.ToString("NOT_EQUAL"),
							Key:   "threshold",
							Value: "21",
						},
					},
					ChannelIds: []string{channelID, unknownChannelID},
				},
				Context: pmmapitests.Context,
			}
			_, err = client.UpdateAlertRule(params)
			pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Failed to find all required channels: [%s].", unknownChannelID)
		})
	})

	t.Run("toggle", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			cParams := createAlertRuleParams(templateName, channelID)
			rule, err := client.CreateAlertRule(cParams)
			require.NoError(t, err)
			defer deleteRule(t, client, rule.Payload.RuleID)

			list, err := client.ListAlertRules(&rules.ListAlertRulesParams{Context: pmmapitests.Context})
			require.NoError(t, err)

			var found bool
			for _, r := range list.Payload.Rules {
				if r.RuleID == rule.Payload.RuleID {
					assert.True(t, r.Disabled)
					found = true
				}
			}
			assert.Truef(t, found, "Rule with id %s not found", rule.Payload.RuleID)

			_, err = client.ToggleAlertRule(&rules.ToggleAlertRuleParams{
				Body: rules.ToggleAlertRuleBody{
					RuleID:   rule.Payload.RuleID,
					Disabled: pointer.ToString(rules.ToggleAlertRuleBodyDisabledFALSE),
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)

			list, err = client.ListAlertRules(&rules.ListAlertRulesParams{Context: pmmapitests.Context})
			require.NoError(t, err)

			found = false
			for _, r := range list.Payload.Rules {
				if r.RuleID == rule.Payload.RuleID {
					assert.False(t, r.Disabled)
					found = true
				}
			}
			assert.Truef(t, found, "Rule with id %s not found", rule.Payload.RuleID)
		})
	})

	t.Run("delete", func(t *testing.T) {
		params := createAlertRuleParams(templateName, channelID)
		rule, err := client.CreateAlertRule(params)
		require.NoError(t, err)

		_, err = client.DeleteAlertRule(&rules.DeleteAlertRuleParams{
			Body:    rules.DeleteAlertRuleBody{RuleID: rule.Payload.RuleID},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		list, err := client.ListAlertRules(&rules.ListAlertRulesParams{Context: pmmapitests.Context})
		require.NoError(t, err)

		for _, r := range list.Payload.Rules {
			assert.NotEqual(t, rule.Payload.RuleID, r.RuleID)
		}
	})

	t.Run("list", func(t *testing.T) {
		params := createAlertRuleParams(templateName, channelID)
		rule, err := client.CreateAlertRule(params)
		require.NoError(t, err)
		defer deleteRule(t, client, rule.Payload.RuleID)

		list, err := client.ListAlertRules(&rules.ListAlertRulesParams{Context: pmmapitests.Context})
		require.NoError(t, err)

		var found bool
		for _, r := range list.Payload.Rules {
			if r.RuleID == rule.Payload.RuleID {
				assert.True(t, r.Disabled)
				assert.Equal(t, params.Body.Summary, r.Summary)
				assert.Len(t, r.Params, 1)
				assert.Equal(t, params.Body.Params[0].Type, r.Params[0].Type)
				assert.Equal(t, params.Body.Params[0].Name, r.Params[0].Name)
				assert.Equal(t, params.Body.Params[0].Float, r.Params[0].Float)
				assert.Equal(t, params.Body.Params[0].Bool, r.Params[0].Bool)
				assert.Equal(t, params.Body.Params[0].String, r.Params[0].String)
				assert.Equal(t, params.Body.For, r.For)
				assert.Equal(t, params.Body.Severity, r.Severity)
				assert.Equal(t, params.Body.CustomLabels, r.CustomLabels)
				assert.Len(t, params.Body.Filters, 1)
				assert.Equal(t, params.Body.Filters[0].Type, r.Filters[0].Type)
				assert.Equal(t, params.Body.Filters[0].Key, r.Filters[0].Key)
				assert.Equal(t, params.Body.Filters[0].Value, r.Filters[0].Value)
				assert.Len(t, r.Channels, 1)
				assert.Equal(t, r.Channels[0].ChannelID, channelID)
				found = true
			}
		}
		assert.Truef(t, found, "Rule with id %s not found", rule.Payload.RuleID)
	})
}

func deleteRule(t *testing.T, client rules.ClientService, id string) {
	_, err := client.DeleteAlertRule(&rules.DeleteAlertRuleParams{
		Body:    rules.DeleteAlertRuleBody{RuleID: id},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
}

func createAlertRuleParams(templateName, channelID string) *rules.CreateAlertRuleParams {
	return &rules.CreateAlertRuleParams{
		Body: rules.CreateAlertRuleBody{
			TemplateName: templateName,
			Disabled:     true,
			Summary:      "example summary",
			Params: []*rules.ParamsItems0{
				{
					Name:  "threshold",
					Type:  pointer.ToString("FLOAT"),
					Float: 12,
				},
			},
			For:          "5s",
			Severity:     pointer.ToString("SEVERITY_WARNING"),
			CustomLabels: map[string]string{"foo": "bar"},
			Filters: []*rules.FiltersItems0{
				{
					Type:  pointer.ToString("EQUAL"),
					Key:   "threshold",
					Value: "12",
				},
			},
			ChannelIds: []string{channelID},
		},
		Context: pmmapitests.Context,
	}
}

func createTemplate(t *testing.T) string {
	b, err := ioutil.ReadFile("../../testdata/ia/template.yaml")
	require.NoError(t, err)

	templateName := gofakeit.UUID()
	_, err = client.Default.Templates.CreateTemplate(&templates.CreateTemplateParams{
		Body: templates.CreateTemplateBody{
			Yaml: fmt.Sprintf(string(b), templateName, gofakeit.UUID()),
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)

	return templateName
}

func createChannel(t *testing.T) string {
	resp, err := client.Default.Channels.AddChannel(&channels.AddChannelParams{
		Body: channels.AddChannelBody{
			Summary:  gofakeit.Quote(),
			Disabled: gofakeit.Bool(),
			EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
				SendResolved: false,
				To:           []string{gofakeit.Email()},
			},
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	return resp.Payload.ChannelID
}
