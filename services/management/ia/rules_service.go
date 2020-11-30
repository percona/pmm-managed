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

package ia

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// RulesService represents API for Integrated Alerting Rules.
type RulesService struct {
	db *reform.DB
}

// NewRulesService creates an API for Integrated Alerting Rules.
func NewRulesService(db *reform.DB) *RulesService {
	return &RulesService{
		db: db,
	}
}

// ListAlertRules returns a list of all Integrated Alerting rules.
func (s *RulesService) ListAlertRules(ctx context.Context, req *iav1beta1.ListAlertRulesRequest) (*iav1beta1.ListAlertRulesResponse, error) {
	var rules []models.Rule
	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		rules, err = models.FindRules(tx.Querier)
		return err
	})
	if e != nil {
		return nil, e
	}

	res := make([]*iav1beta1.Rule, len(rules))
	for i, rule := range rules {
		createdAt, err := ptypes.TimestampProto(rule.CreatedAt)
		if err != nil {
			return nil, err
		}

		r := &iav1beta1.Rule{
			RuleId:    rule.ID,
			Disabled:  rule.Disabled,
			Summary:   rule.Summary,
			For:       ptypes.DurationProto(rule.For),
			CreatedAt: createdAt,
			// TODO return updated_at
		}

		// template, params and severity

		var labels map[string]string
		err = json.Unmarshal(rule.CustomLabels, &labels)
		if err != nil {
			return nil, err
		}
		r.CustomLabels = labels

		filters := make([]*iav1beta1.Filter, len(rule.Filters))
		for _, filter := range rule.Filters {
			f := &iav1beta1.Filter{
				Type:  iav1beta1.FilterType(filter.Type),
				Key:   filter.Key,
				Value: filter.Val,
			}
			filters = append(filters, f)
		}
		r.Filters = filters

		channels := make([]*iav1beta1.Channel, len(rule.Channels))
		for i, channel := range rule.Channels {
			c, err := makeChannel(channel)
			if err != nil {
				// TODO
				return nil, err
			}
			channels[i] = c
		}
		r.Channels = channels

		res[i] = r
	}
	return &iav1beta1.ListAlertRulesResponse{Rules: res}, nil
}

func makeTemplate(template *models.Template) *iav1beta1.Template {
	t := &iav1beta1.Template{
		Name:     template.Name,
		Summary:  template.Summary,
		Expr:     template.Expr,
		Severity: managementpb.Severity(managementpb.Severity_value[template.Severity]),
		For:      ptypes.DurationProto(time.Duration(template.For)),
	}
}

// CreateAlertRule creates Integrated Alerting rule.
func (s *RulesService) CreateAlertRule(ctx context.Context, req *iav1beta1.CreateAlertRuleRequest) (*iav1beta1.CreateAlertRuleResponse, error) {
	params := &models.CreateRuleParams{
		TemplateName: req.TemplateName,
		Disabled:     req.Disabled,
		RuleParams:   req.Params,
		For:          req.For,
		Severity:     req.Severity,
		CustomLabels: req.CustomLabels,
		ChannelIDs:   req.ChannelIds,
	}

	filters := make([]*models.Filter, len(req.Filters))
	for _, filter := range req.Filters {
		filters = append(filters, &models.Filter{
			Type: models.FilterType(filter.Type),
			Key:  filter.Key,
			Val:  filter.Value,
		})
	}

	var rule *models.Rule
	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		rule, err = models.CreateRule(tx.Querier, params)
		return err
	})
	if e != nil {
		return nil, e
	}
	return &iav1beta1.CreateAlertRuleResponse{RuleId: rule.ID}, nil
}

// UpdateAlertRule updates Integrated Alerting rule.
func (s *RulesService) UpdateAlertRule(ctx context.Context, req *iav1beta1.UpdateAlertRuleRequest) (*iav1beta1.UpdateAlertRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAlertRule not implemented")
}

// ToggleAlertRule allows to switch between disabled and enabled states of an Alert Rule.
func (s *RulesService) ToggleAlertRule(ctx context.Context, req *iav1beta1.ToggleAlertRuleRequest) (*iav1beta1.ToggleAlertRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToggleAlertRule not implemented")
}

// DeleteAlertRule deletes Integrated Alerting rule.
func (s *RulesService) DeleteAlertRule(ctx context.Context, req *iav1beta1.DeleteAlertRuleRequest) (*iav1beta1.DeleteAlertRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAlertRule not implemented")
}

// Check interfaces.
var (
	_ iav1beta1.RulesServer = (*RulesService)(nil)
)
