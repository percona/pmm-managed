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
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/agentpb"
	"github.com/percona/pmm/api/inventorypb"
	"github.com/percona/pmm/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/pmm-managed/models"
)

func TestPostgresExporterConfig(t *testing.T) {

	pmmAgentVersion := version.MustParse("2.15.1")
	var postgresql *models.Service
	var exporter *models.Agent
	var expected *agentpb.SetStateRequest_AgentProcess

	setup := func() {
		postgresql = &models.Service{
			Address: pointer.ToString("1.2.3.4"),
			Port:    pointer.ToUint16(5432),
		}
		exporter = &models.Agent{
			AgentID:       "agent-id",
			AgentType:     models.PostgresExporterType,
			Username:      pointer.ToString("username"),
			Password:      pointer.ToString("s3cur3 p@$$w0r4."),
			AgentPassword: pointer.ToString("agent-password"),
		}
		expected = &agentpb.SetStateRequest_AgentProcess{
			Type:               inventorypb.AgentType_POSTGRES_EXPORTER,
			TemplateLeftDelim:  "{{",
			TemplateRightDelim: "}}",
			Args: []string{
				"--collect.custom_query.hr",
				"--collect.custom_query.hr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/high-resolution",
				"--collect.custom_query.lr",
				"--collect.custom_query.lr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/low-resolution",
				"--collect.custom_query.mr",
				"--collect.custom_query.mr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/medium-resolution",
				"--web.listen-address=:{{ .listen_port }}",
			},
			Env: []string{
				"DATA_SOURCE_NAME=postgres://username:s3cur3%20p%40$$w0r4.@1.2.3.4:5432/postgres?connect_timeout=1&sslmode=disable",
				"HTTP_AUTH=pmm:agent-password",
			},
			RedactWords: []string{"s3cur3 p@$$w0r4.", "agent-password"},
		}
	}

	run := func(name string, f func(t *testing.T)) {
		setup()
		t.Run(name, f)
	}

	run("Config", func(t *testing.T) {
		actual := postgresExporterConfig(postgresql, exporter, redactSecrets, pmmAgentVersion)
		requireNoDuplicateFlags(t, actual.Args)
		require.Equal(t, expected.Args, actual.Args)
		require.Equal(t, expected.Env, actual.Env)
		require.Equal(t, expected, actual)
	})

	run("DatabaseName", func(t *testing.T) {
		postgresql.DatabaseName = "db1"
		expected.Env[0] = "DATA_SOURCE_NAME=postgres://username:s3cur3%20p%40$$w0r4.@1.2.3.4:5432/db1?connect_timeout=1&sslmode=disable"
		actual := postgresExporterConfig(postgresql, exporter, redactSecrets, pmmAgentVersion)
		require.Equal(t, expected.Env, actual.Env)
	})

	run("EmptyPassword", func(t *testing.T) {
		exporter.Password = nil
		actual := postgresExporterConfig(postgresql, exporter, exposeSecrets, pmmAgentVersion)
		assert.Equal(t, "DATA_SOURCE_NAME=postgres://username@1.2.3.4:5432/postgres?connect_timeout=1&sslmode=disable", actual.Env[0])
	})

	run("EmptyUsername", func(t *testing.T) {
		exporter.Username = nil
		actual := postgresExporterConfig(postgresql, exporter, exposeSecrets, pmmAgentVersion)
		assert.Equal(t, "DATA_SOURCE_NAME=postgres://:s3cur3%20p%40$$w0r4.@1.2.3.4:5432/postgres?connect_timeout=1&sslmode=disable", actual.Env[0])
	})

	run("EmptyUsernameAndPassword", func(t *testing.T) {
		exporter.Username = nil
		exporter.Password = nil
		actual := postgresExporterConfig(postgresql, exporter, exposeSecrets, pmmAgentVersion)
		assert.Equal(t, "DATA_SOURCE_NAME=postgres://1.2.3.4:5432/postgres?connect_timeout=1&sslmode=disable", actual.Env[0])
	})

	run("Socket", func(t *testing.T) {
		exporter.Username = nil
		exporter.Password = nil
		postgresql.Address = nil
		postgresql.Port = nil
		postgresql.Socket = pointer.ToString("/var/run/postgres")
		actual := postgresExporterConfig(postgresql, exporter, exposeSecrets, pmmAgentVersion)
		assert.Equal(t, "DATA_SOURCE_NAME=postgres:///postgres?connect_timeout=1&host=%2Fvar%2Frun%2Fpostgres&sslmode=disable", actual.Env[0])
	})

	run("DisabledCollectors", func(t *testing.T) {
		postgresql.Address = nil
		postgresql.Port = nil
		postgresql.Socket = pointer.ToString("/var/run/postgres")
		exporter.DisabledCollectors = []string{"custom_query.hr", "custom_query.hr.directory"}
		actual := postgresExporterConfig(postgresql, exporter, exposeSecrets, pmmAgentVersion)
		expected := &agentpb.SetStateRequest_AgentProcess{
			Type:               inventorypb.AgentType_POSTGRES_EXPORTER,
			TemplateLeftDelim:  "{{",
			TemplateRightDelim: "}}",
			Args: []string{
				"--collect.custom_query.lr",
				"--collect.custom_query.lr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/low-resolution",
				"--collect.custom_query.mr",
				"--collect.custom_query.mr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/medium-resolution",
				"--web.listen-address=:{{ .listen_port }}",
			},
		}
		requireNoDuplicateFlags(t, actual.Args)
		require.Equal(t, expected.Args, actual.Args)
	})

	run("AutoDiscovery", func(t *testing.T) {
		pmmAgentVersion := version.MustParse("2.16.0")

		postgresql := &models.Service{
			Address: pointer.ToString("1.2.3.4"),
			Port:    pointer.ToUint16(5432),
		}
		exporter := &models.Agent{
			AgentID:   "agent-id",
			AgentType: models.PostgresExporterType,
			Username:  pointer.ToString("username"),
			Password:  pointer.ToString("s3cur3 p@$$w0r4."),
		}

		actual := postgresExporterConfig(postgresql, exporter, redactSecrets, pmmAgentVersion)
		expected = &agentpb.SetStateRequest_AgentProcess{
			Type:               inventorypb.AgentType_POSTGRES_EXPORTER,
			TemplateLeftDelim:  "{{",
			TemplateRightDelim: "}}",
			Args: []string{
				"--auto-discover-databases",
				"--collect.custom_query.hr",
				"--collect.custom_query.hr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/high-resolution",
				"--collect.custom_query.lr",
				"--collect.custom_query.lr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/low-resolution",
				"--collect.custom_query.mr",
				"--collect.custom_query.mr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/medium-resolution",
				"--exclude-databases=template0,template1,postgres,cloudsqladmin,pmm-managed-dev,azure_maintenance",
				"--web.listen-address=:{{ .listen_port }}",
			},
			Env: []string{
				"DATA_SOURCE_NAME=postgres://username:s3cur3%20p%40$$w0r4.@1.2.3.4:5432/postgres?connect_timeout=1&sslmode=disable",
				"HTTP_AUTH=pmm:agent-id",
			},
			RedactWords: []string{"s3cur3 p@$$w0r4."},
		}
		requireNoDuplicateFlags(t, actual.Args)
		require.Equal(t, expected.Args, actual.Args)
		require.Equal(t, expected.Env, actual.Env)
		require.Equal(t, expected, actual)
	})

	run("AzureTimeout", func(t *testing.T) {
		pmmAgentVersion := version.MustParse("2.16.0")

		postgresql := &models.Service{
			Address: pointer.ToString("1.2.3.4"),
			Port:    pointer.ToUint16(5432),
		}
		exporter := &models.Agent{
			AgentID:   "agent-id",
			AgentType: models.PostgresExporterType,
			Username:  pointer.ToString("username"),
			Password:  pointer.ToString("s3cur3 p@$$w0r4."),
			AzureOptions: &models.AzureOptions{
				SubscriptionID: "subscription_id",
				ClientID:       "client_id",
				ClientSecret:   "client_secret",
				TenantID:       "tenant_id",
				ResourceGroup:  "resource_group",
			},
		}

		actual := postgresExporterConfig(postgresql, exporter, redactSecrets, pmmAgentVersion)
		expected = &agentpb.SetStateRequest_AgentProcess{
			Type:               inventorypb.AgentType_POSTGRES_EXPORTER,
			TemplateLeftDelim:  "{{",
			TemplateRightDelim: "}}",
			Args: []string{
				"--auto-discover-databases",
				"--collect.custom_query.hr",
				"--collect.custom_query.hr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/high-resolution",
				"--collect.custom_query.lr",
				"--collect.custom_query.lr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/low-resolution",
				"--collect.custom_query.mr",
				"--collect.custom_query.mr.directory=" + pathsBase(pointer.GetString(exporter.Version), "{{", "}}") + "/collectors/custom-queries/postgresql/medium-resolution",
				"--exclude-databases=template0,template1,postgres,cloudsqladmin,pmm-managed-dev,azure_maintenance",
				"--web.listen-address=:{{ .listen_port }}",
			},
			Env: []string{
				"DATA_SOURCE_NAME=postgres://username:s3cur3%20p%40$$w0r4.@1.2.3.4:5432/postgres?connect_timeout=5&sslmode=disable",
				"HTTP_AUTH=pmm:agent-id",
			},
			RedactWords: []string{"s3cur3 p@$$w0r4.", "client_secret"},
		}
		requireNoDuplicateFlags(t, actual.Args)
		require.Equal(t, expected.Args, actual.Args)
		require.Equal(t, expected.Env, actual.Env)
		require.Equal(t, expected, actual)
	})
}
