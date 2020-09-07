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

package prometheus

import (
	"fmt"
	"net"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	config "github.com/percona/promconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// ScrapeTimeout returns default scrape timeout for given scrape interval.
func ScrapeTimeout(interval time.Duration) config.Duration {
	switch {
	case interval <= 2*time.Second:
		return config.Duration(time.Second)
	case interval <= 10*time.Second:
		return config.Duration(interval - time.Second)
	default:
		return config.Duration(10 * time.Second)
	}
}

func scrapeConfigForPrometheus(interval time.Duration) *config.ScrapeConfig {
	return &config.ScrapeConfig{
		JobName:        "prometheus",
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		MetricsPath:    "/prometheus/metrics",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			StaticConfigs: []*config.Group{{
				Targets: []string{"127.0.0.1:9090"},
				Labels:  map[string]string{"instance": "pmm-server"},
			}},
		},
	}
}

func scrapeConfigForAlertmanager(interval time.Duration) *config.ScrapeConfig {
	return &config.ScrapeConfig{
		JobName:        "alertmanager",
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		MetricsPath:    "/alertmanager/metrics",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			StaticConfigs: []*config.Group{{
				Targets: []string{"127.0.0.1:9093"},
				Labels:  map[string]string{"instance": "pmm-server"},
			}},
		},
	}
}

func scrapeConfigForGrafana(interval time.Duration) *config.ScrapeConfig {
	return &config.ScrapeConfig{
		JobName:        "grafana",
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		MetricsPath:    "/metrics",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			StaticConfigs: []*config.Group{{
				Targets: []string{"127.0.0.1:3000"},
				Labels:  map[string]string{"instance": "pmm-server"},
			}},
		},
	}
}

func scrapeConfigForPMMManaged(interval time.Duration) *config.ScrapeConfig {
	return &config.ScrapeConfig{
		JobName:        "pmm-managed",
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		MetricsPath:    "/debug/metrics",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			StaticConfigs: []*config.Group{{
				Targets: []string{"127.0.0.1:7773"},
				Labels:  map[string]string{"instance": "pmm-server"},
			}},
		},
	}
}

func scrapeConfigForQANAPI2(interval time.Duration) *config.ScrapeConfig {
	return &config.ScrapeConfig{
		JobName:        "qan-api2",
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		MetricsPath:    "/debug/metrics",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			StaticConfigs: []*config.Group{{
				Targets: []string{"127.0.0.1:9933"},
				Labels:  map[string]string{"instance": "pmm-server"},
			}},
		},
	}
}

func mergeLabels(node *models.Node, service *models.Service, agent *models.Agent) (map[string]string, error) {
	res, err := models.MergeLabels(node, service, agent)
	if err != nil {
		return nil, err
	}

	res["instance"] = agent.AgentID

	return res, nil
}

// jobNameMapping replaces runes that can't be present in Prometheus job name with '_'.
func jobNameMapping(r rune) rune {
	switch r {
	case '/', ':', '.':
		return '_'
	default:
		return r
	}
}

func jobName(agent *models.Agent, intervalName string, interval time.Duration) string {
	return fmt.Sprintf("%s%s_%s-%s", agent.AgentType, strings.Map(jobNameMapping, agent.AgentID), intervalName, interval)
}

func httpClientConfig(agent *models.Agent) config.HTTPClientConfig {
	return config.HTTPClientConfig{
		BasicAuth: &config.BasicAuth{
			Username: "pmm",
			Password: agent.AgentID,
		},
	}
}

type scrapeConfigParams struct {
	host    string // Node address where pmm-agent runs
	node    *models.Node
	service *models.Service
	agent   *models.Agent
}

// scrapeConfigForStandardExporter returns scrape config for endpoint with given parameters.
func scrapeConfigForStandardExporter(intervalName string, interval time.Duration, params *scrapeConfigParams, collect []string) (*config.ScrapeConfig, error) {
	labels, err := mergeLabels(params.node, params.service, params.agent)
	if err != nil {
		return nil, err
	}

	cfg := &config.ScrapeConfig{
		JobName:          jobName(params.agent, intervalName, interval),
		ScrapeInterval:   config.Duration(interval),
		ScrapeTimeout:    ScrapeTimeout(interval),
		MetricsPath:      "/metrics",
		HTTPClientConfig: httpClientConfig(params.agent),
	}

	if len(collect) > 0 {
		sort.Strings(collect)
		cfg.Params = url.Values{
			"collect[]": collect,
		}
	}

	port := int(*params.agent.ListenPort)
	hostport := net.JoinHostPort(params.host, strconv.Itoa(port))

	cfg.ServiceDiscoveryConfig = config.ServiceDiscoveryConfig{
		StaticConfigs: []*config.Group{{
			Targets: []string{hostport},
			Labels:  labels,
		}},
	}

	return cfg, nil
}

