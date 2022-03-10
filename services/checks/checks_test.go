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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/percona-platform/saas/pkg/check"
	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/percona/pmm/version"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
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
	devChecksHost      = "check-dev.percona.com"
	devChecksPublicKey = "RWTg+ZmCCjt7O8eWeAmTLAqW+1ozUbpRSKSwNTmO+exlS5KEIPYWuYdX"
	testChecksFile     = "../../testdata/checks/checks.yml"
	issuerURL          = "https://id-dev.percona.com/oauth2/aus15pi5rjdtfrcH51d7/v1"
)

func TestDownloadChecks(t *testing.T) {
	clientID, clientSecret := os.Getenv("OAUTH_PMM_CLIENT_ID"), os.Getenv("OAUTH_PMM_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("Environment variables OAUTH_PMM_CLIENT_ID / OAUTH_PMM_CLIENT_SECRET are not defined, skipping test")
	}

	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

	insertSSODetails := &models.PerconaSSODetailsInsert{
		IssuerURL:    issuerURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        "percona",
	}
	err := models.InsertPerconaSSODetails(db.Querier, insertSSODetails)
	require.NoError(t, err)

	s, err := New(nil, nil, db)
	require.NoError(t, err)
	s.host = devChecksHost
	s.publicKeys = []string{devChecksPublicKey}

	t.Run("normal", func(t *testing.T) {
		checks, err := s.GetChecks()
		require.NoError(t, err)
		assert.Empty(t, checks)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dChecks, err := s.downloadChecks(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, dChecks)

		checks, err = s.GetChecks()
		require.NoError(t, err)
		assert.NotEmpty(t, checks)
	})

	t.Run("disabled telemetry", func(t *testing.T) {
		_, err := models.UpdateSettings(db.Querier, &models.ChangeSettingsParams{
			DisableTelemetry: true,
		})
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dChecks, err := s.downloadChecks(ctx)
		require.NoError(t, err)
		assert.Empty(t, dChecks)

		checks, err := s.GetChecks()
		require.NoError(t, err)
		assert.Empty(t, checks)
	})
}

func TestLoadLocalChecks(t *testing.T) {
	s, err := New(nil, nil, nil)
	require.NoError(t, err)

	checks, err := s.loadLocalChecks(testChecksFile)
	require.NoError(t, err)
	require.Len(t, checks, 5)

	c1, c2, c3, c4, c5 := checks[0], checks[1], checks[2], checks[3], checks[4]

	assert.Equal(t, check.PostgreSQLSelect, c1.Type)
	assert.Equal(t, "good_check_pg", c1.Name)
	assert.Equal(t, uint32(1), c1.Version)
	assert.Equal(t, "rolpassword FROM pg_authid WHERE rolcanlogin", c1.Query)

	assert.Equal(t, check.MySQLShow, c2.Type)
	assert.Equal(t, "bad_check_mysql", c2.Name)
	assert.Equal(t, uint32(1), c2.Version)
	assert.Equal(t, "VARIABLES LIKE 'version%'", c2.Query)

	assert.Equal(t, check.MongoDBBuildInfo, c3.Type)
	assert.Equal(t, "good_check_mongo", c3.Name)
	assert.Equal(t, uint32(1), c3.Version)
	assert.Empty(t, c3.Query)

	assert.Equal(t, check.MongoDBReplSetGetStatus, c4.Type)
	assert.Equal(t, "check_mongo_replSetGetStatus", c4.Name)
	assert.Equal(t, uint32(1), c4.Version)
	assert.Empty(t, c4.Query)

	assert.Equal(t, check.MongoDBGetDiagnosticData, c5.Type)
	assert.Equal(t, "check_mongo_getDiagnosticData", c5.Name)
	assert.Equal(t, uint32(1), c5.Version)
	assert.Empty(t, c5.Query)
}

func TestCollectChecks(t *testing.T) {
	t.Run("collect local checks", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())

		checks, err := s.GetChecks()
		require.NoError(t, err)
		require.Len(t, checks, 5)

		checkNames := make([]string, 0, len(checks))
		for _, c := range checks {
			checkNames = append(checkNames, c.Name)
		}
		assert.ElementsMatch(t, []string{
			"bad_check_mysql",
			"good_check_pg",
			"good_check_mongo",
			"check_mongo_replSetGetStatus",
			"check_mongo_getDiagnosticData",
		}, checkNames)
	})

	t.Run("download checks", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())
		assert.NotEmpty(t, s.checks)
	})
}

