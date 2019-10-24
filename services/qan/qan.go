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

// Package qan contains business logic of working with QAN and qan-agent on PMM Server node.
package qan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	servicelib "github.com/percona/kardianos-service"
	"github.com/percona/pmm/proto"
	"github.com/percona/pmm/proto/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/percona/pmm-managed/models"
	"github.com/percona/pmm-managed/services"
	"github.com/percona/pmm-managed/utils/logger"
)

type Service struct {
	baseDir    string
	supervisor services.Supervisor
	qanAPI     *http.Client
}

func NewService(ctx context.Context, baseDir string, supervisor services.Supervisor) (*Service, error) {
	svc := &Service{
		baseDir:    baseDir,
		supervisor: supervisor,
		qanAPI:     new(http.Client),
	}

	return svc, nil
}

// qanAgentConfigPath returns agent.conf file path.
func (svc *Service) qanAgentConfigPath() string {
	return filepath.Join(svc.baseDir, "config", "agent.conf")
}

// ensureAgentIsRegistered registers a single qan-agent instance on PMM Server node in QAN.
// It does not re-register or change configuration if agent is already registered.
// QAN API URL is always returned when no error is encountered.
func (svc *Service) ensureAgentIsRegistered(ctx context.Context) (*url.URL, error) {
	qanURL, err := getQanURL(ctx)
	if err != nil {
		return nil, err
	}

	l := logger.Get(ctx).WithField("component", "qan")

	// do not change anything if qan-agent is already registered
	path := svc.qanAgentConfigPath()
	if _, err = os.Stat(path); err == nil {
		l.Debugf("qan-agent already registered (%s exists).", path)
		return qanURL, nil
	}

	path = filepath.Join(svc.baseDir, "bin", "percona-qan-agent-installer")
	args := []string{"-debug", "-hostname=pmm-server"}

	if qanURL.User != nil && qanURL.User.Username() != "" {
		args = append(args, "-server-user="+qanURL.User.Username())
		pass, _ := qanURL.User.Password()
		args = append(args, "-server-pass="+pass)
	}

	args = append(args, qanURL.String()) // full URL, with username and password (yes, again! that's how installer is written)
	cmd := exec.Command(path, args...)
	l.Debug(strings.Join(cmd.Args, " "))
	b, err := cmd.CombinedOutput()

	if err != nil {
		l.Infof("%s", b)
		return nil, errors.Wrap(err, "failed to register qan-agent")
	}

	l.Debugf("%s", b)

	// set logging level - very useful for debugging
	logLevel := "info"
	if l.Level == logrus.DebugLevel {
		logLevel = "debug"
	}

	path = filepath.Join(svc.baseDir, "config", "log.conf")

	if err = ioutil.WriteFile(path, []byte(fmt.Sprintf(`{"Level":%q,"Offline":"false"}`, logLevel)), 0666); err != nil {
		return nil, errors.Wrap(err, "failed to write log.conf")
	}

	return qanURL, nil
}

// ensureAgentRuns checks qan-agent process status and starts it if it is not configured or down.
func (svc *Service) ensureAgentRuns(ctx context.Context, nameForSupervisor string, port uint16) error {
	err := svc.supervisor.Status(ctx, nameForSupervisor)
	if err != nil {
		// error can also mean that service status can't be determined, so we always stop it first
		err = svc.supervisor.Stop(ctx, nameForSupervisor)
		if err != nil {
			logger.Get(ctx).WithField("component", "qan").Warn(err)
		}

		config := &servicelib.Config{
			Name:        nameForSupervisor,
			DisplayName: nameForSupervisor,
			Description: nameForSupervisor,
			Executable:  filepath.Join(svc.baseDir, "bin", "percona-qan-agent"),
			Arguments: []string{
				fmt.Sprintf("-listen=127.0.0.1:%d", port),
			},
		}
		err = svc.supervisor.Start(ctx, config)
	}

	return err
}

