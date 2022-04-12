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

package agents

import (
	"github.com/percona/pmm/api/agentpb"
	"github.com/pkg/errors"

	"github.com/percona/pmm-managed/models"
)

type ParseDefaultsFile struct {
	r *Registry
}

func NewParseDefaultsFile(r *Registry) *ParseDefaultsFile {
	return &ParseDefaultsFile{
		r: r,
	}
}

func (p *ParseDefaultsFile) ParseDefaultsFile(pmmAgentID string, filePath string, serviceType models.ServiceType) (string, string, error) {
	// TODO: add logger
	pmmAgent, err := p.r.get(pmmAgentID)
	if err != nil {
		return "", "", err
	}

	request, err := createRequest(filePath, serviceType)
	if err != nil {
		return "", "", err
	}

	resp, err := pmmAgent.channel.SendAndWaitResponse(request)
	if err != nil {
		return "", "", err
	}

	return resp.(*agentpb.ParseDefaultsFileResponse).GetUsername(), resp.(*agentpb.ParseDefaultsFileResponse).GetPassword(), nil
}

func createRequest(filePath string, serviceType models.ServiceType) (*agentpb.ParseDefaultsFileRequest, error) {
	var request *agentpb.ParseDefaultsFileRequest

	switch serviceType {
	case models.MySQLServiceType:
		request = &agentpb.ParseDefaultsFileRequest{
			PathToFile: filePath,
		}
	default:
		return nil, errors.Errorf("unhandled service type %s", serviceType)
	}
	return request, nil
}
