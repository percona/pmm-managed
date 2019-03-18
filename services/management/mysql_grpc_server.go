package management

import (
	"context"

	"github.com/percona/pmm/api/managementpb"
)

type mysqlGrpcServer struct {
	svc *MySQLService
}

func NewMysqlGrpcServer(s *MySQLService) managementpb.MySQLServer {
	return &mysqlGrpcServer{svc: s}
}

func (s *mysqlGrpcServer) Add(ctx context.Context, req *managementpb.AddMySQLRequest) (*managementpb.AddMySQLResponse, error) {
	panic("not implemented")
}
