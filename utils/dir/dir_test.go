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

package dir

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDataDir(t *testing.T) {
	testcases := []struct {
		name   string
		params Params
		err    string
	}{{
		name: "valid params",
		params: Params{
			Path:  "/tmp/testdir",
			Perm:  os.FileMode(0o775),
			User:  "pmm",
			Group: "pmm",
		},
		err: "",
	}, {
		name: "unknown user",
		params: Params{
			Path:  "/tmp/testdir",
			Perm:  os.FileMode(0o775),
			User:  "$",
			Group: "pmm",
		},
		err: "cannot chown datadir user: unknown user $",
	}, {
		name: "unknown group",
		params: Params{
			Path:  "/tmp/testdir",
			Perm:  os.FileMode(0o775),
			User:  "pmm",
			Group: "$",
		},
		err: "cannot chown datadir group: unknown group $",
	},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateDataDir(tc.params)
			if tc.err != "" {
				assert.Equal(t, tc.err, actual.Error())
			} else {
				stat, err := os.Stat(tc.params.Path)
				require.NoError(t, err)
				assert.True(t, stat.IsDir())
				assert.Equal(t, tc.params.Perm, stat.Mode().Perm())
			}
		})
	}
}
