package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// AssertError error by string, where expected string is empty then error should be null.
// Otherwise error message should be matched to string.
func AssertError(t *testing.T, expected string, actual error) {
	t.Helper()

	if len(expected) > 0 {
		require.EqualError(t, actual, expected)
	} else {
		require.NoError(t, actual)
	}
}
