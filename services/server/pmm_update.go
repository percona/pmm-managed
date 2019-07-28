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
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"sync"

	"github.com/percona/pmm/utils/pdeathsig"
	"github.com/percona/pmm/version"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type pmmUpdate struct {
	m sync.Mutex
}

func (p *pmmUpdate) checkUpdates(ctx context.Context) (*version.UpdateCheckResult, error) {
	p.m.Lock()
	defer p.m.Unlock()

	cmd := exec.CommandContext(ctx, "pmm2-update", "-check") //nolint:gosec
	cmd.Stderr = os.Stderr
	pdeathsig.Set(cmd, unix.SIGKILL)
	b, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var res version.UpdateCheckResult
	if err = json.Unmarshal(b, &res); err != nil {
		return nil, errors.WithStack(err)
	}
	return &res, nil
}

func (p *pmmUpdate) startUpdate() error {
	p.m.Lock()
	defer p.m.Unlock()

	cmd := exec.Command("supervisorctl", "start", "pmm-update") //nolint:gosec
	b, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to start pmm-update")
	}
	if bytes.Contains(b, []byte(`ERROR`)) {
		return errors.Errorf("failed to start pmm-update: %s", b)
	}
	return nil
}
