package rules_test

import (
	"context"
	"fisherman/internal/configuration"
	"fisherman/testing/testutils"
	"io"
	"testing"

	. "fisherman/internal/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		require.NoError(t, err)
	})

	t.Run("succeeded check ", func(t *testing.T) {
		rule := makeRule(
			&PrepareMessage{Message: message},
			WithArgs([]string{messageFilePath}),
			WithFileSystem(fs),
		)

		err := rule.Check(context.TODO(), io.Discard)

		require.NoError(t, err)
	})
}

func TestPrepareMessage_Compile(t *testing.T) {
	rule := PrepareMessage{
		Message: "{{var1}}",
	}

	rule.Compile(map[string]any{"var1": "VALUE"})

	assert.Equal(t, PrepareMessage{
		Message: "VALUE",
	}, rule)
}

func makeRule(rule configuration.Rule, options ...RuleOption) configuration.Rule {
	rule.Configure(options...)

	return rule
}
