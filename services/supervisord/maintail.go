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

package supervisord

import (
	"regexp"
	"strings"
	"time"
)

type eventType string

// See http://supervisord.org/subprocess.html#process-states
const (
	exitedExpected   eventType = "EXITED (expected)"
	exitedUnexpected eventType = "EXITED (unexpected)"
	fatal            eventType = "FATAL"
)

var (
	exitedExpectedRE   = regexp.MustCompile(`^exited\: ([\w-]+) \(exit status \d+; expected\)$`)
	exitedUnexpectedRE = regexp.MustCompile(`^exited\: ([\w-]+) \(exit status \d+; not expected\)$`)
	fatalRE            = regexp.MustCompile(`^gave up\: ([\w-]+) entered FATAL state, too many start retries too quickly$`)

	events = map[*regexp.Regexp]eventType{
		exitedExpectedRE:   exitedExpected,
		exitedUnexpectedRE: exitedUnexpected,
		fatalRE:            fatal,
	}
)

// event represents supervisord program event.
type event struct {
	Time    time.Time
	Type    eventType
	Program string
}

// parseEvent returns parsed event from supervisord maintail line, or nil.
func parseEvent(line string) *event {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) != 4 {
		return nil
	}

	// see https://github.com/golang/go/issues/6189
	ts := strings.Replace(parts[0]+" "+parts[1], ",", ".", 1)
	t, err := time.Parse("2006-01-02 15:04:05.000", ts)
	if err != nil {
		return nil
	}

	for re, typ := range events {
		if m := re.FindStringSubmatch(parts[3]); m != nil {
			return &event{
				Time:    t,
				Type:    typ,
				Program: m[1],
			}
		}
	}

	return nil
}
