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

package models

import (
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

// FindBackupLocations returns saved backup locations configuration.
func FindBackupLocations(q *reform.Querier) ([]*BackupLocation, error) {
	rows, err := q.SelectAllFrom(BackupLocationTable, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select backup locations")
	}

	locations := make([]*BackupLocation, len(rows))
	for i, s := range rows {
		locations[i] = s.(*BackupLocation)
	}

	return locations, nil
}
