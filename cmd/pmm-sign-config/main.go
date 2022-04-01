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

func main() {
	logger.SetupGlobalLogger()
	var signKey, signPub, configPath string
	if signKey = os.Getenv("PMM_CONFIG_SIGN_KEY"); signKey == "" {
		signKey = "./.cfg.dev-sign.key"
	}
	if signPub = os.Getenv("PMM_CONFIG_SIGN_PUB"); signPub == "" {
		signPub = "./.cfg.dev-sign.pub"
	}
	if configPath = os.Getenv("PMM_CONFIG_SIGN_PUB"); configPath == "" {
		configPath = "./config/telemetry/dev"
	}

	private, err := minisign.PrivateKeyFromFile("", signKey)
	if err != nil {
		log.Fatalln(err)
	}

	configRegEx, e := regexp.Compile("^.+\\.yml$")
	if e != nil {
		log.Fatal(e)
	}

	e = filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && configRegEx.MatchString(info.Name()) {
			configContent, err := ioutil.ReadFile(path) // the file is inside the local directory
			if err != nil {
				log.Fatal(err)
			}

			configSignature := minisign.Sign(private, configContent)
			if err = ioutil.WriteFile(path+".minisig", configSignature, 0644); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}
