package ia

import "github.com/percona/pmm-managed/models"

// aletringService is a subset of methods of ia.Service used by this package.
// We use it instead of real type for testing and to avoid dependency cycle.
type aletringService interface {
	AddChannel(ch *models.Channel) error
	ChangeChannel(ch *models.Channel) error
	RemoveChannel(id string) error
	ListChannels() ([]models.Channel, error)
}
