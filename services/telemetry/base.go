// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Package telemetry provides Call Home functionality.
package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Telemetry exported and unexported fields
type Service struct {
	Config
	PMMVersion string
}

// Telemetry config
type Config struct {
	URL      string        `yaml:"url"`
	UUID     string        `yaml:"uuid"`
	Interval time.Duration `yaml:"interval"`
}

const (
	defaultURL      = "https://v.percona.com"
	defaultinterval = 24 * 60 * 60 * time.Second
)

var (
	stat     = os.Stat
	readFile = ioutil.ReadFile
	output   = commandOutput
)

func DefaultConfig() *Config {
	uuid, err := generateUUID()
	if err != nil {
		uuid = ""
	}
	return &Config{
		URL:      defaultURL,
		Interval: defaultinterval,
		UUID:     uuid,
	}
}

// NewService creates a new telemetry service given a configuration
func NewService(config *Config, pmmVersion string) (*Service, error) {
	service := &Service{
		Config:     *config,
		PMMVersion: pmmVersion,
	}
	return service, nil
}

// Run runs telemetry service, sending data every Config.Interval until context is canceled.
func (s *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		s.runOnce()

		select {
		case <-ticker.C:
			// continue with next loop iteration
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) runOnce() string {
	data, err := s.collectData()
	payload, err := s.makePayload(data)
	if err != nil {
		return fmt.Sprintf("%s cannot build payload for telemetry info: %s", time.Now(), err)
	}
	err = s.sendRequest(payload)
	if err != nil {
		return fmt.Sprintf("%s error sending telemetry info: %s", time.Now(), err)
	}
	return fmt.Sprintf("%s telemetry data sent.", time.Now())
}

func (s *Service) collectData() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	if osType, err := getOSNameAndVersion(); err == nil {
		data["OS"] = osType
	}

	data["PMM"] = s.PMMVersion

	return data, nil
}

func (s *Service) makePayload(data map[string]interface{}) ([]byte, error) {
	var w bytes.Buffer

	for key, value := range data {
		w.WriteString(fmt.Sprintf("%s;%s;%s\n", s.Config.UUID, key, value))
	}

	return w.Bytes(), nil
}

func (s *Service) sendRequest(data []byte) error {
	client := &http.Client{}
	body := bytes.NewReader(data)

	req, err := http.NewRequest("POST", s.Config.URL, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "plain/text")
	req.Header.Add("X-Percona-Toolkit-Tool", "pmm")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error while sending telemetry data: Status: %d, %s", resp.StatusCode, resp.Status)
	}
	return nil
}
