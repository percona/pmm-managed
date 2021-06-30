package scheduler

import "context"

//go:generate mockery -name=backupService -case=snake -inpkg -testonly

type backupService interface {
	PerformBackup(ctx context.Context, serviceID, locationID, name, scheduleID string) (string, error)
}
