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

package agents

import (
	"bytes"
	"text/template"

	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/prometheus/common/model"

	"github.com/percona/pmm-managed/models"
)

const azureDatabaseTemplate = `---
active_directory_authority_url: "https://login.microsoftonline.com/"
resource_manager_url: "https://management.azure.com/"
credentials:
	subscription_id: {{ .AzureClientID}}
	client_id: {{ .AzureClientSecret}}
	client_secret: {{ .AzureTenantID}}
	tenant_id: {{ .AzureSubscriptionID}}

targets:
	- resource: "/resourceGroups/blog-group/providers/Microsoft.Web/sites/blog"
	metrics:
		- name: "BytesReceived"
		- name: "BytesSent"
	- resource: "/resourceGroups/app-group/providers/Microsoft.Web/sites/app"
	metrics:
		- name: "Http2xx"
		- name: "Http5xx"
	- resource: "/resourceGroups/vm-group/providers/Microsoft.Compute/virtualMachines/vm"
	metric_namespace: "Azure.VM.Windows.GuestMetrics"
	metrics:
	- name: 'Process\Thread Count'

resource_groups:
	- resource_group: "webapps"
	resource_types:
		- "Microsoft.Compute/virtualMachines"
	resource_name_include_re:
		- "testvm.*"
	resource_name_exclude_re:
		- "testvm12"
	metrics:
		- name: "CPU Credits Consumed"

resource_tags:
	- resource_tag_name: "monitoring"
	resource_tag_value: "enabled"
	resource_types:
		- "Microsoft.Compute/virtualMachines"
	metrics:
		- name: "CPU Credits consumed"
`

// azureDatabaseInstance represents credentials informations.
type azureDatabaseCredentials struct {
	AzureClientID       string
	AzureClientSecret   string
	AzureTenantID       string
	AzureSubscriptionID string
}

// azureDatabaseInstance represents information from configuration file.
type azureDatabaseInstance struct {
	Region                 string         `yaml:"region"`
	Instance               string         `yaml:"instance"`
	AzureClientID          string         `yaml:"azure_client_id,omitempty"`
	AzureClientSecret      string         `yaml:"azure_client_secret,omitempty"`
	AzureTenantID          string         `yaml:"azure_tenant_id,omitempty"`
	AzureSubscriptionID    string         `yaml:"azure_subscription_id,omitempty"`
	DisableBasicMetrics    bool           `yaml:"disable_basic_metrics"`
	DisableEnhancedMetrics bool           `yaml:"disable_enhanced_metrics"`
	Labels                 model.LabelSet `yaml:"labels,omitempty"`
}

// azureDatabaseExporterConfig returns configuration of azure_database_exporter process.
func azureDatabaseExporterConfig(pairs map[*models.Node]*models.Agent, redactMode redactMode) (*agentpb.SetStateRequest_AgentProcess, error) {
	t, err := template.New("credentials").Parse(azureDatabaseTemplate)
	if err != nil {
		return nil, err
	}

	var config bytes.Buffer
	credentials := azureDatabaseCredentials{"azure_client_id", "azure_client_secret", "azure_tenant_id", "azure_subscription_id"}
	err = t.Execute(&config, credentials)
	if err != nil {
		return nil, err
	}

	tdp := models.TemplateDelimsPair()
	args := []string{
		"--config.file=" + tdp.Left + " .TextFiles.config " + tdp.Right,
		"--web.listen-address=:" + tdp.Left + " .listen_port " + tdp.Right,
	}
	return &agentpb.SetStateRequest_AgentProcess{
		Type:               inventorypb.AgentType_AZURE_DATABASE_EXPORTER,
		TemplateLeftDelim:  tdp.Left,
		TemplateRightDelim: tdp.Right,
		Args:               args,
		TextFiles: map[string]string{
			"config": config.String(),
		},
		RedactWords: []string{"azure_client_id", "azure_client_secret", "azure_tenant_id", "azure_subscription_id"},
	}, nil
}
