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
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	api "github.com/percona-platform/saas/gen/check/retrieval"
	"github.com/percona-platform/saas/pkg/alert"
	"github.com/percona-platform/saas/pkg/common"
	iav1beta1 "github.com/percona/pmm/api/managementpb/ia"
	"github.com/percona/promconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/data"
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/dir"
	"github.com/percona/pmm-managed/utils/envvars"
	"github.com/percona/pmm-managed/utils/saasdial"
	"github.com/percona/pmm-managed/utils/signatures"
)

const (
	templatesDir = "/srv/ia/templates"
	envPublicKey = "PERCONA_TEST_CHECKS_PUBLIC_KEY"
)

// templateInfo represents alerting rule template information from various sources.
//
// TODO We already have models.Template, iav1beta1.Template, and alert.Template.
//      We probably can remove that type.
type templateInfo struct {
	alert.Template
	Yaml      string
	Source    iav1beta1.TemplateSource
	CreatedAt *time.Time
}

// TemplatesService is responsible for interactions with IA rule templates.
type TemplatesService struct {
	db                *reform.DB
	l                 *logrus.Entry
	userTemplatesPath string

	host       string
	publicKeys []string

	rw        sync.RWMutex
	templates map[string]templateInfo
}

// NewTemplatesService creates a new TemplatesService.
func NewTemplatesService(db *reform.DB) (*TemplatesService, error) {
	l := logrus.WithField("component", "management/ia/templates")

	err := dir.CreateDataDir(templatesDir, "pmm", "pmm", dirPerm)
	if err != nil {
		l.Error(err)
	}

	host, err := envvars.GetSAASHost()
	if err != nil {
		return nil, err
	}

	s := &TemplatesService{
		db:                db,
		l:                 l,
		userTemplatesPath: templatesDir,
		host:              host,
		templates:         make(map[string]templateInfo),
	}

	// TODO: same code exists in check service, move it to a better place.
	if k := os.Getenv(envPublicKey); k != "" {
		s.publicKeys = strings.Split(k, ",")
		l.Warnf("Public keys changed to %q.", k)
	}

	return s, nil
}

// Enabled returns if service is enabled and can be used.
func (s *TemplatesService) Enabled() bool {
	settings, err := models.GetSettings(s.db)
	if err != nil {
		s.l.WithError(err).Error("can't get settings")
		return false
	}
	return settings.IntegratedAlerting.Enabled
}

func newParamTemplate() *template.Template {
	return template.New("").Option("missingkey=error").Delims("[[", "]]")
}

// getTemplates return collected templates.
func (s *TemplatesService) getTemplates() map[string]templateInfo {
	s.rw.RLock()
	defer s.rw.RUnlock()

	res := make(map[string]templateInfo)
	for n, r := range s.templates {
		res[n] = r
	}
	return res
}

// Collect collects IA rule templates from various sources like:
// builtin templates: read from the generated code in bindata.go.
// SaaS templates: templates downloaded from checks service.
// user file templates: read from yaml files created by the user in `/srv/ia/templates`
// user API templates: in the DB created using the API.
func (s *TemplatesService) Collect(ctx context.Context) {
	builtInTemplates, err := s.loadTemplatesFromAssets(ctx)
	if err != nil {
		s.l.Errorf("Failed to load built-in rule templates: %s.", err)
		return
	}

	userDefinedTemplates, err := s.loadTemplatesFromUserFiles(ctx)
	if err != nil {
		s.l.Errorf("Failed to load user-defined rule templates: %s.", err)
		return
	}

	dbTemplates, err := s.loadTemplatesFromDB()
	if err != nil {
		s.l.Errorf("Failed to load rule templates from DB: %s.", err)
		return
	}

	saasTemplates, err := s.downloadTemplates(ctx)
	if err != nil {
		s.l.Errorf("Failed to download rule templates from SaaS: %s.", err)
		return
	}

	templates := make([]templateInfo, 0, len(builtInTemplates)+len(userDefinedTemplates)+len(dbTemplates)+len(saasTemplates))

	for _, t := range builtInTemplates {
		templates = append(templates, templateInfo{
			Template: t,
			Source:   iav1beta1.TemplateSource_BUILT_IN,
		})
	}

	for _, t := range userDefinedTemplates {
		templates = append(templates, templateInfo{
			Template: t,
			Source:   iav1beta1.TemplateSource_USER_FILE,
		})
	}

	for _, t := range saasTemplates {
		templates = append(templates, templateInfo{
			Template: t,
			Source:   iav1beta1.TemplateSource_SAAS,
		})
	}

	templates = append(templates, dbTemplates...)

	// replace previously stored templates with newly collected ones.
	s.rw.Lock()
	defer s.rw.Unlock()
	s.templates = make(map[string]templateInfo, len(templates))
	for _, t := range templates {
		// TODO Check for name clashes? Allow users to re-define built-in templates?
		// Reserve prefix for built-in or user-defined templates?
		// https://jira.percona.com/browse/PMM-7023

		s.templates[t.Name] = t
	}
}

