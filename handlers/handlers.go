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
