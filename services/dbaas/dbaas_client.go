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

// Package dbaas contains logic related to communication with dbaas-controller.
package dbaas

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	"github.com/percona/pmm/version"
)

// Client is a client for dbaas-controller.
type Client struct {
	l                         *logrus.Entry
	kubernetesClient          controllerv1beta1.KubernetesClusterAPIClient
	xtradbClusterClient       controllerv1beta1.XtraDBClusterAPIClient
	psmdbClusterClient        controllerv1beta1.PSMDBClusterAPIClient
	logsClient                controllerv1beta1.LogsAPIClient
	connM                     *sync.RWMutex
	conn                      *grpc.ClientConn
	dbaasControllerAPIAddress string
	wg                        *sync.WaitGroup
}

// NewClient creates new Client object.
func NewClient(dbaasControllerAPIAddress string) *Client {
	c := &Client{
		l:                         logrus.WithField("component", "dbaas.Client"),
		connM:                     new(sync.RWMutex),
		dbaasControllerAPIAddress: dbaasControllerAPIAddress,
		wg:                        new(sync.WaitGroup),
	}
	return c
}

// Connect connects the client to dbaas-controller API.
func (c *Client) Connect(ctx context.Context) error {
	c.connM.Lock()
	defer c.connM.Unlock()
	c.l.Infof("Connecting to dbaas-controller API on %s.", c.dbaasControllerAPIAddress)
	if c.conn != nil {
		c.l.Warnf("Trying to connect to dbaas-controller API but connection is already up.")
		return nil
	}
	backoffConfig := backoff.DefaultConfig
	backoffConfig.MaxDelay = 10 * time.Second
	opts := []grpc.DialOption{
		grpc.WithBlock(), // Dial blocks, we do not connect in background.
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoffConfig, MinConnectTimeout: 10 * time.Second}),
		grpc.WithUserAgent("pmm-managed/" + version.Version),
	}

	conn, err := grpc.DialContext(ctx, c.dbaasControllerAPIAddress, opts...)
	if err != nil {
		return errors.Errorf("failed to connect to dbaas-controller API: %v", err)
	}
	c.conn = conn

	c.kubernetesClient = controllerv1beta1.NewKubernetesClusterAPIClient(conn)
	c.xtradbClusterClient = controllerv1beta1.NewXtraDBClusterAPIClient(conn)
	c.psmdbClusterClient = controllerv1beta1.NewPSMDBClusterAPIClient(conn)
	c.logsClient = controllerv1beta1.NewLogsAPIClient(conn)

	c.l.Info("Connected to dbaas-controller API.")
	return nil
}

// Disconnect disconnects the client from dbaas-controller API.
func (c *Client) Disconnect() error {
	c.connM.Lock()
	defer c.connM.Unlock()
	c.l.Info("Disconnecting from dbaas-controller API.")
	c.wg.Wait()
	if c.conn == nil {
		c.l.Warnf("Trying to disconnect from dbaas-controller API but the connection is not up.")
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return errors.Errorf("failed to close conn to dbaas-controller API: %v", err)
	}
	c.conn = nil
	c.l.Info("Disconected from dbaas-controller API.")
	return nil
}

// CheckKubernetesClusterConnection checks connection with kubernetes cluster.
func (c *Client) CheckKubernetesClusterConnection(ctx context.Context, kubeConfig string) (*controllerv1beta1.CheckKubernetesClusterConnectionResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)

	in := &controllerv1beta1.CheckKubernetesClusterConnectionRequest{
		KubeAuth: &controllerv1beta1.KubeAuth{
			Kubeconfig: kubeConfig,
		},
	}
	out, err := c.kubernetesClient.CheckKubernetesClusterConnection(ctx, in)
	c.wg.Done()
	return out, err
}

func (c *Client) ListXtraDBClusters(ctx context.Context, in *controllerv1beta1.ListXtraDBClustersRequest, opts ...grpc.CallOption) (*controllerv1beta1.ListXtraDBClustersResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.xtradbClusterClient.ListXtraDBClusters(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// CreateXtraDBCluster creates a new XtraDB cluster.
func (c *Client) CreateXtraDBCluster(ctx context.Context, in *controllerv1beta1.CreateXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.CreateXtraDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.xtradbClusterClient.CreateXtraDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// UpdateXtraDBCluster updates existing XtraDB cluster.
func (c *Client) UpdateXtraDBCluster(ctx context.Context, in *controllerv1beta1.UpdateXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.UpdateXtraDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.xtradbClusterClient.UpdateXtraDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// DeleteXtraDBCluster deletes XtraDB cluster.
func (c *Client) DeleteXtraDBCluster(ctx context.Context, in *controllerv1beta1.DeleteXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.DeleteXtraDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.xtradbClusterClient.DeleteXtraDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// RestartXtraDBCluster restarts XtraDB cluster.
func (c *Client) RestartXtraDBCluster(ctx context.Context, in *controllerv1beta1.RestartXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.RestartXtraDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.xtradbClusterClient.RestartXtraDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// GetXtraDBClusterCredentials gets XtraDB cluster credentials.
func (c *Client) GetXtraDBClusterCredentials(ctx context.Context, in *controllerv1beta1.GetXtraDBClusterCredentialsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetXtraDBClusterCredentialsResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.xtradbClusterClient.GetXtraDBClusterCredentials(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// ListPSMDBClusters returns a list of PSMDB clusters.
func (c *Client) ListPSMDBClusters(ctx context.Context, in *controllerv1beta1.ListPSMDBClustersRequest, opts ...grpc.CallOption) (*controllerv1beta1.ListPSMDBClustersResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.psmdbClusterClient.ListPSMDBClusters(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// CreatePSMDBCluster creates a new PSMDB cluster.
func (c *Client) CreatePSMDBCluster(ctx context.Context, in *controllerv1beta1.CreatePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.CreatePSMDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.psmdbClusterClient.CreatePSMDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// UpdatePSMDBCluster updates existing PSMDB cluster.
func (c *Client) UpdatePSMDBCluster(ctx context.Context, in *controllerv1beta1.UpdatePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.UpdatePSMDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.psmdbClusterClient.UpdatePSMDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// DeletePSMDBCluster deletes PSMDB cluster.
func (c *Client) DeletePSMDBCluster(ctx context.Context, in *controllerv1beta1.DeletePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.DeletePSMDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.psmdbClusterClient.DeletePSMDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// RestartPSMDBCluster restarts PSMDB cluster.
func (c *Client) RestartPSMDBCluster(ctx context.Context, in *controllerv1beta1.RestartPSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.RestartPSMDBClusterResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.psmdbClusterClient.RestartPSMDBCluster(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// GetPSMDBClusterCredentials gets PSMDB cluster credentials.
func (c *Client) GetPSMDBClusterCredentials(ctx context.Context, in *controllerv1beta1.GetPSMDBClusterCredentialsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetPSMDBClusterCredentialsResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.psmdbClusterClient.GetPSMDBClusterCredentials(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// GetLogs gets logs out of cluster containers and events out of pods.
func (c *Client) GetLogs(ctx context.Context, in *controllerv1beta1.GetLogsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetLogsResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.logsClient.GetLogs(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}

// GetResources returns all and available resources of a Kubernetes cluster.
func (c *Client) GetResources(ctx context.Context, in *controllerv1beta1.GetResourcesRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetResourcesResponse, error) {
	c.connM.RLock()
	defer c.connM.RUnlock()
	c.wg.Add(1)
	resp, err := c.kubernetesClient.GetResources(ctx, in, opts...)
	c.wg.Done()
	return resp, err
}
