package rules_test

import (
	"fisherman/internal/rules"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunProgram_GetPosition(t *testing.T) {
	rule := rules.Exec{
		BaseRule: rules.BaseRule{Type: rules.ExecType},
	}

	actual := rule.GetPosition()

	assert.Equal(t, actual, rules.Scripts)
}

func TestRunProgram_GetPrefix(t *testing.T) {
	tests := []struct {
		name     string
		ruleName string
		expected string
	}{
		{
			name:     "user defined name",
			ruleName: "Prefix",
			expected: "Prefix",
		},
		{
			name:     "default prefix",
			expected: rules.ExecType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := rules.Exec{
				BaseRule: rules.BaseRule{Type: rules.ExecType},
				Name:     tt.ruleName,
			}

			actual := rule.GetPrefix()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
