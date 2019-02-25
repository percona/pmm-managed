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
	"io/ioutil"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrometheusConfig(t *testing.T) {
	ctx, p, before := SetupTest(t)
	defer TearDownTest(t, p, before)

	// check that we can write it exactly as it was
	c, err := p.loadConfig()
	assert.NoError(t, err)
	assert.NoError(t, p.saveConfigAndReload(ctx, c))
	after, err := ioutil.ReadFile(p.configPath)
	require.NoError(t, err)
	beforeS, afterS := string(before), string(after)
	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(beforeS),
		FromFile: "Before",
		B:        difflib.SplitLines(afterS),
		ToFile:   "After",
		Context:  1,
	})
	require.NoError(t, err)
	require.Equal(t, strings.Split(beforeS, "\n"), strings.Split(afterS, "\n"), "%s", diff)
	require.Len(t, c.ScrapeConfigs, 0)

	// TODO
	/*
		// specifically check that we can read secrets
		require.NotNil(t, c.ScrapeConfigs[2].HTTPClientConfig.BasicAuth)
		assert.Equal(t, "pmm", c.ScrapeConfigs[2].HTTPClientConfig.BasicAuth.Password)

		// check that invalid configuration is reverted
		c.ScrapeConfigs[1].ScrapeInterval = model.Duration(time.Second)
		err = p.saveConfigAndReload(ctx, c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `scrape timeout greater than scrape interval`)
		after, err = ioutil.ReadFile(p.configPath)
		require.NoError(t, err)
		assert.Equal(t, before, after)
	*/
}

func TestPrometheusScrapeConfigs(t *testing.T) {
	ctx, p, before := SetupTest(t)
	defer TearDownTest(t, p, before)

	_ = ctx
}
