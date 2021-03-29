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
		NodeType:  models.RemoteAzureDatabaseNodeType,
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
	creds := `
	{
		"client_id": "azure_database_client_id",
		"client_secret": "azure_database_client_secret",
		"tenant_id": "azure_database_tenant_id",
		"subscription_id": "azure_database_subscription_id",
		"resource_group": "azure_database_resource_group"
	}
	`

	service1 := &models.Service{
		ServiceID:   "/service_id/service1",
		NodeID:      node1.NodeID,
		ServiceName: "service1",
		ServiceType: "mysql",
	}

	agent := &models.Agent{
		AgentID:          "/agent_id/agent1",
		AgentType:        models.AzureDatabaseExporterType,
		NodeID:           &node1.NodeID,
		ServiceID:        &service1.ServiceID,
		AzureCredentials: pointer.ToString(creds),
	}

	actual, err := azureDatabaseExporterConfig(agent, service1, redactSecrets)
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
  client_id: "azure_database_client_id"
  client_secret: "azure_database_client_secret"
  tenant_id: "azure_database_tenant_id"
  subscription_id: "azure_database_subscription_id"

targets:
resource_groups:
  - resource_group: "azure_database_resource_group"
    aggregations:
      - Average
    resource_types:
      - "Microsoft.DBforMySQL/servers"
      - "Microsoft.DBforMySQL/felexibleServers"
    metrics:
      - name: "cpu_percent"
      - name: "memory_percent"
      - name: "storage_percent"
      - name: "storage_used"
      - name: "storage_limit"
      - name: "network_bytes_egress"
      - name: "network_bytes_ingress"
			`) + "\n",
		},
		RedactWords: []string{
			"azure_database_client_id",
			"azure_database_client_secret",
			"azure_database_tenant_id",
			"azure_database_subscription_id",
		},
	}
	require.Equal(t, expected.Args, actual.Args)
	require.Equal(t, expected.Env, actual.Env)
	require.Equal(t, expected.TextFiles["config"], actual.TextFiles["config"])
	require.Equal(t, expected, actual)
}