// scrapeConfigForRDSExporter returns scrape config for single rds_exporter configuration.
func scrapeConfigForRDSExporter(intervalName string, interval time.Duration, hostport string, metricsPath string) *config.ScrapeConfig {
	jobName := fmt.Sprintf("rds_exporter_%s_%s-%s", strings.Map(jobNameMapping, hostport), intervalName, interval)
	return &config.ScrapeConfig{
		JobName:        jobName,
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		MetricsPath:    metricsPath,
		HonorLabels:    true,
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			StaticConfigs: []*config.Group{{
				Targets: []string{hostport},
			}},
		},
	}
}

func scrapeConfigsForNodeExporter(s *models.MetricsResolutions, params *scrapeConfigParams) ([]*config.ScrapeConfig, error) {
	var hr, mr, lr *config.ScrapeConfig
	var err error
	var hrCollect []string

	if params.node.Distro != "darwin" {
		mr, err = scrapeConfigForStandardExporter("mr", s.MR, params, []string{
			"hwmon",
			"textfile.mr",
		})
		if err != nil {
			return nil, err
		}

		lr, err = scrapeConfigForStandardExporter("lr", s.LR, params, []string{
			"bonding",
			"entropy",
			"textfile.lr",
			"uname",
		})
		if err != nil {
			return nil, err
		}

		hrCollect = append(hrCollect,
			"buddyinfo",
			"filefd",
			"meminfo_numa",
			"netstat",
			"processes",
			"standard.go",
			"standard.process",
			"stat",
			"textfile.hr",
			"vmstat",
		)
	}

	hr, err = scrapeConfigForStandardExporter("hr", s.HR, params, append(hrCollect,
		"cpu",
		"diskstats",
		"filesystem",
		"loadavg",
		"meminfo",
		"netdev",
		"time",
	))
	if err != nil {
		return nil, err
	}

	var r []*config.ScrapeConfig
	if hr != nil {
		r = append(r, hr)
	}
	if mr != nil {
		r = append(r, mr)
	}
	if lr != nil {
		r = append(r, lr)
	}
	return r, nil
}

// scrapeConfigsForMySQLdExporter returns scrape config for mysqld_exporter.
func scrapeConfigsForMySQLdExporter(s *models.MetricsResolutions, params *scrapeConfigParams) ([]*config.ScrapeConfig, error) {
	// keep in sync with mysqld_exporter Agent flags generator

	hr, err := scrapeConfigForStandardExporter("hr", s.HR, params, []string{
		"global_status",
		"info_schema.innodb_metrics",
		"custom_query.hr",
		"standard.go",
		"standard.process",
	})
	if err != nil {
		return nil, err
	}

	mrOptions := []string{
		"engine_innodb_status",
		"info_schema.innodb_cmp",
		"info_schema.innodb_cmpmem",
		"info_schema.processlist",
		"info_schema.query_response_time",
		"perf_schema.eventswaits",
		"perf_schema.file_events",
		"slave_status",
		"custom_query.mr",
	}
	if params.agent.IsMySQLTablestatsGroupEnabled() {
		mrOptions = append(mrOptions, "perf_schema.tablelocks")
	}

	mr, err := scrapeConfigForStandardExporter("mr", s.MR, params, mrOptions)
	if err != nil {
		return nil, err
	}

	lrOptions := []string{
		"binlog_size",
		"engine_tokudb_status",
		"global_variables",
		"heartbeat",
		"info_schema.clientstats",
		"info_schema.innodb_tablespaces",
		"info_schema.userstats",
		"perf_schema.eventsstatements",
		"perf_schema.file_instances",
		"custom_query.lr",
	}
	if params.agent.IsMySQLTablestatsGroupEnabled() {
		lrOptions = append(lrOptions,
			"auto_increment.columns",
			"info_schema.tables",
			"info_schema.tablestats",
			"perf_schema.indexiowaits",
			"perf_schema.tableiowaits",
		)
	}

	lr, err := scrapeConfigForStandardExporter("lr", s.LR, params, lrOptions)
	if err != nil {
		return nil, err
	}

	var r []*config.ScrapeConfig
	if hr != nil {
		r = append(r, hr)
	}
	if mr != nil {
		r = append(r, mr)
	}
	if lr != nil {
		r = append(r, lr)
	}
	return r, nil
}

