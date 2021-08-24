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

package backup

import "github.com/hashicorp/go-version"

type compatibility struct {
	dbMinVersion         *version.Version
	dbMaxVersion         *version.Version
	backupToolMinVersion *version.Version
	backupToolMaxVersion *version.Version
}

var mysqlAndXtrabackupCompatibleVersions []compatibility

func init() {
	versionStrings := []struct {
		mysqlMinVersion      string // inclusively
		mysqlMaxVersion      string // exclusively
		xtrabackupMinVersion string // inclusively
		xtrabackupMaxVersion string // exclusively
	}{
		// It can back up data from InnoDB, XtraDB, and MyISAM tables on MySQL 5.5, 5.6 and 5.7 servers,
		// as well as Percona Server for MySQL with XtraDB.
		// https://www.percona.com/doc/percona-xtrabackup/2.4/index.html
		{
			mysqlMinVersion:      "5.5",
			mysqlMaxVersion:      "5.8",
			xtrabackupMinVersion: "2.4",
			xtrabackupMaxVersion: "2.5",
		},
		// Percona XtraBackup 8.0.4 is based on MySQL 8.0.13 and fully supports Percona Server for MySQL 8.0 series and MySQL 8.0 series.
		// https://www.percona.com/doc/percona-xtrabackup/8.0/release-notes/8.0/8.0.4.html
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.14",
			xtrabackupMinVersion: "8.0.4",
			xtrabackupMaxVersion: "8.0.5",
		},
		// Percona XtraBackup 8.0.5 introduces the support of undo tablespaces created using
		// the new syntax (CREATE UNDO TABLESPACE) available since MySQL 8.0.14.
		// Percona XtraBackup also supports the binary log encryption introduced in MySQL 8.0.14.
		// https://www.percona.com/doc/percona-xtrabackup/8.0/release-notes/8.0/8.0.5.html
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.15",
			xtrabackupMinVersion: "8.0.5",
			xtrabackupMaxVersion: "8.0.6",
		},
		// In version 8.0.6, Percona XtraBackup introduces the support of the MyRocks storage engine
		// with Percona Server for MySQL version 8.0.15-6 or higher.
		// https://www.percona.com/doc/percona-xtrabackup/8.0/release-notes/8.0/8.0.6.html
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.20",
			xtrabackupMinVersion: "8.0.6",
			xtrabackupMaxVersion: "8.0.12",
		},
		// Percona XtraBackup 8.0.12 now supports backup and restore processing for all versions of MySQL;
		// previous versions of Percona XtraBackup will not work with MySQL 8.0.20 and higher.
		// https://www.percona.com/doc/percona-xtrabackup/8.0/release-notes/8.0/8.0.12.html
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.21",
			xtrabackupMinVersion: "8.0.12",
			xtrabackupMaxVersion: "8.0.13",
		},
		// Percona XtraBackup 8.0.13 supports backup and restore processing for all versions of MySQL
		// and has been tested with the latest MySQL 8.0.20.
		// https://www.percona.com/doc/percona-xtrabackup/8.0/release-notes/8.0/8.0.13.html
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.21",
			xtrabackupMinVersion: "8.0.13",
			xtrabackupMaxVersion: "8.0.14",
		},
		// Percona XtraBackup 8.0.14 supports backup and restore processing for all versions of MySQL
		// and has been tested with the latest MySQL 8.0.21.
		// https://www.percona.com/doc/percona-xtrabackup/8.0/release-notes/8.0/8.0.14.html
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.22",
			xtrabackupMinVersion: "8.0.14",
			xtrabackupMaxVersion: "8.0.15",
		},
		// If the Percona XtraBackup version is lower than the database version,
		// processing will be stopped and Percona XtraBackup will not be allowed to continue.
		// https://www.percona.com/blog/2020/08/18/aligning-percona-xtrabackup-versions-with-percona-server-for-mysql/
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.23",
			xtrabackupMinVersion: "8.0.22",
			xtrabackupMaxVersion: "8.0.23",
		},
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.24",
			xtrabackupMinVersion: "8.0.23",
			xtrabackupMaxVersion: "8.0.24",
		},
		{
			mysqlMinVersion:      "8.0",
			mysqlMaxVersion:      "8.0.26",
			xtrabackupMinVersion: "8.0.25",
			xtrabackupMaxVersion: "8.0.26",
		},
	}

	mysqlAndXtrabackupCompatibleVersions = make([]compatibility, 0, len(versionStrings))
	for _, s := range versionStrings {
		mysqlMinVersion, err := version.NewVersion(s.mysqlMinVersion)
		if err != nil {
			panic(err)
		}
		mysqlMaxVersion, err := version.NewVersion(s.mysqlMaxVersion)
		if err != nil {
			panic(err)
		}
		xtrabackupMinVersion, err := version.NewVersion(s.xtrabackupMinVersion)
		if err != nil {
			panic(err)
		}
		xtrabackupMaxVersion, err := version.NewVersion(s.xtrabackupMaxVersion)
		if err != nil {
			panic(err)
		}

		mysqlAndXtrabackupCompatibleVersions = append(mysqlAndXtrabackupCompatibleVersions, compatibility{
			dbMinVersion:         mysqlMinVersion,
			dbMaxVersion:         mysqlMaxVersion,
			backupToolMinVersion: xtrabackupMinVersion,
			backupToolMaxVersion: xtrabackupMaxVersion,
		})
	}
}

func mysqlAndXtrabackupCompatible(mysqlVersionString, xtrabackupVersionString string) (bool, error) {
	mysqlVersion, err := version.NewVersion(mysqlVersionString)
	if err != nil {
		return false, err
	}
	mysqlVersion = mysqlVersion.Core()

	xtrabackupVersion, err := version.NewVersion(xtrabackupVersionString)
	if err != nil {
		return false, err
	}
	xtrabackupVersion = xtrabackupVersion.Core()

	for _, cv := range mysqlAndXtrabackupCompatibleVersions {
		if !(mysqlVersion.GreaterThanOrEqual(cv.dbMinVersion) && mysqlVersion.LessThan(cv.dbMaxVersion)) {
			continue
		}

		if xtrabackupVersion.GreaterThanOrEqual(cv.backupToolMinVersion) &&
			xtrabackupVersion.LessThan(cv.backupToolMaxVersion) {
			return true, nil
		}
	}

	return false, nil
}
