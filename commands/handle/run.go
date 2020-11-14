package handle

import (
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal/clicontext"
	"fisherman/utils"
	"fmt"
	"strings"
)

func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *Command) Run(ctx *clicontext.CommandContext) error {
	if hookHandler, ok := c.handlers[strings.ToLower(c.hook)]; ok {
		if hookHandler.IsConfigured(ctx.Config) {
			log.Debugf("handler for '%s' hook founded", c.hook)
			utils.PrintGraphics(log.InfoOutput, constants.HookHeader, map[string]interface{}{
				constants.HookName:                 c.hook,
				constants.GlobalConfigPath:         utils.OriginalOrNA(ctx.App.GlobalConfigPath),
				constants.LocalConfigPath:          utils.OriginalOrNA(ctx.App.LocalConfigPath),
				constants.RepoConfigPath:           utils.OriginalOrNA(ctx.App.RepoConfigPath),
				constants.FishermanVersionVariable: constants.Version,
			})

			return hookHandler.Handle(c.flagSet.Args())
		}

		log.Debugf("hook %s not presented", c.hook)

		return nil
	}

	return fmt.Errorf("'%s' is not valid hook name", c.hook)
}
