// nolint: dupl

package shell_test

import (
	"context"
	"fisherman/testing/testutils"
	"runtime"
	"testing"

	. "fisherman/pkg/shell"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding"
)

func TestBashStrategy(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skip("test can not be runned on windows")
	}

	for _, tt := range scriptTests {
		t.Run(tt.name, func(t *testing.T) {
			host := NewHost(context.TODO(), Bash())

			actual := host.Run(tt.script)

			testutils.AssertError(t, tt.expectedRrr, actual)
		})
	}
}

func TestBashStrategy_GetName(t *testing.T) {
	actual := Bash().GetName()

	assert.Equal(t, "bash", actual)
}

func TestBashStrategy_ArgsWrapper(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "empty slice",
			args:     []string{},
			expected: []string{},
		},
		{
			name:     "additional arguments",
			args:     []string{"arg1", "arg2"},
			expected: []string{"arg1", "arg2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Bash().ArgsWrapper(tt.args)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBashStrategy_EnvWrapper(t *testing.T) {
	for _, tt := range envTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Bash().EnvWrapper(tt.env)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBashStrategy_GetEncoding(t *testing.T) {
	actual := Bash().GetEncoding()

	assert.EqualValues(t, encoding.Nop, actual)
}
