package platform

import (
	"context"

	"github.com/percona/pmm-managed/models"
)

// supervisordService is a subset of methods of supervisord.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type supervisordService interface {
	UpdateConfiguration(settings *models.Settings, ssoDetails *models.PerconaSSODetails) error
}

// type checksService is a subset of methods of checks.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type checksService interface {
	CollectChecks(ctx context.Context)
}

// grafanaClient is a subset of methods of grafana.Client used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type grafanaClient interface {
	GetCurrentUserAccessToken(ctx context.Context) (string, error)
}
