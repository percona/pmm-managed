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
	"net"
	"sort"
	"strconv"

	"github.com/AlekSi/pointer"
	api "github.com/percona/pmm/api/agent"

	"github.com/percona/pmm-managed/models"
)

func mongodbExporterConfig(service *models.Service, exporter *models.Agent) *api.SetStateRequest_AgentProcess {
	tdp := templateDelimsPair(
		pointer.GetString(service.Address),
		pointer.GetString(exporter.Username),
		pointer.GetString(exporter.Password),
		pointer.GetString(exporter.MetricsURL),
	)

	args := []string{
		"--collect.database",
		"--collect.collection",
		"--collect.topmetrics",
		"--collect.indexusage",
		"--web.listen-address=:" + tdp.left + " .listen_port " + tdp.right,
	}

	if pointer.GetString(exporter.MetricsURL) != "" {
		args = append(args, "--web.telemetry-path="+*exporter.MetricsURL)
	}

	sort.Strings(args)

	host := pointer.GetString(service.Address)
	port := pointer.GetUint16(service.Port)

	usr := pointer.GetString(exporter.Username)
	passwd := pointer.GetString(exporter.Password)
	addr := net.JoinHostPort(host, strconv.Itoa(int(port)))
	connString := fmt.Sprintf("mongodb://%s:%s@%s/", usr, passwd, addr)

	return &api.SetStateRequest_AgentProcess{
		Type:               api.Type_MONGODB_EXPORTER,
		TemplateLeftDelim:  tdp.left,
		TemplateRightDelim: tdp.right,
		Args:               args,
		Env: []string{
			fmt.Sprintf("MONGODB_URI=%s", connString),
		},
	}
}
