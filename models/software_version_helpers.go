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
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

// CreateServiceSoftwareVersionsParams are params for creating a new service software versions entry.
type CreateServiceSoftwareVersionsParams struct {
	ServiceID string
	Versions  []SoftwareVersion
	CheckAt   time.Time
}

// Validate validates params used for creating an service software versions entry.
func (p *CreateServiceSoftwareVersionsParams) Validate() error {
	if p.ServiceID == "" {
		return errors.Wrap(ErrInvalidArgument, "service_id shouldn't be empty")
	}

	return nil
}

// CreateServiceSoftwareVersions creates service software versions entry in DB.
func CreateServiceSoftwareVersions(q *reform.Querier, params CreateServiceSoftwareVersionsParams) (*ServiceSoftwareVersions, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	id := "/service_software_versions_id/" + uuid.New().String()
	_, err := FindServiceSoftwareVersions(q, id)
	switch {
	case err == nil:
		return nil, errors.Errorf("service software versions with id '%s' already exists", id)
	case errors.Is(err, ErrNotFound):
	default:
		return nil, errors.WithStack(err)
	}

	row := &ServiceSoftwareVersions{
		ID:               id,
		ServiceID:        params.ServiceID,
		SoftwareVersions: params.Versions,
		CheckAt:          params.CheckAt,
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to insert service software versions")
	}

	return row, nil
}

// UpdateServiceSoftwareVersionsParams represents params for updating service software versions entity.
type UpdateServiceSoftwareVersionsParams struct {
	CheckAt  *time.Time
	Versions *[]SoftwareVersion
}

// UpdateServiceSoftwareVersions updates existing service software versions.
func UpdateServiceSoftwareVersions(
	q *reform.Querier,
	id string,
	params UpdateServiceSoftwareVersionsParams,
) (*ServiceSoftwareVersions, error) {
	row, err := FindServiceSoftwareVersions(q, id)
	if err != nil {
		return nil, err
	}
	if params.CheckAt != nil {
		row.CheckAt = *params.CheckAt
	}
	if params.Versions != nil {
		row.SoftwareVersions = *params.Versions
	}

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update service software versions")
	}

	return row, nil
}

// FindServiceSoftwareVersions returns service software versions entry by given ID if found, ErrNotFound if not.
func FindServiceSoftwareVersions(q *reform.Querier, id string) (*ServiceSoftwareVersions, error) {
	if id == "" {
		return nil, errors.New("service software versions id is empty")
	}

	versions := &ServiceSoftwareVersions{ID: id}
	switch err := q.Reload(versions); err {
	case nil:
		return versions, nil
	case reform.ErrNoRows:
		return nil, errors.Wrapf(ErrNotFound, "service software versions by id '%s'", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// FindServiceSoftwareVersionsByServiceID returns service software versions entry
// by given service ID if found, ErrNotFound if not.
func FindServiceSoftwareVersionsByServiceID(q *reform.Querier, serviceID string) (*ServiceSoftwareVersions, error) {
	if serviceID == "" {
		return nil, errors.New("service id is empty")
	}

	s, err := q.SelectOneFrom(ServiceTable, "WHERE service_id = $1", serviceID)
	switch err {
	case nil:
		return s.(*ServiceSoftwareVersions), nil
	case reform.ErrNoRows:
		return nil, errors.Wrapf(ErrNotFound, "service software versions by service id '%s'", serviceID)
	default:
		return nil, errors.WithStack(err)
	}
}

// FindServicesSoftwareVersions returns services software versions.
func FindServicesSoftwareVersions(q *reform.Querier, limit *int) ([]*ServiceSoftwareVersions, error) {
	var args []interface{}
	var limitStatement string
	if limit != nil {
		limitStatement = " LIMIT $1"
		args = append(args, *limit)
	}

	const orderByField = "check_at"
	structs, err := q.SelectAllFrom(
		ServiceSoftwareVersionsTable,
		fmt.Sprintf("ORDER BY %s %s", orderByField, limitStatement),
		args...,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	versions := make([]*ServiceSoftwareVersions, len(structs))
	for i, s := range structs {
		versions[i] = s.(*ServiceSoftwareVersions)
	}

	return versions, nil
}

// DeleteServiceSoftwareVersions removes entry from the DB by ID.
func DeleteServiceSoftwareVersions(q *reform.Querier, id string) error {
	if _, err := FindServiceSoftwareVersions(q, id); err != nil {
		return err
	}

	if err := q.Delete(&ServiceSoftwareVersions{ID: id}); err != nil {
		return errors.Wrapf(err, "failed to delete services software versions by id '%s'", id)
	}
	return nil
}
