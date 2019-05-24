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

package action

import (
	"context"
)

type service struct {
	r   Runner
	s   Storage
	rsv PMMAgentIDResolver
	dsn DSNResolver
}

// NewService creates default actions service implementation.
func NewService(r Runner, s Storage, rsv PMMAgentIDResolver, dsn DSNResolver) Service {
	return &service{
		r:   r,
		s:   s,
		rsv: rsv,
		dsn: dsn,
	}
}

// GetAction gets an action result.
//nolint:lll
func (s *service) GetActionResult(ctx context.Context, actionID string) (*Result, error) {
	res, err := s.s.Load(ctx, actionID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// StartPTSummaryAction starts pt-summary action.
//nolint:lll
func (s *service) StartPTSummaryAction(ctx context.Context, pmmAgentID, nodeID string) (*PtSummary, error) {
	a := NewPtSummary(pmmAgentID, nodeID)
	var err error
	a.PMMAgentID, err = s.rsv.ResolvePMMAgentIDByNodeID(a.NodeID, a.PMMAgentID)
	if err != nil {
		return a, err
	}

	err = s.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, err
	}

	err = s.r.StartPTSummaryAction(ctx, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// StartPTMySQLSummaryAction starts pt-mysql-summary action.
//nolint:lll
func (s *service) StartPTMySQLSummaryAction(ctx context.Context, pmmAgentID, serviceID string) (*PtMySQLSummary, error) {
	a := &PtMySQLSummary{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Args:       []string{},
	}

	var err error
	a.PMMAgentID, err = s.rsv.ResolvePMMAgentIDByServiceID(a.ServiceID, a.PMMAgentID)
	if err != nil {
		return a, err
	}

	err = s.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, err
	}

	err = s.r.StartPTMySQLSummaryAction(ctx, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// StartMySQLExplainAction starts mysql-explain action.
//nolint:lll
func (s *service) StartMySQLExplainAction(ctx context.Context, pmmAgentID, serviceID, query string) (*MySQLExplain, error) {
	a := &MySQLExplain{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Query:      query,
	}

	var err error
	a.PMMAgentID, err = s.rsv.ResolvePMMAgentIDByServiceID(a.ServiceID, a.PMMAgentID)
	if err != nil {
		return a, err
	}

	a.Dsn, err = s.dsn.ResolveDSNByServiceID(a.ServiceID)
	if err != nil {
		return a, err
	}

	err = s.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLExplainAction(ctx, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// StartMySQLExplainJSONAction starts mysql-explain json action.
//nolint:lll
func (s *service) StartMySQLExplainJSONAction(ctx context.Context, pmmAgentID, serviceID, query string) (*MySQLExplainJSON, error) {
	a := &MySQLExplainJSON{
		ID:         getUUID(),
		ServiceID:  serviceID,
		PMMAgentID: pmmAgentID,
		Query:      query,
	}

	var err error
	a.PMMAgentID, err = s.rsv.ResolvePMMAgentIDByServiceID(a.ServiceID, a.PMMAgentID)
	if err != nil {
		return a, err
	}

	a.Dsn, err = s.dsn.ResolveDSNByServiceID(a.ServiceID)
	if err != nil {
		return a, err
	}

	err = s.s.Store(ctx, &Result{ID: a.ID, PmmAgentID: a.PMMAgentID})
	if err != nil {
		return nil, err
	}

	err = s.r.StartMySQLExplainJSONAction(ctx, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// CancelAction stops an Action.
//nolint:lll
func (s *service) CancelAction(ctx context.Context, actionID string) error {
	ar, err := s.s.Load(ctx, actionID)
	if err != nil {
		return err
	}

	err = s.r.StopAction(ctx, ar.ID, ar.PmmAgentID)
	if err != nil {
		return err
	}
	return nil
}
