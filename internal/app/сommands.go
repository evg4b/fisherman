package app

import (
	"fisherman/internal"
	"strings"

	"github.com/go-errors/errors"
)

type CliCommands []internal.CliCommand

func (commands CliCommands) GetCommand(args []string) (internal.CliCommand, error) {
	commandName := args[0]

	for _, command := range commands {
		if strings.EqualFold(command.Name(), commandName) {
			err := command.Init(args[1:])
			if err != nil {
				return nil, err
			}

			return command, nil
		}
	}

	return nil, errors.Errorf("unknown command: %s", commandName)
}
