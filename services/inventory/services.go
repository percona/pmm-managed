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
	"github.com/percona/pmm/api/inventory"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// ServicesService works with inventory API Services.
type ServicesService struct {
	q *reform.Querier
	// r *agents.Registry
}

func NewServicesService(q *reform.Querier) *ServicesService {
	return &ServicesService{
		q: q,
		// r: r,
	}
}

// makeService converts database row to Inventory API Service.
func makeService(row *models.ServiceRow) inventory.Service {
	switch row.Type {
	case models.MySQLServiceType:
		return &inventory.MySQLService{
			Id:         row.ID,
			Name:       row.Name,
			NodeId:     row.NodeID,
			Address:    pointer.GetString(row.Address),
			Port:       uint32(pointer.GetUint16(row.Port)),
			UnixSocket: pointer.GetString(row.UnixSocket),
		}

	default:
		panic(fmt.Errorf("unhandled ServiceRow type %s", row.Type))
	}
}

func (ss *ServicesService) get(ctx context.Context, id string) (*models.ServiceRow, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Service ID.")
	}

	row := &models.ServiceRow{ID: id}
	switch err := ss.q.Reload(row); err {
	case nil:
		return row, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Service with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

func (ss *ServicesService) checkUniqueName(ctx context.Context, name string) error {
	_, err := ss.q.FindOneFrom(models.ServiceRowTable, "name", name)
	switch err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Service with name %q already exists.", name)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// List selects all Services in a stable order.
func (ss *ServicesService) List(ctx context.Context) ([]inventory.Service, error) {
	structs, err := ss.q.SelectAllFrom(models.ServiceRowTable, "ORDER BY id")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]inventory.Service, len(structs))
	for i, str := range structs {
		row := str.(*models.ServiceRow)
		res[i] = makeService(row)
	}
	return res, nil
}

// Get selects a single Service by ID.
func (ss *ServicesService) Get(ctx context.Context, id string) (inventory.Service, error) {
	row, err := ss.get(ctx, id)
	if err != nil {
		return nil, err
	}
	return makeService(row), nil
}

// AddMySQL inserts MySQL Service with given parameters.
func (ss *ServicesService) AddMySQL(ctx context.Context, name string, nodeID string, address *string, port *uint16, unixSocket *string) (inventory.Service, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// Both address and socket can't be empty, etc.

	if err := ss.checkUniqueName(ctx, name); err != nil {
		return nil, err
	}

	ns := NewNodesService(ss.q)
	if _, err := ns.get(ctx, nodeID); err != nil {
		return nil, err
	}

	row := &models.ServiceRow{
		ID:         makeID(),
		Type:       models.MySQLServiceType,
		Name:       name,
		NodeID:     nodeID,
		Address:    address,
		Port:       port,
		UnixSocket: unixSocket,
	}
	if err := ss.q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}
	return makeService(row), nil
}

// Change updates Service by ID.
func (ss *ServicesService) Change(ctx context.Context, id string, name string) (inventory.Service, error) {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0, name is not empty and valid.

	if err := ss.checkUniqueName(ctx, name); err != nil {
		return nil, err
	}

	row, err := ss.get(ctx, id)
	if err != nil {
		return nil, err
	}

	row.Name = name
	if err = ss.q.Update(row); err != nil {
		return nil, errors.WithStack(err)
	}
	return makeService(row), nil
}

// Remove deletes Service by ID.
func (ss *ServicesService) Remove(ctx context.Context, id string) error {
	// TODO Decide about validation. https://jira.percona.com/browse/PMM-1416
	// ID is not 0.

	// TODO check absence of Agents

	err := ss.q.Delete(&models.ServiceRow{ID: id})
	if err == reform.ErrNoRows {
		return status.Errorf(codes.NotFound, "Service with ID %q not found.", id)
	}
	return errors.WithStack(err)
}
