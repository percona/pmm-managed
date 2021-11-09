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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

const issuerSubdirectoryAndQuery = "/oauth2/aus15pi5rjdtfrcH51d7/v1/token?grant_type=client_credentials&scope="

var perconaSSOMtx sync.Mutex

// GetPerconaSSODetails returns PerconaSSODetails if there are any, error otherwise.
// Access token is automatically refreshed if it is expired.
// Get, check eventually refresh did in one tx.
func GetPerconaSSODetails(ctx context.Context, q *reform.Querier) (*PerconaSSODetails, error) {
	perconaSSOMtx.Lock()
	defer perconaSSOMtx.Unlock()

	ssoDetails, err := q.SelectOneFrom(PerconaSSODetailsView, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Percona SSO Details")
	}

	details := ssoDetails.(*PerconaSSODetails)
	if details.isAccessTokenExpired() {
		refreshedToken, err := details.refreshAndGetAccessToken(ctx, q)
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert Percona SSO Details")
		}
		details.AccessToken = refreshedToken
	}

	return details, nil
}

// GetPerconaSSOAccessToken returns PerconaSSOAccessToken if there is GetPerconaSSODetails, error otherwise.
func GetPerconaSSOAccessToken(ctx context.Context, q *reform.Querier) (*PerconaSSOAccessToken, error) {
	ssoDetails, err := GetPerconaSSODetails(ctx, q)
	if err != nil {
		return nil, err
	}

	return ssoDetails.AccessToken, nil
}

func (sso *PerconaSSODetails) refreshAndGetAccessToken(ctx context.Context, q *reform.Querier) (*PerconaSSOAccessToken, error) {
	url := sso.IssuerURL + fmt.Sprintf("%s%s", issuerSubdirectoryAndQuery, sso.Scope)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(sso.ClientID + ":" + sso.ClientSecret))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")

	timeBeforeRequest := time.Now()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		var accessToken *PerconaSSOAccessToken
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bodyBytes, &accessToken); err != nil {
			return nil, err
		}

		accessToken.ExpiresAt = timeBeforeRequest.Add(time.Duration(accessToken.ExpiresIn) * time.Second)
		sso.AccessToken = accessToken

		if err := q.Insert(sso); err != nil {
			return nil, err
		}

		return accessToken, nil
	}

	return nil, fmt.Errorf("get access token of Percona SSO Details failed, status code: %d", res.StatusCode)
}

func (sso *PerconaSSODetails) isAccessTokenExpired() bool {
	if sso == nil || sso.AccessToken == nil {
		return true
	}

	return sso.AccessToken.ExpiresAt.After(time.Now())
}

// DeletePerconaSSODetails removes all stored DeletePerconaSSODetails.
func DeletePerconaSSODetails(q *reform.Querier) error {
	_, err := q.DeleteFrom(PerconaSSODetailsView, "")
	if err != nil {
		return errors.Wrap(err, "failed to delete Percona SSO Details")
	}
	return nil
}

// InsertPerconaSSODetails inserts a new Percona SSO details.
func InsertPerconaSSODetails(ctx context.Context, q *reform.Querier, ssoDetails *PerconaSSODetailsInsert) error {
	details := &PerconaSSODetails{
		IssuerURL:    ssoDetails.IssuerURL,
		ClientID:     ssoDetails.ClientID,
		ClientSecret: ssoDetails.ClientSecret,
		Scope:        ssoDetails.Scope,
	}

	_, err := details.refreshAndGetAccessToken(ctx, q)
	if err != nil {
		return errors.Wrap(err, "failed to insert Percona SSO Details")
	}
	return nil
}
