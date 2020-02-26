package server

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadLogs(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	logs, err := serverClient.Default.Server.Logs(&server.LogsParams{
		Context: context.TODO(),
	}, buffer)
	require.NoError(t, err)
	require.NotNil(t, logs)

	zipfile, err := ioutil.TempFile("", "*-test.zip")
	assert.NoError(t, err)

	defer zipfile.Close() //nolint:errcheck

	_, err = io.Copy(zipfile, buffer)
	require.NoError(t, err)

	reader, err := zip.OpenReader(zipfile.Name())
	assert.NoError(t, err)

	hasClientDir := false

	for _, file := range reader.File {
		if filepath.Dir(file.Name) == "client" {
			hasClientDir = true
			break
		}
	}

	assert.True(t, hasClientDir)
}
