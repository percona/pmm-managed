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

// Package validators contains environment variables validator.
package validators

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MetricsResolutions contains standard Prometheus metrics resolutions.
type MetricsResolutions struct {
	HR time.Duration
	MR time.Duration
	LR time.Duration
}

// EnvSettings contains PMM Server settings.
type EnvSettings struct {
	DisableUpdates     bool
	DisableTelemetry   bool
	MetricsResolutions MetricsResolutions
	DataRetention      time.Duration
}

// ValidateEnvVars validates given environment variables.
//
// Returns valid setting and two lists with errors and warnings.
// This function is mainly used in pmm-managed-init to early validate passed
// environment variables, and provide user warnings about unknown variables.
// In case of error, the docker run terminates.
// Short description of environment variables:
//  - PATH, HOSTNAME, TERM, HOME are default environment variables that will be ignored;
//  - DISABLE_UPDATES is a boolean flag to enable or disable pmm-server update;
//  - DISABLE_TELEMETRY is a boolean flag to enable or disable pmm telemetry;
//  - METRICS_RESOLUTION, METRICS_RESOLUTION, METRICS_RESOLUTION_HR,
// METRICS_RESOLUTION_LR are durations of metrics resolution;
//  - DATA_RETENTION is the duration of how long keep time-series data in ClickHouse;
//  - the environment variables prefixed with GF_ passed as related to Grafana.
func ValidateEnvVars(envs []string) (envSettings EnvSettings, errs []error, warns []string) {
	for _, env := range envs {
		p := strings.SplitN(env, "=", 2)
		if len(p) != 2 {
			errs = append(errs, fmt.Errorf("failed to parse environment variable %q", env))
			continue
		}

		var err error
		k, v := strings.ToUpper(p[0]), strings.ToLower(p[1])
		switch k {
		// Skip default environment variables.
		case "PATH", "HOSTNAME", "TERM", "HOME":
		case "DISABLE_UPDATES":
			envSettings.DisableUpdates, err = strconv.ParseBool(v)
			if err != nil {
				err = fmt.Errorf("invalid environment variable %q", env)
			}
		case "DISABLE_TELEMETRY":
			envSettings.DisableTelemetry, err = strconv.ParseBool(v)
			if err != nil {
				err = fmt.Errorf("invalid environment variable %q", env)
			}
		case "METRICS_RESOLUTION", "METRICS_RESOLUTION_HR":
			envSettings.MetricsResolutions.HR, err = validateDuration(v, env, time.Second, time.Second)
		case "METRICS_RESOLUTION_MR":
			envSettings.MetricsResolutions.MR, err = validateDuration(v, env, time.Second, time.Second)
		case "METRICS_RESOLUTION_LR":
			envSettings.MetricsResolutions.LR, err = validateDuration(v, env, time.Second, time.Second)
		case "DATA_RETENTION":
			envSettings.DataRetention, err = validateDuration(v, env, 24*time.Hour, 24*time.Hour)
		default:
			if !strings.HasPrefix(k, "GF_") {
				warns = append(warns, fmt.Sprintf("unknown environment variable %q", env))
			}
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return envSettings, errs, warns
}

func validateDuration(value, env string, min, multipleOf time.Duration) (time.Duration, error) {
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q has invalid duration %v", env, value)
	} else if d < min {
		return 0, fmt.Errorf("environment variable %q cannot be less then %s", env, min)
	} else if d.Truncate(multipleOf) != d {
		return 0, fmt.Errorf("environment variable %q should be a natural number of %s", env, multipleOf)
	}
	return d, nil
}