// Restore ensures that agent is registered and running.
func (svc *Service) Restore(ctx context.Context, nameForSupervisor string, agent models.QanAgent) error {
	l := logger.Get(ctx).WithField("component", "qan")

	qanURL, err := getQanURL(ctx)
	if err != nil {
		l.Infof("getQanURL err: %v", err)
		return err
	}

	agentInstance, dbInstance, err := svc.restoreConfigs(ctx, agent)
	if err != nil {
		l.Infof("restoreConfigs err: %v", err)
		return err
	}

	if err := svc.ensureAgentRuns(ctx, nameForSupervisor, *agent.ListenPort); err != nil {
		return errors.WithStack(err)
	}

	command := "StartTool"
	config := map[string]interface{}{
		"UUID":           dbInstance.UUID,
		"CollectFrom":    "perfschema",
		"Interval":       60,
		"ExampleQueries": true,
	}

	b, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	l.Infof("%s %s %s", agentInstance.UUID, command, b)

	return svc.sendQANCommand(ctx, qanURL, agentInstance.UUID, command, b)
}

func (svc *Service) restoreConfigs(ctx context.Context, agent models.QanAgent) (*proto.Instance, *proto.Instance, error) {
	l := logger.Get(ctx).WithField("component", "qan")

	qanURL, err := getQanURL(ctx)
	if err != nil {
		l.Infof("getQanURL err: %v", err)
		return nil, nil, err
	}

	// restore mysql instance.
	instances, err := svc.getInstances(ctx, qanURL)
	if err != nil {
		l.Infof("getInstances err: %v", err)
		return nil, nil, errors.Wrap(err, "failed to get MySQL instance by UUID")
	}

	var (
		osInstance    proto.Instance
		agentInstance proto.Instance
		dbInstance    proto.Instance
	)

	// look for mysql instance
	for _, inst := range instances {
		if inst.UUID == *agent.QANDBInstanceUUID {
			dbInstance = inst
			break
		}
	}

	// look for related agent and os instances
	for _, inst := range instances {
		if inst.UUID == dbInstance.ParentUUID && inst.Subsystem == "os" {
			osInstance = inst
		}

		if inst.ParentUUID == dbInstance.ParentUUID && inst.Subsystem == "agent" {
			agentInstance = inst
		}
	}

	path := filepath.Join(svc.baseDir, "instance")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0750)
	}

	// restore db instance.
	path = filepath.Join(svc.baseDir, "instance", fmt.Sprintf("%s.json", dbInstance.UUID))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dbInstance.DSN = strings.Replace(dbInstance.DSN, "***", *agent.ServicePassword, 1)
		dbInstance.DSN = fmt.Sprintf("%s/?timeout=5s", dbInstance.DSN)
		dbInstanceJSON, err := json.MarshalIndent(dbInstance, "", "    ")

		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		if err = ioutil.WriteFile(path, dbInstanceJSON, 0666); err != nil {
			return nil, nil, errors.Wrapf(err, "failed to write %s", path)
		}

		l.Infof("restored dbInstance: %s.", path)
	}

	// restore OS instance.
	path = filepath.Join(svc.baseDir, "instance", fmt.Sprintf("%s.json", osInstance.UUID))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		osInstanceJSON, err := json.MarshalIndent(osInstance, "", "    ")

		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		if err = ioutil.WriteFile(path, osInstanceJSON, 0666); err != nil {
			return nil, nil, errors.Wrapf(err, "failed to write %s", path)
		}

		l.Infof("restored osInstance: %s.", path)
	}

	// restore Agent instance.
	path = filepath.Join(svc.baseDir, "instance", fmt.Sprintf("%s.json", agentInstance.UUID))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		agentInstanceJSON, err := json.MarshalIndent(agentInstance, "", "    ")

		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		if err = ioutil.WriteFile(path, agentInstanceJSON, 0666); err != nil {
			return nil, nil, errors.Wrapf(err, "failed to write %s", path)
		}

		l.Infof("restored agentInstance: %s.", path)
	}

	path = filepath.Join(svc.baseDir, "config")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0750)
	}

	path = filepath.Join(svc.baseDir, "config", "agent.conf")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		serverUser := "pmm"
		if os.Getenv("SERVER_USER") != "" {
			serverUser = os.Getenv("SERVER_USER")
		}

		agentConf := struct {
			UUID           string `json:"UUID"`
			APIHostname    string `json:"ApiHostname"`
			APIPath        string `json:"ApiPath"`
			ServerUser     string `json:"ServerUser"`
			ServerPassword string `json:"ServerPassword,omitempty"`
		}{
			agentInstance.UUID,
			"127.0.0.1",
			"/qan-api/",
			serverUser,
			os.Getenv("SERVER_PASSWORD"),
		}

		// agentConf := fmt.Sprintf(`{"UUID":"%s","ApiHostname":"127.0.0.1","ApiPath":"/qan-api/","ServerUser":"pmm"}`, agentInstance.UUID)
		b, err := json.Marshal(agentConf)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to Marshal agent.conf")
		}

		if err = ioutil.WriteFile(path, b, 0666); err != nil {
			return nil, nil, errors.Wrap(err, "failed to write agent.conf")
		}

		l.Infof("restored agent config: %s.", path)
	}

	path = filepath.Join(svc.baseDir, "config", "log.conf")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = ioutil.WriteFile(path, []byte(`{"Level":"info","Offline":"false"}`), 0666); err != nil {
			return nil, nil, errors.Wrap(err, "failed to write agent.conf")
		}

		l.Infof("restored log config: %s.", path)
	}

	path = filepath.Join(svc.baseDir, "config", fmt.Sprintf("qan-%s.conf", dbInstance.UUID))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		qanConf := fmt.Sprintf(`{ "UUID": "%s", "CollectFrom": "perfschema", "Interval": 60, "ExampleQueries": true }`, dbInstance.UUID)

		if err = ioutil.WriteFile(path, []byte(qanConf), 0666); err != nil {
			return nil, nil, errors.Wrap(err, "failed to write agent.conf")
		}

		l.Infof("restored qan config: %s.", path)
	}

	return &agentInstance, &dbInstance, nil
}

