package rules_test

import (
	"fisherman/internal/rules"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunProgram_GetPosition(t *testing.T) {
	rule := rules.RunProgram{
		BaseRule: rules.BaseRule{Type: rules.RunProgramType},
	}

	actual := rule.GetPosition()

	assert.Equal(t, actual, rules.Scripts)
}

func TestRunProgram_GetPrefix(t *testing.T) {
	tests := []struct {
		name     string
		ruleName string
		program  string
		args     []string
		expected string
	}{
		{
			name:     "user defined name",
			ruleName: "Prefix",
			program:  "go",
			args:     []string{"test", "./..."},
			expected: "Prefix",
		},
		{
			name:     "generated short prefix",
			program:  "go",
			args:     []string{"version"},
			expected: "go version",
		},
		{
			name:     "generated long prefix",
			program:  "program",
			args:     []string{"arg1", "arg2", "arg3", "arg4", "arg5", "arg6", "arg7"},
			expected: "program arg1 a...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := rules.RunProgram{
				BaseRule: rules.BaseRule{Type: rules.RunProgramType},
				Name:     tt.ruleName,
				Program:  tt.program,
				Args:     tt.args,
			}

			actual := rule.GetPrefix()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
