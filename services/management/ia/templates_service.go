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
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona-platform/saas/pkg/alert"
	saas "github.com/percona-platform/saas/pkg/alert"
	"github.com/percona-platform/saas/pkg/common"
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/percona/promconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

const (
	builtinTemplatesPath = "/tmp/ia1/*.yml"
	userTemplatesPath    = "/tmp/ia2/*.yml"
)

// Rule represents alerting rule/rule template with added source field.
type Rule struct {
	alert.Rule
	Yaml   string
	Source iav1beta1.TemplateSource
}

// TemplatesService is responsible for interactions with IA rule templates.
type TemplatesService struct {
	db                   *reform.DB
	l                    *logrus.Entry
	builtinTemplatesPath string
	userTemplatesPath    string

	rw    sync.RWMutex
	rules map[string]Rule
}

// NewTemplatesService creates a new TemplatesService.
func NewTemplatesService(db *reform.DB) *TemplatesService {
	return &TemplatesService{
		db:                   db,
		l:                    logrus.WithField("component", "management/ia/templates"),
		builtinTemplatesPath: builtinTemplatesPath,
		userTemplatesPath:    userTemplatesPath,
		rules:                make(map[string]Rule),
	}
}

// getCollected return collected templates.
func (s *TemplatesService) getCollected(ctx context.Context) map[string]Rule {
	s.rw.RLock()
	defer s.rw.RUnlock()

	res := make(map[string]Rule)
	for n, r := range s.rules {
		res[n] = r
	}
	return res
}

// collect collects IA rule templates from various sources like
// built-in templates shipped with PMM and defined by the users.
func (s *TemplatesService) collect(ctx context.Context) {
	rules := make([]Rule, 0, len(s.builtinTemplatesPath)+len(s.userTemplatesPath))

	builtInRules, err := s.loadRulesFromFiles(ctx, s.builtinTemplatesPath)
	if err != nil {
		s.l.Errorf("Failed to load built-in rule templates: %s.", err)
		return
	}
	for _, rule := range builtInRules {
		rules = append(rules, Rule{
			Rule:   rule,
			Source: iav1beta1.TemplateSource_BUILT_IN,
		})
	}

	userDefinedRules, err := s.loadRulesFromFiles(ctx, s.userTemplatesPath)
	if err != nil {
		s.l.Errorf("Failed to load user-defined rule templates: %s.", err)
		return
	}
	for _, rule := range userDefinedRules {
		rules = append(rules, Rule{
			Rule:   rule,
			Source: iav1beta1.TemplateSource_USER_FILE,
		})
	}

	dbRules, err := s.loadRulesFromDB()
	if err != nil {
		s.l.Errorf("Failed to load rule templates from DB: %s.", err)
		return
	}
	rules = append(rules, dbRules...)

	// TODO download templates from SAAS.

	// replace previously stored rules with newly collected ones.
	s.rw.Lock()
	defer s.rw.Unlock()
	s.rules = make(map[string]Rule, len(rules))
	for _, r := range rules {
		// TODO Check for name clashes? Allow users to re-define built-in rules?
		// Reserve prefix for built-in or user-defined rules?
		// https://jira.percona.com/browse/PMM-7023

		s.rules[r.Name] = r
	}
}

func (s *TemplatesService) loadRulesFromFiles(ctx context.Context, path string) ([]alert.Rule, error) {
	paths, err := filepath.Glob(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get paths")
	}

	res := make([]alert.Rule, 0, len(paths))
	for _, path := range paths {
		r, err := s.loadFile(ctx, path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load rule template file: %s", path)
		}

		res = append(res, r...)
	}
	return res, nil
}

