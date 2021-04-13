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

package dbaas

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/percona/pmm/version"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

func TestClient(t *testing.T) {
	getClient := func(t *testing.T) *Client {
		err, c := NewClient("127.0.0.1:20201")
		require.NoError(t, err, "Cannot connect to dbaas-controller")
		return c
	}
	t.Run("ValidKubeConfig", func(t *testing.T) {
		v := os.Getenv("ENABLE_DBAAS")
		dbaasEnabled, err = strconv.ParseBool(v)
		if err != nil {
			t.Skipf("Invalid value %q for environment variable ENABLE_DBAAS", v)
		}
		if !dbaasEnabled {
			t.Skip("DBaaS is not enabled")
		}
		kubeConfig := os.Getenv("DBAAS_KUBECONFIG")
		if kubeConfig == "" {
			t.Skip("DBAAS_KUBECONFIG env variable is not provided")
		}
		c := getClient(t)
		_, err := c.CheckKubernetesClusterConnection(context.TODO(), kubeConfig)
		require.NoError(t, err)
	})

	t.Run("InvalidKubeConfig", func(t *testing.T) {
		v := os.Getenv("ENABLE_DBAAS")
		dbaasEnabled, err = strconv.ParseBool(v)
		if err != nil {
			t.Skipf("Invalid value %q for environment variable ENABLE_DBAAS", v)
		}
		if !dbaasEnabled {
			t.Skip("DBaaS is not enabled")
		}

		c := getClient(t)
		_, err := c.CheckKubernetesClusterConnection(context.TODO(), "{}")
		require.Error(t, err)
	})
}
