package app

import (
	"fisherman/internal"
	"strings"

	"github.com/go-errors/errors"
)

// CliCommands is commands collection type.
type CliCommands []internal.CliCommand

// GetCommand returns command object by name.
func (commands CliCommands) GetCommand(commandName string) (internal.CliCommand, error) {
	for _, command := range commands {
		if strings.EqualFold(command.Name(), commandName) {
			return command, nil
		}
	}

	return nil, errors.Errorf("unknown command: %s", commandName)
}
