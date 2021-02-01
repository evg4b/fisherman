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
	rule2 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.AfterScripts)
	rule3 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)
	rule4 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)

	tests := []struct {
		name           string
		ruleCollection []Rule
		want           []Rule
	}{
		{
			name:           "common filtering",
			ruleCollection: []configuration.Rule{rule1, rule2, rule3, rule4},
			want:           []configuration.Rule{rule1, rule3},
		},
		{
			name:           "empty collection",
			ruleCollection: []configuration.Rule{},
			want:           []configuration.Rule{},
		},
		{
			name:           "collection without target rules",
			ruleCollection: []configuration.Rule{rule2, rule4, rule2, rule4},
			want:           []configuration.Rule{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPreScripts(tt.ruleCollection)

			assert.EqualValues(t, tt.want, got)
		})
	}
}

// nolint: dupl
func Test_getPostScriptRules(t *testing.T) {
	rule1 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.AfterScripts)
	rule2 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)
	rule3 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.AfterScripts)
	rule4 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)

	tests := []struct {
		name           string
		ruleCollection []Rule
		want           []Rule
	}{
		{
			name:           "common filtering",
			ruleCollection: []configuration.Rule{rule1, rule2, rule3, rule4},
			want:           []configuration.Rule{rule1, rule3},
		},
		{
			name:           "empty collection",
			ruleCollection: []configuration.Rule{},
			want:           []configuration.Rule{},
		},
		{
			name:           "collection without target rules",
			ruleCollection: []configuration.Rule{rule2, rule4, rule2, rule4},
			want:           []configuration.Rule{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPostScriptRules(tt.ruleCollection)

			assert.EqualValues(t, tt.want, got)
		})
	}
}

// nolint: dupl
func Test_getScriptRules(t *testing.T) {
	rule1 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)
	rule2 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.AfterScripts)
	rule3 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.Scripts)
	rule4 := mocks.NewRuleMock(t).GetPositionMock.Return(rules.PreScripts)

	tests := []struct {
		name           string
		ruleCollection []Rule
		want           []Rule
	}{
		{
			name:           "common filtering",
			ruleCollection: []configuration.Rule{rule1, rule2, rule3, rule4},
			want:           []configuration.Rule{rule1, rule3},
		},
		{
			name:           "empty collection",
			ruleCollection: []configuration.Rule{},
			want:           []configuration.Rule{},
		},
		{
			name:           "collection without target rules",
			ruleCollection: []configuration.Rule{rule2, rule4, rule2, rule4},
			want:           []configuration.Rule{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getScriptRules(tt.ruleCollection)

			assert.EqualValues(t, tt.want, got)
		})
	}
}
