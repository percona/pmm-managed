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
	"text/template"
)

const (
	configPathPrefix = "/etc/supervisord.d/"
)

var tmpl = template.Must(template.New("").Option("missingkey=error").Parse(`
{{define "prometheus"}}
[program:prometheus]
priority = 7
command =
	/usr/sbin/prometheus
		--config.file=/etc/prometheus.yml
		--storage.tsdb.path=/srv/prometheus/data
		--storage.tsdb.retention.time={{ .DataRetentionDays }}d
		--web.listen-address=:9090
		--web.console.libraries=/usr/share/prometheus/console_libraries
		--web.console.templates=/usr/share/prometheus/consoles
		--web.external-url=http://localhost:9090/prometheus/
		--web.enable-admin-api
		--web.enable-lifecycle
user = pmm
autorestart = true
autostart = true
startretries = 3
startsecs = 1
stopsignal = TERM
stopwaitsecs = 300
stdout_logfile = /srv/logs/prometheus.log
stdout_logfile_maxbytes = 10MB
stdout_logfile_backups = 3
redirect_stderr = true
{{end}}

{{define "qan-api2"}}
[program:qan-api2]
priority = 13
command =
	/usr/sbin/percona-qan-api2
		--data-retention={{ .DataRetentionDays }}
user = pmm
autorestart = true
autostart = true
startretries = 1000
startsecs = 1
stopsignal = TERM
stopwaitsecs = 10
stdout_logfile = /srv/logs/qan-api2.log
stdout_logfile_maxbytes = 10MB
stdout_logfile_backups = 3
redirect_stderr = true
{{end}}
`))
