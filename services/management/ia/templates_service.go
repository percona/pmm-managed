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
	"path/filepath"

	saas "github.com/percona-platform/saas/pkg/alert"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	shippedRuleTemplatePath     = "/tmp/ia1/*.yml"
	userDefinedRuleTemplatePath = "/tmp/ia2/*.yml"
)

// TemplatesService is responsible for interactions with IA rule templates.
type TemplatesService struct {
	l                           *logrus.Entry
	shippedRuleTemplatePath     string
	userDefinedRuleTemplatePath string
	rules                       map[string]saas.Rule
}

// NewTemplatesService creates a new TemplatesService.
func NewTemplatesService() *TemplatesService {
	return &TemplatesService{
		l:                           logrus.WithField("component", "templates service"),
		shippedRuleTemplatePath:     shippedRuleTemplatePath,
		userDefinedRuleTemplatePath: userDefinedRuleTemplatePath,
		rules:                       make(map[string]saas.Rule),
	}
}

// Run starts the TemplatesService.
func (svc *TemplatesService) Run() {
	svc.l.Info("Starting...")
	defer svc.l.Info("Done.")

	svc.collectRuleTemplates()
}

// collectRuleTemplates collects IA rule templates from various sources like
// templates shipped with PMM and defined by the users.
func (svc *TemplatesService) collectRuleTemplates() {
	shippedFilePaths, err := filepath.Glob(svc.shippedRuleTemplatePath)
	if err != nil {
		svc.l.Errorf("Failed to get paths of template files shipped with PMM: %s.", err)
		return
	}

	userDefinedFilePaths, err := filepath.Glob(svc.userDefinedRuleTemplatePath)
	if err != nil {
		svc.l.Errorf("Failed to get paths of user-defined template files: %s.", err)
		return
	}

	rules := make([]saas.Rule, 0, len(shippedFilePaths)+len(userDefinedFilePaths))

	for _, path := range shippedFilePaths {
		r, err := svc.loadRuleTemplates(path)
		if err != nil {
			svc.l.Errorf("Failed to load shipped rule template file: %s, reason: %s.", path, err)
			return
		}
		rules = append(rules, r...)
	}

	for _, path := range userDefinedFilePaths {
		r, err := svc.loadRuleTemplates(path)
		if err != nil {
			svc.l.Errorf("Failed to load user-defined rule template file: %s, reason: %s.", path, err)
			return
		}
		rules = append(rules, r...)
	}

	// TODO download templates from SAAS.

	// replace previously stored rules with newly collected ones.
	for k := range svc.rules {
		delete(svc.rules, k)
	}

	for _, r := range rules {
		svc.rules[r.Name] = r
	}
}

// loadRuleTemplates parses IA rule template files.
func (svc *TemplatesService) loadRuleTemplates(file string) ([]saas.Rule, error) {
	data, err := ioutil.ReadFile(file) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, "failed to read test rule template file")
	}

	// be strict about local files
	params := &saas.ParseParams{
		DisallowUnknownFields: true,
		DisallowInvalidRules:  true,
	}
	rules, err := saas.Parse(bytes.NewReader(data), params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse test rule template file")
	}

	return rules, nil
}

// ListTemplates returns a list of all collected Alert Rule Templates.
func (svc *TemplatesService) ListTemplates(context.Context, *iav1beta1.ListTemplatesRequest) (*iav1beta1.ListTemplatesResponse, error) {
	panic("Unimplemented")
}

// CreateTemplate creates a new template.
func (svc *TemplatesService) CreateTemplate(context.Context, *iav1beta1.CreateTemplateRequest) (*iav1beta1.CreateTemplateResponse, error) {
	panic("Unimplemented")
}

// UpdateTemplate updates existing template, previously created via API.
func (svc *TemplatesService) UpdateTemplate(context.Context, *iav1beta1.UpdateTemplateRequest) (*iav1beta1.UpdateTemplateResponse, error) {
	panic("Unimplemented")
}

// DeleteTemplate deletes existing, previously created via API.
func (svc *TemplatesService) DeleteTemplate(context.Context, *iav1beta1.DeleteTemplateRequest) (*iav1beta1.DeleteTemplateResponse, error) {
	panic("Unimplemented")
}

// Check interfaces.
var (
	_ iav1beta1.TemplatesServer = (*TemplatesService)(nil)
)
