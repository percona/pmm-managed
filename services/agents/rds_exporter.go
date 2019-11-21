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
	"fmt"
	"sort"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"

	"github.com/percona/pmm-managed/models"
)

// rdsExporterConfig returns desired configuration of rds_exporter process.
func rdsExporterConfig(nodes map[*models.Node]*models.Agent) map[*models.Node]*agentpb.SetStateRequest_AgentProcess {

	configs := make(map[*models.Node]*agentpb.SetStateRequest_AgentProcess)

	for node, exporter := range nodes {
		tdp := templateDelimsPair(
			fmt.Sprintf("%d", exporter.ListenPort),
		)

		args := []string{
			"--web.listen-address=" + tdp.left + " .listen_port " + tdp.right,
			"--config.file=", tdp.left + " .config_file " + tdp.right,
		}

		if pointer.GetString(exporter.MetricsURL) != "" {
			// TODO https://jira.percona.com/browse/PMM-1901
		}

		sort.Strings(args)

		// TODO
		// configBody := struct {
		// 	AWSRegion          string `yaml:"aws_region"`
		// 	AWSAccessKeyID     string `yaml:"aws_access_key_id"`
		// 	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
		// }{
		// 	AWSRegion:          *node.Region,
		// 	AWSAccessKeyID:     *exporter.AWSAccessKey,
		// 	AWSSecretAccessKey: *exporter.AWSSecretKey,
		// }
		// buf, _ := yaml.Marshal(configBody)

		configs[node] = &agentpb.SetStateRequest_AgentProcess{
			Type:               inventorypb.AgentType_RDS_EXPORTER,
			TemplateLeftDelim:  tdp.left,
			TemplateRightDelim: tdp.right,
			Args:               args,
			TextFiles:          map[string]string{"AWS_REGION": *node.Region},
			// TODO ConfigBody:         string(buf),
		}
	}

	return configs
}
