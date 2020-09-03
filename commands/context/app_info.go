package context

import "fisherman/infrastructure/path"

// AppInfo is application info structure
type AppInfo struct {
	AppPath            string
	IsRegisteredInPath bool
	GlobalConfigPath   string
	RepoConfigPath     string
	LocalConfigPath    string
}

// GetAppInfo returns application info structure
func (ctx *CommandContext) GetAppInfo() (*AppInfo, error) {
	isRegistered, err := path.IsRegisteredInPath(ctx.path, ctx.appPath)
	if err != nil {
		return nil, err
	}

	return &AppInfo{
		GlobalConfigPath:   ctx.globalConfigPath,
		LocalConfigPath:    ctx.localConfigPath,
		RepoConfigPath:     ctx.repoConfigPath,
		IsRegisteredInPath: isRegistered,
		AppPath:            ctx.appPath,
	}, nil
}
