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

package dbaas

import (
	"context"
	"fmt"
	"strings"

	goversion "github.com/hashicorp/go-version"
	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
	dbaasv1beta1 "github.com/percona/pmm/api/managementpb/dbaas"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/utils/stringset"
	pmmversion "github.com/percona/pmm/version"
)

type componentsService struct {
	l                    *logrus.Entry
	db                   *reform.DB
	dbaasClient          dbaasClient
	versionServiceClient versionService
}

// NewComponentsService creates Components Service.
func NewComponentsService(db *reform.DB, dbaasClient dbaasClient, versionServiceClient versionService) dbaasv1beta1.ComponentsServer {
	l := logrus.WithField("component", "components_service")
	return &componentsService{
		l:                    l,
		db:                   db,
		dbaasClient:          dbaasClient,
		versionServiceClient: versionServiceClient,
	}
}

// Enabled returns if service is enabled and can be used.
func (c *componentsService) Enabled() bool {
	settings, err := models.GetSettings(c.db)
	if err != nil {
		c.l.WithError(err).Error("can't get settings")
		return false
	}
	return settings.DBaaS.Enabled
}

func (c componentsService) GetPSMDBComponents(ctx context.Context, req *dbaasv1beta1.GetPSMDBComponentsRequest) (*dbaasv1beta1.GetPSMDBComponentsResponse, error) {
	var kubernetesCluster *models.KubernetesCluster
	params := componentsParams{
		product:        psmdbOperator,
		versionToApply: req.DbVersion,
	}
	if req.KubernetesClusterName != "" {
		var err error
		kubernetesCluster, err = models.FindKubernetesClusterByName(c.db.Querier, req.KubernetesClusterName)
		if err != nil {
			return nil, err
		}

		checkResponse, e := c.dbaasClient.CheckKubernetesClusterConnection(ctx, kubernetesCluster.KubeConfig)
		if e != nil {
			return nil, e
		}

		if checkResponse.Operators.Psmdb != nil {
			params.productVersion = checkResponse.Operators.Psmdb.Version
		}
	}

	versions, err := c.versions(ctx, params, kubernetesCluster)
	if err != nil {
		return nil, err
	}
	return &dbaasv1beta1.GetPSMDBComponentsResponse{Versions: versions}, nil
}

func (c componentsService) GetPXCComponents(ctx context.Context, req *dbaasv1beta1.GetPXCComponentsRequest) (*dbaasv1beta1.GetPXCComponentsResponse, error) {
	var kubernetesCluster *models.KubernetesCluster
	params := componentsParams{
		product:        pxcOperator,
		versionToApply: req.DbVersion,
	}
	if req.KubernetesClusterName != "" {
		var err error
		kubernetesCluster, err = models.FindKubernetesClusterByName(c.db.Querier, req.KubernetesClusterName)
		if err != nil {
			return nil, err
		}

		checkResponse, e := c.dbaasClient.CheckKubernetesClusterConnection(ctx, kubernetesCluster.KubeConfig)
		if e != nil {
			return nil, e
		}

		if checkResponse.Operators.Xtradb != nil {
			params.productVersion = checkResponse.Operators.Xtradb.Version
		}
	}

	versions, err := c.versions(ctx, params, kubernetesCluster)
	if err != nil {
		return nil, err
	}
	return &dbaasv1beta1.GetPXCComponentsResponse{Versions: versions}, nil
}

