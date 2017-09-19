// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package telemetry

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type file struct {
	path    string
	content []byte
}

func newFile() os.FileInfo {
	return &file{}
}

func (f *file) Name() string {
	return f.path
}
func (f *file) Size() int64 {
	return int64(len(f.content))
}
func (f *file) IsDir() bool {
	return false
}
func (f *file) Sys() interface{} {
	return ""
}
func (f *file) ModTime() time.Time {
	return time.Now()
}
func (f *file) Mode() os.FileMode {
	return os.ModePerm
}

func TestService(t *testing.T) {
	var count int
	var lastHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, fmt.Sprintf("cannot decode body: %s", err.Error()))
			return
		}
		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if xHeader, ok := r.Header["X-Percona-Toolkit-Tool"]; ok {
			if len(xHeader) > 0 {
				lastHeader = xHeader[0]
			}
		}
		count++
	}))
	defer ts.Close()

	uuid, _ := generateUUID()
	cfg := &Config{
		URL:      ts.URL,
		Interval: 1 * time.Second,
		UUID:     uuid,
	}

	ctx, cancel := context.WithCancel(context.Background())

	service, err := NewService(cfg)
	require.NoError(t, err)
	service.Start(ctx)
	isRunning := service.IsRunning()
	assert.Equal(t, isRunning, true)

	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, count, 1)
	cancel()
	service.Wait()
	isRunning = service.IsRunning()
	assert.Equal(t, isRunning, false)
	assert.Equal(t, lastHeader, "pmm")

	ctx, cancel = context.WithCancel(context.Background())
	// Test a service restart
	service.Start(ctx)
	isRunning = service.IsRunning()
	assert.Equal(t, isRunning, true)

	time.Sleep(2100 * time.Millisecond)
	assert.Equal(t, count, 3)
	cancel()
	service.Wait()
	isRunning = service.IsRunning()
	assert.Equal(t, isRunning, false)
}

func TestServiceIntegration(t *testing.T) {
	integrationTests := os.Getenv("INTEGRATION_TESTS")
	if integrationTests == "" {
		t.Skipf("Env var INTEGRATION_TESTS is not set. Skipping integration test")
	}

	// Using this env var for compatibility with the Toolkit
	telemetryEnvURL := os.Getenv("PERCONA_VERSION_CHECK_URL")
	if telemetryEnvURL == "" {
		t.Skipf("Env var PERCONA_VERSION_CHECK_URL is not set. Skipping integration test")
	}
	uuid, _ := generateUUID()
	cfg := &Config{
		URL:      telemetryEnvURL,
		Interval: 1,
		UUID:     uuid,
	}

	ctx := context.Background()

	service, err := NewService(cfg)
	require.NoError(t, err)
	service.Start(ctx)
	isRunning := service.IsRunning()
	assert.Equal(t, isRunning, true)

	time.Sleep(2100 * time.Millisecond)
	assert.Contains(t, service.GetLastStatus(), "telemetry data sent")
}

func TestCollectData(t *testing.T) {
	config := &Config{}
	svc, err := NewService(config)
	require.NoError(t, err)

	m, err := svc.collectData()
	require.NoError(t, err)
	assert.NotEmpty(t, m)

	assert.Contains(t, m, "OS")
	assert.Contains(t, m, "PMM")
}

func TestMakePayload(t *testing.T) {
	config := &Config{
		UUID: "ABCDEFG12345",
	}
	svc, err := NewService(config)
	require.NoError(t, err)

	m := map[string]interface{}{
		"OS":  "Kubuntu",
		"pmm": "1.2.3",
	}

	b, err := svc.makePayload(m)
	require.NoError(t, err)
	// Don't remove \n at the end of the strings. They are needed by the API
	// so I want to ensure makePayload adds them
	assert.Contains(t, string(b), "ABCDEFG12345;OS;Kubuntu\n")
	assert.Contains(t, string(b), "ABCDEFG12345;pmm;1.2.3\n")
}

// freedesktop.org and systemd
func TestGetOSNameAndVersion1(t *testing.T) {
	stat = func(filename string) (os.FileInfo, error) {
		var fs file
		return &fs, nil
	}
	readFile = func(filename string) ([]byte, error) {
		return []byte("NAME=CarlOs\nVERSION=1.0"), nil
	}

	osInfo, err := getOSNameAndVersion()
	require.NoError(t, err)
	assert.Equal(t, osInfo, "CarlOs 1.0")

	// Restore original funcs
	stat = os.Stat
	readFile = ioutil.ReadFile
}

// linuxbase.org
func TestGetOSNameAndVersion2(t *testing.T) {
	stat = func(filename string) (os.FileInfo, error) {
		return nil, fmt.Errorf("fake error")
	}
	readFile = func(filename string) ([]byte, error) {
		return []byte(""), nil
	}

	output = func(args ...string) ([]byte, error) {
		if len(args) == 2 {
			if args[1] == "-si" {
				return []byte("CarlOs"), nil
			}
			if args[1] == "-sr" {
				return []byte("version 2.0"), nil
			}
		}
		return nil, fmt.Errorf("invalid parameters")
	}

	osInfo, err := getOSNameAndVersion()
	require.NoError(t, err)
	assert.Equal(t, osInfo, "CarlOs version 2.0")

	// Restore original funcs
	stat = os.Stat
	readFile = ioutil.ReadFile
}

// For some versions of Debian/Ubuntu without lsb_release command
func TestGetOSNameAndVersion3(t *testing.T) {
	stat = func(filename string) (os.FileInfo, error) {
		if filename == "/etc/lsb-release" {
			return &file{}, nil
		}
		return nil, fmt.Errorf("fake error")
	}
	readFile = func(filename string) ([]byte, error) {
		return []byte("DISTRIB_ID=\"CarlOs\"\nDISTRIB_RELEASE=\"version 3.0\""), nil
	}

	output = func(args ...string) ([]byte, error) {
		return nil, fmt.Errorf("invalid parameters")
	}

	osInfo, err := getOSNameAndVersion()
	require.NoError(t, err)
	assert.Equal(t, osInfo, "CarlOs version 3.0")

	// Restore original funcs
	stat = os.Stat
	readFile = ioutil.ReadFile
}

// Older Debian/Ubuntu/etc.
func TestGetOSNameAndVersion4(t *testing.T) {
	stat = func(filename string) (os.FileInfo, error) {
		if filename == "/etc/debian_version" {
			return &file{}, nil
		}
		return nil, fmt.Errorf("fake error")
	}
	readFile = func(filename string) ([]byte, error) {
		return []byte("version 4.0"), nil
	}

	output = func(args ...string) ([]byte, error) {
		return nil, fmt.Errorf("invalid parameters")
	}

	osInfo, err := getOSNameAndVersion()
	require.NoError(t, err)
	assert.Equal(t, osInfo, "Debian version 4.0")

	// Restore original funcs
	stat = os.Stat
	readFile = ioutil.ReadFile
}
