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
	"syscall"
)

// Params represent the input for CreateDataDir
type Params struct {
	Path string
	Perm os.FileMode
	// The path of the directory whose uID and gID will be used
	// to chown the created dir.
	ChownPath string
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

	if params.ChownPath != "" {
		dataDirSysStat := dataDirStat.Sys().(*syscall.Stat_t)
		aUID, aGID := int(dataDirSysStat.Uid), int(dataDirSysStat.Gid)

		chownDirStat, err := os.Stat(params.ChownPath)
		if err != nil {
			return fmt.Errorf("cannot get stat of %q: %v", params.ChownPath, err)
		}

		chownDirSysStat := chownDirStat.Sys().(*syscall.Stat_t)
		bUID, bGID := int(chownDirSysStat.Uid), int(chownDirSysStat.Gid)
		if aUID != bUID || aGID != bGID {
			if err := os.Chown(params.Path, bUID, bGID); err != nil {
				return fmt.Errorf("cannot chown datadir %v", err)
			}
		}
	}
	return nil
}