func (s *TemplatesService) loadRulesFromDB() ([]Rule, error) {
	var templates []models.Template
	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		templates, err = models.FindTemplates(tx.Querier)
		return err
	})
	if e != nil {
		return nil, errors.Wrap(e, "failed to load rule templates form DB")
	}

	res := make([]Rule, 0, len(templates))
	for _, template := range templates {
		params := make([]alert.Parameter, len(template.Params))
		for _, param := range template.Params {
			p := alert.Parameter{
				Name:    param.Name,
				Summary: param.Summary,
				Unit:    param.Unit,
				Type:    alert.Type(param.Type),
			}

			switch alert.Type(param.Type) {
			case alert.Float:
				f := param.FloatParam
				p.Value = f.Default
				p.Range = []interface{}{f.Min, f.Max}

			}

			params = append(params, p)
		}

		labels, err := template.GetLabels()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load template labels")
		}

		annotations, err := template.GetAnnotations()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load template annotations")
		}

		source := iav1beta1.TemplateSource_TEMPLATE_SOURCE_INVALID
		if v, ok := iav1beta1.TemplateSource_value[template.Source]; ok {
			source = iav1beta1.TemplateSource(v)
		}

		res = append(res,
			Rule{
				Rule: alert.Rule{
					Name:        template.Name,
					Version:     template.Version,
					Summary:     template.Summary,
					Tiers:       template.Tiers,
					Expr:        template.Expr,
					Params:      params,
					For:         promconfig.Duration(template.For),
					Severity:    common.ParseSeverity(template.Severity),
					Labels:      labels,
					Annotations: annotations,
				},
				Yaml:   template.Yaml,
				Source: source,
			},
		)
	}

	return res, nil
}

// loadFile parses IA rule template file.
func (s *TemplatesService) loadFile(ctx context.Context, file string) ([]saas.Rule, error) {
	if ctx.Err() != nil {
		return nil, errors.WithStack(ctx.Err())
	}

	data, err := ioutil.ReadFile(file) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "failed to read rule template file")
	}

	// be strict about local files
	params := &saas.ParseParams{
		DisallowUnknownFields: true,
		DisallowInvalidRules:  true,
	}
	rules, err := saas.Parse(bytes.NewReader(data), params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse rule template file")
	}

	return rules, nil
}

func convertParamType(t alert.Type) iav1beta1.ParamType {
	// TODO: add another types.
	switch t {
	case alert.Float:
		return iav1beta1.ParamType_FLOAT
	default:
		return iav1beta1.ParamType_PARAM_TYPE_INVALID
	}
}

func convertParamUnit(u string) iav1beta1.ParamUnit {
	// TODO: check possible variants.
	switch u {
	case "%", "percentage":
		return iav1beta1.ParamUnit_PERCENTAGE
	default:
		return iav1beta1.ParamUnit_PARAM_UNIT_INVALID
	}
}

// ListTemplates returns a list of all collected Alert Rule Templates.
func (s *TemplatesService) ListTemplates(ctx context.Context, req *iav1beta1.ListTemplatesRequest) (*iav1beta1.ListTemplatesResponse, error) {
	if req.Reload {
		s.collect(ctx)
	}

	templates := s.getCollected(ctx)
	res := &iav1beta1.ListTemplatesResponse{
		Templates: make([]*iav1beta1.Template, 0, len(templates)),
	}
	for _, r := range templates {
		t := &iav1beta1.Template{
			Name:        r.Name,
			Summary:     r.Summary,
			Expr:        r.Expr,
			Params:      make([]*iav1beta1.TemplateParam, 0, len(r.Params)),
			For:         ptypes.DurationProto(time.Duration(r.For)),
			Severity:    managementpb.Severity(r.Severity),
			Labels:      r.Labels,
			Annotations: r.Annotations,
			Source:      r.Source,
			Yaml:        r.Yaml,
		}

		for _, p := range r.Params {
			var tp *iav1beta1.TemplateParam
			switch p.Type {
			case alert.Float:
				value, err := p.GetValueForFloat()
				if err != nil {
					return nil, errors.Wrap(err, "failed to get value for float parameter")
				}

				min, max, err := p.GetRangeForFloat()
				if err != nil {
					return nil, errors.Wrap(err, "failed to get range for float parameter")
				}

				tp = &iav1beta1.TemplateParam{
					Name:    p.Name,
					Summary: p.Summary,
					Unit:    convertParamUnit(p.Unit),
					Type:    convertParamType(p.Type),
					Value: &iav1beta1.TemplateParam_Float{
						Float: &iav1beta1.TemplateFloatParam{
							HasDefault: true,           // TODO remove or fill with valid value.
							Default:    float32(value), // TODO eliminate conversion.
							HasMin:     true,           // TODO remove or fill with valid value.
							Min:        float32(min),   // TODO eliminate conversion.,
							HasMax:     true,           // TODO remove or fill with valid value.
							Max:        float32(max),   // TODO eliminate conversion.,
						},
					},
				}
			default:
				s.l.Warnf("Skipping unexpected parameter type %q for %q.", p.Type, r.Name)
			}

			if tp != nil {
				t.Params = append(t.Params, tp)
			}
		}

		res.Templates = append(res.Templates, t)
	}

	sort.Slice(res.Templates, func(i, j int) bool { return res.Templates[i].Name < res.Templates[j].Name })
	return res, nil
}