// getInstances returns all instances from the QAN API.
func (svc *Service) getInstances(ctx context.Context, qanURL *url.URL) ([]proto.Instance, error) {
	url := *qanURL
	// url.Path = path.Join(url.Path, "instances", UUID)
	url.Path = path.Join(url.Path, "instances")
	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	rb, _ := httputil.DumpRequestOut(req, true)
	logger.Get(ctx).WithField("component", "qan").Debugf("UUID request:\n\n%s\n", rb)

	resp, err := svc.qanAPI.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	rb, _ = httputil.DumpResponse(resp, true)

	if resp.StatusCode != 200 {
		logger.Get(ctx).WithField("component", "qan").Errorf("UUID response:\n\n%s\n", rb)
		return nil, errors.Errorf("unexpected QAN response status code %d", resp.StatusCode)
	}

	logger.Get(ctx).WithField("component", "qan").Debugf("UUID response:\n\n%s\n", rb)

	var instances []proto.Instance
	if err = json.NewDecoder(resp.Body).Decode(&instances); err != nil {
		return nil, errors.WithStack(err)
	}

	return instances, nil
}

// getAgentUUID returns agent UUID from the qan-agent configuration file.
func (svc *Service) getAgentUUID() (string, error) {
	path := svc.qanAgentConfigPath()
	f, err := os.Open(path)

	if err != nil {
		return "", errors.WithStack(err)
	}
	defer f.Close()

	var cfg config.Agent
	if err = json.NewDecoder(f).Decode(&cfg); err != nil {
		return "", errors.WithStack(err)
	}

	if cfg.UUID == "" {
		err = errors.Errorf("missing agent UUID in configuration file %s", path)
	}

	return cfg.UUID, err
}

// getOSUUID returns OS UUID from the QAN API.
func (svc *Service) getOSUUID(ctx context.Context, qanURL *url.URL, agentUUID string) (string, error) {
	url := *qanURL
	url.Path = path.Join(url.Path, "instances", agentUUID)
	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return "", errors.WithStack(err)
	}

	rb, _ := httputil.DumpRequestOut(req, true)
	logger.Get(ctx).WithField("component", "qan").Debugf("getOSUUID request:\n\n%s\n", rb)

	resp, err := svc.qanAPI.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer resp.Body.Close()

	rb, _ = httputil.DumpResponse(resp, true)

	if resp.StatusCode != 200 {
		logger.Get(ctx).WithField("component", "qan").Errorf("getOSUUID response:\n\n%s\n", rb)
		return "", errors.Errorf("unexpected QAN response status code %d", resp.StatusCode)
	}

	logger.Get(ctx).WithField("component", "qan").Debugf("getOSUUID response:\n\n%s\n", rb)

	var instance proto.Instance
	if err = json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return "", errors.WithStack(err)
	}

	return instance.ParentUUID, nil
}

