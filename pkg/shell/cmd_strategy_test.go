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

func TestCmdStrategy(t *testing.T) {
	if runtime.GOOS != windowsOS {
		t.Skip("test can not be runned on linux or darwin")
	}

	for _, tt := range scriptTests {
		t.Run(tt.name, func(t *testing.T) {
			host := NewHost(context.TODO(), Cmd())

			actual := host.Run(tt.script)

			testutils.AssertError(t, tt.expectedRrr, actual)
		})
	}
}

func TestCmdStrategy_GetName(t *testing.T) {
	actual := Cmd().GetName()

	assert.Equal(t, "cmd", actual)
}

func TestCmdStrategy_ArgsWrapper(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "empty slice",
			args:     []string{},
			expected: []string{"/Q", "/D", "/K"},
		},
		{
			name:     "additional arguments",
			args:     []string{"arg1", "arg2"},
			expected: []string{"/Q", "/D", "/K", "arg1", "arg2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Cmd().ArgsWrapper(tt.args)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCmdStrategy_EnvWrapper(t *testing.T) {
	for _, tt := range envTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Cmd().EnvWrapper(tt.env)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCmdStrategy_GetEncoding(t *testing.T) {
	assert.NotPanics(t, func() {
		actual := Cmd().GetEncoding()

		assert.NotNil(t, actual)
	})
}
