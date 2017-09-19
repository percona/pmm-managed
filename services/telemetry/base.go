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
	"sync"
	"time"
)

// Telemetry exported and unexported fields
type Service struct {
	Config
	PMMVersion string
	lastStatus string
	ticker     *time.Ticker
	wg         sync.WaitGroup

	lock      sync.Mutex
	isRunning bool
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

// Start sets a new ticker to collect data every Config.Interval seconds.
func (s *Service) Start(ctx context.Context) {
	if s.IsRunning() {
		return
	}
	s.setRunning(true)
	s.wg.Add(1)
	s.ticker = time.NewTicker(s.Interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				data, err := s.collectData()
				payload, err := s.makePayload(data)
				if err != nil {
					s.lastStatus = fmt.Sprintf("%v cannot build payload for telemetry info: %s", time.Now(), err.Error())
					continue
				}
				err = s.sendRequest(payload)
				if err != nil {
					s.lastStatus = fmt.Sprintf("%v error sending telemetry info: %s", time.Now(), err.Error())
					continue
				}
				s.lastStatus = fmt.Sprintf("%v telemetry data sent.", time.Now())
			case <-ctx.Done():
				s.setRunning(false)
				s.wg.Done()
				return
			}
		}
	}()
}

// Wait for the service to stop
func (s *Service) Wait() {
	s.wg.Wait()
}

func (s *Service) setRunning(status bool) {
	s.lock.Lock()
	s.isRunning = status
	s.lock.Unlock()
}

// IsRunning returns true if the service is running
func (s *Service) IsRunning() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.isRunning
}

// GetLastStatus returns a string having the timestamp and result
// of the last call to the telemetry server
func (s *Service) GetLastStatus() string {
	return s.lastStatus
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
