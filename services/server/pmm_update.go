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

package server

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/percona/pmm/utils/pdeathsig"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const checkResultFresh = updateCheckInterval + 10*time.Minute

// pmmUpdate wraps pmm2-update invocations with caching.
type pmmUpdate struct {
	l            *logrus.Entry
	rw           sync.RWMutex
	latestResult *version.UpdateCheckResult
	latestTime   time.Time
}

func newPMMUpdate(l *logrus.Entry) *pmmUpdate {
	return &pmmUpdate{
		l: l,
	}
}

// checkResult returns the latest `pmm-update -check` result.
// It may force re-check if the latest result is empty or too old.
func (p *pmmUpdate) checkResult() *version.UpdateCheckResult {
	p.rw.RLock()
	defer p.rw.RUnlock()

	if time.Since(p.latestTime) > checkResultFresh {
		p.rw.RUnlock()
		_ = p.check()
		p.rw.RLock()
	}

	return p.latestResult
}

// check calls `pmm2-update -check` and fills latestResult/latestTime on success.
func (p *pmmUpdate) check() error {
	p.rw.Lock()
	defer p.rw.Unlock()

	// TODO use `supervisorctl start` and `supervisorctl tail` instead https://jira.percona.com/browse/PMM-4455
	cmdLine := "pmm2-update -check"
	args := strings.Split(cmdLine, " ")
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	pdeathsig.Set(cmd, unix.SIGKILL)

	b, err := cmd.Output()
	if err != nil {
		p.l.Errorf("%s output: %s", cmdLine, stderr.Bytes())
		return errors.WithStack(err)
	}

	var res version.UpdateCheckResult
	if err = json.Unmarshal(b, &res); err != nil {
		p.l.Errorf("%s output: %s", cmdLine, stderr.Bytes())
		return errors.WithStack(err)
	}

	p.l.Debugf("%s output: %s", cmdLine, stderr.Bytes())
	p.latestResult = &res
	p.latestTime = time.Now()
	return nil
}
