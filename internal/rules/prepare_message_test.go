package rules_test

import (
	"context"
	"fisherman/internal/configuration"
	. "fisherman/internal/rules"
	"fisherman/testing/testutils"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMessage_Check(t *testing.T) {
	messageFilePath := "./hooks/MESSAGE"
	message := "custom message"

	fs := testutils.FsFromMap(t, map[string]string{
		messageFilePath: message,
	})

	t.Run("not configured rule", func(t *testing.T) {
		rule := makeRule(
			&PrepareMessage{},
			WithArgs([]string{messageFilePath}),
			WithFileSystem(fs),
		)

		err := rule.Check(context.TODO(), io.Discard)

		assert.NoError(t, err)
	})

	t.Run("succeeded check ", func(t *testing.T) {
		rule := makeRule(
			&PrepareMessage{Message: message},
			WithArgs([]string{messageFilePath}),
			WithFileSystem(fs),
		)

		err := rule.Check(context.TODO(), io.Discard)

		assert.NoError(t, err)
	})
}

func TestPrepareMessage_Compile(t *testing.T) {
	rule := PrepareMessage{
		Message: "{{var1}}",
	}

	rule.Compile(map[string]interface{}{"var1": "VALUE"})

	assert.Equal(t, PrepareMessage{
		Message: "VALUE",
	}, rule)
}

func makeRule(rule configuration.Rule, options ...RuleOption) configuration.Rule {
	rule.Configure(options...)

	return rule
}
