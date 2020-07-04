package handle

import (
	"fisherman/commands/context"
	"fisherman/constants"
)

// HookInfo is structure for storage about hook
type HookInfo struct {
	Hook             string
	GlobalConfigPath string
	RepoConfigPath   string
	LocalConfigPath  string
	Version          string
}

func (c *Command) header(ctx context.Context) error {
	app, err := ctx.GetAppInfo()
	if err != nil {
		return err
	}
	c.reporter.PrintGraphics(constants.HookHeader, HookInfo{
		Hook:             c.hook,
		GlobalConfigPath: formatNA(app.GlobalConfigPath),
		LocalConfigPath:  formatNA(app.LocalConfigPath),
		RepoConfigPath:   formatNA(app.RepoConfigPath),
		Version:          "0.0.1",
	})
	return nil
}

func formatNA(path *string) string {
	if path == nil {
		return "N/A"
	}
	return *path
}
