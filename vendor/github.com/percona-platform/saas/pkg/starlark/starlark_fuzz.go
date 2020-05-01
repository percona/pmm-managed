// +build gofuzz

package starlark

import "encoding/json"

//nolint:gochecknoinits
func init() {
	doRecover = false
}

func Fuzz(b []byte) int {
	var data fuzzData
	if json.Unmarshal(b, &data) != nil {
		return 0
	}

	env, err := NewEnv("fuzz", string(data.Script), nil)
	if err != nil {
		return 0
	}

	if _, err := env.Run("id", data.Input, nil); err != nil {
		return 0
	}

	return 1
}
