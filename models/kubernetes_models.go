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
	"time"

	"gopkg.in/reform.v1"
)

//go:generate reform

// KubernetesCluster represents a kubernetes cluster as stored in database.
//reform:kubernetes_clusters
type KubernetesCluster struct {
	ID         string    `reform:"id,pk"`
	Name       string    `reform:"name"`
	KubeConfig string    `reform:"kube_config"`
	CreatedAt  time.Time `reform:"created_at"`
	UpdatedAt  time.Time `reform:"updated_at"`
}

// BeforeInsert implements reform.BeforeInserter interface.
func (s *KubernetesCluster) BeforeInsert() error {
	now := Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

// BeforeUpdate implements reform.BeforeUpdater interface.
func (s *KubernetesCluster) BeforeUpdate() error {
	s.UpdatedAt = Now()
	return nil
}

// AfterFind implements reform.AfterFinder interface.
func (s *KubernetesCluster) AfterFind() error {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return nil
}

// check interfaces
var (
	_ reform.BeforeInserter = (*ActionResult)(nil)
	_ reform.BeforeUpdater  = (*ActionResult)(nil)
	_ reform.AfterFinder    = (*ActionResult)(nil)
)
