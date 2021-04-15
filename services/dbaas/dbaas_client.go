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
	"google.golang.org/grpc/codes"

	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	"github.com/percona/pmm/version"

	"github.com/percona/pmm-managed/services/supervisord"
)

type handler func() (interface{}, error)

type apiRequest struct {
	handler    handler
	responseCh chan *apiResponse
}
type apiResponse struct {
	out interface{}
	err error
}

type connectRequest struct {
	ctx        context.Context
	responseCh chan errorResponse
}
type disconnectRequest struct {
	responseCh chan errorResponse
}
type errorResponse struct {
	err error
}

// Client is a client for dbaas-controller.
type Client struct {
	l                         *logrus.Entry
	kubernetesClient          controllerv1beta1.KubernetesClusterAPIClient
	xtradbClusterClient       controllerv1beta1.XtraDBClusterAPIClient
	psmdbClusterClient        controllerv1beta1.PSMDBClusterAPIClient
	logsClient                controllerv1beta1.LogsAPIClient
	conn                      *grpc.ClientConn
	dbaasControllerAPIAddress string
	disconnectCh              chan disconnectRequest
	connectCh                 chan connectRequest
	requests                  chan *apiRequest
}

// NewClient creates new Client object.
func NewClient(dbaasControllerAPIAddress string, supervisor *supervisord.Service) *Client {
	c := &Client{
		l:                         logrus.WithField("component", "dbaas.Client"),
		dbaasControllerAPIAddress: dbaasControllerAPIAddress,
		disconnectCh:              make(chan disconnectRequest),
		connectCh:                 make(chan connectRequest),
		requests:                  make(chan *apiRequest, 8),
	}
	events := supervisor.Subscribe("dbaas-controller", supervisord.Running, supervisord.Stopping, supervisord.ExitedUnexpected)
	go c.watchDbaasControllerEvents(events)
	return c
}

func (c *Client) watchDbaasControllerEvents(events chan *supervisord.Event) {
	wg := new(sync.WaitGroup)
	var cancel context.CancelFunc
	for {
		select {
		case r := <-c.disconnectCh:
			c.l.Debug("Disconnecting from dbaas-controller API...")
			wg.Wait()
			c.l.Debug("...all requests done.")

			if c.conn == nil {
				c.l.Warnf("trying to disconnect from dbaas-controller API but the connection is not up")
				if r.responseCh != nil {
					r.responseCh <- errorResponse{}
				}
				break
			}
			err := c.conn.Close()
			if err != nil {
				err := errors.Errorf("failed to close conn to dbaas-controller API: %v", err)
				// Connect request can come from Connect method or from internal watch of supervisord events.
				// That's why we either return error outside or log it.
				if r.responseCh != nil {
					r.responseCh <- errorResponse{err: err}
					break
				}
				c.l.Error(err)
				break
			}
			c.conn = nil
			if r.responseCh != nil {
				r.responseCh <- errorResponse{}
			}
		case r := <-c.connectCh:
			if c.conn != nil {
				c.l.Warnf("trying to connect to dbaas-controller API but connection is already up")
				if r.responseCh != nil {
					r.responseCh <- errorResponse{}
				}
				break
			}
			c.l.Debug("Connecting")
			err := c.connect(r.ctx)
			// If connect request originates from this watch loop, cancel the request's context.
			if cancel != nil {
				cancel()
				cancel = nil
			}
			if err != nil {
				err := errors.Errorf("failed to connect to dbaas-controller API: %v", err)
				// Disconnect request can come from Disconnect method or from internal watch of supervisord events.
				// That's why we either return error outside or log it.
				if r.responseCh != nil {
					r.responseCh <- errorResponse{err: err}
					break
				}
				c.l.Error(err)
			}
			if r.responseCh != nil {
				r.responseCh <- errorResponse{}
			}
		// Handle supervisord events.
		case event := <-events:
			c.l.Debugf("got event type %v", event)
			switch event.Type {
			case supervisord.Running:
				var ctx context.Context
				ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
				c.connectCh <- connectRequest{ctx: ctx}
			case supervisord.Stopping, supervisord.ExitedUnexpected:
				c.disconnectCh <- disconnectRequest{}
			default:
				c.l.Errorf("unsupported event type %v", event)
			}
		// Handle requests to dbaas-controller API.
		case r := <-c.requests:
			if c.conn == nil {
				r.responseCh <- &apiResponse{
					err: grpc.Errorf(codes.Unavailable, "dbaas-controller is not running"),
					out: nil,
				}
				break
			}
			c.l.Debugf("Got request: %v", r)
			wg.Add(1)
			go func(r *apiRequest) {
				out, err := r.handler()
				r.responseCh <- &apiResponse{err: err, out: out}
				wg.Done()
			}(r)
		}
	}
}

func (c *Client) Connect(ctx context.Context) error {
	respCh := make(chan errorResponse)
	c.connectCh <- connectRequest{ctx: ctx, responseCh: respCh}
	resp := <-respCh
	return resp.err
}

func (c *Client) Disconnect() error {
	respCh := make(chan errorResponse)
	c.disconnectCh <- disconnectRequest{responseCh: respCh}
	resp := <-respCh
	return resp.err
}

// connect connects/reconnects to dbaas controller API.
func (c *Client) connect(ctx context.Context) error {
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
		return err
	}
	c.conn = conn

	c.kubernetesClient = controllerv1beta1.NewKubernetesClusterAPIClient(conn)
	c.xtradbClusterClient = controllerv1beta1.NewXtraDBClusterAPIClient(conn)
	c.psmdbClusterClient = controllerv1beta1.NewPSMDBClusterAPIClient(conn)
	c.logsClient = controllerv1beta1.NewLogsAPIClient(conn)

	return nil
}

