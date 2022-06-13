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

package apitestsunmocked

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httptransport "github.com/go-openapi/runtime/client"
	inventoryClient "github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/agents"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	serverClient "github.com/percona/pmm/api/serverpb/json/client"
	"github.com/percona/pmm/utils/tlsconfig"
)

var (
	connAgent string
	nodeID    string
)

// ErrFromNginx is an error type for nginx HTML response.
type ErrFromNginx string

// Error implements error interface.
func (e *ErrFromNginx) Error() string {
	return "response from nginx: " + string(*e)
}

func TestSuccessfulConnection(t *testing.T) {
	username := "pmm-agent"
	password := "pmm-agent-password"
	require.True(t, establishPostgresConnection(t, username, password))
}

func TestUnsuccessfulConnection(t *testing.T) {
	username := "invalid-pmm-agent"
	password := "invalid-pmm-agent-password"
	require.False(t, establishPostgresConnection(t, username, password))
}

func establishPostgresConnection(t *testing.T, username, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service := addPostgreSQLService(services.AddPostgreSQLServiceBody{
		NodeID:      nodeID,
		Address:     "pmm-agent_postgres",
		Port:        5432,
		ServiceName: "Test PostgreSQL",
	}, ctx)
	if service == nil {
		return false
	}
	defer removeService(t, service.Postgresql.ServiceID)

	exporter := addPostgresExporter(agents.AddPostgresExporterBody{
		ServiceID:  service.Postgresql.ServiceID,
		Username:   username,
		Password:   password,
		PMMAgentID: connAgent,
		CustomLabels: map[string]string{
			"custom_label_postgres_exporter": "test_postgres_exporter",
		},
		SkipConnectionCheck: false,
	}, ctx)
	if exporter == nil {
		return false
	}

	agentID := exporter.PostgresExporter.AgentID
	if agentID == "" {
		return false
	}
	defer removeAgent(t, agentID)

	return true
}

func getBaseURL(serverURL string) (*url.URL, error) {
	baseURL, err := url.Parse(serverURL)

	if err != nil {
		logrus.Warnf("Failed to parse PMM Server URL: %s.", err)
		return nil, err
	}

	if baseURL.Host == "" || baseURL.Scheme == "" {
		logrus.Warnf("Invalid PMM Server URL: %s", baseURL.String())
		return nil, err
	}

	if baseURL.Path == "" {
		baseURL.Path = "/"
	}

	logrus.Debugf("PMM Server URL: %s.", baseURL)
	return baseURL, nil
}

func initClients(serverURL string) error {
	url, err := getBaseURL(serverURL)

	if err != nil {
		return err
	}

	transport := Transport(url, true)
	inventoryClient.Default = inventoryClient.New(transport, nil)

	return nil
}

func init() {

	var serverURL string

	flag.StringVar(&serverURL, "pmm.server-url", "https://admin:admin@localhost/", "PMM Server URL [PMM_SERVER_URL].")
	testing.Init()
	flag.Parse()

	for envVar, f := range map[string]*flag.Flag{
		"PMM_SERVER_URL": flag.Lookup("pmm.server-url"),
	} {
		env, ok := os.LookupEnv(envVar)
		if ok {
			err := f.Value.Set(env)
			if err != nil {
				logrus.Fatalf("Invalid ENV variable %s: %s", envVar, env)
			}
		}
	}

	err := initClients(serverURL)
	if err != nil {
		logrus.Fatalf("Failed to initialized clients: %s", err.Error())
	}

	// do not run tests if server is not available
	_, err = serverClient.Default.Server.Readiness(nil)
	if err != nil {
		panic(err)
	}

	connAgent, nodeID = getPMMAgentForDBConnection()
	logrus.Infof("Connection agent: %v", connAgent)
	logrus.Infof("Node ID: %v", nodeID)
}

// Transport returns configured Swagger transport for given URL.
func Transport(baseURL *url.URL, insecureTLS bool) *httptransport.Runtime {
	transport := httptransport.New(baseURL.Host, baseURL.Path, []string{baseURL.Scheme})
	if u := baseURL.User; u != nil {
		password, _ := u.Password()
		transport.DefaultAuthentication = httptransport.BasicAuth(u.Username(), password)
	}
	transport.SetLogger(logrus.WithField("component", "client"))
	transport.SetDebug(logrus.GetLevel() >= logrus.DebugLevel)
	transport.Context = context.Background() // not Context - do not cancel the whole transport

	// set error handlers for nginx responses if pmm-managed is down
	errorConsumer := runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		b, _ := ioutil.ReadAll(reader)
		err := ErrFromNginx(string(b))
		return &err
	})
	transport.Consumers = map[string]runtime.Consumer{
		runtime.JSONMime:    runtime.JSONConsumer(),
		runtime.HTMLMime:    errorConsumer,
		runtime.TextMime:    errorConsumer,
		runtime.DefaultMime: errorConsumer,
	}

	// disable HTTP/2, set TLS config
	httpTransport := transport.Transport.(*http.Transport)
	httpTransport.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	if baseURL.Scheme == "https" {
		httpTransport.TLSClientConfig = tlsconfig.Get()
		httpTransport.TLSClientConfig.ServerName = baseURL.Hostname()
		httpTransport.TLSClientConfig.InsecureSkipVerify = insecureTLS
	}

	return transport
}

func getPMMAgentForDBConnection() (string, string) {
	agentListResponse, err := inventoryClient.Default.Agents.ListAgents(nil)

	if err != nil {
		logrus.Fatalf("Failed to get list of agents: %s", err.Error())
	}

	if agentListResponse == nil {
		logrus.Fatal("Agent response is nil")
	}

	agentCount := len(agentListResponse.Payload.PMMAgent)
	if agentCount != 2 {
		logrus.Fatalf("Expected only 2 PMM agents, found %d", agentCount)
	}

	for _, pmmAgent := range agentListResponse.Payload.PMMAgent {
		if pmmAgent.AgentID == "pmm-server" {
			continue
		}

		return pmmAgent.AgentID, pmmAgent.RunsOnNodeID
	}

	return "", ""
}

// Register Postgres service with PMM server
func addPostgreSQLService(body services.AddPostgreSQLServiceBody, ctx context.Context) *services.AddPostgreSQLServiceOKBody {
	params := &services.AddPostgreSQLServiceParams{
		Body:    body,
		Context: ctx,
	}

	res, err := inventoryClient.Default.Services.AddPostgreSQLService(params)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return nil
	}

	return res.Payload
}

// Register Postgres exporter with PMM server
func addPostgresExporter(body agents.AddPostgresExporterBody, ctx context.Context) *agents.AddPostgresExporterOKBody {
	params := agents.AddPostgresExporterParams{
		Body:    body,
		Context: ctx,
	}

	res, err := inventoryClient.Default.Agents.AddPostgresExporter(&params)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return nil
	}

	return res.Payload
}

// Deregister a service with PMM server
func removeService(t *testing.T, serviceID string) {
	params := &services.RemoveServiceParams{
		Body: services.RemoveServiceBody{
			ServiceID: serviceID,
			Force:     true,
		},
		Context: context.Background(),
	}

	res, err := inventoryClient.Default.Services.RemoveService(params)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

// Deregister an agent with PMM server
func removeAgent(t *testing.T, agentID string) {
	params := &agents.RemoveAgentParams{
		Body: agents.RemoveAgentBody{
			AgentID: agentID,
			Force:   true,
		},
		Context: context.Background(),
	}

	res, err := inventoryClient.Default.Agents.RemoveAgent(params)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
