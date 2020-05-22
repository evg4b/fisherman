package app_test

import (
	"github.com/evg4b/fisherman/internal"
	"github.com/evg4b/fisherman/testing/mocks"
	"github.com/evg4b/fisherman/testing/testutils"
	"testing"

	. "github.com/evg4b/fisherman/internal/app"

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
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := testCase.commands.GetCommand(testCase.commandName)

			testutils.AssertError(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
