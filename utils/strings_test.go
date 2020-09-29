package utils_test

import (
	"fisherman/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint:dupl
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "Return true for empty string", value: "", want: true},
		{name: "Return true for spaces only", value: "  ", want: true},
		{name: "Return true for tabs only", value: "\t\t", want: true},
		{name: "Return true for newline charsets only", value: "\n\n", want: true},
		{name: "Return true for mixed witespace string", value: "\t\n ", want: true},
		{name: "Return false for string without space", value: "2", want: false},
		{name: "Return false for not empty string with space", value: " 2", want: false},
		{name: "Return false for not empty string with tab", value: "\t2", want: false},
		{name: "Return false for not empty string with newline", value: "\n2", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.IsEmpty(tt.value))
		})
	}
}

// nolint:dupl
func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "Return false for empty string", value: "", expected: false},
		{name: "Return false for spaces only", value: "  ", expected: false},
		{name: "Return false for tabs only", value: "\t\t", expected: false},
		{name: "Return false for newline charsets only", value: "\n\n", expected: false},
		{name: "Return false for mixed witespace string", value: "\t\n ", expected: false},
		{name: "Return true for string without space", value: "2", expected: true},
		{name: "Return true for not empty string with space", value: " 2", expected: true},
		{name: "Return true for not empty string with tab", value: "\t2", expected: true},
		{name: "Return true for not empty string with newline", value: "\n2", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.IsNotEmpty(tt.value))
		})
	}
}

func TestOriginalOrNA(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Return 'N/A' for empty string", value: "", want: "N/A"},
		{name: "Return 'N/A' for spaces only", value: "  ", want: "N/A"},
		{name: "Return 'N/A' for tabs only", value: "\t\t", want: "N/A"},
		{name: "Return 'N/A' for newline charsets only", value: "\n\n", want: "N/A"},
		{name: "Return 'N/A' for mixed witespace string", value: "\t\n ", want: "N/A"},
		{name: "Return original string for string without space", value: "2", want: "2"},
		{name: "Return original for not empty string with space", value: " 2", want: " 2"},
		{name: "Return original for not empty string with tab", value: "\t2", want: "\t2"},
		{name: "Return original for not empty string with newline", value: "\n2", want: "\n2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.OriginalOrNA(tt.value))
		})
	}
}
