package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/percona/pmm-managed/services/checks"
	"github.com/percona/pmm/api/agentpb"
	"github.com/stretchr/testify/require"
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
		name    string
		script  string
		version uint32
		result  []map[string]interface{}
		err     bool
	}{
		{
			name:    "invalid version",
			script:  "def check(): return []",
			version: 5,
			result:  invalidQueryActionResult,
			err:     true,
		},
		{
			name:    "invalid starlark syntax",
			script:  "@ + @",
			version: 1,
			result:  invalidQueryActionResult,
			err:     true,
		},
		{
			name:    "invalid query action result",
			script:  "def check(): return []",
			version: 1,
			result:  invalidQueryActionResult,
			err:     true,
		},
		{
			name:    "bad starlark script",
			script:  "def check(rows): return [1] * (1 << 30-1)",
			version: 1,
			result:  validQueryActionResult,
			err:     true,
		},
		{
			name:    "valid starlark script",
			script:  "def check(rows): return []",
			version: 1,
			result:  validQueryActionResult,
			err:     false,
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
