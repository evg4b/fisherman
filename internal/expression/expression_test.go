package expression_test

import (
	"errors"
	. "fisherman/internal/expression"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionEngine_Eval(t *testing.T) {
	engine := NewGoExpressionEngine()

	tests := []struct {
		name        string
		expression  string
		expected    bool
		expectedErr string
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
			name:        "IsEmpty with not empty value",
			expression:  "IsEmpty()",
			expected:    false,
			expectedErr: "expected bool, but got interface {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := engine.Eval(tt.expression, map[string]interface{}{
				"X":          "this is x value",
				"EmptyValue": "",
			})

			assert.Equal(t, tt.expected, actual)
			testutils.AssertError(t, tt.expectedErr, err)
		})
	}

	t.Run("panic in function", func(t *testing.T) {
		_, err := engine.Eval("test()", map[string]interface{}{
			"test": func() bool { panic(errors.New("test error")) },
		})

		assert.EqualError(t, err, "test error (1:1)\n | test()\n | ^")
	})
}
