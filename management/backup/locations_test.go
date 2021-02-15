package backup

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	backupClient "github.com/percona/pmm/api/managementpb/backup/json/client"
	"github.com/percona/pmm/api/managementpb/backup/json/client/locations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestAddLocation(t *testing.T) {
	t.Parallel()
	client := backupClient.Default.Locations

	t.Run("normal pmm client config", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
				PMMClientConfig: &locations.AddLocationParamsBodyPMMClientConfig{
					Path: "/tmp",
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteLocation(t, client, resp.Payload.LocationID)

		assert.NotEmpty(t, resp.Payload.LocationID)
	})

	t.Run("normal pmm server config", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
				PMMServerConfig: &locations.AddLocationParamsBodyPMMServerConfig{
					Path: "/tmp",
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteLocation(t, client, resp.Payload.LocationID)

		assert.NotEmpty(t, resp.Payload.LocationID)
	})

	t.Run("normal s3 config", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
				S3Config: &locations.AddLocationParamsBodyS3Config{
					Endpoint:  "http://example.com",
					AccessKey: "access_key",
					SecretKey: "secret_key",
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteLocation(t, client, resp.Payload.LocationID)

		assert.NotEmpty(t, resp.Payload.LocationID)
	})
}

func TestAddWrongLocation(t *testing.T) {
	t.Parallel()
	client := backupClient.Default.Locations

	t.Run("missing config", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
			},
			Context: pmmapitests.Context,
		})

		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Missing location type.")
		assert.Nil(t, resp)
	})

	t.Run("missing client config path", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:            gofakeit.Name(),
				Description:     gofakeit.Question(),
				PMMClientConfig: &locations.AddLocationParamsBodyPMMClientConfig{},
			},
			Context: pmmapitests.Context,
		})

		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field PmmClientConfig.Path: value '' must not be an empty string")
		assert.Nil(t, resp)
	})
	t.Run("missing name", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
			},
			Context: pmmapitests.Context,
		})

		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Missing location type.")
		assert.Nil(t, resp)
	})

	t.Run("missing s3 endpoint", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
				S3Config: &locations.AddLocationParamsBodyS3Config{
					AccessKey: "access_key",
					SecretKey: "secret_key",
				},
			},
			Context: pmmapitests.Context,
		})

		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field S3Config.Endpoint: value '' must not be an empty string")
		assert.Nil(t, resp)
	})
	t.Run("double config", func(t *testing.T) {
		t.Parallel()

		resp, err := client.AddLocation(&locations.AddLocationParams{
			Body: locations.AddLocationBody{
				Name:        gofakeit.Name(),
				Description: gofakeit.Question(),
				PMMClientConfig: &locations.AddLocationParamsBodyPMMClientConfig{
					Path: "/tmp",
				},
				S3Config: &locations.AddLocationParamsBodyS3Config{
					Endpoint:  "http://example.com",
					AccessKey: "access_key",
					SecretKey: "secret_key",
				},
			},
			Context: pmmapitests.Context,
		})
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Only one config is allowed.")

		assert.Nil(t, resp)

	})
}

func TestListLocations(t *testing.T) {
	t.Parallel()
	client := backupClient.Default.Locations

	body := locations.AddLocationBody{
		Name:        gofakeit.Name(),
		Description: gofakeit.Question(),
		PMMClientConfig: &locations.AddLocationParamsBodyPMMClientConfig{
			Path: "/tmp",
		},
	}
	addResp, err := client.AddLocation(&locations.AddLocationParams{
		Body:    body,
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	defer deleteLocation(t, client, addResp.Payload.LocationID)

	resp, err := client.ListLocations(&locations.ListLocationsParams{Context: pmmapitests.Context})
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Payload.Locations)
	var found bool
	for _, loc := range resp.Payload.Locations {
		if loc.LocationID == addResp.Payload.LocationID {
			assert.Equal(t, body.Name, loc.Name)
			assert.Equal(t, body.Description, loc.Description)
			assert.Equal(t, body.PMMClientConfig.Path, loc.PMMClientConfig.Path)
			found = true
		}
	}
	assert.True(t, found, "Expected location not found")
}

func deleteLocation(t *testing.T, client locations.ClientService, id string) {
	t.Helper()
	// @TODO call Delete https://jira.percona.com/browse/PMM-7383
}
