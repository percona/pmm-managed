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

// Package dbaas contains all logic related to dbaas services.
package dbaas

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/reform.v1"

	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"

	"github.com/percona/pmm-managed/models"

	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
)

// xtraDBClusterService implements XtraDBClusterServer methods.
type xtraDBClusterService struct {
	db               *reform.DB
	l                *logrus.Entry
	controllerClient controllerv1beta1.XtraDBClusterAPIClient
}

// NewXtraDBClusterService creates XtraDB Service.
func NewXtraDBClusterService(db *reform.DB, conn *grpc.ClientConn) dbaasv1beta1.XtraDBClusterServer {
	l := logrus.WithField("component", "xtradb_cluster")
	c := controllerv1beta1.NewXtraDBClusterAPIClient(conn)
	return &xtraDBClusterService{db: db, l: l, controllerClient: c}
}

// ListXtraDBClusters returns a list of all XtraDB clusters.
func (s xtraDBClusterService) ListXtraDBClusters(ctx context.Context, req *dbaasv1beta1.ListXtraDBClustersRequest) (*dbaasv1beta1.ListXtraDBClustersResponse, error) {
	kubernetesCluster, err := models.FindKubernetesClusterByName(s.db.Querier, req.KubernetesClusterName)
	if err != nil {
		return nil, err
	}

	in := controllerv1beta1.ListXtraDBClustersRequest{
		KubeAuth: &controllerv1beta1.KubeAuth{
			Kubeconfig: kubernetesCluster.KubeConfig,
		},
	}

	out, err := s.controllerClient.ListXtraDBClusters(ctx, &in)
	if err != nil {
		return nil, err
	}

	clusters := make([]*dbaasv1beta1.ListXtraDBClustersResponse_Cluster, len(out.Clusters))
	for _, c := range out.Clusters {
		cluster := dbaasv1beta1.ListXtraDBClustersResponse_Cluster{
			Name: c.Name,
			Params: &dbaasv1beta1.XtraDBClusterParams{
				ClusterSize: c.Params.ClusterSize,
				Pxc: &dbaasv1beta1.XtraDBClusterParams_PXC{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        c.Params.Pxc.ComputeResources.CpuM,
						MemoryBytes: c.Params.Pxc.ComputeResources.MemoryBytes,
					},
				},
				Proxysql: &dbaasv1beta1.XtraDBClusterParams_ProxySQL{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        c.Params.Proxysql.ComputeResources.CpuM,
						MemoryBytes: c.Params.Proxysql.ComputeResources.MemoryBytes,
					},
				},
			},
		}

		clusters = append(clusters, &cluster)
	}

	return &dbaasv1beta1.ListXtraDBClustersResponse{Clusters: clusters}, nil
}

// CreateXtraDBCluster creates XtraDB cluster with given parameters.
func (s xtraDBClusterService) CreateXtraDBCluster(ctx context.Context, req *dbaasv1beta1.CreateXtraDBClusterRequest) (*dbaasv1beta1.CreateXtraDBClusterResponse, error) {
	kubernetesCluster, err := models.FindKubernetesClusterByName(s.db.Querier, req.KubernetesClusterName)
	if err != nil {
		return nil, err
	}

	in := controllerv1beta1.CreateXtraDBClusterRequest{
		KubeAuth: &controllerv1beta1.KubeAuth{
			Kubeconfig: kubernetesCluster.KubeConfig,
		},
		Name: req.Name,
		Params: &controllerv1beta1.XtraDBClusterParams{
			ClusterSize: req.Params.ClusterSize,
			Pxc: &controllerv1beta1.XtraDBClusterParams_PXC{
				ComputeResources: &controllerv1beta1.ComputeResources{
					CpuM:        req.Params.Pxc.ComputeResources.CpuM,
					MemoryBytes: req.Params.Pxc.ComputeResources.MemoryBytes,
				},
			},
			Proxysql: &controllerv1beta1.XtraDBClusterParams_ProxySQL{
				ComputeResources: &controllerv1beta1.ComputeResources{
					CpuM:        req.Params.Proxysql.ComputeResources.CpuM,
					MemoryBytes: req.Params.Proxysql.ComputeResources.MemoryBytes,
				},
			},
		},
	}

	_, err = s.controllerClient.CreateXtraDBCluster(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.CreateXtraDBClusterResponse{}, nil
}

// UpdateXtraDBCluster updates XtraDB cluster.
func (s xtraDBClusterService) UpdateXtraDBCluster(ctx context.Context, req *dbaasv1beta1.UpdateXtraDBClusterRequest) (*dbaasv1beta1.UpdateXtraDBClusterResponse, error) {
	kubernetesCluster, err := models.FindKubernetesClusterByName(s.db.Querier, req.KubernetesClusterName)
	if err != nil {
		return nil, err
	}

	in := controllerv1beta1.UpdateXtraDBClusterRequest{
		KubeAuth: &controllerv1beta1.KubeAuth{
			Kubeconfig: kubernetesCluster.KubeConfig,
		},
		Name: req.Name,
		Params: &controllerv1beta1.XtraDBClusterParams{
			ClusterSize: req.Params.ClusterSize,
			Pxc: &controllerv1beta1.XtraDBClusterParams_PXC{
				ComputeResources: &controllerv1beta1.ComputeResources{
					CpuM:        req.Params.Pxc.ComputeResources.CpuM,
					MemoryBytes: req.Params.Pxc.ComputeResources.MemoryBytes,
				},
			},
			Proxysql: &controllerv1beta1.XtraDBClusterParams_ProxySQL{
				ComputeResources: &controllerv1beta1.ComputeResources{
					CpuM:        req.Params.Proxysql.ComputeResources.CpuM,
					MemoryBytes: req.Params.Proxysql.ComputeResources.MemoryBytes,
				},
			},
		},
	}

	_, err = s.controllerClient.UpdateXtraDBCluster(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.UpdateXtraDBClusterResponse{}, nil
}

// DeleteXtraDBCluster deletes XtraDB cluster by given name.
func (s xtraDBClusterService) DeleteXtraDBCluster(ctx context.Context, req *dbaasv1beta1.DeleteXtraDBClusterRequest) (*dbaasv1beta1.DeleteXtraDBClusterResponse, error) {
	kubernetesCluster, err := models.FindKubernetesClusterByName(s.db.Querier, req.KubernetesClusterName)
	if err != nil {
		return nil, err
	}

	in := controllerv1beta1.DeleteXtraDBClusterRequest{
		Name: req.Name,
		KubeAuth: &controllerv1beta1.KubeAuth{
			Kubeconfig: kubernetesCluster.KubeConfig,
		},
	}

	_, err = s.controllerClient.DeleteXtraDBCluster(ctx, &in)
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.DeleteXtraDBClusterResponse{}, nil
}
