package management

import (
	"testing"

	inventoryClient "github.com/percona/pmm/api/inventorypb/json/client"
	"github.com/percona/pmm/api/inventorypb/json/client/nodes"
	"github.com/percona/pmm/api/inventorypb/json/client/services"
	"github.com/percona/pmm/api/managementpb/json/client"
	"github.com/percona/pmm/api/managementpb/json/client/annotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pmmapitests "github.com/percona/pmm-managed/api-tests"
)

func TestAddAnnotation(t *testing.T) {
	t.Run("Add Basic Annotation", func(t *testing.T) {
		params := &annotation.AddAnnotationParams{
			Body: annotation.AddAnnotationBody{
				Text: "Annotation Text",
				Tags: []string{"tag1", "tag2"},
			},
			Context: pmmapitests.Context,
		}
		_, err := client.Default.Annotation.AddAnnotation(params)
		require.NoError(t, err)
	})

	t.Run("Add Empty Annotation", func(t *testing.T) {
		params := &annotation.AddAnnotationParams{
			Body: annotation.AddAnnotationBody{
				Text: "",
				Tags: []string{},
			},
			Context: pmmapitests.Context,
		}
		_, err := client.Default.Annotation.AddAnnotation(params)
		pmmapitests.AssertAPIErrorf(t, err, 400, codes.InvalidArgument, "invalid field Text: value '' must not be an empty string")
	})

	t.Run("Non-existing service", func(t *testing.T) {
		params := &annotation.AddAnnotationParams{
			Body: annotation.AddAnnotationBody{
				Text:         "Some text",
				ServiceNames: []string{"no-service"},
			},
			Context: pmmapitests.Context,
		}
		_, err := client.Default.Annotation.AddAnnotation(params)
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, `Service with name "no-service" not found.`)
	})

	t.Run("Non-existing node", func(t *testing.T) {
		params := &annotation.AddAnnotationParams{
			Body: annotation.AddAnnotationBody{
				Text:     "Some text",
				NodeName: "no-node",
			},
			Context: pmmapitests.Context,
		}
		_, err := client.Default.Annotation.AddAnnotation(params)
		pmmapitests.AssertAPIErrorf(t, err, 404, codes.NotFound, `Node with name "no-node" not found.`)
	})

	t.Run("Existing service", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "annotation-node")
		paramsNode := &nodes.AddGenericNodeParams{
			Body: nodes.AddGenericNodeBody{
				NodeName: nodeName,
				Address:  "10.0.0.1",
			},
			Context: pmmapitests.Context,
		}
		resNode, err := inventoryClient.Default.Nodes.AddGenericNode(paramsNode)
		assert.NoError(t, err)
		genericNodeID := resNode.Payload.Generic.NodeID
		defer pmmapitests.RemoveNodes(t, genericNodeID)

		serviceName := pmmapitests.TestString(t, "annotation-service")
		paramsService := &services.AddMySQLServiceParams{
			Body: services.AddMySQLServiceBody{
				NodeID:      genericNodeID,
				Address:     "localhost",
				Port:        3306,
				ServiceName: serviceName,
			},
			Context: pmmapitests.Context,
		}
		resService, err := inventoryClient.Default.Services.AddMySQLService(paramsService)
		assert.NoError(t, err)
		require.NotNil(t, resService)
		serviceID := resService.Payload.Mysql.ServiceID
		defer pmmapitests.RemoveServices(t, serviceID)

		paramsAdd := &annotation.AddAnnotationParams{
			Body: annotation.AddAnnotationBody{
				Text:         "Some text",
				ServiceNames: []string{serviceName},
			},
			Context: pmmapitests.Context,
		}
		_, err = client.Default.Annotation.AddAnnotation(paramsAdd)
		require.NoError(t, err)
	})

	t.Run("Existing node", func(t *testing.T) {
		nodeName := pmmapitests.TestString(t, "annotation-node")
		params := &nodes.AddGenericNodeParams{
			Body: nodes.AddGenericNodeBody{
				NodeName: nodeName,
				Address:  "10.0.0.1",
			},
			Context: pmmapitests.Context,
		}
		res, err := inventoryClient.Default.Nodes.AddGenericNode(params)
		assert.NoError(t, err)
		defer pmmapitests.RemoveNodes(t, res.Payload.Generic.NodeID)

		paramsAdd := &annotation.AddAnnotationParams{
			Body: annotation.AddAnnotationBody{
				Text:     "Some text",
				NodeName: nodeName,
			},
			Context: pmmapitests.Context,
		}
		_, err = client.Default.Annotation.AddAnnotation(paramsAdd)
		require.NoError(t, err)
	})
}
