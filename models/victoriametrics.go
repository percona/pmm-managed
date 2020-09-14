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
	"reflect"
	"strconv"

	config "github.com/percona/promconfig"
	"gopkg.in/yaml.v2"
)

const (
	vmCacheDisableEnv = "VM_CACHE_DISABLE"
	vmTestEnableEnv   = "PERCONA_TEST_VM"
)

// VictoriaMetricsParams - defines flags and settings for victoriametrics.
type VictoriaMetricsParams struct {
	VMDB           []string
	VMAlert        []string
	Enabled        bool
	baseConfigPath string
}

// NewVictoriaMetricsParams - returns configuration params for VictoriaMetrics
// error can be suppressed, this func always return inited struct with default
// params.
func NewVictoriaMetricsParams(basePath string) (*VictoriaMetricsParams, error) {
	vmp := &VictoriaMetricsParams{
		VMDB:           []string{},
		VMAlert:        []string{},
		baseConfigPath: basePath,
	}
	if enabledVar := os.Getenv(vmTestEnableEnv); enabledVar != "" {
		parsedBool, err := strconv.ParseBool(enabledVar)
		if err != nil {
			return vmp, fmt.Errorf("cannot parse PERCONA_TEST_VM, as bool, value: %s, err: %w", enabledVar, err)
		}
		vmp.Enabled = parsedBool
	}

	err := vmp.UpdateParams()
	if err != nil {
		return vmp, err
	}

	return vmp, nil
}

// UpdateParams - updates params for VictoriaMetrics services
// reads configuration file and updates corresponding flags.
func (vmp *VictoriaMetricsParams) UpdateParams() error {
	if err := vmp.loadVMAlertParams(); err != nil {
		return err
	}
	if err := vmp.loadVMDBParams(); err != nil {
		return err
	}

	return nil
}

// loadVMAlertParams - load params and converts it to vmalert flags.
func (vmp *VictoriaMetricsParams) loadVMAlertParams() error {
	buf, err := ioutil.ReadFile(vmp.baseConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		return nil
	}
	cfg := &config.Config{}
	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		return err
	}
	vmalertFlags := make([]string, 0, len(vmp.VMAlert))
	for _, r := range cfg.RuleFiles {
		vmalertFlags = append(vmalertFlags, "--rule="+r)
	}
	if cfg.GlobalConfig.EvaluationInterval != 0 {
		vmalertFlags = append(vmalertFlags, "--evaluationInterval="+cfg.GlobalConfig.EvaluationInterval.String())
	}

	if !reflect.DeepEqual(vmalertFlags, vmp.VMAlert) {
		vmp.VMAlert = vmalertFlags
	}

	return nil
}

// loadVMDBParams - flags for VictoriaMetrics database can be added here.
func (vmp *VictoriaMetricsParams) loadVMDBParams() error {
	cacheDisabled := os.Getenv(vmCacheDisableEnv)
	if cacheDisabled == "" {
		return nil
	}
	parsedBool, err := strconv.ParseBool(cacheDisabled)
	if err != nil {
		return err
	}
	vmdbFlags := make([]string, 0, len(vmp.VMDB))
	vmdbFlags = append(vmdbFlags, fmt.Sprintf("--search.disableCache=%t", parsedBool))
	if !reflect.DeepEqual(vmdbFlags, vmp.VMDB) {
		vmp.VMDB = vmdbFlags
	}

	return nil
}
