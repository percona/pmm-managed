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

// Package saasdial provides gRPC connection setup for Percona Platform.
package saasdial

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/percona/pmm/utils/tlsconfig"
)

const dialTimeout = 10 * time.Second

// Dial creates gRPC connection to Percona Platform
func Dial(ctx context.Context, hostPort string) ([]byte, error) {
	u, err := url.Parse(hostPort)
	if err != nil {
		return nil, err
	}

	tlsConfig := tlsconfig.Get()
	tlsConfig.ServerName = u.Host

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hostPort, nil)
	if err != nil {
		return nil, err
	}

	h := req.Header
	h.Add("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: dialTimeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close() //nolint:errcheck

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to dial %s, response body: %s", hostPort, bodyBytes)
	}

	return bodyBytes, nil
}
