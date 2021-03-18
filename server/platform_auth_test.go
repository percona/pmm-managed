package server

import (
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/api/serverpb/json/client/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

// Tests in this file cover Percona Platform authentication.

func TestPlatform(t *testing.T) {
	client := serverClient.Default.Server

	t.Run("signUp", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			email, password := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:    email,
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)
		})

		t.Run("invalid email", func(t *testing.T) {
			_, password := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:    "not-email",
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Error Creating Your Account.")
		})

		t.Run("invalid password", func(t *testing.T) {
			email, _ := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:    email,
					Password: "weak-pass",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Error Creating Your Account.")
		})

		t.Run("empty email", func(t *testing.T) {
			_, password := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:    "",
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Email: value '' must not be an empty string")
		})

		t.Run("empty password", func(t *testing.T) {
			email, _ := genCredentials(t)
			_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
				Body: server.PlatformSignUpBody{
					Email:    email,
					Password: "",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Password: value '' must not be an empty string")
		})
	})

	t.Run("signIn", func(t *testing.T) {
		email, password := genCredentials(t)

		_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
			Body: server.PlatformSignUpBody{
				Email:    email,
				Password: password,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		t.Run("normal", func(t *testing.T) {
			_, err = client.PlatformSignIn(&server.PlatformSignInParams{
				Body: server.PlatformSignInBody{
					Email:    email,
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)
		})

		t.Run("wrong email", func(t *testing.T) {
			_, err = client.PlatformSignIn(&server.PlatformSignInParams{
				Body: server.PlatformSignInBody{
					Email:    "wrong@example.com",
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Incorrect username or password.")
		})

		t.Run("wrong password", func(t *testing.T) {
			_, err = client.PlatformSignIn(&server.PlatformSignInParams{
				Body: server.PlatformSignInBody{
					Email:    email,
					Password: "WrongPassword12345",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Incorrect username or password.")
		})

		t.Run("empty email", func(t *testing.T) {
			_, err = client.PlatformSignIn(&server.PlatformSignInParams{
				Body: server.PlatformSignInBody{
					Email:    "",
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Email: value '' must not be an empty string")
		})

		t.Run("empty password", func(t *testing.T) {
			_, err = client.PlatformSignIn(&server.PlatformSignInParams{
				Body: server.PlatformSignInBody{
					Email:    email,
					Password: "",
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Password: value '' must not be an empty string")
		})
	})

	t.Run("signOut", func(t *testing.T) {
		email, password := genCredentials(t)

		_, err := client.PlatformSignUp(&server.PlatformSignUpParams{
			Body: server.PlatformSignUpBody{
				Email:    email,
				Password: password,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		t.Run("normal", func(t *testing.T) {
			_, err = client.PlatformSignIn(&server.PlatformSignInParams{
				Body: server.PlatformSignInBody{
					Email:    email,
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)

			_, err = client.PlatformSignOut(&server.PlatformSignOutParams{
				Body: server.PlatformSignInBody{
					Email:    email,
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			require.NoError(t, err)
		})

		t.Run("no active session", func(t *testing.T) {
			_, err = client.PlatformSignOut(&server.PlatformSignOutParams{
				Body: server.PlatformSignInBody{
					Email:    email,
					Password: password,
				},
				Context: pmmapitests.Context,
			})
			pmmapitests.AssertAPIErrorf(t, err, 400, codes.FailedPrecondition, "No active sessions.")
		})
	})
}

// genCredentials creates test user email and password.
func genCredentials(t *testing.T) (string, string) {
	hostname, err := os.Hostname()
	require.NoError(t, err)

	u, err := user.Current()
	require.NoError(t, err)

	email := strings.Join([]string{u.Username, hostname, gofakeit.Email(), "test"}, ".")
	password := gofakeit.Password(true, true, true, false, false, 14)
	return email, password
}
