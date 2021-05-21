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
	"encoding/json"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionServiceClient(t *testing.T) {
	c := NewVersionServiceClient(versionServiceURL)

	for _, tt := range []struct {
		params componentsParams
	}{
		{params: componentsParams{operator: psmdbOperator}},
		{params: componentsParams{operator: psmdbOperator, operatorVersion: "1.6.0"}},
		{params: componentsParams{operator: psmdbOperator, operatorVersion: "1.7.0", dbVersion: "4.2.8-8"}},
		{params: componentsParams{operator: pxcOperator}},
		{params: componentsParams{operator: pxcOperator, operatorVersion: "1.7.0"}},
		{params: componentsParams{operator: pxcOperator, operatorVersion: "1.7.0", dbVersion: "8.0.20-11.2"}},
	} {
		t.Run("NotEmptyMatrix", func(t *testing.T) {
			response, err := c.Matrix(context.TODO(), tt.params)
			require.NoError(t, err)
			require.NotEmpty(t, response.Versions)
			for _, v := range response.Versions {
				switch tt.params.operator {
				case psmdbOperator:
					assert.NotEmpty(t, v.Matrix.Mongod)
				case pxcOperator:
					assert.NotEmpty(t, v.Matrix.Pxc)
					assert.NotEmpty(t, v.Matrix.Proxysql)
				}
				assert.NotEmpty(t, v.Matrix.Backup)
			}
		})
	}
}

type fakeLatestVersionServer struct {
	response *VersionServiceResponse
}

func (f fakeLatestVersionServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(f.response)
	if err != nil {
		log.Fatal(err)
	}
}

func TestLatestVersionGetting(t *testing.T) {
	t.Parallel()
	t.Run("Invalid url", func(t *testing.T) {
		t.Parallel()
		c := NewVersionServiceClient("https://check.percona.com/versions/invalid")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		v, err := c.getLatestVersion(ctx, psmdbOperator)
		assert.True(t, errors.Is(err, ErrZeroLatestVersion), "err is expected to be ErrZeroLatestVersion")
		assert.Equal(t, "", v)
	})
	t.Run("get latest", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		hostAndPort := "localhost:7453"
		fakeServer := fakeLatestVersionServer{
			response: &VersionServiceResponse{
				Versions: []struct {
					Product  string `json:"product"`
					Operator string `json:"operator"`
					Matrix   matrix `json:"matrix"`
				}{
					{Operator: "0.8.0"},
					{Operator: "1.8.0"},
					{Operator: "1.7.0"},
				},
			},
		}
		var httpServer *http.Server
		waitForListener := make(chan struct{})
		go func() {
			httpServer = &http.Server{Addr: hostAndPort, Handler: fakeServer}
			listener, err := net.Listen("tcp", hostAndPort)
			if err != nil {
				log.Fatal(err)
			}
			close(waitForListener)
			_ = httpServer.Serve(listener)
		}()
		<-waitForListener
		defer httpServer.Shutdown(ctx)

		c := NewVersionServiceClient("http://" + hostAndPort)
		v, err := c.getLatestVersion(ctx, psmdbOperator)
		require.NoError(t, err, "request to fakeserver for latest version should not fail")
		assert.Equal(t, "1.8.0", v)
	})

}
