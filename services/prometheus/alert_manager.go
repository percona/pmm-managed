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

package prometheus

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/percona/pmm/utils/pdeathsig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const alertingRulesFile = "/srv/prometheus/rules/pmm.rules.yml"

// AlertManagerRulesConfigurator contains all logic related to alert manager rules files.
type AlertManagerRulesConfigurator struct {
	l *logrus.Entry
}

// NewAlertManagerRulesConfigurator creates new AlertManagerRulesConfigurator instance.
func NewAlertManagerRulesConfigurator() *AlertManagerRulesConfigurator {
	return &AlertManagerRulesConfigurator{
		l: logrus.WithField("component", "alert_manager"),
	}
}

// ValidateRules validates alert manager rules.
func (s *AlertManagerRulesConfigurator) ValidateRules(ctx context.Context, rules string) error {
	tempFile, err := ioutil.TempFile("", "temp_rules_*.yml")
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.Remove(tempFile.Name()) //nolint:errcheck

	if _, err = tempFile.Write([]byte(rules)); err != nil {
		tempFile.Close() //nolint:errcheck
		return errors.WithStack(err)
	}
	tempFile.Close() //nolint:errcheck

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, "promtool", "check", "rules", tempFile.Name()) //nolint:gosec
	pdeathsig.Set(cmd, unix.SIGKILL)

	b, err := cmd.CombinedOutput()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok && e.ExitCode() != 0 {
			s.l.Infof("%s: %s\n%s", strings.Join(cmd.Args, " "), e, b)
			return status.Errorf(codes.InvalidArgument, "Invalid Alert Manager rules.")
		}
		return errors.WithStack(err)
	}

	if bytes.Contains(b, []byte("SUCCESS: 0 rules found")) {
		return status.Errorf(codes.InvalidArgument, "Zero Alert Manager rules found.")
	}

	s.l.Debugf("%q check passed.", strings.Join(cmd.Args, " "))
	return nil
}

// ReadRules reads current rules from FS.
func (s *AlertManagerRulesConfigurator) ReadRules() (string, error) {
	b, err := ioutil.ReadFile(alertingRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	return string(b), nil
}

// RemoveRulesFile removes rules file from FS.
func (s *AlertManagerRulesConfigurator) RemoveRulesFile() error {
	return os.Remove(alertingRulesFile)
}

// WriteRules writes rules to file.
func (s *AlertManagerRulesConfigurator) WriteRules(rules string) error {
	return ioutil.WriteFile(alertingRulesFile, []byte(rules), 0644)
}
