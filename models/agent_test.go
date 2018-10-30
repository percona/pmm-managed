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
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestNameForSupervisor(t *testing.T) {
	for _, n := range []interface {
		NameForSupervisor() string
	}{
		&Agent{Type: "dummy-type", ListenPort: pointer.ToUint16(12345)},
		&MySQLdExporter{Type: "dummy-type", ListenPort: pointer.ToUint16(12345)},
		&PostgresExporter{Type: "dummy-type", ListenPort: pointer.ToUint16(12345)},
		&RDSExporter{Type: "dummy-type", ListenPort: pointer.ToUint16(12345)},
		&QanAgent{Type: "dummy-type", ListenPort: pointer.ToUint16(12345)},
	} {
		t.Run(fmt.Sprintf("%T", n), func(t *testing.T) {
			assert.Equal(t, "pmm-dummy-type-12345", n.NameForSupervisor())
		})
	}
}
