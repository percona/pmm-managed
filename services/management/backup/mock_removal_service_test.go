// Code generated by mockery v1.0.0. DO NOT EDIT.

package backup

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockRemovalService is an autogenerated mock type for the removalService type
type mockRemovalService struct {
	mock.Mock
}

// DeleteArtifact provides a mock function with given fields: ctx, artifactID
func (_m *mockRemovalService) DeleteArtifact(ctx context.Context, artifactID string) error {
	ret := _m.Called(ctx, artifactID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, artifactID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
