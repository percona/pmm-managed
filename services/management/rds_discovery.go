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
	"github.com/percona/pmm/api/managementpb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/percona/pmm-managed/utils/logger"
)

const (
	// maximum time for AWS discover APIs calls
	awsDiscoverTimeout = 7 * time.Second
)

// RDSService RDS Management Service.
type RDSService struct {
}

// NewRDSService creates new RDS Management Service.
func NewRDSService() *RDSService {
	return &RDSService{}
}

// Discover returns a list of RDS instances from all AWS zones
func (svc *RDSService) Discover(ctx context.Context, accessKey, secretKey string) (*managementpb.DiscoverRDSResponse, error) {
	l := logger.Get(ctx).WithField("component", "rds")

	// do not break our API if some AWS region is slow or down
	ctx, cancel := context.WithTimeout(ctx, awsDiscoverTimeout)
	defer cancel()
	var g errgroup.Group
	instances := make(chan *managementpb.RDSDiscoveryInstance)

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
						instances <- &managementpb.RDSDiscoveryInstance{
							Address:       pointer.GetString(db.Endpoint.Address),
							Engine:        pointer.GetString(db.Engine),
							EngineVersion: pointer.GetString(db.EngineVersion),
							InstanceId:    *db.DbiResourceId,
							Region:        region,
						}
					}
					return lastPage
				})

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
				// This function is running in a go-routine. We should concatenate errors from different zones.
				return errors.WithStack(err)
			}
			return nil
		})
	}

	go func() {
		g.Wait()
		close(instances)
	}()

	// sort by region and id
	res := &managementpb.DiscoverRDSResponse{
		RdsInstances: make([]*managementpb.RDSDiscoveryInstance, 0),
	}

	for i := range instances {
		res.RdsInstances = append(res.RdsInstances, i)
	}
	sort.Slice(res.RdsInstances, func(i, j int) bool {
		if res.RdsInstances[i].Region != res.RdsInstances[j].Region {
			return res.RdsInstances[i].Region < res.RdsInstances[j].Region
		}
		return res.RdsInstances[i].InstanceId < res.RdsInstances[j].InstanceId
	})
	return res, g.Wait()
}
