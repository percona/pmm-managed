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
	"encoding/json"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/api/managementpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/logger"

	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2019-04-01/resourcegraph"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	// https://docs.microsoft.com/en-us/azure/governance/resource-graph/concepts/query-language
	// https://docs.microsoft.com/en-us/azure/azure-monitor/essentials/metrics-supported
	// TODO: add pagination and filtering https://jira.percona.com/browse/PMM-7813
	azureDatabaseResourceQuery string = `
		Resources
			| where type in~ (
				'Microsoft.DBforMySQL/servers',
				'Microsoft.DBforMySQL/flexibleServers',
				'Microsoft.DBforMariaDB/servers',
				'Microsoft.DBforPostgreSQL/servers',
				'Microsoft.DBforPostgreSQL/serversv2',
				'Microsoft.DBforPostgreSQL/flexibleServers'
			)
			| order by name asc
			| limit 1000
	`
)

// AzureDatabaseService represents instance discovery service.
type AzureDatabaseService struct {
	db       *reform.DB
	registry agentsRegistry
}

// NewAzureDatabaseService creates new instance discovery service.
func NewAzureDatabaseService(db *reform.DB, registry agentsRegistry) *AzureDatabaseService {
	return &AzureDatabaseService{
		db:       db,
		registry: registry,
	}
}

// AzureDatabaseInstanceData reflects Azure Database Instance Data of Discovery Response.
type AzureDatabaseInstanceData struct {
	ID            string                 `json:"id"`
	Location      string                 `json:"location"`
	Name          string                 `json:"name"`
	Properties    map[string]interface{} `json:"properties"`
	Tags          map[string]string      `json:"tags"`
	Sku           map[string]interface{} `json:"sku"`
	ResourceGroup string                 `json:"resourceGroup"`
	Type          string                 `json:"type"`
	Zones         string                 `json:"zones"`
}

func (s *AzureDatabaseService) getAzureClient(req *managementpb.DiscoverAzureDatabaseRequest) (*resourcegraph.BaseClient, error) {
	authSettings := auth.EnvironmentSettings{
		Values: map[string]string{
			auth.ClientID:       req.AzureClientId,
			auth.ClientSecret:   req.AzureClientSecret,
			auth.SubscriptionID: req.AzureSubscriptionId,
			auth.TenantID:       req.AzureTenantId,
			auth.Resource:       azure.PublicCloud.ResourceManagerEndpoint,
		},
		Environment: azure.PublicCloud,
	}

	// Create and authorize a ResourceGraph client
	client := resourcegraph.New()
	authorizer, err := authSettings.GetAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = authorizer
	return &client, nil
}

func (s *AzureDatabaseService) fetchAzureDatabaseInstancesData(
	ctx context.Context,
	req *managementpb.DiscoverAzureDatabaseRequest,
	client *resourcegraph.BaseClient,
) ([]AzureDatabaseInstanceData, error) {
	query := azureDatabaseResourceQuery
	request := resourcegraph.QueryRequest{
		Subscriptions: &[]string{req.AzureSubscriptionId},
		Query:         &query,
		Options: &resourcegraph.QueryRequestOptions{
			ResultFormat: "objectArray",
		},
	}

	// Run the query and get the results
	results, err := client.Resources(ctx, request)
	if err != nil {
		return nil, err
	}

	d, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	dataInst := struct {
		Data []AzureDatabaseInstanceData `json:"data"`
	}{}

	err = json.Unmarshal(d, &dataInst)
	if err != nil {
		return nil, err
	}

	return dataInst.Data, nil
}