// loadTemplatesFromAssets loads built-in alerting rule templates from pmm-managed binary's assets.
func (s *TemplatesService) loadTemplatesFromAssets(ctx context.Context) ([]alert.Template, error) {
	paths := data.AssetNames()
	res := make([]alert.Template, 0, len(paths))
	for _, path := range paths {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		data, err := data.Asset(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read rule template asset: %s", path)
		}

		// be strict about built-in templates
		params := &alert.ParseParams{
			DisallowUnknownFields:    true,
			DisallowInvalidTemplates: true,
		}
		templates, err := alert.Parse(bytes.NewReader(data), params)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse rule template asset: %s", path)
		}

		// built-in-specific validations
		// TODO move to some better / common place

		if l := len(templates); l != 1 {
			return nil, errors.Errorf("%q should contain exactly one template, got %d", path, l)
		}

		t := templates[0]

		filename := filepath.Base(path)
		if strings.HasPrefix(filename, "pmm_") {
			return nil, errors.Errorf("%q file name should not start with 'pmm_' prefix", path)
		}
		if !strings.HasPrefix(t.Name, "pmm_") {
			return nil, errors.Errorf("%s %q: template name should start with 'pmm_' prefix", path, t.Name)
		}
		if expected := strings.TrimPrefix(t.Name, "pmm_") + ".yml"; filename != expected {
			return nil, errors.Errorf("template file name %q should be %q", filename, expected)
		}
		if len(t.Annotations) != 2 || t.Annotations["summary"] == "" || t.Annotations["description"] == "" {
			return nil, errors.Errorf("%s %q: template should contain exactly two annotations: summary and description", path, t.Name)
		}

		res = append(res, t)
	}
	return res, nil
}

// loadTemplatesFromUserFiles loads user's alerting rule templates from /srv/ia/templates.
func (s *TemplatesService) loadTemplatesFromUserFiles(ctx context.Context) ([]alert.Template, error) {
	paths, err := dir.FindFilesWithExtensions(s.userTemplatesPath, "yml", "yaml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get paths")
	}

	res := make([]alert.Template, 0, len(paths))
	for _, path := range paths {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		data, err := ioutil.ReadFile(path) //nolint:gosec
		if err != nil {
			s.l.Warnf("Failed to load rule template file %s.", path)
			continue
		}

		// be strict about user files
		params := &alert.ParseParams{
			DisallowUnknownFields:    true,
			DisallowInvalidTemplates: true,
		}
		templates, err := alert.Parse(bytes.NewReader(data), params)
		if err != nil {
			s.l.Warnf("Failed to parse rule template file %s.", path)
			continue
		}

		for _, t := range templates {
			if err = validateUserTemplate(&t); err != nil { //nolint:gosec
				s.l.Warnf("%s %s", path, err)
				continue
			}

			res = append(res, t)
		}
	}
	return res, nil
}

