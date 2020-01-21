package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func doRequest(t testing.TB, client *http.Client, req *http.Request) (*http.Response, []byte) {
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close() //nolint:errcheck

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, b
}

func TestAuth(t *testing.T) {
	t.Run("AuthErrors", func(t *testing.T) {
		for user, code := range map[*url.Userinfo]int{
			nil:                              401,
			url.UserPassword("bad", "wrong"): 401,
		} {
			user := user
			code := code
			t.Run(fmt.Sprintf("%s/%d", user, code), func(t *testing.T) {
				t.Parallel()

				// copy BaseURL and replace auth
				baseURL, err := url.Parse(pmmapitests.BaseURL.String())
				require.NoError(t, err)
				baseURL.User = user

				uri := baseURL.ResolveReference(&url.URL{
					Path: "v1/version",
				})
				t.Logf("URI: %s", uri)
				resp, err := http.Get(uri.String())
				require.NoError(t, err)
				defer resp.Body.Close() //nolint:errcheck

				b, err := httputil.DumpResponse(resp, true)
				require.NoError(t, err)
				assert.Equal(t, code, resp.StatusCode, "response:\n%s", b)
				require.False(t, bytes.Contains(b, []byte(`<html>`)), "response:\n%s", b)
			})
		}
	})

	t.Run("NormalErrors", func(t *testing.T) {
		for grpcCode, httpCode := range map[codes.Code]int{
			codes.Unauthenticated:  401,
			codes.PermissionDenied: 403,
		} {
			grpcCode := grpcCode
			httpCode := httpCode
			t.Run(fmt.Sprintf("%s/%d", grpcCode, httpCode), func(t *testing.T) {
				t.Parallel()

				res, err := serverClient.Default.Server.Version(&server.VersionParams{
					Dummy:   pointer.ToString(fmt.Sprintf("grpccode-%d", grpcCode)),
					Context: pmmapitests.Context,
				})
				assert.Empty(t, res)
				pmmapitests.AssertAPIErrorf(t, err, httpCode, grpcCode, "gRPC code %d (%s)", grpcCode, grpcCode)
			})
		}
	})
}

func TestSetup(t *testing.T) {
	// make a BaseURL without authentication
	baseURL, err := url.Parse(pmmapitests.BaseURL.String())
	require.NoError(t, err)
	baseURL.User = nil

	// make client that does not follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	t.Run("WebPage", func(t *testing.T) {
		t.Parallel()

		uri := baseURL.ResolveReference(&url.URL{
			Path: "/setup",
		})
		t.Logf("URI: %s", uri)
		req, err := http.NewRequest("GET", uri.String(), nil)
		require.NoError(t, err)
		req.Header.Set("X-Test-Must-Setup", "1")

		resp, b := doRequest(t, client, req)
		assert.Equal(t, 200, resp.StatusCode, "response:\n%s", b)
		assert.True(t, strings.HasPrefix(string(b), `<!doctype html>`), string(b))
	})

	t.Run("Redirect", func(t *testing.T) {
		paths := map[string]int{
			"graph":       303,
			"graph/":      303,
			"prometheus":  303,
			"prometheus/": 303,
			"qan":         303,
			"qan/":        303,
			"swagger":     303,
			"swagger/":    303,

			"v1/readyz":           200,
			"v1/AWSInstanceCheck": 405, // only POST is expected
			"v1/version":          401, // Grafana authentication required
		}
		for path, code := range paths {
			path, code := path, code
			t.Run(fmt.Sprintf("%s=%d", path, code), func(t *testing.T) {
				t.Parallel()

				uri := baseURL.ResolveReference(&url.URL{
					Path: path,
				})
				t.Logf("URI: %s", uri)
				req, err := http.NewRequest("GET", uri.String(), nil)
				require.NoError(t, err)
				req.Header.Set("X-Test-Must-Setup", "1")

				resp, b := doRequest(t, client, req)
				assert.Equal(t, code, resp.StatusCode, "response:\n%s", b)
				if code == 303 {
					assert.Equal(t, "/setup", resp.Header.Get("Location"))
				}
			})
		}
	})

	t.Run("API", func(t *testing.T) {
		t.Parallel()

		uri := baseURL.ResolveReference(&url.URL{
			Path: "v1/AWSInstanceCheck",
		})
		t.Logf("URI: %s", uri)
		b, err := json.Marshal(server.AWSInstanceCheckBody{
			InstanceID: "123",
		})
		require.NoError(t, err)
		req, err := http.NewRequest("POST", uri.String(), bytes.NewReader(b))
		require.NoError(t, err)
		req.Header.Set("X-Test-Must-Setup", "1")

		resp, b := doRequest(t, client, req)
		assert.Equal(t, 200, resp.StatusCode, "response:\n%s", b)
		assert.Equal(t, `{}`, string(b), "response:\n%s", b)
	})
}

func TestSwagger(t *testing.T) {
	// https://jira.percona.com/browse/PMM-5137

	// make client that does not follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, path := range []string{
		"swagger",
		"swagger/",
		"swagger.json",
		"swagger/swagger.json",
	} {
		path := path

		t.Run(path, func(t *testing.T) {
			t.Run("NoAuth", func(t *testing.T) {
				t.Parallel()

				// make a BaseURL without authentication
				baseURL, err := url.Parse(pmmapitests.BaseURL.String())
				require.NoError(t, err)
				baseURL.User = nil

				uri := baseURL.ResolveReference(&url.URL{
					Path: path,
				})
				t.Logf("URI: %s", uri)
				req, err := http.NewRequest("GET", uri.String(), nil)
				require.NoError(t, err)

				resp, _ := doRequest(t, client, req)
				require.NoError(t, err)
				assert.Equal(t, 200, resp.StatusCode)
			})

			t.Run("Auth", func(t *testing.T) {
				t.Parallel()

				uri := pmmapitests.BaseURL.ResolveReference(&url.URL{
					Path: path,
				})
				t.Logf("URI: %s", uri)
				req, err := http.NewRequest("GET", uri.String(), nil)
				require.NoError(t, err)

				resp, _ := doRequest(t, client, req)
				require.NoError(t, err)
				assert.Equal(t, 200, resp.StatusCode)
			})
		})
	}
}
