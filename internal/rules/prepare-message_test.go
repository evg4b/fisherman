package rules_test

import (
	"fisherman/internal/rules"
	"fisherman/mocks"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMessage_Check(t *testing.T) {
	messageFilePath := "./hooks/MESSAGE"
	message := "custom message"

	fs := mocks.NewFileSystemMock(t).
		WriteMock.Expect(messageFilePath, message).Return(nil)

	ctx := mocks.NewExecutionContextMock(t).
		ArgsMock.Return([]string{messageFilePath}).
		FilesMock.Return(fs)

	rule := rules.PrepareMessage{Message: message}
	err := rule.Check(ioutil.Discard, ctx)

	assert.NoError(t, err)
}

func TestPrepareMessage_Check_NotConfigured(t *testing.T) {
	ctx := mocks.NewExecutionContextMock(t)
	rule := rules.PrepareMessage{}

	err := rule.Check(ioutil.Discard, ctx)
	assert.NoError(t, err)
}
