package handle

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fisherman/utils"
)

func (c *Command) header(ctx *context.CommandContext, hook string) error {
	app, err := ctx.GetAppInfo()
	if err != nil {
		return err
	}

	utils.PrintGraphics(ctx.Logger, constants.HookHeader, map[string]string{
		"Hook":             hook,
		"GlobalConfigPath": formatNA(app.GlobalConfigPath),
		"LocalConfigPath":  formatNA(app.LocalConfigPath),
		"RepoConfigPath":   formatNA(app.RepoConfigPath),
		"Version":          constants.Version,
	})
	return nil
}

func formatNA(path string) string {
	if utils.IsEmpty(path) {
		return "N/A"
	}

	return path
}
