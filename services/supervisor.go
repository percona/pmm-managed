// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package services

import (
	"context"

	servicelib "github.com/percona/kardianos-service"
)

//go:generate mockery -name=Supervisor
type Supervisor interface {
	// Start installs, and starts job
	Start(ctx context.Context, config *servicelib.Config) error
	// Stop stops job, and removes it
	Stop(ctx context.Context, name string) error
	// Status returns nil if job is running
	Status(ctx context.Context, name string) error
}
