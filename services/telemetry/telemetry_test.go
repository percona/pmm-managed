package telemetry

import (
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

const testdata = "../../testdata/telemetry/"

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
		Interval: 1,
		UUID:     uuid,
	}

	service, err := NewService(cfg)
	require.NoError(t, err)
	service.Start()
	isRunning := service.IsRunning()
	assert.Equal(t, isRunning, true)

	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, count, 1)
	service.Stop()
	isRunning = service.IsRunning()
	assert.Equal(t, isRunning, false)
	assert.Equal(t, lastHeader, "pmm")

	// Test a service restart
	service.Start()
	isRunning = service.IsRunning()
	assert.Equal(t, isRunning, true)

	time.Sleep(2100 * time.Millisecond)
	assert.Equal(t, count, 3)
	service.Stop()
	isRunning = service.IsRunning()
	assert.Equal(t, isRunning, false)
}

func TestCollectData(t *testing.T) {
	config := &Config{}
	svc, err := NewService(config)
	require.NoError(t, err)

	m, err := svc.collectData()
	require.NoError(t, err)
	assert.NotEmpty(t, m)

	assert.Contains(t, m, "os_type")
	assert.Contains(t, m, "PMM")
}

func TestMakePayload(t *testing.T) {
	config := &Config{
		UUID: "ABCDEFG12345",
	}
	svc, err := NewService(config)
	require.NoError(t, err)

	m := map[string]interface{}{
		"os_type": "Kubuntu",
		"pmm":     "1.2.3",
	}

	b, err := svc.makePayload(m)
	require.NoError(t, err)
	want := "ABCDEFG12345;os_type;Kubuntu\nABCDEFG12345;pmm;1.2.3\n"
	assert.Equal(t, string(b), want)

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
}

// Reserved to test config save
func copyToTmp(filename string) (string, error) {
	tmpfile, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	content, err := ioutil.ReadAll(tmpfile)
	if err != nil {
		return "", nil
	}

	err = ioutil.WriteFile(tmpfile.Name(), content, os.ModePerm)
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
