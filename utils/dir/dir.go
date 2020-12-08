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
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
)

// CreateDataDir creates/updates directories with the given permissions in the persistent volume.
func CreateDataDir(path, username, groupname string, perm os.FileMode) error {
	// try to create data directory
	if err := os.MkdirAll(path, perm); err != nil {
		return errors.Wrap(err, "cannot create datadir")
	}

	var storedErr error // store the first encountered error
	// check and fix directory permissions
	dataDirStat, err := os.Stat(path)
	if err != nil {
		storedErr = errors.Wrapf(err, "cannot get stat of %q", path)
	}

	if err := os.Chmod(path, perm); err != nil {
		err = errors.Wrapf(err, "cannot chmod path %q", path)
		if storedErr == nil {
			storedErr = err
		}
	}

	dataDirSysStat := dataDirStat.Sys().(*syscall.Stat_t)
	aUID, aGID := int(dataDirSysStat.Uid), int(dataDirSysStat.Gid)

	dirUser, err := user.Lookup(username)
	if err != nil {
		return errors.Wrap(err, "cannot chown datadir")
	}
	bUID, err := strconv.Atoi(dirUser.Uid)
	if err != nil {
		if storedErr != nil {
			return storedErr
		}
		storedErr = errors.Wrap(err, "cannot chown datadir")
	}

	group, err := user.LookupGroup(groupname)
	if err != nil {
		return errors.Wrap(err, "cannot chown datadir")
	}
	bGID, err := strconv.Atoi(group.Gid)
	if err != nil {
		if storedErr != nil {
			return storedErr
		}
		storedErr = errors.Wrap(err, "cannot chown datadir")
	}

	if aUID != bUID || aGID != bGID {
		if err := os.Chown(path, bUID, bGID); err != nil {
			if storedErr != nil {
				return storedErr
			}
			storedErr = errors.Wrap(err, "cannot chown datadir")
		}
	}
	return storedErr
}
