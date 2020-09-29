package handle

import (
	"fisherman/commands"
	"fisherman/constants"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"fmt"
	"log"
	"strings"
)

// Init initialize handle command
func (c *Command) Init(args []string) error {
	err := c.fs.Parse(args)
	log.Println(c.fs.Args())

	return err
}

// Run executes handle command
func (c *Command) Run(ctx *commands.CommandContext) error {
	if hookHandler, ok := c.handlers[strings.ToLower(c.hook)]; ok {
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
