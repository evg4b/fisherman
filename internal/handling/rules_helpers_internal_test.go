package handling

import (
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/rules"
	"github.com/evg4b/fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint: dupl
func Test_getPreScripts(t *testing.T) {
	rule1 := getRule(t, rules.PreScripts)
	rule2 := getRule(t, rules.PostScripts)
	rule3 := getRule(t, rules.PreScripts)
	rule4 := getRule(t, rules.Scripts)

	tests := []struct {
		name           string
		ruleCollection []Rule
		expected       []Rule
	}{
		{
			name:           "common filtering",
			ruleCollection: []configuration.Rule{rule1, rule2, rule3, rule4},
			expected:       []configuration.Rule{rule1, rule3},
		},
		{
			name:           "empty collection",
			ruleCollection: []configuration.Rule{},
			expected:       []configuration.Rule{},
		},
		{
			name:           "collection without target rules",
			ruleCollection: []configuration.Rule{rule2, rule4, rule2, rule4},
			expected:       []configuration.Rule{},
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			actual := getPreScriptRules(testCase.ruleCollection)

			assert.EqualValues(t, testCase.expected, actual)
		})
	}
}

// nolint: dupl
func Test_getPostScriptRules(t *testing.T) {
	rule1 := getRule(t, rules.PostScripts)
	rule2 := getRule(t, rules.Scripts)
	rule3 := getRule(t, rules.PostScripts)
	rule4 := getRule(t, rules.PreScripts)

	tests := []struct {
		name           string
		ruleCollection []Rule
		expected       []Rule
	}{
		{
			name:           "common filtering",
			ruleCollection: []configuration.Rule{rule1, rule2, rule3, rule4},
			expected:       []configuration.Rule{rule1, rule3},
		},
		{
			name:           "empty collection",
			ruleCollection: []configuration.Rule{},
			expected:       []configuration.Rule{},
		},
		{
			name:           "collection without target rules",
			ruleCollection: []configuration.Rule{rule2, rule4, rule2, rule4},
			expected:       []configuration.Rule{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getPostScriptRules(tt.ruleCollection)

			assert.EqualValues(t, tt.expected, actual)
		})
	}
}

// nolint: dupl
func Test_getScriptRules(t *testing.T) {
	rule1 := getRule(t, rules.Scripts)
	rule2 := getRule(t, rules.PostScripts)
	rule3 := getRule(t, rules.Scripts)
	rule4 := getRule(t, rules.PreScripts)

	tests := []struct {
		name           string
		ruleCollection []Rule
		expected       []Rule
	}{
		{
			name:           "common filtering",
			ruleCollection: []configuration.Rule{rule1, rule2, rule3, rule4},
			expected:       []configuration.Rule{rule1, rule3},
		},
		{
			name:           "empty collection",
			ruleCollection: []configuration.Rule{},
			expected:       []configuration.Rule{},
		},
		{
			name:           "collection without target rules",
			ruleCollection: []configuration.Rule{rule2, rule4, rule2, rule4},
			expected:       []configuration.Rule{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getScriptRules(tt.ruleCollection)

			assert.EqualValues(t, tt.expected, actual)
		})
	}
}

func getRule(t *testing.T, ruleType byte) Rule {
	t.Helper()

	return mocks.NewRuleMock(t).
		GetTypeMock.Return(rules.ShellScriptType).
		GetPositionMock.Return(ruleType).
		CompileMock.Return()
}
