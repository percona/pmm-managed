package inventory

import (
	"testing"

	"github.com/percona/pmm/api/inventory/json/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/Percona-Lab/pmm-api-tests" // init default client
)

func TestNodes(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		res, err := client.Default.Nodes.ListNodes(nil)
		require.NoError(t, err)
		assert.Len(t, res.Payload.Generic, 1)
	})
}
