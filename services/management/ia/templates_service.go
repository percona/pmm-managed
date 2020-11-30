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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/percona-platform/saas/pkg/alert"
	saas "github.com/percona-platform/saas/pkg/alert"
	"github.com/percona/pmm/api/managementpb"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/percona/promconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
	"gopkg.in/yaml.v2"
)

const (
	builtinTemplatesPath = "/tmp/ia1/*.yml"
	userTemplatesPath    = "/tmp/ia2/*.yml"

	ruleFileDir = "/tmp/ia1/"

	// TODO remove once we start using real values
	dummyParamValue = "80"
)

var paramRegex = regexp.MustCompile(`\[\[.*?\]\]`)

// TemplatesService is responsible for interactions with IA rule templates.
type TemplatesService struct {
	db                   *reform.DB
	l                    *logrus.Entry
	builtinTemplatesPath string
	userTemplatesPath    string

	rw    sync.RWMutex
	rules map[string]saas.Rule
}

// NewTemplatesService creates a new TemplatesService.
func NewTemplatesService(db *reform.DB) *TemplatesService {
	return &TemplatesService{
		db:                   db,
		l:                    logrus.WithField("component", "management/ia/templates"),
		builtinTemplatesPath: builtinTemplatesPath,
		userTemplatesPath:    userTemplatesPath,
		rules:                make(map[string]saas.Rule),
	}
}

// getCollected return collected templates.
func (svc *TemplatesService) getCollected(ctx context.Context) map[string]saas.Rule {
	svc.rw.RLock()
	defer svc.rw.RUnlock()

	res := make(map[string]saas.Rule)
	for n, r := range svc.rules {
		res[n] = r
	}
	return res
}

// collect collects IA rule templates from various sources like
// built-in templates shipped with PMM and defined by the users.
func (svc *TemplatesService) collect(ctx context.Context) {
	builtinFilePaths, err := filepath.Glob(svc.builtinTemplatesPath)
	if err != nil {
		svc.l.Errorf("Failed to get paths of built-in templates files shipped with PMM: %s.", err)
		return
	}

	userFilePaths, err := filepath.Glob(svc.userTemplatesPath)
	if err != nil {
		svc.l.Errorf("Failed to get paths of user-defined template files: %s.", err)
		return
	}

	rules := make([]saas.Rule, 0, len(builtinFilePaths)+len(userFilePaths))

	for _, path := range builtinFilePaths {
		r, err := svc.loadFile(ctx, path)
		if err != nil {
			svc.l.Errorf("Failed to load shipped rule template file: %s, reason: %s.", path, err)
			return
		}

		rules = append(rules, r...)
	}

	for _, path := range userFilePaths {
		r, err := svc.loadFile(ctx, path)
		if err != nil {
			svc.l.Errorf("Failed to load user-defined rule template file: %s, reason: %s.", path, err)
			return
		}
		rules = append(rules, r...)
	}

	// TODO download templates from SAAS.

	// replace previously stored rules with newly collected ones.
	svc.rw.Lock()
	defer svc.rw.Unlock()
	svc.rules = make(map[string]saas.Rule, len(rules))
	for _, r := range rules {
		// TODO Check for name clashes? Allow users to re-define built-in rules?
		// Reserve prefix for built-in or user-defined rules?
		// https://jira.percona.com/browse/PMM-7023

		svc.rules[r.Name] = r
	}
}