// addInstanceToServer adds instance to QAN API.
// If successful, instance's UUID field will be set.
func (svc *Service) addInstanceToServer(ctx context.Context, qanURL *url.URL, instance *proto.Instance) error {
	b, err := json.Marshal(instance)
	if err != nil {
		return errors.WithStack(err)
	}

	url := *qanURL
	url.Path = path.Join(url.Path, "instances")
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(b))

	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rb, _ := httputil.DumpRequestOut(req, true)
	logger.Get(ctx).WithField("component", "qan").Debugf("addInstanceToServer request:\n\n%s\n", rb)

	resp, err := svc.qanAPI.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	rb, _ = httputil.DumpResponse(resp, true)

	if resp.StatusCode != 201 {
		logger.Get(ctx).WithField("component", "qan").Errorf("addInstanceToServer response:\n\n%s\n", rb)
		return errors.Errorf("unexpected QAN response status code %d", resp.StatusCode)
	}

	logger.Get(ctx).WithField("component", "qan").Debugf("addInstanceToServer response:\n\n%s\n", rb)

	// Response Location header looks like this: http://127.0.0.1/qan-api/instances/6cea8824082d4ade682b94109664e6a9
	// Extract UUID directly from it instead of following it.
	parts := strings.Split(resp.Header.Get("Location"), "/")
	instance.UUID = parts[len(parts)-1]

	return nil
}

// removeInstanceFromServer removes instance from QAN API.
func (svc *Service) removeInstanceFromServer(ctx context.Context, qanURL *url.URL, uuid string) error {
	url := *qanURL
	url.Path = path.Join(url.Path, "instances", uuid)
	req, err := http.NewRequest("DELETE", url.String(), nil)

	if err != nil {
		return errors.WithStack(err)
	}

	rb, _ := httputil.DumpRequestOut(req, true)
	logger.Get(ctx).WithField("component", "qan").Debugf("removeInstanceFromServer request:\n\n%s\n", rb)

	resp, err := svc.qanAPI.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	rb, _ = httputil.DumpResponse(resp, true)

	if resp.StatusCode != 204 {
		logger.Get(ctx).WithField("component", "qan").Errorf("removeInstanceFromServer response:\n\n%s\n", rb)
		return errors.Errorf("unexpected QAN response status code %d", resp.StatusCode)
	}

	logger.Get(ctx).WithField("component", "qan").Debugf("removeInstanceFromServer response:\n\n%s\n", rb)

	return nil
}

func (svc *Service) sendQANCommand(ctx context.Context, qanURL *url.URL, agentUUID string, command string, data []byte) error {
	cmd := proto.Cmd{
		User:      "pmm-managed",
		AgentUUID: agentUUID,
		Service:   "qan",
		Cmd:       command,
		Data:      data,
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		return errors.WithStack(err)
	}

	// Send the command to the API which relays it to the agent, then relays the agent's reply back to here.
	// It takes a few seconds for agent to connect to QAN API once it is started via service manager.
	// QAN API fails to start/stop unconnected agent for QAN, so we retry the request when getting 404 response.
	const attempts = 10

	url := *qanURL
	url.Path = path.Join(url.Path, "agents", agentUUID, "cmd")

	for i := 0; i < attempts; i++ {
		req, err := http.NewRequest("PUT", url.String(), bytes.NewReader(b))
		if err != nil {
			return errors.WithStack(err)
		}

		req.Header.Set("Content-Type", "application/json")
		rb, _ := httputil.DumpRequestOut(req, true)
		logger.Get(ctx).WithField("component", "qan").Debugf("sendQANCommand request:\n\n%s\n", rb)

		resp, err := svc.qanAPI.Do(req)
		if err != nil {
			return errors.WithStack(err)
		}

		rb, _ = httputil.DumpResponse(resp, true)
		resp.Body.Close()

		if resp.StatusCode == 200 {
			logger.Get(ctx).WithField("component", "qan").Debugf("sendQANCommand response:\n\n%s\n", rb)
			return nil
		}

		if resp.StatusCode == 404 {
			logger.Get(ctx).WithField("component", "qan").Debugf("sendQANCommand response:\n\n%s\n", rb)
			time.Sleep(time.Second)

			continue
		}

		logger.Get(ctx).WithField("component", "qan").Errorf("sendQANCommand response:\n\n%s\n", rb)

		return errors.Errorf("%s: unexpected QAN API response status code %d", command, resp.StatusCode)
	}

	return errors.Errorf("%s: failed to send command after %d attempts", command, attempts)
}

