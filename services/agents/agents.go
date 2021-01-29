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
	"fmt"
	"strings"

	"github.com/AlekSi/pointer"

	"github.com/percona/pmm-managed/models"
)

type redactMode int

const (
	redactSecrets redactMode = iota
	exposeSecrets
)

// redactWords returns words that should be redacted from given Agent logs/output.
func redactWords(agent *models.Agent) []string {
	var words []string
	if s := pointer.GetString(agent.Password); s != "" {
		words = append(words, s)
	}
	if s := pointer.GetString(agent.AWSSecretKey); s != "" {
		words = append(words, s)
	}
	return words
}

// FilterOutCollectors removes from exporter's flags disabled collectors.
// DisableCollector values should  match collector flag till end of string or till `=` sign.
// Examples:
// 1. if we pass `meminfo` then only "--collector.meminfo" but not "--collector.meminfo_numa"
// 2. if we pass `netstat.field` then "--collector.netstat.fields=^(.*_(InErrors|InErrs|InCsumErrors)..." shold be disabled.
// 3. To disable "--collect.custom_query.hr" with directory ""--collect.custom_query.lr.directory" user should pass both names.
func FilterOutCollectors(prefix string, args, disabledCollectors []string) []string {
	argsMap := make(map[string]string)
	for _, arg := range args {
		flagName := strings.Split(arg, "=")[0]
		argsMap[flagName] = arg
	}

	for _, disabledCollector := range disabledCollectors {
		key := fmt.Sprintf("%s%s", prefix, disabledCollector)
		_, ok := argsMap[key]
		if ok {
			delete(argsMap, key)
		}
	}

	enabledArgs := []string{}
	for _, arg := range argsMap {
		enabledArgs = append(enabledArgs, arg)
	}

	return enabledArgs
}