// CheckKubernetesClusterConnection checks connection with kubernetes cluster.
func (c *Client) CheckKubernetesClusterConnection(ctx context.Context, kubeConfig string) (*controllerv1beta1.CheckKubernetesClusterConnectionResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			in := &controllerv1beta1.CheckKubernetesClusterConnectionRequest{
				KubeAuth: &controllerv1beta1.KubeAuth{
					Kubeconfig: kubeConfig,
				},
			}
			return c.kubernetesClient.CheckKubernetesClusterConnection(ctx, in)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.CheckKubernetesClusterConnectionResponse), resp.err
}

// ListXtraDBClusters returns a list of XtraDB clusters.
func (c *Client) ListXtraDBClusters(ctx context.Context, in *controllerv1beta1.ListXtraDBClustersRequest, opts ...grpc.CallOption) (*controllerv1beta1.ListXtraDBClustersResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.xtradbClusterClient.ListXtraDBClusters(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.ListXtraDBClustersResponse), resp.err
}

// CreateXtraDBCluster creates a new XtraDB cluster.
func (c *Client) CreateXtraDBCluster(ctx context.Context, in *controllerv1beta1.CreateXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.CreateXtraDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.xtradbClusterClient.CreateXtraDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.CreateXtraDBClusterResponse), resp.err
}

// UpdateXtraDBCluster updates existing XtraDB cluster.
func (c *Client) UpdateXtraDBCluster(ctx context.Context, in *controllerv1beta1.UpdateXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.UpdateXtraDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.xtradbClusterClient.UpdateXtraDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.UpdateXtraDBClusterResponse), resp.err
}

// DeleteXtraDBCluster deletes XtraDB cluster.
func (c *Client) DeleteXtraDBCluster(ctx context.Context, in *controllerv1beta1.DeleteXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.DeleteXtraDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.xtradbClusterClient.DeleteXtraDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.DeleteXtraDBClusterResponse), resp.err
}

// RestartXtraDBCluster restarts XtraDB cluster.
func (c *Client) RestartXtraDBCluster(ctx context.Context, in *controllerv1beta1.RestartXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.RestartXtraDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.xtradbClusterClient.RestartXtraDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.RestartXtraDBClusterResponse), resp.err
}

// GetXtraDBClusterCredentials gets XtraDB cluster credentials.
func (c *Client) GetXtraDBClusterCredentials(ctx context.Context, in *controllerv1beta1.GetXtraDBClusterCredentialsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetXtraDBClusterCredentialsResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.xtradbClusterClient.GetXtraDBClusterCredentials(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.GetXtraDBClusterCredentialsResponse), resp.err
}

// ListPSMDBClusters returns a list of PSMDB clusters.
func (c *Client) ListPSMDBClusters(ctx context.Context, in *controllerv1beta1.ListPSMDBClustersRequest, opts ...grpc.CallOption) (*controllerv1beta1.ListPSMDBClustersResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.psmdbClusterClient.ListPSMDBClusters(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.ListPSMDBClustersResponse), resp.err
}

// CreatePSMDBCluster creates a new PSMDB cluster.
func (c *Client) CreatePSMDBCluster(ctx context.Context, in *controllerv1beta1.CreatePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.CreatePSMDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.psmdbClusterClient.CreatePSMDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.CreatePSMDBClusterResponse), resp.err
}

// UpdatePSMDBCluster updates existing PSMDB cluster.
func (c *Client) UpdatePSMDBCluster(ctx context.Context, in *controllerv1beta1.UpdatePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.UpdatePSMDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.psmdbClusterClient.UpdatePSMDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.UpdatePSMDBClusterResponse), resp.err
}

// DeletePSMDBCluster deletes PSMDB cluster.
func (c *Client) DeletePSMDBCluster(ctx context.Context, in *controllerv1beta1.DeletePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.DeletePSMDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.psmdbClusterClient.DeletePSMDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.DeletePSMDBClusterResponse), resp.err
}

// RestartPSMDBCluster restarts PSMDB cluster.
func (c *Client) RestartPSMDBCluster(ctx context.Context, in *controllerv1beta1.RestartPSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.RestartPSMDBClusterResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.psmdbClusterClient.RestartPSMDBCluster(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.RestartPSMDBClusterResponse), resp.err
}

// GetPSMDBClusterCredentials gets PSMDB cluster credentials.
func (c *Client) GetPSMDBClusterCredentials(ctx context.Context, in *controllerv1beta1.GetPSMDBClusterCredentialsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetPSMDBClusterCredentialsResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.psmdbClusterClient.GetPSMDBClusterCredentials(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.GetPSMDBClusterCredentialsResponse), resp.err
}

// GetLogs gets logs out of cluster containers and events out of pods.
func (c *Client) GetLogs(ctx context.Context, in *controllerv1beta1.GetLogsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetLogsResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.logsClient.GetLogs(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.GetLogsResponse), resp.err
}

// GetResources returns all and available resources of a Kubernetes cluster.
func (c *Client) GetResources(ctx context.Context, in *controllerv1beta1.GetResourcesRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetResourcesResponse, error) {
	responseCh := make(chan *apiResponse)
	c.requests <- &apiRequest{
		handler: func() (interface{}, error) {
			return c.kubernetesClient.GetResources(ctx, in, opts...)
		},
		responseCh: responseCh,
	}
	resp := <-responseCh
	return resp.out.(*controllerv1beta1.GetResourcesResponse), resp.err
}
