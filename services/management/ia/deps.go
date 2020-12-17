package ia

import (
	"context"

	"github.com/percona/pmm/api/alertmanager/ammodels"
)

type alertManager interface {
	GetAlerts(ctx context.Context) ([]*ammodels.GettableAlert, error)
	Silence(ctx context.Context, id string) error
	Unsilence(ctx context.Context, id string) error
}
