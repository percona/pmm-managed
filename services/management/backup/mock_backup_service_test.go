// Code generated by mockery v1.0.0. DO NOT EDIT.

package backup

import (
	context "context"

	models "github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/backup"
	mock "github.com/stretchr/testify/mock"
)

// mockBackupService is an autogenerated mock type for the backupService type
type mockBackupService struct {
	mock.Mock
}

// FindArtifactCompatibleServices provides a mock function with given fields: ctx, artifactID
func (_m *mockBackupService) FindArtifactCompatibleServices(ctx context.Context, artifactID string) ([]*models.Service, error) {
	ret := _m.Called(ctx, artifactID)

	var r0 []*models.Service
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.Service); ok {
		r0 = rf(ctx, artifactID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, artifactID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PerformBackup provides a mock function with given fields: ctx, params
func (_m *mockBackupService) PerformBackup(ctx context.Context, params backup.PerformBackupParams) (string, error) {
	ret := _m.Called(ctx, params)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, backup.PerformBackupParams) string); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, backup.PerformBackupParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RestoreBackup provides a mock function with given fields: ctx, serviceID, artifactID
func (_m *mockBackupService) RestoreBackup(ctx context.Context, serviceID string, artifactID string) (string, error) {
	ret := _m.Called(ctx, serviceID, artifactID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, serviceID, artifactID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, serviceID, artifactID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
