package expression_test

import (
	"github.com/evg4b/fisherman/internal/constants"
	"testing"

	. "github.com/evg4b/fisherman/internal/expression"

	"github.com/stretchr/testify/assert"
)

func TestEnvVars_IsEmpty(t *testing.T) {
	var tests []struct {
		name     string
		vars     EnvVars
		value    string
		expected bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.vars.IsEmpty(tt.value)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestEnvVars_IsWindows(t *testing.T) {
	tests := []struct {
		name     string
		vars     EnvVars
		value    string
		expected bool
	}{
		{
			name:     "where current os is windows",
			vars:     os(Windows),
			expected: true,
		},
		{
			name:     "where current os is linux",
			vars:     os(Linux),
			expected: false,
		},
		{
			name:     "where current os is macos",
			vars:     os(Macos),
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.vars.IsWindows()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestEnvVars_IsLinux(t *testing.T) {
	tests := []struct {
		name     string
		vars     EnvVars
		expected bool
	}{
		{
			name:     "where current os is windows",
			vars:     os(Windows),
			expected: false,
		},
		{
			name:     "where current os is linux",
			vars:     os(Linux),
			expected: true,
		},
		{
			name:     "where current os is macos",
			vars:     os(Macos),
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.vars.IsLinux()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestEnvVars_IsMacOs(t *testing.T) {
	tests := []struct {
		name     string
		vars     EnvVars
		expected bool
	}{
		{
			name:     "where current os is windows",
			vars:     os(Windows),
			expected: false,
		},
		{
			name:     "where current os is linux",
			vars:     os(Linux),
			expected: false,
		},
		{
			name:     "where current os is macos",
			vars:     os(Macos),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.vars.IsMacOs()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func os(os string) EnvVars {
	return EnvVars{constants.OsVariable: os}
}
