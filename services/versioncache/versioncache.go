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

// Package versioncache provides service software version cache functionality.
package versioncache

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services/agents"
)

const (
	serviceCheckInterval = time.Hour
	relaxDuration        = 5 * time.Second
)

//go:generate mockery -name=Versioner -case=snake -inpkg -testonly

// Versioner contains method for retrieving versions of different software.
type Versioner interface {
	GetVersions(pmmAgentID string, softwares []agents.Software) ([]agents.Version, error)
}

// Service is responsible for caching service software versions in the DB.
type Service struct {
	db       *reform.DB
	l        *logrus.Entry
	v        Versioner
	updateCh chan struct{}
}

// New creates new service.
func New(db *reform.DB, v Versioner) *Service {
	return &Service{
		db:       db,
		l:        logrus.WithField("component", "version-cache"),
		v:        v,
		updateCh: make(chan struct{}, 1),
	}
}

func (s *Service) syncServices() error {
	err := s.db.InTransaction(func(tx *reform.TX) error {
		serviceType := models.MySQLServiceType
		services, err := models.FindServices(tx.Querier, models.ServiceFilters{ServiceType: &serviceType})
		if err != nil {
			return err
		}

		serviceIDs := make(map[string]struct{}, len(services))
		for _, s := range services {
			serviceIDs[s.ServiceID] = struct{}{}
		}

		serviceVersions, err := models.FindServicesSoftwareVersions(tx.Querier, nil)
		if err != nil {
			return err
		}

		// remove services software versions from the cache which are no longer exist
		for _, sv := range serviceVersions {
			if _, ok := serviceIDs[sv.ServiceID]; !ok {
				if err := models.DeleteServiceSoftwareVersions(tx.Querier, sv.ServiceID); err != nil {
					return err
				}
			}
		}

		// add new services software versions to the cache
		cacheServiceIDs := make(map[string]struct{}, len(serviceVersions))
		for _, sv := range serviceVersions {
			cacheServiceIDs[sv.ServiceID] = struct{}{}
		}
		for _, s := range services {
			if _, ok := cacheServiceIDs[s.ServiceID]; !ok {
				if _, err := models.CreateServiceSoftwareVersions(tx.Querier, models.CreateServiceSoftwareVersionsParams{
					ServiceID:        s.ServiceID,
					SoftwareVersions: []models.SoftwareVersion{},
					CheckAt:          time.Now().UTC(),
				}); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

type prepareResults struct {
	ServiceID   string
	CheckAfter  time.Duration
	NeedsUpdate bool
	AgentID     string
}

func (s *Service) prepareUpdateVersions() (*prepareResults, error) {
	results := &prepareResults{CheckAfter: relaxDuration}

	if err := s.db.InTransaction(func(tx *reform.TX) error {
		limit := 1
		servicesVersions, err := models.FindServicesSoftwareVersions(tx.Querier, &limit)
		if err != nil {
			return err
		}
		if len(servicesVersions) == 0 {
			results.CheckAfter = serviceCheckInterval

			return nil
		}
		if servicesVersions[0].CheckAt.After(time.Now().UTC()) {
			results.CheckAfter = servicesVersions[0].CheckAt.Sub(time.Now().UTC())
			if results.CheckAfter < relaxDuration {
				results.CheckAfter = relaxDuration
			}

			return nil
		}

		results.NeedsUpdate = true
		results.ServiceID = servicesVersions[0].ServiceID

		service, err := models.FindServiceByID(tx.Querier, servicesVersions[0].ServiceID)
		if err != nil {
			return err
		}
		if service.ServiceType != models.MySQLServiceType {
			return nil
		}

		pmmAgents, err := models.FindPMMAgentsForService(tx.Querier, servicesVersions[0].ServiceID)
		if err != nil {
			return err
		}
		if len(pmmAgents) == 0 {
			return errors.Errorf("pmmAgent not found for service")
		}
		results.AgentID = pmmAgents[0].AgentID

		checkAt := time.Now().UTC().Add(serviceCheckInterval) // TODO
		if _, err := models.UpdateServiceSoftwareVersions(tx.Querier, servicesVersions[0].ServiceID,
			models.UpdateServiceSoftwareVersionsParams{CheckAt: &checkAt},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Service) updateVersions() (time.Duration, error) {
	r, err := s.prepareUpdateVersions()
	if err != nil {
		return relaxDuration, err
	}

	if !r.NeedsUpdate {
		return r.CheckAfter, nil
	}

	softwares := []agents.Software{&agents.Mysqld{}, &agents.Xtrabackup{}, &agents.Xbcloud{}, &agents.Qpress{}}
	versions, err := s.v.GetVersions(r.AgentID, softwares)
	if err != nil {
		return relaxDuration, err
	}
	if len(versions) != len(softwares) {
		return relaxDuration, errors.Errorf("slices length mismatch: versions len %d != softwares len %d",
			len(versions), len(softwares))
	}

	svs := make([]models.SoftwareVersion, 0, len(softwares))
	for i, software := range softwares {
		var softwareName models.SoftwareName
		switch software := software.(type) {
		case *agents.Mysqld:
			softwareName = models.MysqldSoftwareName
		case *agents.Xtrabackup:
			softwareName = models.XtrabackupSoftwareName
		case *agents.Xbcloud:
			softwareName = models.XbcloudSoftwareName
		case *agents.Qpress:
			softwareName = models.QpressSoftwareName
		default:
			return relaxDuration, errors.Errorf("invalid software type %T", software)
		}

		if versions[i].Error != "" {
			s.l.Warnf("failed to get version of %q software: %s", softwareName, versions[i].Error)
			continue
		}
		if versions[i].Version == "" {
			continue
		}

		svs = append(svs, models.SoftwareVersion{
			Name:    softwareName,
			Version: versions[i].Version,
		})
	}

	if _, err := models.UpdateServiceSoftwareVersions(s.db.Querier, r.ServiceID,
		models.UpdateServiceSoftwareVersionsParams{SoftwareVersions: &svs},
	); err != nil {
		return relaxDuration, err
	}

	return relaxDuration, err
}

// SyncAndUpdate triggers sync and update service software versions.
func (s *Service) SyncAndUpdate() {
	select {
	case s.updateCh <- struct{}{}:
	default:
	}
}

// Run runs software version cache service.
func (s *Service) Run(ctx context.Context) {
	defer close(s.updateCh)

	if err := s.syncServices(); err != nil {
		s.l.Warn(err)
	}

	syncTicker := time.NewTicker(serviceCheckInterval)
	defer syncTicker.Stop()

	var checkAfter time.Duration
	for {
		select {
		case <-time.After(checkAfter):
			ca, err := s.updateVersions()
			if err != nil {
				s.l.Warn(err)
			}

			checkAfter = ca
		case <-syncTicker.C:
			if err := s.syncServices(); err != nil {
				s.l.Warn(err)
			}
		case <-s.updateCh:
			if err := s.syncServices(); err != nil {
				s.l.Warn(err)
			}

			ca, err := s.updateVersions()
			if err != nil {
				s.l.Warn(err)
			}

			checkAfter = ca
		case <-ctx.Done():
			return
		}
	}
}
