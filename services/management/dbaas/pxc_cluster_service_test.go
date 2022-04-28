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
	"testing"

	"github.com/google/uuid"
	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/percona/pmm-managed/utils/testdb"
	"github.com/percona/pmm-managed/utils/tests"
)

const pxcKubeconfigTest = `
{
	"apiVersion": "v1",
	"kind": "Config",
	"users": [
		{
			"name": "percona-xtradb-cluster-operator",
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
				"user": "percona-xtradb-cluster-operator"
			},
			"name": "svcs-acct-context"
		}
	],
	"current-context": "svcs-acct-context"
}
`
const pxcKubernetesClusterNameTest = "test-k8s-cluster-name"

func TestPXCClusterService(t *testing.T) {
	setup := func(t *testing.T) (ctx context.Context, db *reform.DB, dbaasClient *mockDbaasClient, grafanaClient *mockGrafanaClient,
		componentsService *mockComponentsService, teardown func(t *testing.T),
	) {
		t.Helper()

		ctx = logger.Set(context.Background(), t.Name())
		uuid.SetRand(&tests.IDReader{})

		sqlDB := testdb.Open(t, models.SetupFixtures, nil)
		db = reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(t.Logf))
		dbaasClient = &mockDbaasClient{}
		grafanaClient = &mockGrafanaClient{}
		componentsService = &mockComponentsService{}

		teardown = func(t *testing.T) {
			uuid.SetRand(nil)
			dbaasClient.AssertExpectations(t)
		}

		return
	}

	ctx, db, dbaasClient, grafanaClient, componentsClient, teardown := setup(t)
	defer teardown(t)
	versionService := NewVersionServiceClient(versionServiceURL)

	ks := NewKubernetesServer(db, dbaasClient, grafanaClient, versionService)
	dbaasClient.On("CheckKubernetesClusterConnection", ctx, pxcKubeconfigTest).Return(&controllerv1beta1.CheckKubernetesClusterConnectionResponse{
		Operators: &controllerv1beta1.Operators{
			PxcOperatorVersion:   "",
			PsmdbOperatorVersion: onePointEight,
		},
		Status: controllerv1beta1.KubernetesClusterStatus_KUBERNETES_CLUSTER_STATUS_OK,
	}, nil)

	dbaasClient.On("InstallPXCOperator", mock.Anything, mock.Anything).Return(&controllerv1beta1.InstallPXCOperatorResponse{}, nil)

	registerKubernetesClusterResponse, err := ks.RegisterKubernetesCluster(ctx, &dbaasv1beta1.RegisterKubernetesClusterRequest{
		KubernetesClusterName: pxcKubernetesClusterNameTest,
		KubeAuth:              &dbaasv1beta1.KubeAuth{Kubeconfig: pxcKubeconfigTest},
	})
	require.NoError(t, err)
	assert.NotNil(t, registerKubernetesClusterResponse)

	//nolint:dupl
	t.Run("BasicCreatePXCClusters", func(t *testing.T) {
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
		mockReq := controllerv1beta1.CreatePXCClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: "third-pxc-test",
			Params: &controllerv1beta1.PXCClusterParams{
				ClusterSize: 5,
				Pxc: &controllerv1beta1.PXCClusterParams_PXC{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        3,
						MemoryBytes: 256,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
				Proxysql: &controllerv1beta1.PXCClusterParams_ProxySQL{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        2,
						MemoryBytes: 124,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
				VersionServiceUrl: versionService.GetVersionServiceURL(),
			},
		}

		dbaasClient.On("CreatePXCCluster", ctx, &mockReq).Return(&controllerv1beta1.CreatePXCClusterResponse{}, nil)

		in := dbaasv1beta1.CreatePXCClusterRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  "third-pxc-test",
			Params: &dbaasv1beta1.PXCClusterParams{
				ClusterSize: 5,
				Pxc: &dbaasv1beta1.PXCClusterParams_PXC{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        3,
						MemoryBytes: 256,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
				Proxysql: &dbaasv1beta1.PXCClusterParams_ProxySQL{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        2,
						MemoryBytes: 124,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
			},
		}

		_, err := s.CreatePXCCluster(ctx, &in)
		assert.NoError(t, err)
	})
	t.Run("BasicCreatePXCClusters", func(t *testing.T) {
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
		mockReq := controllerv1beta1.CreatePXCClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: "third-pxc-test",
			Params: &controllerv1beta1.PXCClusterParams{
				ClusterSize: 5,
				Pxc: &controllerv1beta1.PXCClusterParams_PXC{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        3,
						MemoryBytes: 256,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
				Proxysql: &controllerv1beta1.PXCClusterParams_ProxySQL{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        2,
						MemoryBytes: 124,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
				VersionServiceUrl: versionService.GetVersionServiceURL(),
			},
		}

		dbaasClient.On("CreatePXCCluster", ctx, &mockReq).Return(&controllerv1beta1.CreatePXCClusterResponse{}, nil)
		componentsClient.On("GetPXCComponents", ctx, mock.Anything).Return(
			&dbaasv1beta1.GetPXCComponentsResponse{
				Versions: []*dbaasv1beta1.OperatorVersion{
					{
						Product:  "pxc-operator",
						Operator: "1.10.0",
						Matrix: &dbaasv1beta1.Matrix{
							Pxc: map[string]*dbaasv1beta1.Component{
								"8.0.19-10.1": {
									ImagePath: "percona/percona-xtradb-cluster:8.0.19-10.1",
									ImageHash: "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a",
									Status:    "available",
									Critical:  false,
									Default:   false,
									Disabled:  false,
								},
								"8.0.20-11.1": {
									ImagePath: "percona/percona-xtradb-cluster:8.0.20-11.1",
									ImageHash: "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349",
									Status:    "available",
									Critical:  false,
									Default:   false,
									Disabled:  false,
								},
								"8.0.25-15.1": {
									ImagePath: "percona/percona-xtradb-cluster:8.0.25-15.1",
									ImageHash: "529e979c86442429e6feabef9a2d9fc362f4626146f208fbfac704e145a492dd",
									Status:    "recommended",
									Critical:  false,
									Default:   true,
									Disabled:  false,
								},
							},
						},
					},
				},
			},
		)
		in := dbaasv1beta1.CreatePXCClusterRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  "third-pxc-test",
			Params: &dbaasv1beta1.PXCClusterParams{
				ClusterSize: 5,
				Pxc: &dbaasv1beta1.PXCClusterParams_PXC{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        3,
						MemoryBytes: 256,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
				Proxysql: &dbaasv1beta1.PXCClusterParams_ProxySQL{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        2,
						MemoryBytes: 124,
					},
					DiskSize: 1024 * 1024 * 1024,
				},
			},
		}

		_, err := s.CreatePXCCluster(ctx, &in)
		assert.NoError(t, err)
	})
	/*
		&dbaasv1beta1.GetPXCComponentsResponse{
			state:         impl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			Versions:      {
				&dbaasv1beta1.OperatorVersion{
					state:         impl.MessageState{},
					sizeCache:     0,
					unknownFields: nil,
					Product:       "pxc-operator",
					Operator:      "1.10.0",
					Matrix:        &dbaasv1beta1.Matrix{
						state:         impl.MessageState{},
						sizeCache:     0,
						unknownFields: nil,
						Mongod:        {
						},
						Pxc: {
							"8.0.19-10.1": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.19-10.1",
								ImageHash:     "1058ae8eded735ebdf664807aad7187942fc9a1170b3fd0369574cb61206b63a",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.20-11.1": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.20-11.1",
								ImageHash:     "54b1b2f5153b78b05d651034d4603a13e685cbb9b45bfa09a39864fa3f169349",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.20-11.2": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.20-11.2",
								ImageHash:     "feda5612db18da824e971891d6084465aa9cdc9918c18001cd95ba30916da78b",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.21-12.1": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.21-12.1",
								ImageHash:     "d95cf39a58f09759408a00b519fe0d0b19c1b28332ece94349dd5e9cdbda017e",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.22-13.1": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.22-13.1",
								ImageHash:     "1295af1153c1d02e9d40131eb0945b53f7f371796913e64116bf2caa77dc186d",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.23-14.1": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.23-14.1",
								ImageHash:     "8109f7ca4fc465ba862c08021df12e77b65d384395078e31e270d14b77810d79",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.25-15.1": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster:8.0.25-15.1",
								ImageHash:     "529e979c86442429e6feabef9a2d9fc362f4626146f208fbfac704e145a492dd",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
						Pmm: {
							"2.23.0": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/pmm-client:2.23.0",
								ImageHash:     "8fa0e45f740fa8564cbfbdf5d9a5507a07e331f8f40ea022d3a64d7278478eac",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
						Proxysql: {
							"2.0.18": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-proxysql",
								ImageHash:     "f109a62eb316732d59dd80ed0e013fc9594cbae601586b94023b8c068f7ced7b",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"2.0.18-2": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-proxysql-8.0.25",
								ImageHash:     "b84701c47a11c6f5ca46481f25f1b6086c0a30014d05584c7987f1d42a17b584",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
						Haproxy: {
							"2.3.14": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-haproxy",
								ImageHash:     "2f06ac4a0f39b2c0253421c3d024291d5ba19d41e35e633ff6ddcf4ba67fd51a",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"2.3.15": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-haproxy-8.0.25",
								ImageHash:     "62479be2a21192a3215f03d3f9541decd5ef1737741245ac33ee439915a15128",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
						Backup: {
							"2.4.24": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-pxc5.7-backup",
								ImageHash:     "2ff5992220ba251cf064cc2b4d5929e0fdb963db18e35d6c672f9aacb0be3bed",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"2.4.24-2": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-pxc5.7.35-backup",
								ImageHash:     "ac9fcd3078107c6492c687eb98215d4e5daf27a02fb3c78ba4b9e9c01f2078b3",
								Status:        "recommended",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.23": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-pxc8.0-backup",
								ImageHash:     "6ab8efb3804d1e519e49ee10eb46b428a837cfdcee222cc5ae2089cc1dc02a6d",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"8.0.25": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-pxc8.0.25-backup",
								ImageHash:     "c3991f0959a3b4114d7ff629d9d3cdf0dc200c58443ca8ebb1446d8b1cbe416d",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
						Operator: {
							"1.10.0": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0",
								ImageHash:     "73d2266258b700a691db6196f4b5c830845d34d57bdef5be5ffbd45e88407309",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
						LogCollector: {
							"1.10.0": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-1-logcollector",
								ImageHash:     "8f106b1e9134812b77f4e210ad0fcd7d8d3515a90fe53554d24cd49defc9e044",
								Status:        "available",
								Critical:      false,
								Default:       false,
								Disabled:      false,
							},
							"1.10.0-2": &dbaasv1beta1.Component{
								state:         impl.MessageState{},
								sizeCache:     0,
								unknownFields: nil,
								ImagePath:     "percona/percona-xtradb-cluster-operator:1.10.0-logcollector-8.0.25",
								ImageHash:     "d69dad98900532e2ad6d0bf12c34a148462816fa3ee4697e9b73efef7583901a",
								Status:        "recommended",
								Critical:      false,
								Default:       true,
								Disabled:      false,
							},
						},
					},
				},
			},
		}
	*/
	t.Run("CreatePXCClusterMinimumParams", func(t *testing.T) {
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())

		dbaasClient.On("CreatePXCCluster", ctx, mock.Anything).Return(&controllerv1beta1.CreatePXCClusterResponse{}, nil)

		in := dbaasv1beta1.CreatePXCClusterRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  "fourth-pxc-test",
		}

		_, err := s.CreatePXCCluster(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("BasicGetPXCClusterCredentials", func(t *testing.T) {
		name := "third-pxc-test"
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
		mockReq := controllerv1beta1.GetPXCClusterCredentialsRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: name,
		}

		dbaasClient.On("GetPXCClusterCredentials", ctx, &mockReq).Return(&controllerv1beta1.GetPXCClusterCredentialsResponse{
			Credentials: &controllerv1beta1.PXCCredentials{
				Username: "root",
				Password: "root_password",
				Host:     "hostname",
				Port:     3306,
			},
		}, nil)

		in := dbaasv1beta1.GetPXCClusterCredentialsRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  name,
		}

		actual, err := s.GetPXCClusterCredentials(ctx, &in)
		assert.NoError(t, err)
		assert.Equal(t, actual.ConnectionCredentials.Username, "root")
		assert.Equal(t, actual.ConnectionCredentials.Password, "root_password")
		assert.Equal(t, actual.ConnectionCredentials.Host, "hostname", name)
		assert.Equal(t, actual.ConnectionCredentials.Port, int32(3306))
	})

	t.Run("BasicGetPXCClusterCredentialsWithHost", func(t *testing.T) { // Real kubernetes will have ingress
		name := "another-third-pxc-test"
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
		mockReq := controllerv1beta1.GetPXCClusterCredentialsRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: name,
		}

		mockCluster := &controllerv1beta1.GetPXCClusterCredentialsResponse{
			Credentials: &controllerv1beta1.PXCCredentials{
				Username: "root",
				Password: "root_password",
				Host:     "amazing.com",
				Port:     3306,
			},
		}

		dbaasClient.On("GetPXCClusterCredentials", ctx, &mockReq).Return(mockCluster, nil)

		in := dbaasv1beta1.GetPXCClusterCredentialsRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  name,
		}

		actual, err := s.GetPXCClusterCredentials(ctx, &in)
		assert.NoError(t, err)
		assert.Equal(t, "root", actual.ConnectionCredentials.Username)
		assert.Equal(t, "root_password", actual.ConnectionCredentials.Password)
		assert.Equal(t, mockCluster.Credentials.Host, actual.ConnectionCredentials.Host)
		assert.Equal(t, int32(3306), actual.ConnectionCredentials.Port)
	})

	//nolint:dupl
	t.Run("BasicUpdatePXCCluster", func(t *testing.T) {
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
		mockReq := controllerv1beta1.UpdatePXCClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: "third-pxc-test",
			Params: &controllerv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams{
				ClusterSize: 8,
				Pxc: &controllerv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams_PXC{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        1,
						MemoryBytes: 256,
					},
					Image: "path",
				},
				Proxysql: &controllerv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams_ProxySQL{
					ComputeResources: &controllerv1beta1.ComputeResources{
						CpuM:        1,
						MemoryBytes: 124,
					},
				},
			},
		}

		dbaasClient.On("UpdatePXCCluster", ctx, &mockReq).Return(&controllerv1beta1.UpdatePXCClusterResponse{}, nil)

		in := dbaasv1beta1.UpdatePXCClusterRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  "third-pxc-test",
			Params: &dbaasv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams{
				ClusterSize: 8,
				Pxc: &dbaasv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams_PXC{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        1,
						MemoryBytes: 256,
					},
					Image: "path",
				},
				Proxysql: &dbaasv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams_ProxySQL{
					ComputeResources: &dbaasv1beta1.ComputeResources{
						CpuM:        1,
						MemoryBytes: 124,
					},
				},
			},
		}

		_, err := s.UpdatePXCCluster(ctx, &in)
		assert.NoError(t, err)
	})

	//nolint:dupl
	t.Run("BasicSuspendResumePXCCluster", func(t *testing.T) {
		s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
		mockReqSuspend := controllerv1beta1.UpdatePXCClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: "forth-pxc-test",
			Params: &controllerv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams{
				Suspend: true,
			},
		}

		dbaasClient.On("UpdatePXCCluster", ctx, &mockReqSuspend).Return(&controllerv1beta1.UpdatePXCClusterResponse{}, nil)

		in := dbaasv1beta1.UpdatePXCClusterRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  "forth-pxc-test",
			Params: &dbaasv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams{
				Suspend: true,
			},
		}
		_, err := s.UpdatePXCCluster(ctx, &in)
		assert.NoError(t, err)

		mockReqResume := controllerv1beta1.UpdatePXCClusterRequest{
			KubeAuth: &controllerv1beta1.KubeAuth{
				Kubeconfig: pxcKubeconfigTest,
			},
			Name: "forth-pxc-test",
			Params: &controllerv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams{
				Resume: true,
			},
		}
		dbaasClient.On("UpdatePXCCluster", ctx, &mockReqResume).Return(&controllerv1beta1.UpdatePXCClusterResponse{}, nil)

		in = dbaasv1beta1.UpdatePXCClusterRequest{
			KubernetesClusterName: pxcKubernetesClusterNameTest,
			Name:                  "forth-pxc-test",
			Params: &dbaasv1beta1.UpdatePXCClusterRequest_UpdatePXCClusterParams{
				Resume: true,
			},
		}
		_, err = s.UpdatePXCCluster(ctx, &in)
		assert.NoError(t, err)
	})

	t.Run("BasicGetXtraDBClusterResources", func(t *testing.T) {
		t.Parallel()
		t.Run("ProxySQL", func(t *testing.T) {
			t.Parallel()
			s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
			v := int64(1000000000)
			r := int64(2000000000)

			in := dbaasv1beta1.GetPXCClusterResourcesRequest{
				Params: &dbaasv1beta1.PXCClusterParams{
					ClusterSize: 1,
					Pxc: &dbaasv1beta1.PXCClusterParams_PXC{
						ComputeResources: &dbaasv1beta1.ComputeResources{
							CpuM:        1000,
							MemoryBytes: v,
						},
						DiskSize: v,
					},
					Proxysql: &dbaasv1beta1.PXCClusterParams_ProxySQL{
						ComputeResources: &dbaasv1beta1.ComputeResources{
							CpuM:        1000,
							MemoryBytes: v,
						},
						DiskSize: v,
					},
				},
			}

			actual, err := s.GetPXCClusterResources(ctx, &in)
			assert.NoError(t, err)
			assert.Equal(t, uint64(r), actual.Expected.MemoryBytes)
			assert.Equal(t, uint64(2000), actual.Expected.CpuM)
			assert.Equal(t, uint64(r), actual.Expected.DiskSize)
		})

		t.Run("HAProxy", func(t *testing.T) {
			t.Parallel()
			s := NewPXCClusterService(db, dbaasClient, grafanaClient, componentsClient, versionService.GetVersionServiceURL())
			v := int64(1000000000)

			in := dbaasv1beta1.GetPXCClusterResourcesRequest{
				Params: &dbaasv1beta1.PXCClusterParams{
					ClusterSize: 1,
					Pxc: &dbaasv1beta1.PXCClusterParams_PXC{
						ComputeResources: &dbaasv1beta1.ComputeResources{
							CpuM:        1000,
							MemoryBytes: v,
						},
						DiskSize: v,
					},
					Haproxy: &dbaasv1beta1.PXCClusterParams_HAProxy{
						ComputeResources: &dbaasv1beta1.ComputeResources{
							CpuM:        1000,
							MemoryBytes: v,
						},
					},
				},
			}

			actual, err := s.GetPXCClusterResources(ctx, &in)
			assert.NoError(t, err)
			assert.Equal(t, uint64(2000000000), actual.Expected.MemoryBytes)
			assert.Equal(t, uint64(2000), actual.Expected.CpuM)
			assert.Equal(t, uint64(v), actual.Expected.DiskSize)
		})
	})
}
