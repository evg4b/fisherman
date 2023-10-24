package rules_test

import (
	"fmt"
	"testing"

	. "github.com/evg4b/fisherman/internal/rules"

	"github.com/stretchr/testify/assert"
)

func TestBaseRule_GetType(t *testing.T) {
	rule := BaseRule{Type: "demo-rule", Condition: "rule-condition"}

	assert.Equal(t, "demo-rule", rule.GetType())
}

func TestBaseRule_GetContition(t *testing.T) {
	rule := BaseRule{Type: "demo-rule", Condition: "rule-condition"}

	assert.Equal(t, "rule-condition", rule.GetContition())
}

func TestBaseRule_GetPosition(t *testing.T) {
	rule := BaseRule{Type: "demo-rule", Condition: "rule-condition"}

	assert.Equal(t, PreScripts, rule.GetPosition())
}

func TestBaseRule_Compile(t *testing.T) {
	rule := BaseRule{
		Condition: "{{var1}}",
		Type:      "{{var1}}",
	}

	rule.Compile(map[string]any{"var1": "VALUE"})

	assert.Equal(t, BaseRule{Condition: "VALUE", Type: "{{var1}}"}, rule)
}

func TestBaseRule_GetPrefix(t *testing.T) {
	expected := "test-type"
	rule := BaseRule{
		Type: expected,
	}

	actual := rule.GetPrefix()

	assert.Equal(t, expected, actual)
}

func errorMessage(typeString, message string) string {
	return fmt.Sprintf("[%s] %s", typeString, message)
}
