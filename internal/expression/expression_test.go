package expression_test

import (
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
		name        string
		expression  string
		expected    bool
		expectedErr error
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := engine.Eval(tt.expression)

			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
