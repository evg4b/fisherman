package rules_test

import (
	"bytes"
	"context"
	"errors"
	"fisherman/internal/rules"
	"fisherman/pkg/shell"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellScript_Check(t *testing.T) {
	baseRule := rules.BaseRule{Type: rules.ShellScriptType}
	tests := []struct {
		name           string
		config         rules.ShellScript
		execErr        error
		expectedOutput string
		expectedErr    string
		shellOutput    string
		expectedShell  string
		expectedScript *shell.Script
	}{
		{
			name: "script with output",
			config: rules.ShellScript{
				BaseRule: baseRule,
				BaseShell: rules.BaseShell{
					Name:   "testScript",
					Output: true,
				},
			},
			expectedOutput: "test",
			shellOutput:    "test",
			expectedShell:  "",
		},
		{
			name: "script with error",
			config: rules.ShellScript{
				BaseRule: baseRule,
				BaseShell: rules.BaseShell{
					Name:   "testScript",
					Output: true,
				},
			},
			execErr:     errors.New("execution error"),
			expectedErr: "failed to exec shell script: execution error",
		},
		{
			name: "script with exec.ExitError",
			config: rules.ShellScript{
				BaseRule: baseRule,
				BaseShell: rules.BaseShell{
					Name:   "testScript",
					Output: true,
				},
			},
			execErr: &exec.ExitError{
				ProcessState: &os.ProcessState{},
			},
			expectedErr: "[shell-script] script finished with exit code 0",
		},
		{
			name: "script with out output",
			config: rules.ShellScript{
				BaseRule: baseRule,
				BaseShell: rules.BaseShell{
					Name:     "testScript",
					Output:   false,
					Commands: []string{"demo"},
					Env: map[string]string{
						"demo":  "demo",
						"demo2": "demo2",
					},
					Dir: "~",
				},
			},
			expectedOutput: "",
			shellOutput:    "test",
			expectedShell:  "",
			expectedScript: shell.NewScript([]string{"demo"}).
				SetEnvironmentVariables(map[string]string{
					"demo":  "demo",
					"demo2": "demo2",
				}).
				SetDirectory("~"),
		},
		{
			name: "script with with custom shell",
			config: rules.ShellScript{
				BaseRule: baseRule,
				BaseShell: rules.BaseShell{
					Name:   "zsh-script",
					Output: true,
					Shell:  "zsh",
				},
			},
			expectedOutput: "demo",
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

					return tt.execErr
				})

			ctx.ShellMock.Return(sh)

			err := tt.config.Check(ctx, output)

			testutils.CheckError(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedOutput, output.String())
		})
	}
}

func TestShellScript_GetPosition(t *testing.T) {
	rule := rules.ShellScript{BaseRule: rules.BaseRule{Type: rules.ShellScriptType}}

	actual := rule.GetPosition()

	assert.Equal(t, actual, rules.Scripts)
}

func TestShellScript_Compile(t *testing.T) {
	rule := rules.ShellScript{
		BaseRule: rules.BaseRule{Type: rules.ShellScriptType},
		BaseShell: rules.BaseShell{
			Name:     "{{var1}}",
			Shell:    "{{var1}}",
			Commands: []string{"{{var1}}1", "{{var1}}2"},
			Env: map[string]string{
				"{{var1}}": "{{var1}}",
			},
			Dir:    "{{var1}}",
			Output: true,
		},
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, rules.ShellScript{
		BaseRule: rules.BaseRule{Type: rules.ShellScriptType},
		BaseShell: rules.BaseShell{
			Name:     "VALUE",
			Shell:    "{{var1}}",
			Commands: []string{"VALUE1", "VALUE2"},
			Env: map[string]string{
				"{{var1}}": "VALUE",
			},
			Dir:    "VALUE",
			Output: true,
		},
	}, rule)
}

func TestShellScript_GetPrefix(t *testing.T) {
	expectedValue := "TestName"
	rule := rules.ShellScript{
		BaseRule:  rules.BaseRule{Type: rules.ShellScriptType},
		BaseShell: rules.BaseShell{Name: expectedValue},
	}

	actual := rule.GetPrefix()

	assert.Equal(t, actual, expectedValue)
}

func TestShellScript_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expected    *rules.ShellScript
		expectedErr string
	}{
		{
			name: "crossplatform script",
			config: `
type: shell-script
when: 1=1
name: TestName
`,
			expected: &rules.ShellScript{
				BaseRule: rules.BaseRule{
					Type:      rules.ShellScriptType,
					Condition: "1=1",
				},
				BaseShell: rules.BaseShell{
					Name: "TestName",
				},
			},
		},
		{
			name: "platform related script",
			config: `
type: shell-script
when: 1=1
windows:
  name: windows
  shell: bash
  commands: [ 'echo test' ]
  dir: test
linux:
  name: linux
  shell: bash
  commands: [ 'echo test' ]
  dir: test
macos:
  name: darwin
  shell: bash
  commands: [ 'echo test' ]
  dir: test
`,
			expected: &rules.ShellScript{
				BaseRule: rules.BaseRule{
					Type:      rules.ShellScriptType,
					Condition: "1=1",
				},
				BaseShell: rules.BaseShell{
					Name:     runtime.GOOS,
					Shell:    "bash",
					Commands: []string{"echo test"},
					Dir:      "test",
				},
			},
		},
		{
			name: "incorrect yaml",
			config: `
type: shell-script
when: '
`,
			expected:    nil,
			expectedErr: "yaml: line 3: found unexpected end of stream",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config rules.ShellScript

			err := testutils.DecodeYaml(tt.config, &config)

			testutils.CheckError(t, tt.expectedErr, err)
			if len(tt.expectedErr) > 0 {
				assert.Equal(t, config, rules.ShellScript{})
			} else {
				assert.Equal(t, *tt.expected, config)
			}
		})
	}
}
