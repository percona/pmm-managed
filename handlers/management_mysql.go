package handlers

import (
	"context"

	"github.com/percona/pmm-managed/services/management"
	"github.com/percona/pmm/api/managementpb"
)

type mysqlGrpcServer struct {
	svc *management.MySQLService
}

func NewManagementMysqlServer(s *management.MySQLService) managementpb.MySQLServer {
	return &mysqlGrpcServer{svc: s}
}

func (s *mysqlGrpcServer) Add(ctx context.Context, req *managementpb.AddMySQLRequest) (*managementpb.AddMySQLResponse, error) {
	return s.svc.Add(ctx, req)
}
