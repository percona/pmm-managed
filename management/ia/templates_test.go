package ia

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/percona-platform/saas/pkg/alert"
	templatesClient "github.com/percona/pmm/api/managementpb/ia/json/client"
	"github.com/percona/pmm/api/managementpb/ia/json/client/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

// Note: Even though the IA services check for alerting enabled or disabled before returning results
// we don't enable or disable IA explicit in our tests since it is enabled by default through
// ENABLE_ALERTING env var.
func TestAddTemplate(t *testing.T) {
	client := templatesClient.Default.Templates

	b, err := ioutil.ReadFile("../../testdata/ia/template.yaml")
	require.NoError(t, err)

	t.Run("normal", func(t *testing.T) {
		name := gofakeit.UUID()
		expr := gofakeit.UUID()
		yml := formatTemplateYaml(t, fmt.Sprintf(string(b), name, expr, "%"))
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: yml,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteTemplate(t, client, name)

		resp, err := client.ListTemplates(&templates.ListTemplatesParams{
			Body: templates.ListTemplatesBody{
				Reload: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		var found bool
		for _, template := range resp.Payload.Templates {
			if template.Name == name {
				assert.Equal(t, yml, template.Yaml)
				assert.Equal(t, "Test summary", template.Summary)
				assert.Equal(t, expr, template.Expr)
				assert.Len(t, template.Params, 1)
				param := template.Params[0]
				assert.Equal(t, "threshold", param.Name)
				assert.Equal(t, "test param summary", param.Summary)
				assert.Equal(t, "PERCENTAGE", *param.Unit)
				assert.Equal(t, "FLOAT", *param.Type)
				assert.True(t, param.Float.HasDefault)
				assert.Equal(t, float32(80), param.Float.Default)
				assert.True(t, param.Float.HasMax)
				assert.Equal(t, float32(100), param.Float.Max)
				assert.True(t, param.Float.HasMin)
				assert.Equal(t, float32(0), param.Float.Min)
				found = true
			}
		}
		assert.Truef(t, found, "Template with id %s not found", name)
	})

	t.Run("duplicate", func(t *testing.T) {
		name := gofakeit.UUID()
		yaml := fmt.Sprintf(string(b), name, gofakeit.UUID())
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: yaml,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteTemplate(t, client, name)

		_, err = client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: yaml,
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, fmt.Sprintf("Template with name \"%s\" already exists.", name))
	})

	t.Run("invalid yaml", func(t *testing.T) {
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: "not a yaml",
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Failed to parse rule template.")
	})

	t.Run("invalid template", func(t *testing.T) {
		b, err := ioutil.ReadFile("../../testdata/ia/invalid-template.yaml")
		require.NoError(t, err)
		name := gofakeit.UUID()
		_, err = client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID()),
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Failed to parse rule template.")
	})
}

func TestChangeTemplate(t *testing.T) {
	client := templatesClient.Default.Templates

	b, err := ioutil.ReadFile("../../testdata/ia/template.yaml")
	require.NoError(t, err)

	t.Run("normal", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID()),
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteTemplate(t, client, name)

		newExpr := gofakeit.UUID()
		yml := formatTemplateYaml(t, fmt.Sprintf(string(b), name, newExpr, "s"))
		_, err = client.UpdateTemplate(&templates.UpdateTemplateParams{
			Body: templates.UpdateTemplateBody{
				Yaml: yml,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		resp, err := client.ListTemplates(&templates.ListTemplatesParams{
			Body: templates.ListTemplatesBody{
				Reload: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		var found bool
		for _, template := range resp.Payload.Templates {
			if template.Name == name {
				assert.Equal(t, newExpr, template.Expr)
				assert.Equal(t, yml, template.Yaml)
				assert.Equal(t, "Test summary", template.Summary)
				assert.Len(t, template.Params, 1)
				param := template.Params[0]
				assert.Equal(t, "threshold", param.Name)
				assert.Equal(t, "test param summary", param.Summary)
				assert.Equal(t, "SECONDS", *param.Unit)
				assert.Equal(t, "FLOAT", *param.Type)
				assert.True(t, param.Float.HasDefault)
				assert.Equal(t, float32(80), param.Float.Default)
				assert.True(t, param.Float.HasMax)
				assert.Equal(t, float32(100), param.Float.Max)
				assert.True(t, param.Float.HasMin)
				assert.Equal(t, float32(0), param.Float.Min)
				found = true
			}
		}
		assert.Truef(t, found, "Template with id %s not found", name)
	})

	t.Run("unknown template", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err = client.UpdateTemplate(&templates.UpdateTemplateParams{
			Body: templates.UpdateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID()),
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, fmt.Sprintf("Template with name \"%s\" not found.", name))
	})

	t.Run("invalid yaml", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID()),
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteTemplate(t, client, name)

		_, err = client.UpdateTemplate(&templates.UpdateTemplateParams{
			Body: templates.UpdateTemplateBody{
				Yaml: "not a yaml",
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Failed to parse rule template.")
	})

	t.Run("invalid template", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err = client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID()),
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteTemplate(t, client, name)

		b, err = ioutil.ReadFile("../../testdata/ia/invalid-template.yaml")
		_, err = client.UpdateTemplate(&templates.UpdateTemplateParams{
			Body: templates.UpdateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID()),
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Failed to parse rule template.")
	})
}

func TestDeleteTemplate(t *testing.T) {
	client := templatesClient.Default.Templates

	b, err := ioutil.ReadFile("../../testdata/ia/template.yaml")
	require.NoError(t, err)

	t.Run("normal", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID(), "s"),
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		_, err = client.DeleteTemplate(&templates.DeleteTemplateParams{
			Body: templates.DeleteTemplateBody{
				Name: name,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		resp, err := client.ListTemplates(&templates.ListTemplatesParams{
			Body: templates.ListTemplatesBody{
				Reload: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		for _, template := range resp.Payload.Templates {
			assert.NotEqual(t, name, template.Name)
		}
	})

	t.Run("template in use", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err := client.CreateTemplate(&templates.CreateTemplateParams{
			Body: templates.CreateTemplateBody{
				Yaml: fmt.Sprintf(string(b), name, gofakeit.UUID(), "s"),
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		channelID := createChannel(t)
		defer deleteChannel(t, templatesClient.Default.Channels, channelID)

		params := createAlertRuleParams(name, channelID)

		rule, err := templatesClient.Default.Rules.CreateAlertRule(params)
		require.NoError(t, err)

		_, err = client.DeleteTemplate(&templates.DeleteTemplateParams{
			Body: templates.DeleteTemplateBody{
				Name: name,
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 500, codes.Internal, "Internal server error.")

		defer deleteTemplate(t, templatesClient.Default.Templates, name)
		defer deleteRule(t, templatesClient.Default.Rules, rule.Payload.RuleID)

		resp, err := client.ListTemplates(&templates.ListTemplatesParams{
			Body: templates.ListTemplatesBody{
				Reload: true,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		var found bool
		for _, template := range resp.Payload.Templates {
			if name == template.Name {
				found = true
			}
		}
		assert.Truef(t, found, "Template with id %s not found", name)
	})

	t.Run("unknown template", func(t *testing.T) {
		name := gofakeit.UUID()
		_, err = client.DeleteTemplate(&templates.DeleteTemplateParams{
			Body: templates.DeleteTemplateBody{
				Name: name,
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, fmt.Sprintf("Template with name \"%s\" not found.", name))
	})
}

func TestListTemplate(t *testing.T) {
	client := templatesClient.Default.Templates

	b, err := ioutil.ReadFile("../../testdata/ia/template.yaml")
	require.NoError(t, err)

	name := gofakeit.UUID()
	expr := gofakeit.UUID()
	yml := formatTemplateYaml(t, fmt.Sprintf(string(b), name, expr, "%"))
	_, err = client.CreateTemplate(&templates.CreateTemplateParams{
		Body: templates.CreateTemplateBody{
			Yaml: yml,
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	defer deleteTemplate(t, client, name)

	resp, err := client.ListTemplates(&templates.ListTemplatesParams{
		Body: templates.ListTemplatesBody{
			Reload: true,
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)

	var found bool
	for _, template := range resp.Payload.Templates {
		if template.Name == name {
			assert.Equal(t, expr, template.Expr)
			assert.Equal(t, "Test summary", template.Summary)
			assert.Equal(t, "USER_API", *template.Source)
			assert.Equal(t, "SEVERITY_WARNING", *template.Severity)
			assert.Equal(t, "300s", template.For)
			assert.Len(t, template.Params, 1)

			param := template.Params[0]
			assert.Equal(t, "threshold", param.Name)
			assert.Equal(t, "test param summary", param.Summary)
			assert.Equal(t, "FLOAT", *param.Type)
			assert.Equal(t, "PERCENTAGE", *param.Unit)
			assert.Nil(t, param.Bool)
			assert.Nil(t, param.String)
			assert.NotNil(t, param.Float)

			float := param.Float
			assert.True(t, float.HasDefault)
			assert.Equal(t, float32(80), float.Default)
			assert.True(t, float.HasMax)
			assert.Equal(t, float32(100), float.Max)
			assert.True(t, float.HasMin)
			assert.Equal(t, float32(0), float.Min)

			assert.Equal(t, map[string]string{"foo": "bar"}, template.Labels)
			assert.Equal(t, map[string]string{"description": "test description", "summary": "test summary"}, template.Annotations)
			assert.Equal(t, yml, template.Yaml)
			assert.NotEmpty(t, template.CreatedAt)
			found = true
		}
	}
	assert.Truef(t, found, "Template with id %s not found", name)
}

func deleteTemplate(t *testing.T, client templates.ClientService, name string) {
	_, err := client.DeleteTemplate(&templates.DeleteTemplateParams{
		Body: templates.DeleteTemplateBody{
			Name: name,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
}

func formatTemplateYaml(t *testing.T, yml string) string {
	params := &alert.ParseParams{
		DisallowUnknownFields:    true,
		DisallowInvalidTemplates: true,
	}
	r, err := alert.Parse(strings.NewReader(yml), params)
	require.NoError(t, err)
	type templates struct {
		Templates []alert.Template `yaml:"templates"`
	}

	s, err := yaml.Marshal(&templates{Templates: r})
	require.NoError(t, err)

	return string(s)
}
