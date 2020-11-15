package handle

import (
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/utils"
	"fmt"
	"strings"
)

func (command *Command) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *Command) Run() error {
	if hookHandler, ok := command.handlers[strings.ToLower(command.hook)]; ok {
		if hookHandler.IsConfigured(command.config) {
			log.Debugf("handler for '%s' hook founded", command.hook)
			utils.PrintGraphics(log.InfoOutput, constants.HookHeader, map[string]interface{}{
				constants.HookName:                 command.hook,
				constants.GlobalConfigPath:         utils.OriginalOrNA(command.app.GlobalConfigPath),
				constants.LocalConfigPath:          utils.OriginalOrNA(command.app.LocalConfigPath),
				constants.RepoConfigPath:           utils.OriginalOrNA(command.app.RepoConfigPath),
				constants.FishermanVersionVariable: constants.Version,
			})

			return hookHandler.Handle(command.flagSet.Args())
		}

		log.Debugf("hook %s not presented", command.hook)

		return nil
	}

	return fmt.Errorf("'%s' is not valid hook name", command.hook)
}
