package server

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestReadyz(t *testing.T) {
	paths := []string{
		"ping",
		"v1/readyz",
	}
	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			// make a BaseURL without authentication
			baseURL, err := url.Parse(pmmapitests.BaseURL.String())
			require.NoError(t, err)
			baseURL.User = nil

			uri := baseURL.ResolveReference(&url.URL{
				Path: path,
			})

			t.Logf("URI: %s", uri)
			resp, err := http.Get(uri.String())
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck
			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode, "response:\n%s", b)
			assert.Equal(t, "{\n\n}", string(b))
		})
	}
}
