package grafana

import (
	"context"
	"fmt"

	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/metadata"
)

type AnnotationService struct {
	grafanaClient *Client
}

func NewAnnotationService(grafanaClient *Client) *AnnotationService {
	return &AnnotationService{
		grafanaClient: grafanaClient,
	}
}

func (as *AnnotationService) AddAnnotation(ctx context.Context, req *managementpb.AddAnnotationRequest) (*managementpb.AddAnnotationResponse, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("cannot get auth token from headers %v", headers)
	}
	message, err := as.grafanaClient.CreateAnnotation(ctx, req.Tags, req.Text, headers["authorization"][0])
	if err != nil {
		return nil, err
	}
	return &managementpb.AddAnnotationResponse{Message: message}, nil
}
