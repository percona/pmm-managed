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

package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	inventorypb "github.com/percona/pmm/api/inventory"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// FindServiceByID finds service by ID.
func FindServiceByID(q *reform.Querier, id string) (*Service, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Service ID.")
	}

	row := &Service{ServiceID: id}
	switch err := q.Reload(row); err {
	case nil:
		return row, nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Service with ID %q not found.", id)
	default:
		return nil, errors.WithStack(err)
	}
}

func checkServiceUniqueID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Service ID")
	}

	row := &Service{ServiceID: id}
	switch err := q.Reload(row); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Service with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

func checkServiceUniqueName(q *reform.Querier, name string) error {
	_, err := q.FindOneFrom(ServiceTable, "service_name", name)
	switch err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Service with name %q already exists.", name)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// AddDBMSServiceParams contains parameters for adding DBMS (MySQL, PostgreSQL, MongoDB) Services.
type AddDBMSServiceParams struct {
	ServiceName  string
	NodeID       string
	CustomLabels map[string]string
	Address      *string
	Port         *uint16
}

// AddNewService adds new service to storage.
func AddNewService(q *reform.Querier, serviceType ServiceType, params *AddDBMSServiceParams) (*Service, error) {
	id := "/service_id/" + uuid.New().String()
	if err := checkServiceUniqueID(q, id); err != nil {
		return nil, err
	}
	if err := checkServiceUniqueName(q, params.ServiceName); err != nil {
		return nil, err
	}

	if _, err := FindNodeByID(q, params.NodeID); err != nil {
		return nil, err
	}

	row := &Service{
		ServiceID:   id,
		ServiceType: serviceType,
		ServiceName: params.ServiceName,
		NodeID:      params.NodeID,
		Address:     params.Address,
		Port:        params.Port,
	}
	if err := row.SetCustomLabels(params.CustomLabels); err != nil {
		return nil, err
	}
	if err := q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	return row, nil
}

// FindAllServices finds all nodes and loads it from persistent store.
func FindAllServices(q *reform.Querier) ([]*Service, error) {
	structs, err := q.SelectAllFrom(ServiceTable, "ORDER BY service_id")
	if err != nil {
		return nil, err
	}

	services := make([]*Service, len(structs))
	for i, s := range structs {
		services[i] = s.(*Service)
	}

	return services, nil
}

// RemoveService removes single node prom persistent store.
func RemoveService(q *reform.Querier, id string) error {
	err := q.Delete(&Service{ServiceID: id})
	if err == reform.ErrNoRows {
		return status.Errorf(codes.NotFound, "Service with ID %q not found.", id)
	}
	return nil
}

// ServicesForAgent returns all Services for which Agent with given ID provides insights.
func ServicesForAgent(q *reform.Querier, agentID string) ([]*Service, error) {
	structs, err := q.FindAllFrom(AgentServiceView, "agent_id", agentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Service IDs")
	}

	serviceIDs := make([]interface{}, len(structs))
	for i, s := range structs {
		serviceIDs[i] = s.(*AgentService).ServiceID
	}
	if len(serviceIDs) == 0 {
		return []*Service{}, nil
	}

	p := strings.Join(q.Placeholders(1, len(serviceIDs)), ", ")
	tail := fmt.Sprintf("WHERE service_id IN (%s) ORDER BY service_id", p) //nolint:gosec
	structs, err = q.SelectAllFrom(ServiceTable, tail, serviceIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select Services")
	}

	res := make([]*Service, len(structs))
	for i, s := range structs {
		res[i] = s.(*Service)
	}
	return res, nil
}

// ToInventoryService converts database row to Inventory API Service.
func ToInventoryService(row *Service) (inventorypb.Service, error) {
	labels, err := row.GetCustomLabels()
	if err != nil {
		return nil, err
	}

	switch row.ServiceType {
	case MySQLServiceType:
		return &inventorypb.MySQLService{
			ServiceId:    row.ServiceID,
			ServiceName:  row.ServiceName,
			NodeId:       row.NodeID,
			Address:      pointer.GetString(row.Address),
			Port:         uint32(pointer.GetUint16(row.Port)),
			CustomLabels: labels,
		}, nil
	case MongoDBServiceType:
		return &inventorypb.MongoDBService{
			ServiceId:    row.ServiceID,
			ServiceName:  row.ServiceName,
			NodeId:       row.NodeID,
			Address:      pointer.GetString(row.Address),
			Port:         uint32(pointer.GetUint16(row.Port)),
			CustomLabels: labels,
		}, nil
	case PostgreSQLServiceType:
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

// ToInventoryServices converts database rows to Inventory API Services.
func ToInventoryServices(services []*Service) ([]inventorypb.Service, error) {
	var err error
	res := make([]inventorypb.Service, len(services))
	for i, srv := range services {
		res[i], err = ToInventoryService(srv)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

//go:generate reform

// ServiceType represents Service type as stored in database.
type ServiceType string

// Service types.
const (
	MySQLServiceType      ServiceType = "mysql"
	MongoDBServiceType    ServiceType = "mongodb"
	PostgreSQLServiceType ServiceType = "postgresql"
)

// Service represents Service as stored in database.
//reform:services
type Service struct {
	ServiceID    string      `reform:"service_id,pk"`
	ServiceType  ServiceType `reform:"service_type"`
	ServiceName  string      `reform:"service_name"`
	NodeID       string      `reform:"node_id"`
	CustomLabels []byte      `reform:"custom_labels"`
	CreatedAt    time.Time   `reform:"created_at"`
	UpdatedAt    time.Time   `reform:"updated_at"`

	Address *string `reform:"address"`
	Port    *uint16 `reform:"port"`
}

// BeforeInsert implements reform.BeforeInserter interface.
//nolint:unparam
func (s *Service) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
//nolint:unparam
func (s *Service) BeforeUpdate() error {
	s.UpdatedAt = Now()
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// AfterFind implements reform.AfterFinder interface.
//nolint:unparam
func (s *Service) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	if len(s.CustomLabels) == 0 {
		s.CustomLabels = nil
	}
	return nil
}

// GetCustomLabels decodes custom labels.
func (s *Service) GetCustomLabels() (map[string]string, error) {
	if len(s.CustomLabels) == 0 {
		return nil, nil
	}
	m := make(map[string]string)
	if err := json.Unmarshal(s.CustomLabels, &m); err != nil {
		return nil, errors.Wrap(err, "failed to decode custom labels")
	}
	return m, nil
}

// SetCustomLabels encodes custom labels.
func (s *Service) SetCustomLabels(m map[string]string) error {
	if len(m) == 0 {
		s.CustomLabels = nil
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "failed to encode custom labels")
	}
	s.CustomLabels = b
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*Service)(nil)
	_ reform.BeforeUpdater  = (*Service)(nil)
	_ reform.AfterFinder    = (*Service)(nil)
)