// CreateTemplate creates a new template.
func (s *TemplatesService) CreateTemplate(ctx context.Context, req *iav1beta1.CreateTemplateRequest) (*iav1beta1.CreateTemplateResponse, error) {
	pParams := &alert.ParseParams{
		DisallowUnknownFields: true,
		DisallowInvalidRules:  true,
	}

	fmt.Println(req.Yaml)

	rules, err := alert.Parse(strings.NewReader(req.Yaml), pParams)
	if err != nil {
		s.l.Errorf("failed to parse rule template form request: +%v", err)
		return nil, status.Error(codes.InvalidArgument, "Failed to parse rule template.")
	}

	if len(rules) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Request should contain exactly one rule template.")
	}

	params := &models.CreateTemplateParams{
		Rule:   &rules[0],
		Yaml:   req.Yaml,
		Source: iav1beta1.TemplateSource_USER_API.String(),
	}

	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		_, err = models.CreateTemplate(tx.Querier, params)
		return err
	})
	if e != nil {
		return nil, e
	}

	return &iav1beta1.CreateTemplateResponse{}, nil
}

// UpdateTemplate updates existing template, previously created via API.
func (s *TemplatesService) UpdateTemplate(ctx context.Context, req *iav1beta1.UpdateTemplateRequest) (*iav1beta1.UpdateTemplateResponse, error) {
	pParams := &alert.ParseParams{
		DisallowUnknownFields: true,
		DisallowInvalidRules:  true,
	}

	rules, err := alert.Parse(strings.NewReader(req.Yaml), pParams)
	if err != nil {
		s.l.Errorf("failed to parse rule template form request: +%v", err)
		return nil, status.Error(codes.InvalidArgument, "Failed to parse rule template.")
	}

	if len(rules) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Request should contain exactly one rule template.")
	}

	params := &models.ChangeTemplateParams{
		Rule: &rules[0],
	}

	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		_, err = models.ChangeTemplate(tx.Querier, params)
		return err
	})
	if e != nil {
		return nil, e
	}

	return &iav1beta1.UpdateTemplateResponse{}, nil
}

// DeleteTemplate deletes existing, previously created via API.
func (s *TemplatesService) DeleteTemplate(ctx context.Context, req *iav1beta1.DeleteTemplateRequest) (*iav1beta1.DeleteTemplateResponse, error) {
	e := s.db.InTransaction(func(tx *reform.TX) error {
		return models.RemoveTemplate(tx.Querier, req.Name)
	})
	if e != nil {
		return nil, e
	}
	return &iav1beta1.DeleteTemplateResponse{}, nil
}

// Check interfaces.
var (
	_ iav1beta1.TemplatesServer = (*TemplatesService)(nil)
)
