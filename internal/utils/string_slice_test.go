package utils_test

import (
	. "fisherman/internal/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		source   []string
		expected []string
	}{
		{name: "", source: []string{"#1", "2", "#3", "4"}, expected: []string{"#1", "#3"}},
		{name: "", source: []string{"1", "2", "3", "4"}, expected: []string{}},
		{name: "", source: []string{"#1", "#2", "#3", "#4"}, expected: []string{"#1", "#2", "#3", "#4"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := Filter(tt.source, func(value string) bool {
				return strings.Contains(value, "#")
			})

			assert.Equal(t, tt.expected, filtered)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name       string
		collection []string
		value      string
		expected   bool
	}{
		{
			name:       "empty slice",
			collection: []string{},
			value:      "demo",
			expected:   false,
		},
		{
			name:       "single value slice",
			collection: []string{"demo"},
			value:      "demo",
			expected:   true,
		},
		{
			name:       "slice with target value",
			collection: []string{"value1", "value2", "value3", "value4"},
			value:      "value3",
			expected:   true,
		},
		{
			name:       "slice withot target value",
			collection: []string{"value1", "value2", "value3", "value4"},
			value:      "value0",
			expected:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Contains(tt.collection, tt.value)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
