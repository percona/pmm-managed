package server

import (
	"strings"
	"testing"
	"time"

	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdates(t *testing.T) {
	t.Run("CheckUpdates", func(t *testing.T) {
		t.Parallel()

		version, err := serverClient.Default.Server.Version(nil)
		require.NoError(t, err)
		if !strings.HasSuffix(version.Payload.FullVersion, ".el7") {
			t.Skip("skipping test in developer's environment")
		}

		res, err := serverClient.Default.Server.CheckUpdates(nil)
		require.NoError(t, err)

		assert.True(t, strings.HasPrefix(res.Payload.Version, "2.0.0-"), "version = %q should has '2.0.0-' prefix", res.Payload.Version)
		assert.True(t, strings.HasSuffix(res.Payload.FullVersion, ".el7"), "version = %q should has '.el7' suffix", res.Payload.FullVersion)
		ts := time.Time(res.Payload.Timestamp)
		assert.Equal(t, time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC), ts, "Timestamp should contain only date")

		assert.True(t, strings.HasPrefix(res.Payload.LatestVersion, "2.0.0-"), "version = %q should has '2.0.0-' prefix", res.Payload.LatestVersion)
		assert.True(t, strings.HasSuffix(res.Payload.LatestFullVersion, ".el7"), "version = %q should has '.el7' suffix", res.Payload.LatestFullVersion)
		ts = time.Time(res.Payload.LatestTimestamp)
		assert.Equal(t, time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.UTC), ts, "LatestTimestamp should contain only date")
		assert.Empty(t, res.Payload.LatestNewsURL)
	})
}
