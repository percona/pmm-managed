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

package management

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/kr/pretty"
	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

const (
	// maximum time for AWS discover APIs calls
	awsDiscoverTimeout = 7 * time.Second
)

// RDSService RDS Management Service.
type RDSService struct {
	db       *reform.DB
	registry agentsRegistry
}

// InstanceID uniquely identifies RDS instance.
// http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.DBInstance.html
// Each DB instance has a DB instance identifier. This customer-supplied name uniquely identifies the DB instance when interacting
// with the Amazon RDS API and AWS CLI commands. The DB instance identifier must be unique for that customer in an AWS Region.
type InstanceID struct {
	Region string
	Name   string // DBInstanceIdentifier
}

type Instance struct {
	Node    models.Node
	Service models.Service
}

// NewRDSService creates new RDS Management Service.
func NewRDSService(db *reform.DB, registry agentsRegistry) *RDSService {
	return &RDSService{
		db:       db,
		registry: registry,
	}
}

func (svc *RDSService) Discover(ctx context.Context, accessKey, secretKey string) ([]Instance, error) {
	l := logger.Get(ctx).WithField("component", "rds")

	// do not break our API if some AWS region is slow or down
	ctx, cancel := context.WithTimeout(ctx, awsDiscoverTimeout)
	defer cancel()
	var g errgroup.Group
	instances := make(chan Instance)

	for _, r := range endpoints.AwsPartition().Services()[endpoints.RdsServiceID].Regions() {
		region := r.ID()
		g.Go(func() error {
			var creds *credentials.Credentials
			// use given credentials, or default credential chain
			if accessKey != "" || secretKey != "" {
				creds = credentials.NewStaticCredentials(accessKey, secretKey, "")
			}
			config := &aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
				Credentials:                   creds,
				Region:                        aws.String(region),
				HTTPClient:                    new(http.Client),
				Logger:                        aws.LoggerFunc(l.Debug),
				EnableEndpointDiscovery:       pointer.ToBool(true),
				DisableEndpointHostPrefix:     pointer.ToBool(true),
			}
			if l.Logger.GetLevel() >= logrus.DebugLevel {
				config.LogLevel = aws.LogLevel(aws.LogDebug)
			}
			s := session.Must(session.NewSession(config))
			err := rds.New(s).DescribeDBInstancesPagesWithContext(ctx, new(rds.DescribeDBInstancesInput),
				func(out *rds.DescribeDBInstancesOutput, lastPage bool) bool {
					for _, db := range out.DBInstances {
						instances <- Instance{
							Node: models.Node{
								Address:   pointer.GetString(db.Endpoint.Address),
								CreatedAt: *db.InstanceCreateTime,
								NodeID:    *db.DbiResourceId,
								NodeName:  *db.DBInstanceIdentifier,
								NodeType:  models.RDSNodeType,
								Region:    pointer.ToString(region),
							},

							Service: models.Service{
								Address:     db.Endpoint.Address,
								Cluster:     pointer.GetString(db.DBClusterIdentifier),
								CreatedAt:   *db.InstanceCreateTime,
								NodeID:      *db.DbiResourceId,
								Port:        pointer.ToUint16(uint16(*db.Endpoint.Port)),
								ServiceType: models.RDSServiceType,
							},
						}
					}
					return lastPage
				})
			pretty.Println(region)
			if err != nil {
				l.Error(errors.Wrap(err, region))

				if err, ok := err.(awserr.Error); ok {
					if err.OrigErr() != nil && err.OrigErr() == ctx.Err() {
						// ignore timeout, let other goroutines return partial data
						return nil
					}
					switch err.Code() {
					case "InvalidClientTokenId", "EmptyStaticCreds":
						return status.Error(codes.InvalidArgument, region+": "+err.Message())
					default:
						return err
					}
				}
				return errors.WithStack(err)
			}
			return nil
		})
	}

	go func() {
		g.Wait()
		close(instances)
	}()

	// sort by region and name
	var res []Instance
	for i := range instances {
		res = append(res, i)
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Node.Region != res[j].Node.Region {
			return pointer.GetString(res[i].Node.Region) < pointer.GetString(res[j].Node.Region)
		}
		return res[i].Node.NodeName < res[j].Node.NodeName
	})
	return res, g.Wait()
}
