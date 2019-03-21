package pmmapitests

import (
	"fmt"
	"math/rand"
	"testing"
)

// TestString returns semi-random string that can be used as a test data.
func TestString(t *testing.T, name string) string {
	t.Helper()

	n := rand.Int() //nolint:gosec
	return fmt.Sprintf("pmm-api-tests/%s/%s/%s/%d", Hostname, t.Name(), name, n)
}
