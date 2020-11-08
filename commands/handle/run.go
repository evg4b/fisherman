package handle

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/utils"
	"fmt"
	"strings"
)

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

// Run executes handle command
func (c *Command) Run(ctx *clicontext.CommandContext) error {
	if hookHandler, ok := c.handlers[strings.ToLower(c.hook)]; ok {
		log.Debugf("Handler for '%s' hook founded", c.hook)
		utils.PrintGraphics(log.Writer(), constants.HookHeader, map[string]interface{}{
			constants.HookName:                 c.hook,
			constants.GlobalConfigPath:         utils.OriginalOrNA(ctx.App.GlobalConfigPath),
			constants.LocalConfigPath:          utils.OriginalOrNA(ctx.App.LocalConfigPath),
			constants.RepoConfigPath:           utils.OriginalOrNA(ctx.App.RepoConfigPath),
			constants.FishermanVersionVariable: constants.Version,
		})

		return hookHandler(ctx, c.flagSet.Args())
	}

	return fmt.Errorf("'%s' is not valid hook name", c.hook)
}
