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
	"net/url"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/percona/pmm-managed/models"
	config_util "github.com/percona/pmm-managed/services/prometheus/internal/common/config"
	"github.com/percona/pmm-managed/services/prometheus/internal/prometheus/config"
	sd_config "github.com/percona/pmm-managed/services/prometheus/internal/prometheus/discovery/config"
	"github.com/percona/pmm-managed/services/prometheus/internal/prometheus/discovery/targetgroup"
)

func TestScrapeConfig(t *testing.T) {
	s := &models.MetricsResolutions{
		HR: time.Second,
		MR: 5 * time.Second,
		LR: 60 * time.Second,
	}

	t.Run("scrapeConfigsForNodeExporter", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			node := &models.Node{
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				NodeName:     "node_name",
				Address:      "1.2.3.4",
				CustomLabels: []byte(`{"_some_node_label": "foo"}`),
			}
			agent := &models.Agent{
				AgentID:      "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
				AgentType:    models.NodeExporterType,
				RunsOnNodeID: nil,
				CustomLabels: []byte(`{"_some_agent_label": "baz"}`),
				ListenPort:   pointer.ToUint16(12345),
			}

			expected := []*config.ScrapeConfig{
				{
					JobName:        "node_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1s",
					ScrapeInterval: model.Duration(s.HR),
					ScrapeTimeout:  scrapeTimeout(s.HR),
					MetricsPath:    "/metrics",
					HTTPClientConfig: config_util.HTTPClientConfig{
						BasicAuth: &config_util.BasicAuth{
							Username: "pmm",
							Password: "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
						},
					},
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label": "baz",
								"_some_node_label":  "foo",
								"agent_id":          "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":        "node_exporter",
								"instance":          "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":           "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":         "node_name",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"diskstats",
						"filefd",
						"filesystem",
						"loadavg",
						"meminfo",
						"meminfo_numa",
						"netdev",
						"netstat",
						"standard.go",
						"standard.process",
						"stat",
						"textfile.hr",
						"time",
						"vmstat",
					}},
				},
				{
					JobName:        "node_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_5s",
					ScrapeInterval: model.Duration(s.MR),
					ScrapeTimeout:  scrapeTimeout(s.MR),
					MetricsPath:    "/metrics",
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label": "baz",
								"_some_node_label":  "foo",
								"agent_id":          "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":        "node_exporter",
								"instance":          "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":           "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":         "node_name",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"textfile.mr",
					}},
				},
				{
					JobName:        "node_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1m0s",
					ScrapeInterval: model.Duration(s.LR),
					ScrapeTimeout:  scrapeTimeout(s.LR),
					MetricsPath:    "/metrics",
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label": "baz",
								"_some_node_label":  "foo",
								"agent_id":          "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":        "node_exporter",
								"instance":          "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":           "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":         "node_name",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"bonding",
						"entropy",
						"textfile.lr",
						"uname",
					}},
				},
			}

			actual, err := scrapeConfigsForNodeExporter(s, "", node, agent)
			require.NoError(t, err)

			require.NoError(t, err)
			require.Len(t, actual, len(expected))
			for i := 0; i < len(expected); i++ {
				assertScrappedConfigsEqual(t, expected[i], actual[i])
			}
		})
	})

	t.Run("scrapeConfigsForMySQLdExporter", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			node := &models.Node{
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				NodeName:     "node_name",
				Address:      "1.2.3.4",
				CustomLabels: []byte(`{"_some_node_label": "foo"}`),
			}
			service := &models.Service{
				ServiceID:    "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				Address:      pointer.ToString("5.6.7.8"),
				CustomLabels: []byte(`{"_some_service_label": "bar"}`),
			}
			agent := &models.Agent{
				AgentID:      "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
				AgentType:    models.MySQLdExporterType,
				RunsOnNodeID: nil,
				CustomLabels: []byte(`{"_some_agent_label": "baz"}`),
				ListenPort:   pointer.ToUint16(12345),
			}

			expected := []*config.ScrapeConfig{{
				JobName:        "mysqld_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1s",
				ScrapeInterval: model.Duration(s.HR),
				ScrapeTimeout:  scrapeTimeout(s.HR),
				MetricsPath:    "/metrics",
				HTTPClientConfig: config_util.HTTPClientConfig{
					BasicAuth: &config_util.BasicAuth{
						Username: "pmm",
						Password: "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
					},
				},
				ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
					StaticConfigs: []*targetgroup.Group{{
						Targets: []model.LabelSet{{"__address__": "4.5.6.7:12345"}},
						Labels: model.LabelSet{
							"_some_agent_label":   "baz",
							"_some_node_label":    "foo",
							"_some_service_label": "bar",
							"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
							"agent_type":          "mysqld_exporter",
							"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
							"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
							"node_name":           "node_name",
							"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
						},
					}},
				},
				Params: url.Values{"collect[]": []string{
					"global_status",
					"info_schema.innodb_metrics",
					"custom_query.hr",
					"standard.process",
					"standard.go",
				}},
			}, {
				JobName:        "mysqld_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_5s",
				ScrapeInterval: model.Duration(s.MR),
				ScrapeTimeout:  scrapeTimeout(s.MR),
				MetricsPath:    "/metrics",
				ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
					StaticConfigs: []*targetgroup.Group{{
						Targets: []model.LabelSet{{"__address__": "4.5.6.7:12345"}},
						Labels: model.LabelSet{
							"_some_agent_label":   "baz",
							"_some_node_label":    "foo",
							"_some_service_label": "bar",
							"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
							"agent_type":          "mysqld_exporter",
							"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
							"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
							"node_name":           "node_name",
							"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
						},
					}},
				},
				Params: url.Values{"collect[]": []string{
					"engine_innodb_status",
					"info_schema.innodb_cmp",
					"info_schema.innodb_cmpmem",
					"info_schema.processlist",
					"info_schema.query_response_time",
					"perf_schema.eventswaits",
					"perf_schema.file_events",
					"perf_schema.tablelocks",
					"slave_status",
					"custom_query.mr",
				}},
			}, {
				JobName:        "mysqld_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1m0s",
				ScrapeInterval: model.Duration(s.LR),
				ScrapeTimeout:  scrapeTimeout(s.LR),
				MetricsPath:    "/metrics",
				ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
					StaticConfigs: []*targetgroup.Group{{
						Targets: []model.LabelSet{{"__address__": "4.5.6.7:12345"}},
						Labels: model.LabelSet{
							"_some_agent_label":   "baz",
							"_some_node_label":    "foo",
							"_some_service_label": "bar",
							"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
							"agent_type":          "mysqld_exporter",
							"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
							"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
							"node_name":           "node_name",
							"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
						},
					}},
				},
				Params: url.Values{"collect[]": []string{
					"auto_increment.columns",
					"binlog_size",
					"engine_tokudb_status",
					"global_variables",
					"heartbeat",
					"info_schema.clientstats",
					"info_schema.innodb_tablespaces",
					"info_schema.tables",
					"info_schema.tablestats",
					"info_schema.userstats",
					"perf_schema.eventsstatements",
					"perf_schema.file_instances",
					"perf_schema.indexiowaits",
					"perf_schema.tableiowaits",
					"perf_schema.tablestats",
					"custom_query.lr",
				}},
			}}

			actual, err := scrapeConfigsForMySQLdExporter(s, "", node, service, agent)
			require.NoError(t, err)
			require.Len(t, actual, len(expected))
			for i := 0; i < len(expected); i++ {
				assertScrappedConfigsEqual(t, expected[i], actual[i])
			}
		})

		t.Run("BadCustomLabels", func(t *testing.T) {
			node := &models.Node{}
			service := &models.Service{}
			agent := &models.Agent{
				CustomLabels: []byte("{"),
				ListenPort:   pointer.ToUint16(12345),
			}

			_, err := scrapeConfigsForMySQLdExporter(s, "", node, service, agent)
			require.EqualError(t, err, "failed to decode custom labels: unexpected end of JSON input")
		})
	})

	t.Run("scrapeConfigsForMongoDBExporter", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			node := &models.Node{
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				NodeName:     "node_name",
				Address:      "1.2.3.4",
				CustomLabels: []byte(`{"_some_node_label": "foo"}`),
			}
			service := &models.Service{
				ServiceID:    "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				Address:      pointer.ToString("5.6.7.8"),
				CustomLabels: []byte(`{"_some_service_label": "bar"}`),
			}
			agent := &models.Agent{
				AgentID:      "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
				AgentType:    models.MongoDBExporterType,
				RunsOnNodeID: nil,
				CustomLabels: []byte(`{"_some_agent_label": "baz"}`),
				ListenPort:   pointer.ToUint16(12345),
			}

			expected := []*config.ScrapeConfig{
				{
					JobName:        "mongodb_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1s",
					ScrapeInterval: model.Duration(s.HR),
					ScrapeTimeout:  scrapeTimeout(s.HR),
					MetricsPath:    "/metrics",
					HTTPClientConfig: config_util.HTTPClientConfig{
						BasicAuth: &config_util.BasicAuth{
							Username: "pmm",
							Password: "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
						},
					},
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label":   "baz",
								"_some_node_label":    "foo",
								"_some_service_label": "bar",
								"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":          "mongodb_exporter",
								"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":           "node_name",
								"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"collection",
						"database",
						"standard.process",
						"standard.go",
					}},
				},
			}

			actual, err := scrapeConfigsForMongoDBExporter(s, "", node, service, agent)
			require.NoError(t, err)
			require.Len(t, actual, len(expected))
			for i := 0; i < len(expected); i++ {
				assertScrappedConfigsEqual(t, expected[i], actual[i])
			}
		})

		t.Run("BadCustomLabels", func(t *testing.T) {
			node := &models.Node{}
			service := &models.Service{}
			agent := &models.Agent{
				CustomLabels: []byte("{"),
				ListenPort:   pointer.ToUint16(12345),
			}

			_, err := scrapeConfigsForMongoDBExporter(s, "", node, service, agent)
			require.EqualError(t, err, "failed to decode custom labels: unexpected end of JSON input")
		})
	})

	t.Run("scrapeConfigsForPostgresExporter", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			node := &models.Node{
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				NodeName:     "node_name",
				Address:      "1.2.3.4",
				CustomLabels: []byte(`{"_some_node_label": "foo"}`),
			}
			service := &models.Service{
				ServiceID:    "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
				NodeID:       "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
				Address:      pointer.ToString("5.6.7.8"),
				CustomLabels: []byte(`{"_some_service_label": "bar"}`),
			}
			agent := &models.Agent{
				AgentID:      "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
				AgentType:    models.PostgresExporterType,
				RunsOnNodeID: nil,
				CustomLabels: []byte(`{"_some_agent_label": "baz"}`),
				ListenPort:   pointer.ToUint16(12345),
			}

			expected := []*config.ScrapeConfig{
				{
					JobName:        "postgres_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1s",
					ScrapeInterval: model.Duration(s.HR),
					ScrapeTimeout:  scrapeTimeout(s.HR),
					MetricsPath:    "/metrics",
					HTTPClientConfig: config_util.HTTPClientConfig{
						BasicAuth: &config_util.BasicAuth{
							Username: "pmm",
							Password: "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
						},
					},
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label":   "baz",
								"_some_node_label":    "foo",
								"_some_service_label": "bar",
								"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":          "postgres_exporter",
								"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":           "node_name",
								"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"exporter",
						"custom_query.hr",
						"standard.process",
						"standard.go",
					}},
				},
				{
					JobName:        "postgres_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_5s",
					ScrapeInterval: model.Duration(s.MR),
					ScrapeTimeout:  scrapeTimeout(s.MR),
					MetricsPath:    "/metrics",
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label":   "baz",
								"_some_node_label":    "foo",
								"_some_service_label": "bar",
								"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":          "postgres_exporter",
								"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":           "node_name",
								"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"custom_query.mr",
					}},
				},
				{
					JobName:        "postgres_exporter_agent_id_75bb30d3-ef4a-4147-97a8-621a996611dd_1m0s",
					ScrapeInterval: model.Duration(s.LR),
					ScrapeTimeout:  scrapeTimeout(s.LR),
					MetricsPath:    "/metrics",
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label":   "baz",
								"_some_node_label":    "foo",
								"_some_service_label": "bar",
								"agent_id":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"agent_type":          "postgres_exporter",
								"instance":            "/agent_id/75bb30d3-ef4a-4147-97a8-621a996611dd",
								"node_id":             "/node_id/cc663f36-18ca-40a1-aea9-c6310bb4738d",
								"node_name":           "node_name",
								"service_id":          "/service_id/014647c3-b2f5-44eb-94f4-d943260a968c",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"custom_query.lr",
					}},
				},
			}

			actual, err := scrapeConfigsForPostgresExporter(s, "", node, service, agent)
			require.NoError(t, err)
			require.Len(t, actual, len(expected))
			for i := 0; i < len(expected); i++ {
				assertScrappedConfigsEqual(t, expected[i], actual[i])
			}
		})

		t.Run("BadCustomLabels", func(t *testing.T) {
			node := &models.Node{}
			service := &models.Service{}
			agent := &models.Agent{
				CustomLabels: []byte("{"),
				ListenPort:   pointer.ToUint16(12345),
			}

			_, err := scrapeConfigsForPostgresExporter(s, "", node, service, agent)
			require.EqualError(t, err, "failed to decode custom labels: unexpected end of JSON input")
		})
	})

	t.Run("scrapeConfigsForProxySQLExporter", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			node := &models.Node{
				NodeID:       "/node_id/7cc6ec12-4951-48c6-a4d5-7c3141fa4107",
				NodeName:     "node_name",
				Address:      "1.2.3.4",
				CustomLabels: []byte(`{"_some_node_label": "foo"}`),
			}
			service := &models.Service{
				ServiceID:    "/service_id/56fa3285-4476-49cc-95ff-96b36c24b0b6",
				NodeID:       "/node_id/7cc6ec12-4951-48c6-a4d5-7c3141fa4107",
				Address:      pointer.ToString("5.6.7.8"),
				CustomLabels: []byte(`{"_some_service_label": "bar"}`),
			}
			agent := &models.Agent{
				AgentID:      "/agent_id/782589c6-d3af-45e5-aa20-7f664a690940",
				AgentType:    models.ProxySQLExporterType,
				RunsOnNodeID: nil,
				CustomLabels: []byte(`{"_some_agent_label": "baz"}`),
				ListenPort:   pointer.ToUint16(12345),
			}

			expected := []*config.ScrapeConfig{
				{
					JobName:        "proxysql_exporter_agent_id_782589c6-d3af-45e5-aa20-7f664a690940_1s",
					ScrapeInterval: model.Duration(s.HR),
					ScrapeTimeout:  scrapeTimeout(s.HR),
					MetricsPath:    "/metrics",
					HTTPClientConfig: config_util.HTTPClientConfig{
						BasicAuth: &config_util.BasicAuth{
							Username: "pmm",
							Password: "/agent_id/782589c6-d3af-45e5-aa20-7f664a690940",
						},
					},
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						StaticConfigs: []*targetgroup.Group{{
							Targets: []model.LabelSet{{"__address__": "1.2.3.4:12345"}},
							Labels: model.LabelSet{
								"_some_agent_label":   "baz",
								"_some_node_label":    "foo",
								"_some_service_label": "bar",
								"agent_id":            "/agent_id/782589c6-d3af-45e5-aa20-7f664a690940",
								"agent_type":          "proxysql_exporter",
								"instance":            "/agent_id/782589c6-d3af-45e5-aa20-7f664a690940",
								"node_id":             "/node_id/7cc6ec12-4951-48c6-a4d5-7c3141fa4107",
								"node_name":           "node_name",
								"service_id":          "/service_id/56fa3285-4476-49cc-95ff-96b36c24b0b6",
							},
						}},
					},
					Params: url.Values{"collect[]": []string{
						"mysql_connection_pool",
						"mysql_status",
						"standard.process",
						"standard.go",
					}},
				},
			}

			actual, err := scrapeConfigsForProxySQLExporter(s, "", node, service, agent)
			require.NoError(t, err)
			require.Len(t, actual, len(expected))
			for i := 0; i < len(expected); i++ {
				assertScrappedConfigsEqual(t, expected[i], actual[i])
			}
		})

		t.Run("BadCustomLabels", func(t *testing.T) {
			node := &models.Node{}
			service := &models.Service{}
			agent := &models.Agent{
				CustomLabels: []byte("{"),
				ListenPort:   pointer.ToUint16(12345),
			}

			_, err := scrapeConfigsForProxySQLExporter(s, "", node, service, agent)
			require.EqualError(t, err, "failed to decode custom labels: unexpected end of JSON input")
		})
	})
}

func assertScrappedConfigsEqual(t *testing.T, expected, actual *config.ScrapeConfig) bool {
	if !assert.Equal(t, expected, actual) {
		e, err := yaml.Marshal(expected)
		require.NoError(t, err)
		a, err := yaml.Marshal(actual)
		require.NoError(t, err)

		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(e)),
			FromFile: "Expected",
			B:        difflib.SplitLines(string(a)),
			ToFile:   "Actual",
			Context:  3,
		})
		require.NoError(t, err)
		t.Logf("Diff:\n%s", diff)
		return false
	}
	return true
}
