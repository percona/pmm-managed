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
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	config "github.com/percona/promconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (

	// enables cache for vmdb responses.
	vmCacheEnableEnv = "VM_CACHE_ENABLE"
	// enables victoriametrics services at supervisor config.
	vmTestEnableEnv = "VM_TEST_ENABLE"
)

// VictoriaMetricsParams - defines flags and settings for victoriametrics.
type VictoriaMetricsParams struct {
	// VMDBFlags defines additional flags for VictoriaMetrics DB.
	VMDBFlags []string
	// VMAlertFlags additional flags for VMAlert.
	VMAlertFlags []string
	// Enables VictoriaMetrics.
	Enabled bool
	// BaseConfigPath defines path for basic prometheus config.
	BaseConfigPath string
}

// NewVictoriaMetricsParams - returns configuration params for VictoriaMetrics.
func NewVictoriaMetricsParams(basePath string) (*VictoriaMetricsParams, error) {
	vmp := &VictoriaMetricsParams{
		BaseConfigPath: basePath,
	}
	if enabledVar := os.Getenv(vmTestEnableEnv); enabledVar != "" {
		parsedBool, err := strconv.ParseBool(enabledVar)
		if err != nil {
			return vmp, errors.Wrapf(err, "cannot parse %s as bool", vmCacheEnableEnv)
		}
		vmp.Enabled = parsedBool
	}

	if err := vmp.UpdateParams(); err != nil {
		return vmp, err
	}

	return vmp, nil
}

// UpdateParams - reads configuration file and updates corresponding flags.
func (vmp *VictoriaMetricsParams) UpdateParams() error {
	if err := vmp.loadVMAlertParams(); err != nil {
		return errors.Wrap(err, "cannot update VMAlertFlags config param")
	}
	if err := vmp.loadVMDBParams(); err != nil {
		return errors.Wrap(err, "cannot update VMDBFlags config param")
	}

	return nil
}

// loadVMAlertParams - load params and converts it to vmalert flags.
func (vmp *VictoriaMetricsParams) loadVMAlertParams() error {
	buf, err := ioutil.ReadFile(vmp.BaseConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "cannot read baseConfigPath for VMAlertParams")
		}

		return nil
	}
	cfg := &config.Config{}
	if err = yaml.Unmarshal(buf, cfg); err != nil {
		return errors.Wrap(err, "cannot unmarshal baseConfigPath for VMAlertFlags")
	}
	vmalertFlags := make([]string, 0, len(vmp.VMAlertFlags))
	for _, r := range cfg.RuleFiles {
		vmalertFlags = append(vmalertFlags, "--rule="+r)
	}
	if cfg.GlobalConfig.EvaluationInterval != 0 {
		vmalertFlags = append(vmalertFlags, "--evaluationInterval="+cfg.GlobalConfig.EvaluationInterval.String())
	}
	vmp.VMAlertFlags = vmalertFlags

	return nil
}

// loadVMDBParams - flags for VictoriaMetrics database can be added here.
func (vmp *VictoriaMetricsParams) loadVMDBParams() error {
	cacheDisabled := true
	if cacheEnable := os.Getenv(vmCacheEnableEnv); cacheEnable != "" {
		parsedBool, err := strconv.ParseBool(os.Getenv(vmCacheEnableEnv))
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s as bool", vmCacheEnableEnv)
		}
		// we have to invert parsed variable
		cacheDisabled = !parsedBool
	}
	vmdbFlags := make([]string, 0, len(vmp.VMDBFlags))
	vmdbFlags = append(vmdbFlags, fmt.Sprintf("--search.disableCache=%t", cacheDisabled))
	vmp.VMDBFlags = vmdbFlags

	return nil
}
