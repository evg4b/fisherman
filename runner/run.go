package runner

import (
	"fisherman/commands"
	"fisherman/constants"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"flag"
	"fmt"
	"strings"
)

// Run executes application
func (r *Runner) Run(args []string) error {
	if len(args) < 1 {
		logger.Debug("No command detected")
		utils.PrintGraphics(logger.Writer(), constants.Logo, constants.Version)
		flag.Parse()
		flag.PrintDefaults()
		return nil
	}

	commandName := args[0]
	logger.Debugf("Runned program from binary '%s'", r.app.Executable)
	logger.Debugf("Called command '%s'", commandName)

	for _, command := range r.commandList {
		if strings.EqualFold(command.Name(), commandName) {
			ctx := commands.NewContext(commands.CliCommandContextParams{
				FileAccessor: r.fileAccessor,
				Usr:          r.systemUser,
				App:          r.app,
				Config:       r.config,
			})
			logger.Debugf("Context for command '%s' was created", commandName)

			err := command.Init(args[1:])
			utils.HandleCriticalError(err)
			logger.Debugf("Command '%s' was initialized", commandName)

			if commandError := command.Run(ctx); commandError != nil {
				logger.Debugf("Command '%s' finished with error", commandName)
				return commandError
			}

			logger.Debugf("Command '%s' finished witout error", commandName)
			return nil
		}
	}

	logger.Debugf("Command %s not found", commandName)
	return fmt.Errorf("Unknown command: %s", commandName)
}
