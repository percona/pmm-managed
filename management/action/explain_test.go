package action

import (
	"fmt"
	"testing"
	"time"

	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/actions"
	"github.com/stretchr/testify/require"

	pmmapitests "github.com/Percona-Lab/pmm-api-tests"
)

func TestRunExplain(t *testing.T) {
	t.Skip("not implemented yet")

	explainActionOK, err := client.Default.Actions.StartMySQLExplainAction(&actions.StartMySQLExplainActionParams{
		Context: pmmapitests.Context,
		Body: actions.StartMySQLExplainActionBody{
			//PMMAgentID: "/agent_id/f235005b-9cca-4b73-bbbd-1251067c3138",
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
	fmt.Println(actionOK.Payload.Output)
}
