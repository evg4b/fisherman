package runner

import (
	"fisherman/config"
	"fisherman/constants"
	"fisherman/utils"
	"flag"
	"fmt"
	"strings"
)

// Run executes application
func (runner *Runner) Run(conf *config.LoadInfo, args []string) error {
	if len(args) < 2 {
		utils.PrintGraphics(runner.logger, constants.Logo, constants.Version)
		flag.Parse()
		flag.PrintDefaults()
		return nil
	}

	appPath := args[0]
	commandName := args[1]

	ctx, err := runner.createContext(conf, appPath)
	utils.HandleCriticalError(err)

	for _, command := range runner.commandList {
		if strings.EqualFold(command.Name(), commandName) {
			return command.Run(ctx, args[2:])
		}
	}

	return fmt.Errorf("Unknown command: %s", commandName)
}
