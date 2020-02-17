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

// EnvVarValidator validates given environment variables.
// Returns two lists with errors and warnings.
func EnvVarValidator(envs []string) (envSettings EnvSettings, errs []error, warns []string) {
	for _, env := range envs {
		p := strings.SplitN(env, "=", 2)
		if len(p) != 2 {
			errs = append(errs, fmt.Errorf("failed to parse environment variable %q", env))
			continue
		}

		k, v := strings.ToUpper(p[0]), strings.ToLower(p[1])
		switch k {
		case "PATH", "HOSTNAME", "TERM", "HOME":
		case "DISABLE_UPDATES":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", env))
				continue
			}
			envSettings.DisableUpdates = b
		case "DISABLE_TELEMETRY":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", env))
				continue
			}
			envSettings.DisableTelemetry = b
		case "METRICS_RESOLUTION", "METRICS_RESOLUTION_HR":
			d, err := time.ParseDuration(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", env))
				continue
			}
			envSettings.MetricsResolutions.HR = d
		case "METRICS_RESOLUTION_MR":
			d, err := time.ParseDuration(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", env))
				continue
			}
			envSettings.MetricsResolutions.MR = d
		case "METRICS_RESOLUTION_LR":
			d, err := time.ParseDuration(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", env))
				continue
			}
			envSettings.MetricsResolutions.LR = d
		case "DATA_RETENTION":
			d, err := time.ParseDuration(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", env))
			} else if d < 24*time.Hour {
				errs = append(errs, fmt.Errorf("data_retention: minimal resolution is 24h. received: %q", env))
			} else if d.Truncate(24*time.Hour) != d {
				errs = append(errs, fmt.Errorf("data_retention: should be a natural number of days. received: %q", env))
			} else {
				envSettings.DataRetention = d
			}
		default:
			if !strings.HasPrefix(k, "GF_") {
				warns = append(warns, fmt.Sprintf("unknown environment variable %q", env))
			}
		}
	}
	return envSettings, errs, warns
}
