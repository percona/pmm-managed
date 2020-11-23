// package ia implements integrated alerting logic.
package ia

import (
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

// Service is responsible for integrated alerting.
type Service struct {
	db *reform.DB
}

func New(db *reform.DB) *Service {
	return &Service{
		db: db,
	}
}

// AddChannel adds new notification channel.
func (s *Service) AddChannel(ch *models.Channel) error {
	return models.SaveChannel(s.db, ch)
}

// ChangeChannel changes existing notification channel.
func (s *Service) ChangeChannel(ch *models.Channel) error {
	return models.UpdateChannel(s.db, ch)
}

// RemoveChannel removes notification channel.
func (s *Service) RemoveChannel(id string) error {
	return models.RemoveChannel(s.db, id)
}

// ListChannels returns list of available channels.
func (s *Service) ListChannels() ([]models.Channel, error) {
	return models.GetChannels(s.db)
}
