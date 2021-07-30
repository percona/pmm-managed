package action

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestRunExplain(t *testing.T) {
	t.Skip("not implemented yet")

	explainActionOK, err := client.Default.Actions.StartMySQLExplainAction(&actions.StartMySQLExplainActionParams{
		Context: pmmapitests.Context,
		Body: actions.StartMySQLExplainActionBody{
			// PMMAgentID: "/agent_id/f235005b-9cca-4b73-bbbd-1251067c3138",
			ServiceID: "/service_id/5a9a7aa6-7af4-47be-817c-6d88e955bff2",
			Query:     "SELECT `t` . * FROM `test` . `key_value` `t`",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, explainActionOK.Payload.ActionID)

	time.Sleep(2 * time.Second)

	actionOK, err := client.Default.Actions.GetAction(&actions.GetActionParams{
		Context: pmmapitests.Context,
		Body: actions.GetActionBody{
			ActionID: explainActionOK.Payload.ActionID,
		},
	})
	require.NoError(t, err)
	require.Empty(t, actionOK.Payload.Error)
	t.Log(actionOK.Payload.Output)
}

func TestRunMongoDBExplain(t *testing.T) {
	// When we have an pmm-agent in dev-container and we can remove this skip, please remove the t.Logf at the end
	// of this test and replace it with a proper test that checks the results.
	t.Skip("pmm-agent in dev-container is not fully implemented yet")

	explainActionOK, err := client.Default.Actions.StartMongoDBExplainAction(&actions.StartMongoDBExplainActionParams{
		Context: pmmapitests.Context,
		Body: actions.StartMongoDBExplainActionBody{
			ServiceID: "/service_id/2402bf45-19c2-4bee-931a-307b26ed5300",
			Query:     `{"ns":"test.coll","op":"query","query":{"k":{"$lte":{"$numberInt":"1"}}}}`,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, explainActionOK.Payload.ActionID)

	var actionOK *actions.GetActionOK

	for i := 0; i < 6; i++ {
		var err error
		actionOK, err = client.Default.Actions.GetAction(&actions.GetActionParams{
			Context: pmmapitests.Context,
			Body: actions.GetActionBody{
				ActionID: explainActionOK.Payload.ActionID,
			},
		})
		require.NoError(t, err)
		require.Empty(t, actionOK.Payload.Error)

		if actionOK.Payload.Done {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	assert.True(t, actionOK.Payload.Done)

	want := map[string]interface{}{
		"winningPlan": map[string]interface{}{
			"stage": "EOF",
		},
		"rejectedPlans": []interface{}{},
		"plannerVersion": map[string]interface{}{
			"$numberInt": "1",
		},
		"namespace":      "test.coll",
		"indexFilterSet": bool(false),
		"parsedQuery": map[string]interface{}{
			"k": map[string]interface{}{
				"$lte": map[string]interface{}{
					"$numberInt": "1",
				},
			},
		},
	}
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(actionOK.Payload.Output), &m)
	assert.NoError(t, err)
	assert.Equal(t, m["queryPlanner"], want)
}
