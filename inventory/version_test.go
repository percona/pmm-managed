package inventory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Percona-Lab/pmm-api-tests"
)

func prepareUrl(path string) *url.URL {
	return pmmapitests.BaseURL.ResolveReference(&url.URL{
		Path: path,
	})
}

func TestVersion(t *testing.T) {
	type VersionResponse struct {
		Version string
	}

	paths := []string{
		"ping",
		"managed/v1/version",
		"v1/version",
	}
	for _, path := range paths {
		t.Run(fmt.Sprintf("Get %s", path), func(t *testing.T) {
			uri := prepareUrl(path)
			var versionResponse VersionResponse

			resp, err := http.Get(uri.String())
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Body)
			require.Equal(t, resp.StatusCode, 200)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			err = json.Unmarshal(body, &versionResponse)
			require.NoError(t, err)
			require.NotNil(t, versionResponse)
			require.NotNil(t, versionResponse.Version)
			require.Equal(t, "2.0.0-dev", versionResponse.Version)
		})
	}
}
