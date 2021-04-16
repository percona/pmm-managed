package server

import (
	"archive/zip"
	"bytes"
	"os"
	"sort"
	"testing"

	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestDownloadLogs(t *testing.T) {
	var buf bytes.Buffer
	res, err := serverClient.Default.Server.Logs(&server.LogsParams{
		Context: pmmapitests.Context,
	}, &buf)
	require.NoError(t, err)
	require.NotNil(t, res)

	r := bytes.NewReader(buf.Bytes())
	zipR, err := zip.NewReader(r, r.Size())
	assert.NoError(t, err)

	expected := []string{
		"alertmanager.base.yml",
		"alertmanager.ini",
		"alertmanager.log",
		"alertmanager.yml",
		"clickhouse-server.err.log",
		"clickhouse-server.log",
		"clickhouse-server.startup.log",
		"client/list.txt",
		"client/pmm-admin-version.txt",
		"client/pmm-agent-config.yaml",
		"client/pmm-agent-version.txt",
		"client/status.json",
		"cron.log",
		"dashboard-upgrade.log",
		"grafana.log",
		"installed.json",
		"nginx.access.log",
		"nginx.conf",
		"nginx.error.log",
		"nginx.startup.log",
		"pmm-agent.log",
		"pmm-agent.yaml",
		"pmm-managed.log",
		"pmm-ssl.conf",
		"pmm-version.txt",
		"pmm.conf",
		"pmm.ini",
		"postgresql.log",
		"postgresql.startup.log",
		"prometheus.base.yml",
		"prometheus.ini",
		"prometheus.yml",
		"qan-api2.ini",
		"qan-api2.log",
		"supervisorctl_status.log",
		"supervisord.conf",
		"supervisord.log",
		"systemctl_status.log",
		"victoriametrics-promscrape.yml",
		"victoriametrics.ini",
		"victoriametrics.log",
		"victoriametrics_targets.json",
		"vmalert.ini",
		"vmalert.log",
	}

	if os.Getenv("PERCONA_TEST_DBAAS") == "1" {
		expected = append(expected, "dbaas-controller.log")
		sort.Strings(expected)
	}

	actual := make([]string, len(zipR.File))
	for i, file := range zipR.File {
		actual[i] = file.Name
	}

	sort.Strings(actual)
	assert.Equal(t, expected, actual)
}
