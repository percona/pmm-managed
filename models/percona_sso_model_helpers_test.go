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

package models_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func setupDB(t *testing.T) (*reform.DB, func()) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	cleanup := func() {
		require.NoError(t, sqlDB.Close())
	}
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	return db, cleanup
}

func TestPerconaSSODetails(t *testing.T) {
	ctx := context.Background()

	t.Run("CorrectCredentials", func(t *testing.T) {
		clientID, clientSecret := os.Getenv("OAUTH_PMM_CLIENT_ID"), os.Getenv("OAUTH_PMM_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			t.Skip("Environment variables OAUTH_PMM_CLIENT_ID / OAUTH_PMM_CLIENT_SECRET are not defined, skipping test")
		}

		db, cleanup := setupDB(t)
		defer cleanup()

		expectedSSODetails := &models.PerconaSSODetails{
			IssuerURL:    "https://id-dev.percona.com/oauth2/aus15pi5rjdtfrcH51d7/v1/token?grant_type=client_credentials&scope=",
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scope:        "percona",
		}
		InsertSSODetails := &models.PerconaSSODetailsInsert{
			IssuerURL:    expectedSSODetails.IssuerURL,
			ClientID:     expectedSSODetails.ClientID,
			ClientSecret: expectedSSODetails.ClientSecret,
			Scope:        expectedSSODetails.Scope,
		}
		err := models.InsertPerconaSSODetails(ctx, db.Querier, InsertSSODetails)
		require.NoError(t, err)
		ssoDetails, err := models.GetPerconaSSODetails(ctx, db.Querier)
		require.NoError(t, err)
		assert.NotNil(t, ssoDetails)
		assert.Equal(t, expectedSSODetails.ClientID, ssoDetails.ClientID)
		assert.Equal(t, expectedSSODetails.ClientSecret, ssoDetails.ClientSecret)
		assert.Equal(t, expectedSSODetails.IssuerURL, ssoDetails.IssuerURL)
		assert.Equal(t, expectedSSODetails.Scope, ssoDetails.Scope)
		err = models.DeletePerconaSSODetails(db.Querier)
		require.NoError(t, err)
		ssoDetails, err = models.GetPerconaSSODetails(ctx, db.Querier)
		assert.Error(t, err)
		assert.Nil(t, ssoDetails)
		// See https://github.com/percona/pmm-managed/pull/852#discussion_r738178192
		ssoDetails, err = models.GetPerconaSSODetails(ctx, db.Querier)
		assert.Error(t, err)
		assert.Nil(t, ssoDetails)
	})

	t.Run("WrongCredentials", func(t *testing.T) {
		db, cleanup := setupDB(t)
		defer cleanup()

		InsertSSODetails := &models.PerconaSSODetailsInsert{
			IssuerURL:    "https://id-dev.percona.com",
			ClientID:     "wrongClientID",
			ClientSecret: "wrongClientSecret",
			Scope:        "percona",
		}
		err := models.InsertPerconaSSODetails(ctx, db.Querier, InsertSSODetails)
		require.Error(t, err)
	})
}
