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

package supervisord

import (
	"archive/zip"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// see devcontainer_test.go for more logs tests

func TestReadLog(t *testing.T) {
	f, err := ioutil.TempFile("", "pmm-managed-supervisord-tests-")
	require.NoError(t, err)
	for i := 0; i < 10; i++ {
		fmt.Fprintf(f, "line #%03d\n", i) // 10 bytes
	}
	require.NoError(t, f.Close())
	defer os.Remove(f.Name()) //nolint:errcheck

	t.Run("LimitByLines", func(t *testing.T) {
		b, m, err := readLog(f.Name(), 5, 500)
		require.NoError(t, err)
		assert.WithinDuration(t, time.Now(), m, 5*time.Second)
		expected := []string{"line #005", "line #006", "line #007", "line #008", "line #009"}
		actual := strings.Split(strings.TrimSpace(string(b)), "\n")
		assert.Equal(t, expected, actual)
	})

	t.Run("LimitByBytes", func(t *testing.T) {
		b, m, err := readLog(f.Name(), 500, 5)
		require.NoError(t, err)
		assert.WithinDuration(t, time.Now(), m, 5*time.Second)
		expected := []string{"#009"}
		actual := strings.Split(strings.TrimSpace(string(b)), "\n")
		assert.Equal(t, expected, actual)
	})
}

func TestZipFunctions(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "pmm-admin-summary")
	assert.NoError(t, err)

	// create a "client" dir to emulate the output of pmm-admin-summary
	clientDir := filepath.Join(tmpDir, "client")
	assert.NoError(t, os.Mkdir(clientDir, os.ModePerm))

	for i := 0; i < 3; i++ {
		tmpfile, err := ioutil.TempFile(clientDir, "*-test.txt")
		assert.NoError(t, err)
		assert.NoError(t, tmpfile.Close())

		buf := []byte(strings.Repeat(fmt.Sprintf("%d", i), 10))

		err = ioutil.WriteFile(tmpfile.Name(), buf, os.ModePerm)
		assert.NoError(t, err)
	}

	zipfile, err := ioutil.TempFile("", "*-test.zip")
	assert.NoError(t, err)
	zw := zip.NewWriter(zipfile)

	err = addToZip(tmpDir, zw)
	assert.NoError(t, err)
	assert.NoError(t, zw.Close())

	outTmpDir, err := ioutil.TempDir("", "pmm-admin-summary")
	assert.NoError(t, err)
	err = unzip(zipfile.Name(), outTmpDir)
	assert.NoError(t, err)
}

func TestAddAdminSummary(t *testing.T) {
	zipfile, err := ioutil.TempFile("", "*-test.zip")
	assert.NoError(t, err)

	zw := zip.NewWriter(zipfile)
	err = addAdminSummary(context.Background(), zw)
	assert.NoError(t, err)

	assert.NoError(t, zw.Close())

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
