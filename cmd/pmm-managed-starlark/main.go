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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/percona-platform/saas/pkg/starlark"
	"github.com/percona/pmm-managed/services/checks"

	"github.com/percona/pmm/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("stdlog: ")

	kingpin.Version(version.FullInfo())
	kingpin.HelpFlag.Short('h')

	kingpin.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000-07:00",
	})
	if on, _ := strconv.ParseBool(os.Getenv("PMM_DEBUG")); on {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if on, _ := strconv.ParseBool(os.Getenv("PMM_TRACE")); on {
		logrus.SetLevel(logrus.TraceLevel)
	}

	l := logrus.WithField("component", "pmm-managed-starlark")

	decoder := json.NewDecoder(os.Stdin)
	var data checks.StarlarkScriptData
	err := decoder.Decode(&data)
	if err != nil {
		l.Error("Error decoding json data: ", err)
		os.Exit(1)
	}

	funcs, err := checks.GetFuncsForVersion(data.CheckVersion)
	if err != nil {
		l.Error("Error getting funcs: ", err)
		os.Exit(1)
	}

	env, err := starlark.NewEnv(data.CheckName, data.Script, funcs)
	if err != nil {
		l.Error("Error initializing starlark env: ", err)
		os.Exit(1)
	}

	results, err := env.Run(data.CheckName, data.ScriptInput, l.Debugln)
	if err != nil {
		l.Error("Error running starlark env: ", err)
		os.Exit(1)
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		l.Error("Error marshalling JSON: ", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonResults))
}
