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

package okta

// SSOService includes logic for handling access_token for Okta SSO.
type SSOService struct{}

// NewSSOService creates a new Okta single sign-on service.
func NewSSOService() *SSOService {
	return &SSOService{}
}

// RenewAccessToken creates a request for the new access token. After it's successful the new token replaces the
// old one and is returned. As a side effect, when we don't get error, we know Okta secret is still
// valid and PMM is still connected to the Portal.
func (o *SSOService) RenewAccessToken() (string, error) {
	return "", nil
}

// AccessToken returns valid access token. It either returns one that's still valid or requests a new one.
func (o *SSOService) AccessToken() (string, error) {
	return "", nil
}
