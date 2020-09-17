package handle

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fisherman/utils"
)

func header(ctx *context.CommandContext, hook string) {
	appInfo, err := ctx.GetAppInfo()
	if err != nil {
		panic(err)
	}

	utils.PrintGraphics(ctx.Logger, constants.HookHeader, map[string]string{
		"Hook":             hook,
		"GlobalConfigPath": formatNA(appInfo.GlobalConfigPath),
		"LocalConfigPath":  formatNA(appInfo.LocalConfigPath),
		"RepoConfigPath":   formatNA(appInfo.RepoConfigPath),
		"Version":          constants.Version,
	})
}

func formatNA(path string) string {
	if utils.IsEmpty(path) {
		return "N/A"
	}

	return path
}
