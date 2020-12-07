package dir

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
			}
		})
	}
}
