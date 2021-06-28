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
	"strings"
	"testing"

	goversion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionServiceClient(t *testing.T) {
	c := NewVersionServiceClient("https://check.percona.com/versions/v1")

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
	response   *VersionServiceResponse
	components []string
}

func (f fakeLatestVersionServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	var response *VersionServiceResponse
	var certainVersionRequested bool
	var component string
	for _, c := range f.components {
		if strings.Contains(r.URL.Path, c) {
			component = c
			certainVersionRequested = strings.Contains(r.URL.Path, component+"/")
			break
		}
	}
	if certainVersionRequested {
		segments := strings.Split(r.URL.Path, "/")
		version := segments[len(segments)-2]
		var dbVersion string
		// handle product/version/applyversion
		if _, err := goversion.NewVersion(version); err == nil {
			dbVersion = segments[len(segments)-1]
		} else {
			version = segments[len(segments)-1]
		}
		for _, v := range f.response.Versions {
			if v.Operator == version && v.Product == component {
				if dbVersion != "" {
					var database map[string]componentVersion
					switch component {
					case pxcOperator:
						database = v.Matrix.Pxc
					case psmdbOperator:
						database = v.Matrix.Mongod
					default:
						panic(component + " not supported")
					}
					if _, ok := database[dbVersion]; !ok {
						response = nil
						break
					}
				}
				response = &VersionServiceResponse{
					Versions: []struct {
						Product  string `json:"product"`
						Operator string `json:"operator"`
						Matrix   matrix `json:"matrix"`
					}{v},
				}
				break
			}
		}
	} else if component != "" {
		response = &VersionServiceResponse{}
		for _, v := range f.response.Versions {
			if v.Product == component {
				response.Versions = append(response.Versions, v)
			}
		}
	} else {
		panic("path " + r.URL.Path + " not expected")
	}
	err := encoder.Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

func newFakeVersionService(response *VersionServiceResponse, port string, components ...string) (versionService, func(*testing.T)) {
	var httpServer *http.Server
	waitForListener := make(chan struct{})
	server := fakeLatestVersionServer{
		response:   response,
		components: components,
	}
	fakeHostAndPort := "localhost:" + port
	go func() {
		httpServer = &http.Server{Addr: fakeHostAndPort, Handler: server}
		listener, err := net.Listen("tcp", fakeHostAndPort)
		if err != nil {
			log.Fatal(err)
		}
		close(waitForListener)
		_ = httpServer.Serve(listener)
	}()
	<-waitForListener

	return NewVersionServiceClient("http://" + fakeHostAndPort + "/versions/v1"), func(t *testing.T) {
		assert.NoError(t, httpServer.Shutdown(context.TODO()))
	}
}
