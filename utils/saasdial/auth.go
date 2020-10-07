package saasdial

import (
	"context"

	"google.golang.org/grpc/credentials"
)

type platformAuth struct {
	sessionID string
}

// GetRequestMetadata implements credentials.PerRPCCredentials interface.
func (b *platformAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": platformAuthType + " " + b.sessionID,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials interface.
func (*platformAuth) RequireTransportSecurity() bool {
	return false
}

// check interfaces
var (
	_ credentials.PerRPCCredentials = (*platformAuth)(nil)
)
