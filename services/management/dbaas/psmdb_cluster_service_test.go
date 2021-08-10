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
	"fmt"
	"testing"

	"github.com/google/uuid"
	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

const kubeconfTest = `
	{
		"apiVersion": "v1",
		"kind": "Config",
		"users": [
			{
				"name": "percona-server-mongodb-operator",
				"user": {
					"token": "some-token"
				}
			}
		],
		"clusters": [
			{
				"cluster": {
					"certificate-authority-data": "some-certificate-authority-data",
					"server": "https://192.168.0.42:8443"
				},
				"name": "self-hosted-cluster"
			}
		],
		"contexts": [
			{
				"context": {
					"cluster": "self-hosted-cluster",
					"user": "percona-server-mongodb-operator"
				},
				"name": "svcs-acct-context"
			}
		],
		"current-context": "svcs-acct-context"
	}
`
const kubernetesClusterNameTest = "test-k8s-cluster-name"

func TestPSMDBClusterService(t *testing.T) {
	setup := func(t *testing.T) (ctx context.Context, db *reform.DB, dbaasClient *mockDbaasClient, grafanaClient *mockGrafanaClient, teardown func(t *testing.T)) {
		t.Helper()

		ctx = logger.Set(context.Background(), t.Name())
		uuid.SetRand(new(tests.IDReader))

		sqlDB := testdb.Open(t, models.SetupFixtures, nil)
		db = reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		dbaasClient = new(mockDbaasClient)
		grafanaClient = new(mockGrafanaClient)

		teardown = func(t *testing.T) {
			uuid.SetRand(nil)
			dbaasClient.AssertExpectations(t)
			require.NoError(t, sqlDB.Close())
		}

		return
	}

	ctx, db, dbaasClient, grafanaClient, teardown := setup(t)
	defer teardown(t)

	ks := NewKubernetesServer(db, dbaasClient, grafanaClient)
	dbaasClient.On("CheckKubernetesClusterConnection", ctx, kubeconfTest).Return(&controllerv1beta1.CheckKubernetesClusterConnectionResponse{
		Operators: &controllerv1beta1.Operators{
			Xtradb: &controllerv1beta1.Operator{Status: controllerv1beta1.OperatorsStatus_OPERATORS_STATUS_NOT_INSTALLED},
			Psmdb:  &controllerv1beta1.Operator{Status: controllerv1beta1.OperatorsStatus_OPERATORS_STATUS_OK},
		},
		Status: controllerv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_OK,
	}, nil)

	dbaasClient.On("InstallXtraDBOperator", mock.Anything, mock.Anything).Return(&controllerv1beta1.InstallXtraDBOperatorResponse{}, nil)
	dbaasClient.On("InstallPSMDBOperator", mock.Anything, mock.Anything).Return(&controllerv1beta1.InstallPSMDBOperatorResponse{}, nil)

	registerKubernetesClusterResponse, err := ks.RegisterKubernetesCluster(ctx, &dbaasv1beta1.RegisterKubernetesClusterRequest{
		KubernetesClusterName: kubernetesClusterNameTest,
		KubeAuth:              &dbaasv1beta1.KubeAuth{Kubeconfig: kubeconfTest},
	})
	require.NoError(t, err)
	assert.NotNil(t, registerKubernetesClusterResponse)

	t.Run("BasicListPSMDBClusters", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)
		mockResp := controllerv1beta1.ListPSMDBClustersResponse{
			Clusters: []*controllerv1beta1.ListPSMDBClustersResponse_Cluster{
				{
					Name: "first-psmdb-test",
					Params: &controllerv1beta1.PSMDBClusterParams{
						ClusterSize: 5,
						Replicaset: &controllerv1beta1.PSMDBClusterParams_ReplicaSet{
							ComputeResources: &controllerv1beta1.ComputeResources{
								CpuM:        3,
								MemoryBytes: 256,
							},
						},
					},
					Operation: &controllerv1beta1.RunningOperation{
						TotalSteps:    int32(10),
						FinishedSteps: int32(10),
					},
				},
			},
		}

		dbaasClient.On("ListPSMDBClusters", ctx, mock.Anything).Return(&mockResp, nil)

		resp, err := s.ListPSMDBClusters(ctx, &dbaasv1beta1.ListPSMDBClustersRequest{KubernetesClusterName: kubernetesClusterNameTest})
		assert.NoError(t, err)
		require.NotNil(t, resp.Clusters[0])
		assert.Equal(t, resp.Clusters[0].Name, "first-psmdb-test")
		assert.Equal(t, int32(5), resp.Clusters[0].Params.ClusterSize)
		assert.Equal(t, int32(3), resp.Clusters[0].Params.Replicaset.ComputeResources.CpuM)
		assert.Equal(t, int64(256), resp.Clusters[0].Params.Replicaset.ComputeResources.MemoryBytes)
		assert.Equal(t, int32(10), resp.Clusters[0].Operation.TotalSteps)
		assert.Equal(t, int32(10), resp.Clusters[0].Operation.FinishedSteps)
	})

	//nolint:dupl
	t.Run("BasicCreatePSMDBClusters", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)
		mockReq := controllerv1beta1.CreatePSMDBClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: kubeconfTest,
			},
			Name: "third-psmdb-test",
			Params: &controllerv1beta1.PSMDBClusterParams{
				ClusterSize: 5,
				Replicaset: &controllerv1beta1.PSMDBClusterParams_ReplicaSet{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        3,
						MemoryBytes: 256,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
			},
		}

		dbaasClient.On("CreatePSMDBCluster", ctx, &mockReq).Return(&controllerv1beta1.CreatePSMDBClusterResponse{}, nil)

		in := dbaasv1beta1.CreatePSMDBClusterRequest{
			KubernetesClusterName: kubernetesClusterNameTest,
			Name:                  "third-psmdb-test",
			Params: &dbaasv1beta1.PSMDBClusterParams{
				ClusterSize: 5,
				Replicaset: &dbaasv1beta1.PSMDBClusterParams_ReplicaSet{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        3,
						MemoryBytes: 256,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
			},
		}

		_, err := s.CreatePSMDBCluster(ctx, &in)
		assert.NoError(t, err)
	})

	//nolint:dupl
	t.Run("BasicUpdatePSMDBCluster", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)
		mockReq := controllerv1beta1.UpdatePSMDBClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: kubeconfTest,
			},
			Name: "third-psmdb-test",
			Params: &controllerv1beta1.UpdatePSMDBClusterRequest_UpdatePSMDBClusterParams{
				ClusterSize: 8,
				Replicaset: &controllerv1beta1.UpdatePSMDBClusterRequest_UpdatePSMDBClusterParams_ReplicaSet{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        1,
						MemoryBytes: 256,
					},
				},
			},
		}

		dbaasClient.On("UpdatePSMDBCluster", ctx, &mockReq).Return(&controllerv1beta1.UpdatePSMDBClusterResponse{}, nil)

		in := dbaasv1beta1.UpdatePSMDBClusterRequest{
			KubernetesClusterName: kubernetesClusterNameTest,
			Name:                  "third-psmdb-test",
			Params: &dbaasv1beta1.UpdatePSMDBClusterRequest_UpdatePSMDBClusterParams{
				ClusterSize: 8,
				Replicaset: &dbaasv1beta1.UpdatePSMDBClusterRequest_UpdatePSMDBClusterParams_ReplicaSet{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        1,
						MemoryBytes: 256,
					},
				},
			},
		}

		_, err := s.UpdatePSMDBCluster(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("BasicGetPSMDBClusterCredentials", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)

		mockReq := controllerv1beta1.GetPSMDBClusterCredentialsRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: kubeconfTest,
			},
			Name: "third-psmdb-test",
		}

		dbaasClient.On("GetPSMDBClusterCredentials", ctx, &mockReq).Return(&controllerv1beta1.GetPSMDBClusterCredentialsResponse{
			Credentials: &controllerv1beta1.PSMDBCredentials{
				Username:   "userAdmin",
				Password:   "userAdmin123",
				Host:       "hostname",
				Port:       27017,
				Replicaset: "rs0",
			},
		}, nil)

		in := dbaasv1beta1.GetPSMDBClusterCredentialsRequest{
			KubernetesClusterName: kubernetesClusterNameTest,
			Name:                  "third-psmdb-test",
		}

		cluster, err := s.GetPSMDBClusterCredentials(ctx, &in)

		assert.NoError(t, err)
		assert.Equal(t, "hostname", cluster.ConnectionCredentials.Host)
	})

	t.Run("BasicGetPSMDBClusterCredentialsWithHost", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)
		name := "another-third-psmdb-test"

		mockReq := controllerv1beta1.GetPSMDBClusterCredentialsRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: kubeconfTest,
			},
			Name: name,
		}

		resp := controllerv1beta1.GetPSMDBClusterCredentialsResponse{
			Credentials: &controllerv1beta1.PSMDBCredentials{
				Host: "host",
			},
		}
		dbaasClient.On("GetPSMDBClusterCredentials", ctx, &mockReq).Return(&resp, nil)

		in := dbaasv1beta1.GetPSMDBClusterCredentialsRequest{
			KubernetesClusterName: kubernetesClusterNameTest,
			Name:                  name,
		}

		cluster, err := s.GetPSMDBClusterCredentials(ctx, &in)

		assert.NoError(t, err)
		assert.Equal(t, resp.Credentials.Host, cluster.ConnectionCredentials.Host)
	})

	t.Run("BasicRestartPSMDBCluster", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)
		mockReq := controllerv1beta1.RestartPSMDBClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: kubeconfTest,
			},
			Name: "third-psmdb-test",
		}

		dbaasClient.On("RestartPSMDBCluster", ctx, &mockReq).Return(&controllerv1beta1.RestartPSMDBClusterResponse{}, nil)

		in := dbaasv1beta1.RestartPSMDBClusterRequest{
			KubernetesClusterName: kubernetesClusterNameTest,
			Name:                  "third-psmdb-test",
		}

		_, err := s.RestartPSMDBCluster(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("BasicDeletePSMDBCluster", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)
		dbClusterName := "delete-psmdb-test"
		mockReq := controllerv1beta1.DeletePSMDBClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: kubeconfTest,
			},
			Name: dbClusterName,
		}

		dbaasClient.On("DeletePSMDBCluster", ctx, &mockReq).Return(&controllerv1beta1.DeletePSMDBClusterResponse{}, nil)
		grafanaClient.On("DeleteAPIKeysWithPrefix", ctx, fmt.Sprintf("psmdb-%s-%s", kubernetesClusterNameTest, dbClusterName)).Return(nil)

		in := dbaasv1beta1.DeletePSMDBClusterRequest{
			KubernetesClusterName: kubernetesClusterNameTest,
			Name:                  dbClusterName,
		}

		_, err := s.DeletePSMDBCluster(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("BasicGetPSMDBClusterResources", func(t *testing.T) {
		s := NewPSMDBClusterService(db, dbaasClient, grafanaClient)

		in := dbaasv1beta1.GetPSMDBClusterResourcesRequest{
			Params: &dbaasv1beta1.PSMDBClusterParams{
				ClusterSize: 4,
				Replicaset: &dbaasv1beta1.PSMDBClusterParams_ReplicaSet{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        2000,
						MemoryBytes: 2000000000,
					},
					DiskSize: 2000000000,
				},
			},
		}

		actual, err := s.GetPSMDBClusterResources(ctx, &in)
		assert.NoError(t, err)
		assert.Equal(t, uint64(16000000000), actual.Expected.MemoryBytes)
		assert.Equal(t, uint64(16000), actual.Expected.CpuM)
		assert.Equal(t, uint64(14000000000), actual.Expected.DiskSize)
	})
}
