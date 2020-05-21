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
	"strings"
	"testing"

	api "github.com/percona-platform/saas/gen/check/retrieval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services"
	"github.com/percona/pmm-managed/utils/testdb"
)

const (
	devChecksHost      = "check-dev.percona.com:443"
	devChecksPublicKey = "RWTg+ZmCCjt7O8eWeAmTLAqW+1ozUbpRSKSwNTmO+exlS5KEIPYWuYdX"
)

func TestDownloadChecks(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		s := New(nil, nil, nil, "2.5.0")
		s.host = devChecksHost
		s.publicKeys = []string{devChecksPublicKey}

		assert.Empty(t, s.getChecks())
		ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
		defer cancel()

		checks, err := s.downloadChecks(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, checks)
	})
}

func TestVerifySignatures(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		s := New(nil, nil, nil, "2.5.0")
		s.host = devChecksHost

		validKey := "RWSdGihBPffV2c4IysqHAIxc5c5PLfmQStbRPkuLXDr3igJOqFWt7aml"
		invalidKey := "RWSdGihBPffV2c4IysqHAIxc5c5PLfmQStbRPkuLXDr3igJO+INVALID"

		s.publicKeys = []string{invalidKey, validKey}

		validSign := strings.TrimSpace(`
untrusted comment: signature from minisign secret key
RWSdGihBPffV2W/zvmIiTLh8UnocoF3OcwmczGdZ+zM13eRnm2Qq9YxfQ9cLzAp1dA5w7C5a3Cp5D7jlYiydu5hqZhJUxJt/ugg=
trusted comment: some comment
uEF33ScMPYpvHvBKv8+yBkJ9k4+DCfV4nDs6kKYwGhalvkkqwWkyfJffO+KW7a1m3y42WHpOnzBxLJeU/AuzDw==
`)

		invalidSign := strings.TrimSpace(`
untrusted comment: signature from minisign secret key
RWSdGihBPffV2W/zvmIiTLh8UnocoF3OcwmczGdZ+zM13eRnm2Qq9YxfQ9cLzAp1dA5w7C5a3Cp5D7jlYiydu5hqZhJ+INVALID=
trusted comment: some comment
uEF33ScMPYpvHvBKv8+yBkJ9k4+DCfV4nDs6kKYwGhalvkkqwWkyfJffO+KW7a1m3y42WHpOnzBxLJ+INVALID==
`)

		resp := api.GetAllChecksResponse{
			File:       "random data",
			Signatures: []string{invalidSign, validSign},
		}

		err := s.verifySignatures(&resp)
		assert.NoError(t, err)
	})

	t.Run("empty signatures", func(t *testing.T) {
		s := New(nil, nil, nil, "2.5.0")
		s.host = devChecksHost
		s.publicKeys = []string{"RWSdGihBPffV2c4IysqHAIxc5c5PLfmQStbRPkuLXDr3igJOqFWt7aml"}

		resp := api.GetAllChecksResponse{
			File:       "random data",
			Signatures: []string{},
		}

		err := s.verifySignatures(&resp)
		assert.EqualError(t, err, "zero signatures received")
	})
}

func TestStartChecks(t *testing.T) {
	t.Run("stt disabled", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

		defer func() {
			require.NoError(t, sqlDB.Close())
		}()

		s := New(nil, nil, db, "2.5.0")
		err := s.StartChecks(context.Background())
		assert.EqualError(t, err, services.ErrSTTDisabled.Error())
	})

	t.Run("stt enabled", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

		defer func() {
			require.NoError(t, sqlDB.Close())
		}()

		var ar mockAlertRegistry
		ar.On("RemovePrefix", mock.Anything, mock.Anything).Return()

		s := New(nil, &ar, db, "2.5.0")
		settings, err := models.GetSettings(db)
		require.NoError(t, err)

		settings.SaaS.STTEnabled = true
		err = models.SaveSettings(db, settings)
		require.NoError(t, err)

		err = s.StartChecks(context.Background())
		require.NoError(t, err)
	})
}
