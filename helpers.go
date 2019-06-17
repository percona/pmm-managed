package pmmapitests

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/nodes"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestString returns semi-random string that can be used as a test data.
func TestString(t *testing.T, name string) string {
	t.Helper()

	n := rand.Int() //nolint:gosec
	return fmt.Sprintf("pmm-api-tests/%s/%s/%s/%d", Hostname, t.Name(), name, n)
}

type ErrorResponse interface {
	Code() int
}

type ServerResponse struct {
	Code  int
	Error string
}

func AssertEqualAPIError(t require.TestingT, err error, expected ServerResponse) bool {
	if n, ok := t.(interface {
		Helper()
	}); ok {
		n.Helper()
	}

	if !assert.Error(t, err) {
		return false
	}

	require.Implementsf(t, new(ErrorResponse), err, "Wrong response type. Expected %T, got %T.\nError message: %v", new(ErrorResponse), err, err)

	assert.Equal(t, expected.Code, err.(ErrorResponse).Code())

	// Have to use reflect because there are a lot of types with the same structure and different names.
	val := reflect.ValueOf(err)

	payload := val.Elem().FieldByName("Payload")
	require.True(t, payload.IsValid(), "Wrong response structure. There is no field Payload.")

	errorField := payload.Elem().FieldByName("Error")
	require.True(t, errorField.IsValid(), "Wrong response structure. There is no field Error in Payload.")

	return assert.Equal(t, expected.Error, errorField.String())
}

func ExpectFailure(t *testing.T, link string) (failureTestingT *expectedFailureTestingT) {
	failureTestingT = &expectedFailureTestingT{
		t:    t,
		link: link,
	}
	return failureTestingT
}

// expectedFailureTestingT expects that test will fail.
// if test is failed we skip it
// if it doesn't we call Fail
type expectedFailureTestingT struct {
	t      *testing.T
	errors []string
	failed bool
	link   string
}

func (tt *expectedFailureTestingT) Helper()      { tt.t.Helper() }
func (tt *expectedFailureTestingT) Name() string { return tt.t.Name() }

func (tt *expectedFailureTestingT) FailNow() {
	tt.failed = true

	// We have to set unexported testing.T.finished = true to make everything work,
	// but we can't call tt.t.FailNow() as it calls Fail().
	tt.t.SkipNow()
}

func (tt *expectedFailureTestingT) Errorf(format string, args ...interface{}) {
	tt.errors = append(tt.errors, fmt.Sprintf(format, args...))
	tt.failed = true
}

func (tt *expectedFailureTestingT) Check() {
	tt.t.Helper()

	if tt.failed {
		for _, v := range tt.errors {
			tt.t.Log(v)
		}
		tt.t.Skipf("Expected failure: %s", tt.link)
		return
	}

	tt.t.Fatalf("%s expected to fail, but didn't: %s", tt.Name(), tt.link)
}

func RemoveNodes(t *testing.T, nodeIDs ...string) {
	t.Helper()

	for _, nodeID := range nodeIDs {
		params := &nodes.RemoveNodeParams{
			Body: nodes.RemoveNodeBody{
				NodeID: nodeID,
			},
			Context: context.Background(),
		}
		res, err := client.Default.Nodes.RemoveNode(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

func RemoveServices(t *testing.T, serviceIDs ...string) {
	t.Helper()

	for _, serviceID := range serviceIDs {
		params := &services.RemoveServiceParams{
			Body: services.RemoveServiceBody{
				ServiceID: serviceID,
				//Force:     true,
			},
			Context: context.Background(),
		}
		res, err := client.Default.Services.RemoveService(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

func RemoveAgents(t *testing.T, agentIDs ...string) {
	t.Helper()

	for _, agentID := range agentIDs {
		params := &agents.RemoveAgentParams{
			Body: agents.RemoveAgentBody{
				AgentID: agentID,
			},
			Context: context.Background(),
		}
		res, err := client.Default.Agents.RemoveAgent(params)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	}
}

// check interfaces
var (
	_ assert.TestingT  = (*expectedFailureTestingT)(nil)
	_ require.TestingT = (*expectedFailureTestingT)(nil)
)