// DiscoverAzureDatabase discovers database instances on Azure.
func (s *AzureDatabaseService) DiscoverAzureDatabase(
	ctx context.Context,
	req *managementpb.DiscoverAzureDatabaseRequest,
) (*managementpb.DiscoverAzureDatabaseResponse, error) {
	client, err := s.getAzureClient(req)
	if err != nil {
		return nil, err
	}

	dataInstData, err := s.fetchAzureDatabaseInstancesData(ctx, req, client)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	resp := managementpb.DiscoverAzureDatabaseResponse{}

	for _, instance := range dataInstData {
		inst := managementpb.DiscoverAzureDatabaseInstance{
			InstanceId:         instance.ID,
			Region:             instance.Location,
			ServiceName:        instance.Name,
			AzureResourceGroup: instance.ResourceGroup,
			Environment:        instance.Tags["environment"],
			Az:                 instance.Zones,
		}
		switch instance.Type {
		case "microsoft.dbformysql/servers",
			"microsoft.dbformysql/flexibleservers":
			inst.Type = managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_MYSQL
		case "microsoft.dbforpostgresql/servers",
			"microsoft.dbforpostgresql/flexibleservers",
			"microsoft.dbforpostgresql/serversv2":
			inst.Type = managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_POSTGRESQL
		case "microsoft.dbformariadb/servers":
			inst.Type = managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_MARIADB
		default:
			inst.Type = managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_INVALID
		}

		if val, ok := instance.Properties["administratorLogin"].(string); ok {
			inst.Username = fmt.Sprintf("%s@%s", val, instance.Name)
		}
		if val, ok := instance.Properties["fullyQualifiedDomainName"].(string); ok {
			inst.Address = val
		}
		if val, ok := instance.Sku["name"].(string); ok {
			inst.NodeModel = val
		}

		resp.AzureDatabaseInstance = append(resp.AzureDatabaseInstance, &inst)
	}

	return &resp, nil
}

