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

// AddChannel adds new notification channel.
func (s *Service) AddChannel(params *models.CreateChannelParams) error {
	_, err := models.CreateChannel(s.db.Querier, params)
	if err != nil {
		return err
	}
	return nil
}

// ChangeChannel changes existing notification channel.
func (s *Service) ChangeChannel(channelID string, params *models.ChangeChannelParams) error {
	_, err := models.ChangeChannel(s.db.Querier, channelID, params)
	if err != nil {
		return err
	}
	return nil
}

// RemoveChannel removes notification channel.
func (s *Service) RemoveChannel(id string) error {
	return models.RemoveChannel(s.db.Querier, id)
}

// ListChannels returns list of available channels.
func (s *Service) ListChannels() ([]models.Channel, error) {
	return models.FindChannels(s.db.Querier)
}
