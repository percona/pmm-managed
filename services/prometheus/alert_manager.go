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

type AlertManager struct {
	// To make this testeable. To run API tests we need to write the rules file but on dev envs
	// there is no /srv/prometheus/rules/ directory
	alertManagerFile string
	l                *logrus.Entry
}

func NewAlertManager(alertManagerFile string) *AlertManager {
	return &AlertManager{
		alertManagerFile: alertManagerFile,
		l:                logrus.WithField("component", "alert_manager"),
	}
}

func (s *AlertManager) ValidateRules(ctx context.Context, rules string) error {
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

func (s *AlertManager) ReadRules() (string, error) {
	b, err := ioutil.ReadFile(s.alertManagerFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	return string(b), err
}

func (s *AlertManager) RemoveRulesFile() error {
	return os.Remove(s.alertManagerFile)
}

func (s *AlertManager) WriteRules(rules string) error {
	return ioutil.WriteFile(s.alertManagerFile, []byte(rules), 0644)
}
