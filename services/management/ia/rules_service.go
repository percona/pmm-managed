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
	"errors"
	"fmt"
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
			Severity:  managementpb.Severity(managementpb.Severity_value[rule.Severity]),
			For:       ptypes.DurationProto(rule.For),
			CreatedAt: createdAt,
		}

		template, err := makeTemplate(rule.Template)
		if err != nil {
			return nil, err
		}
		r.Template = template

		params, err := makeRuleParams(rule.Params)
		if err != nil {
			//panic(fmt.Errorf("rule param filed %s", err))
			return nil, err
		}
		r.Params = params

		var labels map[string]string
		err = json.Unmarshal(rule.CustomLabels, &labels)
		if err != nil {
			//panic(fmt.Errorf("custom label filed %s", err))
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
				//panic(fmt.Errorf("channel failed %s", err))
				return nil, err
			}
			channels[i] = c
		}
		r.Channels = channels

		res[i] = r
	}
	return &iav1beta1.ListAlertRulesResponse{Rules: res}, nil
}

// CreateAlertRule creates Integrated Alerting rule.
func (s *RulesService) CreateAlertRule(ctx context.Context, req *iav1beta1.CreateAlertRuleRequest) (*iav1beta1.CreateAlertRuleResponse, error) {
	params := &models.CreateRuleParams{
		TemplateName: req.TemplateName,
		Disabled:     req.Disabled,
		For:          req.For,
		Severity:     req.Severity,
		CustomLabels: req.CustomLabels,
		ChannelIDs:   req.ChannelIds,
	}

	ruleParams, err := makeModelRuleParams(req.Params)
	if err != nil {
		return nil, err
	}
	params.RuleParams = ruleParams

	params.Filters = makeFilters(req.Filters)

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
	params := &models.UpdateRuleParams{
		RuleID:       req.RuleId,
		Disabled:     req.Disabled,
		For:          req.For,
		Severity:     req.Severity,
		CustomLabels: req.CustomLabels,
		ChannelIDs:   req.ChannelIds,
	}

	ruleParams, err := makeModelRuleParams(req.Params)
	if err != nil {
		return nil, err
	}
	params.RuleParams = ruleParams

	params.Filters = makeFilters(req.Filters)

	e := s.db.InTransaction(func(tx *reform.TX) error {
		_, err := models.UpdateRule(tx.Querier, req.RuleId, params)
		return err
	})
	if e != nil {
		return nil, e
	}
	return &iav1beta1.UpdateAlertRuleResponse{}, nil
}

// ToggleAlertRule allows to switch between disabled and enabled states of an Alert Rule.
func (s *RulesService) ToggleAlertRule(ctx context.Context, req *iav1beta1.ToggleAlertRuleRequest) (*iav1beta1.ToggleAlertRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToggleAlertRule not implemented")
}

// DeleteAlertRule deletes Integrated Alerting rule.
func (s *RulesService) DeleteAlertRule(ctx context.Context, req *iav1beta1.DeleteAlertRuleRequest) (*iav1beta1.DeleteAlertRuleResponse, error) {
	e := s.db.InTransaction(func(tx *reform.TX) error {
		return models.RemoveRule(tx.Querier, req.RuleId)
	})
	if e != nil {
		return nil, e
	}
	return &iav1beta1.DeleteAlertRuleResponse{}, nil
}

func makeTemplate(template *models.Template) (*iav1beta1.Template, error) {
	t := &iav1beta1.Template{
		Name:     template.Name,
		Summary:  template.Summary,
		Expr:     template.Expr,
		Params:   makeTemplateParams(template.Params),
		Severity: managementpb.Severity(managementpb.Severity_value[template.Severity]),
		For:      ptypes.DurationProto(time.Duration(template.For)),
		Source:   iav1beta1.TemplateSource(iav1beta1.TemplateSource_value[template.Source]),
	}

	createdAt, err := ptypes.TimestampProto(template.CreatedAt)
	if err != nil {
		panic(fmt.Errorf("temp timestamp failed %s", err))
		return nil, err
	}

	t.CreatedAt = createdAt

	labels, err := byteToMap(template.Labels)
	if err != nil {
		panic(fmt.Errorf("temp map unmarshall failed %s", err))
		return nil, err
	}
	t.Labels = labels

	annotations, err := byteToMap(template.Annotations)
	if err != nil {
		panic(fmt.Errorf("temp annotations unmarshall failed %s", err))
		return nil, err
	}
	t.Annotations = annotations
	return t, nil
}

func byteToMap(b []byte) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func makeTemplateParams(params models.Params) []*iav1beta1.TemplateParam {
	templateParams := make([]*iav1beta1.TemplateParam, len(params))
	for i, p := range params {
		param := &iav1beta1.TemplateParam{
			Name:    p.Name,
			Summary: p.Summary,
			Unit:    iav1beta1.ParamUnit(iav1beta1.ParamUnit_value[p.Unit]),
			Type:    iav1beta1.ParamType(iav1beta1.ParamType_value[p.Type]),
			Value: &iav1beta1.TemplateParam_Float{
				Float: &iav1beta1.TemplateFloatParam{
					Default: float32(p.FloatParam.Default),
					Min:     float32(p.FloatParam.Min),
					Max:     float32(p.FloatParam.Max),
				},
			},
		}
		templateParams[i] = param
	}
	return templateParams
}

func makeRuleParams(params models.RuleParams) ([]*iav1beta1.RuleParam, error) {
	ruleParams := make([]*iav1beta1.RuleParam, len(params))
	for i, param := range params {
		p := &iav1beta1.RuleParam{
			Name: param.Name,
			Type: iav1beta1.ParamType(param.Type),
		}

		switch p.Type {
		case iav1beta1.ParamType_BOOL:
			p.Value = &iav1beta1.RuleParam_Bool{
				Bool: param.BoolVal,
			}
		case iav1beta1.ParamType_FLOAT:
			p.Value = &iav1beta1.RuleParam_Float{
				Float: param.FloatVal,
			}
		case iav1beta1.ParamType_STRING:
			p.Value = &iav1beta1.RuleParam_String_{
				String_: param.StringVal,
			}
		default:
			return nil, errors.New("invalid rule param value type")
		}
		ruleParams[i] = p
	}
	return ruleParams, nil
}

func makeModelRuleParams(params []*iav1beta1.RuleParam) (models.RuleParams, error) {
	ruleParams := make(models.RuleParams, len(params))
	for i, param := range params {
		p := models.RuleParam{
			Name: param.Name,
			Type: models.ParamType(param.Type),
		}

		switch p.Type {
		case models.BoolRuleParam:
			p.BoolVal = param.GetBool()
		case models.FloatRuleParam:
			p.FloatVal = param.GetFloat()
		case models.StringRuleParam:
			p.StringVal = param.GetString_()
		default:
			return nil, errors.New("invalid model rule param value type")
		}
		ruleParams[i] = p
	}
	return ruleParams, nil
}

func makeFilters(filters []*iav1beta1.Filter) models.Filters {
	mFilters := make(models.Filters, len(filters))
	for i, filter := range filters {
		f := models.Filter{
			Type: models.FilterType(filter.Type),
			Key:  filter.Key,
			Val:  filter.Value,
		}
		mFilters[i] = f
	}
	return mFilters
}

// Check interfaces.
var (
	_ iav1beta1.RulesServer = (*RulesService)(nil)
)
