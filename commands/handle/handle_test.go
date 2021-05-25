package handle_test

import (
	"fisherman/commands/handle"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := handle.NewCommand(
		mocks.NewFactoryMock(t),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	assert.Equal(t, "handle", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := handle.NewCommand(
		mocks.NewFactoryMock(t),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	assert.NotEmpty(t, command.Description())
}
