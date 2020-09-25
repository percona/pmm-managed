package dbaas

import (
	"context"

	dbaasController "github.com/percona-platform/dbaas-api/gen/controller"
	"google.golang.org/grpc"
)

type Client struct {
	kubernetesClient dbaasController.KubernetesClusterAPIClient
}

func NewClient(con grpc.ClientConnInterface) *Client {
	return &Client{
		kubernetesClient: dbaasController.NewKubernetesClusterAPIClient(con),
	}
}

func (c *Client) CheckKubernetesClusterConnection(ctx context.Context, kubeConfig string) error {
	_, err := c.kubernetesClient.CheckKubernetesClusterConnection(ctx, &dbaasController.CheckKubernetesClusterConnectionRequest{
		KubeAuth: &dbaasController.KubeAuth{Kubeconfig: kubeConfig},
	})
	return err
}
