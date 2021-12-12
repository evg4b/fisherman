package app_test

import (
	"fisherman/internal"
	. "fisherman/internal/app"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliCommands_GetCommand(t *testing.T) {
	command1 := mocks.NewCliCommandMock(t).
		NameMock.Return("test")

	command2 := mocks.NewCliCommandMock(t).
		NameMock.Return("demo")

	command3 := mocks.NewCliCommandMock(t).
		NameMock.Return("fail")

	tests := []struct {
		name        string
		commands    CliCommands
		commandName string
		expected    internal.CliCommand
		expectedErr string
	}{
		{
			name:        "returns target command correctly",
			commands:    CliCommands{command1, command2, command3},
			expected:    command1,
			commandName: "test",
		},
		{
			name:        "unregistered command",
			commands:    CliCommands{command1, command2, command3},
			expected:    nil,
			commandName: "unregistered-command",
			expectedErr: "unknown command: unregistered-command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.commands.GetCommand(tt.commandName)

			testutils.AssertError(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
