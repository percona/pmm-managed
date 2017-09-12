package telemetry

import (
	"bytes"
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
	ticker     *time.Ticker
	lastStatus string
	stopChan   chan bool
	wg         sync.WaitGroup
	isRunning  bool
	lock       sync.Mutex
}

// Telemetry config
type Config struct {
	URL      string `yaml:"url'`
	UUID     string `yaml:"uuid"`
	Interval int    `yaml:"interval"`
}

const (
	defaultURL      = "https://v.percona.com"
	defaultUUID     = ""
	defaultinterval = 24 * 60 * 60
)

var (
	stat     = os.Stat
	readFile = ioutil.ReadFile
	output   = commandOutput
	Version  = "1.3.0"
)

// NewService creates a new telemetry service given a configuration
func NewService(config *Config) (*Service, error) {
	service := &Service{
		Config: *config,
	}
	return service, nil
}

// Start sets a new ticker to collect data every Config.Interval seconds.
func (s *Service) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.isRunning {
		return
	}
	s.stopChan = make(chan bool)
	s.isRunning = true
	s.ticker = time.NewTicker(time.Second * time.Duration(s.Interval))
	s.wg.Add(1)
	go func() {
	LOOP:
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
			case <-s.stopChan:
				break LOOP
			}
		}
		s.wg.Done()
	}()
}

// Stop the service
func (s *Service) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.isRunning {
		return
	}
	s.ticker.Stop()
	close(s.stopChan)
	s.wg.Wait()
	s.isRunning = false
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

	data["PMM"] = Version

	return data, nil
}

func (s *Service) makePayload(data map[string]interface{}) ([]byte, error) {
	w := bytes.NewBuffer(nil)

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
	req.Header.Add("X-Percona-Toolkit-Tool", "pmm")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error while sending telemetry data: Status: %d, %s", resp.StatusCode, resp.Status)
	}
	return nil
}
