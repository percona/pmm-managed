// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package platform

import (
	"context"

	api "github.com/percona-platform/saas/gen/auth"
	"google.golang.org/grpc"
)

//go:generate mockery -name=saasService -case=snake -inpkg -testonly

// saasService represents a wrapper for Platform API calls
type saasService interface {
	SignUp(cc *grpc.ClientConn, ctx context.Context, req *api.SignUpRequest) (*api.SignUpResponse, error)
	SignIn(cc *grpc.ClientConn, ctx context.Context, req *api.SignInRequest) (*api.SignInResponse, error)
	SignOut(cc *grpc.ClientConn, ctx context.Context, req *api.SignOutRequest) (*api.SignOutResponse, error)
	RefreshSession(cc *grpc.ClientConn, ctx context.Context, req *api.RefreshSessionRequest) (*api.RefreshSessionResponse, error)
}
