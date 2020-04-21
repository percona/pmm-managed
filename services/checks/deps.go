package checks

import "context"

//go:generate mockery -name=registryService -case=snake -inpkg -testonly

// registryService is a subset of methods of agents.Registry used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type registryService interface {
	StartMySQLQueryShowAction(ctx context.Context, id, pmmAgentID, dsn, query string) error
	StartMySQLQuerySelectAction(ctx context.Context, id, pmmAgentID, dsn, query string) error
	StartPostgreSQLQueryShowAction(ctx context.Context, id, pmmAgentID, dsn string) error
	StartPostgreSQLQuerySelectAction(ctx context.Context, id, pmmAgentID, dsn, query string) error
	StartMongoDBQueryGetParameterAction(ctx context.Context, id, pmmAgentID, dsn string) error
}
