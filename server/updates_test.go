package server

import (
	"io"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestCheckUpdates(t *testing.T) {
	// do not run this test in parallel with other tests as it also tests timings

	const fast, slow = 3 * time.Second, 60 * time.Second

	// that call should always be fast
	version, err := serverClient.Default.Server.Version(server.NewVersionParamsWithTimeout(fast))
	require.NoError(t, err)
	if version.Payload.Server == nil || version.Payload.Server.Version == "" {
		t.Skip("skipping test in developer's environment")
	}

	params := &server.CheckUpdatesParams{
		Context: pmmapitests.Context,
	}
	params.SetTimeout(slow) // that call can be slow with a cold cache
	res, err := serverClient.Default.Server.CheckUpdates(params)
	require.NoError(t, err)

	require.NotEmpty(t, res.Payload.Installed)
	assert.True(t, strings.HasPrefix(res.Payload.Installed.Version, "2.0.0-"),
		"installed.version = %q should have '2.0.0-' prefix", res.Payload.Installed.Version)
	assert.NotEmpty(t, res.Payload.Installed.FullVersion)
	require.NotEmpty(t, res.Payload.Installed.Timestamp)
	ts := time.Time(res.Payload.Installed.Timestamp)
	hour, min, _ := ts.Clock()
	assert.Zero(t, hour, "installed.timestamp should contain only date")
	assert.Zero(t, min, "installed.timestamp should contain only date")

	require.NotEmpty(t, res.Payload.Latest)
	assert.True(t, strings.HasPrefix(res.Payload.Latest.Version, "2.0.0-"),
		"latest.version = %q should have '2.0.0-' prefix", res.Payload.Latest.Version)
	assert.NotEmpty(t, res.Payload.Latest.FullVersion)
	require.NotEmpty(t, res.Payload.Latest.Timestamp)
	ts = time.Time(res.Payload.Latest.Timestamp)
	hour, min, _ = ts.Clock()
	assert.Zero(t, hour, "latest.timestamp should contain only date")
	assert.Zero(t, min, "latest.timestamp should contain only date")

	assert.Equal(t, res.Payload.Installed.FullVersion != res.Payload.Latest.FullVersion, res.Payload.UpdateAvailable)
	assert.NotEmpty(t, res.Payload.LastCheck)

	t.Run("HotCache", func(t *testing.T) {
		params = &server.CheckUpdatesParams{
			Context: pmmapitests.Context,
		}
		params.SetTimeout(fast) // that call should be fast with hot cache
		resAgain, err := serverClient.Default.Server.CheckUpdates(params)
		require.NoError(t, err)

		assert.Equal(t, res.Payload, resAgain.Payload)
	})

	t.Run("Force", func(t *testing.T) {
		params = &server.CheckUpdatesParams{
			Body: server.CheckUpdatesBody{
				Force: true,
			},
			Context: pmmapitests.Context,
		}
		params.SetTimeout(slow) // that call with force can be slow
		resForce, err := serverClient.Default.Server.CheckUpdates(params)
		require.NoError(t, err)

		assert.Equal(t, res.Payload.Installed, resForce.Payload.Installed)
		assert.Equal(t, resForce.Payload.Installed.FullVersion != resForce.Payload.Latest.FullVersion, resForce.Payload.UpdateAvailable)
		assert.NotEqual(t, res.Payload.LastCheck, resForce.Payload.LastCheck)
	})
}

// sync test name with Makefile
func TestUpdate(t *testing.T) {
	// do not run this test in parallel with other tests

	if !pmmapitests.RunUpdateTest {
		t.Skip("skipping PMM Server update test")
	}

	// make a new client without authentication
	baseURL, err := url.Parse(pmmapitests.BaseURL.String())
	require.NoError(t, err)
	baseURL.User = nil
	noAuthClient := serverClient.New(pmmapitests.Transport(baseURL, true), nil)

	// without authentication
	_, err = noAuthClient.Server.StartUpdate(nil)
	pmmapitests.AssertAPIErrorf(t, err, 401, codes.Unauthenticated, "Unauthorized")

	// with authentication
	startRes, err := serverClient.Default.Server.StartUpdate(nil)
	require.NoError(t, err)
	authToken := startRes.Payload.AuthToken
	logOffset := startRes.Payload.LogOffset
	require.NotEmpty(t, authToken)
	assert.Zero(t, logOffset)

	_, err = serverClient.Default.Server.StartUpdate(nil)
	pmmapitests.AssertAPIErrorf(t, err, 400, codes.FailedPrecondition, "Update is already running.")

	// without token
	_, err = noAuthClient.Server.UpdateStatus(&server.UpdateStatusParams{
		Body: server.UpdateStatusBody{
			LogOffset: logOffset,
		},
		Context: pmmapitests.Context,
	})
	pmmapitests.AssertAPIErrorf(t, err, 403, codes.PermissionDenied, "Invalid authentication token.")

	// read log lines like UI would do, but without delays to increase a chance for race detector to spot something
	var retries int
	for {
		start := time.Now()
		statusRes, err := noAuthClient.Server.UpdateStatus(&server.UpdateStatusParams{
			Body: server.UpdateStatusBody{
				AuthToken: authToken,
				LogOffset: logOffset,
			},
			Context: pmmapitests.Context,
		})
		if err != nil {
			// check that we know and understand all possible errors
			switch err := err.(type) {
			case *pmmapitests.ErrFromNginx:
				// nothing
			case *url.Error:
				switch err2 := err.Err.(type) {
				case *net.OpError:
					t.Logf("%[1]T %[1]s", err2)
				default:
					assert.Equal(t, io.EOF, err.Err, "%[1]T %[1]s", err)
				}
			case *server.UpdateStatusDefault:
				assert.Equal(t, 503, err.Code(), "%[1]T %[1]s", err)
			default:
				t.Fatalf("%#v", err)
			}
			continue
		}
		dur := time.Since(start)
		t.Logf("%s, offset = %d->%d, done = %t:\n%s", dur, logOffset, statusRes.Payload.LogOffset,
			statusRes.Payload.Done, strings.Join(statusRes.Payload.LogLines, "\n"))

		if statusRes.Payload.LogOffset == logOffset {
			// pmm-managed waits up to 30 seconds for new log lines. Usually, that's more than enough for
			// Ansible playbook to produce a new output, and that test checks that. However, our Jenkins node
			// is very slow, so we try several times.
			// That code should be removed once Jenkins performance is fixed.
			t.Logf("retries = %d", retries)
			if !statusRes.Payload.Done {
				retries++
				if retries < 5 {
					assert.InDelta(t, (30 * time.Second).Seconds(), dur.Seconds(), (3 * time.Second).Seconds())
					continue
				}
			}

			assert.Empty(t, statusRes.Payload.LogLines, "lines should be empty for the same offset")
			require.True(t, statusRes.Payload.Done, "lines should be empty only when done")
			break
		}

		retries = 0
		assert.NotEmpty(t, statusRes.Payload.LogLines, "pmm-managed should delay response until some lines are available")
		assert.True(t, statusRes.Payload.LogOffset > logOffset,
			"expected statusRes.Payload.LogOffset (%d) > logOffset (%d)",
			statusRes.Payload.LogOffset, logOffset,
		)
		logOffset = statusRes.Payload.LogOffset
	}

	// extra check for done
	statusRes, err := noAuthClient.Server.UpdateStatus(&server.UpdateStatusParams{
		Body: server.UpdateStatusBody{
			AuthToken: authToken,
			LogOffset: logOffset,
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	assert.True(t, statusRes.Payload.Done, "should be done")
	assert.Empty(t, statusRes.Payload.LogLines, "lines should be empty when done")
	assert.Equal(t, logOffset, statusRes.Payload.LogOffset)

	// whole log
	statusRes, err = noAuthClient.Server.UpdateStatus(&server.UpdateStatusParams{
		Body: server.UpdateStatusBody{
			AuthToken: authToken,
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	assert.True(t, statusRes.Payload.Done, "should be done")
	assert.Equal(t, int(logOffset), len(strings.Join(statusRes.Payload.LogLines, "\n")+"\n"))
	assert.Equal(t, logOffset, statusRes.Payload.LogOffset)
}
