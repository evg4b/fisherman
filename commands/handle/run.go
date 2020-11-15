package handle

import (
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/utils"
)

func (command *Command) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *Command) Run() error {
	handler, err := command.handlers.Get(command.hook)
	if err != nil {
		return err
	}

	if handler == nil {
		log.Debugf("hook %s not presented", command.hook)

		return nil
	}

	log.Debugf("handler for '%s' hook founded", command.hook)
	utils.PrintGraphics(log.InfoOutput, constants.HookHeader, map[string]interface{}{
		constants.HookName:                 command.hook,
		constants.GlobalConfigPath:         utils.OriginalOrNA(command.app.GlobalConfigPath),
		constants.LocalConfigPath:          utils.OriginalOrNA(command.app.LocalConfigPath),
		constants.RepoConfigPath:           utils.OriginalOrNA(command.app.RepoConfigPath),
		constants.FishermanVersionVariable: constants.Version,
	})

	return handler.Handle(command.flagSet.Args())
}
