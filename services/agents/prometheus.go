package agents

import (
	"context"
)

//go:generate mockery -name=prometheus -inpkg -testonly

// prometheus is a subset of methods of prometheus.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type prometheus interface {
	UpdateConfiguration(ctx context.Context) error
}
