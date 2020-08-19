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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

func checkUniqueKubernetesClusterID(q *reform.Querier, id string) error {
	if id == "" {
		panic("empty Kubernetes Cluster ID")
	}

	cluster := &KubernetesCluster{ID: id}
	switch err := q.Reload(cluster); err {
	case nil:
		return status.Errorf(codes.AlreadyExists, "Cluster with ID %q already exists.", id)
	case reform.ErrNoRows:
		return nil
	default:
		return errors.WithStack(err)
	}
}

// FindAllKubernetesClusters returns all kubernetes clusters.
func FindAllKubernetesClusters(q *reform.Querier) ([]*KubernetesCluster, error) {
	structs, err := q.SelectAllFrom(KubernetesClusterTable, "ORDER BY id")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clusters := make([]*KubernetesCluster, len(structs))
	for i, s := range structs {
		clusters[i] = s.(*KubernetesCluster)
	}

	return clusters, nil
}

func FindKubernetesClusterByName(q *reform.Querier, name string) (*KubernetesCluster, error) {
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Kubernetes Cluster Name.")
	}

	switch cluster, err := q.FindOneFrom(KubernetesClusterTable, "kubernetes_cluster_name", name); err {
	case nil:
		return cluster.(*KubernetesCluster), nil
	case reform.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "Cluster with name %q not found.", name)
	default:
		return nil, errors.WithStack(err)
	}
}

type CreateKubernetesClusterParams struct {
	KubernetesClusterName string
	KubeConfig            string
}

func CreateKubernetesCluster(q *reform.Querier, params CreateKubernetesClusterParams) (*KubernetesCluster, error) {
	id := "/kubernetes_cluster_id/" + uuid.New().String()
	if err := checkUniqueKubernetesClusterID(q, id); err != nil {
		return nil, err
	}

	row := &KubernetesCluster{
		ID:                    id,
		KubernetesClusterName: params.KubernetesClusterName,
		KubeConfig:            params.KubeConfig,
	}
	if err := q.Insert(row); err != nil {
		return nil, errors.WithStack(err)
	}

	return row, nil
}

func RemoveKubernetesCluster(q *reform.Querier, name string) error {
	c, err := FindKubernetesClusterByName(q, name)
	if err != nil {
		return err
	}

	return errors.Wrap(q.Delete(c), "failed to delete Kubernetes Cluster")
}
