package rules_test

import (
	"bytes"
	"context"
	"fisherman/internal/rules"
	"fisherman/pkg/shell"
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
		expectedScript *shell.Script
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
			expectedScript: shell.NewScript().
				SetCommands([]string{"demo"}).
				SetEnvironmentVariables(map[string]string{
					"demo":  "demo",
					"demo2": "demo2",
				}).
				SetDirectory("~"),
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
				Set(func(c1 context.Context, w1 io.Writer, s1 string, s2 *shell.Script) error {
					fmt.Fprint(w1, tt.shellOutput)

					assert.Equal(t, tt.expectedShell, s1)
					assert.Equal(t, ctx, c1)
					if tt.expectedScript != nil {
						assert.EqualValues(t, *tt.expectedScript, *s2)
					}

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

func TestShellScript_Compile(t *testing.T) {
	rule := rules.ShellScript{
		Name:     "{{var1}}",
		Shell:    "{{var1}}",
		Commands: []string{"{{var1}}1", "{{var1}}2"},
		Env: map[string]string{
			"{{var1}}": "{{var1}}",
		},
		Dir:    "{{var1}}",
		Output: true,
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, rules.ShellScript{
		Name:     "VALUE",
		Shell:    "{{var1}}",
		Commands: []string{"VALUE1", "VALUE2"},
		Env: map[string]string{
			"{{var1}}": "VALUE",
		},
		Dir:    "VALUE",
		Output: true,
	}, rule)
}