// AddAzureDatabase add azure database to monitoring.
//nolint:funlen
func (s *AzureDatabaseService) AddAzureDatabase(ctx context.Context, req *managementpb.AddAzureDatabaseRequest) (*managementpb.AddAzureDatabaseResponse, error) {
	l := logger.Get(ctx).WithField("component", "discover/azureDatabase")
	if e := s.db.InTransaction(func(tx *reform.TX) error {
		// tweak according to API docs
		if req.NodeName == "" {
			req.NodeName = req.InstanceId
		}
		if req.ServiceName == "" {
			req.ServiceName = req.InstanceId
		}

		// tweak according to API docs
		tablestatsGroupTableLimit := req.TablestatsGroupTableLimit
		if tablestatsGroupTableLimit == 0 {
			tablestatsGroupTableLimit = defaultTablestatsGroupTableLimit
		}
		if tablestatsGroupTableLimit < 0 {
			tablestatsGroupTableLimit = -1
		}

		// add Remote Azure Database Node
		node, err := models.CreateNode(tx.Querier, models.RemoteAzureDatabaseNodeType, &models.CreateNodeParams{
			NodeName:     req.NodeName,
			NodeModel:    req.NodeModel,
			AZ:           req.Az,
			Address:      req.InstanceId,
			Region:       &req.Region,
			CustomLabels: req.CustomLabels,
		})
		if err != nil {
			return err
		}
		l.Infof("Created Azure Database Node with NodeID: %s", node.NodeID)

		// add Azure Database Agent
		if req.AzureDatabaseExporter {

			creds := models.AzureCredentials{
				SubscriptionID: req.AzureSubscriptionId,
				ClientID:       req.AzureClientId,
				ClientSecret:   req.AzureClientSecret,
				TenantID:       req.AzureTenantId,
				ResourceGroup:  req.AzureResourceGroup,
			}

			azureCredentials, err := json.Marshal(creds)
			if err != nil {
				return err
			}

			//nolint:exhaustive
			switch req.Type {
			case managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_MYSQL,
				managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_MARIADB:

				// add MySQL Service
				service, err := models.AddNewService(tx.Querier, models.MySQLServiceType, &models.AddDBMSServiceParams{
					ServiceName:  req.ServiceName,
					NodeID:       node.NodeID,
					Environment:  req.Environment,
					CustomLabels: req.CustomLabels,
					Address:      &req.Address,
					Port:         pointer.ToUint16(uint16(req.Port)),
				})
				if err != nil {
					return err
				}
				l.Infof("Added Azure Database Service with ServiceID: %s", service.ServiceID)

				azureDatabaseExporter, err := models.CreateAgent(tx.Querier, models.AzureDatabaseExporterType, &models.CreateAgentParams{
					PMMAgentID:       models.PMMServerAgentID,
					ServiceID:        service.ServiceID,
					AzureCredentials: string(azureCredentials),
				})
				if err != nil {
					return err
				}
				l.Infof("Created Azure Database Exporter with AgentID: %s", azureDatabaseExporter.AgentID)

				// add MySQL Exporter
				mysqldExporter, err := models.CreateAgent(tx.Querier, models.MySQLdExporterType, &models.CreateAgentParams{
					PMMAgentID:                     models.PMMServerAgentID,
					ServiceID:                      service.ServiceID,
					Username:                       req.Username,
					Password:                       req.Password,
					TLS:                            req.Tls,
					TLSSkipVerify:                  req.TlsSkipVerify,
					TableCountTablestatsGroupLimit: tablestatsGroupTableLimit,
				})
				if err != nil {
					return err
				}
				l.Infof("Added Azure Database Exporter with AgentID: %s", mysqldExporter.AgentID)

				if !req.SkipConnectionCheck {
					if err = s.registry.CheckConnectionToService(ctx, tx.Querier, service, mysqldExporter); err != nil {
						return err
					}
				}

				// add MySQL PerfSchema QAN Agent
				if req.Qan {
					qanAgent, err := models.CreateAgent(tx.Querier, models.QANMySQLPerfSchemaAgentType, &models.CreateAgentParams{
						PMMAgentID:            models.PMMServerAgentID,
						ServiceID:             service.ServiceID,
						Username:              req.Username,
						Password:              req.Password,
						TLS:                   req.Tls,
						TLSSkipVerify:         req.TlsSkipVerify,
						QueryExamplesDisabled: req.DisableQueryExamples,
					})
					if err != nil {
						return err
					}
					l.Infof("Added Azure Database QAN with AgentID: %s", qanAgent.AgentID)
				}

				return nil

			case managementpb.DiscoverAzureDatabaseType_DISCOVER_AZURE_DATABASE_TYPE_POSTGRESQL:
				// add PostgreSQL Service
				service, err := models.AddNewService(tx.Querier, models.PostgreSQLServiceType, &models.AddDBMSServiceParams{
					ServiceName:  req.ServiceName,
					NodeID:       node.NodeID,
					Environment:  req.Environment,
					CustomLabels: req.CustomLabels,
					Address:      &req.Address,
					Port:         pointer.ToUint16(uint16(req.Port)),
				})
				if err != nil {
					return err
				}
				l.Infof("Added Azure Database Service with ServiceID: %s", service.ServiceID)

				azureDatabaseExporter, err := models.CreateAgent(tx.Querier, models.AzureDatabaseExporterType, &models.CreateAgentParams{
					PMMAgentID:       models.PMMServerAgentID,
					ServiceID:        service.ServiceID,
					AzureCredentials: string(azureCredentials),
				})
				if err != nil {
					return err
				}
				l.Infof("Created Azure Database Exporter with AgentID: %s", azureDatabaseExporter.AgentID)

				// add PostgreSQL Exporter
				postgresqlExporter, err := models.CreateAgent(tx.Querier, models.PostgresExporterType, &models.CreateAgentParams{
					PMMAgentID:    models.PMMServerAgentID,
					ServiceID:     service.ServiceID,
					Username:      req.Username,
					Password:      req.Password,
					TLS:           req.Tls,
					TLSSkipVerify: req.TlsSkipVerify,
				})
				if err != nil {
					return err
				}
				l.Infof("Added Azure Database Exporter with AgentID: %s", postgresqlExporter.AgentID)

				if !req.SkipConnectionCheck {
					if err = s.registry.CheckConnectionToService(ctx, tx.Querier, service, postgresqlExporter); err != nil {
						return err
					}
				}

				// add MySQL PerfSchema QAN Agent
				if req.Qan {
					qanAgent, err := models.CreateAgent(tx.Querier, models.QANPostgreSQLPgStatementsAgentType, &models.CreateAgentParams{
						PMMAgentID:            models.PMMServerAgentID,
						ServiceID:             service.ServiceID,
						Username:              req.Username,
						Password:              req.Password,
						TLS:                   req.Tls,
						TLSSkipVerify:         req.TlsSkipVerify,
						QueryExamplesDisabled: req.DisableQueryExamples,
					})
					if err != nil {
						return err
					}
					l.Infof("Added Azure Database QAN with AgentID: %s", qanAgent.AgentID)
				}

			default:
				return status.Errorf(codes.InvalidArgument, "Unsupported Azure Database type %q.", req.Type)
			}
		}

		return nil
	}); e != nil {
		return nil, e
	}
	return &managementpb.AddAzureDatabaseResponse{}, nil
}
