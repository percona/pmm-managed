package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestVersion(t *testing.T) {
	paths := []string{
		"ping",
		"managed/v1/version",
		"v1/version",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			uri := pmmapitests.BaseURL.ResolveReference(&url.URL{
				Path: path,
			})

			resp, err := http.Get(uri.String())
			require.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck
			assert.Equal(t, resp.StatusCode, 200)

			var res struct {
				Version string
			}
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)

			assert.True(t, strings.HasPrefix(res.Version, "2.0.0-"), "version = %q", res.Version)
		})
	}
}
