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

package prometheus

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"

	"github.com/percona/pmm-managed/models"
	config_util "github.com/percona/pmm-managed/services/prometheus/internal/common/config"
	"github.com/percona/pmm-managed/services/prometheus/internal/prometheus/config"
	sd_config "github.com/percona/pmm-managed/services/prometheus/internal/prometheus/discovery/config"
	"github.com/percona/pmm-managed/services/prometheus/internal/prometheus/discovery/targetgroup"
)

const addressLabel = model.LabelName(model.AddressLabel)

func scrapeConfig(node *models.Node, service *models.Service, agent *models.Agent) (*config.ScrapeConfig, error) {
	cfg, err := scrapeConfigBase(agent)
	if err != nil {
		return nil, err
	}

	switch agent.AgentType {
	case models.NodeExporterType:
		err = scrapeConfigForNodeExporter(cfg, node, agent)
	case models.MySQLdExporterType:
		err = scrapeConfigForMySQLdExporter(cfg, service, agent)
	default:
		panic(fmt.Errorf("unhandled Agent type %s", agent.AgentType))
	}
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func scrapeConfigBase(agent *models.Agent) (*config.ScrapeConfig, error) {
	cfg := &config.ScrapeConfig{
		JobName: strings.Replace(agent.AgentID, "/", "_", -1),
	}

	var labels model.LabelSet
	if err := labels.UnmarshalJSON(agent.CustomLabels); err != nil {
		return nil, errors.Wrap(err, "failed to set custom labels")
	}
	cfg.ServiceDiscoveryConfig = sd_config.ServiceDiscoveryConfig{
		StaticConfigs: []*targetgroup.Group{{
			Labels: labels,
		}},
	}

	username, password := pointer.GetString(agent.Username), pointer.GetString(agent.Password)
	if username != "" {
		cfg.HTTPClientConfig = config_util.HTTPClientConfig{
			BasicAuth: &config_util.BasicAuth{
				Username: username,
				Password: password,
			},
		}
	}

	return cfg, nil
}

func scrapeConfigForNodeExporter(cfg *config.ScrapeConfig, node *models.Node, agent *models.Agent) error {
	panic(fmt.Errorf("unhandled Agent type %s", agent.AgentType))
}

func scrapeConfigForMySQLdExporter(cfg *config.ScrapeConfig, service *models.Service, agent *models.Agent) error {
	port := pointer.GetUint16(agent.ListenPort)
	if port == 0 {
		return errors.New("listen port is not known")
	}
	hostport := net.JoinHostPort(pointer.GetString(service.Address), strconv.Itoa(int(port)))
	target := model.LabelSet{addressLabel: model.LabelValue(hostport)}
	if err := target.Validate(); err != nil {
		return errors.Wrap(err, "failed to set targets")
	}
	cfg.ServiceDiscoveryConfig.StaticConfigs[0].Targets = []model.LabelSet{target}

	return nil
}
