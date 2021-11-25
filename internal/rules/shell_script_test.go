package rules_test

import (
	"fisherman/internal/rules"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellScript_GetPosition(t *testing.T) {
	rule := rules.ShellScript{
		BaseRule: rules.BaseRule{Type: rules.ShellScriptType},
	}

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
