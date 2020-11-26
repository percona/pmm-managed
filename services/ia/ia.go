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

// Package ia implements integrated alerting logic.
package ia

import (
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// Service is responsible for integrated alerting.
type Service struct {
	db *reform.DB
}

// New creates new IA service.
func New(db *reform.DB) *Service {
	return &Service{
		db: db,
	}
}

// AddRule adds new alert Rule.
func (s *Service) AddRule(r *models.Rule) error {
	return models.SaveRule(s.db.Querier, r)
}

// UpdateRule updates existing alert Rule.
func (s *Service) UpdateRule(r *models.Rule) error {
	return models.UpdateRule(s.db.Querier, r)
}

// RemoveRule removes alert Rule.
func (s *Service) RemoveRule(id string) error {
	return models.RemoveRule(s.db.Querier, id)
}

// ListRules returns list of available alert Rules.
func (s *Service) ListRules() ([]models.Rule, error) {
	return models.GetRules(s.db.Querier)
}
