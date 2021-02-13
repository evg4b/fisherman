package hookfactory

import (
	"fisherman/configuration"
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint: dupl
func Test_getPreScripts(t *testing.T) {
	rule1 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)
	rule2 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PostScripts)
	rule3 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)
	rule4 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)

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
			actual := getPreScriptRules(tt.ruleCollection)

			assert.EqualValues(t, tt.expected, actual)
		})
	}
}

// nolint: dupl
func Test_getPostScriptRules(t *testing.T) {
	rule1 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PostScripts)
	rule2 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)
	rule3 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PostScripts)
	rule4 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)

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
	rule1 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)
	rule2 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PostScripts)
	rule3 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)
	rule4 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)

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
