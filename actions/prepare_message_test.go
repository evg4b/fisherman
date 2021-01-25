package actions_test

import (
	"fisherman/actions"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMessage_NotConfigured(t *testing.T) {
	next, err := actions.PrepareMessage(mocks.NewAsyncContextMock(t), "")

	assert.NoError(t, err)
	assert.True(t, next)
}

func TestPrepareMessage_CorrectWrite(t *testing.T) {
	messageFilePath := "./hooks/MESSAGE"
	message := "custom message"

	fs := mocks.NewFileSystemMock(t).
		WriteMock.Expect(messageFilePath, message).Return(nil)

	ctx := mocks.NewAsyncContextMock(t).
		ArgsMock.Return([]string{messageFilePath}).
		FilesMock.Return(fs)

	next, err := actions.PrepareMessage(ctx, message)

	assert.NoError(t, err)
	assert.False(t, next)
}
