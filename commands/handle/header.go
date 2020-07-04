package handle

import (
	"fisherman/commands/context"
	"fisherman/constants"
)

func (c *Command) header(ctx context.Context, hook string) error {
	app, err := ctx.GetAppInfo()
	if err != nil {
		return err
	}
	c.reporter.PrintGraphics(constants.HookHeader, map[string]string{
		"Hook":             hook,
		"GlobalConfigPath": formatNA(app.GlobalConfigPath),
		"LocalConfigPath":  formatNA(app.LocalConfigPath),
		"RepoConfigPath":   formatNA(app.RepoConfigPath),
		"Version":          constants.Version,
	})
	return nil
}

func formatNA(path *string) string {
	if path == nil {
		return "N/A"
	}
	return *path
}
