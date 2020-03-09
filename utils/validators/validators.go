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

// Package validators contains settings validators.
package validators

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// MetricsResolutionMin is the smallest value metric resolution can accept.
	MetricsResolutionMin = time.Second
	// MetricsResolutionMultipleOf is value metrics resolution should be multiple of.
	MetricsResolutionMultipleOf = time.Second
	// DataRetentionMin is the smallest value data retention can accept.
	DataRetentionMin = 24 * time.Hour
	// DataRetentionMultipleOf is a value of data retention should be multiple of.
	DataRetentionMultipleOf = 24 * time.Hour
)

// MinDurationError minimum allowed duration error.
type MinDurationError struct {
	Msg string
	Min time.Duration
}

func (e MinDurationError) Error() string { return e.Msg }

// AliquotDurationError multiple of duration allowed error.
type AliquotDurationError struct {
	Msg        string
	MultipleOf time.Duration
}

func (e AliquotDurationError) Error() string { return e.Msg }

// ValidateDuration validate duration.
func validateDuration(d, min, multipleOf time.Duration) (time.Duration, error) {
	if d < min {
		return d, MinDurationError{"min duration error", min}
	}

	if d.Truncate(multipleOf) != d {
		return d, AliquotDurationError{"aliquot	duration error", multipleOf}
	}
	return d, nil
}

// ValidateMetricResolution validate metric resolution.
func ValidateMetricResolution(value time.Duration) (time.Duration, error) {
	return validateDuration(value, MetricsResolutionMin, MetricsResolutionMultipleOf)
}

// ValidateDataRetention validate metric resolution.
func ValidateDataRetention(value time.Duration) (time.Duration, error) {
	return validateDuration(value, DataRetentionMin, DataRetentionMultipleOf)
}

// ValidateAWSPartitions validates AWS partitions list.
func ValidateAWSPartitions(partitions []string) error {
	if len(partitions) > len(endpoints.DefaultPartitions()) {
		return status.Errorf(codes.InvalidArgument, "aws_partitions: list is too long")
	}

	for _, p := range partitions {
		var valid bool
		for _, vp := range endpoints.DefaultPartitions() {
			if p == vp.ID() {
				valid = true
				break
			}
		}
		if !valid {
			return status.Errorf(codes.InvalidArgument, "aws_partitions: partition %q is invalid", p)
		}
	}

	return nil
}
