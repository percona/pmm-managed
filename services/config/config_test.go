package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadFailOnDuplicatedId(t *testing.T) {
	err := os.Setenv("PERCONA_PMM_CONFIG_PATH", "pmm-managed-test-duplicated-id.yaml")
	require.NoError(t, err)

	s := NewService()
	err = s.Load()
	require.Error(t, err, "Duplicated id [NodeLoad1min] found in the config file [pmm-managed-test-duplicated-id.yaml]")
}
