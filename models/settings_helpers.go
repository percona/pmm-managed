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
	"encoding/json"
	"sort"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// GetSettings returns current PMM Server settings.
func GetSettings(q reform.DBTX) (*Settings, error) {
	var b []byte
	if err := q.QueryRow("SELECT settings FROM settings").Scan(&b); err != nil {
		return nil, errors.Wrap(err, "failed to select settings")
	}

	var s Settings
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal settings")
	}

	s.fillDefaults()
	return &s, nil
}

// SaveSettings saves PMM Server settings.
// It may modify passed settings to fill defaults.
func SaveSettings(q reform.DBTX, s *Settings) error {
	s.fillDefaults()

	var err error
	if s.AWSPartitions, err = validateAWSPartitions(s.AWSPartitions); err != nil {
		return err
	}

	b, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "failed to marshal settings")
	}

	_, err = q.Exec("UPDATE settings SET settings = $1", b)
	if err != nil {
		return errors.Wrap(err, "failed to update settings")
	}

	return nil
}

// validateAWSPartitions deduplicates and validates AWS partitions list.
func validateAWSPartitions(partitions []string) ([]string, error) {
	if len(partitions) > len(endpoints.DefaultPartitions()) {
		return nil, status.Errorf(codes.InvalidArgument, "aws_partitions: list is too long")
	}

	set := make(map[string]struct{})
	for _, p := range partitions {
		var valid bool
		for _, vp := range endpoints.DefaultPartitions() {
			if p == vp.ID() {
				valid = true
				break
			}
		}
		if !valid {
			return nil, status.Errorf(codes.InvalidArgument, "aws_partitions: partition %q is invalid", p)
		}
		set[p] = struct{}{}
	}

	slice := make([]string, 0, len(set))
	for partition := range set {
		slice = append(slice, partition)
	}
	sort.Strings(slice)

	return slice, nil
}
