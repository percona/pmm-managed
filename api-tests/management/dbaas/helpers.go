package dbaas

import (
	"testing"

	dbaasClient "github.com/percona/pmm/api/managementpb/dbaas/json/client"
	"github.com/percona/pmm/api/managementpb/dbaas/json/client/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func registerKubernetesCluster(t *testing.T, kubernetesClusterName string, kubeconfig string) {
	registerKubernetesClusterResponse, err := dbaasClient.Default.Kubernetes.RegisterKubernetesCluster(
		&kubernetes.RegisterKubernetesClusterParams{
			Body: kubernetes.RegisterKubernetesClusterBody{
				KubernetesClusterName: kubernetesClusterName,
				KubeAuth:              &kubernetes.RegisterKubernetesClusterParamsBodyKubeAuth{Kubeconfig: kubeconfig},
			},
			Context: pmmapitests.Context,
		},
	)
	require.NoError(t, err)
	assert.NotNil(t, registerKubernetesClusterResponse)
	t.Cleanup(func() {
		_, _ = unregisterKubernetesCluster(kubernetesClusterName)
	})
}

func unregisterKubernetesCluster(kubernetesClusterName string) (*kubernetes.UnregisterKubernetesClusterOK, error) {
	return dbaasClient.Default.Kubernetes.UnregisterKubernetesCluster(
		&kubernetes.UnregisterKubernetesClusterParams{
			Body:    kubernetes.UnregisterKubernetesClusterBody{KubernetesClusterName: kubernetesClusterName},
			Context: pmmapitests.Context,
		},
	)
}