func TestDisableChecks(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())

		checks, err := s.GetChecks()
		require.NoError(t, err)
		assert.Len(t, checks, 5)

		disChecks, err := s.GetDisabledChecks()
		require.NoError(t, err)
		assert.Empty(t, disChecks)

		err = s.DisableChecks([]string{checks["bad_check_mysql"].Name})
		require.NoError(t, err)

		disChecks, err = s.GetDisabledChecks()
		require.NoError(t, err)
		assert.Len(t, disChecks, 1)
	})

	t.Run("disable same check twice", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())

		checks, err := s.GetChecks()
		require.NoError(t, err)
		assert.Len(t, checks, 5)

		disChecks, err := s.GetDisabledChecks()
		require.NoError(t, err)
		assert.Empty(t, disChecks)

		err = s.DisableChecks([]string{checks["bad_check_mysql"].Name})
		require.NoError(t, err)

		err = s.DisableChecks([]string{checks["bad_check_mysql"].Name})
		require.NoError(t, err)

		disChecks, err = s.GetDisabledChecks()
		require.NoError(t, err)
		assert.Len(t, disChecks, 1)
	})

	t.Run("disable unknown check", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())

		err = s.DisableChecks([]string{"unknown_check"})
		require.Error(t, err)

		disChecks, err := s.GetDisabledChecks()
		require.NoError(t, err)
		assert.Empty(t, disChecks)
	})
}

func TestEnableChecks(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())

		checks, err := s.GetChecks()
		require.NoError(t, err)
		assert.Len(t, checks, 5)

		err = s.DisableChecks([]string{checks["bad_check_mysql"].Name, checks["good_check_pg"].Name, checks["good_check_mongo"].Name})
		require.NoError(t, err)

		err = s.EnableChecks([]string{checks["good_check_pg"].Name, checks["good_check_mongo"].Name})
		require.NoError(t, err)

		disChecks, err := s.GetDisabledChecks()
		require.NoError(t, err)
		assert.Equal(t, []string{checks["bad_check_mysql"].Name}, disChecks)

		enabledChecksCount := len(checks) - len(disChecks)
		assert.Equal(t, 4, enabledChecksCount)
	})
}

func TestChangeInterval(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		var ams mockAlertmanagerService
		ams.On("SendAlerts", mock.Anything, mock.Anything).Return()
		sqlDB := testdb.Open(t, models.SkipFixtures, nil)
		db := reform.NewDB(sqlDB, postgresql.Dialect, nil)
		s, err := New(nil, &ams, db)
		require.NoError(t, err)
		s.localChecksFile = testChecksFile

		s.CollectChecks(context.Background())

		checks, err := s.GetChecks()
		require.NoError(t, err)
		assert.Len(t, checks, 5)

		// change all check intervals from standard to rare
		params := make(map[string]check.Interval)
		for _, c := range checks {
			params[c.Name] = check.Rare
		}
		err = s.ChangeInterval(params)
		require.NoError(t, err)

		updatedChecks, err := s.GetChecks()
		require.NoError(t, err)
		for _, c := range updatedChecks {
			assert.Equal(t, check.Rare, c.Interval)
		}

		t.Run("preserve intervals on restarts", func(t *testing.T) {
			settings, err := models.GetSettings(db)
			require.NoError(t, err)

			settings.SaaS.STTEnabled = true
			err = models.SaveSettings(db, settings)
			require.NoError(t, err)

			err = s.runChecksGroup(context.Background(), "")
			require.NoError(t, err)

			checks, err := s.GetChecks()
			require.NoError(t, err)
			for _, c := range checks {
				assert.Equal(t, check.Rare, c.Interval)
			}
		})
	})
}

