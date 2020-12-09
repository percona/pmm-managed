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
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/percona-platform/saas/pkg/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

func checkUniqueRuleID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Rule ID")
	}

	rule := &Rule{ID: id}
	switch err := q.Reload(rule); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Rule with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// FindRules returns saved alert rules configuration.
func FindRules(q *reform.Querier) ([]Rule, error) {
	rows, err := q.SelectAllFrom(RuleTable, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select alert rules")
	}

	rules := make([]Rule, len(rows))
	for i, s := range rows {
		c := s.(*Rule)

		rules[i] = *c
	}

	return rules, nil
}

// FindRuleByID finds Rule by ID.
func FindRuleByID(q *reform.Querier, id string) (*Rule, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Rule ID.")
	}

	rule := &Rule{ID: id}
	switch err := q.Reload(rule); err {
	case nil:
		return rule, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Rule with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// CreateRuleParams are params for creating new Rule.
type CreateRuleParams struct {
	TemplateName string
	Disabled     bool
	RuleParams   RuleParams
	For          time.Duration
	Severity     common.Severity
	CustomLabels map[string]string
	Filters      Filters
	ChannelIDs   []string
}

// CreateRule persists alert Rule.
func CreateRule(q *reform.Querier, params *CreateRuleParams) (*Rule, error) {
	id := "/rule_id/" + uuid.New().String()
	if err := checkUniqueRuleID(q, id); err != nil {
		return nil, err
	}

	channelIDs := deduplicateStrings(params.ChannelIDs)
	channels, err := FindChannelsByIDs(q, channelIDs)
	if err != nil {
		return nil, err
	}

	if len(channelIDs) != len(channels) {
		return nil, errors.Errorf("failed to find all required channels %v", channelIDs)
	}

	row := &Rule{
		ID:           id,
		TemplateName: params.TemplateName,
		Disabled:     params.Disabled,
		For:          params.For,
		Severity:     convertSeverity(params.Severity),
		Filters:      params.Filters,
		Params:       params.RuleParams,
		ChannelIDs:   params.ChannelIDs,
	}

	err = row.SetCustomLabels(params.CustomLabels)
	if err != nil {
		return nil, err
	}

	err = q.Insert(row)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create alert rule")
	}

	return row, nil
}

// ChangeRuleParams is params for updating existing Rule.
type ChangeRuleParams struct {
	Disabled     bool
	RuleParams   RuleParams
	For          time.Duration
	Severity     common.Severity
	CustomLabels map[string]string
	Filters      Filters
	ChannelIDs   []string
}

// ChangeRule updates existing alerts Rule.
func ChangeRule(q *reform.Querier, ruleID string, params *ChangeRuleParams) (*Rule, error) {
	row, err := FindRuleByID(q, ruleID)
	if err != nil {
		return nil, err
	}

	channelIDs := deduplicateStrings(params.ChannelIDs)
	channels, err := FindChannelsByIDs(q, channelIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels")
	}

	if len(channelIDs) != len(channels) {
		return nil, errors.Errorf("failed to find all required channels %v", channelIDs)
	}

	row.Disabled = params.Disabled
	row.For = params.For
	row.Severity = convertSeverity(params.Severity)
	row.Filters = params.Filters
	row.Params = params.RuleParams

	labels, err := json.Marshal(params.CustomLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update alert rule")
	}
	row.CustomLabels = labels
	row.ChannelIDs = params.ChannelIDs

	err = q.Update(row)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update alerts Rule")
	}

	return row, nil
}

// RemoveRule removes alert Rule with specified id.
func RemoveRule(q *reform.Querier, id string) error {
	if _, err := FindRuleByID(q, id); err != nil {
		return err
	}

	err := q.Delete(&Rule{ID: id})
	if err != nil {
		return errors.Wrap(err, "failed to delete alert Rule")
	}
	return nil
}
