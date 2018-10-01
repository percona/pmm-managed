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

package logs

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	servicelib "github.com/percona/kardianos-service"
	"gopkg.in/yaml.v2"

	"github.com/percona/pmm-managed/utils/logger"
)

// File represents log file content.
type File struct {
	Name string
	Data []byte
	Err  error
}

type Log struct {
	FilePath string
	UnitName string
	Command  string
}

var DefaultLogs = []Log{
	{"/var/log/consul.log", "consul", ""},
	{"/var/log/createdb.log", "", ""},
	{"/var/log/cron.log", "crond", ""},
	{"/var/log/dashboard-upgrade.log", "", ""},
	{"/var/log/grafana/grafana.log", "", ""},
	{"/var/log/mysql.log", "", ""},
	{"/var/log/mysqld.log", "mysqld", ""},
	{"/var/log/nginx.log", "nginx", ""},
	{"/var/log/nginx/access.log", "", ""},
	{"/var/log/nginx/error.log", "", ""},
	{"/var/log/node_exporter.log", "node_exporter", ""},
	{"/var/log/orchestrator.log", "orchestrator", ""},
	{"/var/log/pmm-manage.log", "pmm-manage", ""},
	{"/var/log/pmm-managed.log", "pmm-managed", ""},
	{"/var/log/prometheus1.log", "prometheus1", ""},
	{"/var/log/prometheus.log", "prometheus", ""},
	{"/var/log/qan-api.log", "percona-qan-api", ""},
	{"/var/log/supervisor/supervisord.log", "", ""},
	{"/etc/prometheus.yml", "cat", ""},
	{"/etc/supervisord.d/pmm.ini", "cat", ""},
	{"/etc/nginx/conf.d/pmm.conf", "cat", ""},
	{"prometheus_targets.html", "http", "http://localhost/prometheus/targets"},
	{"consul_nodes.json", "http", "http://localhost/v1/internal/ui/nodes?dc=dc1"},
	{"qan-api_instances.json", "http", "http://localhost/qan-api/instances"},
	{"managed_RDS-Aurora.json", "http", "http://localhost/managed/v0/rds"},
	{"pmm-version.txt", "pmmVersion", ""},
	{"supervisorctl_status.log", "exec", "supervisorctl status"},
	{"systemctl_status.log", "exec", "systemctl -l status"},
	{"pt-summary.log", "exec", "pt-summary"},
}

// Logs is responsible for interactions with logs.
type Logs struct {
	n              int
	logs           []Log
	journalctlPath string
	ctx            context.Context
}

