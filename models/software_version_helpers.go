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
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

// CreateServiceSoftwareVersionsParams are params for creating a new service software versions entry.
type CreateServiceSoftwareVersionsParams struct {
	ServiceID        string
	SoftwareVersions []SoftwareVersion
	NextCheckAt      time.Time
}

// Validate validates params used for creating a service software versions entry.
func (p *CreateServiceSoftwareVersionsParams) Validate() error {
	if p.ServiceID == "" {
		return errors.Wrap(ErrInvalidArgument, "service_id shouldn't be empty")
	}

	for _, sv := range p.SoftwareVersions {
		switch sv.Name {
		case MysqldSoftwareName:
		case XtrabackupSoftwareName:
		case XbcloudSoftwareName:
		case QpressSoftwareName:
		default:
			return errors.Wrapf(ErrInvalidArgument, "invalid software name %q", sv.Name)
		}
	}

	return nil
}

// CreateServiceSoftwareVersions creates service software versions entry in DB.
func CreateServiceSoftwareVersions(q *reform.Querier, params CreateServiceSoftwareVersionsParams) (*ServiceSoftwareVersions, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	row := &ServiceSoftwareVersions{
		ServiceID:        params.ServiceID,
		SoftwareVersions: params.SoftwareVersions,
		NextCheckAt:      params.NextCheckAt,
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to insert service software versions")
	}

	return row, nil
}

// UpdateServiceSoftwareVersionsParams represents params for updating service software versions entity.
type UpdateServiceSoftwareVersionsParams struct {
	NextCheckAt      *time.Time
	SoftwareVersions *[]SoftwareVersion
}

// Validate validates params used for updating a service software versions entry.
func (u *UpdateServiceSoftwareVersionsParams) Validate() error {
	if u.SoftwareVersions == nil {
		return nil
	}

	for _, sv := range *u.SoftwareVersions {
		switch sv.Name {
		case MysqldSoftwareName:
		case XtrabackupSoftwareName:
		case XbcloudSoftwareName:
		case QpressSoftwareName:
		default:
			return errors.Wrapf(ErrInvalidArgument, "invalid software name %q", sv.Name)
		}
	}

	return nil
}

// UpdateServiceSoftwareVersions updates existing service software versions.
func UpdateServiceSoftwareVersions(
	q *reform.Querier,
	serviceID string,
	params UpdateServiceSoftwareVersionsParams,
) (*ServiceSoftwareVersions, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	row, err := FindServiceSoftwareVersions(q, serviceID)
	if err != nil {
		return nil, err
	}
	if params.NextCheckAt != nil {
		row.NextCheckAt = *params.NextCheckAt
	}
	if params.SoftwareVersions != nil {
		row.SoftwareVersions = *params.SoftwareVersions
	}

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update service software versions")
	}

	return row, nil
}

// FindServiceSoftwareVersions returns service software versions entry by given service ID if found, ErrNotFound if not.
func FindServiceSoftwareVersions(q *reform.Querier, serviceID string) (*ServiceSoftwareVersions, error) {
	if serviceID == "" {
		return nil, errors.New("service id is empty")
	}

	versions := &ServiceSoftwareVersions{ServiceID: serviceID}
	switch err := q.Reload(versions); err {
	case nil:
		return versions, nil
	case reform.ErrNoRows:
		return nil, errors.Wrapf(ErrNotFound, "service software versions by service id '%s'", serviceID)
	default:
		return nil, errors.WithStack(err)
	}
}

// FindServicesSoftwareVersionsFilter represents a filter for finding service software versions.
type FindServicesSoftwareVersionsFilter struct {
	Limit *int
}

// FindServicesSoftwareVersions returns all services software versions.
func FindServicesSoftwareVersions(q *reform.Querier, filter FindServicesSoftwareVersionsFilter) ([]*ServiceSoftwareVersions, error) {
	var args []interface{}
	var tail strings.Builder
	tail.WriteString("ORDER BY check_at ")
	if filter.Limit != nil {
		tail.WriteString("LIMIT $1")
		args = append(args, *filter.Limit)
	}

	structs, err := q.SelectAllFrom(ServiceSoftwareVersionsTable, tail.String(), args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	versions := make([]*ServiceSoftwareVersions, len(structs))
	for i, s := range structs {
		versions[i] = s.(*ServiceSoftwareVersions)
	}

	return versions, nil
}

// DeleteServiceSoftwareVersions removes entry from the DB by service ID.
func DeleteServiceSoftwareVersions(q *reform.Querier, serviceID string) error {
	if _, err := FindServiceSoftwareVersions(q, serviceID); err != nil {
		return err
	}

	if err := q.Delete(&ServiceSoftwareVersions{ServiceID: serviceID}); err != nil {
		return errors.Wrapf(err, "failed to delete services software versions by service id '%s'", serviceID)
	}
	return nil
}