// AddMySQL adds MySQL instance to QAN, configuring and enabling it.
// It sets MySQL instance UUID to qanAgent.QANDBInstanceUUID.
func (svc *Service) AddMySQL(ctx context.Context, nodeName string, mySQLService *models.MySQLService, qanAgent *models.QanAgent) error {
	qanURL, err := svc.ensureAgentIsRegistered(ctx)
	if err != nil {
		return err
	}

	agentUUID, err := svc.getAgentUUID()
	if err != nil {
		return err
	}

	osUUID, err := svc.getOSUUID(ctx, qanURL, agentUUID)
	if err != nil {
		return err
	}

	instance := &proto.Instance{
		Subsystem:  "mysql",
		ParentUUID: osUUID,
		Name:       nodeName,
		DSN:        sanitizeDSN(qanAgent.DSN(mySQLService)),
		Version:    *mySQLService.EngineVersion,
	}
	if err = svc.addInstanceToServer(ctx, qanURL, instance); err != nil {
		return err
	}

	qanAgent.QANDBInstanceUUID = pointer.ToString(instance.UUID)

	// we need real DSN (with password) for qan-agent to work, and it seems to be the only way to pass it
	path := filepath.Join(svc.baseDir, "instance", fmt.Sprintf("%s.json", instance.UUID))
	instance.DSN = qanAgent.DSN(mySQLService)

	b, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return errors.WithStack(err)
	}

	if err = ioutil.WriteFile(path, b, 0666); err != nil {
		return errors.WithStack(err)
	}

	if err = svc.ensureAgentRuns(ctx, models.NameForSupervisor(qanAgent.Type, *qanAgent.ListenPort), *qanAgent.ListenPort); err != nil {
		return err
	}

	command := "StartTool"
	config := map[string]interface{}{
		"UUID":           instance.UUID,
		"CollectFrom":    "perfschema",
		"Interval":       60,
		"ExampleQueries": true,
	}

	b, err = json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	logger.Get(ctx).WithField("component", "qan").Debugf("%s %s %s", agentUUID, command, b)

	return svc.sendQANCommand(ctx, qanURL, agentUUID, command, b)
}

func (svc *Service) RemoveMySQL(ctx context.Context, qanAgent *models.QanAgent) error {
	qanURL, err := svc.ensureAgentIsRegistered(ctx)
	if err != nil {
		return err
	}

	// agent should be running to remove instance from it
	if err = svc.ensureAgentRuns(ctx, models.NameForSupervisor(qanAgent.Type, *qanAgent.ListenPort), *qanAgent.ListenPort); err != nil {
		return err
	}

	agentUUID, err := svc.getAgentUUID()
	if err != nil {
		return err
	}

	command := "StopTool"
	b := []byte(*qanAgent.QANDBInstanceUUID)
	logger.Get(ctx).WithField("component", "qan").Debugf("%s %s %s", agentUUID, command, b)

	if err = svc.sendQANCommand(ctx, qanURL, agentUUID, command, b); err != nil {
		return err
	}

	// we do not stop qan-agent even if it has zero MySQL instances now - to be safe
	return svc.removeInstanceFromServer(ctx, qanURL, *qanAgent.QANDBInstanceUUID)
}
