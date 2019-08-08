package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestVersion(t *testing.T) {
	paths := []string{
		"managed/v1/version",
		"v1/version",
	}
	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			uri := pmmapitests.BaseURL.ResolveReference(&url.URL{
				Path: path,
			})

			t.Logf("URI: %s", uri)
			resp, err := http.Get(uri.String())
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck
			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode, "response:\n%s", b)

			var res server.VersionOKBody
			err = json.Unmarshal(b, &res)
			require.NoError(t, err, "response:\n%s", b)

			require.True(t, strings.HasPrefix(res.Version, "2.0.0-"),
				"version = %q must have '2.0.0-' prefix for PMM 1.x's pmm-client compatibility checking", res.Version)

			require.NotEmpty(t, res.Managed)
			assert.True(t, strings.HasPrefix(res.Managed.Version, "2.0.0-"),
				"managed.version = %q should have '2.0.0-' prefix", res.Managed.Version)
			assert.NotEmpty(t, res.Managed.FullVersion)
			assert.NotEmpty(t, res.Managed.Timestamp)
			ts := time.Time(res.Managed.Timestamp)
			hour, min, _ := ts.Clock()
			assert.NotZero(t, hour, "managed timestamp should not contain only date")
			assert.NotZero(t, min, "managed timestamp should not contain only date")

			if res.Server == nil || res.Server.Version == "" {
				t.Skip("skipping the rest of the test in developer's environment")
			}

			require.NotEmpty(t, res.Server)
			assert.True(t, strings.HasPrefix(res.Server.Version, "2.0.0-"),
				"server.version = %q should have '2.0.0-' prefix", res.Server.Version)
			assert.NotEmpty(t, res.Server.FullVersion)
			require.NotEmpty(t, res.Server.Timestamp)
			ts = time.Time(res.Server.Timestamp)
			hour, min, _ = ts.Clock()
			assert.NotZero(t, hour, "server timestamp should not contain only date")
			assert.NotZero(t, min, "server timestamp should not contain only date")
		})
	}
}
