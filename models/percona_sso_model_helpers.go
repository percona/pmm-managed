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
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

const issuerSubdirectoryAndQuery = "/oauth2/aus15pi5rjdtfrcH51d7/v1/token?grant_type=client_credentials&scope="

// GetPerconaSSODetails returns PerconaSSODetails if there are any, error otherwise.
// Access token is automatically refresh if it is expired.
func GetPerconaSSODetails(q *reform.Querier) (*PerconaSSODetails, error) {
	ssoDetails, err := q.SelectOneFrom(PerconaSSODetailsView, "")
	if err != nil {
		return nil, err
	}

	details := ssoDetails.(*PerconaSSODetails)

	if details.isAccessTokenExpired() {
		refreshedToken, err := details.refreshAndGetAccessToken(q)
		if err != nil {
			return nil, err
		}
		details.AccessToken = refreshedToken
	}

	return details, nil
}

// GetPerconaSSOAccessToken returns PerconaSSOAccessToken if there are any, error otherwise.
func GetPerconaSSOAccessToken(q *reform.Querier) (*PerconaSSOAccessToken, error) {
	ssoDetails, err := GetPerconaSSODetails(q)
	if err != nil {
		return nil, err
	}

	if ssoDetails.isAccessTokenExpired() {
		refreshedToken, err := ssoDetails.refreshAndGetAccessToken(q)
		if err != nil {
			return nil, err
		}
		return refreshedToken, nil
	}

	return ssoDetails.AccessToken, nil
}

func (sso *PerconaSSODetails) refreshAndGetAccessToken(q *reform.Querier) (*PerconaSSOAccessToken, error) {
	url := sso.IssuerURL + fmt.Sprintf("%s%s", issuerSubdirectoryAndQuery, sso.Scope)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(sso.ClientID + ":" + sso.ClientSecret))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")

	timeBeforeRequest := time.Now()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to refresh access")
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
func InsertPerconaSSODetails(q *reform.Querier, ssoDetails *PerconaSSODetailsInsert) error {
	details := &PerconaSSODetails{
		IssuerURL:    ssoDetails.IssuerURL,
		ClientID:     ssoDetails.ClientID,
		ClientSecret: ssoDetails.ClientSecret,
		Scope:        ssoDetails.Scope,
	}

	accessToken, err := GetPerconaSSOAccessToken(q)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	details.AccessToken = accessToken
	if details.isAccessTokenExpired() {
		refreshedToken, err := details.refreshAndGetAccessToken(q)
		if err != nil {
			return err
		}
		details.AccessToken = refreshedToken
	}

	if err := q.Insert(details); err != nil {
		return err
	}
	return nil
}
