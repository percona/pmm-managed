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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestOktaSSODetails(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	expectedSSODetails := &models.OktaSSODetails{}
	err := models.InsertOktaSSODetails(db.Querier, expectedSSODetails)
	require.NoError(t, err)
	ssoDetails, err := models.GetOktaSSODetails(db.Querier)
	require.NoError(t, err)
	assert.NotNil(t, ssoDetails)
	assert.Equal(t, expectedSSODetails.ClientID, ssoDetails.ClientID)
	assert.Equal(t, expectedSSODetails.ClientSecret, ssoDetails.ClientSecret)
	assert.Equal(t, expectedSSODetails.IssuerURL, ssoDetails.IssuerURL)
	assert.Equal(t, expectedSSODetails.Scope, ssoDetails.Scope)
	err = models.DeleteOktaSSODetails(db.Querier)
	require.NoError(t, err)
	ssoDetails, err = models.GetOktaSSODetails(db.Querier)
	assert.Error(t, err)
	assert.Nil(t, ssoDetails)
}
