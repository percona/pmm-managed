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

package supervisord

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseEvent(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		t.Parallel()

		log := strings.Split(`
			2019-08-08 17:09:41,806 INFO spawned: 'pmm-update-perform' with pid 12983
			2019-08-08 17:09:43,509 INFO success: pmm-update-perform entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:09:48,494 INFO exited: pmm-update-perform (exit status 1; not expected)
			2019-08-08 17:09:48,506 INFO spawned: 'pmm-update-perform' with pid 13000
			2019-08-08 17:09:49,506 INFO success: pmm-update-perform entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:09:57,284 INFO received SIGUSR2 indicating log reopen request
			2019-08-08 17:09:57,284 INFO supervisord logreopen
			2019-08-08 17:09:57,854 INFO waiting for pmm-managed to stop
			2019-08-08 17:09:59,854 INFO waiting for pmm-managed to stop
			2019-08-08 17:10:00,863 INFO stopped: pmm-managed (exit status 0)
			2019-08-08 17:10:01,932 INFO spawned: 'pmm-managed' with pid 13191
			2019-08-08 17:10:03,006 INFO success: pmm-managed entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:08,138 INFO waiting for nginx to stop
			2019-08-08 17:10:08,258 INFO stopped: nginx (exit status 0)
			2019-08-08 17:10:09,321 INFO spawned: 'nginx' with pid 13335
			2019-08-08 17:10:09,848 INFO waiting for postgresql to stop
			2019-08-08 17:10:09,877 INFO stopped: postgresql (exit status 0)
			2019-08-08 17:10:09,878 INFO reaped unknown pid 12411
			2019-08-08 17:10:10,435 INFO success: nginx entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:10,980 INFO spawned: 'postgresql' with pid 13397
			2019-08-08 17:10:11,525 INFO waiting for prometheus to stop
			2019-08-08 17:10:11,535 INFO stopped: prometheus (exit status 0)
			2019-08-08 17:10:12,109 INFO success: postgresql entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:12,640 INFO spawned: 'prometheus' with pid 13481
			2019-08-08 17:10:13,234 INFO waiting for clickhouse to stop
			2019-08-08 17:10:13,639 INFO success: prometheus entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:15,234 INFO waiting for clickhouse to stop
			2019-08-08 17:10:17,234 INFO waiting for clickhouse to stop
			2019-08-08 17:10:19,234 INFO waiting for clickhouse to stop
			2019-08-08 17:10:21,028 INFO stopped: clickhouse (exit status 0)
			2019-08-08 17:10:22,108 INFO spawned: 'clickhouse' with pid 13567
			2019-08-08 17:10:22,624 INFO waiting for grafana to stop
			2019-08-08 17:10:22,730 INFO stopped: grafana (exit status 0)
			2019-08-08 17:10:23,267 INFO success: clickhouse entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:23,799 INFO spawned: 'grafana' with pid 13677
			2019-08-08 17:10:24,370 INFO waiting for qan-api2 to stop
			2019-08-08 17:10:24,397 INFO stopped: qan-api2 (exit status 0)
			2019-08-08 17:10:24,935 INFO success: grafana entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:25,473 INFO spawned: 'qan-api2' with pid 13751
			2019-08-08 17:10:26,024 INFO waiting for pmm-agent to stop
			2019-08-08 17:10:26,032 INFO stopped: pmm-agent (exit status 0)
			2019-08-08 17:10:26,557 INFO success: qan-api2 entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:27,083 INFO spawned: 'pmm-agent' with pid 13828
			2019-08-08 17:10:27,686 INFO spawned: 'dashboard-upgrade' with pid 13888
			2019-08-08 17:10:27,686 INFO success: dashboard-upgrade entered RUNNING state, process has stayed up for > than 0 seconds (startsecs)
			2019-08-08 17:10:27,761 INFO exited: dashboard-upgrade (exit status 0; expected)
			2019-08-08 17:10:28,149 INFO success: pmm-agent entered RUNNING state, process has stayed up for > than 1 seconds (startsecs)
			2019-08-08 17:10:28,975 INFO exited: pmm-update-perform (exit status 0; expected)
		`, "\n")

		var actual []*event
		for _, line := range log {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if e := parseEvent(line); e != nil {
				actual = append(actual, e)
			}
		}
		expected := []*event{
			{Time: time.Date(2019, 8, 8, 17, 9, 48, 494000000, time.UTC), Type: exitedUnexpected, Program: "pmm-update-perform"},
			{Time: time.Date(2019, 8, 8, 17, 10, 27, 761000000, time.UTC), Type: exitedExpected, Program: "dashboard-upgrade"},
			{Time: time.Date(2019, 8, 8, 17, 10, 28, 975000000, time.UTC), Type: exitedExpected, Program: "pmm-update-perform"},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("Fatal", func(t *testing.T) {
		t.Parallel()

		log := strings.Split(`
			2019-08-09 09:18:25,667 INFO spawned: 'pmm-update-check' with pid 11410
			2019-08-09 09:18:26,539 INFO exited: pmm-update-check (exit status 0; not expected)
			2019-08-09 09:18:27,543 INFO spawned: 'pmm-update-check' with pid 11421
			2019-08-09 09:18:28,324 INFO exited: pmm-update-check (exit status 0; not expected)
			2019-08-09 09:18:30,335 INFO spawned: 'pmm-update-check' with pid 11432
			2019-08-09 09:18:31,109 INFO exited: pmm-update-check (exit status 0; not expected)
			2019-08-09 09:18:34,119 INFO spawned: 'pmm-update-check' with pid 11443
			2019-08-09 09:18:34,883 INFO exited: pmm-update-check (exit status 0; not expected)
			2019-08-09 09:18:35,885 INFO gave up: pmm-update-check entered FATAL state, too many start retries too quickly
		`, "\n")

		var actual []*event
		for _, line := range log {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if e := parseEvent(line); e != nil {
				actual = append(actual, e)
			}
		}
		expected := []*event{
			{Time: time.Date(2019, 8, 9, 9, 18, 26, 539000000, time.UTC), Type: exitedUnexpected, Program: "pmm-update-check"},
			{Time: time.Date(2019, 8, 9, 9, 18, 28, 324000000, time.UTC), Type: exitedUnexpected, Program: "pmm-update-check"},
			{Time: time.Date(2019, 8, 9, 9, 18, 31, 109000000, time.UTC), Type: exitedUnexpected, Program: "pmm-update-check"},
			{Time: time.Date(2019, 8, 9, 9, 18, 34, 883000000, time.UTC), Type: exitedUnexpected, Program: "pmm-update-check"},
			{Time: time.Date(2019, 8, 9, 9, 18, 35, 885000000, time.UTC), Type: fatal, Program: "pmm-update-check"},
		}
		assert.Equal(t, expected, actual)
	})
}
