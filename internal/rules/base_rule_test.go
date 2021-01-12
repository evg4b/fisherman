package rules

import (
	"testing"

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

	assert.Equal(t, BeforeScripts, rule.GetPosition())
}
