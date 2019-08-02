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

			assert.True(t, strings.HasPrefix(res.Version, "2.0.0-"), "version = %q should have '2.0.0-' prefix", res.Version)
			assert.True(t, strings.HasSuffix(res.FullVersion, ".el7"), "version = %q should have '.el7' suffix", res.FullVersion)
			ts := time.Time(res.Timestamp)
			assert.Equal(t, time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC), ts, "timestamp should contain only date")
		})
	}
}
