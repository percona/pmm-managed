package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceEnabled interface {
	Enabled() bool
}

// UnaryServiceEnabledInterceptor returns a new unary server interceptor that checks if service is enabled.
//
// Request on disabled service will be rejected with `FailedPrecondition` before reaching any userspace handlers.
func UnaryServiceEnabledInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if svc, ok := info.Server.(serviceEnabled); ok {
			if !svc.Enabled() {
				return nil, status.Errorf(codes.FailedPrecondition, "Service %s is disabled.", extractServiceName(info.FullMethod))
			}
		}
		return handler(ctx, req)
	}
}

// StreamServiceEnabledInterceptor returns a new unary server interceptor that checks if service is enabled.
func StreamServiceEnabledInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if svc, ok := srv.(serviceEnabled); ok {
			if !svc.Enabled() {
				return status.Errorf(codes.FailedPrecondition, "Service %s is disabled.", extractServiceName(info.FullMethod))
			}
		}
		return handler(srv, stream)
	}
}

func extractServiceName(fullMethod string) string {
	split := strings.Split(fullMethod, "/")
	if len(split) < 2 {
		return fullMethod
	}
	return split[1]
}
