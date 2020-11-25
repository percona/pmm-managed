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

	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	reform "gopkg.in/reform.v1"
)

// SaveRule persists alert rule.
func SaveRule(q *reform.Querier, r *Rule) error {
	if err := ValidateRule(r); err != nil {
		return err
	}

	nc, err := ruleToAlertRule(r)
	if err != nil {
		return err
	}

	err = q.Insert(nc)
	if err != nil {
		return errors.Wrap(err, "failed to create alert rule")
	}

	return nil
}

// UpdateRule updates existing alert rule.
func UpdateRule(q *reform.Querier, r *Rule) error {
	if err := ValidateRule(r); err != nil {
		return errors.Wrap(err, "channel validation failed")
	}

	nc, err := ruleToAlertRule(r)
	if err != nil {
		return err
	}

	err = q.Update(nc)
	if err != nil {
		return errors.Wrap(err, "failed to create alert rule")
	}

	return nil
}

// RemoveRule removes alert rule with specified id.
func RemoveRule(q *reform.Querier, id string) error {
	err := q.Delete(&alertRule{ID: id})
	if err != nil {
		return errors.Wrap(err, "failed to delete alert rule")
	}
	return nil
}

// GetRules returns saved alert rules configuration.
func GetRules(q *reform.Querier) ([]Rule, error) {
	structs, err := q.SelectAllFrom(alertRulesTable, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select alert rules")

	}

	rules := make([]Rule, len(structs))
	for i, s := range structs {
		c, err := alertRuleToRule(s.(*alertRule))
		if err != nil {
			return nil, err
		}
		rules[i] = *c
	}

	return rules, nil
}

// ValidateRule validates alert rule.
func ValidateRule(r *Rule) error {
	if r.ID == "" {
		return status.Error(codes.InvalidArgument, "alert rule id is empty")
	}
	return nil
}

func ruleToAlertRule(r *Rule) (*alertRule, error) {
	ar := &alertRule{
		ID:        r.ID,
		Type:      r.Type,
		Disabled:  r.Disabled,
		For:       r.For.String(),
		CreatedAt: r.CreatedAt.String(),
	}

	t, err := json.Marshal(r.Template)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall template")
	}
	ar.Template = &t

	p, err := json.Marshal(r.Params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall params")
	}
	ar.Params = &p

	cl, err := json.Marshal(r.CustomLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall custom labels")
	}
	ar.CustomLabels = &cl

	f, err := json.Marshal(r.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall filters")
	}
	ar.Filters = &f

	c, err := json.Marshal(r.Channels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall channels")
	}
	ar.Channels = &c

	return ar, nil
}

func alertRuleToRule(ar *alertRule) (*Rule, error) {
	r := &Rule{
		ID:       ar.ID,
		Type:     ar.Type,
		Disabled: ar.Disabled,
	}

	r.Template = &iav1beta1.Template{}
	err := json.Unmarshal(*ar.Template, r.Template)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall template")
	}

	r.Params = &[]iav1beta1.Params{}
	err = json.Unmarshal(*ar.Params, r.Params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall params")
	}

	r.Filters = &[]Filter{}
	err = json.Unmarshal(*ar.Filters, r.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall filters")
	}

	r.Channels = &[]iav1beta1.Channels{}
	err = json.Unmarshal(*ar.Channels, r.Channels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall channels")
	}

	return r, nil
}
