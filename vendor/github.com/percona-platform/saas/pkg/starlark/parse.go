package starlark

import (
	"github.com/pkg/errors"
	"go.starlark.net/starlark"

	"github.com/percona-platform/saas/pkg/check"
)

// parseScriptOutput parses and validates script output and returns slice of Results.
func parseScriptOutput(v starlark.Value) ([]check.Result, error) {
	switch v := v.(type) {
	case starlark.Tuple:
		if v.Len() != 2 {
			return nil, errors.New("script has invalid output")
		}

		errMsg, err := parseErrorMessage(v.Index(1))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse error message")
		}

		if errMsg != "" {
			return nil, errors.Errorf("script error: %s", errMsg)
		}

		results, err := parseResults(v.Index(0))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse results list")
		}

		for _, result := range results {
			if err = result.Validate(); err != nil {
				return nil, err
			}
		}

		return results, nil
	default:
		return nil, errors.Errorf("unhandled result type %T", v)
	}
}

// parseResults returns slice of results parsed from starlark value.
func parseResults(v starlark.Value) ([]check.Result, error) {
	val, err := starlarkToGo(v)
	if err != nil {
		return nil, err
	}

	rs, ok := val.([]interface{})
	if !ok {
		return nil, errors.Errorf("results list has wrong type: %T", val)
	}

	results := make([]check.Result, len(rs))
	for i, r := range rs {
		m := r.(map[string]interface{})
		var sum, desc string
		var sev check.Severity

		if v, ok := m["summary"]; ok {
			if sum, ok = v.(string); !ok {
				return nil, errors.Errorf("summary field has wrong type: %T", v)
			}
		}

		if v, ok := m["description"]; ok {
			if desc, ok = v.(string); !ok {
				return nil, errors.Errorf("description field has wrong type: %T", v)
			}
		}

		if v, ok := m["severity"]; ok {
			sevS, ok := v.(string)
			if !ok {
				return nil, errors.Errorf("severity field has wrong type: %T", v)
			}

			sev = check.StrToSeverity(sevS)
		}

		results[i] = check.Result{
			Summary:     sum,
			Description: desc,
			Severity:    sev,
		}
	}

	return results, nil
}

// parseErrorMessage returns error message parsed from starlark value.
func parseErrorMessage(v starlark.Value) (string, error) {
	val, err := starlarkToGo(v)
	if err != nil {
		return "", err
	}

	str, ok := val.(string)
	if !ok {
		return "", errors.Errorf("error message has wrong type: %T", val)
	}

	return str, nil
}