// A proper unit test could not be written due
// to problems with the code responsible for locating agents
// Once it is fixed rewrite this test to actually run `executeChecks`
// method and test for recorded metrics.
func TestSTTMetrics(t *testing.T) {
	t.Run("check for recorded metrics", func(t *testing.T) {
		s, err := New(nil, nil, nil)
		require.NoError(t, err)
		expected := strings.NewReader(`
		    # HELP pmm_managed_checks_alerts_generated_total Counter of alerts generated per service type per check type
		    # TYPE pmm_managed_checks_alerts_generated_total counter
		    pmm_managed_checks_alerts_generated_total{check_type="MONGODB_BUILDINFO",service_type="mongodb"} 0
		    pmm_managed_checks_alerts_generated_total{check_type="MONGODB_GETCMDLINEOPTS",service_type="mongodb"} 0
			pmm_managed_checks_alerts_generated_total{check_type="MONGODB_GETDIAGNOSTICDATA",service_type="mongodb"} 0
		    pmm_managed_checks_alerts_generated_total{check_type="MONGODB_GETPARAMETER",service_type="mongodb"} 0
			pmm_managed_checks_alerts_generated_total{check_type="MONGODB_REPLSETGETSTATUS",service_type="mongodb"} 0
		    pmm_managed_checks_alerts_generated_total{check_type="MYSQL_SELECT",service_type="mysql"} 0
		    pmm_managed_checks_alerts_generated_total{check_type="MYSQL_SHOW",service_type="mysql"} 0
		    pmm_managed_checks_alerts_generated_total{check_type="POSTGRESQL_SELECT",service_type="postgresql"} 0
		    pmm_managed_checks_alerts_generated_total{check_type="POSTGRESQL_SHOW",service_type="postgresql"} 0
		    # HELP pmm_managed_checks_scripts_executed_total Counter of check scripts executed per service type
		    # TYPE pmm_managed_checks_scripts_executed_total counter
		    pmm_managed_checks_scripts_executed_total{service_type="mongodb"} 0
		    pmm_managed_checks_scripts_executed_total{service_type="mysql"} 0
		    pmm_managed_checks_scripts_executed_total{service_type="postgresql"} 0
		`)
		assert.NoError(t, promtest.CollectAndCompare(s, expected))
	})
}

func TestGetSecurityCheckResults(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

	t.Run("STT disabled", func(t *testing.T) {
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		results, err := s.GetSecurityCheckResults()
		assert.Nil(t, results)
		assert.EqualError(t, err, services.ErrSTTDisabled.Error())
	})

	t.Run("STT enabled", func(t *testing.T) {
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		settings, err := models.GetSettings(db)
		require.NoError(t, err)

		settings.SaaS.STTEnabled = true
		err = models.SaveSettings(db, settings)
		require.NoError(t, err)

		results, err := s.GetSecurityCheckResults()
		assert.Empty(t, results)
		require.NoError(t, err)
	})
}

func TestStartChecks(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

	t.Run("stt disabled", func(t *testing.T) {
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		err = s.runChecksGroup(context.Background(), "")
		assert.EqualError(t, err, services.ErrSTTDisabled.Error())
	})

	t.Run("unknown interval", func(t *testing.T) {
		s, err := New(nil, nil, db)
		require.NoError(t, err)
		settings, err := models.GetSettings(db)
		require.NoError(t, err)

		settings.SaaS.STTEnabled = true
		err = models.SaveSettings(db, settings)
		require.NoError(t, err)

		err = s.runChecksGroup(context.Background(), check.Interval("unknown"))
		assert.EqualError(t, err, "unknown check interval: unknown")
	})

	t.Run("stt enabled", func(t *testing.T) {
		var ams mockAlertmanagerService
		ams.On("SendAlerts", mock.Anything, mock.Anything).Return()

		s, err := New(nil, &ams, db)
		require.NoError(t, err)
		settings, err := models.GetSettings(db)
		require.NoError(t, err)

		settings.SaaS.STTEnabled = true
		err = models.SaveSettings(db, settings)
		require.NoError(t, err)

		err = s.runChecksGroup(context.Background(), "")
		require.NoError(t, err)
	})
}

func TestFilterChecks(t *testing.T) {
	t.Parallel()

	valid := []check.Check{
		{Name: "MySQLShow", Version: 1, Type: check.MySQLShow},
		{Name: "MySQLSelect", Version: 1, Type: check.MySQLSelect},
		{Name: "PostgreSQLShow", Version: 1, Type: check.PostgreSQLShow},
		{Name: "PostgreSQLSelect", Version: 1, Type: check.PostgreSQLSelect},
		{Name: "MongoDBGetParameter", Version: 1, Type: check.MongoDBGetParameter},
		{Name: "MongoDBBuildInfo", Version: 1, Type: check.MongoDBBuildInfo},
		{Name: "MongoDBGetCmdLineOpts", Version: 1, Type: check.MongoDBGetCmdLineOpts},
		{Name: "MongoDBReplSetGetStatus", Version: 1, Type: check.MongoDBReplSetGetStatus},
		{Name: "MongoDBGetDiagnosticData", Version: 1, Type: check.MongoDBGetDiagnosticData},
	}

	invalid := []check.Check{
		{Name: "unsupported version", Version: maxSupportedVersion + 1, Type: check.MySQLShow},
		{Name: "unsupported type", Version: 1, Type: check.Type("RedisInfo")},
		{Name: "missing type", Version: 1},
	}

	checks := append(valid, invalid...)

	s, err := New(nil, nil, nil)
	require.NoError(t, err)
	actual := s.filterSupportedChecks(checks)
	assert.ElementsMatch(t, valid, actual)
}

