package handle

import (
	"fisherman/commands"
	"fisherman/constants"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"fmt"
	"strings"
)

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

// Run executes handle command
func (c *Command) Run(ctx *commands.CommandContext) error {
	if hookHandler, ok := c.handlers[strings.ToLower(c.hook)]; ok {
		logger.Debugf("Handler for '%s' hook founded", c.hook)
		utils.PrintGraphics(logger.Writer(), constants.HookHeader, map[string]interface{}{
			"Hook":             c.hook,
			"GlobalConfigPath": utils.OriginalOrNA(ctx.App.GlobalConfigPath),
			"LocalConfigPath":  utils.OriginalOrNA(ctx.App.LocalConfigPath),
			"RepoConfigPath":   utils.OriginalOrNA(ctx.App.Cwd),
			"Version":          constants.Version,
		})

		return hookHandler(ctx, c.fs.Args())
	}

	return fmt.Errorf("'%s' is not valid hook name", c.hook)
}
