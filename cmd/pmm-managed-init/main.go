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
	"os"

	"github.com/percona/pmm-managed/utils/validators"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logrus.WithField("component", "pmm-managed-init")
	envVars := os.Environ()
	_, errs, warns := validators.EnvVarValidator(envVars)
	for _, warn := range warns {
		l.Warnln(warn)
	}
	for _, err := range errs {
		l.Errorln(err)
	}

	if len(errs) > 0 {
		os.Exit(1)
	}
}