func (s *TemplatesService) loadTemplatesFromDB() ([]templateInfo, error) {
	var templates []models.Template
	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		templates, err = models.FindTemplates(tx.Querier)
		return err
	})
	if e != nil {
		return nil, errors.Wrap(e, "failed to load rule templates form DB")
	}

	res := make([]templateInfo, 0, len(templates))
	for _, t := range templates {
		t := t
		params := make([]alert.Parameter, 0, len(t.Params))
		for _, param := range t.Params {
			p := alert.Parameter{
				Name:    param.Name,
				Summary: param.Summary,
				Unit:    alert.Unit(param.Unit),
				Type:    alert.Type(param.Type),
			}

			switch alert.Type(param.Type) {
			case alert.Float:
				f := param.FloatParam

				if f.Default != nil {
					p.Value = *f.Default
				}

				if f.Min != nil && f.Max != nil {
					p.Range = []interface{}{*f.Min, *f.Max}
				}
			}

			params = append(params, p)
		}

		labels, err := t.GetLabels()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load template labels")
		}

		annotations, err := t.GetAnnotations()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load template annotations")
		}

		res = append(res,
			templateInfo{
				Template: alert.Template{
					Name:        t.Name,
					Version:     t.Version,
					Summary:     t.Summary,
					Tiers:       t.Tiers,
					Expr:        t.Expr,
					Params:      params,
					For:         promconfig.Duration(t.For),
					Severity:    common.Severity(t.Severity),
					Labels:      labels,
					Annotations: annotations,
				},
				Yaml:      t.Yaml,
				Source:    convertSource(t.Source),
				CreatedAt: &t.CreatedAt,
			},
		)
	}

	return res, nil
}

// downloadTemplates downloads IA templates from SaaS.
func (s *TemplatesService) downloadTemplates(ctx context.Context) ([]alert.Template, error) {
	s.l.Infof("Downloading templates from %s ...", s.host)

	settings, err := models.GetSettings(s.db)
	if err != nil {
		return nil, err
	}

	cc, err := saasdial.Dial(ctx, settings.SaaS.SessionID, s.host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}
	defer cc.Close() //nolint:errcheck

	resp, err := api.NewRetrievalAPIClient(cc).GetAllAlertRuleTemplates(ctx, &api.GetAllAlertRuleTemplatesRequest{})
	if err != nil {
		// if credentials are invalid then force a logout so that the next check download
		// attempt can be successful.
		logoutErr := saasdial.LogoutIfInvalidAuth(s.db, s.l, err)
		if logoutErr != nil {
			s.l.Warnf("Failed to force logout: %v", logoutErr)
		}
		return nil, errors.Wrap(err, "failed to request checks service")
	}

	if err = signatures.Verify(s.l, resp.File, resp.Signatures, s.publicKeys); err != nil {
		return nil, err
	}

	// be liberal about files from SaaS for smooth transition to future versions
	params := &alert.ParseParams{
		DisallowUnknownFields:    false,
		DisallowInvalidTemplates: false,
	}
	templates, err := alert.Parse(strings.NewReader(resp.File), params)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// validateUserTemplate validates user-provided template (API or file).
func validateUserTemplate(t *alert.Template) error {
	// TODO move to some better place

	if strings.HasPrefix(t.Name, "pmm_") || strings.HasPrefix(t.Name, "saas_") {
		return errors.Errorf("%s: template name should not start with 'pmm_' or 'saas_' prefix", t.Name)
	}

	// TODO more validations

	// validate expression template with fake parameter values
	params := make(map[string]string, len(t.Params))
	for _, p := range t.Params {
		var value string
		switch p.Type {
		case alert.Float:
			value = "0"
		case alert.Bool:
			value = "false"
		case alert.String:
			value = "param_text"
		default:
			return errors.Errorf("invalid parameter type %s", p.Type)
		}

		params[p.Name] = value
	}

	if _, err := templateRuleExpr(t.Expr, params); err != nil {
		return err
	}

	return nil
}

func convertSource(source models.Source) iav1beta1.TemplateSource {
	switch source {
	case models.BuiltInSource:
		return iav1beta1.TemplateSource_BUILT_IN
	case models.SAASSource:
		return iav1beta1.TemplateSource_SAAS
	case models.UserFileSource:
		return iav1beta1.TemplateSource_USER_FILE
	case models.UserAPISource:
		return iav1beta1.TemplateSource_USER_API
	default:
		return iav1beta1.TemplateSource_TEMPLATE_SOURCE_INVALID
	}
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

// ListTemplates returns a list of all collected Alert Rule Templates.
func (s *TemplatesService) ListTemplates(ctx context.Context, req *iav1beta1.ListTemplatesRequest) (*iav1beta1.ListTemplatesResponse, error) {
	var pageIndex int
	var pageSize int
	if req.PageParams != nil {
		pageIndex = int(req.PageParams.Index)
		pageSize = int(req.PageParams.PageSize)
	}

	if req.Reload {
		s.Collect(ctx)
	}

	templates := s.getTemplates()
	res := &iav1beta1.ListTemplatesResponse{
		Templates: make([]*iav1beta1.Template, 0, len(templates)),
		Totals: &iav1beta1.PageTotals{
			TotalItems: int32(len(templates)),
			TotalPages: 1,
		},
	}

	if pageSize > 0 {
		res.Totals.TotalPages = int32(len(templates) / pageSize)
		if len(templates)%pageSize > 0 {
			res.Totals.TotalPages++
		}
	}

	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	from, to := pageIndex*pageSize, (pageIndex+1)*pageSize
	if to > len(names) || to == 0 {
		to = len(names)
	}

	if from > len(names) {
		from = len(names)
	}

	for _, name := range names[from:to] {
		t, err := convertTemplate(s.l, templates[name])
		if err != nil {
			return nil, err
		}

		res.Templates = append(res.Templates, t)
	}

	return res, nil
}

// CreateTemplate creates a new template.
func (s *TemplatesService) CreateTemplate(ctx context.Context, req *iav1beta1.CreateTemplateRequest) (*iav1beta1.CreateTemplateResponse, error) {
	pParams := &alert.ParseParams{
		DisallowUnknownFields:    true,
		DisallowInvalidTemplates: true,
	}

	templates, err := alert.Parse(strings.NewReader(req.Yaml), pParams)
	if err != nil {
		s.l.Errorf("failed to parse rule template form request: +%v", err)
		return nil, status.Error(codes.InvalidArgument, "Failed to parse rule template.")
	}

	if len(templates) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Request should contain exactly one rule template.")
	}

	for _, t := range templates {
		if err = validateUserTemplate(&t); err != nil { //nolint:gosec
			return nil, status.Errorf(codes.InvalidArgument, "%s.", err)
		}
	}

	params := &models.CreateTemplateParams{
		Template: &templates[0],
		Yaml:     req.Yaml,
		Source:   models.UserAPISource,
	}

	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		_, err = models.CreateTemplate(tx.Querier, params)
		return err
	})
	if e != nil {
		return nil, e
	}

	s.Collect(ctx)

	return &iav1beta1.CreateTemplateResponse{}, nil
}

