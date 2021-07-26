// Code generated by mockery v1.0.0. DO NOT EDIT.

package management

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockAgentsRegistry is an autogenerated mock type for the agentsRegistry type
type mockAgentsRegistry struct {
	mock.Mock
}

// IsConnected provides a mock function with given fields: pmmAgentID
func (_m *mockAgentsRegistry) IsConnected(pmmAgentID string) bool {
	ret := _m.Called(pmmAgentID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(pmmAgentID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Kick provides a mock function with given fields: ctx, pmmAgentID
func (_m *mockAgentsRegistry) Kick(ctx context.Context, pmmAgentID string) {
	_m.Called(ctx, pmmAgentID)
}

// RequestStateUpdate provides a mock function with given fields: ctx, pmmAgentID
func (_m *mockAgentsRegistry) RequestStateUpdate(ctx context.Context, pmmAgentID string) {
	_m.Called(ctx, pmmAgentID)
}
