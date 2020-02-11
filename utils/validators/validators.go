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

// EnvVarValidator validates given environment variables.
// Returns two lists with errors and warnings.
func EnvVarValidator(env []string) (envVars map[string]string, errs []error, warns []string) {
	envVars = make(map[string]string)
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		if len(p) != 2 {
			errs = append(errs, fmt.Errorf("failed to parse environment variable %q", e))
			continue
		}

		k, v := strings.ToUpper(p[0]), strings.ToLower(p[1])
		switch k {
		case "PATH", "HOSTNAME", "TERM", "HOME":
		case "DISABLE_UPDATES":
			if _, err := strconv.ParseBool(v); err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", e))
			}
		case "DISABLE_TELEMETRY":
			if _, err := strconv.ParseBool(v); err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", e))
			}
		case "METRICS_RESOLUTION", "METRICS_RESOLUTION_HR", "METRICS_RESOLUTION_MR", "METRICS_RESOLUTION_LR":
			if _, err := time.ParseDuration(v); err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", e))
			}
		case "DATA_RETENTION":
			d, err := time.ParseDuration(v)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid environment variable %q", e))
			} else if d < 24*time.Hour {
				warns = append(warns, fmt.Sprintf("retention period with the value less than a day can be wrong (%q)", e))
			}
		default:
			if !strings.HasPrefix(k, "GF_") {
				warns = append(warns, fmt.Sprintf("unknown environment variable %q", e))
			}
		}
		envVars[k] = v
	}
	return
}
