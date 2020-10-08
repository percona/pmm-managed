package dbaas

import (
	"fmt"
	"testing"

	dbaasClient "github.com/percona/pmm/api/managementpb/dbaas/json/client"
	"github.com/percona/pmm/api/managementpb/dbaas/json/client/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestKubernetesServer(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		kubernetesClusterName := pmmapitests.TestString(t, "api-test-cluster")
		clusters, err := dbaasClient.Default.Kubernetes.ListKubernetesClusters(nil)
		require.NoError(t, err)
		require.NotContains(t, clusters.Payload.KubernetesClusters, &kubernetes.KubernetesClustersItems0{KubernetesClusterName: kubernetesClusterName})

		registerKubernetesCluster(t, kubernetesClusterName, "")

		clusters, err = dbaasClient.Default.Kubernetes.ListKubernetesClusters(nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(clusters.Payload.KubernetesClusters), 1)
		assert.Contains(t, clusters.Payload.KubernetesClusters, &kubernetes.KubernetesClustersItems0{KubernetesClusterName: kubernetesClusterName})

		unregisterKubernetesClusterResponse, err := dbaasClient.Default.Kubernetes.UnregisterKubernetesCluster(
			&kubernetes.UnregisterKubernetesClusterParams{
				Body:    kubernetes.UnregisterKubernetesClusterBody{KubernetesClusterName: kubernetesClusterName},
				Context: pmmapitests.Context,
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, unregisterKubernetesClusterResponse)

		clusters, err = dbaasClient.Default.Kubernetes.ListKubernetesClusters(nil)
		assert.NoError(t, err)
		require.NotContains(t, clusters.Payload.KubernetesClusters, &kubernetes.KubernetesClustersItems0{KubernetesClusterName: kubernetesClusterName})
	})

	t.Run("DuplicateClusterName", func(t *testing.T) {
		kubernetesClusterName := pmmapitests.TestString(t, "api-test-cluster-duplicate")
		registerKubernetesCluster(t, kubernetesClusterName, "")
		registerKubernetesClusterResponse, err := dbaasClient.Default.Kubernetes.RegisterKubernetesCluster(
			&kubernetes.RegisterKubernetesClusterParams{
				Body: kubernetes.RegisterKubernetesClusterBody{
					KubernetesClusterName: kubernetesClusterName,
					KubeAuth:              &kubernetes.RegisterKubernetesClusterParamsBodyKubeAuth{Kubeconfig: "{}"},
				},
				Context: pmmapitests.Context,
			},
		)
		pmmapitests.AssertAPIErrorf(t, err, 409, codes.AlreadyExists, fmt.Sprintf("Kubernetes Cluster with Name %q already exists.", kubernetesClusterName))
		require.Nil(t, registerKubernetesClusterResponse)
	})

	t.Run("EmptyKubernetesClusterName", func(t *testing.T) {
		registerKubernetesClusterResponse, err := dbaasClient.Default.Kubernetes.RegisterKubernetesCluster(
			&kubernetes.RegisterKubernetesClusterParams{
				Body: kubernetes.RegisterKubernetesClusterBody{
					KubernetesClusterName: "",
					KubeAuth:              &kubernetes.RegisterKubernetesClusterParamsBodyKubeAuth{Kubeconfig: "{}"},
				},
				Context: pmmapitests.Context,
			},
		)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field KubernetesClusterName: value '' must not be an empty string")
		require.Nil(t, registerKubernetesClusterResponse)
	})

	t.Run("EmptyKubeConfig", func(t *testing.T) {
		registerKubernetesClusterResponse, err := dbaasClient.Default.Kubernetes.RegisterKubernetesCluster(
			&kubernetes.RegisterKubernetesClusterParams{
				Body: kubernetes.RegisterKubernetesClusterBody{
					KubernetesClusterName: "empty-kube-config",
					KubeAuth:              &kubernetes.RegisterKubernetesClusterParamsBodyKubeAuth{},
				},
				Context: pmmapitests.Context,
			},
		)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field KubeAuth.Kubeconfig: value '' must not be an empty string")
		require.Nil(t, registerKubernetesClusterResponse)
	})

	t.Run("UnregisterNotExistCluster", func(t *testing.T) {
		unregisterKubernetesClusterOK, err := unregisterKubernetesCluster("not-exist-cluster")
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, "Kubernetes Cluster with name \"not-exist-cluster\" not found.")
		require.Nil(t, unregisterKubernetesClusterOK)
	})

	t.Run("UnregisterEmptyClusterName", func(t *testing.T) {
		unregisterKubernetesClusterOK, err := unregisterKubernetesCluster("")
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field KubernetesClusterName: value '' must not be an empty string")
		require.Nil(t, unregisterKubernetesClusterOK)
	})
}
