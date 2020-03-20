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

// Package validators contains settings validators.
package validators

import (
	"github.com/pkg/errors"
)

// ValidateMySQLConnectionOptions validates MySQL connection options.
func ValidateMySQLConnectionOptions(socket, host *string, port *uint16) error {
	if (host == nil || port == nil) && socket == nil {
		return errors.New("Address (with port) or socket is required.")
	}

	if (host != nil || port != nil) && socket != nil {
		return errors.New("Setting both address (with port) and socket in once is disallowed.")
	}

	if host != nil && port == nil {
		return errors.New("Port is required.")
	}

	return nil
}
