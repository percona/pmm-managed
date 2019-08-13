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

// Package supervisord provides facilities for working with Supervisord.
package supervisord

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/percona/pmm/utils/pdeathsig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type Service struct {
	supervisorctlPath string
	l                 *logrus.Entry
	m                 sync.Mutex
	subs              map[chan *event]sub
}

type sub struct {
	program    string
	eventTypes []eventType
}

// values from supervisord configuration
const (
	pmmUpdatePerformProgram = "pmm-update-perform"
	pmmUpdatePerformLog     = "/srv/logs/pmm-update-perform.log"
)

func New() *Service {
	return &Service{
		supervisorctlPath: "supervisorctl",
		l:                 logrus.WithField("component", "supervisord"),
		subs:              make(map[chan *event]sub),
	}
}

func (s *Service) Run(ctx context.Context) {
	var lastEvent *event
	for ctx.Err() == nil {
		cmd := exec.CommandContext(ctx, s.supervisorctlPath, "maintail", "-f") //nolint:gosec
		pdeathsig.Set(cmd, unix.SIGKILL)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Start(); err != nil {
			s.l.Error(err)
			time.Sleep(time.Second)
			continue
		}

		scanner := bufio.NewScanner(&stdout)
		for scanner.Scan() {
			e := parseEvent(scanner.Text())
			if e == nil {
				continue
			}

			// skip old events (and events with exactly the same time as old events) if maintail was restarted
			if lastEvent != nil && !lastEvent.Time.Before(e.Time) {
				continue
			}
			lastEvent = e

			s.m.Lock()
			var toDelete []chan *event
			for ch, sub := range s.subs {
				if e.Program == sub.program {
					var found bool
					for _, t := range sub.eventTypes {
						if e.Type == t {
							found = true
							break
						}
					}
					if found {
						ch <- e
						close(ch)
						toDelete = append(toDelete, ch)
					}
				}
			}
			for _, ch := range toDelete {
				delete(s.subs, ch)
			}
			s.m.Unlock()
		}
		if err := scanner.Err(); err != nil {
			s.l.Error(err)
		}

		if err := cmd.Wait(); err != nil {
			s.l.Error(err)
		}
	}
}

func (s *Service) subscribe(program string, eventTypes ...eventType) chan *event {
	ch := make(chan *event, 1)
	s.m.Lock()
	s.subs[ch] = sub{
		program:    program,
		eventTypes: eventTypes,
	}
	s.m.Unlock()
	return ch
}

func (s *Service) supervisorctl(args ...string) ([]byte, error) {
	cmd := exec.Command(s.supervisorctlPath, args...) //nolint:gosec
	pdeathsig.Set(cmd, unix.SIGKILL)
	b, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, nil
}

func (s *Service) StartPMMUpdate() error {
	ch := s.subscribe("supervisord", logReopen)

	s.m.Lock()
	defer s.m.Lock()

	// TODO comment

	if err := os.Remove(pmmUpdatePerformLog); err != nil {
		s.l.Warn(err)
	}

	b, err := s.supervisorctl("pid")
	if err != nil {
		return err
	}
	pid, err := strconv.Atoi(string(b))
	if err != nil {
		return errors.WithStack(err)
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = p.Signal(unix.SIGUSR2); err != nil {
		s.l.Warn(err)
	}
	<-ch

	fi, err := os.Stat(pmmUpdatePerformLog)
	if err != nil {
		s.l.Warn(err)
	}
	if fi.Size() != 0 {
		s.l.Warnf("%+v", fi)
	}

	_, err = s.supervisorctl("start", pmmUpdatePerformProgram)
	return err
}
