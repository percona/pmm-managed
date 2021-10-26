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
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

// GetOktaSSODetails return OktaSSODetails if there are any, error otherwise.
func GetOktaSSODetails(q *reform.Querier) (*OktaSSODetails, error) {
	ssoDetails, err := q.SelectOneFrom(OktaSSODetailsView, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Okta SSO Details")
	}
	return ssoDetails.(*OktaSSODetails), nil
}

// DeleteOktaSSODetails removes all stored DeleteOktaSSODetails.
func DeleteOktaSSODetails(q *reform.Querier) error {
	_, err := q.DeleteFrom(OktaSSODetailsView, "")
	if err != nil {
		return errors.Wrap(err, "failed to delete Okta SSO Details")
	}
	return nil
}

// InsertOktaSSODetails inserts a new Okta SSO details.
func InsertOktaSSODetails(q *reform.Querier, ssoDetails *OktaSSODetails) error {
	if err := q.Insert(ssoDetails); err != nil {
		return errors.Wrap(err, "failed to insert Okta SSO Details")
	}
	return nil
}
