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
	"strings"
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

	tags := req.Tags
	if len(req.ServiceNames) == 0 && req.NodeName == "" {
		tags = append([]string{"pmm_annotation"}, tags...)
	}
	postfix := []string{}
	if len(req.ServiceNames) > 0 {
		for _, sn := range req.ServiceNames {
			_, err := models.FindServiceByName(as.db.Querier, sn)
			if err != nil {
				return nil, err
			}
		}

		tags = append(tags, req.ServiceNames...)
		postfix = append(postfix, "Service Name: "+strings.Join(req.ServiceNames, ","))
	}

	if req.NodeName != "" {
		_, err := models.FindNodeByName(as.db.Querier, req.NodeName)
		if err != nil {
			return nil, err
		}

		tags = append(tags, req.NodeName)
		postfix = append(postfix, "Node Name: "+req.NodeName)
	}

	if len(postfix) > 0 {
		req.Text += "(" + strings.Join(postfix, ",") + ")"
	}

	_, err := as.grafanaClient.CreateAnnotation(ctx, tags, time.Now(), req.Text, authorizationHeaders[0])
	if err != nil {
		return nil, err
	}

	return &managementpb.AddAnnotationResponse{}, nil
}