func TestGroupChecksByDB(t *testing.T) {
	t.Parallel()

	checks := map[string]check.Check{
		"MySQLShow":                {Name: "MySQLShow", Version: 1, Type: check.MySQLShow},
		"MySQLSelect":              {Name: "MySQLSelect", Version: 1, Type: check.MySQLSelect},
		"PostgreSQLShow":           {Name: "PostgreSQLShow", Version: 1, Type: check.PostgreSQLShow},
		"PostgreSQLSelect":         {Name: "PostgreSQLSelect", Version: 1, Type: check.PostgreSQLSelect},
		"MongoDBGetParameter":      {Name: "MongoDBGetParameter", Version: 1, Type: check.MongoDBGetParameter},
		"MongoDBBuildInfo":         {Name: "MongoDBBuildInfo", Version: 1, Type: check.MongoDBBuildInfo},
		"MongoDBGetCmdLineOpts":    {Name: "MongoDBGetCmdLineOpts", Version: 1, Type: check.MongoDBGetCmdLineOpts},
		"MongoDBReplSetGetStatus":  {Name: "MongoDBReplSetGetStatus", Version: 1, Type: check.MongoDBReplSetGetStatus},
		"MongoDBGetDiagnosticData": {Name: "MongoDBGetDiagnosticData", Version: 1, Type: check.MongoDBGetDiagnosticData},
		"unsupported type":         {Name: "unsupported type", Version: 1, Type: check.Type("RedisInfo")},
		"missing type":             {Name: "missing type", Version: 1},
	}

	s, err := New(nil, nil, nil)
	require.NoError(t, err)
	mySQLChecks, postgreSQLChecks, mongoDBChecks := s.groupChecksByDB(checks)

	require.Len(t, mySQLChecks, 2)
	require.Len(t, postgreSQLChecks, 2)
	require.Len(t, mongoDBChecks, 5)

	assert.Equal(t, check.MySQLShow, mySQLChecks["MySQLShow"].Type)
	assert.Equal(t, check.MySQLSelect, mySQLChecks["MySQLSelect"].Type)

	assert.Equal(t, check.PostgreSQLShow, postgreSQLChecks["PostgreSQLShow"].Type)
	assert.Equal(t, check.PostgreSQLSelect, postgreSQLChecks["PostgreSQLSelect"].Type)

	assert.Equal(t, check.MongoDBGetParameter, mongoDBChecks["MongoDBGetParameter"].Type)
	assert.Equal(t, check.MongoDBBuildInfo, mongoDBChecks["MongoDBBuildInfo"].Type)
	assert.Equal(t, check.MongoDBGetCmdLineOpts, mongoDBChecks["MongoDBGetCmdLineOpts"].Type)
	assert.Equal(t, check.MongoDBReplSetGetStatus, mongoDBChecks["MongoDBReplSetGetStatus"].Type)
	assert.Equal(t, check.MongoDBGetDiagnosticData, mongoDBChecks["MongoDBGetDiagnosticData"].Type)
}

func setup(t *testing.T, db *reform.DB, serviceName, nodeID, pmmAgentVersion string) {
	pmmAgent, err := models.CreatePMMAgent(db.Querier, nodeID, nil)
	require.NoError(t, err)

	pmmAgent.Version = pointer.ToStringOrNil(pmmAgentVersion)
	err = db.Update(pmmAgent)
	require.NoError(t, err)

	mysql, err := models.AddNewService(db.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
		ServiceName: serviceName,
		NodeID:      nodeID,
		Address:     pointer.ToString("127.0.0.1"),
		Port:        pointer.ToUint16(3306),
	})
	require.NoError(t, err)

	_, err = models.CreateAgent(db.Querier, models.MySQLdExporterType, &models.CreateAgentParams{
		PMMAgentID: pmmAgent.AgentID,
		ServiceID:  mysql.ServiceID,
	})
	require.NoError(t, err)
}

