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
	"testing"

	"github.com/percona/pmm/api/agentpb"
	"github.com/stretchr/testify/require"

	"github.com/percona/pmm-managed/services/checks"
)

var (
	validQueryActionResult = []map[string]interface{}{
		{"Value": "5.7.30-33-log", "Variable_name": "version"},
		{"Value": "Percona Server (GPL), Release 33, Revision 6517692", "Variable_name": "version_comment"},
		{"Value": "x86_64", "Variable_name": "version_compile_machine"},
		{"Value": "Linux", "Variable_name": "version_compile_os"},
		{"Value": "-log", "Variable_name": "version_suffix"},
	}

	invalidQueryActionResult = []map[string]interface{}{
		{"key-1": "val-1", "key-2": "val-2"},
	}
)

func TestRunChecks(t *testing.T) {
	testCases := []struct {
		err     bool
		version uint32
		name    string
		script  string
		result  []map[string]interface{}
	}{
		{
			err:     true,
			version: 5,
			name:    "invalid version",
			script:  "def check(): return []",
			result:  invalidQueryActionResult,
		},
		{
			err:     true,
			version: 1,
			name:    "invalid starlark syntax",
			script:  "@ + @",
			result:  invalidQueryActionResult,
		},
		{
			err:     true,
			version: 1,
			name:    "invalid query action result",
			script:  "def check(): return []",
			result:  invalidQueryActionResult,
		},
		{
			err:     true,
			version: 1,
			name:    "bad starlark script",
			script:  "def check(rows): return [1] * (1 << 30-1)",
			result:  validQueryActionResult,
		},
		{
			err:     false,
			version: 1,
			name:    "valid starlark script",
			script:  "def check(rows): return []",
			result:  validQueryActionResult,
		},
	}

	err := exec.Command("/bin/sh", "-c", "cd ../../; make release").Run()
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

			var stdin bytes.Buffer
			cmd.Stdin = &stdin

			encoder := json.NewEncoder(&stdin)
			err = encoder.Encode(data)
			require.NoError(t, err)

			err = cmd.Run()
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
