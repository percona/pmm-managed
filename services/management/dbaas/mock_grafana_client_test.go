// Code generated by mockery v1.0.0. DO NOT EDIT.

package dbaas

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockGrafanaClient is an autogenerated mock type for the grafanaClient type
type mockGrafanaClient struct {
	mock.Mock
}

// CreateAdminAPIKey provides a mock function with given fields: ctx, name
func (_m *mockGrafanaClient) CreateAdminAPIKey(ctx context.Context, name string) (int64, string, error) {
	ret := _m.Called(ctx, name)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, name)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteAPIKeyByID provides a mock function with given fields: ctx, id
func (_m *mockGrafanaClient) DeleteAPIKeyByID(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAPIKeysWithPrefix provides a mock function with given fields: ctx, name
func (_m *mockGrafanaClient) DeleteAPIKeysWithPrefix(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
