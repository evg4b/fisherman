package rules_test

import (
	"context"
	"github.com/evg4b/fisherman/testing/testutils"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/evg4b/fisherman/internal/rules"

	"github.com/stretchr/testify/assert"
)

func TestShellScript_GetPosition(t *testing.T) {
	rule := ShellScript{
		BaseRule: BaseRule{Type: ShellScriptType},
	}

	actual := rule.GetPosition()

	assert.Equal(t, actual, Scripts)
}

func TestShellScript_Compile(t *testing.T) {
	rule := ShellScript{
		BaseRule: BaseRule{Type: ShellScriptType},
		Name:     "{{var1}}",
		Shell:    "{{var1}}",
		Commands: []string{"{{var1}}1", "{{var1}}2"},
		Env: map[string]string{
			"{{var1}}": "{{var1}}",
		},
		Dir:    "{{var1}}",
		Output: true,
	}

	rule.Compile(map[string]any{"var1": "VALUE"})

	assert.Equal(t, ShellScript{
		BaseRule: BaseRule{Type: ShellScriptType},
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

func TestShellScript_GetPrefix(t *testing.T) {
	expectedValue := "TestName"
	rule := ShellScript{
		BaseRule: BaseRule{Type: ShellScriptType},
		Name:     expectedValue,
	}

	actual := rule.GetPrefix()

	assert.Equal(t, expectedValue, actual)
}

func TestShellScript_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expected    *ShellScript
		expectedErr string
	}{
		{
			name: "crossplatform script",
			config: `
type: shell-script
when: 1=1
name: TestName
`,
			expected: &ShellScript{
				BaseRule: BaseRule{
					Type:      ShellScriptType,
					Condition: "1=1",
				},
				Name: "TestName",
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
			var config ShellScript

			err := testutils.DecodeYaml(tt.config, &config)

			testutils.AssertError(t, tt.expectedErr, err)
			if len(tt.expectedErr) > 0 {
				assert.Equal(t, ShellScript{}, config)
			} else {
				assert.Equal(t, *tt.expected, config)
			}
		})
	}
}

func TestShellScript_Check(t *testing.T) {
	t.Run("return error for unknown shell", func(t *testing.T) {
		rule := ShellScript{
			BaseRule: BaseRule{Type: ShellScriptType},
			Shell:    "unknown-shell",
		}

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "failed to cheate shell host: unsupported shell")
	})

	t.Run("return error for incorrect encoding", func(t *testing.T) {
		rule := ShellScript{
			BaseRule: BaseRule{Type: ShellScriptType},
			Shell:    "bash",
			Encoding: "incorrect-encoding",
		}

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "failed to cheate shell host: 'incorrect-encoding' is unknown encoding")
	})

	t.Run("executed successful", func(t *testing.T) {
		rule := ShellScript{
			BaseRule: BaseRule{Type: ShellScriptType},
			Commands: []string{
				"echo test",
			},
		}

		err := rule.Check(context.TODO(), io.Discard)

		require.NoError(t, err)
	})

	t.Run("executed with non zero exit code ", func(t *testing.T) {
		rule := ShellScript{
			BaseRule: BaseRule{Type: ShellScriptType},
			Commands: []string{
				"exit 33",
			},
		}

		err := rule.Check(context.TODO(), io.Discard)

		require.EqualError(t, err, "[shell-script] script finished with exit code 33")
	})
}
