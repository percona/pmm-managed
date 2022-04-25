package agents

import (
	"github.com/AlekSi/pointer"
	"github.com/percona/pmm/version"
)

// Log level available in exporters with pmm 2.28
var exporterLogLevelCommandVersion = version.MustParse("2.28.0")

func withLogLevel(args []string, logLevel *string, pmmAgentVersion *version.Parsed) []string {
	if pointer.GetString(logLevel) != "" && !pmmAgentVersion.Less(exporterLogLevelCommandVersion) {
		args = append(args, "-log.level="+*logLevel)
	}

	return args
}
