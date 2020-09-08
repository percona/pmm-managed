package main

import (
	"testing"

	"github.com/percona/pmm-managed/services/checks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRunChecks(t *testing.T) {
	l := logrus.WithField("component", "pmm-managed-starlark")

	testCases := []struct {
		data checks.StarlarkScriptData
		err  bool
	}{
		{
			data: checks.StarlarkScriptData{
				CheckName:         "invalid version",
				Script:            "def check(): return []",
				CheckVersion:      5,
				QueryActionResult: "some result",
			},
			err: true,
		},
		{
			data: checks.StarlarkScriptData{
				CheckName:         "invalid starlark syntax",
				Script:            "@ + @",
				CheckVersion:      1,
				QueryActionResult: "some result",
			},
			err: true,
		},
		{
			data: checks.StarlarkScriptData{
				CheckName:         "invalid query action result",
				Script:            "def check(): return []",
				CheckVersion:      1,
				QueryActionResult: "some result",
			},
			err: true,
		},
		{
			data: checks.StarlarkScriptData{
				CheckName:         "bad starlark script",
				Script:            "def check(rows): return [1] * (1 << 30-1)",
				CheckVersion:      1,
				QueryActionResult: "\n\rVariable_name\n\x05Value\x12\x1c\n\t2\aversion\n\x0f2\r5.7.30-33-log\x12I\n\x112\x0fversion_comment\n422Percona Server (GPL), Release 33, Revision 6517692\x12%\n\x192\x17version_compile_machine\n\b2\x06x86_64\x12\x1f\n\x142\x12version_compile_os\n\a2\x05Linux\x12\x1a\n\x102\x0eversion_suffix\n\x062\x04-log",
			},
			err: true,
		},
		{
			data: checks.StarlarkScriptData{
				CheckName:         "valid starlark script",
				Script:            "def check(rows): return []",
				CheckVersion:      1,
				QueryActionResult: "\n\rVariable_name\n\x05Value\x12\x1c\n\t2\aversion\n\x0f2\r5.7.30-33-log\x12I\n\x112\x0fversion_comment\n422Percona Server (GPL), Release 33, Revision 6517692\x12%\n\x192\x17version_compile_machine\n\b2\x06x86_64\x12\x1f\n\x142\x12version_compile_os\n\a2\x05Linux\x12\x1a\n\x102\x0eversion_suffix\n\x062\x04-log",
			},
			err: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.data.CheckName, func(t *testing.T) {
			err := runChecks(l, tc.data)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			t.Log(err)
		})
	}
}
