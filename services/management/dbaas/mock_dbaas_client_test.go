// Code generated by mockery v1.0.0. DO NOT EDIT.

package dbaas

import (
	context "context"

	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// mockDbaasClient is an autogenerated mock type for the dbaasClient type
type mockDbaasClient struct {
	mock.Mock
}

// CheckKubernetesClusterConnection provides a mock function with given fields: ctx, kubeConfig
func (_m *mockDbaasClient) CheckKubernetesClusterConnection(ctx context.Context, kubeConfig string) (*controllerv1beta1.CheckKubernetesClusterConnectionResponse, error) {
	ret := _m.Called(ctx, kubeConfig)

	var r0 *controllerv1beta1.CheckKubernetesClusterConnectionResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) *controllerv1beta1.CheckKubernetesClusterConnectionResponse); ok {
		r0 = rf(ctx, kubeConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.CheckKubernetesClusterConnectionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, kubeConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePSMDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) CreatePSMDBCluster(ctx context.Context, in *controllerv1beta1.CreatePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.CreatePSMDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.CreatePSMDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.CreatePSMDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.CreatePSMDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.CreatePSMDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.CreatePSMDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateXtraDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) CreateXtraDBCluster(ctx context.Context, in *controllerv1beta1.CreateXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.CreateXtraDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.CreateXtraDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.CreateXtraDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.CreateXtraDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.CreateXtraDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.CreateXtraDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePSMDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) DeletePSMDBCluster(ctx context.Context, in *controllerv1beta1.DeletePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.DeletePSMDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.DeletePSMDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.DeletePSMDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.DeletePSMDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.DeletePSMDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.DeletePSMDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteXtraDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) DeleteXtraDBCluster(ctx context.Context, in *controllerv1beta1.DeleteXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.DeleteXtraDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.DeleteXtraDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.DeleteXtraDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.DeleteXtraDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.DeleteXtraDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.DeleteXtraDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPSMDBClusterCredentials provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) GetPSMDBClusterCredentials(ctx context.Context, in *controllerv1beta1.GetPSMDBClusterCredentialsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetPSMDBClusterCredentialsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.GetPSMDBClusterCredentialsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.GetPSMDBClusterCredentialsRequest, ...grpc.CallOption) *controllerv1beta1.GetPSMDBClusterCredentialsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.GetPSMDBClusterCredentialsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.GetPSMDBClusterCredentialsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetXtraDBClusterCredentials provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) GetXtraDBClusterCredentials(ctx context.Context, in *controllerv1beta1.GetXtraDBClusterCredentialsRequest, opts ...grpc.CallOption) (*controllerv1beta1.GetXtraDBClusterCredentialsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.GetXtraDBClusterCredentialsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.GetXtraDBClusterCredentialsRequest, ...grpc.CallOption) *controllerv1beta1.GetXtraDBClusterCredentialsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.GetXtraDBClusterCredentialsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.GetXtraDBClusterCredentialsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPSMDBClusters provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) ListPSMDBClusters(ctx context.Context, in *controllerv1beta1.ListPSMDBClustersRequest, opts ...grpc.CallOption) (*controllerv1beta1.ListPSMDBClustersResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.ListPSMDBClustersResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.ListPSMDBClustersRequest, ...grpc.CallOption) *controllerv1beta1.ListPSMDBClustersResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.ListPSMDBClustersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.ListPSMDBClustersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListXtraDBClusters provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) ListXtraDBClusters(ctx context.Context, in *controllerv1beta1.ListXtraDBClustersRequest, opts ...grpc.CallOption) (*controllerv1beta1.ListXtraDBClustersResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.ListXtraDBClustersResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.ListXtraDBClustersRequest, ...grpc.CallOption) *controllerv1beta1.ListXtraDBClustersResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.ListXtraDBClustersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.ListXtraDBClustersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RestartPSMDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) RestartPSMDBCluster(ctx context.Context, in *controllerv1beta1.RestartPSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.RestartPSMDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.RestartPSMDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.RestartPSMDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.RestartPSMDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.RestartPSMDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.RestartPSMDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RestartXtraDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) RestartXtraDBCluster(ctx context.Context, in *controllerv1beta1.RestartXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.RestartXtraDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.RestartXtraDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.RestartXtraDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.RestartXtraDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.RestartXtraDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.RestartXtraDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePSMDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) UpdatePSMDBCluster(ctx context.Context, in *controllerv1beta1.UpdatePSMDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.UpdatePSMDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.UpdatePSMDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.UpdatePSMDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.UpdatePSMDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.UpdatePSMDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.UpdatePSMDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateXtraDBCluster provides a mock function with given fields: ctx, in, opts
func (_m *mockDbaasClient) UpdateXtraDBCluster(ctx context.Context, in *controllerv1beta1.UpdateXtraDBClusterRequest, opts ...grpc.CallOption) (*controllerv1beta1.UpdateXtraDBClusterResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *controllerv1beta1.UpdateXtraDBClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, *controllerv1beta1.UpdateXtraDBClusterRequest, ...grpc.CallOption) *controllerv1beta1.UpdateXtraDBClusterResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*controllerv1beta1.UpdateXtraDBClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *controllerv1beta1.UpdateXtraDBClusterRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
