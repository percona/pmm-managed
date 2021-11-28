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
	"time"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/version"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/collectors"
)

var (
	// New MongoDB Exporter will be released with PMM agent v2.10.0.
	newMongoExporterPMMVersion = version.MustParse("2.9.99")
	v2_25                      = version.MustParse("2.25.0")
)

// mongodbExporterConfig returns desired configuration of mongodb_exporter process.
func mongodbExporterConfig(service *models.Service, exporter *models.Agent, redactMode redactMode,
	pmmAgentVersion *version.Parsed) *agentpb.SetStateRequest_AgentProcess {
	tdp := exporter.TemplateDelimiters(service)

	var args []string
	// Starting with PMM 2.10.0, we are shipping the new mongodb_exporter
	switch {
	case !pmmAgentVersion.Less(v2_25): // >= 2.25
		args = []string{
			"--mongodb.global-conn-pool",
			"--compatible-mode",
			"--web.listen-address=:" + tdp.Left + " .listen_port " + tdp.Right,
			"--discovering-mode",
			"--collector.dbstats",
			"--enable.top",
		}
	case !pmmAgentVersion.Less(newMongoExporterPMMVersion): // >= 2.10
		args = []string{
			"--mongodb.global-conn-pool",
			"--compatible-mode",
			"--web.listen-address=:" + tdp.Left + " .listen_port " + tdp.Right,
		}
	default:
		args = []string{
			"--collect.collection",
			"--collect.database",
			"--collect.topmetrics",
			"--no-collect.connpoolstats",
			"--no-collect.indexusage",
			"--web.listen-address=:" + tdp.Left + " .listen_port " + tdp.Right,
		}
	}

	args = collectors.FilterOutCollectors("--collect.", args, exporter.DisabledCollectors)

	if pointer.GetString(exporter.MetricsPath) != "" {
		args = append(args, "--web.telemetry-path="+*exporter.MetricsPath)
	}

	sort.Strings(args)

	database := ""
	if exporter.MongoDBOptions != nil {
		database = exporter.MongoDBOptions.AuthenticationDatabase
	}
	env := []string{
		fmt.Sprintf("MONGODB_URI=%s", exporter.DSN(service, time.Second, database, tdp)),
		fmt.Sprintf("HTTP_AUTH=pmm:%s", exporter.GetAgentPassword()),
	}

	res := &agentpb.SetStateRequest_AgentProcess{
		Type:               inventorypb.AgentType_MONGODB_EXPORTER,
		TemplateLeftDelim:  tdp.Left,
		TemplateRightDelim: tdp.Right,
		Args:               args,
		Env:                env,
		TextFiles:          exporter.Files(),
	}

	if redactMode != exposeSecrets {
		res.RedactWords = redactWords(exporter)
	}
	return res
}

// qanMongoDBProfilerAgentConfig returns desired configuration of qan-mongodb-profiler-agent built-in agent.
func qanMongoDBProfilerAgentConfig(service *models.Service, agent *models.Agent) *agentpb.SetStateRequest_BuiltinAgent {
	tdp := agent.TemplateDelimiters(service)
	return &agentpb.SetStateRequest_BuiltinAgent{
		Type:                 inventorypb.AgentType_QAN_MONGODB_PROFILER_AGENT,
		Dsn:                  agent.DSN(service, time.Second, "", nil),
		DisableQueryExamples: agent.QueryExamplesDisabled,
		TextFiles: &agentpb.TextFiles{
			Files:              agent.Files(),
			TemplateLeftDelim:  tdp.Left,
			TemplateRightDelim: tdp.Right,
		},
	}
}
