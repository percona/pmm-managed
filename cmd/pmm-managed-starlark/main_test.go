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

package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"

	"github.com/percona/pmm/api/agentpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/pmm-managed/services/checks"
)

const (
	invalidStarlarkScriptStderr = "Error running starlark script: thread invalid starlark script: failed to execute function check: function check accepts no arguments (1 given)"
	invalidVersionStderr        = "Error running starlark script: unsupported check version: 5"
	memoryConsumingScriptStderr = "fatal error: runtime: out of memory"
)

var (
	validQueryActionResult = []map[string]interface{}{
		{"Value": "5.7.30-33-log", "Variable_name": "version"},
		{"Value": "Percona Server (GPL), Release 33, Revision 6517692", "Variable_name": "version_comment"},
		{"Value": "x86_64", "Variable_name": "version_compile_machine"},
		{"Value": "Linux", "Variable_name": "version_compile_os"},
		{"Value": "-log", "Variable_name": "version_suffix"},
	}
)

func TestRunChecks(t *testing.T) {
	testCases := []struct {
		err          bool
		version      uint32
		name         string
		script       string
		errorMessage string
		stderr       string
		stdout       string
		result       []map[string]interface{}
	}{
		{
			err:          true,
			version:      1,
			name:         "invalid starlark script",
			errorMessage: "exit status 1",
			stderr:       invalidStarlarkScriptStderr,
			stdout:       "",
			script:       "def check(): return []",
			result:       validQueryActionResult,
		},
		{
			err:          true,
			version:      5,
			name:         "invalid version",
			errorMessage: "exit status 1",
			stderr:       invalidVersionStderr,
			stdout:       "",
			script:       "def check(): return []",
			result:       validQueryActionResult,
		},
		{
			err:          true,
			version:      1,
			name:         "memory consuming starlark script",
			script:       "def check(rows): return [1] * (1 << 30-1)",
			errorMessage: "exit status 2",
			stderr:       memoryConsumingScriptStderr,
			stdout:       "",
			result:       validQueryActionResult,
		},
		{
			err:     true,
			version: 1,
			name:    "cpu consuming starlark script",
			script: `def check(rows):
						while True:
							pass`,
			errorMessage: "signal: killed",
			stderr:       "",
			stdout:       "",
			result:       validQueryActionResult,
		},
		{
			err:     false,
			version: 1,
			name:    "valid starlark script",
			script: `def check(rows):
						results = []
						results.append({
							"summary": "Fake check",
							"description": "Fake check description",
							"severity": "warning",
						})
						return results`,
			errorMessage: "",
			stderr:       "",
			stdout:       "[{\"summary\":\"Fake check\",\"description\":\"Fake check description\",\"severity\":5,\"labels\":null}]\n",
			result:       validQueryActionResult,
		},
	}

	// since we run the binary as a child process to test it we need to build it first.
	err := exec.Command("make", "-C", "../..", "release").Run()
	require.NoError(t, err)

	for _, tc := range testCases {
		result, err := agentpb.MarshalActionQueryDocsResult(tc.result)
		require.NoError(t, err)

		data := checks.StarlarkScriptData{
			CheckName:         tc.name,
			CheckVersion:      tc.version,
			Script:            tc.script,
			QueryActionResult: result,
		}

		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("./../../bin/pmm-managed-starlark")

			var stdin, stderr bytes.Buffer
			cmd.Stdin = &stdin
			cmd.Stderr = &stderr
			cmd.Env = []string{"PERCONA_TEST_STARLARK_ALLOW_RECURSION=true"}

			encoder := json.NewEncoder(&stdin)
			err = encoder.Encode(data)
			require.NoError(t, err)

			stdout, err := cmd.Output()
			stderrContent := stderr.String()
			if tc.err {
				require.Error(t, err)
				require.Empty(t, tc.stdout)
				require.Equal(t, tc.errorMessage, err.Error())
				assert.True(t, strings.Contains(stderrContent, tc.stderr))
				// make sure that the limits were set
				assert.False(t, strings.Contains(stderrContent, cpuUsageWarning))
				assert.False(t, strings.Contains(stderrContent, memoryUsageWarning))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.stdout, string(stdout))
			}
		})
	}
}