func (c componentsService) ChangePSMDBComponents(ctx context.Context, req *dbaasv1beta1.ChangePSMDBComponentsRequest) (*dbaasv1beta1.ChangePSMDBComponentsResponse, error) {
	err := c.db.InTransaction(func(tx *reform.TX) error {
		kubernetesCluster, e := models.FindKubernetesClusterByName(tx.Querier, req.KubernetesClusterName)
		if e != nil {
			return e
		}

		if req.Mongod != nil {
			kubernetesCluster.Mongod, e = setComponent(kubernetesCluster.Mongod, req.Mongod)
			if e != nil {
				message := fmt.Sprintf("%s, cluster: %s, component: mongod", e.Error(), kubernetesCluster.KubernetesClusterName)
				return status.Errorf(codes.InvalidArgument, message)
			}
		}

		e = tx.Save(kubernetesCluster)
		if e != nil {
			return e
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.ChangePSMDBComponentsResponse{}, nil
}

func (c componentsService) ChangePXCComponents(ctx context.Context, req *dbaasv1beta1.ChangePXCComponentsRequest) (*dbaasv1beta1.ChangePXCComponentsResponse, error) {
	err := c.db.InTransaction(func(tx *reform.TX) error {
		kubernetesCluster, e := models.FindKubernetesClusterByName(tx.Querier, req.KubernetesClusterName)
		if e != nil {
			return e
		}

		if req.Pxc != nil {
			kubernetesCluster.PXC, e = setComponent(kubernetesCluster.PXC, req.Pxc)
			if e != nil {
				message := fmt.Sprintf("%s, cluster: %s, component: pxc", e.Error(), kubernetesCluster.KubernetesClusterName)
				return status.Errorf(codes.InvalidArgument, message)
			}
		}

		if req.Proxysql != nil {
			kubernetesCluster.ProxySQL, e = setComponent(kubernetesCluster.ProxySQL, req.Proxysql)
			if e != nil {
				message := fmt.Sprintf("%s, cluster: %s, component: proxySQL", e.Error(), kubernetesCluster.KubernetesClusterName)
				return status.Errorf(codes.InvalidArgument, message)
			}
		}

		if req.Haproxy != nil {
			kubernetesCluster.HAProxy, e = setComponent(kubernetesCluster.HAProxy, req.Haproxy)
			if e != nil {
				message := fmt.Sprintf("%s, cluster: %s, component: HAProxy", e.Error(), kubernetesCluster.KubernetesClusterName)
				return status.Errorf(codes.InvalidArgument, message)
			}
		}
		e = tx.Save(kubernetesCluster)
		if e != nil {
			return e
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dbaasv1beta1.ChangePXCComponentsResponse{}, nil
}

type checkResponse struct {
	err                      error
	installedOperatorVersion string
}

func (c componentsService) CheckForOperatorUpdate(ctx context.Context, req *dbaasv1beta1.CheckForOperatorUpdateRequest) (*dbaasv1beta1.CheckForOperatorUpdateResponse, error) {
	if req.OperatorType != pxcOperator && req.OperatorType != psmdbOperator {
		return nil, errors.Errorf("%q is an unsupported operator type", req.OperatorType)
	}
	if pmmversion.PMMVersion == "" {
		return nil, status.Error(codes.Internal, "failed to get current PMM version")
	}

	responseCh := make(chan checkResponse)
	// Fetch installed version of the operator.
	go func(responseCh chan checkResponse, db *reform.Querier) {
		kubernetesCluster, err := models.FindKubernetesClusterByName(db, req.KubernetesClusterName)
		if err != nil {
			responseCh <- checkResponse{
				err: err,
			}
			return
		}
		resp, err := c.dbaasClient.CheckKubernetesClusterConnection(ctx, kubernetesCluster.KubeConfig)
		if err != nil {
			responseCh <- checkResponse{
				err: err,
			}
			return
		}
		var operator *controllerv1beta1.Operator
		switch req.OperatorType {
		case pxcOperator:
			operator = resp.Operators.Xtradb
		case psmdbOperator:
			operator = resp.Operators.Psmdb
		}
		responseCh <- checkResponse{
			installedOperatorVersion: operator.Version,
		}
	}(responseCh, c.db.Querier)

	pmmVersionParts := strings.Split(pmmversion.PMMVersion, "-")
	latestOperator, latestPMM, err := c.versionServiceClient.GetLatestOperatorVersion(ctx, req.OperatorType, pmmVersionParts[0])
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if latestOperator == nil || latestPMM == nil {
		return nil, status.Error(codes.Internal, "latest versions are of invalid values")
	}

	operatorCheck := <-responseCh
	if operatorCheck.err != nil {
		return nil, status.Error(codes.Internal, operatorCheck.err.Error())
	}
	installedOperatorVersion, err := goversion.NewVersion(operatorCheck.installedOperatorVersion)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if latestOperator.GreaterThan(installedOperatorVersion) {
		return &dbaasv1beta1.CheckForOperatorUpdateResponse{
			Status:                   dbaasv1beta1.OperatorUpdateStatus_UPDATE_AVAILABLE,
			AvailableOperatorVersion: latestOperator.String(),
		}, nil
	} else if latestOperator.Equal(installedOperatorVersion) {
		// We are running the latest operator for installed PMM.
		// Omit PMM version to get latest operator for latest PMM version.
		latestOperator, latestPMM, err := c.versionServiceClient.GetLatestOperatorVersion(ctx, req.OperatorType, "")
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if latestOperator == nil || latestPMM == nil {
			return nil, status.Error(codes.Internal, "latest versions are of invalid values")
		}
		if latestOperator.GreaterThan(installedOperatorVersion) {
			return &dbaasv1beta1.CheckForOperatorUpdateResponse{
				Status:                    dbaasv1beta1.OperatorUpdateStatus_UPDATE_AVAILABLE_BUT_NOT_COMPATIBLE,
				AvailableOperatorVersion:  latestOperator.String(),
				AvailablePmmServerVersion: latestPMM.String(),
			}, nil
		}
		return &dbaasv1beta1.CheckForOperatorUpdateResponse{
			Status: dbaasv1beta1.OperatorUpdateStatus_UPDATE_NOT_AVAILABLE,
		}, nil
	}
	return &dbaasv1beta1.CheckForOperatorUpdateResponse{
			Status: dbaasv1beta1.OperatorUpdateStatus_UPDATE_NOT_AVAILABLE,
		},
		status.Errorf(
			codes.NotFound,
			"no more up-to-date version found, %s is ahead of latest compatible version of the operator",
			installedOperatorVersion.String(),
		)
}

func (c componentsService) versions(ctx context.Context, params componentsParams, cluster *models.KubernetesCluster) ([]*dbaasv1beta1.OperatorVersion, error) {
	components, err := c.versionServiceClient.Matrix(ctx, params)
	if err != nil {
		return nil, err
	}

	var mongod, pxc, proxySQL, haproxy *models.Component
	if cluster != nil {
		mongod = cluster.Mongod
		pxc = cluster.PXC
		proxySQL = cluster.ProxySQL
		haproxy = cluster.HAProxy
	}

	versions := make([]*dbaasv1beta1.OperatorVersion, 0, len(components.Versions))
	mongodMinimalVersion, _ := goversion.NewVersion("4.2.0")
	pxcMinimalVersion, _ := goversion.NewVersion("8.0.0")
	for _, v := range components.Versions {
		respVersion := &dbaasv1beta1.OperatorVersion{
			Product:  v.Product,
			Operator: v.ProductVersion,
			Matrix: &dbaasv1beta1.Matrix{
				Mongod:       c.matrix(v.Matrix.Mongod, mongodMinimalVersion, mongod),
				Pxc:          c.matrix(v.Matrix.Pxc, pxcMinimalVersion, pxc),
				Pmm:          c.matrix(v.Matrix.Pmm, nil, nil),
				Proxysql:     c.matrix(v.Matrix.Proxysql, nil, proxySQL),
				Haproxy:      c.matrix(v.Matrix.Haproxy, nil, haproxy),
				Backup:       c.matrix(v.Matrix.Backup, nil, nil),
				Operator:     c.matrix(v.Matrix.Operator, nil, nil),
				LogCollector: c.matrix(v.Matrix.LogCollector, nil, nil),
			},
		}
		versions = append(versions, respVersion)
	}

	return versions, nil
}

func (c componentsService) matrix(m map[string]componentVersion, minimalVersion *goversion.Version, kc *models.Component) map[string]*dbaasv1beta1.Component {
	result := make(map[string]*dbaasv1beta1.Component)

	var lastVersion string
	lastVersionParsed, _ := goversion.NewVersion("0.0.0")
	for v, component := range m {
		parsedVersion, err := goversion.NewVersion(v)
		if err != nil {
			c.l.Warnf("couldn't parse version %s: %s", v, err.Error())
			continue
		}
		if minimalVersion != nil && parsedVersion.LessThan(minimalVersion) {
			continue
		}
		result[v] = &dbaasv1beta1.Component{
			ImagePath: component.ImagePath,
			ImageHash: component.ImageHash,
			Status:    component.Status,
			Critical:  component.Critical,
		}
		if lastVersionParsed.LessThan(parsedVersion) && component.Status == "recommended" {
			lastVersionParsed = parsedVersion
			lastVersion = v
		}
	}

	defaultVersionSet := false
	if kc != nil {
		if _, ok := result[kc.DefaultVersion]; ok {
			result[kc.DefaultVersion].Default = true
			defaultVersionSet = true
		}
		for _, v := range kc.DisabledVersions {
			if _, ok := result[v]; ok {
				result[v].Disabled = true
			}
		}
	}
	if lastVersion != "" && !defaultVersionSet {
		result[lastVersion].Default = true
	}
	return result
}

func setComponent(kc *models.Component, rc *dbaasv1beta1.ChangeComponent) (*models.Component, error) {
	if kc == nil {
		kc = new(models.Component)
	}
	if rc.DefaultVersion != "" {
		kc.DefaultVersion = rc.DefaultVersion
	}

	disabledVersions := make(map[string]struct{})
	for _, v := range kc.DisabledVersions {
		disabledVersions[v] = struct{}{}
	}
	for _, v := range rc.Versions {
		if v.Enable && v.Disable {
			return nil, fmt.Errorf("enable and disable for version %s can't be passed together", v.Version)
		}
		if v.Enable {
			delete(disabledVersions, v.Version)
		}
		if v.Disable {
			disabledVersions[v.Version] = struct{}{}
		}
	}
	if _, ok := disabledVersions[kc.DefaultVersion]; ok {
		return nil, fmt.Errorf("default version can't be disabled")
	}
	kc.DisabledVersions = stringset.ToSlice(disabledVersions)
	return kc, nil
}
