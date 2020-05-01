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
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
)

func parseVersion(args ...interface{}) (interface{}, error) {
	if l := len(args); l != 1 {
		return nil, errors.Errorf("expected 1 argument, got %d", l)
	}

	s, ok := args[0].(string)
	if !ok {
		return nil, errors.Errorf("expected string argument, got %[1]T (%[1]v)", args[0])
	}

	p, err := version.Parse(s)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"major": p.Major,
		"minor": p.Minor,
		"patch": p.Patch,
		"rest":  p.Rest,
		"num":   p.Num,
	}, nil
}

func formatVersion(args ...interface{}) (interface{}, error) {
	if l := len(args); l != 1 {
		return nil, errors.Errorf("expected 1 argument, got %d", l)
	}

	d, ok := args[0].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("expected dict argument, got %[1]T (%[1]v)", args[0])
	}

	// FIXME handle type assertion panics
	p := &version.Parsed{
		Major: d["major"].(int64),
		Minor: d["minor"].(int64),
		Patch: d["patch"].(int64),
		Rest:  d["rest"].(string),
		Num:   d["num"].(int64),
	}
	return p.String(), nil
}
