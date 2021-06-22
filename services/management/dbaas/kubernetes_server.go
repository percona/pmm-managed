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
	"regexp"
	"sync"

	dbaascontrollerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"
	pmmversion "github.com/percona/pmm/version"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

var (
	operatorIsForbiddenRegexp  = regexp.MustCompile(`.*\.percona\.com is forbidden`)
	resourceDoesntExistsRegexp = regexp.MustCompile(`the server doesn't have a resource type "(PerconaXtraDBCluster|PerconaServerMongoDB)"`)
)

type kubernetesServer struct {
	l              *logrus.Entry
	db             *reform.DB
	dbaasClient    dbaasClient
	versionService versionService
}

// NewKubernetesServer creates Kubernetes Server.
func NewKubernetesServer(db *reform.DB, dbaasClient dbaasClient, vs versionService) dbaasv1beta1.KubernetesServer {
	l := logrus.WithField("component", "kubernetes_server")
	return &kubernetesServer{l: l, db: db, dbaasClient: dbaasClient, versionService: vs}
}

// Enabled returns if service is enabled and can be used.
func (k *kubernetesServer) Enabled() bool {
	settings, err := models.GetSettings(k.db)
	if err != nil {
		k.l.WithError(err).Error("can't get settings")
		return false
	}
	return settings.DBaaS.Enabled
}

// convertOperatorStatus exists mainly to assign appropriate status when installed operator is unsupported.
// dbaas-controller does not have a clue what's supported, so we have to do it here.
func (k kubernetesServer) convertOperatorStatus(ctx context.Context, inputOperator *dbaascontrollerv1beta1.Operator, outputOperator *dbaasv1beta1.Operator) error {
	switch inputOperator.Status {
	case dbaascontrollerv1beta1.OperatorsStatus_OPERATORS_STATUS_NOT_INSTALLED:
		outputOperator.Status = dbaasv1beta1.OperatorsStatus_OPERATORS_STATUS_NOT_INSTALLED
	case dbaascontrollerv1beta1.OperatorsStatus_OPERATORS_STATUS_OK:
		supported, err := k.versionService.IsOperatorVersionSupported(ctx, pmmversion.PMMVersion, pxcOperator, inputOperator.Version)
		if err != nil {
			return err
		}
		if supported {
			outputOperator.Status = dbaasv1beta1.OperatorsStatus_OPERATORS_STATUS_OK
		} else {
			outputOperator.Status = dbaasv1beta1.OperatorsStatus_OPERATORS_STATUS_UNSUPPORTED
		}
	}
	return nil
}

// ListKubernetesClusters returns a list of all registered Kubernetes clusters.
func (k kubernetesServer) ListKubernetesClusters(ctx context.Context, _ *dbaasv1beta1.ListKubernetesClustersRequest) (*dbaasv1beta1.ListKubernetesClustersResponse, error) {
	kubernetesClusters, err := models.FindAllKubernetesClusters(k.db.Querier)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	clusters := make([]*dbaasv1beta1.ListKubernetesClustersResponse_Cluster, len(kubernetesClusters))
	for i, cluster := range kubernetesClusters {
		i := i
		cluster := cluster
		wg.Add(1)
		go func() {
			defer wg.Done()
			clusters[i] = &dbaasv1beta1.ListKubernetesClustersResponse_Cluster{
				KubernetesClusterName: cluster.KubernetesClusterName,
				Operators: &dbaasv1beta1.Operators{
					Xtradb: new(dbaasv1beta1.Operator),
					Psmdb:  new(dbaasv1beta1.Operator),
				},
			}
			resp, e := k.dbaasClient.CheckKubernetesClusterConnection(ctx, cluster.KubeConfig)
			if e != nil {
				clusters[i].Status = dbaasv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_UNAVAILABLE
				return
			}

			clusters[i].Status = dbaasv1beta1.KubernetesClusterStatus(resp.Status)

			if resp.Operators == nil {
				return
			}
			if err := k.convertOperatorStatus(ctx, resp.Operators.Xtradb, clusters[i].Operators.Xtradb); err != nil {
				k.l.Errorf("failed to convert dbaas-controller operator status to PMM status: %v", err)
			}
			if err := k.convertOperatorStatus(ctx, resp.Operators.Psmdb, clusters[i].Operators.Psmdb); err != nil {
				k.l.Errorf("failed to convert dbaas-controller operator status to PMM status: %v", err)
			}
		}()
	}
	wg.Wait()
	return &dbaasv1beta1.ListKubernetesClustersResponse{KubernetesClusters: clusters}, nil
}

