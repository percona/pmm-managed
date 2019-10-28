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
	"os"
	"sort"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"

	"github.com/percona/pmm-managed/models"
)

// rdsExporterConfig returns desired configuration of rds_exporter process.
func rdsExporterConfig(service *models.Service, exporter *models.Agent) *agentpb.SetStateRequest_AgentProcess {
	tdp := templateDelimsPair(
		pointer.GetString(service.Address),
		pointer.GetString(exporter.Username),
		pointer.GetString(exporter.Password),
		pointer.GetString(exporter.MetricsURL),
	)

	args := []string{
		"--web.listen-address=:" + tdp.left + " .listen_port " + tdp.right,
		fmt.Sprintf("--config.file=%s", os.Getenv("RDS_CONFIG")),
	}

	if pointer.GetString(exporter.MetricsURL) != "" {
		args = append(args, "-web.telemetry-path="+*exporter.MetricsURL)
	}

	sort.Strings(args)

	return &agentpb.SetStateRequest_AgentProcess{
		Type:               inventorypb.AgentType_RDS_EXPORTER,
		TemplateLeftDelim:  tdp.left,
		TemplateRightDelim: tdp.right,
		Args:               args,
		//TODO: get real parameters from config
		Env: []string{
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", "<access_key_id_from_config>"),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", "<secret_access_key_from_config>"),
			fmt.Sprintf("AWS_REGION=%s", "<region_from_config>"),
		},
	}
}
