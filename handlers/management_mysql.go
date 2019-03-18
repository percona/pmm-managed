package handlers

import (
	"context"

	"github.com/percona/pmm/api/managementpb"

	"github.com/percona/pmm-managed/services/management"
)

type mysqlGrpcServer struct {
	svc *management.MySQLService
}

// NewManagementMysqlServer creates Management MySQL Server.
func NewManagementMysqlServer(s *management.MySQLService) managementpb.MySQLServer {
	return &mysqlGrpcServer{svc: s}
}

// Add adds "MySQL Service", "MySQL Exporter Agent" and "QAN MySQL PerfSchema Agent".
func (s *mysqlGrpcServer) Add(ctx context.Context, req *managementpb.AddMySQLRequest) (*managementpb.AddMySQLResponse, error) {
	return s.svc.Add(ctx, req)
}
