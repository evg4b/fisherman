package runner

import (
	"fisherman/constants"
	"flag"
	"fmt"
	"strings"
)

// Run executes application
func (runner *Runner) Run(args []string) error {
	if len(args) < 2 {
		runner.reporter.Info(constants.Logo)
		flag.Parse()
		flag.PrintDefaults()
		return nil
	}
	return runner.runInternal(args[1:], args[0])
}

func (runner *Runner) runInternal(args []string, appPath string) error {
	commandName := args[0]
	context, err := runner.buildContext(appPath)
	if err != nil {
		return err
	}

	for _, command := range runner.commandList {
		if strings.EqualFold(command.Name(), commandName) {
			return command.Run(context, args[1:])
		}
	}

	return fmt.Errorf("unknown command: %s", commandName)
}
