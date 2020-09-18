package handle

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fisherman/utils"
	"fmt"
	"strings"
)

// Run executes handle command
func (c *Command) Run(ctx *context.CommandContext, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	if hookHandler, ok := c.handlers[strings.ToLower(c.hook)]; ok {
		utils.PrintGraphics(ctx.Logger, constants.HookHeader, map[string]string{
			"Hook":             c.hook,
			"GlobalConfigPath": utils.OriginalOrNA(ctx.AppInfo.GlobalConfigPath),
			"LocalConfigPath":  utils.OriginalOrNA(ctx.AppInfo.LocalConfigPath),
			"RepoConfigPath":   utils.OriginalOrNA(ctx.AppInfo.RepoConfigPath),
			"Version":          constants.Version,
		})
		return hookHandler(ctx, c.fs.Args())
	}

	return fmt.Errorf("'%s' is not valid hook name", c.hook)
}
