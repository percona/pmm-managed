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
	"bytes"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/percona-platform/saas/pkg/alert"
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/percona/pmm-managed/models"
)

const (
	dirPerm = os.FileMode(0o775)
)

func convertParamUnit(u alert.Unit) iav1beta1.ParamUnit {
	switch u {
	case alert.Percentage:
		return iav1beta1.ParamUnit_PERCENTAGE
	case alert.Seconds:
		return iav1beta1.ParamUnit_SECONDS
	}

	// do not add `default:` to make exhaustive linter do its job

	return iav1beta1.ParamUnit_PARAM_UNIT_INVALID
}

func convertRule(l *logrus.Entry, rule *models.Rule, channels []*models.Channel) (*iav1beta1.Rule, error) {
	r := &iav1beta1.Rule{
		RuleId:          rule.ID,
		Disabled:        rule.Disabled,
		Summary:         rule.TemplateSummary,
		Name:            rule.Summary,
		ExprTemplate:    rule.Expr,
		DefaultSeverity: managementpb.Severity(rule.DefaultSeverity),
		Severity:        managementpb.Severity(rule.Severity),
		DefaultFor:      durationpb.New(rule.DefaultFor),
		For:             durationpb.New(rule.For),
		CreatedAt:       timestamppb.New(rule.CreatedAt),
	}

	if err := r.CreatedAt.CheckValid(); err != nil {
		return nil, errors.Wrap(err, "failed to convert timestamp")
	}

	params, err := processExprParameters(rule.ParamsDefinitions, rule.ParamsValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process expression parameters")
	}

	r.Expr, err = fillExprWithParams(rule.Expr, params.AsStringMap())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fill expression template with parameters values")
	}

	r.ParamsDefinitions, err = convertModelToParamsDefinitions(rule.ParamsDefinitions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert parameters definitions")
	}

	r.ParamsValues, err = convertModelToParamValues(rule.ParamsValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert parameters values")
	}

	r.CustomLabels, err = rule.GetCustomLabels()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load rule custom labels")
	}

	r.Labels, err = rule.GetLabels()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load rule labels")
	}

	r.Annotations, err = rule.GetAnnotations()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load rule annotations")
	}

	r.Filters = make([]*iav1beta1.Filter, len(rule.Filters))
	for i, filter := range rule.Filters {
		r.Filters[i] = &iav1beta1.Filter{
			Type:  convertModelToFilterType(filter.Type),
			Key:   filter.Key,
			Value: filter.Val,
		}
	}

	cm := make(map[string]*models.Channel)
	for _, channel := range channels {
		cm[channel.ID] = channel
	}

	for _, id := range rule.ChannelIDs {
		channel, ok := cm[id]
		if !ok {
			l.Warningf("Skip missing channel with ID %s", id)
			continue
		}

		c, err := convertChannel(channel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert channel")
		}
		r.Channels = append(r.Channels, c)
	}

	return r, nil
}

func fillExprWithParams(expr string, values map[string]string) (string, error) {
	var buf bytes.Buffer
	t, err := newParamTemplate().Parse(expr)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse expression")
	}
	if err = t.Execute(&buf, values); err != nil {
		return "", errors.Wrap(err, "failed to fill expression placeholders")
	}
	return buf.String(), nil
}

func processExprParameters(definitions models.ParamsDefinitions, values models.ParamsValues) (models.ParamsValues, error) {
	unknownParams := make(map[string]struct{}, len(values))
	for _, p := range values {
		unknownParams[p.Name] = struct{}{}
	}

	params := make(models.ParamsValues, 0, len(definitions))
	for _, d := range definitions {
		var filled bool
		for _, rp := range values {
			if rp.Name == d.Name {
				if string(d.Type) != string(rp.Type) {
					return nil, status.Errorf(codes.InvalidArgument, "Parameter %s has type %s instead of %s.", d.Name, rp.Type, d.Type)
				}
				delete(unknownParams, rp.Name)
				filled = true
				params = append(params, rp)
				break
			}
		}

		if !filled {
			p := models.ParamValue{
				Name: d.Name,
				Type: d.Type,
			}

			switch d.Type {
			case models.Float:
				v := d.FloatParam
				if v.Default == nil {
					return nil, status.Errorf(codes.InvalidArgument, "Parameter %s doesn't have "+
						"default value, so it should be specified explicitly", d.Name)
				}
				p.FloatValue = float32(pointer.GetFloat64(v.Default))
			}
			params = append(params, p)
		}
	}

	names := make([]string, 0, len(unknownParams))
	for name := range unknownParams {
		names = append(names, name)
	}
	if len(names) != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Unknown parameters %s.", names)
	}

	return params, nil
}
