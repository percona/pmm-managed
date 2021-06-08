package scheduler

import "context"

//go:generate mockery -name=backupsLogicService -case=snake -inpkg -testonly

type backupsLogicService interface {
	PerformBackup(ctx context.Context, serviceID, locationID, name string) (string, error)
}
