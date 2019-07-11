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

package grafana

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// logrus.SetLevel(logrus.TraceLevel)

	ctx := context.Background()
	c := NewClient("127.0.0.1:3000")

	req, err := http.NewRequest("GET", "/dummy", nil)
	require.NoError(t, err)
	req.SetBasicAuth("admin", "admin")
	authHeaders := req.Header

	t.Run("isGrafanaAdmin", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			admin, err := c.isGrafanaAdmin(ctx, authHeaders)
			assert.NoError(t, err)
			assert.True(t, admin)
		})

		t.Run("Unauthorized", func(t *testing.T) {
			admin, err := c.isGrafanaAdmin(ctx, nil)
			apiError, _ := err.(*apiError)
			require.NotNil(t, apiError)
			assert.Equal(t, 401, apiError.code)
			assert.False(t, admin)
		})
	})

	t.Run("getRole", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			role, err := c.getRole(ctx, authHeaders)
			assert.NoError(t, err)
			assert.Equal(t, admin, role)
			assert.Equal(t, "admin", role.String())
		})

		t.Run("Unauthorized", func(t *testing.T) {
			role, err := c.getRole(ctx, nil)
			apiError, _ := err.(*apiError)
			require.NotNil(t, apiError)
			assert.Equal(t, 401, apiError.code)
			assert.Equal(t, none, role)
			assert.Equal(t, "none", role.String())
		})
	})

	t.Run("CreateAnnotation", func(t *testing.T) {
		t.Skip("https://jira.percona.com/browse/PMM-3812")

		from := time.Now()

		t.Run("Normal", func(t *testing.T) {
			msg, err := c.CreateAnnotation(ctx, []string{"tag1", "tag2"}, "Normal")
			require.NoError(t, err)
			assert.Equal(t, "Annotation added", msg)

			annotations, err := c.findAnnotations(ctx, from, from.Add(time.Second))
			require.NoError(t, err)
			for _, a := range annotations {
				if a.Text == "Normal" {
					assert.Equal(t, []string{"pmm_annotation", "tag1", "tag2"}, a.Tags)
					assert.InDelta(t, from.Unix(), a.Time.Unix(), 1)
					return
				}
			}
			assert.Fail(t, "annotation not found", "%s", annotations)
		})

		t.Run("Empty", func(t *testing.T) {
			msg, err := c.CreateAnnotation(ctx, nil, "")
			require.NoError(t, err)
			assert.Equal(t, "Failed to save annotation", msg)
		})

		t.Run("No tags", func(t *testing.T) {
			msg, err := c.CreateAnnotation(ctx, nil, "No tags")
			require.NoError(t, err)
			assert.Equal(t, "Annotation added", msg)

			annotations, err := c.findAnnotations(ctx, from, from.Add(time.Second))
			require.NoError(t, err)
			for _, a := range annotations {
				if a.Text == "No tags" {
					assert.Equal(t, []string{"pmm_annotation"}, a.Tags)
					assert.InDelta(t, from.Unix(), a.Time.Unix(), 1)
					return
				}
			}
			assert.Fail(t, "annotation not found", "%s", annotations)
		})
	})
}
