package models

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"
)

// FindCheckStateByName finds ChecksState by name.
func FindCheckStateByName(q *reform.Querier, name string) (*ChecksState, error) {
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty Check name.")
	}

	cs := &ChecksState{Name: name}
	switch err := q.Reload(cs); err {
	case nil:
		return cs, nil
	case reform.ErrNoRows:
		return nil, err
	default:
		return nil, errors.WithStack(err)
	}
}

// CreateCheckState persists ChecksState.
func CreateCheckState(q *reform.Querier, name string, interval Interval) (*ChecksState, error) {
	row := &ChecksState{
		Name:     name,
		Interval: interval,
	}

	if err := q.Insert(row); err != nil {
		return nil, errors.Wrap(err, "failed to create checks state")
	}

	return row, nil
}

// ChangeCheckState updates the interval of a check state if already present.
func ChangeCheckState(q *reform.Querier, name string, interval Interval) (*ChecksState, error) {
	row, err := FindCheckStateByName(q, name)
	if err != nil {
		return nil, err
	}

	row.Interval = interval

	if err := q.Update(row); err != nil {
		return nil, errors.Wrap(err, "failed to update checks state")
	}

	return row, nil
}
