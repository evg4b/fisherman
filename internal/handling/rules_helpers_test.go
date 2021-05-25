package handling_test

import (
	"fisherman/configuration"
	"fisherman/internal/constants"
	"fisherman/internal/handling"
	"fisherman/internal/rules"
	"fisherman/testing/mocks"
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
		ruleCollection []handling.Rule
		expected       []handling.Rule
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
			factory := handling.NewHookHandlerFactory(mocks.NewEngineMock(t), configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: tt.ruleCollection,
				},
			})

			actual, _ := factory.GetHook(constants.PreCommitHook, globalVars)

			handler := actual.(*handling.HookHandler)

			assert.EqualValues(t, tt.expected, handler.Rules)
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
		ruleCollection []handling.Rule
		expected       []handling.Rule
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
			factory := handling.NewHookHandlerFactory(mocks.NewEngineMock(t), configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: tt.ruleCollection,
				},
			})

			actual, _ := factory.GetHook(constants.PreCommitHook, globalVars)

			handler := actual.(*handling.HookHandler)

			assert.EqualValues(t, tt.expected, handler.PostScriptRules)
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
		ruleCollection []handling.Rule
		expected       []handling.Rule
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
			factory := handling.NewHookHandlerFactory(mocks.NewEngineMock(t), configuration.HooksConfig{
				PreCommitHook: &configuration.HookConfig{
					Rules: tt.ruleCollection,
				},
			})

			actual, _ := factory.GetHook(constants.PreCommitHook, globalVars)

			handler := actual.(*handling.HookHandler)

			assert.EqualValues(t, tt.expected, handler.Scripts)
		})
	}
}

func getRule(t *testing.T, ruleType byte) handling.Rule {
	return mocks.NewRuleMock(t).
		GetTypeMock.Return(rules.ShellScriptType).
		GetPositionMock.Return(ruleType).
		CompileMock.Return()
}
