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

package models

import (
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-sql-driver/mysql"
)

func MySQLDSN(service *Service, exporter *Agent) string {
	// TODO TLSConfig: "true", https://jira.percona.com/browse/PMM-1727
	// TODO Other parameters?

	cfg := mysql.NewConfig()
	cfg.User = pointer.GetString(exporter.Username)
	cfg.Passwd = pointer.GetString(exporter.Password)
	cfg.Net = "tcp"
	host := pointer.GetString(service.Address)
	port := pointer.GetUint16(service.Port)
	cfg.Addr = net.JoinHostPort(host, strconv.Itoa(int(port)))
	cfg.Timeout = 1 * time.Second

	// QAN code in pmm-agent uses reform which requires those fields
	cfg.ClientFoundRows = true
	cfg.ParseTime = true

	return cfg.FormatDSN()
}

func PostgreSQLDSN(service *Service, exporter *Agent) string {
	q := make(url.Values)
	q.Set("sslmode", "disable") // TODO: make it configurable
	q.Set("connect_timeout", "1")

	host := pointer.GetString(service.Address)
	port := pointer.GetUint16(service.Port)
	username := pointer.GetString(exporter.Username)
	password := pointer.GetString(exporter.Password)

	u := &url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(host, strconv.Itoa(int(port))),
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
	switch {
	case password != "":
		u.User = url.UserPassword(username, password)
	case username != "":
		u.User = url.User(username)
	}

	return u.String()
}

func MongoDBDSN(service *Service, exporter *Agent) string {
	host := pointer.GetString(service.Address)
	port := pointer.GetUint16(service.Port)
	username := pointer.GetString(exporter.Username)
	password := pointer.GetString(exporter.Password)

	u := &url.URL{
		Scheme: "mongodb",
		Host:   net.JoinHostPort(host, strconv.Itoa(int(port))),
	}
	switch {
	case password != "":
		u.User = url.UserPassword(username, password)
	case username != "":
		u.User = url.User(username)
	}

	return u.String()
}