// loadFile parses IA rule template file.
func (svc *TemplatesService) loadFile(ctx context.Context, file string) ([]saas.Rule, error) {
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

type ruleFile struct {
	Group []ruleGroup `yaml:"groups"`
}

type ruleGroup struct {
	Name  string `yaml:"name"`
	Rules []rule `yaml:"rules"`
}

type rule struct {
	Alert       string              `yaml:"alert"` // same as alert name in template file
	Expr        string              `yaml:"expr"`
	Duration    promconfig.Duration `yaml:"for"`
	Labels      map[string]string   `yaml:"labels,omitempty"`
	Annotations map[string]string   `yaml:"annotations,omitempty"`
}

// converts an alert template rule to a rule file. generates one file per rule.
func (svc *TemplatesService) convertTemplates(ctx context.Context) {
	templates := svc.getCollected(ctx)
	for _, template := range templates {
		r := rule{
			Alert:    template.Name,
			Duration: template.For,
			Labels:   template.Labels,
		}

		res := transformExpr(template.Expr)
		r.Expr = res.transformedExpr

		for t := range res.templateSet {
			key := strings.Trim(t, "[[ . ]]")
			r.Labels[key] = dummyParamValue
		}
		r.Labels["ia"] = "1"
		r.Labels["severity"] = template.Severity.String()

		transformAnnotations(template.Annotations, res.templateSet)
		r.Annotations = template.Annotations

		rf := ruleFile{
			Group: []ruleGroup{{
				Name:  "PMM Server Integrated Alerting",
				Rules: []rule{r},
			}},
		}
		err := dumpRule(rf)

		if err != nil {
			svc.l.Error(err)
		}
	}
}

type parsedExpr struct {
	transformedExpr string
	// stores unique templates found in expr.
	templateSet map[string]struct{}
}

// extracts unique occurences of templates in the expression
// and replaces all templates with a dummy value.
func transformExpr(expr string) parsedExpr {
	params := paramRegex.FindAll([]byte(expr), -1)
	set := make(map[string]struct{}, len(params))
	for _, p := range params {
		set[string(p)] = struct{}{}
	}

	// TODO use real values instead of dummy.
	tExpr := string(paramRegex.ReplaceAll([]byte(expr), []byte(dummyParamValue)))

	return parsedExpr{
		transformedExpr: tExpr,
		templateSet:     set,
	}
}

// replaces any occurence of a template in annotations with a dummy value.
func transformAnnotations(annotations map[string]string, templateSet map[string]struct{}) {
	var val string
	for k, v := range annotations {
		for param := range templateSet {
			if strings.Contains(v, param) {
				// TODO use real values instead of dummy.
				val = strings.ReplaceAll(v, param, dummyParamValue)
			}
		}
		annotations[k] = val
	}
}

// dump the transformed IA rules to a file.
func dumpRule(rule ruleFile) error {
	b, err := yaml.Marshal(rule)
	if err != nil {
		return errors.Errorf("failed to marshal rule %s", err)
	}
	b = append([]byte("---\n"), b...)

	filepath := ruleFileDir + rule.Group[0].Rules[0].Alert + ".yml"

	_, err = os.Stat(ruleFileDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(ruleFileDir, 0755)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(filepath, b, 0644)
	if err != nil {
		return errors.Errorf("failed to dump rule to file %s: %s", ruleFileDir, err)

	}
	return nil
}

// ListTemplates returns a list of all collected Alert Rule Templates.
func (svc *TemplatesService) ListTemplates(ctx context.Context, req *iav1beta1.ListTemplatesRequest) (*iav1beta1.ListTemplatesResponse, error) {
	if req.Reload {
		svc.collect(ctx)
	}

	templates := svc.getCollected(ctx)
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
			Source:      iav1beta1.TemplateSource_TEMPLATE_SOURCE_INVALID, // TODO
		}

		for _, p := range r.Params {
			var tp *iav1beta1.TemplateParam
			switch p.Type {
			case alert.Float:
				tp = &iav1beta1.TemplateParam{
					Name:    p.Name,
					Summary: p.Summary,
					Unit:    iav1beta1.ParamUnit_PARAM_UNIT_INVALID, // TODO
					Type:    iav1beta1.ParamType_FLOAT,
					Value:   nil, // TODO
				}
			default:
				svc.l.Warnf("Skipping unexpected parameter type %q for %q.", p.Type, r.Name)
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
func (svc *TemplatesService) CreateTemplate(ctx context.Context, req *iav1beta1.CreateTemplateRequest) (*iav1beta1.CreateTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTemplate not implemented")
}

// UpdateTemplate updates existing template, previously created via API.
func (svc *TemplatesService) UpdateTemplate(ctx context.Context, req *iav1beta1.UpdateTemplateRequest) (*iav1beta1.UpdateTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTemplate not implemented")
}

// DeleteTemplate deletes existing, previously created via API.
func (svc *TemplatesService) DeleteTemplate(ctx context.Context, req *iav1beta1.DeleteTemplateRequest) (*iav1beta1.DeleteTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTemplate not implemented")
}

// Check interfaces.
var (
	_ iav1beta1.TemplatesServer = (*TemplatesService)(nil)
)
