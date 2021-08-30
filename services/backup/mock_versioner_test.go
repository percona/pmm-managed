// Code generated by mockery v1.0.0. DO NOT EDIT.

package backup

import (
	agents "github.com/percona/pmm-managed/services/agents"
	mock "github.com/stretchr/testify/mock"
)

// mockVersioner is an autogenerated mock type for the versioner type
type mockVersioner struct {
	mock.Mock
}

// GetVersions provides a mock function with given fields: pmmAgentID, softwares
func (_m *mockVersioner) GetVersions(pmmAgentID string, softwares []agents.Software) ([]agents.Version, error) {
	ret := _m.Called(pmmAgentID, softwares)

	var r0 []agents.Version
	if rf, ok := ret.Get(0).(func(string, []agents.Software) []agents.Version); ok {
		r0 = rf(pmmAgentID, softwares)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]agents.Version)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []agents.Software) error); ok {
		r1 = rf(pmmAgentID, softwares)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
