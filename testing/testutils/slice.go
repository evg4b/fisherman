package testutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertSlice(t *testing.T, expected, actual []string) {
	t.Helper()

	assert.Equal(t, len(expected), len(actual))
	for _, value := range expected {
		assert.Contains(t, actual, value, fmt.Sprintf("slice should contains value '%s'", value))
	}
}
