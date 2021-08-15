package utils_test

import (
	"fisherman/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint:dupl
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "Return true for empty string", value: "", expected: true},
		{name: "Return true for spaces only", value: "  ", expected: true},
		{name: "Return true for tabs only", value: "\t\t", expected: true},
		{name: "Return true for newline charsets only", value: "\n\n", expected: true},
		{name: "Return true for mixed witespace string", value: "\t\n ", expected: true},
		{name: "Return false for string without space", value: "2", expected: false},
		{name: "Return false for not empty string with space", value: " 2", expected: false},
		{name: "Return false for not empty string with tab", value: "\t2", expected: false},
		{name: "Return false for not empty string with newline", value: "\n2", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.IsEmpty(tt.value))
		})
	}
}

func TestOriginalOrNA(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "Return 'N/A' for empty string", value: "", expected: "N/A"},
		{name: "Return 'N/A' for spaces only", value: "  ", expected: "N/A"},
		{name: "Return 'N/A' for tabs only", value: "\t\t", expected: "N/A"},
		{name: "Return 'N/A' for newline charsets only", value: "\n\n", expected: "N/A"},
		{name: "Return 'N/A' for mixed witespace string", value: "\t\n ", expected: "N/A"},
		{name: "Return original string for string without space", value: "2", expected: "2"},
		{name: "Return original for not empty string with space", value: " 2", expected: " 2"},
		{name: "Return original for not empty string with tab", value: "\t2", expected: "\t2"},
		{name: "Return original for not empty string with newline", value: "\n2", expected: "\n2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.OriginalOrNA(tt.value))
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	defaultValue := "default-value"

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "Return default value for empty string", value: "", expected: defaultValue},
		{name: "Return default value for spaces only", value: "  ", expected: defaultValue},
		{name: "Return default value for tabs only", value: "\t\t", expected: defaultValue},
		{name: "Return default value for newline charsets only", value: "\n\n", expected: defaultValue},
		{name: "Return default value for mixed witespace string", value: "\t\n ", expected: defaultValue},
		{name: "Return original string for string without space", value: "2", expected: "2"},
		{name: "Return original for not empty string with space", value: " 2", expected: " 2"},
		{name: "Return original for not empty string with tab", value: "\t2", expected: "\t2"},
		{name: "Return original for not empty string with newline", value: "\n2", expected: "\n2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.GetOrDefault(tt.value, defaultValue))
		})
	}
}
