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
	"hash"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/utils/validators"
)

const (
	editableAlertingRulesFileDir = "/srv/prometheus/rules/*.yml"
	generatedAlertingRulesDir    = "/etc/ia/rules/*.yml"
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
	// TODO: this method supposed to return external alertmanager rules.
	// TODO: For now we have common rules for both, external and internal alertmanagers.
	return "", nil
}

// RemoveRulesFile removes rules file from FS.
func (s *AlertingRules) RemoveRulesFile() error {
	// TODO: same as ReadRules()
	return nil
}

// WriteRules writes rules to file.
func (s *AlertingRules) WriteRules(rules string) error {
	// TODO: same as ReadRules()
	return nil
}

// GetRulesHash returns current rules files hash sum.
func (s *AlertingRules) GetRulesHash() ([]byte, error) {
	h := sha256.New()
	var err error

	if err = addFilesToHash(editableAlertingRulesFileDir, h); err != nil {
		return nil, err
	}

	if err = addFilesToHash(generatedAlertingRulesDir, h); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func addFilesToHash(pattern string, hash hash.Hash) error {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return errors.Wrap(err, "failed to get paths")
	}

	var b []byte
	for _, path := range paths {
		b, err = ioutil.ReadFile(path) //nolint:gosec
		if err != nil {
			return err
		}
		hash.Write(b) //nolint:errcheck
	}
	return nil
}
