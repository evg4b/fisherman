package rules_test

import (
	"fisherman/internal/rules"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseRule_GetType(t *testing.T) {
	rule := rules.BaseRule{Type: "demo-rule", Condition: "rule-condition"}

	assert.Equal(t, "demo-rule", rule.GetType())
}

func TestBaseRule_GetContition(t *testing.T) {
	rule := rules.BaseRule{Type: "demo-rule", Condition: "rule-condition"}

	assert.Equal(t, "rule-condition", rule.GetContition())
}

func TestBaseRule_GetPosition(t *testing.T) {
	rule := rules.BaseRule{Type: "demo-rule", Condition: "rule-condition"}

	assert.Equal(t, rules.PreScripts, rule.GetPosition())
}

func TestBaseRule_Compile(t *testing.T) {
	rule := rules.BaseRule{
		Condition: "{{var1}}",
		Type:      "{{var1}}",
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, rules.BaseRule{Condition: "VALUE", Type: "{{var1}}"}, rule)
}

func TestBaseRule_GetPrefix(t *testing.T) {
	expected := "test-type"
	rule := rules.BaseRule{
		Type: expected,
	}

	actual := rule.GetPrefix()

	assert.Equal(t, expected, actual)
}
