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
func (runner *Runner) Run(hooksConfig *config.LoadInfo, args []string) error {
	if len(args) < 2 {
		utils.PrintGraphics(runner.logger, constants.Logo, map[string]string{
			"Version": constants.Version,
		})
		flag.Parse()
		flag.PrintDefaults()
		return nil
	}

	return runner.runInternal(hooksConfig, args[1:], args[0])
}

func (runner *Runner) runInternal(conf *config.LoadInfo, args []string, appPath string) error {
	commandName := args[0]
	ctx, err := runner.createContext(conf, appPath)
	if err != nil {
		return err
	}

	for _, command := range runner.commandList {
		if strings.EqualFold(command.Name(), commandName) {
			return command.Run(ctx, args[1:])
		}
	}

	return fmt.Errorf("Unknown command: %s", commandName)
}
