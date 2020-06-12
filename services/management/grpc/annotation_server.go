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

// Package grpc provides gRPC servers.
package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/grafana"
)

// AnnotationServer is a server for making annotations in Grafana.
type AnnotationServer struct {
	db            *reform.DB
	grafanaClient *grafana.Client
}

// NewAnnotationServer creates Annotation Server.
func NewAnnotationServer(db *reform.DB, grafanaClient *grafana.Client) *AnnotationServer {
	return &AnnotationServer{
		db:            db,
		grafanaClient: grafanaClient,
	}
}

// AddAnnotation adds annotation to Grafana.
func (as *AnnotationServer) AddAnnotation(ctx context.Context, req *managementpb.AddAnnotationRequest) (*managementpb.AddAnnotationResponse, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("cannot get headers from metadata %v", headers)
	}
	// get authorization from headers.
	authorizationHeaders := headers.Get("Authorization")
	if len(authorizationHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Authorization error.")
	}

	var service, node bool
	var serviceErr, nodeErr error
	if len(req.ServiceName) > 0 {
		service = true
	}
	if req.NodeName != "" {
		node = true
	}

	if service {
		serviceErr = findService(as.db, req.ServiceName)
	}

	if node {
		nodeErr = findNode(as.db, req.NodeName)
	}

	if serviceErr != nil && nodeErr != nil || serviceErr != nil {
		return nil, serviceErr
	}

	if nodeErr != nil {
		return nil, nodeErr
	}

	if service {
		for _, sn := range req.ServiceName {
			text := fmt.Sprintf("%s (Service Name: %s)", req.Text, sn)
			_, err := as.grafanaClient.CreateAnnotation(ctx, req.Tags, time.Now(), text, authorizationHeaders[0])
			if err != nil {
				return nil, err
			}
		}
	}

	if node {
		text := fmt.Sprintf("%s (Node Name: %s)", req.Text, req.NodeName)
		_, err := as.grafanaClient.CreateAnnotation(ctx, req.Tags, time.Now(), text, authorizationHeaders[0])
		if err != nil {
			return nil, err
		}
	}

	if !node && !service {
		_, err := as.grafanaClient.CreateAnnotation(ctx, req.Tags, time.Now(), req.Text, authorizationHeaders[0])
		if err != nil {
			return nil, err
		}
	}
	return &managementpb.AddAnnotationResponse{}, nil
}

func findService(db *reform.DB, serviceName []string) error {
	for _, sn := range serviceName {
		err := db.InTransaction(func(tx *reform.TX) error {
			var err error
			_, err = models.FindServiceByName(tx.Querier, sn)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func findNode(db *reform.DB, nodeName string) error {
	err := db.InTransaction(func(tx *reform.TX) error {
		var err error
		_, err = models.FindNodeByName(tx.Querier, nodeName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
