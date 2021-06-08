// Code generated by mockery v1.0.0. DO NOT EDIT.

package management

import (
	context "context"

	check "github.com/percona-platform/saas/pkg/check"

	mock "github.com/stretchr/testify/mock"

	utils "github.com/percona/pmm-managed/utils"
)

// mockChecksService is an autogenerated mock type for the checksService type
type mockChecksService struct {
	mock.Mock
}

// ChangeInterval provides a mock function with given fields: params
func (_m *mockChecksService) ChangeInterval(params map[string]check.Interval) error {
	ret := _m.Called(params)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]check.Interval) error); ok {
		r0 = rf(params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DisableChecks provides a mock function with given fields: checkNames
func (_m *mockChecksService) DisableChecks(checkNames []string) error {
	ret := _m.Called(checkNames)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(checkNames)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnableChecks provides a mock function with given fields: checkNames
func (_m *mockChecksService) EnableChecks(checkNames []string) error {
	ret := _m.Called(checkNames)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(checkNames)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllChecks provides a mock function with given fields:
func (_m *mockChecksService) GetAllChecks() map[string]check.Check {
	ret := _m.Called()

	var r0 map[string]check.Check
	if rf, ok := ret.Get(0).(func() map[string]check.Check); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]check.Check)
		}
	}

	return r0
}

// GetDisabledChecks provides a mock function with given fields:
func (_m *mockChecksService) GetDisabledChecks() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecurityCheckResults provides a mock function with given fields:
func (_m *mockChecksService) GetSecurityCheckResults() ([]utils.STTCheckResult, error) {
	ret := _m.Called()

	var r0 []utils.STTCheckResult
	if rf, ok := ret.Get(0).(func() []utils.STTCheckResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]utils.STTCheckResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartChecks provides a mock function with given fields: ctx, group, checkNames
func (_m *mockChecksService) StartChecks(ctx context.Context, group check.Interval, checkNames []string) error {
	ret := _m.Called(ctx, group, checkNames)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, check.Interval, []string) error); ok {
		r0 = rf(ctx, group, checkNames)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
