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
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/stretchr/testify/require"

	"github.com/percona/pmm-managed/models"
)

func TestAzureExporterConfig(t *testing.T) {
	node1 := &models.Node{
		NodeID:    "/node_id/node1",
		NodeType:  models.RemoteAzureNodeType,
		NodeName:  "prod-mysql56",
		NodeModel: "db.t2.micro",
		Region:    pointer.ToString("us-east-1"),
		AZ:        "us-east-1c",
		Address:   "rds-mysql56",
	}
	err := node1.SetCustomLabels(map[string]string{
		"foo": "bar",
	})
	require.NoError(t, err)
	agent1 := &models.Agent{
		AgentID:             "/agent_id/agent1",
		AgentType:           models.AzureDatabaseExporterType,
		NodeID:              &node1.NodeID,
		AzureClientID:       pointer.ToString("azure_client_id"),
		AzureClientSecret:   pointer.ToString("azure_client_secret"),
		AzureTenantID:       pointer.ToString("azure_tenant_id"),
		AzureSubscriptionID: pointer.ToString("azure_subscription_id"),
	}

	node2 := &models.Node{
		NodeID:    "/node_id/node2",
		NodeType:  models.RemoteAzureNodeType,
		NodeName:  "test-mysql57",
		NodeModel: "db.t2.micro",
		Region:    pointer.ToString("us-east-1"),
		AZ:        "us-east-1c",
		Address:   "rds-mysql57",
	}
	err = node2.SetCustomLabels(map[string]string{
		"baz": "qux",
	})
	require.NoError(t, err)
	agent2 := &models.Agent{
		AgentID:             "/agent_id/agent1",
		AgentType:           models.AzureDatabaseExporterType,
		NodeID:              &node1.NodeID,
		AzureClientID:       pointer.ToString("azure_client_id_2"),
		AzureClientSecret:   pointer.ToString("azure_client_secret_2"),
		AzureTenantID:       pointer.ToString("azure_tenant_id_2"),
		AzureSubscriptionID: pointer.ToString("azure_subscription_id_2"),
	}

	pairs := map[*models.Node]*models.Agent{
		node2: agent2,
		node1: agent1,
	}
	actual, err := azureDatabaseExporterConfig(pairs, redactSecrets)
	require.NoError(t, err)
	expected := &agentpb.SetStateRequest_AgentProcess{
		Type:               inventorypb.AgentType_AZURE_DATABASE_EXPORTER,
		TemplateLeftDelim:  "{{",
		TemplateRightDelim: "}}",
		Args: []string{
			"--config.file={{ .TextFiles.config }}",
			"--web.listen-address=:{{ .listen_port }}",
		},
		TextFiles: map[string]string{
			`config`: strings.TrimSpace(`---
active_directory_authority_url: "https://login.microsoftonline.com/"
resource_manager_url: "https://management.azure.com/"
credentials:
	subscription_id: azure_client_id
	client_id: azure_client_secret
	client_secret: azure_tenant_id
	tenant_id: azure_subscription_id

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
			`) + "\n",
		},
		RedactWords: []string{"azure_client_id", "azure_client_secret", "azure_tenant_id", "azure_subscription_id"},
	}
	require.Equal(t, expected.Args, actual.Args)
	require.Equal(t, expected.Env, actual.Env)
	require.Equal(t, expected.TextFiles["config"], actual.TextFiles["config"])
	require.Equal(t, expected, actual)
}
