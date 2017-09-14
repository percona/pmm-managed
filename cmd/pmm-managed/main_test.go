// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package consul provides facilities for working with Consul.
package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestRunTelemetry(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	logrus.SetOutput(w)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	ch := make(chan bool)
	go func() {
		runTelemetryService(ctx, nil)
		ch <- true
	}()

	select {
	case <-ch:
		break
	case <-time.After(1100 * time.Millisecond):
		t.Fatalf("Service didn't stop upon context cancellation")
	}

	w.Flush()
	assert.Contains(t, b.String(), `msg="Telemetry is enabled. Send data interval = 24h0m0s" component=TELEMETRY`)

	logrus.SetOutput(os.Stdout)
}
