package runner

import (
	"fisherman/infrastructure/log"
	"fisherman/utils"
	"fmt"
	"strings"
)

func (r *Runner) Run(args []string) error {
	if len(args) < 1 {
		log.Debug("No command detected")
		r.PrintDefaults()

		return nil
	}

	commandName := args[0]
	log.Debugf("Runned program from binary '%s'", r.app.Executable)
	log.Debugf("Runned runned in folder '%s'", r.app.Cwd)
	log.Debugf("Called command '%s'", commandName)

	for _, command := range r.commands {
		if strings.EqualFold(command.Name(), commandName) {
			err := command.Init(args[1:])
			utils.HandleCriticalError(err)
			log.Debugf("Command '%s' was initialized", commandName)

			if commandError := command.Run(); commandError != nil {
				log.Debugf("Command '%s' finished with error", commandName)

				return commandError
			}

			log.Debugf("Command '%s' finished witout error", commandName)

			return nil
		}
	}

	log.Debugf("Command %s not found", commandName)

	return fmt.Errorf("unknown command: %s", commandName)
}
