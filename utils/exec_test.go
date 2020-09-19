package utils_test

import (
	"fisherman/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCommandExists(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected bool
	}{
		{name: "ping", cmd: "ping", expected: true},
		{name: "ping", cmd: "hot-exist-command", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := utils.IsCommandExists(tt.cmd)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
