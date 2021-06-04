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
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/percona/pmm-managed/utils/irt"
)

const (
	psmdbOperator = "psmdb-operator"
	pxcOperator   = "pxc-operator"
)

// componentVersion contains info about exact component version.
type componentVersion struct {
	ImagePath string `json:"imagePath"`
	ImageHash string `json:"imageHash"`
	Status    string `json:"status"`
	Critical  bool   `json:"critical"`
}

type matrix struct {
	Mongod        map[string]componentVersion `json:"mongod"`
	Pxc           map[string]componentVersion `json:"pxc"`
	Pmm           map[string]componentVersion `json:"pmm"`
	Proxysql      map[string]componentVersion `json:"proxysql"`
	Haproxy       map[string]componentVersion `json:"haproxy"`
	Backup        map[string]componentVersion `json:"backup"`
	Operator      map[string]componentVersion `json:"operator"`
	LogCollector  map[string]componentVersion `json:"logCollector"`
	PXCOperator   map[string]componentVersion `json:"pxcOperator,omitempty"`
	PSMDBOperator map[string]componentVersion `json:"psmdbOperator,omitempty"`
}

// VersionServiceResponse represents response from version service API.
type VersionServiceResponse struct {
	Versions []struct {
		Product        string `json:"product"`
		ProductVersion string `json:"operator"`
		Matrix         matrix `json:"matrix"`
	} `json:"versions"`
}

// componentsParams contains params to filter components in version service API.
type componentsParams struct {
	product        string
	productVersion string
	versionToApply string
}

// VersionServiceClient represents a client for Version Service API.
type VersionServiceClient struct {
	url  string
	http *http.Client
	irtm prom.Collector
	l    *logrus.Entry
}

// NewVersionServiceClient creates a new client for given version service URL.
func NewVersionServiceClient(url string) *VersionServiceClient {
	var t http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          50,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if logrus.GetLevel() >= logrus.TraceLevel {
		t = irt.WithLogger(t, logrus.WithField("component", "versionService/client").Tracef)
	}
	t, irtm := irt.WithMetrics(t, "versionService_client")

	return &VersionServiceClient{
		url: url,
		http: &http.Client{
			Transport: t,
		},
		irtm: irtm,
		l:    logrus.WithField("component", "VersionServiceClient"),
	}
}

// Describe implements prometheus.Collector.
func (c *VersionServiceClient) Describe(ch chan<- *prom.Desc) {
	c.irtm.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *VersionServiceClient) Collect(ch chan<- prom.Metric) {
	c.irtm.Collect(ch)
}

// Matrix calls version service with given params and returns components matrix.
func (c *VersionServiceClient) Matrix(ctx context.Context, params componentsParams) (*VersionServiceResponse, error) {
	paths := []string{c.url, params.product}
	if params.operatorVersion != "" {
		paths = append(paths, params.operatorVersion)
		if params.dbVersion != "" {
			paths = append(paths, params.dbVersion)
		}
	}
	url := strings.Join(paths, "/")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vsResponse VersionServiceResponse
	err = json.Unmarshal(body, &vsResponse)
	if err != nil {
		return nil, err
	}

	return &vsResponse, nil
}

// GetLatestOperatorVersion return latest oprator for given pmm version, if it's empty, it returns latest of all versions.
func (v *VersionServiceClient) GetLatestOperatorVersion(ctx context.Context, operatorType, pmmVersion string) (*goversion.Version, error) {
	params := componentsParams{
		product: "pmm-server",
	}
	if pmmVersion != "" {
		params.productVersion = pmmVersion
	}
	resp, err := v.Matrix(ctx, params)
	if err != nil {
		return nil, err
	}
	latestPerPMMVersion := make([]string, 0, len(resp.Versions))
	for _, pmmVersionDeps := range resp.Versions {
		var operatorVersions componentVersion
		switch operatorType {
		case pxcOperator:
			operatorVersions = pmmVersionDeps.Matrix.PXCOperator
		case psmdbOperator:
			operatorVersions = pmmVersionDeps.Matrix.PSMDBOperator
		default:
			return nil, errors.New("unsupported operator type")
		}
		latest, err := getLatest(keysToSliceOfVersions(operatorVersions))
		if err != nil {
			return "", err
		}
		latestPerPMMVersion = append(latestPerPMMVersion, latest)
	}
	return getLatest(latestPerPMMVersion)
}

func keysToSliceOfVersions(m map[string]interface{}) ([]*goversion.Version, error) {
	keys := make([]*goversion.Version, len(m))
	i := 0
	var err error
	for k := range m {
		keys[i], err = goversion.NewVersion(k)
		if err != nil {
			return []*goversion.Version{}, err
		}
		i++
	}
	return keys, nil
}

func (v *VersionServiceClient) getLatest(versions []*goversion.Version) (*goversion.Version, error) {
	if len(versions) == 0 {
		return "", errors.New("no latest version found")
	}
	latest := goversion.Must("v0.0.0")
	for _, version := range versions {
		if version.GreaterThan(latest) {
			latest = version
		}
	}
	return latest
}
