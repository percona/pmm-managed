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

package dir

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// Params represent the input for CreateDataDir
type Params struct {
	Path  string
	Perm  os.FileMode
	User  string
	Group string
}

// CreateDataDir creates/updates directories with the given permissions in the persistent volume.
func CreateDataDir(params Params) error {
	// try to create data directory
	if err := os.MkdirAll(params.Path, params.Perm); err != nil {
		return fmt.Errorf("cannot create datadir %v", err)
	}

	// check and fix directory permissions
	dataDirStat, err := os.Stat(params.Path)
	if err != nil {
		return fmt.Errorf("cannot get stat of %q: %v", params.Path, err)
	}

	if dataDirStat.Mode()&os.ModePerm != params.Perm {
		if err := os.Chmod(params.Path, params.Perm); err != nil {
			return fmt.Errorf("cannot chmod datadir %v", err)
		}
	}

	dataDirSysStat := dataDirStat.Sys().(*syscall.Stat_t)
	aUID, aGID := int(dataDirSysStat.Uid), int(dataDirSysStat.Gid)

	dirUser, err := user.Lookup(params.User)
	if err != nil {
		return fmt.Errorf("cannot chown datadir %v", err)
	}
	bUID, err := strconv.Atoi(dirUser.Uid)
	if err != nil {
		return fmt.Errorf("cannot chown datadir %v", err)
	}

	group, err := user.LookupGroup(params.Group)
	if err != nil {
		return fmt.Errorf("cannot chown datadir %v", err)
	}
	bGID, err := strconv.Atoi(group.Gid)
	if err != nil {
		return fmt.Errorf("cannot chown datadir %v", err)
	}

	if aUID != bUID || aGID != bGID {
		if err := os.Chown(params.Path, bUID, bGID); err != nil {
			return fmt.Errorf("cannot chown datadir %v", err)
		}
	}
	return nil
}
