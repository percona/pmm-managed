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
	"bytes"
	"compress/gzip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	t.Parallel()
	t.Run("Heap test", func(t *testing.T) {
		heapBytes, err := Heap(true)

		// read gzip
		reader, err := gzip.NewReader(bytes.NewBuffer(heapBytes))
		assert.NoError(t, err)

		var resB bytes.Buffer
		_, err = resB.ReadFrom(reader)
		assert.NoError(t, err)
		assert.True(t, len(resB.Bytes()) != 0)
	})
}

func TestProfile(t *testing.T) {
	t.Parallel()
	t.Run("Profile test", func(t *testing.T) {
		profileBytes, err := Profile(1 * time.Second)

		assert.NoError(t, err)
		assert.True(t, len(profileBytes) != 0)

		// read gzip
		reader, err := gzip.NewReader(bytes.NewBuffer(profileBytes))
		assert.NoError(t, err)

		var resB bytes.Buffer
		_, err = resB.ReadFrom(reader)
		assert.NoError(t, err)

		assert.True(t, len(resB.Bytes()) != 0)
	})
}

func TestTrace(t *testing.T) {
	t.Parallel()
	t.Run("Trace test", func(t *testing.T) {
		traceBytes, err := Trace(1 * time.Second)

		assert.NoError(t, err)
		assert.True(t, len(traceBytes) != 0)
	})
}
