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

package inventory

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	inventorypb "github.com/percona/pmm/api/inventory"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// ServicesService works with inventory API Services.
type ServicesService struct {
	r registry
}

// NewServicesService creates new ServicesService
func NewServicesService(r registry) *ServicesService {
	return &ServicesService{
		r: r,
	}
}

// toInventoryService converts database row to Inventory API Service.
func toInventoryService(row *models.Service) (inventorypb.Service, error) {
	labels, err := row.GetCustomLabels()
	if err != nil {
		return nil, err
	}

	switch row.ServiceType {
	case models.MySQLServiceType:
		return &inventorypb.MySQLService{
			ServiceId:    row.ServiceID,
			ServiceName:  row.ServiceName,
			NodeId:       row.NodeID,
			Address:      pointer.GetString(row.Address),
			Port:         uint32(pointer.GetUint16(row.Port)),
			CustomLabels: labels,
		}, nil
	case models.MongoDBServiceType:
		return &inventorypb.MongoDBService{
			ServiceId:    row.ServiceID,
			ServiceName:  row.ServiceName,
			NodeId:       row.NodeID,
			Address:      pointer.GetString(row.Address),
			Port:         uint32(pointer.GetUint16(row.Port)),
			CustomLabels: labels,
		}, nil
	case models.PostgreSQLServiceType:
		return &inventorypb.PostgreSQLService{
			ServiceId:    row.ServiceID,
			ServiceName:  row.ServiceName,
			NodeId:       row.NodeID,
			Address:      pointer.GetString(row.Address),
			Port:         uint32(pointer.GetUint16(row.Port)),
			CustomLabels: labels,
		}, nil

	default:
		panic(fmt.Errorf("unhandled Service type %s", row.ServiceType))
	}
}

func toInventoryServices(services []*models.Service) ([]inventorypb.Service, error) {
	var err error
	res := make([]inventorypb.Service, len(services))
	for i, srv := range services {
		res[i], err = toInventoryService(srv)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// List selects all Services in a stable order.
//nolint:unparam
func (ss *ServicesService) List(ctx context.Context, q *reform.Querier) ([]inventorypb.Service, error) {
	services, err := models.FindAllServices(q)
	if err != nil {
		return nil, err
	}
	return toInventoryServices(services)
}

// Get selects a single Service by ID.
//nolint:unparam
func (ss *ServicesService) Get(ctx context.Context, q *reform.Querier, id string) (inventorypb.Service, error) {
	row, err := models.FindServiceByID(q, id)
	if err != nil {
		return nil, err
	}
	return toInventoryService(row)
}

// AddMySQL inserts MySQL Service with given parameters.
//nolint:dupl
func (ss *ServicesService) AddMySQL(ctx context.Context, q *reform.Querier, params *models.AddDBMSServiceParams) (*inventorypb.MySQLService, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// Both address and socket can't be empty, etc.

	row, err := models.AddNewService(q, models.MySQLServiceType, params)
	if err != nil {
		return nil, err
	}

	res, err := toInventoryService(row)
	if err != nil {
		return nil, err
	}
	return res.(*inventorypb.MySQLService), nil
}

// AddMongoDB inserts MongoDB Service with given parameters.
//nolint:dupl
func (ss *ServicesService) AddMongoDB(ctx context.Context, q *reform.Querier, params *models.AddDBMSServiceParams) (*inventorypb.MongoDBService, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416

	row, err := models.AddNewService(q, models.MongoDBServiceType, params)
	if err != nil {
		return nil, err
	}

	res, err := toInventoryService(row)
	if err != nil {
		return nil, err
	}
	return res.(*inventorypb.MongoDBService), nil
}

// AddPostgreSQL inserts PostgreSQL Service with given parameters.
func (ss *ServicesService) AddPostgreSQL(ctx context.Context, q *reform.Querier, params *models.AddDBMSServiceParams) (*inventorypb.PostgreSQLService, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// Both address and socket can't be empty, etc.

	row, err := models.AddNewService(q, models.PostgreSQLServiceType, params)
	if err != nil {
		return nil, err
	}
	res, err := toInventoryService(row)
	if err != nil {
		return nil, err
	}
	return res.(*inventorypb.PostgreSQLService), nil
}

// Remove deletes Service by ID.
//nolint:unparam
func (ss *ServicesService) Remove(ctx context.Context, q *reform.Querier, id string) error {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0.

	// TODO check absence of Agents

	err := models.RemoveService(q, id)
	if err != nil {
		return err
	}
	return nil
}
