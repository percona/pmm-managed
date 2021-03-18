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

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"

	"github.com/percona/pmm-managed/models"
)

const azureDatabaseTemplate = `---
active_directory_authority_url: "https://login.microsoftonline.com/"
resource_manager_url: "https://management.azure.com/"
credentials:
  client_id: "{{ .AzureDatabaseClientID}}"
  client_secret: "{{ .AzureDatabaseClientSecret}}"
  tenant_id: "{{ .AzureDatabaseTenantID}}"
  subscription_id: "{{ .AzureDatabaseSubscriptionID}}"

targets:
resource_groups:
  - resource_group: "pmmdemo"
    aggregations:
      - Average
{{ .AzureDatabaseResourceTypes }}
    metrics:
      - name: "cpu_percent"
      - name: "memory_percent"
      - name: "storage_percent"
      - name: "storage_used"
      - name: "storage_limit"
      - name: "network_bytes_egress"
      - name: "network_bytes_ingress"
`

// azureDatabaseInstance represents credentials informations.
type azureDatabaseCredentials struct {
	AzureDatabaseClientID       string
	AzureDatabaseClientSecret   string
	AzureDatabaseTenantID       string
	AzureDatabaseSubscriptionID string
	AzureDatabaseResourceTypes  string
}

// azureDatabaseExporterConfig returns configuration of azure_database_exporter process.
func azureDatabaseExporterConfig(exporter *models.Agent, redactMode redactMode) (*agentpb.SetStateRequest_AgentProcess, error) {
	t, err := template.New("credentials").Parse(azureDatabaseTemplate)
	if err != nil {
		return nil, err
	}

	exporterType := pointer.GetString(exporter.AzureDatabaseExporterType)
	var resourceTypes string
	switch exporterType {
	case "mysql":
		resourceTypes = `    resource_types:
      - "Microsoft.DBforMySQL/servers"
      - "Microsoft.DBforMySQL/felexibleServers"`
	case "maria":
		resourceTypes = `    resource_types:
      - "Microsoft.DBforMariaDB/servers"`
	case "postgres":
		resourceTypes = `    resource_types:
      - "Microsoft.DBforPostgreSQL/servers"
      - "Microsoft.DBforPostgreSQL/flexibleServers"
      - "Microsoft.DBforPostgreSQL/serversv2"`
	}

	var config bytes.Buffer
	credentials := azureDatabaseCredentials{
		pointer.GetString(exporter.AzureDatabaseClientID),
		pointer.GetString(exporter.AzureDatabaseClientSecret),
		pointer.GetString(exporter.AzureDatabaseTenantID),
		pointer.GetString(exporter.AzureDatabaseSubscriptionID),
		resourceTypes,
	}
	err = t.Execute(&config, credentials)
	if err != nil {
		return nil, err
	}

	var words []string
	if redactMode != exposeSecrets {
		words = append(words, redactWords(exporter)...)
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
		RedactWords: words,
	}, nil
}
