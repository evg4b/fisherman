package expression_test

import (
	"errors"
	"fisherman/internal/expression"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionEngine_Eval(t *testing.T) {
	engine := expression.NewExpressionEngine(map[string]interface{}{
		"X":          "this is x value",
		"EmptyValue": "",
	})

	tests := []struct {
		name       string
		expression string
		expected   bool
		err        error
	}{
		{
			name:       "arithmetic expression",
			expression: "1 == 1",
			expected:   true,
		},
		{
			name:       "IsEmpty with empty value",
			expression: "IsEmpty(EmptyValue)",
			expected:   true,
		},
		{
			name:       "IsEmpty with not empty value",
			expression: "IsEmpty(X)",
			expected:   false,
		},
		{
			name:       "IsNotEmpty with not empty value",
			expression: "IsNotEmpty(X)",
			expected:   true,
		},
		{
			name:       "IsNotEmpty with empty value",
			expression: "IsNotEmpty(EmptyValue)",
			expected:   false,
		},
		{
			name:       "Envalid ex",
			expression: "IsNotEmpty(EmptyValue",
			expected:   false,
			err:        errors.New("Unbalanced parenthesis"),
		},
		{
			name:       "Envalid ex",
			expression: "IsNotEmpty(EmptyValue, X)",
			expected:   false,
			err:        errors.New("incorrect arguments for isNotEmpty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := engine.Eval(tt.expression)

			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.err, err)
		})
	}
}
