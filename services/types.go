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

package services

import (
	"github.com/percona-platform/saas/pkg/check"

	"github.com/percona/pmm-managed/models"
)

// Target contains required info about STT check target.
type Target struct {
	AgentID       string
	ServiceID     string
	ServiceName   string
	Labels        map[string]string
	DSN           string
	Files         map[string]string
	TDP           *models.DelimiterPair
	TLSSkipVerify bool
}

// STTCheckResult contains the output from the check file and other information.
type STTCheckResult struct {
	CheckName string
	Silenced  bool
	AlertID   string
	Interval  check.Interval
	Target    Target
	Result    check.Result
}

// CheckSummary contains the summary of failed checks for a service.
type CheckSummary struct {
	ServiceName   string
	ServiceID     string
	CriticalCount uint32
	MajorCount    uint32
	TrivialCount  uint32
}
