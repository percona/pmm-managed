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

package server

import (
	"testing"

	"github.com/percona/pmm/api/serverpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/testdb"
)

func TestAWSInstanceChecker(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures)
	defer func() {
		require.NoError(t, sqlDB.Close())
	}()
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))

	t.Run("Docker", func(t *testing.T) {
		telemetry := new(mockTelemetryService)
		telemetry.Test(t)
		telemetry.On("DistributionMethod").Return(serverpb.DistributionMethod_DOCKER)
		defer telemetry.AssertExpectations(t)

		checker := NewAWSInstanceChecker(db, telemetry)
		assert.False(t, checker.MustCheck())
		assert.Error(t, checker.check("foo"))
	})

	t.Run("AMI", func(t *testing.T) {
		telemetry := new(mockTelemetryService)
		telemetry.Test(t)
		telemetry.On("DistributionMethod").Return(serverpb.DistributionMethod_AMI)
		defer telemetry.AssertExpectations(t)

		checker := NewAWSInstanceChecker(db, telemetry)
		assert.True(t, checker.MustCheck())
		assert.Error(t, checker.check("foo"))
	})
}
