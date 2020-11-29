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

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
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
	rules, err := models.GetRules(s.db.Querier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get alert rules")
	}
	res := make([]*iav1beta1.Rule, len(rules))
	for i, rule := range rules {
		createdAt, err := ptypes.TimestampProto(rule.CreatedAt)
		if err != nil {
			return nil, err
		}

		r := &iav1beta1.Rule{
			RuleId:       rule.ID,
			Template:     rule.Template,
			Disabled:     rule.Disabled,
			Summary:      rule.Summary,
			Params:       rule.Params,
			For:          ptypes.DurationProto(rule.For),
			Severity:     rule.Severity,
			CustomLabels: rule.CustomLabels,
			Channels:     rule.Channels,
			CreatedAt:    createdAt,
		}

		filters := make([]*iav1beta1.Filter, len(rule.Filters))
		for _, filter := range rule.Filters {
			f := &iav1beta1.Filter{
				Type:  iav1beta1.FilterType(filter.Type),
				Key:   filter.Key,
				Value: filter.Value,
			}
			filters = append(filters, f)
		}
		r.Filters = filters
		res[i] = r
	}
	return &iav1beta1.ListAlertRulesResponse{Rules: res}, nil
}

// CreateAlertRule creates Integrated Alerting rule.
func (s *RulesService) CreateAlertRule(ctx context.Context, req *iav1beta1.CreateAlertRuleRequest) (*iav1beta1.CreateAlertRuleResponse, error) {
	r := &models.Rule{
		// Add Template, CreatedAt, etc
		ID:           "/ia/rule_id/" + uuid.New().String(),
		Disabled:     req.GetDisabled(),
		Params:       req.GetParams(),
		For:          req.For.AsDuration(),
		Severity:     req.GetSeverity(),
		CustomLabels: req.GetCustomLabels(),
	}
	err := models.SaveRule(s.db.Querier, r)
	if err != nil {
		return nil, err
	}
	return &iav1beta1.CreateAlertRuleResponse{}, nil
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
