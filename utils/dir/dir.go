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

// Package dir contains utilities for creating directories.
package dir

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// CreateDataDir creates/updates directories with the given permissions in the persistent volume.
func CreateDataDir(path, username, groupname string, perm os.FileMode) error {
	// try to create data directory
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("cannot create datadir %v", err)
	}

	var storedErr error // store the first encountered error
	// check and fix directory permissions
	dataDirStat, err := os.Stat(path)
	if err != nil {
		storedErr = fmt.Errorf("cannot get stat of %q: %v", path, err)
	}

	if dataDirStat.Mode()&os.ModePerm != perm {
		if err := os.Chmod(path, perm); err != nil {
			if storedErr != nil {
				return storedErr
			}
			storedErr = fmt.Errorf("cannot chmod datadir %v", err)
		}
	}

	dataDirSysStat := dataDirStat.Sys().(*syscall.Stat_t)
	aUID, aGID := int(dataDirSysStat.Uid), int(dataDirSysStat.Gid)

	dirUser, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("cannot chown datadir %v", err)
	}
	bUID, err := strconv.Atoi(dirUser.Uid)
	if err != nil {
		if storedErr != nil {
			return storedErr
		}
		storedErr = fmt.Errorf("cannot chown datadir %v", err)
	}

	group, err := user.LookupGroup(groupname)
	if err != nil {
		return fmt.Errorf("cannot chown datadir %v", err)
	}
	bGID, err := strconv.Atoi(group.Gid)
	if err != nil {
		if storedErr != nil {
			return storedErr
		}
		storedErr = fmt.Errorf("cannot chown datadir %v", err)
	}

	if aUID != bUID || aGID != bGID {
		if err := os.Chown(path, bUID, bGID); err != nil {
			if storedErr != nil {
				return storedErr
			}
			storedErr = fmt.Errorf("cannot chown datadir %v", err)
		}
	}
	return storedErr
}
