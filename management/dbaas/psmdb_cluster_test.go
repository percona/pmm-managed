package dbaas

import (
	"testing"

	dbaasClient "github.com/percona/pmm/api/managementpb/dbaas/json/client"
	psmdbcluster "github.com/percona/pmm/api/managementpb/dbaas/json/client/psmdb_cluster"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

const (
	psmdbKubernetesClusterName = "api-test-k8s-mongodb-cluster"
)

//nolint:funlen
func TestPSMDBClusterServer(t *testing.T) {
	if pmmapitests.Kubeconfig == "" {
		t.Skip("Skip tests of PSMDBClusterServer without kubeconfig")
	}
	registerKubernetesCluster(t, psmdbKubernetesClusterName, pmmapitests.Kubeconfig)

	t.Run("BasicPSMDBCluster", func(t *testing.T) {
		paramsFirstPSMDB := psmdbcluster.CreatePSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.CreatePSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "first-psmdb-test",
				Params: &psmdbcluster.CreatePSMDBClusterParamsBodyParams{
					ClusterSize: 3,
					Replicaset: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicaset{
						ComputeResources: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicasetComputeResources{
							CPUm:        500,
							MemoryBytes: "1000000000",
						},
						DiskSize: "1000000000",
					},
				},
			},
		}

		_, err := dbaasClient.Default.PSMDBCluster.CreatePSMDBCluster(&paramsFirstPSMDB)
		assert.NoError(t, err)
		// Create one more PSMDB Cluster.
		paramsSecondPSMDB := psmdbcluster.CreatePSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.CreatePSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "second-psmdb-test",
				Params: &psmdbcluster.CreatePSMDBClusterParamsBodyParams{
					ClusterSize: 1,
					Replicaset: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicaset{
						ComputeResources: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicasetComputeResources{
							CPUm:        500,
							MemoryBytes: "1000000000",
						},
						DiskSize: "1000000000",
					},
				},
			},
		}
		_, err = dbaasClient.Default.PSMDBCluster.CreatePSMDBCluster(&paramsSecondPSMDB)
		assert.NoError(t, err)

		listPSMDBClustersParamsParam := psmdbcluster.ListPSMDBClustersParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.ListPSMDBClustersBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
			},
		}
		xtraDBClusters, err := dbaasClient.Default.PSMDBCluster.ListPSMDBClusters(&listPSMDBClustersParamsParam)
		assert.NoError(t, err)

		for _, name := range []string{"first-psmdb-test", "second-psmdb-test"} {
			foundPSMDB := false
			for _, psmdb := range xtraDBClusters.Payload.Clusters {
				if name == psmdb.Name {
					foundPSMDB = true

					break
				}
			}
			assert.True(t, foundPSMDB, "Cannot find PSMDB with name %s in cluster list", name)
		}

		paramsUpdatePSMDB := psmdbcluster.UpdatePSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.UpdatePSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "second-psmdb-test",
				Params: &psmdbcluster.UpdatePSMDBClusterParamsBodyParams{
					ClusterSize: 2,
					Replicaset: &psmdbcluster.UpdatePSMDBClusterParamsBodyParamsReplicaset{
						ComputeResources: &psmdbcluster.UpdatePSMDBClusterParamsBodyParamsReplicasetComputeResources{
							CPUm:        2,
							MemoryBytes: "128",
						},
					},
				},
			},
		}

		_, err = dbaasClient.Default.PSMDBCluster.UpdatePSMDBCluster(&paramsUpdatePSMDB)
		pmmapitests.AssertAPIErrorf(t, err, 500, codes.Internal, `state is initializing: PSMDB cluster is not ready`)

		for _, psmdb := range xtraDBClusters.Payload.Clusters {
			if psmdb.Name == "" {
				continue
			}
			deletePSMDBClusterParamsParam := psmdbcluster.DeletePSMDBClusterParams{
				Context: pmmapitests.Context,
				Body: psmdbcluster.DeletePSMDBClusterBody{
					KubernetesClusterName: psmdbKubernetesClusterName,
					Name:                  psmdb.Name,
				},
			}
			_, err := dbaasClient.Default.PSMDBCluster.DeletePSMDBCluster(&deletePSMDBClusterParamsParam)
			assert.NoError(t, err)
		}

		cluster, err := dbaasClient.Default.PSMDBCluster.GetPSMDBCluster(&psmdbcluster.GetPSMDBClusterParams{
			Body: psmdbcluster.GetPSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "second-psmdb-test",
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		assert.Equal(t, &psmdbcluster.GetPSMDBClusterOKBodyConnectionCredentials{
			Username:   "userAdmin",
			Password:   "userAdmin123456",
			Host:       "second-psmdb-test-rs0.default.svc.cluster.local",
			Port:       27017,
			Replicaset: "rs0",
		}, cluster.Payload.ConnectionCredentials)

		t.Skip("Skip restart till better implementation. https://jira.percona.com/browse/PMM-6980")
		_, err = dbaasClient.Default.PSMDBCluster.RestartPSMDBCluster(&psmdbcluster.RestartPSMDBClusterParams{
			Body: psmdbcluster.RestartPSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "first-psmdb-test",
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
	})

	t.Run("CreatePSMDBClusterEmptyName", func(t *testing.T) {
		paramsPSMDBEmptyName := psmdbcluster.CreatePSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.CreatePSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "",
				Params: &psmdbcluster.CreatePSMDBClusterParamsBodyParams{
					ClusterSize: 3,
					Replicaset: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicaset{
						ComputeResources: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicasetComputeResources{
							CPUm:        1,
							MemoryBytes: "64",
						},
					},
				},
			},
		}
		_, err := dbaasClient.Default.PSMDBCluster.CreatePSMDBCluster(&paramsPSMDBEmptyName)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `invalid field Name: value '' must be a string conforming to regex "^[a-z]([-a-z0-9]*[a-z0-9])?$"`)
	})

	t.Run("CreatePSMDBClusterInvalidName", func(t *testing.T) {
		paramsPSMDBInvalidName := psmdbcluster.CreatePSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.CreatePSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "123_asd",
				Params: &psmdbcluster.CreatePSMDBClusterParamsBodyParams{
					ClusterSize: 3,
					Replicaset: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicaset{
						ComputeResources: &psmdbcluster.CreatePSMDBClusterParamsBodyParamsReplicasetComputeResources{
							CPUm:        1,
							MemoryBytes: "64",
						},
					},
				},
			},
		}
		_, err := dbaasClient.Default.PSMDBCluster.CreatePSMDBCluster(&paramsPSMDBInvalidName)
		assert.Error(t, err)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, `invalid field Name: value '123_asd' must be a string conforming to regex "^[a-z]([-a-z0-9]*[a-z0-9])?$"`)
	})

	t.Run("ListUnknownCluster", func(t *testing.T) {
		listPSMDBClustersParamsParam := psmdbcluster.ListPSMDBClustersParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.ListPSMDBClustersBody{
				KubernetesClusterName: "Unknown-kubernetes-cluster-name",
			},
		}
		_, err := dbaasClient.Default.PSMDBCluster.ListPSMDBClusters(&listPSMDBClustersParamsParam)
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, `Kubernetes Cluster with name "Unknown-kubernetes-cluster-name" not found.`)
	})

	t.Run("RestartUnknownPSMDBCluster", func(t *testing.T) {
		restartPSMDBClusterParamsParam := psmdbcluster.RestartPSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.RestartPSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "Unknown-psmdb-name",
			},
		}
		_, err := dbaasClient.Default.PSMDBCluster.RestartPSMDBCluster(&restartPSMDBClusterParamsParam)
		require.Error(t, err)
		assert.Equal(t, 500, err.(pmmapitests.ErrorResponse).Code())
	})

	t.Run("DeleteUnknownPSMDBCluster", func(t *testing.T) {
		deletePSMDBClusterParamsParam := psmdbcluster.DeletePSMDBClusterParams{
			Context: pmmapitests.Context,
			Body: psmdbcluster.DeletePSMDBClusterBody{
				KubernetesClusterName: psmdbKubernetesClusterName,
				Name:                  "Unknown-psmdb-name",
			},
		}
		_, err := dbaasClient.Default.PSMDBCluster.DeletePSMDBCluster(&deletePSMDBClusterParamsParam)
		require.Error(t, err)
		assert.Equal(t, 500, err.(pmmapitests.ErrorResponse).Code())
	})
}