func scrapeConfigsForMongoDBExporter(s *models.MetricsResolutions, params *scrapeConfigParams) ([]*config.ScrapeConfig, error) {
	hr, err := scrapeConfigForStandardExporter("hr", s.HR, params, nil)
	if err != nil {
		return nil, err
	}

	var r []*config.ScrapeConfig
	if hr != nil {
		r = append(r, hr)
	}
	return r, nil
}

func scrapeConfigsForPostgresExporter(s *models.MetricsResolutions, params *scrapeConfigParams) ([]*config.ScrapeConfig, error) {
	hr, err := scrapeConfigForStandardExporter("hr", s.HR, params, []string{
		"exporter",
		"custom_query.hr",
		"standard.go",
		"standard.process",
	})
	if err != nil {
		return nil, err
	}

	mr, err := scrapeConfigForStandardExporter("mr", s.MR, params, []string{
		"custom_query.mr",
	})
	if err != nil {
		return nil, err
	}

	lr, err := scrapeConfigForStandardExporter("lr", s.LR, params, []string{
		"custom_query.lr",
	})
	if err != nil {
		return nil, err
	}

	var r []*config.ScrapeConfig
	if hr != nil {
		r = append(r, hr)
	}
	if mr != nil {
		r = append(r, mr)
	}
	if lr != nil {
		r = append(r, lr)
	}
	return r, nil
}

func scrapeConfigsForProxySQLExporter(s *models.MetricsResolutions, params *scrapeConfigParams) ([]*config.ScrapeConfig, error) {
	hr, err := scrapeConfigForStandardExporter("hr", s.HR, params, nil) // TODO https://jira.percona.com/browse/PMM-4619
	if err != nil {
		return nil, err
	}

	var r []*config.ScrapeConfig
	if hr != nil {
		r = append(r, hr)
	}
	return r, nil
}

func scrapeConfigsForRDSExporter(s *models.MetricsResolutions, params []*scrapeConfigParams) []*config.ScrapeConfig {
	hostportSet := make(map[string]struct{}, len(params))
	for _, p := range params {
		port := int(*p.agent.ListenPort)
		hostport := net.JoinHostPort(p.host, strconv.Itoa(port))
		hostportSet[hostport] = struct{}{}
	}

	hostports := make([]string, 0, len(hostportSet))
	for hostport := range hostportSet {
		hostports = append(hostports, hostport)
	}
	sort.Strings(hostports)

	r := make([]*config.ScrapeConfig, 0, len(hostports)*2)
	for _, hostport := range hostports {
		mr := scrapeConfigForRDSExporter("mr", s.MR, hostport, "/enhanced")
		lr := scrapeConfigForRDSExporter("lr", s.LR, hostport, "/basic")
		r = append(r, mr, lr)
	}

	return r
}

func scrapeConfigsForExternalExporter(s *models.MetricsResolutions, params *scrapeConfigParams) ([]*config.ScrapeConfig, error) {
	labels, err := mergeLabels(params.node, params.service, params.agent)
	if err != nil {
		return nil, err
	}

	interval := s.MR
	cfg := &config.ScrapeConfig{
		JobName:        jobName(params.agent, "mr", interval),
		ScrapeInterval: config.Duration(interval),
		ScrapeTimeout:  ScrapeTimeout(interval),
		Scheme:         pointer.GetString(params.agent.MetricsScheme),
		MetricsPath:    pointer.GetString(params.agent.MetricsPath),
	}

	if pointer.GetString(params.agent.Username) != "" {
		cfg.HTTPClientConfig = config.HTTPClientConfig{
			BasicAuth: &config.BasicAuth{
				Username: pointer.GetString(params.agent.Username),
				Password: pointer.GetString(params.agent.Password),
			},
		}
	}

	port := int(*params.agent.ListenPort)
	hostport := net.JoinHostPort(params.host, strconv.Itoa(port))

	cfg.ServiceDiscoveryConfig = config.ServiceDiscoveryConfig{
		StaticConfigs: []*config.Group{{
			Targets: []string{hostport},
			Labels:  labels,
		}},
	}

	return []*config.ScrapeConfig{cfg}, nil
}

// PopulateScrapeConfigs populates scrape configs for basic running service and agents.
func PopulateScrapeConfigs(cfg *config.Config, l *logrus.Entry, q *reform.Querier, s *models.MetricsResolutions) error {
	agentConfigs, err := agentScrapeConfigs(l, q, s)
	if err != nil {
		return err
	}
	cfg.ScrapeConfigs = append(cfg.ScrapeConfigs,
		scrapeConfigForAlertmanager(s.MR),
		scrapeConfigForGrafana(s.MR),
		scrapeConfigForPMMManaged(s.MR),
		scrapeConfigForQANAPI2(s.MR),
	)
	cfg.ScrapeConfigs = append(cfg.ScrapeConfigs, agentConfigs...)

	return nil
}

