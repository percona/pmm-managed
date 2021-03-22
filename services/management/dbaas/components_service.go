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

	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

type componentsService struct {
	l                    *logrus.Entry
	db                   *reform.DB
	dbaasClient          dbaasClient
	versionServiceClient *VersionServiceClient
}

// NewComponentsService creates Components Service.
func NewComponentsService(db *reform.DB, dbaasClient dbaasClient, versionService *VersionServiceClient) dbaasv1beta1.ComponentsServer {
	l := logrus.WithField("component", "components_service")
	return &componentsService{
		l:                    l,
		db:                   db,
		dbaasClient:          dbaasClient,
		versionServiceClient: versionService,
	}
}

func (c componentsService) GetPSMDBComponents(ctx context.Context, req *dbaasv1beta1.GetPSMDBComponentsRequest) (*dbaasv1beta1.GetPSMDBComponentsResponse, error) {
	params := componentsParams{
		operator:  PSMDBOperator,
		dbVersion: req.DbVersion,
	}
	if req.KubernetesClusterName != "" {
		kubernetesCluster, err := models.FindKubernetesClusterByName(c.db.Querier, req.KubernetesClusterName)
		if err != nil {
			return nil, err
		}

		checkResponse, e := c.dbaasClient.CheckKubernetesClusterConnection(ctx, kubernetesCluster.KubeConfig)
		if e != nil {
			return nil, e
		}

		params.operatorVersion = checkResponse.Operators.Psmdb.Version
	}

	versions, err := c.versions(ctx, params)
	if err != nil {
		return nil, err
	}
	return &dbaasv1beta1.GetPSMDBComponentsResponse{Versions: versions}, nil
}

func (c componentsService) GetPXCComponents(ctx context.Context, req *dbaasv1beta1.GetPXCComponentsRequest) (*dbaasv1beta1.GetPXCComponentsResponse, error) {
	params := componentsParams{
		operator:  PXCOperator,
		dbVersion: req.DbVersion,
	}
	if req.KubernetesClusterName != "" {
		kubernetesCluster, err := models.FindKubernetesClusterByName(c.db.Querier, req.KubernetesClusterName)
		if err != nil {
			return nil, err
		}

		checkResponse, e := c.dbaasClient.CheckKubernetesClusterConnection(ctx, kubernetesCluster.KubeConfig)
		if e != nil {
			return nil, e
		}

		params.operatorVersion = checkResponse.Operators.Psmdb.Version
	}

	versions, err := c.versions(ctx, params)
	if err != nil {
		return nil, err
	}
	return &dbaasv1beta1.GetPXCComponentsResponse{Versions: versions}, nil
}

func (c componentsService) versions(ctx context.Context, params componentsParams) ([]*dbaasv1beta1.Version, error) {
	var versions []*dbaasv1beta1.Version
	components, err := c.versionServiceClient.Matrix(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, v := range components.Versions {
		version := &dbaasv1beta1.Version{
			Product:  v.Product,
			Operator: v.Operator,
			Matrix: &dbaasv1beta1.Matrix{
				Mongod:       c.matrix(v.Matrix.Mongod),
				Pxc:          c.matrix(v.Matrix.Pxc),
				Pmm:          c.matrix(v.Matrix.Pmm),
				Proxysql:     c.matrix(v.Matrix.Proxysql),
				Haproxy:      c.matrix(v.Matrix.Haproxy),
				Backup:       c.matrix(v.Matrix.Backup),
				Operator:     c.matrix(v.Matrix.Operator),
				LogCollector: c.matrix(v.Matrix.LogCollector),
			},
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func (c componentsService) matrix(m map[string]component) map[string]*dbaasv1beta1.Component {
	result := make(map[string]*dbaasv1beta1.Component)

	for v, component := range m {
		result[v] = &dbaasv1beta1.Component{
			ImagePath: component.ImagePath,
			ImageHash: component.ImageHash,
			Status:    component.Status,
			Critical:  component.Critical,
			Default:   component.Status == "recommended",
		}
	}
	return result
}