type manageConfig struct {
	Users []struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"users"`
}

// PMM version
var Version string

// getCredential fetchs PMM credential
func getCredential(ctx context.Context) (string, error) {
	var u string
	f, err := os.Open("/srv/update/pmm-manage.yml")
	if err != nil {
		return u, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return u, err
	}

	var config manageConfig
	if err = yaml.Unmarshal(b, &config); err != nil {
		return u, err
	}
	if len(config.Users) > 0 && config.Users[0].Username != "" {
		u = strings.Join([]string{config.Users[0].Username, config.Users[0].Password}, ":")
	}
	return u, err
}

// New creates a new Logs service.
// n is a number of last lines of log to read.
func New(ctx context.Context, pmmVersion string, logs []Log, n int) *Logs {
	l := &Logs{
		n:    n,
		logs: logs,
		ctx:  ctx,
	}

	Version = pmmVersion

	// PMM Server Docker image contails journalctl, so we can't use exec.LookPath("journalctl") alone for detection.
	// TODO Probably, that check should be moved to supervisor service.
	//      Or the whole logs service should be merged with it.
	if servicelib.Platform() == "linux-systemd" {
		l.journalctlPath, _ = exec.LookPath("journalctl")
	}

	return l
}

// Zip creates .zip archive with all logs.
func (l *Logs) Zip(ctx context.Context, w io.Writer) error {
	zw := zip.NewWriter(w)

	now := time.Now().UTC()
	for _, log := range l.logs {
		name, content, err := l.readLog(ctx, &log)
		if name == "" {
			continue
		}

		if err != nil {
			logger.Get(l.ctx).WithField("component", "logs").Error(err)

			// do not let a single error break the whole archive
			if len(content) > 0 {
				content = append(content, "\n\n"...)
			}
			content = append(content, []byte(err.Error())...)
		}

		f, err := zw.CreateHeader(&zip.FileHeader{
			Name:     name,
			Method:   zip.Deflate,
			Modified: now,
		})
		if err != nil {
			return err
		}
		if _, err = f.Write(content); err != nil {
			return err
		}
	}

	// make sure to check the error on Close
	return zw.Close()
}

// Files returns list of logs and their content.
func (l *Logs) Files(ctx context.Context) []File {
	files := make([]File, len(l.logs))

	for i, log := range l.logs {
		var file File
		file.Name, file.Data, file.Err = l.readLog(ctx, &log)
		files[i] = file
	}

	return files
}

// readLog reads last l.n lines from defined Log configuration.
func (l *Logs) readLog(ctx context.Context, log *Log) (name string, data []byte, err error) {
	if log.UnitName == "exec" {
		name = filepath.Base(log.FilePath)
		data, err = l.collectExec(ctx, log.FilePath, log.Command)
		return name, data, err
	}

	if log.UnitName == "pmmVersion" {
		name = filepath.Base(log.FilePath)
		data = []byte(Version)
		return
	}

	if log.UnitName == "http" {
		s := strings.Split(log.Command, "//")
		credential, err1 := getCredential(ctx)
		if len(s) > 1 && len(credential) > 1 {
			log.Command = fmt.Sprintf("%s//%s@%s", s[0], credential, s[1])
		}
		name = filepath.Base(log.FilePath)
		data, err := l.readURL(log.Command)
		if err1 != nil {
			return name, data, fmt.Errorf("%v; %v", err1, err)
		}
		return name, data, err
	}

	if log.UnitName == "cat" {
		name = filepath.Base(log.FilePath)
		data, err = l.readFile(log.FilePath)
		return
	}

	if log.UnitName != "" && l.journalctlPath != "" {
		name = log.UnitName
		data, err = l.journalctlN(ctx, log.UnitName)
		return
	}

	if log.FilePath != "" {
		name = filepath.Base(log.FilePath)
		data, err = l.tailN(ctx, log.FilePath)
		return
	}

	return
}

// journalctlN reads last l.n lines from systemd unit u using `journalctl` command.
func (l *Logs) journalctlN(ctx context.Context, u string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, l.journalctlPath, "-n", strconv.Itoa(l.n), "-u", u)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	b, err := cmd.Output()
	if err != nil {
		return b, fmt.Errorf("%s: %s: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}
	return b, nil
}

// tailN reads last l.n lines from log file at given path using `tail` command.
func (l *Logs) tailN(ctx context.Context, path string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "/usr/bin/tail", "-n", strconv.Itoa(l.n), path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	b, err := cmd.Output()
	if err != nil {
		return b, fmt.Errorf("%s: %s: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}
	return b, nil
}

// collectExec collects output from various commands
func (l *Logs) collectExec(ctx context.Context, path string, command string) ([]byte, error) {
	cmd := &exec.Cmd{}
	if filepath.Dir(path) != "." {
		cmd = exec.CommandContext(ctx, command, path)
	} else {
		command := strings.Split(command, " ")
		cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	}
	var stderr bytes.Buffer
	cmd.Stderr = new(bytes.Buffer)
	b, err := cmd.Output()
	if err != nil {
		return b, fmt.Errorf("%s: %s: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}
	return b, nil
}

// readFile reads content of a file
func (l *Logs) readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// readUrl reads content of a page
func (l *Logs) readURL(url string) ([]byte, error) {
	u, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer u.Body.Close()
	b, err := ioutil.ReadAll(u.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
