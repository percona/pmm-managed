// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

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
