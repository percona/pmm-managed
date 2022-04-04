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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"aead.dev/minisign"

	"github.com/percona/pmm-managed/utils/logger"
)

func main() { //nolint:cyclop
	logger.SetupGlobalLogger()
	var signKey, signPub, configPath string
	if signKey = os.Getenv("PMM_CONFIG_SIGN_KEY"); signKey == "" {
		signKey = "./.cfg.dev-sign.key"
	}
	if signPub = os.Getenv("PMM_CONFIG_SIGN_PUB"); signPub == "" {
		signPub = "./.cfg.dev-sign.pub" //nolint:ineffassign,wastedassign
	}
	if configPath = os.Getenv("PMM_CONFIG_SIGN_PUB"); configPath == "" {
		configPath = "./config/telemetry/dev"
	}

	private, err := minisign.PrivateKeyFromFile("", signKey)
	if err != nil {
		log.Fatalln(err)
	}

	configRegEx, e := regexp.Compile("^.+\\.yml$") //nolint:gocritic,gosimple
	if e != nil {
		log.Fatal(e)
	}

	e = filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && configRegEx.MatchString(info.Name()) {
			// the file is inside the local directory
			configContent, err := ioutil.ReadFile(path) //nolint:gosec
			if err != nil {
				log.Fatal(err)
			}

			configSignature := minisign.Sign(private, configContent)
			// TODO: review if we can use 0600 or less
			if err = ioutil.WriteFile(path+".minisig", configSignature, 0o644); err != nil { //nolint:gosec
				log.Fatal(err)
			}
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}
