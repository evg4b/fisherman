// nolint: dupl
package shell_test

import (
	"context"
	. "fisherman/pkg/shell"
	"fisherman/testing/testutils"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPowerShellStrategy(t *testing.T) {
	if runtime.GOOS != windowsOS {
		t.Skip("test can not be runned on linux or darwin")
	}

	for _, tt := range scriptTests {
		t.Run(tt.name, func(t *testing.T) {
			host := NewHost(context.TODO(), PowerShell())

			actual := host.Run(tt.script)

			testutils.AssertError(t, tt.expectedRrr, actual)
		})
	}
}

func TestPowerShellStrategy_GetName(t *testing.T) {
	actual := PowerShell().GetName()

	assert.Equal(t, "powershell", actual)
}

func TestPowerShellStrategy_ArgsWrapper(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "empty slice",
			args:     []string{},
			expected: []string{"-NoProfile", "-NonInteractive", "-NoLogo"},
		},
		{
			name:     "additional arguments",
			args:     []string{"arg1", "arg2"},
			expected: []string{"-NoProfile", "-NonInteractive", "-NoLogo", "arg1", "arg2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := PowerShell().ArgsWrapper(tt.args)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPowerShellStrategy_EnvWrapper(t *testing.T) {
	for _, tt := range envTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := PowerShell().EnvWrapper(tt.env)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPowerShellStrategy_GetEncoding(t *testing.T) {
	assert.NotPanics(t, func() {
		actual := PowerShell().GetEncoding()

		assert.NotNil(t, actual)
	})
}
