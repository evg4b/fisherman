package validation_test

import (
	"fisherman/internal/validation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "Message with prefix",
			err:      validation.Errorf("prefix1", "message"),
			expected: "[prefix1] message",
		},
		{
			name:     "Message with prefix and arguments",
			err:      validation.Errorf("prefix2", "message %s %d", "test", 13),
			expected: "[prefix2] message test 13",
		},
		{
			name:     "Message without prefix",
			err:      validation.Errorf("", "message"),
			expected: "message",
		},
		{
			name:     "Message with prefix and arguments",
			err:      validation.Errorf("", "message %s %d", "test", 13),
			expected: "message test 13",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.err.Error()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
