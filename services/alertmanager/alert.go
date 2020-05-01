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

package alertmanager

import (
	"github.com/percona-platform/saas/pkg/check"
	"github.com/pkg/errors"

	"github.com/percona/pmm-managed/models"
)

// AlertParams defines alert parameters.
type AlertParams struct {
	Name        string
	Summary     string
	Description string
	Severity    check.Severity

	Node    *models.Node
	Service *models.Service
	Agent   *models.Agent
}

// validate checks parameters and fills defaults.
func (ap *AlertParams) validate() error {
	if ap.Name == "" {
		return errors.New("empty Name")
	}
	if ap.Summary == "" {
		return errors.New("empty Summary")
	}
	if ap.Description == "" {
		return errors.New("empty Description")
	}

	if ap.Severity < check.Emergency || ap.Severity > check.Debug {
		return errors.Errorf("invalid severity level: %s", ap.Severity)
	}

	return nil
}
