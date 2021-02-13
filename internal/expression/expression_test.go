package expression_test

import (
	"fisherman/internal/expression"
	"fisherman/testing/testutils"
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
			expectedErr: "incorrect arguments for IsEmpty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := engine.Eval(tt.expression)

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

func TestGovaluateEngine_EvalMap(t *testing.T) {
	engine := expression.NewExpressionEngine(map[string]interface{}{
		"X":          "this is x value",
		"EmptyValue": "",
		"variable1":  "refs/heads/master",
		"variable2":  "this is overwrites values",
	})

	tests := []struct {
		name        string
		expression  string
		expected    map[string]interface{}
		expectedErr string
	}{
		{
			name:       "static expression",
			expression: "Extract(\"refs/heads/master\", \"refs/heads/(?P<CurrentBranch>.*)\")",
			expected: map[string]interface{}{
				"CurrentBranch": "master",
			},
		},
		{
			name:        "incurrect expression",
			expression:  "Extract(\"refs/heads/master\", \"refs/heads/(?P<CurrentBranch>.*))",
			expected:    nil,
			expectedErr: "Unclosed string literal",
		},
		{
			name:        "eval error",
			expression:  "Extract(\"refs/heads/master\")",
			expected:    nil,
			expectedErr: "incorrect arguments for Extract",
		},
		{
			name:        "eval matching error",
			expression:  "Extract(\"demo\", \"refs/heads/(?P<CurrentBranch>.*)\")",
			expected:    nil,
			expectedErr: "filed match 'demo' to expression 'refs/heads/(?P<CurrentBranch>.*)'",
		},
		{
			name:        "eval regexp error",
			expression:  "Extract(\"demo\", \"refs/heads/(?P<CurrentBranch>\")",
			expected:    nil,
			expectedErr: "error parsing regexp: missing closing ): `refs/heads/(?P<CurrentBranch>`",
		},
		{
			name:       "static expression",
			expression: "Extract(variable1, \"refs/heads/(?P<CurrentBranch>.*)\")",
			expected: map[string]interface{}{
				"CurrentBranch": "master",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := engine.EvalMap(tt.expression, map[string]interface{}{
				"variable2": "refs/heads/master",
			})

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}