func TestFindTargets(t *testing.T) {
	sqlDB := testdb.Open(t, models.SetupFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
	s, err := New(nil, nil, db)
	require.NoError(t, err)

	t.Run("unknown service", func(t *testing.T) {
		t.Parallel()

		targets, err := s.findTargets(models.PostgreSQLServiceType, nil)
		require.NoError(t, err)
		assert.Len(t, targets, 0)
	})

	t.Run("different pmm agent versions", func(t *testing.T) {
		t.Parallel()

		node, err := models.CreateNode(db.Querier, models.GenericNodeType, &models.CreateNodeParams{
			NodeName: "test-node",
		})
		require.NoError(t, err)

		setup(t, db, "mysql1", node.NodeID, "")
		setup(t, db, "mysql2", node.NodeID, "2.5.0")
		setup(t, db, "mysql3", node.NodeID, "2.6.0")
		setup(t, db, "mysql4", node.NodeID, "2.6.1")
		setup(t, db, "mysql5", node.NodeID, "2.7.0")

		tests := []struct {
			name               string
			minRequiredVersion *version.Parsed
			count              int
		}{
			{"without version", nil, 5},
			{"version 2.5.0", version.MustParse("2.5.0"), 4},
			{"version 2.6.0", version.MustParse("2.6.0"), 3},
			{"version 2.6.1", version.MustParse("2.6.1"), 2},
			{"version 2.7.0", version.MustParse("2.7.0"), 1},
			{"version 2.9.0", version.MustParse("2.9.0"), 0},
		}

		for _, test := range tests {
			test := test

			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				targets, err := s.findTargets(models.MySQLServiceType, test.minRequiredVersion)
				require.NoError(t, err)
				assert.Len(t, targets, test.count)
			})
		}
	})
}

func TestFilterChecksByInterval(t *testing.T) {
	t.Parallel()
	s, err := New(nil, nil, nil)
	require.NoError(t, err)

	rareCheck := check.Check{Name: "rareCheck", Interval: check.Rare}
	standardCheck := check.Check{Name: "standardCheck", Interval: check.Standard}
	frequentCheck := check.Check{Name: "frequentCheck", Interval: check.Frequent}
	emptyCheck := check.Check{Name: "emptyCheck"}

	checks := map[string]check.Check{
		rareCheck.Name:     rareCheck,
		standardCheck.Name: standardCheck,
		frequentCheck.Name: frequentCheck,
		emptyCheck.Name:    emptyCheck,
	}

	rareChecks := s.filterChecks(checks, check.Rare, nil, nil)
	assert.Equal(t, map[string]check.Check{"rareCheck": rareCheck}, rareChecks)

	standardChecks := s.filterChecks(checks, check.Standard, nil, nil)
	assert.Equal(t, map[string]check.Check{"standardCheck": standardCheck, "emptyCheck": emptyCheck}, standardChecks)

	frequentChecks := s.filterChecks(checks, check.Frequent, nil, nil)
	assert.Equal(t, map[string]check.Check{"frequentCheck": frequentCheck}, frequentChecks)
}

func TestGetFailedChecks(t *testing.T) {
	sqlDB := testdb.Open(t, models.SkipFixtures, nil)
	db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

	t.Run("STT disabled", func(t *testing.T) {
		var ams mockAlertmanagerService
		ams.On("GetFilteredAlerts", mock.Anything, mock.Anything).Return(nil, services.ErrSTTDisabled)

		s, err := New(nil, &ams, db)
		require.NoError(t, err)
		results, err := s.GetFailedChecks(context.Background(), "test_svc")
		assert.Nil(t, results)
		assert.EqualError(t, err, services.ErrSTTDisabled.Error())
	})

	t.Run("STT enabled", func(t *testing.T) {
		var ams mockAlertmanagerService
		ams.On("GetFilteredAlerts", mock.Anything, mock.Anything).Return([]*ammodels.GettableAlert{}, nil)

		s, err := New(nil, &ams, db)
		require.NoError(t, err)
		settings, err := models.GetSettings(db)
		require.NoError(t, err)

		settings.SaaS.STTEnabled = true
		err = models.SaveSettings(db, settings)
		require.NoError(t, err)

		results, err := s.GetFailedChecks(context.Background(), "test_svc")
		assert.Empty(t, results)
		require.NoError(t, err)
	})
}
