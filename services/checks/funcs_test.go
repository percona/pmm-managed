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
	"strings"
	"testing"

	"github.com/percona-platform/saas/pkg/check"
	"github.com/percona-platform/saas/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	script := strings.TrimSpace(`
def check(rows):
    v = parse_version(rows[0].get("version"))
    s = format_version_num(v["num"])
    return [{
		"summary": s,
		"severity": "warning",
        "labels": {
            "major": str(v["major"]),
            "minor": str(v["minor"]),
            "patch": str(v["patch"]),
            "rest":  str(v["rest"]),
            "num":   str(v["num"]),
        }
    }]
	`)
	funcs := map[string]starlark.GoFunc{
		"parse_version":      parseVersion,
		"format_version_num": formatVersionNum,
	}
	env, err := starlark.NewEnv(t.Name(), script, funcs)
	require.NoError(t, err)

	input := []map[string]interface{}{
		{"version": int64(1)},
	}
	res, err := env.Run("type", input, nil)
	expectedErr := strings.TrimSpace(`
thread type: failed to execute function check: parse_version: expected string argument, got int64 (1)
Traceback (most recent call last):
  TestVersion:2:22: in check
  <builtin>: in parse_version
	`) + "\n"
	assert.EqualError(t, err, expectedErr)
	assert.Empty(t, res)

	input = []map[string]interface{}{
		{"version": "foo"},
	}
	res, err = env.Run("foo", input, nil)
	expectedErr = strings.TrimSpace(`
thread foo: failed to execute function check: parse_version: failed to parse "foo"
Traceback (most recent call last):
  TestVersion:2:22: in check
  <builtin>: in parse_version
	`) + "\n"
	assert.EqualError(t, err, expectedErr)
	assert.Empty(t, res)

	input = []map[string]interface{}{
		{"version": "8.0.19-10"},
	}
	res, err = env.Run("valid", input, nil)
	require.NoError(t, err)
	expected := []check.Result{{
		Summary:  "8.0.19",
		Severity: check.Warning,
		Labels: map[string]string{
			"major": "8",
			"minor": "0",
			"patch": "19",
			"rest":  "-10",
			"num":   "80019",
		},
	}}
	assert.Equal(t, expected, res)
}
