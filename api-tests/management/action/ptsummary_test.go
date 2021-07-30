package action

import (
	"context"
	"testing"
	"time"

	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/actions"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestPTSummary(t *testing.T) {
	ctx, cancel := context.WithTimeout(pmmapitests.Context, 30*time.Second)
	defer cancel()

	explainActionOK, err := client.Default.Actions.StartPTSummaryAction(&actions.StartPTSummaryActionParams{
		Context: ctx,
		Body: actions.StartPTSummaryActionBody{
			NodeID: "pmm-server",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, explainActionOK.Payload.ActionID)

	for {
		actionOK, err := client.Default.Actions.GetAction(&actions.GetActionParams{
			Context: ctx,
			Body: actions.GetActionBody{
				ActionID: explainActionOK.Payload.ActionID,
			},
		})
		require.NoError(t, err)

		if !actionOK.Payload.Done {
			time.Sleep(1 * time.Second)

			continue
		}

		require.True(t, actionOK.Payload.Done)
		require.Empty(t, actionOK.Payload.Error)
		require.NotEmpty(t, actionOK.Payload.Output)
		t.Log(actionOK.Payload.Output)

		break
	}
}