// RegisterKubernetesCluster registers an existing Kubernetes cluster in PMM.
func (k kubernetesServer) RegisterKubernetesCluster(ctx context.Context, req *dbaasv1beta1.RegisterKubernetesClusterRequest) (*dbaasv1beta1.RegisterKubernetesClusterResponse, error) {
	err := k.db.InTransaction(func(t *reform.TX) error {
		_, e := k.dbaasClient.CheckKubernetesClusterConnection(ctx, req.KubeAuth.Kubeconfig)
		if e != nil {
			return e
		}

		_, err := models.CreateKubernetesCluster(k.db.Querier, &models.CreateKubernetesClusterParams{
			KubernetesClusterName: req.KubernetesClusterName,
			KubeConfig:            req.KubeAuth.Kubeconfig,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	pxcOperatorVersion, psmdbOperatorVersion, err := k.versionService.GetLatestOperatorVersion(ctx, pmmversion.PMMVersion)
	if err != nil {
		return nil, err
	}

	_, err = k.dbaasClient.InstallXtraDBOperator(ctx, &dbaascontrollerv1beta1.InstallXtraDBOperatorRequest{
		KubeAuth: &dbaascontrollerv1beta1.KubeAuth{
			Kubeconfig: req.KubeAuth.Kubeconfig,
		},
		Version: pxcOperatorVersion.String(),
	})
	if err != nil {
		return nil, err
	}
	_, err = k.dbaasClient.InstallPSMDBOperator(ctx, &dbaascontrollerv1beta1.InstallPSMDBOperatorRequest{
		KubeAuth: &dbaascontrollerv1beta1.KubeAuth{
			Kubeconfig: req.KubeAuth.Kubeconfig,
		},
		Version: psmdbOperatorVersion.String(),
	})
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.RegisterKubernetesClusterResponse{}, nil
}

// UnregisterKubernetesCluster removes a registered Kubernetes cluster from PMM.
func (k kubernetesServer) UnregisterKubernetesCluster(ctx context.Context, req *dbaasv1beta1.UnregisterKubernetesClusterRequest) (*dbaasv1beta1.UnregisterKubernetesClusterResponse, error) {
	err := k.db.InTransaction(func(t *reform.TX) error {
		kubernetesCluster, err := models.FindKubernetesClusterByName(k.db.Querier, req.KubernetesClusterName)
		if err != nil {
			return err
		}

		if req.Force {
			return models.RemoveKubernetesCluster(k.db.Querier, req.KubernetesClusterName)
		}

		xtraDBClusters, err := k.dbaasClient.ListXtraDBClusters(ctx,
			&dbaascontrollerv1beta1.ListXtraDBClustersRequest{
				KubeAuth: &dbaascontrollerv1beta1.KubeAuth{
					Kubeconfig: kubernetesCluster.KubeConfig,
				},
			})
		switch {
		case err != nil && accessError(err):
			k.l.Warn(err)
		case err != nil:
			return err
		case len(xtraDBClusters.Clusters) > 0:
			return status.Errorf(codes.FailedPrecondition, "Kubernetes cluster %s has XtraDB clusters", req.KubernetesClusterName)
		}

		psmdbClusters, err := k.dbaasClient.ListPSMDBClusters(ctx, &dbaascontrollerv1beta1.ListPSMDBClustersRequest{
			KubeAuth: &dbaascontrollerv1beta1.KubeAuth{
				Kubeconfig: kubernetesCluster.KubeConfig,
			},
		})
		switch {
		case err != nil && accessError(err):
			k.l.Warn(err)
		case err != nil:
			return err
		case len(psmdbClusters.Clusters) > 0:
			return status.Errorf(codes.FailedPrecondition, "Kubernetes cluster %s has PSMDB clusters", req.KubernetesClusterName)
		}
		return models.RemoveKubernetesCluster(k.db.Querier, req.KubernetesClusterName)
	})
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.UnregisterKubernetesClusterResponse{}, nil
}

// GetKubernetesCluster return KubeAuth with Kubernetes config.
func (k kubernetesServer) GetKubernetesCluster(ctx context.Context, req *dbaasv1beta1.GetKubernetesClusterRequest) (*dbaasv1beta1.GetKubernetesClusterResponse, error) {
	kubernetesCluster, err := models.FindKubernetesClusterByName(k.db.Querier, req.KubernetesClusterName)
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.GetKubernetesClusterResponse{
		KubeAuth: &dbaasv1beta1.KubeAuth{
			Kubeconfig: kubernetesCluster.KubeConfig,
		},
	}, nil
}

func accessError(err error) bool {
	if err == nil {
		return false
	}
	accessErrors := []*regexp.Regexp{
		operatorIsForbiddenRegexp,
		resourceDoesntExistsRegexp,
	}

	for _, regex := range accessErrors {
		if regex.MatchString(err.Error()) {
			logrus.Warn(err.Error())
			return true
		}
	}
	return false
}

// GetResources returns all and available resources of a Kubernetes cluster.
func (k kubernetesServer) GetResources(ctx context.Context, req *dbaasv1beta1.GetResourcesRequest) (*dbaasv1beta1.GetResourcesResponse, error) {
	kubernetesCluster, err := models.FindKubernetesClusterByName(k.db.Querier, req.KubernetesClusterName)
	if err != nil {
		return nil, err
	}
	in := &dbaascontrollerv1beta1.GetResourcesRequest{
		KubeAuth: &dbaascontrollerv1beta1.KubeAuth{
			Kubeconfig: kubernetesCluster.KubeConfig,
		},
	}
	response, err := k.dbaasClient.GetResources(ctx, in)
	if err != nil {
		return nil, err
	}
	return &dbaasv1beta1.GetResourcesResponse{
		All: &dbaasv1beta1.Resources{
			CpuM:        response.All.CpuM,
			MemoryBytes: response.All.MemoryBytes,
			DiskSize:    response.All.DiskSize,
		},
		Available: &dbaasv1beta1.Resources{
			CpuM:        response.Available.CpuM,
			MemoryBytes: response.Available.MemoryBytes,
			DiskSize:    response.Available.DiskSize,
		},
	}, nil
}
