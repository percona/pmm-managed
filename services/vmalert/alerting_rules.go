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

package vmalert

import (
	"context"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/utils/validators"
)

const (
	// TODO: Currently that file can be edited via Settings API.
	// TODO: It seems that Settings API should edit configuration for external AlertManager.
	alertingRulesFile = "/srv/prometheus/rules/pmm.rules.yml"

	generatedAlertingRulesDir = "/etc/ia/rules/*.yml"
)

// AlertingRules contains all logic related to alerting rules files.
type AlertingRules struct {
	l *logrus.Entry
}

// NewAlertingRules creates new AlertingRules instance.
func NewAlertingRules() *AlertingRules {
	return &AlertingRules{
		l: logrus.WithField("component", "alerting_rules"),
	}
}

// ValidateRules validates alerting rules.
func (s *AlertingRules) ValidateRules(ctx context.Context, rules string) error {
	err := validators.ValidateAlertingRules(ctx, rules)
	if e, ok := err.(*validators.InvalidAlertingRuleError); ok {
		return status.Errorf(codes.InvalidArgument, e.Msg)
	}
	return err
}

// ReadRules reads current rules from FS.
func (s *AlertingRules) ReadRules() (string, error) {
	b, err := ioutil.ReadFile(alertingRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	return string(b), nil
}

// RemoveRulesFile removes rules file from FS.
func (s *AlertingRules) RemoveRulesFile() error {
	return os.Remove(alertingRulesFile)
}

// WriteRules writes rules to file.
func (s *AlertingRules) WriteRules(rules string) error {
	return ioutil.WriteFile(alertingRulesFile, []byte(rules), 0o644) //nolint:gosec
}

// GetRulesHash returns current rules files hash sum.
func (s *AlertingRules) GetRulesHash() ([]byte, error) {
	h := sha256.New()
	b, err := ioutil.ReadFile(alertingRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	h.Write(b) //nolint:errcheck

	paths, err := filepath.Glob(generatedAlertingRulesDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get paths")
	}

	for _, path := range paths {
		b, err = ioutil.ReadFile(path) //nolint:gosec
		if err != nil {
			return nil, err
		}
		h.Write(b) //nolint:errcheck
	}

	return h.Sum(nil), nil
}
