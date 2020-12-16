package ia

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	channelsClient "github.com/percona/pmm/api/managementpb/ia/json/client"
	"github.com/percona/pmm/api/managementpb/ia/json/client/channels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestAddChannel(t *testing.T) {
	client := channelsClient.Default.Channels

	t.Run("normal", func(t *testing.T) {
		resp, err := client.AddChannel(&channels.AddChannelParams{
			Body: channels.AddChannelBody{
				Summary:  gofakeit.Quote(),
				Disabled: gofakeit.Bool(),
				EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
					SendResolved: false,
					To:           []string{gofakeit.Email()},
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteChannel(t, client, resp.Payload.ChannelID)

		assert.NotEmpty(t, resp.Payload.ChannelID)
	})

	t.Run("invalid request", func(t *testing.T) {
		resp, err := client.AddChannel(&channels.AddChannelParams{
			Body: channels.AddChannelBody{
				Summary:  gofakeit.Quote(),
				Disabled: gofakeit.Bool(),
				EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
					SendResolved: false,
				},
			},
			Context: pmmapitests.Context,
		})

		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field EmailConfig.To: value '[]' must contain at least 1 elements")
		assert.Nil(t, resp)
	})

	t.Run("missing config", func(t *testing.T) {
		resp, err := client.AddChannel(&channels.AddChannelParams{
			Body: channels.AddChannelBody{
				Summary:  gofakeit.Quote(),
				Disabled: gofakeit.Bool(),
			},
			Context: pmmapitests.Context,
		})

		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "Missing channel configuration.")
		assert.Nil(t, resp)
	})
}

func TestChangeChannel(t *testing.T) {
	client := channelsClient.Default.Channels

	t.Run("normal", func(t *testing.T) {
		resp1, err := client.AddChannel(&channels.AddChannelParams{
			Body: channels.AddChannelBody{
				Summary:  gofakeit.Quote(),
				Disabled: false,
				EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
					SendResolved: false,
					To:           []string{gofakeit.Email()},
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)
		defer deleteChannel(t, client, resp1.Payload.ChannelID)

		newEmail := []string{gofakeit.Email()}
		_, err = client.ChangeChannel(&channels.ChangeChannelParams{
			Body: channels.ChangeChannelBody{
				ChannelID: resp1.Payload.ChannelID,
				Disabled:  true,
				EmailConfig: &channels.ChangeChannelParamsBodyEmailConfig{
					SendResolved: true,
					To:           newEmail,
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		resp2, err := client.ListChannels(&channels.ListChannelsParams{Context: pmmapitests.Context})
		require.NoError(t, err)

		assert.NotEmpty(t, resp2.Payload.Channels)
		var found bool
		for _, channel := range resp2.Payload.Channels {
			if channel.ChannelID == resp1.Payload.ChannelID {
				assert.True(t, channel.Disabled)
				assert.Equal(t, newEmail, channel.EmailConfig.To)
				assert.True(t, channel.EmailConfig.SendResolved)
				found = true
			}
		}

		assert.True(t, found, "Expected channel not found")
	})
}

func TestRemoveChannel(t *testing.T) {
	client := channelsClient.Default.Channels

	t.Run("normal", func(t *testing.T) {
		summary := gofakeit.UUID()
		resp1, err := client.AddChannel(&channels.AddChannelParams{
			Body: channels.AddChannelBody{
				Summary:  summary,
				Disabled: gofakeit.Bool(),
				EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
					SendResolved: false,
					To:           []string{gofakeit.Email()},
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		_, err = client.RemoveChannel(&channels.RemoveChannelParams{
			Body: channels.RemoveChannelBody{
				ChannelID: resp1.Payload.ChannelID,
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		resp2, err := client.ListChannels(&channels.ListChannelsParams{Context: pmmapitests.Context})
		require.NoError(t, err)

		for _, channel := range resp2.Payload.Channels {
			assert.NotEqual(t, resp1, channel.ChannelID)
		}
	})
	t.Run("unknown id", func(t *testing.T) {
		_, err := client.AddChannel(&channels.AddChannelParams{
			Body: channels.AddChannelBody{
				Summary:  gofakeit.Quote(),
				Disabled: gofakeit.Bool(),
				EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
					SendResolved: false,
					To:           []string{gofakeit.Email()},
				},
			},
			Context: pmmapitests.Context,
		})
		require.NoError(t, err)

		_, err = client.RemoveChannel(&channels.RemoveChannelParams{
			Body: channels.RemoveChannelBody{
				ChannelID: gofakeit.UUID(),
			},
			Context: pmmapitests.Context,
		})
		require.Error(t, err)
	})
}

func TestListChannels(t *testing.T) {
	client := channelsClient.Default.Channels

	summary := gofakeit.UUID()
	email := gofakeit.Email()
	disabled := gofakeit.Bool()
	resp1, err := client.AddChannel(&channels.AddChannelParams{
		Body: channels.AddChannelBody{
			Summary:  summary,
			Disabled: disabled,
			EmailConfig: &channels.AddChannelParamsBodyEmailConfig{
				SendResolved: true,
				To:           []string{email},
			},
		},
		Context: pmmapitests.Context,
	})
	require.NoError(t, err)
	defer deleteChannel(t, client, resp1.Payload.ChannelID)

	resp, err := client.ListChannels(&channels.ListChannelsParams{Context: pmmapitests.Context})
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Payload.Channels)
	var found bool
	for _, channel := range resp.Payload.Channels {
		if channel.ChannelID == resp1.Payload.ChannelID {
			assert.Equal(t, summary, channel.Summary)
			assert.Equal(t, disabled, channel.Disabled)
			assert.Equal(t, []string{email}, channel.EmailConfig.To)
			assert.True(t, channel.EmailConfig.SendResolved)
			found = true
		}
	}

	assert.True(t, found, "Expected channel not found")
}

func deleteChannel(t *testing.T, client channels.ClientService, id string) {
	_, err := client.RemoveChannel(&channels.RemoveChannelParams{
		Body: channels.RemoveChannelBody{
			ChannelID: id,
		},
		Context: pmmapitests.Context,
	})
	assert.NoError(t, err)
}
