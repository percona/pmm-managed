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

package server

import (
	"context"
	"time"

	"github.com/percona/pmm/version"
)

//go:generate mockery -name=prometheusService -case=snake -inpkg -testonly
//go:generate mockery -name=supervisordService -case=snake -inpkg -testonly

// prometheusService is a subset of methods of prometheus.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type prometheusService interface {
	UpdateConfiguration()
	Check(ctx context.Context) error
}

// supervisordService is a subset of methods of supervisord.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type supervisordService interface {
	InstalledPackageInfo() *version.PackageInfo
	CheckResult() (*version.UpdateCheckResult, time.Time)
	Check() error

	StartPMMUpdate() (uint32, error)
	PMMUpdateRunning() bool
	PMMUpdateLog(offset uint32) ([]string, uint32, error)
}
