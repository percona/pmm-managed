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
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPMMUpdate(t *testing.T) {
	if os.Getenv("PMM_SERVER_IMAGE") == "" {
		t.Skip("can be tested only inside devcontainer")
	}

	c := newPMMUpdateCheck(logrus.WithField("test", t.Name()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.run(ctx)

	t.Run("installedPackageInfo", func(t *testing.T) {
		info := c.installedPackageInfo()
		assert.True(t, strings.HasPrefix(info.Version, "2.0.0-beta"), "%s", info.Version)
		assert.True(t, strings.HasPrefix(info.FullVersion, "1:2.0.0"), "%s", info.FullVersion)
		require.NotEmpty(t, info.BuildTime)
		assert.True(t, time.Since(*info.BuildTime) < 60*24*time.Hour, "InstalledTime = %s", info.BuildTime)
		assert.Equal(t, "local", info.Repo)
	})
}