// UpdateTemplate updates existing template, previously created via API.
func (s *TemplatesService) UpdateTemplate(ctx context.Context, req *iav1beta1.UpdateTemplateRequest) (*iav1beta1.UpdateTemplateResponse, error) {
	parseParams := &alert.ParseParams{
		DisallowUnknownFields:    true,
		DisallowInvalidTemplates: true,
	}

	templates, err := alert.Parse(strings.NewReader(req.Yaml), parseParams)
	if err != nil {
		s.l.Errorf("failed to parse rule template form request: +%v", err)
		return nil, status.Error(codes.InvalidArgument, "Failed to parse rule template.")
	}

	if len(templates) != 1 {
		return nil, status.Error(codes.InvalidArgument, "Request should contain exactly one rule template.")
	}

	tmpl := templates[0]

	if err = validateUserTemplate(&tmpl); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s.", err)
	}

	changeParams := &models.ChangeTemplateParams{
		Template: &tmpl,
		Name:     req.Name,
		Yaml:     req.Yaml,
	}

	e := s.db.InTransaction(func(tx *reform.TX) error {
		var err error
		_, err = models.ChangeTemplate(tx.Querier, changeParams)
		return err
	})
	if e != nil {
		return nil, e
	}

	s.Collect(ctx)

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

	s.Collect(ctx)

	return &iav1beta1.DeleteTemplateResponse{}, nil
}

// Check interfaces.
var (
	_ iav1beta1.TemplatesServer = (*TemplatesService)(nil)
)
