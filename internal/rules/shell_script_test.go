package rules_test

import (
	"bytes"
	"context"
	"fisherman/infrastructure/shell"
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellScript_Check(t *testing.T) {
	tests := []struct {
		name           string
		config         rules.ShellScript
		expectedOutput string
		expectedErr    error
		shellOutput    string
		expectedShell  string
		expectedScript shell.ShScript
	}{
		{
			name: "script with output",
			config: rules.ShellScript{
				Name:   "testScript",
				Output: true,
			},
			expectedOutput: "[testScript] test",
			expectedErr:    nil,
			shellOutput:    "test",
			expectedShell:  "",
		},
		{
			name: "script with out output",
			config: rules.ShellScript{
				Name:     "testScript",
				Output:   false,
				Commands: []string{"demo"},
				Env: map[string]string{
					"demo":  "demo",
					"demo2": "demo2",
				},
				Dir: "~",
			},
			expectedOutput: "",
			expectedErr:    nil,
			shellOutput:    "test",
			expectedShell:  "",
			expectedScript: shell.ShScript{
				Commands: []string{"demo"},
				Env: map[string]string{
					"demo":  "demo",
					"demo2": "demo2",
				},
				Dir: "~",
			},
		},
		{
			name: "script with with custom shell",
			config: rules.ShellScript{
				Name:   "zsh-script",
				Output: true,
				Shell:  "zsh",
			},
			expectedOutput: "[zsh-script] demo",
			expectedErr:    nil,
			shellOutput:    "demo",
			expectedShell:  "zsh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			ctx := mocks.NewExecutionContextMock(t)
			sh := mocks.NewShellMock(t).
				ExecMock.
				Set(func(c1 context.Context, w1 io.Writer, s1 string, s2 shell.ShScript) error {
					fmt.Fprint(w1, tt.shellOutput)

					assert.Equal(t, tt.expectedShell, s1)
					assert.Equal(t, ctx, c1)
					assert.EqualValues(t, tt.expectedScript, s2)

					return tt.expectedErr
				})

			ctx.ShellMock.Return(sh)

			err := tt.config.Check(ctx, output)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedOutput, output.String())
		})
	}
}

func TestShellScript_GetPosition(t *testing.T) {
	rule := rules.ShellScript{}

	actual := rule.GetPosition()

	assert.Equal(t, actual, rules.Scripts)
}
