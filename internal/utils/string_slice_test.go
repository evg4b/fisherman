package utils

import (
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
