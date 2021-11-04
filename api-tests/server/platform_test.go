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

package server

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

// Tests in this file cover Percona Platform authentication.

func TestPlatform(t *testing.T) {
	client := serverClient.Default.Server

	t.Run("signUp", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			email, _, firstName, lastName := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:     email,
					FirstName: firstName,
					LastName:  lastName,
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)
		})

		t.Run("invalid email", func(t *testing.T) {
			_, _, firstName, lastName := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:     "not-email",
					FirstName: firstName,
					LastName:  lastName,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Error Creating Your Account.")
		})

		t.Run("empty email", func(t *testing.T) {
			_, _, firstName, lastName := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:     "",
					FirstName: firstName,
					LastName:  lastName,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Email: value '' must not be an empty string")
		})

		t.Run("empty first name", func(t *testing.T) {
			email, _, _, lastName := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:     email,
					FirstName: "",
					LastName:  lastName,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Error Creating Your Account.")
		})

		t.Run("empty last name", func(t *testing.T) {
			email, _, firstName, _ := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:     email,
					FirstName: firstName,
					LastName:  "",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Error Creating Your Account.")
		})
	})

	t.Run("connect", func(t *testing.T) {
		// TODO Change this test once real API for Portal is ready.
		if os.Getenv("PERCONA_SSO_CLIENT_ID") == "" || os.Getenv("PERCONA_SSO_CLIENT_SECRET") == "" || os.Getenv("PERCONA_SSO_ISSUER_URL") == "" || os.Getenv("PERCONA_SSO_SCOPE") == "" {
			t.Skip("Skipping - some of the required Percona SSO environment variables are not set.")
		}

		const serverName string = "my PMM"

		// TODO make sure we are using existing credentials when using real Platform API for connecting.
		email, password, _, _ := genCredentials(t)

		t.Run("wrong email", func(t *testing.T) {
			t.Skip("Skip until we use the real credentials for dev")
			_, err := client.PlatformConnect(&server.PlatformConnectParams{
				Body: server.PlatformConnectBody{
					Email:      "wrong@example.com",
					Password:   password,
					ServerName: serverName,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Incorrect username or password.")
		})

		t.Run("wrong password", func(t *testing.T) {
			t.Skip("Skip until we use the real credentials for dev")
			_, err := client.PlatformConnect(&server.PlatformConnectParams{
				Body: server.PlatformConnectBody{
					Email:      email,
					Password:   "WrongPassword12345",
					ServerName: serverName,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Incorrect username or password.")
		})

		t.Run("empty email", func(t *testing.T) {
			_, err := client.PlatformConnect(&server.PlatformConnectParams{
				Body: server.PlatformConnectBody{
					Email:      "",
					Password:   password,
					ServerName: serverName,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Email: value '' must not be an empty string")
		})

		t.Run("empty password", func(t *testing.T) {
			_, err := client.PlatformConnect(&server.PlatformConnectParams{
				Body: server.PlatformConnectBody{
					Email:      email,
					ServerName: serverName,
					Password:   "",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Password: value '' must not be an empty string")
		})

		t.Run("empty server name", func(t *testing.T) {
			_, err := client.PlatformConnect(&server.PlatformConnectParams{
				Body: server.PlatformConnectBody{
					Email:      email,
					Password:   password,
					ServerName: "",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field ServerName: value '' must not be an empty string")
		})

		t.Run("normal", func(t *testing.T) {
			_, err := client.PlatformConnect(&server.PlatformConnectParams{
				Body: server.PlatformConnectBody{
					ServerName: "my PMM server",
					Email:      email,
					Password:   password,
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)

			// Check SSO is setup in the grafana.ini supervisord config.
			grafanaConfig, err := ioutil.ReadFile("/etc/supervisord.d/grafana.ini")
			require.NoError(t, err)
			assert.True(t, strings.Contains(string(grafanaConfig), "cfg:default.auth.generic_oauth.enabled=true"), "generic_oauth should have been enabled")

			// Confirm we are connected to Portal.
			settings, err := client.GetSettings(nil)
			require.NoError(t, err)
			require.NotNil(t, settings)
			assert.True(t, settings.GetPayload().Settings.ConnectedToPortal)
		})

	})
}

// genCredentials creates test user email, password, firstName and lastName.
func genCredentials(t *testing.T) (string, string, string, string) {
	hostname, err := os.Hostname()
	require.NoError(t, err)

	u, err := user.Current()
	require.NoError(t, err)

	email := strings.Join([]string{u.Username, hostname, gofakeit.Email(), "test"}, ".")
	password := gofakeit.Password(true, true, true, false, false, 14)
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	return email, password, firstName, lastName
}
