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
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
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
	PXCOperator   map[string]componentVersion `json:"pxcOperator"`
	PSMDBOperator map[string]componentVersion `json:"psmdbOperator"`
	LogCollector  map[string]componentVersion `json:"logCollector"`
}

// VersionServiceResponse represents response from version service API.
type VersionServiceResponse struct {
	Versions []struct {
		Product  string `json:"product"`
		Operator string `json:"operator"`
		Matrix   matrix `json:"matrix"`
	} `json:"versions"`
}

// componentsParams contains params to filter components in version service API.
type componentsParams struct {
	operator        string
	operatorVersion string
	dbVersion       string
}

// VersionServiceClient represents a client for Version Service API.
type VersionServiceClient struct {
	url  string
	http *http.Client
	irtm prom.Collector
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
	paths := []string{c.url, params.operator}
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

func (c *VersionServiceClient) IsDatabaseVersionSupportedByOperator(ctx context.Context, operatorType, operatorVersion, databaseVersion string) (bool, error) {
	m, err := c.Matrix(ctx, componentsParams{
		operator:        operatorType,
		operatorVersion: operatorVersion,
		dbVersion:       databaseVersion,
	})
	if err != nil {
		return false, err
	}
	return len(m.Versions) != 0, nil
}

func (c *VersionServiceClient) IsOperatorVersionSupported(ctx context.Context, operatorType string, pmmVersion string, operatorVersion string) (bool, error) {
	pmm, err := goversion.NewVersion(pmmVersion)
	if err != nil {
		return false, err
	}
	resp, err := c.Matrix(ctx, componentsParams{operator: "pmm-server", operatorVersion: pmm.Core().String()})
	if err != nil {
		return false, err
	}
	if len(resp.Versions) == 0 {
		return false, nil
	}
	var operator map[string]componentVersion
	switch operatorType {
	case pxcOperator:
		operator = resp.Versions[0].Matrix.PXCOperator
	case psmdbOperator:
		operator = resp.Versions[0].Matrix.PSMDBOperator
	default:
		return false, errors.Errorf("%q is an unknown operator type", operatorType)
	}

	for version, _ := range operator {
		if version == operatorVersion {
			return true, nil
		}
	}
	return false, nil
}

// TODO THIS IS JUST COPY PASTED CODE SO I DON'T HAVE TO REBASE EVERYTIME THE ORIGIN OF THIS CODE CHANGES -> IT'S IN REVEIW
// --------------------
var errNoVersionsFound error = errors.New("no version found")

func getLatest(m map[string]componentVersion) (*goversion.Version, error) {
	if len(m) == 0 {
		return nil, errNoVersionsFound
	}
	latest := goversion.Must(goversion.NewVersion("v0.0.0"))
	for version := range m {
		parsedVersion, err := goversion.NewVersion(version)
		if err != nil {
			return nil, err
		}
		if parsedVersion.GreaterThan(latest) {
			latest = parsedVersion
		}
	}
	return latest, nil
}

// GetLatestOperatorVersion return latest PXC and PSMDB operators for given PMM version.
func (c *VersionServiceClient) GetLatestOperatorVersion(ctx context.Context, pmmVersion string) (*goversion.Version, *goversion.Version, error) {
	if pmmVersion == "" {
		return nil, nil, errors.New("given PMM version is empty")
	}
	params := componentsParams{
		operator:        "pmm-server",
		operatorVersion: pmmVersion,
	}
	resp, err := c.Matrix(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	if len(resp.Versions) != 1 {
		return nil, nil, nil // no deps for the PMM version passed to c.Matrix
	}
	pmmVersionDeps := resp.Versions[0]
	latestPSMDBOperator, err := getLatest(pmmVersionDeps.Matrix.PSMDBOperator)
	if err != nil {
		return nil, nil, err
	}
	latestPXCOperator, err := getLatest(pmmVersionDeps.Matrix.PXCOperator)
	return latestPXCOperator, latestPSMDBOperator, err
}

// -------------------- TODO REMOVE THIS SECTION
