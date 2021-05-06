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
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

var (
	// ErrNotFound returned when entity is not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidArgument returned when some passed argument is invalid.
	ErrInvalidArgument = errors.New("invalid argument")
)

// ArtifactFilters represents filters for artifacts list.
type ArtifactFilters struct {
	// Return only Agents that provide insights for that Service.
	ServiceID string
}

// FindArtifacts returns artifacts list.
func FindArtifacts(q *reform.Querier, filters *ArtifactFilters) ([]*Artifact, error) {
	var conditions []string
	var args []interface{}
	if filters != nil && filters.ServiceID != "" {
		if _, err := FindServiceByID(q, filters.ServiceID); err != nil {
			return nil, err
		}
		conditions = append(conditions, fmt.Sprintf("service_id = %s", q.Placeholder(1)))
		args = append(args, filters.ServiceID)
	}

	var whereClause string
	if len(conditions) != 0 {
		whereClause = fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
	}
	rows, err := q.SelectAllFrom(ArtifactTable, fmt.Sprintf("%s ORDER BY created_at DESC", whereClause), args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select artifacts")
	}

	artifacts := make([]*Artifact, 0, len(rows))
	for _, r := range rows {
		artifacts = append(artifacts, r.(*Artifact))
	}

	return artifacts, nil
}

// FindArtifactsByIDs finds artifacts by IDs.
func FindArtifactsByIDs(q *reform.Querier, ids []string) (map[string]*Artifact, error) {
	if len(ids) == 0 {
		return map[string]*Artifact{}, nil
	}

	p := strings.Join(q.Placeholders(1, len(ids)), ", ")
	tail := fmt.Sprintf("WHERE id IN (%s)", p)
	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	all, err := q.SelectAllFrom(ArtifactTable, tail, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	artifacts := make(map[string]*Artifact, len(all))
	for _, l := range all {
		artifact := l.(*Artifact)
		artifacts[artifact.ID] = artifact
	}
	return artifacts, nil
}

// FindArtifactByID returns artifact by given ID if found, ErrNotFound if not.
func FindArtifactByID(q *reform.Querier, id string) (*Artifact, error) {
	if id == "" {
		return nil, errors.New("provided artifact id is empty")
	}

	artifact := &Artifact{ID: id}
	switch err := q.Reload(artifact); err {
	case nil:
		return artifact, nil
	case reform.ErrNoRows:
		return nil, errors.Wrapf(ErrNotFound, "artifact by id '%s'", id)
	default:
		return nil, errors.WithStack(err)
	}
}

// CreateArtifactParams are params for creating a new artifact.
type CreateArtifactParams struct {
	Name       string
	Vendor     string
	LocationID string
	ServiceID  string
	DataModel  DataModel
	Status     BackupStatus
}

// Validate validates params used for creating an artifact entry.
func (p *CreateArtifactParams) Validate() error {
	if p.Name == "" {
		return errors.Wrap(ErrInvalidArgument, "name shouldn't be empty")
	}
	if p.Vendor == "" {
		return errors.Wrap(ErrInvalidArgument, "vendor shouldn't be empty")
	}
	if p.LocationID == "" {
		return errors.Wrap(ErrInvalidArgument, "location_id shouldn't be empty")
	}
	if p.ServiceID == "" {
		return errors.Wrap(ErrInvalidArgument, "service_id shouldn't be empty")
	}

	if err := p.DataModel.Validate(); err != nil {
		return err
	}

	return p.Status.Validate()
}

// CreateArtifact creates artifact entry in DB.
func CreateArtifact(q *reform.Querier, params CreateArtifactParams) (*Artifact, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	id := "/artifact_id/" + uuid.New().String()
	_, err := FindArtifactByID(q, id)
	switch {
	case err == nil:
		return nil, errors.Errorf("artifact with id '%s' already exists", id)
	case errors.Is(err, ErrNotFound):
	default:
		return nil, errors.WithStack(err)
	}

	row := &Artifact{
		ID:         id,
		Name:       params.Name,
		Vendor:     params.Vendor,
		LocationID: params.LocationID,
		ServiceID:  params.ServiceID,
		DataModel:  params.DataModel,
		Status:     params.Status,
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to insert artifact")
	}

	return row, nil
}

// ChangeArtifactParams are params for changing existing artifact.
type ChangeArtifactParams struct {
	Status BackupStatus
}

// ChangeArtifact updates existing artifact.
func ChangeArtifact(q *reform.Querier, artifactID string, params ChangeArtifactParams) (*Artifact, error) {
	row, err := FindArtifactByID(q, artifactID)
	if err != nil {
		return nil, err
	}
	row.Status = params.Status

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update backup artifact")
	}

	return row, nil
}

// RemoveArtifact removes artifact by ID.
func RemoveArtifact(q *reform.Querier, id string) error {
	if _, err := FindArtifactByID(q, id); err != nil {
		return err
	}

	if err := q.Delete(&Artifact{ID: id}); err != nil {
		return errors.Wrapf(err, "failed to delete artifact by id '%s'", id)
	}
	return nil
}