// agentScrapeConfigs generates Prometheus scrape configs for all Agents.
func agentScrapeConfigs(l *logrus.Entry, q *reform.Querier, s *models.MetricsResolutions) ([]*config.ScrapeConfig, error) {
	agents, err := q.SelectAllFrom(models.AgentTable, "WHERE NOT disabled AND listen_port IS NOT NULL ORDER BY agent_type, agent_id")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	configs := make([]*config.ScrapeConfig, 0, len(agents))
	var rdsParams []*scrapeConfigParams
	for _, str := range agents {
		agent := str.(*models.Agent)

		if agent.AgentType == models.PMMAgentType {
			// TODO https://jira.percona.com/browse/PMM-4087
			continue
		}

		// sanity check
		if (agent.NodeID != nil) && (agent.ServiceID != nil) {
			l.Panicf("Both agent.NodeID and agent.ServiceID are present: %s", agent)
		}

		// find Service for this Agent
		var paramsService *models.Service
		if agent.ServiceID != nil {
			paramsService, err = models.FindServiceByID(q, pointer.GetString(agent.ServiceID))
			if err != nil {
				return nil, err
			}
		}

		// find Node for this Agent or Service
		var paramsNode *models.Node
		switch {
		case agent.NodeID != nil:
			paramsNode, err = models.FindNodeByID(q, pointer.GetString(agent.NodeID))
		case paramsService != nil:
			paramsNode, err = models.FindNodeByID(q, paramsService.NodeID)
		}
		if err != nil {
			return nil, err
		}

		// find Node address where the agent runs
		var paramsHost string
		switch {
		case agent.PMMAgentID != nil:
			// extract node address through pmm-agent
			pmmAgent, err := models.FindAgentByID(q, *agent.PMMAgentID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pmmAgentNode := &models.Node{NodeID: pointer.GetString(pmmAgent.RunsOnNodeID)}
			if err = q.Reload(pmmAgentNode); err != nil {
				return nil, errors.WithStack(err)
			}
			paramsHost = pmmAgentNode.Address
		case agent.RunsOnNodeID != nil:
			externalExporterNode := &models.Node{NodeID: pointer.GetString(agent.RunsOnNodeID)}
			if err = q.Reload(externalExporterNode); err != nil {
				return nil, errors.WithStack(err)
			}
			paramsHost = externalExporterNode.Address
		default:
			l.Warnf("It's not possible to get host, skipping scrape config for %s.", agent)

			continue
		}

		var scfgs []*config.ScrapeConfig
		switch agent.AgentType {
		case models.NodeExporterType:
			scfgs, err = scrapeConfigsForNodeExporter(s, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: nil,
				agent:   agent,
			})

		case models.MySQLdExporterType:
			scfgs, err = scrapeConfigsForMySQLdExporter(s, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: paramsService,
				agent:   agent,
			})

		case models.MongoDBExporterType:
			scfgs, err = scrapeConfigsForMongoDBExporter(s, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: paramsService,
				agent:   agent,
			})

		case models.PostgresExporterType:
			scfgs, err = scrapeConfigsForPostgresExporter(s, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: paramsService,
				agent:   agent,
			})

		case models.ProxySQLExporterType:
			scfgs, err = scrapeConfigsForProxySQLExporter(s, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: paramsService,
				agent:   agent,
			})

		case models.QANMySQLPerfSchemaAgentType, models.QANMySQLSlowlogAgentType:
			continue
		case models.QANMongoDBProfilerAgentType:
			continue
		case models.QANPostgreSQLPgStatementsAgentType:
			continue

		case models.RDSExporterType:
			rdsParams = append(rdsParams, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: paramsService,
				agent:   agent,
			})
			continue

		case models.ExternalExporterType:
			scfgs, err = scrapeConfigsForExternalExporter(s, &scrapeConfigParams{
				host:    paramsHost,
				node:    paramsNode,
				service: paramsService,
				agent:   agent,
			})

		default:
			l.Warnf("Skipping scrape config for %s.", agent)
			continue
		}

		if err != nil {
			l.Warnf("Failed to add %s %q, skipping: %s.", agent.AgentType, agent.AgentID, err)
		}
		configs = append(configs, scfgs...)
	}

	return append(configs, scrapeConfigsForRDSExporter(s, rdsParams)...), nil
}
