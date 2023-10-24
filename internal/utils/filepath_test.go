package utils_test

import (
	"fisherman/testing/testutils"
	"testing"

	. "fisherman/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestMatchToGlobs(t *testing.T) {
	tests := []struct {
		name        string
		globs       []string
		file        string
		expected    bool
		expectedErr string
	}{
		{
			name:     "empty",
			globs:    []string{},
			file:     "",
			expected: false,
		},
		{
			name:     "direct path matched",
			globs:    []string{"some/path.json"},
			file:     "some/path.json",
			expected: true,
		},
		{
			name:     "direct path not matched",
			globs:    []string{"some/path.json"},
			file:     "some/path.json",
			expected: true,
		},
		{
			name:     "multi matches",
			globs:    []string{"some/path.json", "some/*.json"},
			file:     "some/path.json",
			expected: true,
		},
		{
			name:     "matched with glob",
			globs:    []string{"some/*.json"},
			file:     "some/path.json",
			expected: true,
		},
		{
			name:     "not matched with glob",
			globs:    []string{"some/*.md"},
			file:     "some/path.json",
			expected: false,
		},
		{
			name:        "invalid pattern",
			globs:       []string{"some/[*"},
			file:        "some/path.json",
			expected:    false,
			expectedErr: "syntax error in pattern",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := MatchToGlobs(tt.globs, tt.file)

			testutils.AssertError(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
