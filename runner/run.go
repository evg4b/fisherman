package runner

import (
	"fisherman/commands"
	"fisherman/constants"
	"fisherman/utils"
	"flag"
	"fmt"
	"strings"
)

// Run executes application
func (r *Runner) Run(args []string) error {
	if len(args) < 2 {
		r.logger.Debug("No command detected")
		utils.PrintGraphics(r.logger, constants.Logo, constants.Version)
		flag.Parse()
		flag.PrintDefaults()
		return nil
	}

	commandName := args[0]
	r.logger.Debugf("Runned program from binary '%s'", r.app.Executable)
	r.logger.Debugf("Called command '%s'", commandName)

	for _, command := range r.commandList {
		if strings.EqualFold(command.Name(), commandName) {
			ctx := commands.NewContext(commands.CliCommandContextParams{
				FileAccessor: r.fileAccessor,
				Usr:          r.systemUser,
				Logger:       r.logger,
				App:          r.app,
				Config:       r.config,
			})
			r.logger.Debugf("Context for command '%s' was created", commandName)

			err := command.Init(args[1:])
			utils.HandleCriticalError(err)
			r.logger.Debugf("Command '%s' was initialized", commandName)

			if commandError := command.Run(ctx); commandError != nil {
				r.logger.Debugf("Command '%s' finished with error", commandName)
				return commandError
			}

			r.logger.Debugf("Command '%s' finished witout error", commandName)
			return nil
		}
	}

	r.logger.Debugf("Command %s not found", commandName)
	return fmt.Errorf("Unknown command: %s", commandName)
}
