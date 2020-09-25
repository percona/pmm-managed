package dbaas

import "context"

type dbaasClient interface {
	CheckKubernetesClusterConnection(ctx context.Context, kubeConfig string) error
}
