package app_test

import (
	"errors"
	"fisherman/internal"
	"fisherman/internal/app"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliCommands_GetCommand(t *testing.T) {
	command1 := mocks.NewCliCommandMock(t).
		NameMock.Return("test").
		InitMock.Return(nil)
	command2 := mocks.NewCliCommandMock(t).
		NameMock.Return("demo").
		InitMock.When([]string{"arg1", "arg2"}).Then(nil)
	command3 := mocks.NewCliCommandMock(t).
		NameMock.Return("fail").
		InitMock.Return(errors.New("init failed"))

	tests := []struct {
		name        string
		commands    app.CliCommands
		args        []string
		expected    internal.CliCommand
		expectedErr string
	}{
		{
			name:     "Returns target command correctly",
			commands: app.CliCommands{command1, command2, command3},
			expected: command1,
			args:     []string{"test"},
		},
		{
			name:     "Returns target command correctly with arguments",
			commands: app.CliCommands{command1, command2, command3},
			expected: command2,
			args:     []string{"demo", "arg1", "arg2"},
		},
		{
			name:        "init returns error",
			commands:    app.CliCommands{command1, command2, command3},
			expected:    nil,
			args:        []string{"fail"},
			expectedErr: "init failed",
		},
		{
			name:        "unregistered command",
			commands:    app.CliCommands{command1, command2, command3},
			expected:    nil,
			args:        []string{"unregistered-command"},
			expectedErr: "unknown command: unregistered-command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.commands.GetCommand(tt.args)

			testutils.CheckError(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
