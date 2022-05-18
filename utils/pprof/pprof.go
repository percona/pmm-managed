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

package pprof

import (
	"fmt"
	"io"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// Profile responds with the pprof-formatted cpu profile.
// Profiling lasts for duration specified in seconds.
func Profile(writer io.Writer, duration int64) error {
	if err := pprof.StartCPUProfile(writer); err != nil {
		return err
	}

	time.Sleep(time.Duration(duration) * time.Second)
	pprof.StopCPUProfile()

	return nil
}

// Trace responds with the execution trace in binary form.
// Tracing lasts for duration specified in seconds.
func Trace(writer io.Writer, duration int64) error {
	if err := trace.Start(writer); err != nil {
		return err
	}

	time.Sleep(time.Duration(duration) * time.Second)
	trace.Stop()

	return nil
}

// Heap responds with the pprof-formatted profile named "heap".
// listing the available profiles.
// You can specify the gc parameter to run gc before taking the heap sample.
func Heap(writer io.Writer, gc bool) error {
	debug := 0
	profile := "heap"

	p := pprof.Lookup(profile)
	if p == nil {
		return fmt.Errorf("profile cannot be found: %s", profile)
	}

	if gc {
		runtime.GC()
	}

	return p.WriteTo(writer, debug)
}
