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

package supervisord

import (
	"archive/zip"
	"bufio"
	"container/ring"
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/percona/pmm/utils/pdeathsig"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	"github.com/percona/pmm-managed/utils/logger"
)

// logInfo represents log file information, or the way to read log.
type logInfo struct {
	FilePath string
}

// fileContent represents log or configuration file content.
type fileContent struct {
	Name string
	Data []byte
	Err  error
}

var defaultLogs = map[string]logInfo{
	// system
	"cron.log":        {"/srv/logs/cron.log"},
	"supervisord.log": {"/srv/logs/supervisord.log"},

	// storages
	"clickhouse-server.log":     {"/srv/logs/clickhouse-server.log"},
	"clickhouse-server.err.log": {"/srv/logs/clickhouse-server.err.log"},
	"postgresql.log":            {"/srv/logs/postgresql.log"},

	// nginx
	"nginx.log":        {"/srv/logs/nginx.startup.log"},
	"nginx_access.log": {"/srv/logs/nginx.access.log"},
	"nginx_error.log":  {"/srv/logs/nginx.error.log"},

	// metrics
	"prometheus.log": {"/srv/logs/prometheus.log"},
	"grafana.log":    {"/srv/logs/grafana.log"},

	// core PMM components
	"pmm-managed.log": {"/srv/logs/pmm-managed.log"},
	"qan-api2.log":    {"/srv/logs/qan-api2.log"},

	// upgrades
	"dashboard-upgrade.log": {"/srv/logs/dashboard-upgrade.log"},
}

// Logs is responsible for interactions with logs.
type Logs struct {
	pmmVersion string
	logs       map[string]logInfo // for testing
}

// NewLogs creates a new Logs service.
// n is a number of last lines of log to read.
func NewLogs(pmmVersion string) *Logs {
	return &Logs{
		pmmVersion: pmmVersion,
		logs:       defaultLogs,
	}
}

// Zip creates .zip archive with all logs.
func (l *Logs) Zip(ctx context.Context, w io.Writer) error {
	zw := zip.NewWriter(w)
	now := time.Now().UTC()
	for _, file := range l.files(ctx) {
		if file.Err != nil {
			logger.Get(ctx).WithField("component", "logs").Errorf("%s: %s", file.Name, file.Err)

			// do not let a single error break the whole archive
			if len(file.Data) > 0 {
				file.Data = append(file.Data, "\n\n"...)
			}
			file.Data = append(file.Data, file.Err.Error()...)
		}

		f, err := zw.CreateHeader(&zip.FileHeader{
			Name:     file.Name,
			Method:   zip.Deflate,
			Modified: now,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create zip file header")
		}
		if _, err = f.Write(file.Data); err != nil {
			return errors.Wrap(err, "failed to write zip file data")
		}
	}
	return errors.Wrap(zw.Close(), "failed to close zip file")
}

// files reads log/config files and returns content.
func (l *Logs) files(ctx context.Context) []fileContent {
	files := make([]fileContent, 0, len(l.logs))

	for name, log := range l.logs {
		f := fileContent{
			Name: name,
		}
		f.Data, f.Err = readLog(log.FilePath, 1000, 1024*1024)
		files = append(files, f)
	}

	// add PMM version
	files = append(files, fileContent{
		Name: "pmm-version.txt",
		Data: []byte(l.pmmVersion + "\n"),
	})

	// add configs
	for _, f := range []string{
		"/etc/prometheus.yml",
		"/etc/supervisord.d/pmm.ini",
		"/etc/nginx/conf.d/pmm.conf",
	} {
		b, err := ioutil.ReadFile(f) //nolint:gosec
		files = append(files, fileContent{
			Name: filepath.Base(f),
			Data: b,
			Err:  err,
		})
	}

	// add supervisord status
	cmd := exec.CommandContext(ctx, "supervisorctl", "status") //nolint:gosec
	pdeathsig.Set(cmd, unix.SIGKILL)
	b, err := cmd.CombinedOutput() //nolint:gosec
	files = append(files, fileContent{
		Name: "supervisorctl_status.log",
		Data: b,
		Err:  err,
	})

	// add systemd status for OVF/AMI
	cmd = exec.CommandContext(ctx, "systemctl", "-l", "status") //nolint:gosec
	pdeathsig.Set(cmd, unix.SIGKILL)
	b, err = cmd.CombinedOutput() //nolint:gosec
	files = append(files, fileContent{
		Name: "systemctl_status.log",
		Data: b,
		Err:  err,
	})

	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	return files
}

// readLog reads last lines (up to given number of lines and bytes) from given file.
func readLog(name string, maxLines int, maxBytes int64) ([]byte, error) {
	f, err := os.Open(name) //nolint:gosec
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close() //nolint:errcheck

	fi, err := f.Stat()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if fi.Size() > maxBytes {
		if _, err = f.Seek(-maxBytes, io.SeekEnd); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	r := ring.New(maxLines)
	s := bufio.NewScanner(f)
	for s.Scan() {
		r.Value = []byte(s.Text() + "\n")
		r = r.Next()
	}
	if err = s.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]byte, 0, maxBytes)
	r.Do(func(v interface{}) {
		if v != nil {
			res = append(res, v.([]byte)...)
		}
	})
	return res, nil
}
