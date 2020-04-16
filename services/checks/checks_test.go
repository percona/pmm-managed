package checks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	devChecksHost      = "check-dev.percona.com:443"
	devChecksPublicKey = "RWS69zYk2LOS7gWnSQNgnPRbBEwaoG3N/ITwDqfowUItfHvrpfQ++D0g"
)

func TestDownloadChecks(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		s := New("2.5.0")
		s.host = devChecksHost
		s.publicKey = devChecksPublicKey

		assert.Empty(t, s.Checks())
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := s.downloadChecks(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, s.Checks())
	})
}
