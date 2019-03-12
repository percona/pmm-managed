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

package handlers

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toBadRequestWithFieldViolation(err error) error {
	st := status.New(codes.InvalidArgument, err.Error())

	errSlice := strings.SplitN(err.Error(), ":", 3)
	if len(errSlice) > 1 {
		st = status.New(codes.InvalidArgument, "validation error")
		fieldName := strcase.ToSnake(strings.TrimSpace(strings.Replace(errSlice[0], "invalid field", "", 1)))
		errorMsg := strings.TrimSpace(errSlice[1])

		v := &errdetails.BadRequest_FieldViolation{
			Field:       fieldName,
			Description: errorMsg,
		}
		br := &errdetails.BadRequest{}
		br.FieldViolations = append(br.FieldViolations, v)
		var err error
		st, err = st.WithDetails(br)
		if err != nil {
			panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
		}
	}

	return st.Err()
}
