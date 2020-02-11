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
	"errors"
	"reflect"
	"testing"
)

func TestEnvVarValidator(t *testing.T) {
	type args struct {
		env []string
	}
	tests := []struct {
		name        string
		args        args
		wantEnvVars map[string]string
		wantErrs    []error
		wantWarns   []string
	}{
		{
			"Valid env variables",
			args{[]string{
				"DISABLE_UPDATES=True",
				"DISABLE_TELEMETRY=False",
				"METRICS_RESOLUTION=5m",
				"METRICS_RESOLUTION_HR=5s",
				"METRICS_RESOLUTION_LR=1h",
				"DATA_RETENTION=72h",
			}},
			map[string]string{
				"DATA_RETENTION":        "72h",
				"DISABLE_TELEMETRY":     "false",
				"DISABLE_UPDATES":       "true",
				"METRICS_RESOLUTION":    "5m",
				"METRICS_RESOLUTION_HR": "5s",
				"METRICS_RESOLUTION_LR": "1h",
			},
			nil,
			nil,
		},
		{
			"Unknown env variables",
			args{[]string{"UNKNOWN_VAR=VAL", "ANOTHER_UNKNOWN_VAR=VAL"}},
			map[string]string{"ANOTHER_UNKNOWN_VAR": "val", "UNKNOWN_VAR": "val"},
			nil,
			[]string{
				`unknown environment variable "UNKNOWN_VAR=VAL"`,
				`unknown environment variable "ANOTHER_UNKNOWN_VAR=VAL"`,
			},
		},
		{
			"Invalid env variables values",
			args{[]string{
				"DISABLE_UPDATES=5",
				"DISABLE_TELEMETRY=X",
				"METRICS_RESOLUTION=5f",
				"METRICS_RESOLUTION_HR=s5",
				"METRICS_RESOLUTION_LR=1hour",
				"DATA_RETENTION=keep one week",
			}},
			map[string]string{
				"DATA_RETENTION":        "keep one week",
				"DISABLE_TELEMETRY":     "x",
				"DISABLE_UPDATES":       "5",
				"METRICS_RESOLUTION":    "5f",
				"METRICS_RESOLUTION_HR": "s5",
				"METRICS_RESOLUTION_LR": "1hour",
			},
			[]error{
				errors.New(`invalid environment variable "DISABLE_UPDATES=5"`),
				errors.New(`invalid environment variable "DISABLE_TELEMETRY=X"`),
				errors.New(`invalid environment variable "METRICS_RESOLUTION=5f"`),
				errors.New(`invalid environment variable "METRICS_RESOLUTION_HR=s5"`),
				errors.New(`invalid environment variable "METRICS_RESOLUTION_LR=1hour"`),
				errors.New(`invalid environment variable "DATA_RETENTION=keep one week"`),
			},
			nil,
		},
		{
			"Default env vars",
			args{[]string{
				"PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin",
				"HOSTNAME=host",
				"TERM=xterm-256color",
				"HOME=/home/user/",
			}},
			map[string]string{
				"HOME":     "/home/user/",
				"HOSTNAME": "host",
				"PATH":     "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin",
				"TERM":     "xterm-256color",
			},
			nil,
			nil,
		},
		{
			"Grafana env vars",
			args{[]string{
				`GF_AUTH_GENERIC_OAUTH_ALLOWED_DOMAINS='example.com'`,
				`GF_AUTH_GENERIC_OAUTH_ENABLED='true'`,
				`GF_PATHS_CONFIG="/etc/grafana/grafana.ini"`,
				`GF_PATHS_DATA="/var/lib/grafana"`,
				`GF_PATHS_HOME="/usr/share/grafana"`,
				`GF_PATHS_LOGS="/var/log/grafana"`,
				`GF_PATHS_PLUGINS="/var/lib/grafana/plugins"`,
				`GF_PATHS_PROVISIONING="/etc/grafana/provisioning"`,
			}},
			map[string]string{
				"GF_AUTH_GENERIC_OAUTH_ALLOWED_DOMAINS": "'example.com'",
				"GF_AUTH_GENERIC_OAUTH_ENABLED":         "'true'",
				"GF_PATHS_CONFIG":                       `"/etc/grafana/grafana.ini"`,
				"GF_PATHS_DATA":                         `"/var/lib/grafana"`,
				"GF_PATHS_HOME":                         `"/usr/share/grafana"`,
				"GF_PATHS_LOGS":                         `"/var/log/grafana"`,
				"GF_PATHS_PLUGINS":                      `"/var/lib/grafana/plugins"`,
				"GF_PATHS_PROVISIONING":                 `"/etc/grafana/provisioning"`,
			},
			nil,
			nil,
		},
		{
			"Warnings",
			args{[]string{
				"DATA_RETENTION=72m",
			}},
			map[string]string{
				"DATA_RETENTION": "72m",
			},
			nil,
			[]string{
				`retention period with the value less than a day can be wrong ("DATA_RETENTION=72m")`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnvVars, gotErrs, gotWarns := EnvVarValidator(tt.args.env)
			if !reflect.DeepEqual(gotEnvVars, tt.wantEnvVars) {
				t.Errorf("EnvVarValidator() gotEnvVars = %v, want %v", gotEnvVars, tt.wantEnvVars)
			}
			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("EnvVarValidator() gotErrs = %v, want %v", gotErrs, tt.wantErrs)
			}
			if !reflect.DeepEqual(gotWarns, tt.wantWarns) {
				t.Errorf("EnvVarValidator() gotWarns = %v, want %v", gotWarns, tt.wantWarns)
			}
		})
	}
}
