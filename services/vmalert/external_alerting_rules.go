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
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/utils/validators"
)

const externalAlertingRulesFile = "/srv/prometheus/rules/pmm.rules.yml"

// ExternalAlertingRules wraps external alerting rules handling.
type ExternalAlertingRules struct {
	l *logrus.Entry
}

// NewExternalAlertingRules creates new ExternalAlertingRules instance.
func NewExternalAlertingRules() *ExternalAlertingRules {
	return &ExternalAlertingRules{
		l: logrus.WithField("component", "external_alerting_rules"),
	}
}

// ValidateRules validates alerting rules.
func (s *ExternalAlertingRules) ValidateRules(ctx context.Context, rules string) error {
	err := validators.ValidateAlertingRules(ctx, rules)
	if e, ok := err.(*validators.InvalidAlertingRuleError); ok {
		return status.Errorf(codes.InvalidArgument, e.Msg)
	}
	return err
}

// ReadRules reads current rules from FS.
func (s *ExternalAlertingRules) ReadRules() (string, error) {
	b, err := ioutil.ReadFile(externalAlertingRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return "", errors.Wrap(err, "failed to read external alerting rules")
	}
	return string(b), nil
}

// RemoveRulesFile removes rules file from FS.
func (s *ExternalAlertingRules) RemoveRulesFile() error {
	err := os.Remove(externalAlertingRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to remove external alerting rules")
	}
	return nil
}

// WriteRules writes rules to file.
func (s *ExternalAlertingRules) WriteRules(rules string) error {
	err := ioutil.WriteFile(externalAlertingRulesFile, []byte(rules), 0o644) //nolint:gosec
	if err != nil {
		return errors.Wrap(err, "failed to write external alerting rules")
	}
	return nil
}
