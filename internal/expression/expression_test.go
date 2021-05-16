package expression_test

import (
	"fisherman/internal/expression"
	"fisherman/testing/testutils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionEngine_Eval(t *testing.T) {
	engine := expression.NewGoExpressionEngine()

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
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

func TestGovaluateEngine_EvalMap(t *testing.T) {
	engine := expression.NewGoExpressionEngine()

	tests := []struct {
		name        string
		expression  string
		expected    map[string]interface{}
		expectedErr []string
	}{
		{
			name:       "static expression",
			expression: "Extract(\"refs/heads/master\", \"refs/heads/(?P<CurrentBranch>.*)\")",
			expected: map[string]interface{}{
				"CurrentBranch": "master",
			},
		},
		{
			name:       "incurrect expression",
			expression: "Extract(\"refs/heads/master\", \"refs/heads/(?P<CurrentBranch>.*))",
			expected:   nil,
			expectedErr: []string{
				"literal not terminated (1:64)",
				" | Extract(\"refs/heads/master\", \"refs/heads/(?P<CurrentBranch>.*))",
				" | ...............................................................^",
			},
		},
		{
			name:       "eval error",
			expression: "Extract(\"refs/heads/master\")",
			expected:   nil,
			expectedErr: []string{
				"incorrect arguments for Extract (1:1)",
				" | Extract(\"refs/heads/master\")",
				" | ^",
			},
		},
		{
			name:       "eval matching error",
			expression: "Extract(\"demo\", \"refs/heads/(?P<CurrentBranch>.*)\")",
			expected:   nil,
			expectedErr: []string{
				"filed match 'demo' to expression 'refs/heads/(?P<CurrentBranch>.*)' (1:1)",
				" | Extract(\"demo\", \"refs/heads/(?P<CurrentBranch>.*)\")",
				" | ^",
			},
		},
		{
			name:       "eval regexp error",
			expression: "Extract(\"demo\", \"refs/heads/(?P<CurrentBranch>\")",
			expected:   nil,
			expectedErr: []string{
				"error parsing regexp: missing closing ): `refs/heads/(?P<CurrentBranch>` (1:1)",
				" | Extract(\"demo\", \"refs/heads/(?P<CurrentBranch>\")",
				" | ^",
			},
		},
		{
			name:       "static expression",
			expression: "Extract(variable1, \"refs/heads/(?P<CurrentBranch>.*)\")",
			expected:   map[string]interface{}{"CurrentBranch": "master"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := engine.EvalMap(tt.expression, map[string]interface{}{
				"X":          "this is x value",
				"EmptyValue": "",
				"variable1":  "refs/heads/master",
				"variable2":  "refs/heads/master",
			})

			assert.Equal(t, tt.expected, actual)
			testutils.CheckError(t, strings.Join(tt.expectedErr, "\n"), err)
		})
	}
}
