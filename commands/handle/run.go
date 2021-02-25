package handle

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal/handling"
	"fisherman/utils"
)

func (command *Command) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *Command) Run() error {
	handler, err := command.hookFactory.GetHook(command.hook)
	if err != nil {
		if errors.Is(err, handling.ErrNotPresented) {
			log.Debugf("hook %s not presented", command.hook)

			return nil
		}

		return err
	}

	log.Debugf("handler for '%s' hook founded", command.hook)
	files := command.app.Configs
	utils.PrintGraphics(log.InfoOutput, constants.HookHeader, map[string]interface{}{
		constants.HookName:                 command.hook,
		constants.GlobalConfigPath:         utils.OriginalOrNA(files[configuration.GlobalMode]),
		constants.RepoConfigPath:           utils.OriginalOrNA(files[configuration.RepoMode]),
		constants.LocalConfigPath:          utils.OriginalOrNA(files[configuration.LocalMode]),
		constants.FishermanVersionVariable: constants.Version,
	})

	ctx := command.ctxFactory(command.flagSet.Args(), log.Stdout())

	return handler.Handle(ctx, command.flagSet.Args())
}
