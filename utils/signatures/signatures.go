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

// Package signatures verifies signatures received from Percona Platform.
package signatures

import (
	"github.com/percona-platform/saas/pkg/check"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// VerifySignatures verifies checks signatures and returns error in case of verification problem.
func VerifySignatures(l *logrus.Entry, file string, signatures, publicKeys []string) error {
	if len(signatures) == 0 {
		return errors.New("zero signatures received")
	}

	var err error
	for _, sign := range signatures {
		for _, key := range publicKeys {
			if err = check.Verify([]byte(file), key, sign); err == nil {
				l.Debugf("Key %q matches signature %q.", key, sign)
				return nil
			}
			l.Debugf("Key %q doesn't match signature %q: %s.", key, sign, err)
		}
	}

	return errors.New("no verified signatures")
}
