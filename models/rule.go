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
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
)

//go:generate reform

// alertRule represents an IA rule to be stored in the database.
//reform:alert_rules
type alertRule struct {
	Template     *[]byte `reform:"template"`
	ID           string  `reform:"id,pk"`
	Summary      string  `reform:"summary"`
	Disabled     bool    `reform:"disabled"`
	Params       *[]byte `reform:"params"`
	For          string  `reform:"for"`
	Severity     string  `reform:"severity"`
	CustomLabels *[]byte `reform:"custom_labels"`
	Filters      *[]byte `reform:"filters"`
	Channels     *[]byte `reform:"channels"`
	CreatedAt    string  `reform:"created_at"`
}

// Rule represents alertRule configuration.
type Rule struct {
	Template     *iav1beta1.Template
	ID           string
	Summary      string
	Disabled     bool
	Params       *[]iav1beta1.RuleParam
	For          *duration.Duration
	Severity     string
	CustomLabels map[string]string
	Filters      *[]Filter
	Channels     *[]iav1beta1.Channel // replace with Channels type defined in models
	CreatedAt    *timestamp.Timestamp
}

type FilterType int32

const (
	Invalid FilterType = 0
	// =
	equal FilterType = 1
	// !=
	NotEqual FilterType = 2
	// =~
	Regex FilterType = 3
	// !~
	NotRegex FilterType = 4
)

type Filter struct {
	Type  FilterType `json:"type"`
	Key   string     `json:"key"`
	Value string     `json:"value"`
}
