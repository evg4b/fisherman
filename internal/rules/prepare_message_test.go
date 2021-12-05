package rules_test

import (
	. "fisherman/internal/rules"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMessage_Check(t *testing.T) {
	messageFilePath := "./hooks/MESSAGE"
	message := "custom message"

	fs := testutils.FsFromMap(t, map[string]string{
		messageFilePath: message,
	})

	ctx := mocks.NewExecutionContextMock(t).
		ArgsMock.Return([]string{messageFilePath}).
		FilesMock.Return(fs)

	t.Run("not configured rule", func(t *testing.T) {
		rule := PrepareMessage{}

		err := rule.Check(ctx, ioutil.Discard)

		assert.NoError(t, err)
	})

	t.Run("succeeded check ", func(t *testing.T) {
		rule := PrepareMessage{Message: message}

		err := rule.Check(ctx, ioutil.Discard)

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
