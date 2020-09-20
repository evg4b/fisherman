package runner

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/utils"
	"flag"
	"fmt"
	"strings"
)

// Run executes application
func (runner *Runner) Run(conf *config.FishermanConfig, args []string) error {
	if len(args) < 2 {
		runner.logger.Debug("No command detected.")
		utils.PrintGraphics(runner.logger, constants.Logo, constants.Version)
		flag.Parse()
		flag.PrintDefaults()
		return nil
	}

	appPath := args[0]
	commandName := args[1]
	runner.logger.Debugf("Runned program from binary '%s'", appPath)
	runner.logger.Debugf("Called command '%s'", commandName)

	for _, command := range runner.commandList {
		if strings.EqualFold(command.Name(), commandName) {
			ctx := context.NewContext(context.CliCommandContextParams{
				FileAccessor: runner.fileAccessor,
				Usr:          runner.systemUser,
				Cwd:          runner.cwd,
				AppPath:      appPath,
				ConfigInfo:   runner.configInfo,
				Logger:       runner.logger,
			})
			runner.logger.Debugf("Context for command '%s' was created", commandName)
			if commandError := command.Run(ctx, args[2:]); commandError != nil {
				runner.logger.Debugf("Command '%s' finished with error", commandName)
				return commandError
			}

			runner.logger.Debugf("Command '%s' finished witout error", commandName)
			return nil
		}
	}

	runner.logger.Debugf("Command %s not found", commandName)
	return fmt.Errorf("Unknown command: %s", commandName)
}
