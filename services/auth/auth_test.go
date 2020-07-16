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

package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/brianvoe/gofakeit"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

const devAuthHost = "check-dev.percona.com:443"

func TestAuth(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

	defer func() {
		require.NoError(t, sqlDB.Close())
	}()

	s := New(db)
	s.host = devAuthHost

	login := gofakeit.Email()
	password := "Password12345"

	// SignUp test
	err := s.SignUp(context.Background(), login, password)
	require.NoError(t, err)

	// SignIn test
	settings, err := models.GetSettings(s.db)
	require.NoError(t, err)
	require.Empty(t, settings.SessionID)

	err = s.SignIn(context.Background(), login, password)
	require.NoError(t, err)

	settings, err = models.GetSettings(s.db)
	require.NoError(t, err)
	require.NotEmpty(t, settings.SessionID)

	// RefreshSession test
	err = s.RefreshSession(context.Background())
	require.NoError(t, err)
}

func init() { //nolint:gochecknoinits
	gofakeit.Seed(time.Now().UnixNano())
}
