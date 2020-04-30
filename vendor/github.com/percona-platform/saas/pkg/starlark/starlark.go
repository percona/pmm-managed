// Package starlark is executor for starklark.
package starlark

import (
	"github.com/pkg/errors"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/percona-platform/saas/pkg/check"
)

// Run executes the script with given name and input data.
func Run(name, script string, input []map[string]interface{}) ([]check.Result, error) {
	thread := &starlark.Thread{
		Name: name,
	}

	rows, err := prepareRows(input)
	if err != nil {
		return nil, err
	}

	globals, err := starlark.ExecFile(thread, script, []byte(script), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute starlark script")
	}

	f := globals["check"]
	if f == nil {
		return nil, errors.New("check function is not defined")
	}

	v, err := starlark.Call(thread, f, starlark.Tuple{rows}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute check function")
	}

	return parseScriptOutput(v)
}

func prepareRows(input []map[string]interface{}) (starlark.Tuple, error) {
	rows := make(starlark.Tuple, len(input))
	for i, v := range input {
		sv, err := goToStarlark(v)
		if err != nil {
			return nil, err
		}
		rows[i] = sv
	}
	rows.Freeze()

	return rows, nil
}

// modify unavoidable global state once on package initialization to avoid race conditions
//nolint:gochecknoinits
func init() {
	resolve.AllowFloat = true
	resolve.AllowSet = true
}
